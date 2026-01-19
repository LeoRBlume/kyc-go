package interfaces

import (
	"kyc-sim/internal/domain/models"
	"kyc-sim/internal/dto/requests"
)

type KycService interface {
	EnqueueRunChecks(customerID string, req requests.RunChecksRequest) (*models.Job, bool, error)
}
