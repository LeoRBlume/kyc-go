package requests

import "kyc-sim/internal/domain"

type CreateCustomerRequest struct {
	Type domain.CustomerType `json:"type" binding:"required,oneof=INDIVIDUAL BUSINESS"`
}
