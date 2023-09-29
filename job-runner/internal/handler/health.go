package handler

import (
	"job-runner-app/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler interface {
	Health(ctx *gin.Context)
}

type healthHandler struct {
	healthService service.HealthService
}

func NewHealthHandler(s service.HealthService) HealthHandler {
	return &healthHandler{healthService: s}
}

func (h *healthHandler) Health(ctx *gin.Context) {
	if err := h.healthService.Health(ctx); err != nil {
		log.Printf("failed to check health, error: %v", err)
		ctx.JSON(toErrorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
