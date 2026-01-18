package gormrepo

import (
	"kyc-sim/internal/domain/models"
	repoif "kyc-sim/internal/repository/interfaces"

	"gorm.io/gorm"
)

type DocumentRepo struct {
	db *gorm.DB
}

func NewDocumentRepo(db *gorm.DB) repoif.DocumentRepository {
	return &DocumentRepo{db: db}
}

func (r *DocumentRepo) Create(doc *models.Document) error {
	return r.db.Create(doc).Error
}

func (r *DocumentRepo) FindByCustomer(customerID string) ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Where("customer_id = ?", customerID).Find(&docs).Error
	return docs, err
}

func (r *DocumentRepo) ExistsByCustomerAndKind(customerID string, kind string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Document{}).
		Where("customer_id = ? AND kind = ?", customerID, kind).
		Count(&count).Error
	return count > 0, err
}
