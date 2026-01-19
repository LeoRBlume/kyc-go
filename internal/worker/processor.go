package worker

import (
	"time"

	"kyc-sim/internal/domain"
	"kyc-sim/internal/domain/models"
	repoif "kyc-sim/internal/repository/interfaces"
	"kyc-sim/internal/worker/checks"
)

type Processor struct {
	jobRepo      repoif.JobRepository
	checkRepo    repoif.CheckRepository
	customerRepo repoif.CustomerRepository
}

func NewProcessor(
	jobRepo repoif.JobRepository,
	checkRepo repoif.CheckRepository,
	customerRepo repoif.CustomerRepository,
) *Processor {
	return &Processor{
		jobRepo:      jobRepo,
		checkRepo:    checkRepo,
		customerRepo: customerRepo,
	}
}

func (p *Processor) Process(job *models.Job) error {
	// Carrega itens do job
	items, err := p.jobRepo.ListItems(job.ID)
	if err != nil {
		return err
	}

	// Carrega checks já existentes do customer
	checksList, err := p.checkRepo.FindByCustomer(job.CustomerID)
	if err != nil {
		return err
	}

	checkByType := make(map[domain.CheckType]*models.Check)
	for i := range checksList {
		chk := checksList[i]
		checkByType[chk.Type] = &chk
	}

	now := time.Now()

	// Flags para decisão
	allPass := true
	hasFail := false
	sanctionsFail := false

	for i := range items {
		item := &items[i]

		if item.Status != "PENDING" {
			continue
		}

		item.StartedAt = &now

		runner := checks.Build(item.CheckType)
		if runner == nil {
			item.Status = "FAILED"
			errMsg := "unsupported check type"
			item.Error = &errMsg
			item.FinishedAt = &now
			_ = p.jobRepo.UpdateItem(item)

			hasFail = true
			allPass = false
			continue
		}

		status, score, details := runner.Run()

		// Atualiza Check correspondente
		chk, ok := checkByType[item.CheckType]
		if !ok {
			// isso não deveria acontecer, mas não quebra o worker
			item.Status = "FAILED"
			errMsg := "check record not found"
			item.Error = &errMsg
			item.FinishedAt = &now
			_ = p.jobRepo.UpdateItem(item)

			hasFail = true
			allPass = false
			continue
		}

		chk.Status = status
		chk.Score = score
		chk.Details = details
		chk.UpdatedAt = now
		_ = p.checkRepo.Update(chk)

		// Atualiza flags de decisão
		switch status {
		case domain.CheckFail:
			allPass = false
			hasFail = true
			if chk.Type == domain.CheckSanctions {
				sanctionsFail = true
			}
		case domain.CheckInconclusive:
			allPass = false
		case domain.CheckPass:
			// ok
		}

		item.Status = "DONE"
		item.FinishedAt = &now
		_ = p.jobRepo.UpdateItem(item)
	}

	// Decisão automática
	customer, err := p.customerRepo.FindByID(job.CustomerID)
	if err != nil {
		return err
	}

	switch {
	case sanctionsFail:
		customer.Status = domain.StatusRejected
	case hasFail:
		customer.Status = domain.StatusInReview
	case allPass:
		customer.Status = domain.StatusApproved
	default:
		customer.Status = domain.StatusInReview
	}

	customer.UpdatedAt = time.Now()
	_ = p.customerRepo.Update(customer)

	// Finaliza job
	return p.jobRepo.MarkDone(job.ID)
}
