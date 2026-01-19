package impl

import (
	"time"

	"github.com/google/uuid"
	"kyc-sim/internal/domain"
	"kyc-sim/internal/domain/models"
	"kyc-sim/internal/dto/requests"
	repoif "kyc-sim/internal/repository/interfaces"
	svcif "kyc-sim/internal/service/interfaces"
)

type KycService struct {
	customerRepo repoif.CustomerRepository
	jobRepo      repoif.JobRepository
	checkRepo    repoif.CheckRepository
}

func NewKycService(
	customerRepo repoif.CustomerRepository,
	jobRepo repoif.JobRepository,
	checkRepo repoif.CheckRepository,
) svcif.KycService {
	return &KycService{
		customerRepo: customerRepo,
		jobRepo:      jobRepo,
		checkRepo:    checkRepo,
	}
}

func (s *KycService) EnqueueRunChecks(customerID string, req requests.RunChecksRequest) (*models.Job, bool, error) {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil {
		return nil, false, err
	}

	if customer.Status != domain.StatusSubmitted && customer.Status != domain.StatusInReview {
		return nil, false, &domain.DomainError{
			Code:    domain.ErrConflict,
			Message: "checks can only be run after submit",
			Details: map[string]any{
				"currentStatus": customer.Status,
			},
		}
	}

	if len(req.Checks) == 0 {
		return nil, false, domain.NewValidation("at least one check is required", nil)
	}

	// Idempotência
	if existing, err := s.jobRepo.FindActiveRunChecksJob(customerID); err != nil {
		return nil, false, err
	} else if existing != nil {
		return existing, true, nil
	}

	now := time.Now()
	job := &models.Job{
		ID:          "job_" + uuid.NewString(),
		CustomerID:  customerID,
		Kind:        domain.JobRunChecks,
		Status:      domain.JobPending,
		Attempt:     0,
		MaxAttempts: 5,
		NextRunAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	items := make([]models.JobItem, 0, len(req.Checks))
	for _, chk := range req.Checks {
		items = append(items, models.JobItem{
			ID:        uuid.NewString(),
			JobID:     job.ID,
			CheckType: chk,
			Status:    "PENDING",
		})

		// cria check inicial
		s.checkRepo.Create(&models.Check{
			ID:         uuid.NewString(),
			CustomerID: customerID,
			Type:       chk,
			Status:     domain.CheckPending,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
	}

	if err := s.jobRepo.CreateJobWithItems(job, items); err != nil {
		return nil, false, domain.NewInternal("failed to enqueue job", nil)
	}

	customer.Status = domain.StatusInReview
	customer.UpdatedAt = now
	_ = s.customerRepo.Update(customer)

	return job, false, nil
}
