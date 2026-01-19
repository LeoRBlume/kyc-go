package checks

import "kyc-sim/internal/domain"

type Pep struct{}

func NewPep() *Pep {
	return &Pep{}
}

func (p *Pep) Run() (domain.CheckStatus, *int, string) {
	score := 70
	return domain.CheckPass, &score, "not politically exposed"
}
