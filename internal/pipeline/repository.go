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

func GetAllPipelineRuns() ([]PipelineRun, error) {
	rows, err := storage.DB.Query(context.Background(), `
		SELECT id, commit_hash, branch, status, duration, created_at
		FROM pipeline_runs
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []PipelineRun
	for rows.Next() {
		var run PipelineRun
		if err := rows.Scan(
			&run.ID,
			&run.CommitHash,
			&run.Branch,
			&run.Status,
			&run.Duration,
			&run.CreatedAt,
		); err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, nil
}
