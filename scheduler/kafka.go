package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

func initKafka() {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "jobs.ready",
		Balancer: &kafka.LeastBytes{},
	}
}

func publishJob(jobID string) {
	err := kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(jobID),
			Value: []byte(jobID),
		},
	)

	if err != nil {
		log.Println("Kafka publish failed:", err)
	} else {
		log.Println("Published job to Kafka:", jobID)
	}
}

// This function publishes a job ID to the Kafka topic "jobs.ready".
