package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"job-runner-app/internal/config"
	"job-runner-app/internal/exception"
	"job-runner-app/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type WorkerRepository interface {
	Get(context.Context, string) (*model.Worker, error)
	Create(context.Context, string, string) (string, error)
	Delete(context.Context, string) (bool, error)
}

type workerRepository struct {
	config config.DBConfig
}

func NewWorkerRepository(c config.DBConfig) WorkerRepository {
	return &workerRepository{config: c}
}

func (r *workerRepository) Get(ctx context.Context, id string) (*model.Worker, error) {
	dbConn, err := sqlx.Connect("postgres", BuildDataSourceName(r.config))
	if err != nil {
		return nil, err
	}

	var w model.Worker
	if err = dbConn.GetContext(ctx, &w, "SELECT * FROM get_worker($1)", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &w, nil
}

func (r *workerRepository) Create(ctx context.Context, name, description string) (string, error) {
	dbConn, err := sqlx.Connect("postgres", BuildDataSourceName(r.config))
	if err != nil {
		return "", err
	}
	var id string
	err = dbConn.GetContext(ctx, &id, "SELECT create_worker($1, $2)", name, description)
	if err != nil {
		var pErr *pq.Error
		if errors.As(err, &pErr) && pErr.Code == "23505" {
			// duplicate ke value violates unique constraint "workers_name_key"
			return "", exception.JobError{
				Code:    exception.WorkerAlreadyExists,
				Message: fmt.Sprintf("worker with name '%s' already exists", name),
			}
		}

		return "", err
	}

	return id, err
}

func (r *workerRepository) Delete(ctx context.Context, id string) (bool, error) {
	dbConn, err := sqlx.Connect("postgres", BuildDataSourceName(r.config))
	if err != nil {
		return false, err
	}

	log.Println("...repository.worker.Delete() id:", id)

	var deletedId *string
	err = dbConn.GetContext(ctx, &deletedId, "SELECT delete_worker($1)", id)

	log.Println("...repository.worker.Delete() deletedId:", deletedId)

	return deletedId != nil && *deletedId == id, err
}
