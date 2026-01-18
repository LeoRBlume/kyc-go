package interfaces

import "kyc-sim/internal/domain/models"

type CheckService interface {
	ListByCustomer(customerID string) ([]models.Check, error)
}
