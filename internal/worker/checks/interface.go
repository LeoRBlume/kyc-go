package checks

import "kyc-sim/internal/domain"

type Runner interface {
	Run() (domain.CheckStatus, *int, string)
}
