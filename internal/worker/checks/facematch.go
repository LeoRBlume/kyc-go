package checks

import "kyc-sim/internal/domain"

type FaceMatch struct{}

func NewFaceMatch() *FaceMatch {
	return &FaceMatch{}
}

func (f *FaceMatch) Run() (domain.CheckStatus, *int, string) {
	score := 85
	return domain.CheckPass, &score, "face match successful"
}
