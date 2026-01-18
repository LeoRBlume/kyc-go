package requests

import "kyc-sim/internal/domain"

type PatchCustomerRequest struct {
	Type *domain.CustomerType `json:"type,omitempty"`
}
