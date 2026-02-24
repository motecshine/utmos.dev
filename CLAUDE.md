# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
# Build all services
make build

# Run tests with race detection and coverage
make test

# Run linting (golangci-lint with 70+ linters)
make lint

# Format code
make fmt

# Start local infrastructure (PostgreSQL, InfluxDB, RabbitMQ, VerneMQ, observability stack)
docker-compose up -d

# Run individual services
make run-api        # HTTP API service (iot-api)
make run-ws         # WebSocket service (iot-ws)
make run-uplink     # Uplink message processor (iot-uplink)
make run-downlink   # Downlink message processor (iot-downlink)
make run-gateway    # MQTT gateway (iot-gateway)
```

## Architecture Overview

UMOS is a microservices-based IoT platform with 5 core services communicating via RabbitMQ:

```
Devices ←→ VerneMQ ←→ iot-gateway ←→ RabbitMQ ←→ [iot-uplink, iot-downlink, iot-api, iot-ws] ←→ Clients
```

### Service Responsibilities

| Service | Role | Data Access |
|---------|------|-------------|
| **iot-gateway** | ONLY service connecting to VerneMQ; MQTT↔RabbitMQ bridge | PostgreSQL (device auth) |
| **iot-api** | HTTP RESTful API | PostgreSQL, InfluxDB |
| **iot-ws** | WebSocket real-time push | None (via RabbitMQ only) |
| **iot-uplink** | Device→Platform message processing | PostgreSQL (read), InfluxDB (write) |
| **iot-downlink** | Platform→Device command routing | PostgreSQL |

### Critical Constraints

- **MQTT Isolation**: Only `iot-gateway` may connect to VerneMQ. Other services MUST use RabbitMQ.
- **No Direct Service Calls**: All inter-service communication via RabbitMQ. No HTTP/gRPC between services.
- **RabbitMQ Routing Keys**: Format `iot.{vendor}.{service}.{action}` (e.g., `iot.dji.uplink.property.report`)
- **Message Format**: All messages must include `device_sn`, `tid`, `bid`, `timestamp`, and W3C Trace Context headers

## Code Standards

- **Style**: Uber Go Style Guide (strictly enforced)
- **ORM**: GORM only. Raw SQL prohibited. Use AutoMigrate for migrations.
- **Logging**: logrus with JSON format. Include `trace_id` and `span_id` in all logs.
- **Metrics**: Prometheus format at `/metrics`. Naming: `iot_{component}_{metric_type}_{unit}`
- **Required Labels**: `service`, `vendor`, `message_type`, `status`

### Naming Conventions

- **Packages**: lowercase, no underscores (e.g., `thingmodel`, `iotapi`)
- **Files**: snake_case (e.g., `device_manager.go`, `thing_model.go`)
- **Variables/Functions**: camelCase for private, PascalCase for exported
- **Constants**: ALL_CAPS with underscores

## Tech Stack (Non-Negotiable)

- **Language**: Go 1.22+
- **Web**: Gin Framework
- **ORM**: GORM
- **Databases**: PostgreSQL (business), InfluxDB (time-series)
- **Message Queue**: RabbitMQ (services), VerneMQ (MQTT devices)
- **Observability**: Prometheus, Loki, Tempo, Grafana

## Project Structure

```
cmd/                    # Service entry points (main.go for each service)
internal/shared/        # Shared internal code (config, logger)
pkg/                    # Public packages (models, metrics, tracer, rabbitmq)
specs/                  # Feature specifications (spec.md, plan.md, tasks.md)
deployments/            # Docker and Kubernetes configs
```

## Constitution Reference

See `.specify/memory/constitution.md` for complete architectural principles. Key NON-NEGOTIABLE items:
- Thing Model Driven Architecture (TSL JSON format)
- Test-First Development (TDD mandatory)
- Microservice Architecture (5 services, RabbitMQ communication)

## Recent Changes
- 003-dji-protocol-implementation: Added [if applicable, e.g., PostgreSQL, CoreData, files or N/A]
- 003-dji-protocol-implementation: Added Go 1.22+
