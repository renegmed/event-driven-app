package dependency

import (
	"job-worker-app/internal/config"
	"job-worker-app/internal/consumer"
)

type Dependency struct {
	JobConsumer consumer.Consumer
}

func NewDependency() (*Dependency, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	jobConsumer := consumer.NewJobLaunchConsumer(cfg.Jobs)
	return &Dependency{JobConsumer: jobConsumer}, nil
}
