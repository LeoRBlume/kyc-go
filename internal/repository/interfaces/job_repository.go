package interfaces

import "kyc-sim/internal/domain/models"

type JobRepository interface {
	FindActiveRunChecksJob(customerID string) (*models.Job, error)
	CreateJobWithItems(job *models.Job, items []models.JobItem) error
}
