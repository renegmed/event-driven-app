package model

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type WorkerAPI struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateWorkerRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=255"`
	Description string `json:"description" binding:"required,min=3,max=255"`
}

type CreateWorkerResponse struct {
	ID string `json:"id"`
}

type DeleteWorkerResponse struct {
	Delete bool `json:"deleted"`
}

type GetWorkerResponse struct {
	Worker WorkerAPI `json:"worker"`
}