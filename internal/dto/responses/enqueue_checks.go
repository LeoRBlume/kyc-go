package responses

import "time"

type EnqueueChecksResponse struct {
	JobID      string    `json:"jobId"`
	CustomerID string    `json:"customerId"`
	Status     string    `json:"status"` // QUEUED | ALREADY_RUNNING
	QueuedAt   time.Time `json:"queuedAt"`
}
