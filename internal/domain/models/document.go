package models

import (
	"kyc-sim/internal/domain"
	"time"
)

type Document struct {
	ID         string                `gorm:"primaryKey;type:text"`
	CustomerID string                `gorm:"type:text;not null;index"`
	Kind       domain.DocumentKind   `gorm:"type:text;not null"`
	FileURL    string                `gorm:"type:text;not null"`
	Status     domain.DocumentStatus `gorm:"type:text;not null"`
	UploadedAt time.Time             `gorm:"not null"`
}
