package service_test

import (
	"context"
	"time"

	"job-runner-app/internal/broker"
	"job-runner-app/internal/model"
	"job-runner-app/internal/repository"
	pkgmodel "job-runner-app/pkg/model"

	"github.com/stretchr/testify/mock"
)

type workerRepoMock struct {
	mock.Mock
}

func (m *workerRepoMock) Create(ctx context.Context, name, topic string) (string, error) {
	args := m.Called(ctx, name, topic)
	return args.String(0), args.Error(1)
}

func (m *workerRepoMock) Get(ctx context.Context, id string) (*model.Worker, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Worker), args.Error(1)
}

func (m *workerRepoMock) Delete(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

/*
This line of code creates a static type assertion.

It tells the Go compiler to check if the workerRepoMock type
satisfies the repository.WorkerRepository interface

It doesn't create a variable or have any runtime effect; it's
solely for the purpose of static type checking during compilation
*/
var _ repository.WorkerRepository = (*workerRepoMock)(nil)

///////////////////////////////

type brokerMock[T any] struct {
	mock.Mock
}

func (m *brokerMock[T]) Publish(ctx context.Context, topic string, data T) error {
	args := m.Called(ctx, topic, data)
	return args.Error(0)
}

/*
This line of code creates a static type assertion.

It tells the Go compiler to check if the brokerMock type
satisfies the broker.Broker interface

It doesn't create a variable or have any runtime effect; it's
solely for the purpose of static type checking during compilation
*/
var _ broker.Broker[pkgmodel.JobLaunchMessage] = (*brokerMock[pkgmodel.JobLaunchMessage])(nil)

//////////////////////////////

type jobRepoMock struct {
	mock.Mock
}

func (r *jobRepoMock) AppendStatus(
	ctx context.Context,
	jobID string,
	message string,
	timestamp time.Time,
	output map[string]any,
) (int64, error) {
	args := r.Called(ctx, jobID, message, timestamp, output)
	return args.Get(0).(int64), args.Error(1)
}

func (r *jobRepoMock) GetStatuses(
	ctx context.Context,
	jobID string) ([]model.JobStatus, error) {
	args := r.Called(ctx, jobID)
	return args.Get(0).([]model.JobStatus), args.Error(1)
}

/*
This line of code creates a static type assertion.

It tells the Go compiler to check if the jobRepoMock type
satisfies the repository.JobRepository interface

It doesn't create a variable or have any runtime effect; it's
solely for the purpose of static type checking during compilation
*/
var _ repository.JobRepository = (*jobRepoMock)(nil)
