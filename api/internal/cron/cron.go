package cron

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Job struct {
	Name     string
	Schedule string
	Function any
	Params   []any
}

type JobStatus struct {
	Name     string    `json:"name"`
	Schedule string    `json:"schedule"`
	LastRun  time.Time `json:"last_run"`
	NextRun  time.Time `json:"next_run"`
	IsRunning bool     `json:"is_running"`
}

type Scheduler struct {
	scheduler    gocron.Scheduler
	jobs         map[string]gocron.Job
	jobSchedules map[string]string
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewScheduler(db *sql.DB) (*Scheduler, error) {
	locker, err := NewPostgresLocker(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres locker: %w", err)
	}

	scheduler, err := gocron.NewScheduler(
		gocron.WithDistributedLocker(locker),
		gocron.WithLocation(time.UTC),
		gocron.WithGlobalJobOptions(
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		scheduler:    scheduler,
		jobs:         make(map[string]gocron.Job),
		jobSchedules: make(map[string]string),
		ctx:          ctx,
		cancel:       cancel,
	}, nil
}

func (s *Scheduler) RegisterJob(job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.Name]; exists {
		return fmt.Errorf("job %s already registered", job.Name)
	}

	var scheduledJob gocron.Job
	var err error

	switch job.Schedule {
	case "@every 1s":
		scheduledJob, err = s.scheduler.NewJob(
			gocron.DurationJob(1*time.Second),
			gocron.NewTask(job.Function, job.Params...),
			gocron.WithName(job.Name),
			gocron.WithEventListeners(
				gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
					fmt.Printf("Job %s (ID: %s) is starting\n", jobName, jobID.String())
				}),
				gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
					fmt.Printf("Job %s (ID: %s) completed\n", jobName, jobID.String())
				}),
				gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
					fmt.Printf("Job %s (ID: %s) failed with error: %v\n", jobName, jobID.String(), err)
				}),
			),
		)
	default:
		if isCronExpression(job.Schedule) {
			scheduledJob, err = s.scheduler.NewJob(
				gocron.CronJob(job.Schedule, false),
				gocron.NewTask(job.Function, job.Params...),
				gocron.WithName(job.Name),
				gocron.WithEventListeners(
					gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
						slog.Info("Job is starting", "name", jobName, "id", jobID.String())
					}),
					gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
						slog.Info("Job completed", "name", jobName, "id", jobID.String())
					}),
					gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
						slog.Error("Job failed", "name", jobName, "id", jobID.String(), "error", err)
					}),
				),
			)
		} else if duration, parseErr := time.ParseDuration(job.Schedule); parseErr == nil {
			scheduledJob, err = s.scheduler.NewJob(
				gocron.DurationJob(duration),
				gocron.NewTask(job.Function, job.Params...),
				gocron.WithName(job.Name),
				gocron.WithEventListeners(
					gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
						slog.Info("Job is starting", "name", jobName, "id", jobID.String())
					}),
					gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
						slog.Info("Job completed", "name", jobName, "id", jobID.String())
					}),
					gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
						slog.Error("Job failed", "name", jobName, "id", jobID.String(), "error", err)
					}),
				),
			)
		} else {
			return fmt.Errorf("invalid schedule format: %s", job.Schedule)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to create job %s: %w", job.Name, err)
	}

	s.jobs[job.Name] = scheduledJob
	s.jobSchedules[job.Name] = job.Schedule
	return nil
}

func (s *Scheduler) Start() error {
	s.scheduler.Start()
	return nil
}

func (s *Scheduler) Stop() error {
	s.cancel()
	return s.scheduler.Shutdown()
}

func (s *Scheduler) RemoveJob(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.jobs[name]
	if !exists {
		return fmt.Errorf("job %s not found", name)
	}

	if err := s.scheduler.RemoveJob(job.ID()); err != nil {
		return fmt.Errorf("failed to remove job %s: %w", name, err)
	}

	delete(s.jobs, name)
	delete(s.jobSchedules, name)
	return nil
}

func (s *Scheduler) GetJob(name string) (gocron.Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, exists := s.jobs[name]
	return job, exists
}

func (s *Scheduler) ListJobs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	names := make([]string, 0, len(s.jobs))
	for name := range s.jobs {
		names = append(names, name)
	}
	return names
}

func (s *Scheduler) GetJobStatuses() ([]JobStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	statuses := make([]JobStatus, 0, len(s.jobs))
	for name, job := range s.jobs {
		status := JobStatus{
			Name:     name,
			Schedule: s.jobSchedules[name],
		}
		
		// Get last run time
		if lastRun, err := job.LastRun(); err == nil {
			status.LastRun = lastRun
		}
		
		// Get next run time
		if nextRun, err := job.NextRun(); err == nil {
			status.NextRun = nextRun
		}
		
		// Check if job is currently running
		// Note: gocron v2 doesn't expose running state directly
		// We approximate by checking if current time is close to last run
		now := time.Now()
		if !status.LastRun.IsZero() && now.Sub(status.LastRun) < 5*time.Second {
			status.IsRunning = true
		}
		
		statuses = append(statuses, status)
	}
	
	return statuses, nil
}

func (s *Scheduler) RunJobNow(name string) error {
	s.mu.RLock()
	job, exists := s.jobs[name]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("job %s not found", name)
	}

	return job.RunNow()
}

func isCronExpression(schedule string) bool {
	fields := 0
	for i := 0; i < len(schedule); i++ {
		if schedule[i] == ' ' {
			fields++
		}
	}
	return fields >= 4 && fields <= 6
}
