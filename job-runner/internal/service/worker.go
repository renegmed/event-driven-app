package service

import (
	"context"
	"job-runner-app/internal/model"
	"job-runner-app/internal/repository"
)

type WorkerService interface {
	Get(context.Context, string) (*model.Worker, error)
	Create(context.Context, string, string) (string, error)
	Delete(context.Context, string) (bool, error)
}
type workerService struct {
	repository repository.WorkerRepository
}

func NewWorkerService(r repository.WorkerRepository) WorkerService {
	return &workerService{repository: r}
}

func (s *workerService) Get(ctx context.Context, id string) (*model.Worker, error) {
	return s.repository.Get(ctx, id)
}

func (s *workerService) Create(ctx context.Context, name, description string) (string, error) {
	return s.repository.Create(ctx, name, description)
}

func (s *workerService) Delete(ctx context.Context, id string) (bool, error) {
	return s.repository.Delete(ctx, id)
}
