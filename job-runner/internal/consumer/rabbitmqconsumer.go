package consumer

import "job-runner-app/internal/config"

type rabbitMQConsumer[T any] struct {
	rabbitConfig config.RabbitMQConfig
}

func newRabbitMQConsumer[T any](rc config.RabbitMQConfig) *rabbitMQConsumer[T] {
	return &rabbitMQConsumer[T]{rabbitConfig: rc}
}
