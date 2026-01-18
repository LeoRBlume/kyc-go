package interfaces

import "kyc-sim/internal/domain/models"

type CustomerRepository interface {
	Create(customer *models.Customer) error
	FindByID(id string) (*models.Customer, error)
	Update(customer *models.Customer) error
}
