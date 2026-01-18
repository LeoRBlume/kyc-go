package requests

import "kyc-sim/internal/domain"

type CreateDocumentRequest struct {
	Kind    domain.DocumentKind `json:"kind"`
	FileURL string              `json:"fileUrl"`
}
