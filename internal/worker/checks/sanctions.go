package checks

import "kyc-sim/internal/domain"

type Sanctions struct{}

func NewSanctions() *Sanctions {
	return &Sanctions{}
}

func (s *Sanctions) Run() (domain.CheckStatus, *int, string) {
	score := 98
	return domain.CheckPass, &score, "no sanctions found"
}
