package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/edwinjordan/micro_golang/rabbitmq"
)

const (
	defaultRabbitMQURL = "amqp://guest:guest@localhost:5672/"
	queueName          = "microservice_queue"
)

func main() {
	// Get RabbitMQ URL from environment or use default
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = defaultRabbitMQURL
	}

	// Connect to RabbitMQ
	conn, err := rabbitmq.NewConnection(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Declare queue
	q, err := conn.DeclareQueue(queueName)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Start consuming messages
	msgs, err := conn.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Printf("Consumer started. Waiting for messages on queue: %s", queueName)

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	go func() {
		for msg := range msgs {
			log.Printf("Received message: %s", string(msg.Body))

			// Acknowledge the message
			if err := msg.Ack(false); err != nil {
				log.Printf("Failed to acknowledge message: %v", err)
			} else {
				log.Println("Message acknowledged")
			}
		}
	}()

	// Wait for termination signal
	<-sigChan
	log.Println("Shutting down consumer...")
}
