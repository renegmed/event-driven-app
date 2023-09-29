package server

import (
	"job-runner-app/internal/handler"
	//"job-runner-app/docs"

	"github.com/gin-gonic/gin"
)

func NewEngine(
	workerHandler handler.WorkerHandler,
	jobHandler handler.JobHandler,
	healthHandler handler.HealthHandler,
) *gin.Engine {
	e := gin.Default()
	//docs.SwaggerInfo.BasePath = "/"

	w := e.Group("/workers")
	{
		w.GET("/:id", workerHandler.Get)
		w.PUT("/", workerHandler.Create)
		w.DELETE("/:id", workerHandler.Delete)
	}

	j := e.Group("/jobs")
	{
		j.POST("/", jobHandler.Launch)
		j.GET("/:id/statuses", jobHandler.JobStatuses)
	}

	h := e.Group("/health")
	{
		h.GET("/", healthHandler.Health)
	}

	//e.GET("/swagger/*any", ginSwagger.WrapHander(swaggerfiles.Handler))

	return e
}
