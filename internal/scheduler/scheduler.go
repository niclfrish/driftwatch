package scheduler

import (
	"context"
	"log"
	"time"
)

// Job is a function that can be scheduled to run periodically.
type Job func(ctx context.Context) error

// Scheduler runs a job at a fixed interval.
type Scheduler struct {
	interval time.Duration
	job      Job
	logger   *log.Logger
}

// New creates a new Scheduler with the given interval and job.
func New(interval time.Duration, job Job, logger *log.Logger) *Scheduler {
	if logger == nil {
		logger = log.Default()
	}
	return &Scheduler{
		interval: interval,
		job:      job,
		logger:   logger,
	}
}

// Run starts the scheduler loop, executing the job immediately and then
// at each tick of the configured interval. It blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) error {
	if err := s.runJob(ctx); err != nil {
		s.logger.Printf("scheduler: job error: %v", err)
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Println("scheduler: shutting down")
			return ctx.Err()
		case <-ticker.C:
			if err := s.runJob(ctx); err != nil {
				s.logger.Printf("scheduler: job error: %v", err)
			}
		}
	}
}

func (s *Scheduler) runJob(ctx context.Context) error {
	s.logger.Println("scheduler: running job")
	return s.job(ctx)
}
