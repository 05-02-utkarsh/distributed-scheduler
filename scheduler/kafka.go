package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

type JobMessage struct {
	JobID       string `json:"job_id"`
	ExecutionID string `json:"execution_id"`
}

func initKafka() {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "jobs.ready",
		Balancer: &kafka.LeastBytes{},
	}
	log.Println("Kafka producer initialized")
}

func publishJob(jobID, executionID string) error {
	msg := JobMessage{
		JobID:       jobID,
		ExecutionID: executionID,
	}

	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(jobID),
			Value: value,
		},
	)
}

// This function publishes a job message to the Kafka topic "jobs.ready".
