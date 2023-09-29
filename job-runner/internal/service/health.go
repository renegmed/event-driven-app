package service

import (
	"context"
	"job-runner-app/internal/exception"
	"job-runner-app/internal/repository"
)

type HealthService interface {
	Health(context.Context) error
}

type healthService struct {
	healthRepository repository.HealthRepository
}

func NewHealthService(r repository.HealthRepository) HealthService {
	return &healthService{healthRepository: r}
}

func (s *healthService) Health(ctx context.Context) error {
	err := s.healthRepository.Health(ctx)
	if err != nil {
		return exception.JobError{
			Code:    exception.UnhealthService,
			Message: "failed to check health",
			Err:     err,
		}
	}
	return nil
}
