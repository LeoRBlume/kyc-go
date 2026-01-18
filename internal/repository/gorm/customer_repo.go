package gormrepo

import (
	"kyc-sim/internal/domain/models"
	repoif "kyc-sim/internal/repository/interfaces"

	"gorm.io/gorm"
)

type CustomerRepo struct {
	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) repoif.CustomerRepository {
	return &CustomerRepo{db: db}
}

func (r *CustomerRepo) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepo) FindByID(id string) (*models.Customer, error) {
	var c models.Customer
	err := r.db.First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
