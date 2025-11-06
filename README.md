# micro_golang

A simple microservice in Go that communicates with RabbitMQ. This project demonstrates the basic patterns for building producer and consumer services using RabbitMQ as a message broker.

## Features

- **Producer Service**: Sends messages to RabbitMQ queue every 5 seconds
- **Consumer Service**: Consumes and processes messages from RabbitMQ queue
- **RabbitMQ Integration**: Uses RabbitMQ as message broker with persistent messages
- **Docker Support**: Easy setup with Docker Compose
- **Graceful Shutdown**: Both services handle SIGINT/SIGTERM signals properly

## Architecture

```
Producer → RabbitMQ Queue → Consumer
```

The producer sends timestamped messages to a RabbitMQ queue, and the consumer reads and acknowledges these messages.

## Prerequisites

- Go 1.24+ installed
- Docker and Docker Compose (for running RabbitMQ)
- Make (optional, for using Makefile commands)

## Quick Start

### 1. Start RabbitMQ

```bash
docker compose up -d
```

(Or `docker-compose up -d` if using the older Docker Compose v1)

This starts RabbitMQ with:
- AMQP port: 5672
- Management UI: http://localhost:15672 (username: guest, password: guest)

### 2. Install Dependencies

```bash
go mod download
```

Or using Make:
```bash
make deps
```

### 3. Run the Consumer

In one terminal:
```bash
go run consumer/main.go
```

Or using Make:
```bash
make run-consumer
```

### 4. Run the Producer

In another terminal:
```bash
go run producer/main.go
```

Or using Make:
```bash
make run-producer
```

## Building

Build both services:
```bash
make build
```

Build individually:
```bash
make build-producer
make build-consumer
```

Run the built binaries:
```bash
./producer/producer
./consumer/consumer
```

## Configuration

Both services can be configured using environment variables:

- `RABBITMQ_URL`: RabbitMQ connection URL (default: `amqp://guest:guest@localhost:5672/`)

Example:
```bash
RABBITMQ_URL="amqp://user:password@rabbitmq-host:5672/" go run producer/main.go
```

## Project Structure

```
.
├── producer/           # Producer service
│   └── main.go
├── consumer/           # Consumer service
│   └── main.go
├── rabbitmq/           # Shared RabbitMQ connection package
│   └── connection.go
├── docker-compose.yml  # RabbitMQ setup
├── Makefile           # Build and run commands
├── go.mod             # Go module definition
└── README.md          # This file
```

## Available Make Commands

- `make help` - Display available commands
- `make build` - Build both producer and consumer
- `make build-producer` - Build only the producer
- `make build-consumer` - Build only the consumer
- `make run-producer` - Run the producer
- `make run-consumer` - Run the consumer
- `make docker-up` - Start RabbitMQ with Docker Compose
- `make docker-down` - Stop RabbitMQ
- `make clean` - Clean build artifacts
- `make test` - Run tests
- `make deps` - Download and tidy dependencies

## RabbitMQ Management UI

Access the RabbitMQ management interface at http://localhost:15672
- Username: `guest`
- Password: `guest`

From the UI, you can:
- Monitor queues and messages
- View connections and channels
- Check message rates

## Stopping Services

Stop the producer/consumer with `Ctrl+C` (they will shut down gracefully).

Stop RabbitMQ:
```bash
docker compose down
```

(Or `docker-compose down` if using the older Docker Compose v1)

Or using Make:
```bash
make docker-down
```

## How It Works

1. **Connection**: Both services connect to RabbitMQ using the shared `rabbitmq` package
2. **Queue Declaration**: A durable queue named `microservice_queue` is declared
3. **Producer**: Sends a message every 5 seconds with an ID, message text, and timestamp
4. **Consumer**: Listens for messages, logs them, and acknowledges receipt
5. **Graceful Shutdown**: Both services listen for interrupt signals and close connections cleanly

## License

MIT
