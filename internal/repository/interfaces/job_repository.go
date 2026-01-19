package interfaces

import (
	"kyc-sim/internal/domain/models"
	"time"
)

type JobRepository interface {
	FindActiveRunChecksJob(customerID string) (*models.Job, error)
	CreateJobWithItems(job *models.Job, items []models.JobItem) error
	ClaimNext(now time.Time, workerID string) (*models.Job, error)
	UpdateHeartbeat(jobID string, now time.Time) error
	MarkDone(jobID string) error
	MarkFailed(jobID string, errMsg string) error
	ListItems(jobID string) ([]models.JobItem, error)
	UpdateItem(item *models.JobItem) error
	GetByID(jobID string) (*models.Job, error)
}
