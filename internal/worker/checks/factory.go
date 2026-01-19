package checks

import "kyc-sim/internal/domain"

func Build(checkType domain.CheckType) Runner {
	switch checkType {
	case domain.CheckCpfCnpj:
		return NewCpfCnpj()
	case domain.CheckSanctions:
		return NewSanctions()
	case domain.CheckPep:
		return NewPep()
	case domain.CheckFaceMatch:
		return NewFaceMatch()
	default:
		return nil
	}
}
