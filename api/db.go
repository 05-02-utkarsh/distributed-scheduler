package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func initDB() {
	connStr := "postgres://scheduler:scheduler@localhost:5432/scheduler_db"

	var err error
	db, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to connect to DB:", err)
	}

	log.Println("Connected to PostgreSQL")
}

// This connects Go → Docker PostgreSQL.

func updateJobStatus(jobID string, from, to JobStatus) error {
	if err := validateTransition(from, to); err != nil {
		return err
	}

	cmd, err := db.Exec(context.Background(), `
		UPDATE jobs
		SET status = $1, updated_at = NOW()
		WHERE job_id = $2 AND status = $3
	`, to, jobID, from)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("job status update failed (possible race condition)")
	}

	return nil
}

// This function updates the job status in the database after validating the state transition.
