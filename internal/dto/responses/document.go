package responses

import "time"

type DocumentResponse struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customerId"`
	Kind       string    `json:"kind"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploadedAt"`
}
