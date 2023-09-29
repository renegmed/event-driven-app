package service_test

import (
	"context"

	"job-runner-app/internal/model"

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
