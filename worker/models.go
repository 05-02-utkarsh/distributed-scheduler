package main

import "time"

type Worker struct {
	WorkerID      string
	Hostname      string
	Status        string
	LastHeartbeat time.Time
	RegisteredAt  time.Time
}
