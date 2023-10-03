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

	pkgmodel "job-runner-app/pkg/model"

	"github.com/gin-gonic/gin"
)

type Dependency struct {
	CFG         config.Config
	JobConsumer consumer.Consumer

	JobBroker broker.Broker[pkgmodel.JobLaunchMessage]

	// -----------------
	JobRepository    repository.JobRepository
	JobService       service.JobService
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

	jobRepository := repository.NewJobRepository(cfg.DB)

	healthRepository := repository.NewHealthRepository(cfg.DB)

	healthService := service.NewHealthService(healthRepository)
	healthHandler := handler.NewHealthHandler(healthService)

	workerRepository := repository.NewWorkerRepository(cfg.DB)
	workerService := service.NewWorkerService(workerRepository)
	workerHandler := handler.NewWorkerHandler(workerService)

	jobBroker := broker.NewJobLaunchBroker(cfg.RabbitMQ, cfg.Jobs)
	jobService := service.NewJobService(workerRepository, jobRepository, jobBroker)
	jobHandler := handler.NewJobHandler(jobService)

	ginEngine := server.NewEngine(workerHandler, jobHandler, healthHandler)

	jobConsumer := consumer.NewJobStatusConsumer(cfg.RabbitMQ, cfg.JobEvents, jobService) // concrete rabbitMQConsumer

	return &Dependency{
		CFG:              cfg,
		JobConsumer:      jobConsumer,
		JobBroker:        jobBroker,
		JobRepository:    jobRepository,
		JobService:       jobService,
		HealthRepository: healthRepository,
		GinEngine:        ginEngine,
	}, nil
}
