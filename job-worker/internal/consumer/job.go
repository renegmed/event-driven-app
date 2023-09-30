package consumer

import (
	"context"
	"job-worker-app/internal/broker"
	"job-worker-app/internal/config"
	"job-worker-app/internal/helper"
	"job-worker-app/internal/model"
	"job-worker-app/internal/service"
	pkgmodel "job-worker-app/pkg/model"
	"log"
	"time"
)

type jobLaunchConsumer struct {
	*rabbitMQConsumer[pkgmodel.JobLaunchMessage] // generic type JobLaunchMessage
	workerService                                service.WorkerService
	broker                                       broker.Broker[pkgmodel.JobEventMessage]
}

func NewJobLaunchConsumer(
	rc config.RabbitMQConfig,
	cc config.SubscriberConfig,
	ws service.WorkerService,
	broker broker.Broker[pkgmodel.JobEventMessage],
) Consumer {

	//log.Println("...consumer.NewJobLaunchConsumer() config.SubscriberConfig Group:", cc.Group, " Queue:", cc.Queue)
	return &jobLaunchConsumer{
		rabbitMQConsumer: newRabbitMQConsumer[pkgmodel.JobLaunchMessage](rc, cc),
		workerService:    ws,
		broker:           broker,
	}
}

func (j *jobLaunchConsumer) Consume() error {
	return j.consume(j.consumeCallback)
}

func (j *jobLaunchConsumer) consumeCallback(ctx context.Context, msg pkgmodel.JobLaunchMessage, err error) error {

	if err != nil {
		log.Printf("job launch message error: %v", err)
		j.broker.Publish(
			ctx,
			"event",
			pkgmodel.JobEventMessage{
				JobID:     msg.JobID,
				Message:   "random worker error",
				Timestamp: time.Now().UnixMilli(),
				Output: map[string]any{
					"error": err.Error(),
				},
			},
		)

		return nil
	}

	log.Println("...consumer.consumeCallback() msg:", msg)

	err = j.broker.Publish(
		ctx,
		"event",
		pkgmodel.JobEventMessage{
			JobID:     msg.JobID,
			Message:   "random worker started",
			Timestamp: time.Now().UnixMilli(),
			Output:    msg.Input,
		},
	)
	if err != nil {
		return err
	}

	args, err := helper.To[model.JobArgs](msg.Input)
	if err != nil {
		return err
	}

	err = j.workerService.Launch(ctx, args)
	if err != nil {
		return err
	}

	err = j.broker.Publish(
		ctx,
		"event",
		pkgmodel.JobEventMessage{
			JobID:     msg.JobID,
			Message:   "random worker finished",
			Timestamp: time.Now().UnixMilli(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (j *jobLaunchConsumer) Shutdown(ctx context.Context) error {
	return j.shutdown()
}
