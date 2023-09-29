package consumer

import (
	"context"
	"job-worker-app/internal/config"
	pkgmodel "job-worker-app/pkg/model"
	"log"
)

type jobLaunchConsumer struct {
	*rabbitMQConsumer[pkgmodel.JobLaunchMessage] // generic type JobLaunchMessage
}

func NewJobLaunchConsumer(cc config.SubscriberConfig) Consumer {
	log.Println("...NewJobLaunchConsumer() config.SubscriberConfig Group:", cc.Group, " Queue:", cc.Queue)
	return &jobLaunchConsumer{}
}

func (j *jobLaunchConsumer) Consume() error {
	return j.consume(j.consumeCallback)
}

func (j *jobLaunchConsumer) consumeCallback(ctx context.Context, msg pkgmodel.JobLaunchMessage, err error) error {

	if err != nil {
		log.Printf("job launch message error: %v", err)

		return nil
	}

	log.Println("...consumeCallback() msg:", msg)

	return nil
}
