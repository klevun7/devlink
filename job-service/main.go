package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// define job json struct
type Job struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Company  string   `json:"company"`
	Tags     []string `json:"tags"`
	Location string   `json:"location"`
}

func main() {
	// kafka writer for job_posted 
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9093"},
		Topic:    "job_posted",
		Balancer: &kafka.LeastBytes{},
	
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	})
	
	// Ensure writer is closed when function exits
	defer writer.Close()

	// mock job 
	job := Job{
		ID:       "job-001",
		Title:    "Backend Engineer",
		Company:  "Stripe",
		Tags:     []string{"go", "kafka", "SQL"},
		Location: "San Francisco",
	}

	// marshall json and catch error
	data, err := json.Marshal(job)
	if err != nil {
		log.Fatal("error marshalling job:", err)
	}

	// Add context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// write job message to kafka with context param; key: jobID, value: JSON data
	err = writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(job.ID),
		Value: data,
	})

	if err != nil {
		log.Fatal("failed to write message:", err)
	}

	fmt.Println("Job sent to kafka topic!")
}