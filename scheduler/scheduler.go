package main

import (
	"context"
)

type Job struct {
	JobID string
}

func fetchReadyJobs(limit int) ([]Job, error) {
	rows, err := db.Query(context.Background(), `
		SELECT job_id
		FROM jobs
		WHERE status = 'SUBMITTED'
		  AND next_run_time <= NOW()
		ORDER BY created_at
		LIMIT $1
	`, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var j Job
		if err := rows.Scan(&j.JobID); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

// This function fetches jobs that are ready to be processed from the database.
