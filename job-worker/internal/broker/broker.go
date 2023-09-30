package broker

import (
	"context"
	"encoding/json"
	"job-worker-app/internal/config"

	pkgmodel "job-worker-app/pkg/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker[T any] interface {
	Publish(context.Context, string, T) error
}

type jobEventsBroker struct {
	rabbitMQConfig  config.RabbitMQConfig
	publisherConfig config.PublisherConfig
}

func NewJobEventsBroker(
	rc config.RabbitMQConfig,
	pc config.PublisherConfig,
) Broker[pkgmodel.JobEventMessage] {
	return &jobEventsBroker{
		rabbitMQConfig:  rc,
		publisherConfig: pc,
	}
}

func (b *jobEventsBroker) Publish(
	ctx context.Context,
	routingKey string,
	msg pkgmodel.JobEventMessage,
) error {
	ch, close, err := NewRabbitMQChannel(b.rabbitMQConfig)
	if err != nil {
		return err
	}

	defer close()

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		ctx,
		b.publisherConfig.Exchange,
		routingKey,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		},
	)
}
