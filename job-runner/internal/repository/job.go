package repository

import (
	"context"
	"encoding/json"
	"job-runner-app/internal/config"
	"job-runner-app/internal/model"
	"log"
	"time"

	pie "github.com/elliotchance/pie/v2"
	"github.com/jmoiron/sqlx"
)

type JobRepository interface {
	AppendStatus(context.Context, string, string, time.Time, map[string]any) (int64, error)
	GetStatuses(context.Context, string) ([]model.JobStatus, error)
}

type jobRepository struct {
	config config.DBConfig
}

func NewJobRepository(c config.DBConfig) JobRepository {
	log.Println("...repository.NewJobRepository() config.DBConfig:", c)
	return &jobRepository{config: c}
}

func (r *jobRepository) AppendStatus(
	ctx context.Context,
	jobID string,
	message string,
	timestamp time.Time,
	output map[string]any,
) (int64, error) {

	db, err := sqlx.Connect("postgres", BuildDataSourceName(r.config))
	if err != nil {
		return 0, err
	}
	b, err := json.Marshal(output)
	if err != nil {
		return 0, err
	}

	var id int64
	err = db.GetContext(ctx, &id, "SELECT append_job_status ($1, $2, $3, $4)", jobID, message, timestamp, b)
	return id, err
}

func (r *jobRepository) GetStatuses(ctx context.Context, jobID string) ([]model.JobStatus, error) {
	db, err := sqlx.Connect("postgres", BuildDataSourceName(r.config))
	if err != nil {
		return nil, err
	}

	log.Println("...repository.job.GetStatuses() jobID:", jobID)
	var statuses []model.JobStatusDTO
	if err = db.SelectContext(ctx, &statuses, "SELECT * FROM job_statuses($1)", jobID); err != nil {
		log.Println("...repository.job.GetStatuses() SELECT * FROM job_statuses error:", err)
		return nil, err
	}

	var mapErr error
	result := pie.Map(statuses, func(s model.JobStatusDTO) model.JobStatus {
		var output map[string]any
		if err := json.Unmarshal([]byte(s.Output), &output); err != nil {
			mapErr = err
		}

		return model.JobStatus{
			Message:   s.Message,
			Timestamp: s.CreatedAt,
			Output:    output,
			ID:        s.ID,
			JobID:     s.JobID,
		}
	})

	return result, mapErr

}
