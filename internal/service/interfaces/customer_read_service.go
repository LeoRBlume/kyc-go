package interfaces

import (
	"time"
)

type CustomerView struct {
	ID        string
	Type      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time

	Documents []DocumentView
	Checks    []CheckView
	Summary   ChecksSummaryView

	LastJob *JobView
}

type DocumentView struct {
	ID         string
	CustomerID string
	Kind       string
	Status     string
	UploadedAt time.Time
	FileURL    string
}

type CheckView struct {
	CheckType string
	Status    string
	Score     *int
	Details   string
	UpdatedAt time.Time
}

type ChecksSummaryView struct {
	Pending      int
	Pass         int
	Fail         int
	Inconclusive int
}

type JobView struct {
	JobID      string
	Kind       string
	Status     string
	CustomerID string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Progress   JobProgressView
}

type JobProgressView struct {
	Total  int
	Done   int
	Failed int
}

type CustomerReadService interface {
	GetCustomerView(customerID string) (*CustomerView, error)
}
