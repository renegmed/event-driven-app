package consumer

import (
	"context"
	"job-runner-app/internal/config"
	"job-runner-app/internal/service"
	pkgmodel "job-runner-app/pkg/model"
	"log"
	"time"
)

type jobStatusConsumer struct {
	*rabbitMQConsumer[pkgmodel.JobEventMessage]
	jobService service.JobService
}

func NewJobStatusConsumer(
	rc config.RabbitMQConfig,
	sc config.SubscriberConfig,
	jobService service.JobService,
) Consumer {

	return &jobStatusConsumer{
		newRabbitMQConsumer[pkgmodel.JobEventMessage](rc, sc), // ??????? explain ????
		jobService, // ??????? explain ????
	}
}

func (c *jobStatusConsumer) Consume() error {
	return c.consume(c.consumeCallback)
}

func (c *jobStatusConsumer) Shutdown(ctx context.Context) error {
	return c.shutdown()
}

func (c *jobStatusConsumer) consumeCallback(ctx context.Context, m pkgmodel.JobEventMessage, err error) error {
	if err != nil {
		log.Printf("consuming job statuses error: %v", err)
	}
	_, err = c.jobService.AppendJobStatus(
		ctx,
		m.JobID,
		m.Message,
		time.UnixMilli(m.Timestamp),
		m.Output,
	)

	return err
}
