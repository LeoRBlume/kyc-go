package models

import (
	"kyc-sim/internal/domain"
	"time"
)

type Check struct {
	ID         string             `gorm:"primaryKey;type:text"`
	CustomerID string             `gorm:"type:text;not null;index"`
	Type       domain.CheckType   `gorm:"type:text;not null"`
	Status     domain.CheckStatus `gorm:"type:text;not null"`
	Score      *int               `gorm:"type:int"`
	Details    string             `gorm:"type:text"`
	CreatedAt  time.Time          `gorm:"not null"`
	UpdatedAt  time.Time          `gorm:"not null"`
}
