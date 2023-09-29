package service_test

import (
	"context"
	"fmt"
	"job-runner-app/internal/model"
	"job-runner-app/internal/service"
	"testing"
	"time"

	pkgmodel "job-runner-app/pkg/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLaunchJob(t *testing.T) {
	data := []struct {
		testName   string
		worker     *model.Worker
		workerErr  error
		publishErr error
		fails      bool
	}{
		{
			testName:  "worker error",
			workerErr: fmt.Errorf("test error"),
			fails:     true,
		},
		{
			testName: "worker nil",
			worker:   nil,
			fails:    true,
		},
		{
			testName:  "publish error",
			worker:    &model.Worker{Name: uuid.NewString()},
			workerErr: fmt.Errorf("test error"),
			fails:     true,
		},
		{
			testName: "ok",
			worker:   &model.Worker{Name: uuid.NewString()},
			fails:    false,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			workerID := uuid.NewString()

			workerRepo := &workerRepoMock{}
			workerRepo.On("Get", mock.Anything, workerID).Return(d.worker, d.workerErr)

			broker := &brokerMock[pkgmodel.JobLaunchMessage]{}
			if d.worker != nil {
				broker.On("Publish", mock.Anything, fmt.Sprintf("worker.%s", d.worker.Name), mock.Anything).Return(d.publishErr)
			}

			jobRepo := &jobRepoMock{}

			s := service.NewJobService(workerRepo, jobRepo, broker)
			jobID, err := s.LaunchJob(context.Background(), workerID, nil)

			if d.fails {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, jobID)
			}

		})
	}
}

func TestAppendJobStatus(t *testing.T) {
	data := []struct {
		testName string
		jobID    string
		message  string
		output   map[string]any
		fails    bool
	}{
		{
			testName: "ok",
			jobID:    uuid.NewString(),
			message:  "test message",
			output:   map[string]any{"test": "test"},
			fails:    false,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			jobRepo := &jobRepoMock{}
			jobRepo.On("AppendStatus", mock.Anything, d.jobID, mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)

			jobService := service.NewJobService(nil, jobRepo, nil)
			_, err := jobService.AppendJobStatus(context.Background(), d.jobID, d.message, time.Now(), d.output)

			if d.fails {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				jobRepo.AssertCalled(t, "AppendStatus", mock.Anything, d.jobID, d.message, mock.Anything, d.output)
			}
		})
	}
}

func TestGetJobStatuses(t *testing.T) {
	data := []struct {
		testName string
		jobID    string
		statuses []model.JobStatus
		err      error
	}{
		{
			testName: "ok",
			jobID:    uuid.NewString(),
			statuses: []model.JobStatus{
				{
					ID:        1,
					JobID:     uuid.NewString(),
					Message:   "test message",
					Timestamp: time.Now(),
					Output:    map[string]any{"test": "test"},
				},
				// {
				// 	ID:        2,
				// 	JobID:     uuid.NewString(),
				// 	Message:   "test message 2",
				// 	Timestamp: time.Now(),
				// 	Output:    map[string]any{"test": "test 2"},
				// },
			},
			err: nil,
		},
		{
			testName: "error",
			jobID:    uuid.NewString(),
			statuses: nil,
			err:      fmt.Errorf("test error"),
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			jobRepo := &jobRepoMock{}
			jobRepo.On("GetStatuses", mock.Anything, d.jobID).Return(d.statuses, d.err)

			jobService := service.NewJobService(nil, jobRepo, nil)
			statuses, err := jobService.GetJobStatuses(context.Background(), d.jobID)
			if d.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, d.statuses, statuses)
			}
		})
	}

}
