# micro_golang

A collection of simple microservices in Go that demonstrate basic patterns for building distributed systems. This project includes RabbitMQ-based messaging services and a REST API-based stock management service.

## Features

- **Producer Service**: Sends messages to RabbitMQ queue every 5 seconds
- **Consumer Service**: Consumes and processes messages from RabbitMQ queue
- **Stock Management Service**: RESTful API for managing product inventory
- **RabbitMQ Integration**: Uses RabbitMQ as message broker with persistent messages
- **Docker Support**: Easy setup with Docker Compose
- **Graceful Shutdown**: All services handle SIGINT/SIGTERM signals properly

## Architecture

```
Producer → RabbitMQ Queue → Consumer

Stock Management API (HTTP REST)
```

The producer sends timestamped messages to a RabbitMQ queue, and the consumer reads and acknowledges these messages. The stock management service provides a REST API for managing product inventory independently.

## Prerequisites

- Go 1.24+ installed
- Docker and Docker Compose (for running RabbitMQ)
- Make (optional, for using Makefile commands)

## Quick Start

### RabbitMQ Services

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

### Stock Management Service

The stock management service is a standalone REST API that doesn't require RabbitMQ.

Run the stock service:
```bash
go run stockservice/main.go
```

Or using Make:
```bash
make run-stock
```

The service will start on port 8080 by default. You can change the port using the `PORT` environment variable:
```bash
PORT=3000 go run stockservice/main.go
```

### Stock Management API Usage

The stock service provides the following endpoints:

**Get all products:**
```bash
curl http://localhost:8080/api/stock
```

**Get a specific product:**
```bash
curl http://localhost:8080/api/stock/LAPTOP001
```

**Add a new product:**
```bash
curl -X POST http://localhost:8080/api/stock \
  -H "Content-Type: application/json" \
  -d '{
    "id": "MONITOR001",
    "name": "4K Monitor",
    "description": "32 inch 4K display",
    "quantity": 25,
    "price": 499.99
  }'
```

**Update stock quantity:**
```bash
# Add stock
curl -X POST http://localhost:8080/api/stock/LAPTOP001/update \
  -H "Content-Type: application/json" \
  -d '{"quantity": 10, "operation": "add"}'

# Remove stock
curl -X POST http://localhost:8080/api/stock/LAPTOP001/update \
  -H "Content-Type: application/json" \
  -d '{"quantity": 5, "operation": "remove"}'

# Set exact quantity
curl -X POST http://localhost:8080/api/stock/LAPTOP001/update \
  -H "Content-Type: application/json" \
  -d '{"quantity": 100, "operation": "set"}'
```

**Delete a product:**
```bash
curl -X DELETE http://localhost:8080/api/stock/LAPTOP001
```

**Health check:**
```bash
curl http://localhost:8080/health
```

## Building

Build all services:
```bash
make build
```

Build individually:
```bash
make build-producer
make build-consumer
make build-stock
```

Run the built binaries:
```bash
./producer/producer
./consumer/consumer
./stockservice/stockservice
```

## Configuration

### RabbitMQ Services

Both producer and consumer services can be configured using environment variables:

- `RABBITMQ_URL`: RabbitMQ connection URL (default: `amqp://guest:guest@localhost:5672/`)

Example:
```bash
RABBITMQ_URL="amqp://user:password@rabbitmq-host:5672/" go run producer/main.go
```

### Stock Management Service

The stock service can be configured using environment variables:

- `PORT`: HTTP port for the service (default: `8080`)

Example:
```bash
PORT=3000 go run stockservice/main.go
```

## Project Structure

```
.
├── producer/           # Producer service (RabbitMQ)
│   └── main.go
├── consumer/           # Consumer service (RabbitMQ)
│   └── main.go
├── stockservice/       # Stock management service (HTTP REST API)
│   └── main.go
├── stock/              # Stock management package
│   ├── types.go        # Data structures
│   ├── store.go        # Business logic
│   ├── handler.go      # HTTP handlers
│   └── store_test.go   # Tests
├── rabbitmq/           # Shared RabbitMQ connection package
│   ├── connection.go
│   └── connection_test.go
├── docker-compose.yml  # RabbitMQ setup
├── Makefile           # Build and run commands
├── go.mod             # Go module definition
└── README.md          # This file
```

## Available Make Commands

- `make help` - Display available commands
- `make build` - Build all services
- `make build-producer` - Build only the producer
- `make build-consumer` - Build only the consumer
- `make build-stock` - Build only the stock service
- `make run-producer` - Run the producer
- `make run-consumer` - Run the consumer
- `make run-stock` - Run the stock service
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

### RabbitMQ Services

1. **Connection**: Both services connect to RabbitMQ using the shared `rabbitmq` package
2. **Queue Declaration**: A durable queue named `microservice_queue` is declared
3. **Producer**: Sends a message every 5 seconds with an ID, message text, and timestamp
4. **Consumer**: Listens for messages, logs them, and acknowledges receipt
5. **Graceful Shutdown**: Both services listen for interrupt signals and close connections cleanly

### Stock Management Service

1. **In-Memory Store**: Products are stored in memory with thread-safe access using sync.RWMutex
2. **REST API**: Provides HTTP endpoints for CRUD operations on products
3. **Stock Operations**: Supports add, remove, and set operations for stock quantities
4. **Sample Data**: Initializes with 3 sample products on startup
5. **Validation**: Includes error handling for insufficient stock, invalid quantities, and missing products
6. **Graceful Shutdown**: Handles shutdown signals with proper cleanup

## License

MIT
