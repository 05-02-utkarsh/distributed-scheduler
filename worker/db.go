package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func initDB() {
	connStr := "postgres://scheduler:scheduler@localhost:5432/scheduler_db"

	var err error
	db, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
}

func registerWorker(workerID, hostname string) {
	_, err := db.Exec(context.Background(), `
		INSERT INTO workers (
			worker_id, hostname, status, last_heartbeat
		)
		VALUES ($1, $2, 'ALIVE', NOW())
	`,
		workerID,
		hostname,
	)

	if err != nil {
		log.Fatal("Worker registration failed:", err)
	}

	log.Println("Worker registered:", workerID)
}

func updateHeartbeat(workerID string) {
	_, err := db.Exec(context.Background(), `
		UPDATE workers
		SET last_heartbeat = NOW()
		WHERE worker_id = $1
	`, workerID)

	if err != nil {
		log.Println("Heartbeat update failed:", err)
	}
}
