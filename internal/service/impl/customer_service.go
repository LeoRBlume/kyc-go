package impl

import (
	"errors"
	"kyc-sim/internal/domain"
	"kyc-sim/internal/domain/models"
	"kyc-sim/internal/dto/requests"
	repoif "kyc-sim/internal/repository/interfaces"
	"time"

	"gorm.io/gorm"

	svcif "kyc-sim/internal/service/interfaces"

	"github.com/google/uuid"
)

type CustomerService struct {
	repo repoif.CustomerRepository
}

func NewCustomerService(repo repoif.CustomerRepository) svcif.CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Create(req requests.CreateCustomerRequest) (*models.Customer, error) {
	if req.Type != domain.CustomerTypeIndividual && req.Type != domain.CustomerTypeBusiness {
		return nil, domain.NewValidation("invalid customer type", map[string]any{
			"type": req.Type,
		})
	}

	now := time.Now()
	c := &models.Customer{
		ID:        uuid.NewString(),
		Type:      req.Type,
		Status:    domain.StatusDraft,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(c); err != nil {
		return nil, domain.NewInternal("failed to create customer", map[string]any{
			"cause": err.Error(),
		})
	}
	return c, nil
}

func (s *CustomerService) GetByID(id string) (*models.Customer, error) {
	c, err := s.repo.FindByID(id)
	if err == nil {
		return c, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.NewNotFound("customer not found", map[string]any{
			"id": id,
		})
	}

	return nil, domain.NewInternal("failed to load customer", map[string]any{
		"cause": err.Error(),
	})
}

func (s *CustomerService) Patch(id string, req requests.PatchCustomerRequest) (*models.Customer, error) {
	customer, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if customer.Status != domain.StatusDraft {
		return nil, &domain.DomainError{
			Code:    "STATUS_CONFLICT",
			Message: "customer cannot be modified in current status",
			Details: map[string]any{
				"currentStatus": customer.Status,
				"allowedStatus": []string{string(domain.StatusDraft)},
			},
		}
	}

	if req.Type != nil {
		if *req.Type != domain.CustomerTypeIndividual && *req.Type != domain.CustomerTypeBusiness {
			return nil, domain.NewValidation("invalid customer type", map[string]any{
				"type": *req.Type,
			})
		}
		customer.Type = *req.Type
	}

	customer.UpdatedAt = time.Now()

	if err := s.repo.Update(customer); err != nil {
		return nil, domain.NewInternal("failed to update customer", map[string]any{
			"cause": err.Error(),
		})
	}

	return customer, nil
}
