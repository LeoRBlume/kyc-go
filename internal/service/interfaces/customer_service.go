package interfaces

import (
	"kyc-sim/internal/domain/models"
	"kyc-sim/internal/dto/requests"
)

type CustomerService interface {
	Create(req requests.CreateCustomerRequest) (*models.Customer, error)
	GetByID(id string) (*models.Customer, error)
}
