package repository

import (
	"context"
	"job-runner-app/internal/config"
	"log"

	"github.com/jmoiron/sqlx"
)

type HealthRepository interface {
	Health(context.Context) error
}

type healthRepository struct {
	config config.DBConfig
}

func NewHealthRepository(c config.DBConfig) HealthRepository {
	log.Println("...NewHealthRepository() config.DBConfig:", c)
	return &healthRepository{config: c}
}

func (r *healthRepository) Health(ctx context.Context) error {
	dbConn, err := sqlx.Connect("postgres", BuildDataSourceName(r.config))
	if err != nil {
		return err
	}

	return dbConn.PingContext(ctx)
}
