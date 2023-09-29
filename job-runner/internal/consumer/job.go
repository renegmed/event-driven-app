package consumer

import (
	"context"
	"job-runner-app/internal/config"
	"job-runner-app/pkg/model"
)

type jobStatusConsumer struct {
	*rabbitMQConsumer[model.JobEventMessage]
}

func NewJobStatusConsumer(rc config.RabbitMQConfig) Consumer {

	return &jobStatusConsumer{
		newRabbitMQConsumer[model.JobEventMessage](rc),
	}
}

func (c *jobStatusConsumer) Consume() error {
	return nil
}

func (c *jobStatusConsumer) Shutdown(ctx context.Context) error {
	return nil
}
