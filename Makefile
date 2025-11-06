.PHONY: help build build-producer build-consumer run-producer run-consumer docker-up docker-down clean test

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: build-producer build-consumer ## Build both producer and consumer

build-producer: ## Build the producer
	@echo "Building producer..."
	@cd producer && go build -o producer main.go

build-consumer: ## Build the consumer
	@echo "Building consumer..."
	@cd consumer && go build -o consumer main.go

run-producer: ## Run the producer
	@cd producer && go run main.go

run-consumer: ## Run the consumer
	@cd consumer && go run main.go

docker-up: ## Start RabbitMQ with Docker Compose
	docker compose up -d

docker-down: ## Stop RabbitMQ
	docker compose down

clean: ## Clean build artifacts
	@rm -f producer/producer consumer/consumer
	@echo "Cleaned build artifacts"

test: ## Run tests
	go test -v ./...

deps: ## Download dependencies
	go mod download
	go mod tidy
