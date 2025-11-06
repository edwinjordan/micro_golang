.PHONY: help build build-producer build-consumer build-stock run-producer run-consumer run-stock docker-up docker-down clean test

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: build-producer build-consumer build-stock ## Build all services

build-producer: ## Build the producer
	@echo "Building producer..."
	@cd producer && go build -o producer main.go

build-consumer: ## Build the consumer
	@echo "Building consumer..."
	@cd consumer && go build -o consumer main.go

build-stock: ## Build the stock service
	@echo "Building stock service..."
	@cd stockservice && go build -o stockservice main.go

run-producer: ## Run the producer
	@cd producer && go run main.go

run-consumer: ## Run the consumer
	@cd consumer && go run main.go

run-stock: ## Run the stock service
	@cd stockservice && go run main.go

docker-up: ## Start RabbitMQ with Docker Compose
	docker compose up -d

docker-down: ## Stop RabbitMQ
	docker compose down

clean: ## Clean build artifacts
	@rm -f producer/producer consumer/consumer stockservice/stockservice
	@echo "Cleaned build artifacts"

test: ## Run tests
	go test -v ./...

deps: ## Download dependencies
	go mod download
	go mod tidy
