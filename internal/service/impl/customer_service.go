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
	repo         repoif.CustomerRepository
	documentRepo repoif.DocumentRepository
}

func NewCustomerService(
	repo repoif.CustomerRepository,
	documentRepo repoif.DocumentRepository,
) svcif.CustomerService {
	return &CustomerService{
		repo:         repo,
		documentRepo: documentRepo,
	}
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

func (s *CustomerService) Submit(id string) error {
	customer, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if customer.Status != domain.StatusDraft {
		return &domain.DomainError{
			Code:    domain.ErrConflict,
			Message: "customer cannot be submitted in current status",
			Details: map[string]any{
				"currentStatus": customer.Status,
				"allowedStatus": []string{string(domain.StatusDraft)},
			},
		}
	}

	// Documentos obrigatórios por tipo
	required := []domain.DocumentKind{}
	switch customer.Type {
	case domain.CustomerTypeIndividual:
		required = []domain.DocumentKind{
			domain.DocumentIDFront,
			domain.DocumentSelfie,
		}
	case domain.CustomerTypeBusiness:
		required = []domain.DocumentKind{
			domain.DocumentBusinessDoc,
		}
	}

	// Buscar documentos
	docs, err := s.documentRepo.FindByCustomer(customer.ID)
	if err != nil {
		return domain.NewInternal("failed to load documents", nil)
	}

	present := map[domain.DocumentKind]bool{}
	for _, d := range docs {
		present[d.Kind] = true
	}

	missing := []string{}
	for _, r := range required {
		if !present[r] {
			missing = append(missing, string(r))
		}
	}

	if len(missing) > 0 {
		return &domain.DomainError{
			Code:    domain.ErrValidation,
			Message: "missing required documents",
			Details: map[string]any{
				"required": toStringSlice(required),
				"missing":  missing,
			},
		}
	}

	customer.Status = domain.StatusSubmitted
	customer.UpdatedAt = time.Now()

	if err := s.repo.Update(customer); err != nil {
		return domain.NewInternal("failed to submit customer", nil)
	}

	return nil
}

func toStringSlice(kinds []domain.DocumentKind) []string {
	out := make([]string, 0, len(kinds))
	for _, k := range kinds {
		out = append(out, string(k))
	}
	return out
}
