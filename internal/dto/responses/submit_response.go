package responses

import "time"

type SubmitResponse struct {
	ID                string    `json:"id"`
	Status            string    `json:"status"`
	SubmittedAt       time.Time `json:"submittedAt"`
	RequiredDocuments []string  `json:"requiredDocuments"`
	MissingDocuments  []string  `json:"missingDocuments"`
}
