package requests

import "kyc-sim/internal/domain"

type RunChecksRequest struct {
	Checks []domain.CheckType `json:"checks"`
}
