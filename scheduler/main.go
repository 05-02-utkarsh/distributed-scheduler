package main

import (
	"log"
	"time"
)

func main() {
	// 1️⃣ Init dependencies
	initDB()
	initKafka()
	initEtcd()            // 🆕 etcd client
	startLeaderElection() // 🆕 leader election goroutine

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("Scheduler started")

	for {
		<-ticker.C

		// 🔒 ONLY LEADER RUNS THE SCHEDULER LOOP
		if !isLeader {
			continue
		}

		jobs, err := fetchReadyJobs(10)
		if err != nil {
			log.Println("Failed to fetch jobs:", err)
			continue
		}

		for _, job := range jobs {

			// 1️⃣ Try to claim the job (SUBMITTED → QUEUED)
			if err := claimJob(job.JobID); err != nil {
				// Another scheduler probably claimed it
				continue
			}

			// 2️⃣ Create a new execution record
			execID, err := createExecution(job.JobID)
			if err != nil {
				log.Println("Failed to create execution for job", job.JobID, ":", err)
				continue
			}

			// 3️⃣ Publish job + execution to Kafka
			if err := publishJob(job.JobID, execID); err != nil {
				log.Println("Failed to publish to Kafka for job", job.JobID, ":", err)
				continue
			}

			log.Println("Scheduled job:", job.JobID, "execution:", execID)
		}
	}
}
