# Research: Core Services Implementation

**Feature**: 004-core-services-implementation
**Date**: 2026-03-23

## Research Topics

### 1. Message routing and service boundaries

**Decision**: Keep `iot-gateway` as the only service that connects to MQTT/VerneMQ, route raw vendor traffic through RabbitMQ, and preserve `iot-downlink` as the dedicated downlink orchestration service rather than dispatching device commands directly from `iot-api`.

**Rationale**:
- The constitution and repository guidance require strict MQTT isolation and RabbitMQ-only inter-service communication.
- Existing code already establishes the right boundary in `cmd/iot-gateway/main.go` and `internal/gateway/bridge/*.go`, where gateway owns MQTT connectivity and RabbitMQ bridging.
- The current codebase shows architectural drift where `iot-api` instantiates a downlink dispatcher directly, while `iot-downlink` is not yet fully wired as the active consumer path. Planning should treat that as technical debt to resolve, not as the target architecture.

**Alternatives considered**:
- Let `iot-api` publish device commands directly: rejected because it weakens the single-responsibility boundary for `iot-downlink` and conflicts with the feature spec.
- Allow non-gateway services to connect to MQTT: rejected because it violates the constitution and project constraints.

### 2. RabbitMQ contract and routing-key strategy

**Decision**: Treat `pkg/rabbitmq` as the canonical message contract source and align planning artifacts to two routing-key families: standard business routing keys `iot.{vendor}.{service}.{action}` and raw bridge routing keys `iot.raw.{vendor}.{direction}`.

**Rationale**:
- `pkg/rabbitmq/routing.go` defines the intended routing-key contract and helper builders.
- `pkg/rabbitmq/message.go` defines the standard message envelope with `tid`, `bid`, `timestamp`, `device_sn`, business action metadata, and W3C trace context headers via `MessageHeader`.
- Existing internal routers in `internal/uplink/router/router.go` and `internal/downlink/router/router.go` use shortcuts such as `iot.ws.property` and `iot.gateway.downlink`; those patterns are easier to implement but drift from the declared standard and complicate queue binding consistency.

**Alternatives considered**:
- Preserve ad hoc routing keys per service: rejected because it creates terminology drift and weakens interoperability.
- Use only one routing-key family for both normalized and raw vendor traffic: rejected because raw ingress/egress and normalized business messages have different consumers and lifecycle concerns.

### 3. Device and capability modeling

**Decision**: Base the design on the existing GORM-backed device and thing-model records, while extending feature planning to include credentials, topology, telemetry history, and service request lifecycle as first-class artifacts.

**Rationale**:
- `pkg/models/device.go` already models managed devices with `device_sn`, vendor, type, gateway relationship, status, and thing-model linkage.
- `pkg/models/thing_model.go` already models TSL-backed capability definitions using JSON storage and product keys.
- Existing feature artifacts add important operational entities that are not all present in shared migrations yet, especially credentials and service-call tracking. Planning should treat those as necessary design artifacts to keep the spec, plan, and future tasks aligned.

**Alternatives considered**:
- Model service requests only as transient queue messages: rejected because operators need auditable status and retry history.
- Keep capability validation purely vendor-specific in adapter code: rejected because the constitution requires thing-model-driven abstraction at the platform level.

### 4. Realtime delivery and subscription behavior

**Decision**: Keep realtime delivery centered on `iot-ws`, with RabbitMQ-fed push, per-client subscriptions, and heartbeat-based stale connection cleanup, while planning for authorization-aware subscription control.

**Rationale**:
- `internal/ws/service.go`, `internal/ws/hub/client.go`, and `internal/ws/subscription/manager.go` already provide the runtime shape: hub-managed clients, topic subscriptions, ping/pong heartbeats, and push from a RabbitMQ-backed queue.
- The current implementation extracts `user_id` and `device_sn` from the request but does not yet enforce subscription authorization boundaries. The spec requires authorized subscriptions, so planning must call this out as a behavior the design and tasks must preserve.

**Alternatives considered**:
- Poll-based client updates instead of WebSocket push: rejected because the constitution explicitly requires WebSocket support for realtime delivery.
- Broker-specific subscription logic embedded into API or uplink services: rejected because `iot-ws` is the dedicated realtime boundary.

### 5. Health, readiness, observability, and test strategy

**Decision**: Plan for every core service to expose health/readiness endpoints, structured logs, distributed tracing, Prometheus metrics, and layered tests (unit, integration, end-to-end, and contract), with readiness checks reflecting real dependency status rather than only process liveness.

**Rationale**:
- All five core services already expose `/health`; several also expose `/ready`, but the checks are inconsistent in depth.
- The constitution requires observability, API documentation, and test-first development with contract coverage.
- The repository already contains test layers (`tests/integration/*.go`, unit tests alongside packages, and OpenAPI-oriented API patterns), so the plan should standardize and extend those patterns rather than invent new ones.

**Alternatives considered**:
- Simple liveness-only health checks: rejected because they do not satisfy the feature requirement for operational readiness visibility.
- Rely only on integration tests: rejected because the constitution explicitly requires unit, integration, and contract-oriented testing.

## Conclusions

1. `iot-gateway` remains the sole MQTT/VerneMQ boundary and all other services communicate through RabbitMQ.
2. Canonical routing should come from `pkg/rabbitmq` and planning should treat current ad hoc routing keys as drift to be corrected.
3. Shared GORM models for devices and thing models are the foundation, with credentials, service requests, topology, and telemetry persistence explicitly designed around them.
4. `iot-ws` remains the realtime boundary with heartbeat and subscription management, but authorization-aware subscriptions must be accounted for in planning.
5. The implementation plan must preserve constitution-mandated observability, OpenAPI contracts, and TDD-driven verification with readiness checks that represent actual dependency health.
