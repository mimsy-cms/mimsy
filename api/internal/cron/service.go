package cron

import (
	"context"
	"database/sql"
	"fmt"
)

type CronService interface {
	RegisterJob(ctx context.Context, job Job) error
	RemoveJob(ctx context.Context, name string) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	RunJobNow(ctx context.Context, name string) error
	ListJobs(ctx context.Context) []string
	GetJobStatuses(ctx context.Context) ([]JobStatus, error)
}

type cronService struct {
	scheduler *Scheduler
}

func NewCronService(db *sql.DB) (CronService, error) {
	scheduler, err := NewScheduler(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	return &cronService{
		scheduler: scheduler,
	}, nil
}

func (s *cronService) RegisterJob(ctx context.Context, job Job) error {
	return s.scheduler.RegisterJob(job)
}

func (s *cronService) RemoveJob(ctx context.Context, name string) error {
	return s.scheduler.RemoveJob(name)
}

func (s *cronService) Start(ctx context.Context) error {
	return s.scheduler.Start()
}

func (s *cronService) Stop(ctx context.Context) error {
	return s.scheduler.Stop()
}

func (s *cronService) RunJobNow(ctx context.Context, name string) error {
	return s.scheduler.RunJobNow(name)
}

func (s *cronService) ListJobs(ctx context.Context) []string {
	return s.scheduler.ListJobs()
}

func (s *cronService) GetJobStatuses(ctx context.Context) ([]JobStatus, error) {
	return s.scheduler.GetJobStatuses()
}