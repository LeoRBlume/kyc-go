package worker

import (
	"log"
	"time"

	repoif "kyc-sim/internal/repository/interfaces"
)

type Runner struct {
	jobRepo   repoif.JobRepository
	processor *Processor
	workerID  string
}

func NewRunner(j repoif.JobRepository, p *Processor, workerID string) *Runner {
	return &Runner{jobRepo: j, processor: p, workerID: workerID}
}

func (r *Runner) Start() {
	for {
		job, err := r.jobRepo.ClaimNext(time.Now(), r.workerID)
		if err != nil || job == nil {
			time.Sleep(5000 * time.Millisecond)
			continue
		}

		log.Printf("[worker %s] processing job %s", r.workerID, job.ID)

		if err := r.processor.Process(job); err != nil {
			_ = r.jobRepo.MarkFailed(job.ID, err.Error())
		}
	}
}
