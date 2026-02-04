# Getting Started with UMOS Development

## Prerequisites

- **Go 1.22+**: [Download Go](https://go.dev/dl/)
- **Docker & Docker Compose**: For local infrastructure
- **Make**: Build automation
- **Git**: Version control

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/utmos/utmos.git
cd utmos
```

### 2. Start Infrastructure

Start the required services (PostgreSQL, RabbitMQ, InfluxDB, Tempo, Grafana):

```bash
docker-compose up -d
```

Verify services are running:
```bash
docker-compose ps
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Build All Services

```bash
make build
```

Or build individually:
```bash
go build -o bin/iot-api ./cmd/iot-api
go build -o bin/iot-ws ./cmd/iot-ws
go build -o bin/iot-uplink ./cmd/iot-uplink
go build -o bin/iot-downlink ./cmd/iot-downlink
go build -o bin/iot-gateway ./cmd/iot-gateway
```

### 5. Run Tests

```bash
make test
```

Or with coverage:
```bash
go test ./... -cover
```

### 6. Run a Service

```bash
# Set environment
export APP_ENV=dev

# Run iot-api
./bin/iot-api

# Or run directly
go run ./cmd/iot-api
```

## Project Structure

```
utmos/
├── api/v1/                 # OpenAPI specifications
├── cmd/                    # Service entry points
│   ├── iot-api/           # REST API server
│   ├── iot-ws/            # WebSocket server
│   ├── iot-uplink/        # Uplink processor
│   ├── iot-downlink/      # Downlink processor
│   └── iot-gateway/       # MQTT gateway
├── configs/               # Configuration files
│   ├── config.dev.yaml
│   └── config.prod.yaml
├── deployments/           # Deployment configurations
│   └── docker/            # Dockerfiles
├── docs/                  # Documentation
├── internal/              # Private application code
│   └── shared/            # Shared internal packages
│       ├── config/        # Configuration loading
│       ├── database/      # Database connections
│       ├── logger/        # Logging utilities
│       └── server/        # Server utilities
├── pkg/                   # Public packages
│   ├── errors/            # Error handling
│   ├── metrics/           # Prometheus metrics
│   ├── models/            # GORM data models
│   ├── rabbitmq/          # RabbitMQ client
│   ├── repository/        # Data repositories
│   └── tracer/            # OpenTelemetry tracing
├── tests/                 # Integration tests
│   └── integration/
├── docker-compose.yml     # Local development stack
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Environment (dev/prod) | `dev` |
| `SERVER_HOST` | Server bind host | `0.0.0.0` |
| `SERVER_PORT` | Server port | `8080` |
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |
| `RABBITMQ_URL` | RabbitMQ connection URL | `amqp://guest:guest@localhost:5672/` |
| `TRACER_ENDPOINT` | Tempo OTLP endpoint | `http://localhost:4318` |

### Configuration Files

Development configuration (`configs/config.dev.yaml`):
```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

database:
  postgres:
    host: localhost
    port: 5432
    user: postgres
    password: postgres
    dbname: umos
    sslmode: disable

rabbitmq:
  url: "amqp://guest:guest@localhost:5672/"
  exchange_name: "iot"
  exchange_type: "topic"

tracer:
  enabled: true
  endpoint: "http://localhost:4318"
  sampling_rate: 1.0

logger:
  level: debug
  format: json
```

## Development Workflow

### Adding a New Feature

1. Create a feature branch:
   ```bash
   git checkout -b feature/my-feature
   ```

2. Write tests first (TDD):
   ```bash
   # Create test file
   touch pkg/mypackage/myfeature_test.go
   ```

3. Implement the feature

4. Run tests:
   ```bash
   go test ./pkg/mypackage/...
   ```

5. Run linter:
   ```bash
   make lint
   ```

6. Commit and push:
   ```bash
   git add .
   git commit -m "Add my feature"
   git push origin feature/my-feature
   ```

### Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Run `golangci-lint` before committing

### Testing

```bash
# Run all tests
make test

# Run specific package tests
go test ./pkg/rabbitmq/...

# Run with verbose output
go test -v ./...

# Run integration tests
go test ./tests/integration/...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Local Development URLs

| Service | URL |
|---------|-----|
| iot-api | http://localhost:8080 |
| iot-ws | http://localhost:8081 |
| RabbitMQ Management | http://localhost:15672 |
| Grafana | http://localhost:3000 |
| Tempo | http://localhost:3200 |
| PostgreSQL | localhost:5432 |
| InfluxDB | http://localhost:8086 |

### Default Credentials

| Service | Username | Password |
|---------|----------|----------|
| RabbitMQ | guest | guest |
| Grafana | admin | admin |
| PostgreSQL | postgres | postgres |
| InfluxDB | admin | adminpassword |

## Debugging

### View Service Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f rabbitmq
```

### Check RabbitMQ Queues

1. Open http://localhost:15672
2. Login with guest/guest
3. Navigate to Queues tab

### View Traces in Tempo

1. Open Grafana at http://localhost:3000
2. Go to Explore
3. Select Tempo data source
4. Search by trace ID or service name

### Database Access

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U postgres -d umos

# Common queries
\dt                    # List tables
SELECT * FROM devices; # View devices
```

## Troubleshooting

### RabbitMQ Connection Failed

```bash
# Check if RabbitMQ is running
docker-compose ps rabbitmq

# View RabbitMQ logs
docker-compose logs rabbitmq

# Restart RabbitMQ
docker-compose restart rabbitmq
```

### Database Connection Issues

```bash
# Check PostgreSQL status
docker-compose ps postgres

# Reset database
docker-compose down -v
docker-compose up -d postgres
```

### Build Errors

```bash
# Clean and rebuild
go clean -cache
go mod tidy
go build ./...
```

## Next Steps

1. Read the [Architecture Overview](../architecture/overview.md)
2. Review the [API Documentation](../../api/v1/openapi.yaml)
3. Explore the codebase starting with `cmd/iot-api/main.go`
4. Join the team chat for questions
