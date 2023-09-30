package consumer

import (
	"context"
	"encoding/json"
	"job-worker-app/internal/broker"
	"job-worker-app/internal/config"
)

type rabbitMQConsumer[T any] struct {
	rabbitMQConfig   config.RabbitMQConfig
	subscriberConfig config.SubscriberConfig
	close            func() error
	cancel           func()
}

func newRabbitMQConsumer[T any](
	rabbitMQConfig config.RabbitMQConfig,
	subscriberConfig config.SubscriberConfig,
) *rabbitMQConsumer[T] {
	return &rabbitMQConsumer[T]{
		rabbitMQConfig:   rabbitMQConfig,
		subscriberConfig: subscriberConfig,
		cancel:           func() {},
		close:            func() error { return nil },
	}
}

func (r *rabbitMQConsumer[T]) consume(callback consumeCallback[T]) error {
	ch, close, err := broker.NewRabbitMQChannel(r.rabbitMQConfig)
	if err != nil {
		return err
	}
	defer close()

	msgs, err := ch.Consume(
		r.subscriberConfig.Queue,
		r.subscriberConfig.Group,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r.cancel = cancel

	for {
		select {
		case <-ctx.Done():
			return ErrConsumerClosed
		case msg := <-msgs:
			var m T
			if err := json.Unmarshal(msg.Body, &m); err != nil {
				msg.Nack(false, false)
				callback(ctx, m, err)
				continue
			}

			if err := callback(ctx, m, nil); err != nil {
				msg.Nack(false, true)
				callback(ctx, m, err)
				continue
			}

			if err := msg.Ack(false); err != nil {
				msg.Nack(false, false)
				callback(ctx, m, err)
			}
		}
	}
}

func (r *rabbitMQConsumer[T]) shutdown() error {
	r.cancel()
	return r.close()
}
