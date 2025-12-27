-- =========================
-- Enable UUID generation
-- =========================
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =========================
-- ENUM types
-- =========================
CREATE TYPE job_status AS ENUM (
    'SUBMITTED',
    'QUEUED',
    'RUNNING',
    'SUCCESS',
    'FAILED',
    'DEAD'
);

CREATE TYPE worker_status AS ENUM (
    'ALIVE',
    'DEAD'
);

-- =========================
-- Jobs table
-- =========================
CREATE TABLE jobs (
    job_id UUID PRIMARY KEY,
    job_type TEXT NOT NULL,
    payload JSONB NOT NULL,

    status job_status NOT NULL,
    retry_count INT NOT NULL DEFAULT 0,
    max_retries INT NOT NULL DEFAULT 3,

    next_run_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    idempotency_key TEXT UNIQUE
);

-- =========================
-- Job executions table
-- =========================
CREATE TABLE job_executions (
    execution_id UUID PRIMARY KEY,
    job_id UUID NOT NULL,
    worker_id UUID NOT NULL,

    status job_status NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ,
    error_message TEXT,

    CONSTRAINT fk_job
        FOREIGN KEY (job_id)
        REFERENCES jobs(job_id)
        ON DELETE CASCADE
);

-- =========================
-- Workers table
-- =========================
CREATE TABLE workers (
    worker_id UUID PRIMARY KEY,
    hostname TEXT,

    status worker_status NOT NULL,
    last_heartbeat TIMESTAMPTZ NOT NULL,
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =========================
-- Indexes
-- =========================
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_next_run_time ON jobs(next_run_time);

CREATE INDEX idx_job_exec_job_id ON job_executions(job_id);
CREATE INDEX idx_job_exec_status ON job_executions(status);

CREATE INDEX idx_workers_status ON workers(status);

