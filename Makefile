.PHONY: help build test lint clean run-api run-ws run-uplink run-downlink run-gateway coverage

# Variables
GO := go
GOFMT := gofmt
# Try to find golangci-lint in PATH, otherwise fallback to user's home directory
GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null || echo $(HOME)/golangci-lint/golangci-lint)
SERVICES := iot-api iot-ws iot-uplink iot-downlink iot-gateway
COVERAGE_THRESHOLD := 80

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

test-short: ## Run tests in short mode (skip integration tests)
	@echo "Running short tests..."
	$(GO) test -v -short -race -coverprofile=coverage.out ./...

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	$(GO) test -v -race ./tests/integration/...

test-coverage: test ## Run tests with coverage report
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

coverage: ## Verify test coverage meets threshold (TDD-002: >= 80%)
	@echo "Running tests with coverage..."
	@$(GO) test -coverprofile=coverage.out ./... > /dev/null 2>&1 || true
	@echo ""
	@echo "=== Coverage Report ==="
	@$(GO) tool cover -func=coverage.out | tail -1
	@echo ""
	@COVERAGE=$$($(GO) tool cover -func=coverage.out | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	echo "Coverage: $$COVERAGE%"; \
	echo "Threshold: $(COVERAGE_THRESHOLD)%"; \
	if [ $$(echo "$$COVERAGE >= $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "✓ Coverage meets threshold"; \
	else \
		echo "✗ Coverage below threshold"; \
		exit 1; \
	fi

coverage-report: ## Generate detailed coverage report
	@echo "Generating coverage report..."
	@$(GO) test -coverprofile=coverage.out ./... > /dev/null 2>&1 || true
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo ""
	@echo "=== Package Coverage ==="
	@$(GO) tool cover -func=coverage.out | grep -E "^github.com/utmos/utmos/(internal|pkg)" | head -30
	@echo ""
	@echo "Coverage report generated: coverage.html"

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./tests/integration/...

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

