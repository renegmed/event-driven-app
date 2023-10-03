package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"job-runner-app/internal/broker"
	"job-runner-app/internal/config"
	pkgmodel "job-runner-app/pkg/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLaunchJobShouldDoIt(t *testing.T) {
	client := newClient()
	wname, wdesc := fmt.Sprintf("test-%s", uuid.NewString()), fmt.Sprintf("test-%s", uuid.NewString())
	r, err := client.CreateWorker(
		wname,
		wdesc,
	)

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, r.ID)

	cfg, err := config.LoadConfigFrom("../..")
	if err != nil {
		t.Fatal(err)
	}

	ch, close, err := broker.NewRabbitMQChannel(cfg.RabbitMQ)
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	queue, err := ch.QueueDeclare(
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
		queue.Name,
		fmt.Sprintf("worker.%s", wname),
		cfg.Jobs.Exchange,
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	d, err := ch.Consume(
		queue.Name,
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

	input := map[string]any{
		"test": "test",
	}

	launchJobResp, err := client.LaunchJob(r.ID, input)
	if err != nil {
		t.Fatal(err)
	}

	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case msg := <-d:
		var m pkgmodel.JobLaunchMessage
		err := json.Unmarshal(msg.Body, &m)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, launchJobResp.JobID, m.JobID)
		assert.Equal(t, input, m.Input)
		assert.GreaterOrEqual(t, 5*time.Second, time.Since(time.UnixMilli(m.Timestamp)))
	case <-timeout.Done():
		t.Fatal(timeout.Err())
	}
}

func TestLaunchJobShouldReturnBadRequest(t *testing.T) {
	data := []struct {
		testName string
		workerID string
		input    map[string]any
	}{
		{
			testName: "worker not found",
			workerID: uuid.NewString(),
			input: map[string]any{
				"test": "test",
			},
		},
		{
			testName: "worker id is empty",
			workerID: "",
			input: map[string]any{
				"test": "test",
			},
		},
		{
			testName: "input is nil",
			workerID: uuid.NewString(),
			input:    nil,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			client := newClient()
			_, err := client.LaunchJob(d.workerID, d.input)
			assert.Error(t, err)
			assert.Equal(t, "unexpected status code: 400", err.Error())
		})
	}
}
