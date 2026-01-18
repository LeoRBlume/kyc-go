package interfaces

import (
	"kyc-sim/internal/domain/models"
	"kyc-sim/internal/dto/requests"
)

type DocumentService interface {
	Create(customerID string, req requests.CreateDocumentRequest) (*models.Document, error)
	List(customerID string) ([]models.Document, error)
}
