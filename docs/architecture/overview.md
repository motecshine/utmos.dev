# UMOS IoT Platform Architecture

## Overview

UMOS is a distributed IoT platform designed for multi-vendor device management with distributed tracing, message routing, and real-time communication capabilities.

## System Architecture

```
                                    ┌─────────────────┐
                                    │   Web Clients   │
                                    │  (Dashboard/App)│
                                    └────────┬────────┘
                                             │ WebSocket
                                             ▼
┌─────────────┐    REST API    ┌─────────────────────────┐
│  External   │◄──────────────►│       iot-api           │
│  Systems    │                │   (REST API Server)     │
└─────────────┘                └───────────┬─────────────┘
                                           │ RabbitMQ
                                           ▼
                               ┌─────────────────────────┐
                               │       RabbitMQ          │
                               │   (Message Broker)      │
                               │  Topic Exchange: iot    │
                               └─────────────────────────┘
                                     │           │
                    ┌────────────────┘           └────────────────┐
                    ▼                                             ▼
        ┌─────────────────────┐                       ┌─────────────────────┐
        │     iot-uplink      │                       │    iot-downlink     │
        │ (Uplink Processor)  │                       │ (Downlink Processor)│
        └──────────┬──────────┘                       └──────────┬──────────┘
                   │                                             │
                   │ RabbitMQ                         RabbitMQ   │
                   ▼                                             ▼
        ┌─────────────────────┐                       ┌─────────────────────┐
        │      iot-ws         │                       │    iot-gateway      │
        │ (WebSocket Server)  │                       │  (MQTT ↔ RabbitMQ)  │
        └─────────────────────┘                       └──────────┬──────────┘
                                                                 │ MQTT
                                                                 ▼
                                                      ┌─────────────────────┐
                                                      │      VerneMQ        │
                                                      │   (MQTT Broker)     │
                                                      └──────────┬──────────┘
                                                                 │
                                                                 ▼
                                                      ┌─────────────────────┐
                                                      │    IoT Devices      │
                                                      │  (DJI/Tuya/Generic) │
                                                      └─────────────────────┘
```

## Services

### iot-api
- **Purpose**: REST API server for device management and command dispatch
- **Port**: 8080
- **Responsibilities**:
  - Device CRUD operations
  - Thing model management
  - Command dispatch via RabbitMQ
  - Health and metrics endpoints

### iot-ws
- **Purpose**: WebSocket server for real-time client communication
- **Port**: 8081
- **Responsibilities**:
  - WebSocket connection management
  - Real-time message push to clients
  - Subscribe to uplink messages from RabbitMQ

### iot-uplink
- **Purpose**: Uplink message processor (device → cloud)
- **Port**: 8082
- **Responsibilities**:
  - Process messages from devices
  - Route to appropriate services
  - Message transformation and validation

### iot-downlink
- **Purpose**: Downlink message processor (cloud → device)
- **Port**: 8083
- **Responsibilities**:
  - Process commands from API
  - Route to gateway for device delivery
  - Command acknowledgment handling

### iot-gateway
- **Purpose**: MQTT ↔ RabbitMQ bridge
- **Port**: 8084
- **Responsibilities**:
  - MQTT client connection to VerneMQ
  - Message protocol conversion
  - Only service allowed to connect to MQTT broker

## Message Flow

### Uplink (Device → Cloud)
```
Device → MQTT → VerneMQ → iot-gateway → RabbitMQ → iot-uplink → iot-ws → Client
                                                              → iot-api (storage)
```

### Downlink (Cloud → Device)
```
Client → iot-api → RabbitMQ → iot-downlink → RabbitMQ → iot-gateway → MQTT → Device
```

## RabbitMQ Routing

### Exchange Configuration
- **Exchange Name**: `iot`
- **Exchange Type**: `topic`

### Routing Key Format
```
iot.{vendor}.{service}.{action}
```

Examples:
- `iot.dji.device.property.report` - DJI device property report
- `iot.tuya.event.event.report` - Tuya event report
- `iot.generic.service.service.call` - Generic service call

### Vendors
| Vendor | Description |
|--------|-------------|
| `dji` | DJI drones and devices |
| `tuya` | Tuya smart home devices |
| `generic` | Generic MQTT devices |

### Queue Bindings
```
iot-uplink-queue    → iot.*.device.#
iot-downlink-queue  → iot.*.service.#
iot-ws-queue        → iot.#
```

## Distributed Tracing

### Technology
- **OpenTelemetry SDK** for instrumentation
- **Tempo** for trace storage and querying
- **W3C Trace Context** for propagation

### Trace Propagation
1. HTTP requests: `traceparent` and `tracestate` headers
2. RabbitMQ messages: Custom headers with W3C format
3. All services extract and inject trace context

### Sampling
- **Development**: 100% sampling rate
- **Production**: 10% sampling rate

## Metrics

### Technology
- **Prometheus** for metrics collection
- **Grafana** for visualization

### Naming Convention
```
iot_{component}_{metric}_{unit}
```

### Standard Metrics
| Metric | Type | Description |
|--------|------|-------------|
| `iot_http_requests_total` | Counter | Total HTTP requests |
| `iot_http_request_duration_seconds` | Histogram | HTTP request latency |
| `iot_rabbitmq_messages_total` | Counter | Total RabbitMQ messages |
| `iot_devices_active` | Gauge | Active device count |

## Data Storage

### PostgreSQL
- Device registry
- Thing models
- Message logs
- Configuration

### InfluxDB
- Device telemetry time-series data
- Metrics history
- Event logs

## Security Considerations

1. **Authentication**: JWT Bearer tokens for API access
2. **Authorization**: Role-based access control
3. **Encryption**: TLS for all external communication
4. **MQTT Security**: VerneMQ authentication plugins

## Deployment

### Docker Compose (Development)
```bash
docker-compose up -d
```

### Kubernetes (Production)
- Helm charts in `deployments/k8s/`
- Horizontal pod autoscaling
- Service mesh integration (optional)

## Configuration

Configuration files are located in `configs/`:
- `config.dev.yaml` - Development settings
- `config.prod.yaml` - Production settings

Environment variables override file settings.
