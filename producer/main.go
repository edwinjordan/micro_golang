package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edwinjordan/micro_golang/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
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
	_, err = conn.DeclareQueue(queueName)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start producing messages
	go produceMessages(ctx, conn)

	// Wait for termination signal
	<-sigChan
	log.Println("Shutting down producer...")
	cancel()
	time.Sleep(time.Second) // Give some time for graceful shutdown
}

func produceMessages(ctx context.Context, conn *rabbitmq.Connection) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	messageCount := 0

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping message production")
			return
		case <-ticker.C:
			messageCount++
			message := map[string]interface{}{
				"id":        messageCount,
				"message":   "Hello from microservice!",
				"timestamp": time.Now().Format(time.RFC3339),
			}

			body := formatMessage(message)

			err := conn.Channel.PublishWithContext(
				ctx,
				"",        // exchange
				queueName, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         []byte(body),
					Timestamp:    time.Now(),
				},
			)

			if err != nil {
				log.Printf("Failed to publish message: %v", err)
			} else {
				log.Printf("Sent message #%d: %s", messageCount, body)
			}
		}
	}
}

func formatMessage(data map[string]interface{}) string {
	return fmt.Sprintf("ID: %v, Message: %v, Timestamp: %v",
		data["id"], data["message"], data["timestamp"])
}
