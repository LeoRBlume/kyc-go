package interfaces

import "kyc-sim/internal/domain/models"

type DocumentRepository interface {
	Create(doc *models.Document) error
	FindByCustomer(customerID string) ([]models.Document, error)
	ExistsByCustomerAndKind(customerID string, kind string) (bool, error)
}
