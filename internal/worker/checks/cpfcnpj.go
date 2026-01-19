package checks

import "kyc-sim/internal/domain"

type CpfCnpj struct{}

func NewCpfCnpj() *CpfCnpj {
	return &CpfCnpj{}
}

func (c *CpfCnpj) Run() (domain.CheckStatus, *int, string) {
	score := 90
	return domain.CheckPass, &score, "cpf/cnpj valid and active"
}
