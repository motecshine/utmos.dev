.PHONY: help build test lint clean run-api run-ws run-uplink run-downlink run-gateway

# Variables
GO := go
GOFMT := gofmt
GOLANGCI_LINT := golangci-lint
SERVICES := iot-api iot-ws iot-uplink iot-downlink iot-gateway

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all services
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		$(GO) build -o bin/$$service ./cmd/$$service; \
	done

test: ## Run all tests
	@echo "Running tests..."
	$(GO) test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	$(GOLANGCI_LINT) run

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...

tidy: ## Run go mod tidy
	@echo "Running go mod tidy..."
	$(GO) mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf coverage.out coverage.html
	$(GO) clean -cache

run-api: ## Run iot-api service
	$(GO) run ./cmd/iot-api

run-ws: ## Run iot-ws service
	$(GO) run ./cmd/iot-ws

run-uplink: ## Run iot-uplink service
	$(GO) run ./cmd/iot-uplink

run-downlink: ## Run iot-downlink service
	$(GO) run ./cmd/iot-downlink

run-gateway: ## Run iot-gateway service
	$(GO) run ./cmd/iot-gateway

docker-build: ## Build Docker images for all services
	@echo "Building Docker images..."
	@for service in $(SERVICES); do \
		echo "Building Docker image for $$service..."; \
		docker build -f deployments/docker/$$service.Dockerfile -t umos/$$service:latest .; \
	done

docker-compose-up: ## Start services with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop services with docker-compose
	docker-compose down

