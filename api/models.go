package main

import "time"

type Job struct {
	JobID       string    `json:"job_id"`
	JobType     string    `json:"job_type"`
	Payload     any       `json:"payload"`
	Status      string    `json:"status"`
	RetryCount  int       `json:"retry_count"`
	MaxRetries  int       `json:"max_retries"`
	NextRunTime time.Time `json:"next_run_time"`
	CreatedAt   time.Time `json:"created_at"`
}
