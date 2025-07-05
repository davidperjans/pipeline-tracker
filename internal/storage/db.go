package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectToDB() error {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		getEnv("DB_USER", "tracker"),
		getEnv("DB_PASSWORD", "secret"),
		getEnv("DB_PORT", "5433"),
		getEnv("DB_NAME", "pipeline_db"),
	)

	var err error
	DB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
