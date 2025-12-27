package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

// Initialize DB connection
func initDB() {
	connStr := "postgres://scheduler:scheduler@localhost:5432/scheduler_db"

	var err error
	db, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	log.Println("Scheduler connected to PostgreSQL")
}

// claimJob atomically moves a job from SUBMITTED → QUEUED
// Only ONE scheduler instance can succeed
func claimJob(jobID string) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	cmd, err := tx.Exec(context.Background(), `
		UPDATE jobs
		SET status = 'QUEUED', updated_at = NOW()
		WHERE job_id = $1
		  AND status = 'SUBMITTED'
	`, jobID)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("job already claimed or invalid state")
	}

	return tx.Commit(context.Background())
}

// This function claims a job by updating its status from SUBMITTED to QUEUED atomically.
