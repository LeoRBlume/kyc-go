package domain

type CustomerType string
type CustomerStatus string

const (
	CustomerTypeIndividual CustomerType = "INDIVIDUAL"
	CustomerTypeBusiness   CustomerType = "BUSINESS"
)

const (
	StatusDraft     CustomerStatus = "DRAFT"
	StatusSubmitted CustomerStatus = "SUBMITTED"
	StatusInReview  CustomerStatus = "IN_REVIEW"
	StatusApproved  CustomerStatus = "APPROVED"
	StatusRejected  CustomerStatus = "REJECTED"
	StatusExpired   CustomerStatus = "EXPIRED"
)
