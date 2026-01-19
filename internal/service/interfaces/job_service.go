package interfaces

import "kyc-sim/internal/dto/responses"

type JobService interface {
	Get(jobID string) (*responses.GetJobResponse, error)
}
