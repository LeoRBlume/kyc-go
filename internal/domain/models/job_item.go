package models

import (
	"kyc-sim/internal/domain"
	"time"
)

type JobItem struct {
	ID         string           `gorm:"primaryKey;type:text"`
	JobID      string           `gorm:"type:text;not null;index"`
	CheckType  domain.CheckType `gorm:"type:text;not null"`
	Status     string           `gorm:"type:text;not null"` // PENDING/DONE/FAILED
	StartedAt  *time.Time
	FinishedAt *time.Time
	Error      *string `gorm:"type:text"`
}
