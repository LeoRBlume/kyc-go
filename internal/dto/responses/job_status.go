package responses

import "time"

type JobProgress struct {
	Total  int `json:"total"`
	Done   int `json:"done"`
	Failed int `json:"failed"`
}

type JobItemResponse struct {
	CheckType  string     `json:"checkType"`
	Status     string     `json:"status"` // PENDING/DONE/FAILED
	StartedAt  *time.Time `json:"startedAt,omitempty"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
	Error      *string    `json:"error,omitempty"`
}

type GetJobResponse struct {
	JobID       string            `json:"jobId"`
	Kind        string            `json:"kind"`
	Status      string            `json:"status"`
	CustomerID  string            `json:"customerId"`
	Attempt     int               `json:"attempt"`
	MaxAttempts int               `json:"maxAttempts"`
	NextRunAt   time.Time         `json:"nextRunAt"`
	LockedBy    *string           `json:"lockedBy,omitempty"`
	HeartbeatAt *time.Time        `json:"heartbeatAt,omitempty"`
	LastError   *string           `json:"lastError,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Progress    JobProgress       `json:"progress"`
	Items       []JobItemResponse `json:"items"`
}
