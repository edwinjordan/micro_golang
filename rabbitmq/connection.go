package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection wraps RabbitMQ connection
type Connection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewConnection creates a new RabbitMQ connection
func NewConnection(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	log.Println("Successfully connected to RabbitMQ")

	return &Connection{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// Close closes the connection and channel
func (c *Connection) Close() {
	if c.Channel != nil {
		c.Channel.Close()
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
	log.Println("RabbitMQ connection closed")
}

// DeclareQueue declares a queue
func (c *Connection) DeclareQueue(queueName string) (amqp.Queue, error) {
	q, err := c.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare a queue: %w", err)
	}

	log.Printf("Queue declared: %s", queueName)
	return q, nil
}
