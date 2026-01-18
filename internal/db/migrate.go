package db

import (
	"kyc-sim/internal/domain/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Customer{},
		&models.Document{},
		&models.Check{},
	)
}
