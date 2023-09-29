package service_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"job-runner-app/internal/model"
	"job-runner-app/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	data := []struct {
		testName string
		name     string
		desc     string
		res      string
		err      error
	}{
		{
			testName: "error",
			name:     "name",
			desc:     "desc",
			res:      "",
			err:      fmt.Errorf("test error"),
		},
		{
			testName: "ok",
			name:     "name",
			desc:     "desc",
			res:      uuid.NewString(),
			err:      nil,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			r := &workerRepoMock{}
			r.On("Create", mock.Anything, d.name, d.desc).Return(d.res, d.err)

			s := service.NewWorkerService(r)
			res, err := s.Create(context.Background(), d.name, d.desc)

			log.Printf("\n___test.service.Create() '%s'. d.res: %v, res: %v, error: %v", d.testName, d.res, res, err)

			assert.Equal(t, d.res, res)
			if d.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	data := []struct {
		testName string
		id       string
		res      *model.Worker
		err      error
	}{
		{
			testName: "error",
			id:       uuid.NewString(),
			res:      nil,
			err:      fmt.Errorf("test error"),
		},
		{
			testName: "ok",
			id:       uuid.NewString(),
			res: &model.Worker{
				ID: uuid.NewString(),
			},
			err: nil,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			r := &workerRepoMock{}
			r.On("Get", mock.Anything, d.id).Return(d.res, d.err)

			s := service.NewWorkerService(r)
			res, err := s.Get(context.Background(), d.id)

			log.Printf("\n___test.service.Get() '%s'. d.res: %v, res: %v, error: %v", d.testName, d.res, res, err)

			assert.Equal(t, d.res, res)
			if d.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestDelete(t *testing.T) {
	data := []struct {
		testName string
		id       string
		res      bool
		err      error
	}{
		{
			testName: "error",
			id:       uuid.NewString(),
			res:      false,
			err:      fmt.Errorf("test error"),
		},
		{
			testName: "ok",
			id:       uuid.NewString(),
			res:      true,
			err:      nil,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			r := &workerRepoMock{}
			r.On("Delete", mock.Anything, d.id).Return(d.res, d.err)

			s := service.NewWorkerService(r)
			res, err := s.Delete(context.Background(), d.id)

			log.Printf("\n___test.service.Delete() '%s'. d.res: %v, res: %v, error: %v", d.testName, d.res, res, err)

			assert.Equal(t, d.res, res)
			if d.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
