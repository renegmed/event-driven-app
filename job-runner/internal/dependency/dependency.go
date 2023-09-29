package dependency

import (
	"job-runner-app/internal/broker"
	"job-runner-app/internal/config"
	"job-runner-app/internal/consumer"
	"job-runner-app/internal/handler"
	"job-runner-app/internal/repository"
	"job-runner-app/internal/server"
	"job-runner-app/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

type Dependency struct {
	CFG              config.Config
	JobConsumer      consumer.Consumer
	JobRepository    repository.JobRepository
	HealthRepository repository.HealthRepository
	HealthService    service.HealthService
	// gin handlers
	HealthHandler handler.HealthHandler
	GinEngine     *gin.Engine
}

func NewDependency() (*Dependency, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	log.Println("...NewDependency() cfg:", cfg)

	jobConsumer := consumer.NewJobStatusConsumer(cfg.RabbitMQ) // concrete rabbitMQConsumer

	jobRepository := repository.NewJobRepository(cfg.DB)

	healthRepository := repository.NewHealthRepository(cfg.DB)

	healthService := service.NewHealthService(healthRepository)
	healthHandler := handler.NewHealthHandler(healthService)

	workerRepository := repository.NewWorkerRepository(cfg.DB)
	workerService := service.NewWorkerService(workerRepository)
	workerHandler := handler.NewWorkerHandler(workerService)

	broker := broker.NewJobLaunchBroker(cfg.RabbitMQ, cfg.Jobs)
	jobService := service.NewJobService(workerRepository, jobRepository, broker)
	jobHandler := handler.NewJobHandler(jobService)

	ginEngine := server.NewEngine(workerHandler, jobHandler, healthHandler)

	return &Dependency{
		CFG:              cfg,
		JobConsumer:      jobConsumer,
		JobRepository:    jobRepository,
		HealthRepository: healthRepository,
		GinEngine:        ginEngine,
	}, nil
}
