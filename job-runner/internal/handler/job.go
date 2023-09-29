package handler

import (
	"fmt"
	"job-runner-app/internal/exception"
	"job-runner-app/internal/model"
	"job-runner-app/internal/service"
	pkgmodel "job-runner-app/pkg/model"
	"log"
	"net/http"

	pie "github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

type JobHandler interface {
	Launch(*gin.Context)
	JobStatuses(*gin.Context)
}

type jobHandler struct {
	jobService service.JobService
}

func NewJobHandler(s service.JobService) JobHandler {
	return &jobHandler{jobService: s}
}

func (h *jobHandler) Launch(ctx *gin.Context) {
	var req pkgmodel.LaunchJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("failed to bind request. error: %v", err)
		ctx.JSON(toErrorResponse(exception.JobError{
			Code:    exception.InvalidRequest,
			Message: fmt.Sprintf("invalid request. error: %v", err),
			Err:     err,
		}))
	}

	log.Printf("...handler.job.Launch() req.WorkerID: %s  req.Input: %v", req.WorkerID, req.Input)
	jobID, err := h.jobService.LaunchJob(ctx, req.WorkerID, req.Input)
	if err != nil {
		log.Printf("...handler.job.Launch() failed to launch job. error: %v", err)
		ctx.JSON(toErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, pkgmodel.LaunchJobResponse{
		JobID: jobID,
	})

}
func (h *jobHandler) JobStatuses(ctx *gin.Context) {
	jobID := ctx.Param("id")

	log.Println("...handler.job.JobStatuses() jobID:", jobID)

	statuses, err := h.jobService.GetJobStatuses(ctx, jobID)
	if err != nil {
		log.Printf("...handler.job.JobStatuses() failed to get job statuses. error: %v", err)
		ctx.JSON(toErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, pkgmodel.JobStatusesResponse{
		Statuses: pie.Map(statuses, func(s model.JobStatus) pkgmodel.JobStatusAPI {
			return pkgmodel.JobStatusAPI{
				Message:   s.Message,
				Timestamp: s.Timestamp,
				Output:    s.Output,
			}
		}),
	})
}
