package e2e_test

import (
	"encoding/json"
	"fmt"
	"job-worker-app/internal/broker"
	"job-worker-app/internal/config"
	"testing"
	"time"

	"golang.org/x/net/context"

	pkgmodel "job-worker-app/pkg/model"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func TestWorkerShouldTrackJobProgress(t *testing.T) {
	cfg, err := config.LoadConfigFrom("../..")
	if err != nil {
		t.Fatal(err)
	}

	ch, close, err := broker.NewRabbitMQChannel(cfg.RabbitMQ)
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	q, err := ch.QueueDeclare(
		fmt.Sprintf("test-%s", uuid.NewString()),
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = ch.QueueBind(
		q.Name,
		"event",
		cfg.JobEvents.Exchange,
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	jobLaunchMessage := pkgmodel.JobLaunchMessage{
		JobID:     uuid.NewString(),
		Timestamp: time.Now().UnixMilli(),
		Input: map[string]any{
			"key":   "value",
			"value": "key",
		},
	}
	b, err := json.Marshal(jobLaunchMessage)
	if err != nil {
		t.Fatal(err)
	}
	ch.PublishWithContext(
		context.Background(),
		"jobs",
		"worker.random",
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

outer:
	for {
		select {
		case <-ctx.Done():
			t.Fatal("timeout")
		case msg := <-msgs:
			var event pkgmodel.JobEventMessage
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				t.Fatal(err)
			}

			if event.JobID == jobLaunchMessage.JobID {
				assert.GreaterOrEqual(t, event.Timestamp, jobLaunchMessage.Timestamp)
				break outer
			}
		}
	}

}
