package consumer

import "job-runner-app/internal/config"

type rabbitMQConsumer[T any] struct {
	rabbitConfig     config.RabbitMQConfig
	subscriberConfig config.SubscriberConfig
	close            func() error
	cancel           func()
}

func newRabbitMQConsumer[T any](
	rc config.RabbitMQConfig,
	sc config.SubscriberConfig,
) *rabbitMQConsumer[T] {
	return &rabbitMQConsumer[T]{
		rabbitConfig:     rc,
		subscriberConfig: sc,
		cancel:           func() {},
		close:            func() error { return nil },
	}
}

func (r *rabbitMQConsumer[T]) shutdown() error {
	r.cancel()
	return r.close()
}

func (r *rabbitMQConsumer[T]) consume(callback consumeCallback[T]) error {
	return nil
}
