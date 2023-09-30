package broker

import (
	"errors"
	"fmt"
	"job-worker-app/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

var emptyClose = func() error { return nil }

func BuildRabbitMQURL(c config.RabbitMQConfig) string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
	)
}
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
