package responses

type ListDocumentsResponse struct {
	CustomerID string             `json:"customerId"`
	Documents  []DocumentResponse `json:"documents"`
}
