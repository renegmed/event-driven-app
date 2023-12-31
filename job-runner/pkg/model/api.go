package model

import "time"

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type WorkerAPI struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type JobStatusAPI struct {
	Message   string         `json:"message"`
	Timestamp time.Time      `json:"timestamp"`
	Output    map[string]any `json:"output"`
}

type CreateWorkerRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=255"`
	Description string `json:"description" binding:"required,min=3,max=255"`
}

type CreateWorkerResponse struct {
	ID string `json:"id"`
}

type DeleteWorkerResponse struct {
	Deleted bool `json:"deleted"`
}

type GetWorkerResponse struct {
	Worker WorkerAPI `json:"worker"`
}

type LaunchJobRequest struct {
	WorkerID string         `json:"worker_id" binding:"required,min=3,max=255"`
	Input    map[string]any `json:"input" binding:"required"`
}

type LaunchJobResponse struct {
	JobID string `json:"job_id"`
}

type JobStatusesResponse struct {
	Statuses []JobStatusAPI `json:"statuses"`
}
