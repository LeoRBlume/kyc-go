package interfaces

import "kyc-sim/internal/domain/models"

type CheckRepository interface {
	FindByCustomer(customerID string) ([]models.Check, error)
	Create(check *models.Check) error
	Update(check *models.Check) error
}
