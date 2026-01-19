package gormrepo

import (
	"errors"
	"kyc-sim/internal/domain"
	"kyc-sim/internal/domain/models"
	repoif "kyc-sim/internal/repository/interfaces"
	"time"

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

func (r *JobRepo) ClaimNext(now time.Time, workerID string) (*models.Job, error) {
	var job models.Job

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("status = ? AND next_run_at <= ?", domain.JobPending, now).
			Order("created_at").
			First(&job).Error; err != nil {
			return err
		}

		job.Status = domain.JobRunning
		job.LockedBy = &workerID
		job.LockedAt = &now
		job.HeartbeatAt = &now
		job.UpdatedAt = now

		return tx.Save(&job).Error
	})

	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepo) UpdateHeartbeat(jobID string, now time.Time) error {
	return r.db.Model(&models.Job{}).
		Where("id = ?", jobID).
		Update("heartbeat_at", now).Error
}

func (r *JobRepo) MarkDone(jobID string) error {
	return r.db.Model(&models.Job{}).
		Where("id = ?", jobID).
		Updates(map[string]any{
			"status":     domain.JobDone,
			"updated_at": time.Now(),
		}).Error
}

func (r *JobRepo) MarkFailed(jobID string, errMsg string) error {
	return r.db.Model(&models.Job{}).
		Where("id = ?", jobID).
		Updates(map[string]any{
			"status":     domain.JobFailed,
			"last_error": errMsg,
			"updated_at": time.Now(),
		}).Error
}

func (r *JobRepo) ListItems(jobID string) ([]models.JobItem, error) {
	var items []models.JobItem
	err := r.db.Where("job_id = ?", jobID).Find(&items).Error
	return items, err
}

func (r *JobRepo) UpdateItem(item *models.JobItem) error {
	return r.db.Save(item).Error
}
