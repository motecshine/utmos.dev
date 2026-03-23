# Implementation Plan: Core Services Implementation

**Branch**: `004-core-services-implementation` | **Date**: 2026-03-23 | **Spec**: [/Users/laputalaputa/Desktop/code/utmos.dev/specs/004-core-services-implementation/spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-core-services-implementation/spec.md`

## Summary

Implement and align the five UMOS core services so the platform can authenticate devices, process upstream telemetry and events, track downlink service requests, expose inventory and history through the API, and deliver realtime updates through WebSocket while preserving the project's RabbitMQ-based microservice boundaries.

## Technical Context

**Language/Version**: Go 1.22+
**Primary Dependencies**: Gin Framework, GORM, RabbitMQ client utilities in `pkg/rabbitmq`, existing vendor adapter packages under `pkg/adapter/*`, logrus-based logging, OpenTelemetry-based tracing
**Storage**: PostgreSQL for relational platform data, InfluxDB for telemetry/time-series data
**Testing**: Go unit tests, integration tests under `tests/integration`, OpenAPI/contract verification, end-to-end flow validation, coverage gate via `make coverage`
**Target Platform**: Linux containerized services running via Docker Compose and deployable to Kubernetes-compatible environments
**Project Type**: Go microservices repository
**Performance Goals**: 95% of valid telemetry available to downstream consumers within 1 second; support 1,000 connected devices and 10,000 realtime connections per node; recover from transient dependency interruptions within 30 seconds; demonstrate 99.9% availability during controlled resilience testing windows
**Constraints**: Only `iot-gateway` may connect to VerneMQ/MQTT; all inter-service communication must use RabbitMQ; normalized routing should follow `iot.{vendor}.{service}.{action}` and raw bridge traffic should follow `iot.raw.{vendor}.{direction}`; all inter-service messages must preserve `tid`, `bid`, `timestamp`, `device_sn`, and W3C trace context; all database access must use GORM; OpenAPI documentation and contract tests are required for API behavior; consumers handling RabbitMQ traffic must process duplicate deliveries idempotently using `tid` and `bid`
**Scale/Scope**: Five core services (`iot-gateway`, `iot-uplink`, `iot-downlink`, `iot-api`, `iot-ws`) plus their shared contracts, runtime health/readiness behavior, and end-to-end operational flow

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| Thing Model Driven Architecture | PASS | Spec and data model center validation and interpretation on TSL-backed capability models. |
| Multi-Protocol Support | PASS | The feature spans MQTT device ingress, HTTP API access, and WebSocket realtime delivery. |
| Device Abstraction Layer | PASS | Planning keeps vendor-specific parsing in adapters and core services working on normalized device concepts. |
| Extensibility & Plugin Architecture | PASS | Existing vendor adapters remain separate from core service orchestration and routing. |
| Standardized API Design | PASS | OpenAPI contract is included and uses versioned HTTP API paths. |
| Test-First Development | PASS | Plan requires unit, integration, end-to-end, and contract coverage with an 80% coverage gate. |
| Observability & Monitoring | PASS | Plan preserves health/readiness endpoints, metrics, structured logging, and trace propagation. |
| Technology Stack Standards | PASS | Plan stays on Go, Gin, and GORM as required by the constitution. |
| Microservice Architecture | PASS | The plan preserves the five-service split and RabbitMQ-only service-to-service communication. |
| Infrastructure & Middleware Standards | PASS | Plan continues using VerneMQ, RabbitMQ, PostgreSQL, InfluxDB, and Prometheus-compatible observability tooling. |

## Project Structure

### Documentation (this feature)

```text
specs/004-core-services-implementation/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── openapi.yaml
└── tasks.md
```

### Source Code (repository root)

```text
cmd/
├── dji-adapter/main.go
├── iot-api/main.go
├── iot-downlink/main.go
├── iot-gateway/main.go
├── iot-uplink/main.go
└── iot-ws/main.go

internal/
├── api/
│   ├── handler/
│   ├── middleware/
│   └── router.go
├── downlink/
│   ├── dispatcher/
│   ├── model/
│   ├── retry/
│   ├── router/
│   └── service.go
├── gateway/
│   ├── bridge/
│   ├── connection/
│   ├── model/
│   ├── mqtt/
│   └── service.go
├── shared/
│   ├── config/
│   ├── database/
│   └── server/
├── uplink/
│   ├── processor/
│   ├── router/
│   ├── storage/
│   └── service.go
└── ws/
    ├── hub/
    ├── push/
    ├── subscription/
    └── service.go

pkg/
├── adapter/dji/
├── logger/
├── metrics/
├── models/
├── rabbitmq/
└── tracer/

tests/
├── integration/
└── mocks/
```

**Structure Decision**: Use the existing Go microservices layout already present in the repository. Core runtime behavior remains under `internal/<service>/`, shared contracts and infrastructure remain under `pkg/` and `internal/shared/`, and cross-service verification continues under `tests/integration/`.

## Phase 0: Research Summary

Research completed in `research.md` resolved the main planning questions:
- the canonical service boundary keeps MQTT isolated inside `iot-gateway`
- RabbitMQ contract helpers in `pkg/rabbitmq` are the source of truth for routing and message shape
- the current codebase contains routing-key and downlink-boundary drift that planning should explicitly correct
- realtime subscriptions and heartbeat behavior stay in `iot-ws`
- readiness and observability requirements must be standardized across all five services

## Phase 1: Design Artifacts

### Data model

See `data-model.md` for logical entities and lifecycle behavior covering:
- thing models
- managed devices
- device credentials
- device topology
- telemetry records
- service requests
- realtime sessions and subscriptions

### Contracts

See `contracts/openapi.yaml` for the current planning-time API contract covering:
- device inventory operations
- telemetry query
- service request submission and status lookup
- realtime subscription entry points
- health and readiness endpoints

### Quickstart

See `quickstart.md` for the validation flow that brings up the five services, verifies health/readiness, and exercises the primary user stories.

## Post-Design Constitution Re-Check

| Principle | Status | Notes |
|-----------|--------|-------|
| Thing Model Driven Architecture | PASS | Data model and plan keep capability validation tied to thing models. |
| Multi-Protocol Support | PASS | Design artifacts cover MQTT ingress, HTTP APIs, and WebSocket delivery. |
| Device Abstraction Layer | PASS | Vendor-specific concerns remain adapter-facing; core services consume normalized concepts. |
| Extensibility & Plugin Architecture | PASS | No core design change collapses vendor adapters into service-specific logic. |
| Standardized API Design | PASS | OpenAPI contract stays versioned and consistent with standardized HTTP entry points. |
| Test-First Development | PASS | Quickstart and plan preserve layered testing and quality gates. |
| Observability & Monitoring | PASS | Health/readiness, logging, metrics, and tracing remain explicit design concerns. |
| Technology Stack Standards | PASS | No unapproved technology changes introduced. |
| Microservice Architecture | PASS | Design keeps five services and identifies current boundary drift as something implementation must resolve. |
| Infrastructure & Middleware Standards | PASS | Approved infrastructure remains unchanged. |

## Implementation Phases

### Phase 1: Gateway ingress and connectivity
- enforce device authentication and credential lookup
- maintain device connectivity status
- bridge MQTT ingress/egress through canonical raw RabbitMQ routing
- expose meaningful gateway readiness based on MQTT and RabbitMQ health

### Phase 2: Uplink normalization and storage
- consume raw vendor uplink traffic from RabbitMQ
- validate normalized messages against thing-model capabilities
- persist telemetry and event history idempotently
- route normalized updates to downstream consumers using canonical business routing keys

### Phase 3: Downlink orchestration and tracking
- receive service requests through RabbitMQ rather than direct API-to-device dispatch
- persist and update service-request lifecycle state idempotently, including protection against late replies after timeout
- apply retry and dead-letter behavior with audit visibility
- route device-bound requests to gateway through canonical downlink routing

### Phase 4: API surface and inventory workflows
- expose versioned inventory and telemetry APIs aligned with the contract
- submit tracked service requests into the downlink path
- return auditable request status and not-found/error semantics consistently
- expose health/readiness and OpenAPI documentation

### Phase 5: Realtime delivery
- manage WebSocket sessions and subscriptions
- enforce heartbeat-based stale-session cleanup
- consume downstream RabbitMQ messages and push authorized updates to subscribers
- expose operational stats and readiness aligned with actual dependencies

### Phase 6: Cross-service verification
- verify full upstream and downlink flows
- validate resilience, readiness, and recovery behavior
- enforce linting, contract checks, and coverage thresholds

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Routing-key drift between helpers and service-specific routers | High | Standardize all new planning and task generation on `pkg/rabbitmq` routing conventions. |
| `iot-api` bypasses `iot-downlink` for command dispatch | High | Treat downlink orchestration as a dedicated implementation phase and regenerate tasks around that boundary. |
| `iot-downlink` runtime consumer path remains incomplete | High | Ensure tasks explicitly wire subscriber-based consumption and lifecycle persistence. |
| Readiness endpoints report process liveness rather than dependency readiness | Medium | Define readiness checks per service around actual broker, database, and runtime worker state. |
| Realtime subscriptions lack authorization-aware filtering | Medium | Include authorization checks and subscription validation in ws/api task planning. |

## Complexity Tracking

No constitution violations require special justification.
