package handler

import (
	"fmt"
	"job-runner-app/internal/exception"
	"job-runner-app/internal/service"
	"job-runner-app/pkg/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WorkerHandler interface {
	Get(*gin.Context)
	Create(*gin.Context)
	Delete(*gin.Context)
}

type workerHandler struct {
	workerService service.WorkerService
}

func NewWorkerHandler(ws service.WorkerService) WorkerHandler {
	return &workerHandler{workerService: ws}
}

func (h *workerHandler) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	w, err := h.workerService.Get(ctx, id)
	if err != nil {
		log.Printf("handler.worker.Get() failed to get worker. error: %v", err)
		ctx.JSON(toErrorResponse(err))
		return
	}
	if w == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, model.GetWorkerResponse{
		Worker: model.WorkerAPI{
			ID:          w.ID,
			Name:        w.Name,
			Description: w.Description,
		},
	})
}

func (h *workerHandler) Create(ctx *gin.Context) {
	var req model.CreateWorkerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("handler.worker.Create() failed to create worker. error: %v", err)
		ctx.JSON(toErrorResponse(exception.JobError{
			Code:    exception.InvalidRequest,
			Message: fmt.Sprintf("invalid request. error: %v", err),
			Err:     err,
		}))
		return
	}

	log.Println("...handler.worker.Create() req name:", req.Name, " description:", req.Description)

	id, err := h.workerService.Create(ctx, req.Name, req.Description)
	if err != nil {
		log.Printf("handler.Worker.Create() failed to create worker. error: %v", err)
		ctx.JSON(toErrorResponse(err))
		return
	}

	log.Println("...handler.worker.Create() id:", id)

	ctx.JSON(http.StatusCreated, model.CreateWorkerResponse{ID: id})
}

func (h *workerHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	log.Println("...handler.worker.Delete() id:", id)

	ok, err := h.workerService.Delete(ctx, id)
	if err != nil {
		log.Printf("handler.Worker.Delete() failed to delete worker. error: %v", err)
		ctx.JSON(toErrorResponse(exception.JobError{
			Code:    exception.InvalidRequest,
			Message: fmt.Sprintf("invalid request. error: %v", err),
			Err:     err,
		}))
		return
	}

	ctx.JSON(http.StatusOK, model.DeleteWorkerResponse{
		Deleted: ok,
	})
}
