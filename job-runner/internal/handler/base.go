package handler

import (
	"errors"
	"job-runner-app/internal/exception"
	"job-runner-app/pkg/model"
	"net/http"
)

func toErrorResponse(err error) (int, any) {
	var e exception.JobError
	if ok := errors.As(err, &e); ok {
		statusCode := int(e.Code) / 100
		return statusCode, model.ErrorResponse{
			Code:    int(e.Code),
			Message: e.Message,
		}
	}

	return http.StatusInternalServerError, model.ErrorResponse{
		Code:    int(exception.UnknownError),
		Message: "unknown error",
	}
}
