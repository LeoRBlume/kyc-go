package models

import (
	"kyc-sim/internal/domain"
	"time"
)

type Customer struct {
	ID        string                `gorm:"primaryKey;type:text"`
	Type      domain.CustomerType   `gorm:"type:text;not null"`
	Status    domain.CustomerStatus `gorm:"type:text;not null"`
	CreatedAt time.Time             `gorm:"not null"`
	UpdatedAt time.Time             `gorm:"not null"`
}
