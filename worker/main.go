package main

import (
	"log"
	"os"

	"github.com/google/uuid"
)

func main() {
	initDB()

	workerID := uuid.New().String()
	hostname, _ := os.Hostname()

	registerWorker(workerID, hostname)
	startHeartbeat(workerID)

	log.Println("Worker started:", workerID)

	// Block forever (worker will later consume Kafka)
	select {}
}
