# Quickstart: Core Services Implementation

**Feature**: 004-core-services-implementation
**Date**: 2026-03-23

## Purpose

This quickstart verifies that the five core services can be brought up together and support the primary operational flows:
- device connectivity and upstream ingestion
- telemetry storage and query
- tracked device command execution
- realtime client subscriptions

## Prerequisites

- Go 1.22+
- Docker and Docker Compose
- Running PostgreSQL, InfluxDB, RabbitMQ, and VerneMQ from the project's local stack
- Valid local configuration for all five core services
- Existing supported capability models and vendor adapters already available in the repository

## 1. Start shared infrastructure

```bash
docker-compose up -d
```

Verify required infrastructure is healthy before starting services.

## 2. Start the five core services

In separate terminals:

```bash
make run-gateway
make run-uplink
make run-downlink
make run-api
make run-ws
```

## 3. Verify health and readiness

Check that each core service exposes health and readiness endpoints.

```bash
curl http://localhost:8082/health
curl http://localhost:8082/ready

curl http://localhost:8083/health
curl http://localhost:8083/ready

curl http://localhost:8084/health
curl http://localhost:8084/ready

curl http://localhost:8080/health
curl http://localhost:8080/ready

curl http://localhost:8081/health
curl http://localhost:8081/ready
```

Expected outcome:
- all `/health` endpoints return success
- `/ready` returns success only when service dependencies are available

## 4. Verify device inventory flow

Create a managed device through the HTTP API.

```bash
curl -X POST http://localhost:8080/api/v1/devices \
  -H "Content-Type: application/json" \
  -H "X-API-Key: local-dev-key" \
  -d '{
    "deviceSN": "device-001",
    "deviceName": "Dock Aircraft 001",
    "deviceType": "aircraft",
    "vendor": "dji"
  }'
```

Expected outcome:
- the API returns the stored device record
- the device can be listed and retrieved afterward

## 5. Verify upstream message flow

Use a supported simulator or MQTT publisher to send a valid upstream message for a provisioned device.

Expected outcome:
- `iot-gateway` authenticates and bridges the message
- `iot-uplink` processes and validates the message
- telemetry becomes queryable through the API
- realtime consumers become eligible to receive the update

Optional verification:

```bash
curl -H "X-API-Key: local-dev-key" \
  "http://localhost:8080/api/v1/devices/device-001/telemetry?start=-15m&end=now"
```

## 6. Verify command flow

Submit a tracked service request.

```bash
curl -X POST http://localhost:8080/api/v1/service-requests \
  -H "Content-Type: application/json" \
  -H "X-API-Key: local-dev-key" \
  -d '{
    "deviceSN": "device-001",
    "vendor": "dji",
    "method": "flighttask_prepare",
    "callType": "command",
    "params": {
      "file_id": "wayline-001"
    }
  }'
```

Expected outcome:
- the API returns a tracked request identifier
- the request is routed through the downlink path toward gateway delivery
- later status checks reflect the final request outcome

Follow-up status check:

```bash
curl -H "X-API-Key: local-dev-key" \
  http://localhost:8080/api/v1/service-requests/<request-id>
```

## 7. Verify realtime subscription flow

Open a WebSocket session to the realtime service and subscribe to relevant topics.

Example subscription payload:

```json
{
  "topics": ["property.processed", "event.processed", "device.device-001.property.processed"]
}
```

Expected outcome:
- the session is accepted
- matching upstream updates are pushed without polling
- stale connections are closed after missed heartbeats

## 8. Run validation tests

Run the repository test and quality gates.

```bash
make test
make lint
make coverage
```

Expected outcome:
- unit, integration, and end-to-end relevant tests pass
- linting passes
- coverage meets the 80% threshold

## Troubleshooting

### Readiness fails while health passes

Interpretation:
- the process is running but a required dependency or runtime worker is unavailable

Check:
- RabbitMQ connectivity for `iot-gateway`, `iot-uplink`, `iot-downlink`, and `iot-ws`
- PostgreSQL connectivity for `iot-api` and gateway-auth/device records
- InfluxDB connectivity for telemetry query and storage paths

### Commands are accepted but never complete

Check:
- whether `iot-downlink` is consuming the intended command path
- whether gateway downlink routing matches the canonical RabbitMQ routing conventions
- whether the target device is online and capable of handling the requested method

### Realtime clients connect but receive no updates

Check:
- whether `iot-uplink` publishes normalized messages to the expected routing keys
- whether `iot-ws` is consuming the correct RabbitMQ queue bindings
- whether the client subscribed to the exact emitted topic pattern
