package impl

import (
	"errors"

	"kyc-sim/internal/domain"
	"kyc-sim/internal/dto/responses"
	repoif "kyc-sim/internal/repository/interfaces"
	svcif "kyc-sim/internal/service/interfaces"

	"gorm.io/gorm"
)

type JobService struct {
	jobRepo repoif.JobRepository
}

func NewJobService(jobRepo repoif.JobRepository) svcif.JobService {
	return &JobService{jobRepo: jobRepo}
}

func (s *JobService) Get(jobID string) (*responses.GetJobResponse, error) {
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewNotFound("job not found", map[string]any{"jobId": jobID})
		}
		return nil, domain.NewInternal("failed to load job", map[string]any{"cause": err.Error()})
	}

	items, err := s.jobRepo.ListItems(job.ID)
	if err != nil {
		return nil, domain.NewInternal("failed to load job items", map[string]any{"cause": err.Error()})
	}

	respItems := make([]responses.JobItemResponse, 0, len(items))
	progress := responses.JobProgress{Total: len(items)}

	for _, it := range items {
		respItems = append(respItems, responses.JobItemResponse{
			CheckType:  string(it.CheckType),
			Status:     it.Status,
			StartedAt:  it.StartedAt,
			FinishedAt: it.FinishedAt,
			Error:      it.Error,
		})

		switch it.Status {
		case "DONE":
			progress.Done++
		case "FAILED":
			progress.Failed++
		}
	}

	return &responses.GetJobResponse{
		JobID:       job.ID,
		Kind:        string(job.Kind),
		Status:      string(job.Status),
		CustomerID:  job.CustomerID,
		Attempt:     job.Attempt,
		MaxAttempts: job.MaxAttempts,
		NextRunAt:   job.NextRunAt,
		LockedBy:    job.LockedBy,
		HeartbeatAt: job.HeartbeatAt,
		LastError:   job.LastError,
		CreatedAt:   job.CreatedAt,
		UpdatedAt:   job.UpdatedAt,
		Progress:    progress,
		Items:       respItems,
	}, nil
}
