package service

import (
	"context"
	"fmt"
	"job-search-service/internal/models"
	"job-search-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type JobService struct {
	repo *repository.JobRepository
}

func NewJobService(repo *repository.JobRepository) *JobService {
	return &JobService{
		repo: repo,
	}
}

func (s *JobService) CreateJob(ctx context.Context, title, description, company, location string, skills []string, salary float64) (string, error) {
	job := &models.Job{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Company:     company,
		Location:    location,
		Skills:      skills,
		Salary:      salary,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, job); err != nil {
		return "", fmt.Errorf("failed to create job: %w", err)
	}

	return job.ID, nil
}

func (s *JobService) SearchJobs(ctx context.Context, query, location string, skills []string) ([]*models.Job, error) {
	jobs, err := s.repo.Search(ctx, query, location, skills)
	if err != nil {
		return nil, fmt.Errorf("failed to search jobs: %w", err)
	}

	return jobs, nil
}

func (s *JobService) GetJob(ctx context.Context, id string) (*models.Job, error) {
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return job, nil
}

func (s *JobService) DeleteJob(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	return nil
}
