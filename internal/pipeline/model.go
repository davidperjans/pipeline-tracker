package pipeline

import "time"

type PipelineRun struct {
	ID         string    `json:"id"`
	CommitHash string    `json:"commit_hash"`
	Branch     string    `json:"branch"`
	Status     string    `json:"status"`
	Duration   int       `json:"duration"` // in seconds
	CreatedAt  time.Time `json:"created_at,omitempty"`
}
