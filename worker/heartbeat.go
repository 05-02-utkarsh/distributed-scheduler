package main

import (
	"time"
)

func startHeartbeat(workerID string) {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			<-ticker.C
			updateHeartbeat(workerID)
		}
	}()
}
