package service

import (
	"context"
	"fmt"
	"job-runner-app/internal/broker"
	"job-runner-app/internal/exception"
	"job-runner-app/internal/model"
	"job-runner-app/internal/repository"
	"time"

	pkgmodel "job-runner-app/pkg/model"

	"github.com/google/uuid"
)

type JobService interface {
	LaunchJob(context.Context, string, map[string]any) (string, error)
	AppendJobStatus(context.Context, string, string, time.Time, map[string]any) (int64, error)
	GetJobStatuses(context.Context, string) ([]model.JobStatus, error)
}

type jobService struct {
	workerRepository repository.WorkerRepository
	jobRepository    repository.JobRepository
	broker           broker.Broker[pkgmodel.JobLaunchMessage]
}

func NewJobService(
	wr repository.WorkerRepository,
	jr repository.JobRepository,
	broker broker.Broker[pkgmodel.JobLaunchMessage],
) JobService {
	return &jobService{
		workerRepository: wr,
		jobRepository:    jr,
		broker:           broker,
	}
}

func (s *jobService) LaunchJob(ctx context.Context, workerID string, input map[string]any) (string, error) {
	worker, err := s.workerRepository.Get(ctx, workerID)
	if err != nil {
		return "", err
	}
	if worker == nil {
		return "", exception.JobError{
			Code:    exception.WorkerNotFound,
			Message: fmt.Sprintf("worker '%s' not found", workerID),
		}
	}

	jobID := uuid.NewString()
	err = s.broker.Publish(
		ctx,
		fmt.Sprintf("worker.%s", worker.Name),
		pkgmodel.JobLaunchMessage{
			JobID:     jobID,
			Timestamp: time.Now().UnixMilli(),
			Input:     input,
		},
	)

	return jobID, err

}
func (s *jobService) AppendJobStatus(
	ctx context.Context,
	jobID, message string,
	timestamp time.Time,
	output map[string]any) (int64, error) {

	id, err := s.jobRepository.AppendStatus(ctx, jobID, message, timestamp, output)
	if err != nil {
		return 0, exception.JobError{
			Code:    exception.FailedToAppendJobStatus,
			Err:     err,
			Message: fmt.Sprintf("failed to append job status for job '%s'", jobID),
		}
	}

	return id, nil
}

func (s *jobService) GetJobStatuses(ctx context.Context, jobID string) ([]model.JobStatus, error) {
	statuses, err := s.jobRepository.GetStatuses(ctx, jobID)
	if err != nil {
		return nil, exception.JobError{
			Code:    exception.FailedToGetJobStatuses,
			Err:     err,
			Message: fmt.Sprintf("failed to get job statuses for job '%s'", jobID),
		}
	}

	return statuses, nil
}
