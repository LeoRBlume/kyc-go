package impl

import (
	"kyc-sim/internal/domain/models"
	repoif "kyc-sim/internal/repository/interfaces"
	svcif "kyc-sim/internal/service/interfaces"
)

type CheckService struct {
	customerRepo repoif.CustomerRepository
	checkRepo    repoif.CheckRepository
}

func NewCheckService(
	customerRepo repoif.CustomerRepository,
	checkRepo repoif.CheckRepository,
) svcif.CheckService {
	return &CheckService{
		customerRepo: customerRepo,
		checkRepo:    checkRepo,
	}
}

func (s *CheckService) ListByCustomer(customerID string) ([]models.Check, error) {
	// garante 404 se customer não existir
	if _, err := s.customerRepo.FindByID(customerID); err != nil {
		return nil, err
	}

	return s.checkRepo.FindByCustomer(customerID)
}
