package repository_test

import (
	"context"
	"fmt"
	"job-runner-app/internal/config"
	"job-runner-app/internal/repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryAppendStatus(t *testing.T) {
	cfg, err := config.LoadConfigFrom("../../..")
	if err != nil {
		t.Fatal(err)
	}

	r := repository.NewJobRepository(cfg.DB)

	jobID := uuid.NewString()
	message := fmt.Sprintf("test-%s", uuid.NewString())
	timestamp := time.Now()
	output := map[string]any{
		"test": "test",
	}

	id, err := r.AppendStatus(context.Background(), jobID, message, timestamp, output)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	db, err := sqlx.Connect("postgres", repository.BuildDataSourceName(cfg.DB))
	if err != nil {
		t.Fatal(err)
	}

	var count int
	err = db.Get(&count, "SELECT  COUNT(*) FROM jobs WHERE job_id = $1 AND message = $2", jobID, message)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, count)

}

func TestRepositoryGetStatuses(t *testing.T) {
	cfg, err := config.LoadConfigFrom("../../..")
	if err != nil {
		t.Fatal(err)
	}

	r := repository.NewJobRepository(cfg.DB)

	jobID := uuid.NewString()
	message := fmt.Sprintf("test-%s", uuid.NewString())
	timestamp := time.Now()
	output := map[string]any{
		"test": "test",
	}

	id, err := r.AppendStatus(context.Background(), jobID, message, timestamp, output)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	statuses, err := r.GetStatuses(context.Background(), jobID)
	assert.NoError(t, err)
	assert.NotEmpty(t, statuses)
	assert.Equal(t, 1, len(statuses))
	assert.Equal(t, message, statuses[0].Message)
	assert.Equal(t, output, statuses[0].Output)
	assert.Equal(t, id, statuses[0].ID)
	assert.Equal(t, jobID, statuses[0].JobID)
}

func TestRepositoryGetStatusesNoJob(t *testing.T) {
	cfg, err := config.LoadConfigFrom("../../..")
	if err != nil {
		t.Fatal(err)
	}

	r := repository.NewJobRepository(cfg.DB)

	statuses, err := r.GetStatuses(context.Background(), uuid.NewString())
	assert.NoError(t, err)
	assert.Empty(t, statuses)
}
