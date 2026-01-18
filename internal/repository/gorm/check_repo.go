package gormrepo

import (
	"kyc-sim/internal/domain/models"
	repoif "kyc-sim/internal/repository/interfaces"

	"gorm.io/gorm"
)

type CheckRepo struct {
	db *gorm.DB
}

func NewCheckRepo(db *gorm.DB) repoif.CheckRepository {
	return &CheckRepo{db: db}
}

func (r *CheckRepo) FindByCustomer(customerID string) ([]models.Check, error) {
	var checks []models.Check
	err := r.db.Where("customer_id = ?", customerID).Find(&checks).Error
	return checks, err
}
