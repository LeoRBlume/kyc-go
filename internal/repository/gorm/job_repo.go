package gormrepo

import (
	"errors"
	"kyc-sim/internal/domain"
	"kyc-sim/internal/domain/models"
	repoif "kyc-sim/internal/repository/interfaces"

	"gorm.io/gorm"
)

type JobRepo struct {
	db *gorm.DB
}

func NewJobRepo(db *gorm.DB) repoif.JobRepository {
	return &JobRepo{db: db}
}

func (r *JobRepo) FindActiveRunChecksJob(customerID string) (*models.Job, error) {
	var job models.Job
	err := r.db.
		Where("customer_id = ? AND kind = ? AND status IN ?", customerID, domain.JobRunChecks, []domain.JobStatus{domain.JobPending, domain.JobRunning}).
		First(&job).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &job, err
}

func (r *JobRepo) CreateJobWithItems(job *models.Job, items []models.JobItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(job).Error; err != nil {
			return err
		}
		for _, it := range items {
			if err := tx.Create(&it).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
