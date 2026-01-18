package responses

import "time"

type CheckResponse struct {
	Type      string    `json:"checkType"`
	Status    string    `json:"status"`
	Score     *int      `json:"score,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ChecksSummary struct {
	Pending      int `json:"pending"`
	Pass         int `json:"pass"`
	Fail         int `json:"fail"`
	Inconclusive int `json:"inconclusive"`
}

type ListChecksResponse struct {
	CustomerID string          `json:"customerId"`
	Results    []CheckResponse `json:"results"`
	Summary    ChecksSummary   `json:"summary"`
}
