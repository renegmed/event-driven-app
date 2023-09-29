package model

type JobEventMessage struct {
	JobID     string         `json:"job_id"`
	Message   string         `json:"message"`
	Timestamp int64          `json:"timestamp"`
	Output    map[string]any `json:"output"`
}
