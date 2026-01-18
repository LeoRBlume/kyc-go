package impl

import (
	"time"

	"kyc-sim/internal/domain"
	"kyc-sim/internal/domain/models"
	"kyc-sim/internal/dto/requests"
	repoif "kyc-sim/internal/repository/interfaces"
	svcif "kyc-sim/internal/service/interfaces"

	"github.com/google/uuid"
)

type DocumentService struct {
	customerRepo repoif.CustomerRepository
	documentRepo repoif.DocumentRepository
}

func NewDocumentService(
	customerRepo repoif.CustomerRepository,
	documentRepo repoif.DocumentRepository,
) svcif.DocumentService {
	return &DocumentService{
		customerRepo: customerRepo,
		documentRepo: documentRepo,
	}
}

func (s *DocumentService) Create(customerID string, req requests.CreateDocumentRequest) (*models.Document, error) {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil {
		return nil, err
	}

	if customer.Status != domain.StatusDraft {
		return nil, &domain.DomainError{
			Code:    domain.ErrConflict,
			Message: "documents can only be added while status is DRAFT",
			Details: map[string]any{
				"currentStatus": customer.Status,
			},
		}
	}

	if req.FileURL == "" {
		return nil, domain.NewValidation("fileUrl is required", nil)
	}

	exists, err := s.documentRepo.ExistsByCustomerAndKind(customerID, string(req.Kind))
	if err != nil {
		return nil, domain.NewInternal("failed to validate document uniqueness", nil)
	}
	if exists {
		return nil, &domain.DomainError{
			Code:    domain.ErrConflict,
			Message: "document kind already exists for customer",
			Details: map[string]any{
				"kind": req.Kind,
			},
		}
	}

	doc := &models.Document{
		ID:         uuid.NewString(),
		CustomerID: customerID,
		Kind:       req.Kind,
		FileURL:    req.FileURL,
		Status:     domain.DocumentUploaded,
		UploadedAt: time.Now(),
	}

	if err := s.documentRepo.Create(doc); err != nil {
		return nil, domain.NewInternal("failed to create document", nil)
	}

	return doc, nil
}

func (s *DocumentService) List(customerID string) ([]models.Document, error) {
	if _, err := s.customerRepo.FindByID(customerID); err != nil {
		return nil, err
	}
	return s.documentRepo.FindByCustomer(customerID)
}
