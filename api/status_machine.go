package main

import "fmt"

type JobStatus string

const (
	StatusSubmitted JobStatus = "SUBMITTED"
	StatusQueued    JobStatus = "QUEUED"
	StatusRunning   JobStatus = "RUNNING"
	StatusSuccess   JobStatus = "SUCCESS"
	StatusFailed    JobStatus = "FAILED"
	StatusDead      JobStatus = "DEAD"
)

var validTransitions = map[JobStatus][]JobStatus{
	StatusSubmitted: {StatusQueued},
	StatusQueued:    {StatusRunning},
	StatusRunning:   {StatusSuccess, StatusFailed},
	StatusFailed:    {StatusQueued, StatusDead},
}

func isValidTransition(from, to JobStatus) bool {
	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

func validateTransition(from, to JobStatus) error {
	if !isValidTransition(from, to) {
		return fmt.Errorf("invalid job state transition: %s → %s", from, to)
	}
	return nil
}

// This file defines the state machine for job status transitions.
