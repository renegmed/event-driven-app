package broker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"job-runner-app/internal/config"
	pkgmodel "job-runner-app/pkg/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

type jobLaunchBroker struct {
	rabbitMQConfig  config.RabbitMQConfig
	publisherConfig config.PublisherConfig
}

func NewJobLaunchBroker(
	rc config.RabbitMQConfig,
	pc config.PublisherConfig,
) Broker[pkgmodel.JobLaunchMessage] {
	return &jobLaunchBroker{
		rabbitMQConfig:  rc,
		publisherConfig: pc,
	}
}

func (b *jobLaunchBroker) Publish(
	ctx context.Context,
	key string,
	job pkgmodel.JobLaunchMessage,
) error {

	ch, close, err := NewRabbitMQChannel(b.rabbitMQConfig)
	if err != nil {
		return err
	}
	defer close()

	bytes, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		ctx,
		b.publisherConfig.Exchange, // 'jobs'
		key,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		},
	)
}

var emptyClose = func() error { return nil }

func NewRabbitMQChannel(config config.RabbitMQConfig) (*amqp.Channel, func() error, error) {
	conn, err := amqp.Dial(BuildRabbitMQURL(config))
	if err != nil {
		return nil, emptyClose, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, emptyClose, err
	}

	close := func() error {
		return errors.Join(
			conn.Close(),
			ch.Close(),
		)
	}

	return ch, close, nil
}

func BuildRabbitMQURL(c config.RabbitMQConfig) string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
	)
}
