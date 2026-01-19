package impl

import (
	"errors"

	"kyc-sim/internal/domain"
	repoif "kyc-sim/internal/repository/interfaces"
	"kyc-sim/internal/service/interfaces"

	"gorm.io/gorm"
)

type CustomerReadService struct {
	customerRepo repoif.CustomerRepository
	documentRepo repoif.DocumentRepository
	checkRepo    repoif.CheckRepository
	jobRepo      repoif.JobRepository
}

func NewCustomerReadService(
	customerRepo repoif.CustomerRepository,
	documentRepo repoif.DocumentRepository,
	checkRepo repoif.CheckRepository,
	jobRepo repoif.JobRepository,
) interfaces.CustomerReadService {
	return &CustomerReadService{
		customerRepo: customerRepo,
		documentRepo: documentRepo,
		checkRepo:    checkRepo,
		jobRepo:      jobRepo,
	}
}

func (s *CustomerReadService) GetCustomerView(customerID string) (*interfaces.CustomerView, error) {
	cust, err := s.customerRepo.FindByID(customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewNotFound("customer not found", map[string]any{"id": customerID})
		}
		return nil, domain.NewInternal("failed to load customer", map[string]any{"cause": err.Error()})
	}

	docs, err := s.documentRepo.FindByCustomer(customerID)
	if err != nil {
		return nil, domain.NewInternal("failed to load documents", map[string]any{"cause": err.Error()})
	}

	checks, err := s.checkRepo.FindByCustomer(customerID)
	if err != nil {
		return nil, domain.NewInternal("failed to load checks", map[string]any{"cause": err.Error()})
	}

	// Summary
	var summary interfaces.ChecksSummaryView
	checkViews := make([]interfaces.CheckView, 0, len(checks))
	for _, chk := range checks {
		checkViews = append(checkViews, interfaces.CheckView{
			CheckType: string(chk.Type),
			Status:    string(chk.Status),
			Score:     chk.Score,
			Details:   chk.Details,
			UpdatedAt: chk.UpdatedAt,
		})

		switch chk.Status {
		case domain.CheckPending:
			summary.Pending++
		case domain.CheckPass:
			summary.Pass++
		case domain.CheckFail:
			summary.Fail++
		case domain.CheckInconclusive:
			summary.Inconclusive++
		}
	}

	docViews := make([]interfaces.DocumentView, 0, len(docs))
	for _, d := range docs {
		docViews = append(docViews, interfaces.DocumentView{
			ID:         d.ID,
			CustomerID: d.CustomerID,
			Kind:       string(d.Kind),
			Status:     string(d.Status),
			UploadedAt: d.UploadedAt,
			FileURL:    d.FileURL,
		})
	}

	var lastJobView *interfaces.JobView
	active, err := s.jobRepo.FindActiveRunChecksJob(customerID)
	if err != nil {
		return nil, domain.NewInternal("failed to load job", map[string]any{"cause": err.Error()})
	}
	if active != nil {
		items, _ := s.jobRepo.ListItems(active.ID)
		var prog interfaces.JobProgressView
		prog.Total = len(items)
		for _, it := range items {
			if it.Status == "DONE" {
				prog.Done++
			}
			if it.Status == "FAILED" {
				prog.Failed++
			}
		}

		lastJobView = &interfaces.JobView{
			JobID:      active.ID,
			Kind:       string(active.Kind),
			Status:     string(active.Status),
			CustomerID: active.CustomerID,
			CreatedAt:  active.CreatedAt,
			UpdatedAt:  active.UpdatedAt,
			Progress:   prog,
		}
	}

	return &interfaces.CustomerView{
		ID:        cust.ID,
		Type:      string(cust.Type),
		Status:    string(cust.Status),
		CreatedAt: cust.CreatedAt,
		UpdatedAt: cust.UpdatedAt,

		Documents: docViews,
		Checks:    checkViews,
		Summary:   summary,
		LastJob:   lastJobView,
	}, nil
}
