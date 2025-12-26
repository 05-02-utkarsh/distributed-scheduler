package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateJobRequest struct {
	JobType    string `json:"job_type" binding:"required"`
	Payload    any    `json:"payload" binding:"required"`
	MaxRetries int    `json:"max_retries"`
}

func createJob(c *gin.Context) {
	var req CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobID := uuid.New().String()

	_, err := db.Exec(context.Background(), `
		INSERT INTO jobs (
			job_id, job_type, payload, status, max_retries, next_run_time
		)
		VALUES ($1, $2, $3, 'SUBMITTED', $4, $5)
	`,
		jobID,
		req.JobType,
		req.Payload,
		req.MaxRetries,
		time.Now(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"job_id": jobID,
	})
}

func getJob(c *gin.Context) {
	jobID := c.Param("id")

	row := db.QueryRow(context.Background(), `
		SELECT job_id, job_type, payload, status, retry_count, max_retries, next_run_time, created_at
		FROM jobs
		WHERE job_id = $1
	`, jobID)

	var job Job
	err := row.Scan(
		&job.JobID,
		&job.JobType,
		&job.Payload,
		&job.Status,
		&job.RetryCount,
		&job.MaxRetries,
		&job.NextRunTime,
		&job.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}
