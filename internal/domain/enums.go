package domain

type CustomerType string
type CustomerStatus string
type DocumentKind string
type DocumentStatus string

type CheckType string
type CheckStatus string

type JobKind string
type JobStatus string

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

const (
	CheckCpfCnpj   CheckType = "CPF_CNPJ"
	CheckSanctions CheckType = "SANCTIONS"
	CheckPep       CheckType = "PEP"
	CheckFaceMatch CheckType = "FACE_MATCH"
)

const (
	CheckPending      CheckStatus = "PENDING"
	CheckPass         CheckStatus = "PASS"
	CheckFail         CheckStatus = "FAIL"
	CheckInconclusive CheckStatus = "INCONCLUSIVE"
)

const (
	JobRunChecks JobKind = "RUN_CHECKS"
)

const (
	JobPending  JobStatus = "PENDING"
	JobRunning  JobStatus = "RUNNING"
	JobDone     JobStatus = "DONE"
	JobFailed   JobStatus = "FAILED"
	JobCanceled JobStatus = "CANCELED"
)
