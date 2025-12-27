package main

import (
	"log"
	"time"
)

func main() {
	initDB()
	initKafka()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("Scheduler started")

	for {
		<-ticker.C

		jobs, err := fetchReadyJobs(10)
		if err != nil {
			log.Println("Fetch failed:", err)
			continue
		}

		for _, job := range jobs {
			if err := claimJob(job.JobID); err != nil {
				continue
			}

			publishJob(job.JobID)

		}
	}
}

// This is the main scheduler loop that periodically fetches ready jobs, claims them, and publishes to Kafka.
