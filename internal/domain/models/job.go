package models

import (
	"kyc-sim/internal/domain"
	"time"
)

type Job struct {
	ID          string           `gorm:"primaryKey;type:text"`
	CustomerID  string           `gorm:"type:text;not null;index"`
	Kind        domain.JobKind   `gorm:"type:text;not null"`
	Status      domain.JobStatus `gorm:"type:text;not null"`
	Attempt     int              `gorm:"not null"`
	MaxAttempts int              `gorm:"not null"`
	NextRunAt   time.Time        `gorm:"not null"`
	LockedBy    *string          `gorm:"type:text"`
	LockedAt    *time.Time
	HeartbeatAt *time.Time
	LastError   *string   `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}
