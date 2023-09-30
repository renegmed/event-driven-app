package dependency

import (
	"job-worker-app/internal/broker"
	"job-worker-app/internal/config"
	"job-worker-app/internal/consumer"
	"job-worker-app/internal/service"
	pkgmodel "job-worker-app/pkg/model"
	"log"
)

type Dependency struct {
	Config        config.Config
	JobConsumer   consumer.Consumer
	WorkerService service.WorkerService
	Broker        broker.Broker[pkgmodel.JobEventMessage]
}

func NewDependency() (*Dependency, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	log.Println("...dependency.NewDependency() cfg:", cfg)

	workerService := service.NewWorkerService()
	broker := broker.NewJobEventsBroker(cfg.RabbitMQ, cfg.JobEvents)

	jobConsumer := consumer.NewJobLaunchConsumer(cfg.RabbitMQ, cfg.Jobs, workerService, broker)

	return &Dependency{
		Config:        cfg,
		JobConsumer:   jobConsumer,
		WorkerService: workerService,
		Broker:        broker,
	}, nil
}
