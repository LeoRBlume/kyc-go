package domain

type CustomerType string
type CustomerStatus string
type DocumentKind string
type DocumentStatus string

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

const (
	DocumentIDFront      DocumentKind = "ID_FRONT"
	DocumentIDBack       DocumentKind = "ID_BACK"
	DocumentSelfie       DocumentKind = "SELFIE"
	DocumentProofAddress DocumentKind = "PROOF_OF_ADDRESS"
	DocumentBusinessDoc  DocumentKind = "BUSINESS_DOC"
	DocumentOther        DocumentKind = "OTHER"
)

const (
	DocumentUploaded  DocumentStatus = "UPLOADED"
	DocumentValidated DocumentStatus = "VALIDATED"
	DocumentRejected  DocumentStatus = "REJECTED"
)
