CREATE TABLE pipeline_runs (
    id SERIAL PRIMARY KEY,
    commit_hash TEXT NOT NULL,
    branch TEXT NOT NULL,
    status TEXT NOT NULL,
    duration INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
