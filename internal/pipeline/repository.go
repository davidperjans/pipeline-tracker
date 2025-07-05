package pipeline

import (
	"context"

	"github.com/davidperjans/pipeline-tracker/internal/storage"
)

func InsertPipelineRun(run PipelineRun) error {
	query := `
		INSERT INTO pipeline_runs (commit_hash, branch, status, duration)
		VALUES ($1, $2, $3, $4)
	`

	_, err := storage.DB.Exec(context.Background(), query,
		run.CommitHash,
		run.Branch,
		run.Status,
		run.Duration,
	)

	return err
}
