# Tasks: Core Services Implementation

**Input**: Design documents from `/specs/004-core-services-implementation/`
**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/openapi.yaml`, `quickstart.md`

**Tests**: Required. The constitution and `NFR-006` require unit, integration, end-to-end, and contract coverage. Write the failing tests for each story before implementation.

**Organization**: Tasks are grouped by user story so each story can be implemented and verified independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependency on incomplete tasks)
- **[Story]**: User story label (`[US1]`, `[US2]`, `[US3]`, `[US4]`)
- Every task includes the exact file paths it changes when those paths are knowable in advance; repository-wide validation tasks may reference the command-reported files they reconcile

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Align the existing repository scaffolding with the current plan before story work starts.

- [X] T001 Update shared model migration coverage in `pkg/models/migrate.go` and `internal/downlink/model/service_call.go` so the planned core-service entities are migratable.
- [X] T002 [P] Align canonical RabbitMQ routing and message-envelope helpers with the plan in `pkg/rabbitmq/routing.go`, `pkg/rabbitmq/message.go`, `pkg/rabbitmq/routing_test.go`, and `pkg/rabbitmq/message_test.go`.
- [X] T003 [P] Standardize service startup wiring and dependency configuration entry points in `cmd/iot-gateway/main.go`, `cmd/iot-uplink/main.go`, `cmd/iot-downlink/main.go`, `cmd/iot-api/main.go`, and `cmd/iot-ws/main.go`.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that must be complete before user stories can be delivered safely.

**⚠️ CRITICAL**: No user story work should be considered done until this phase is complete.

- [X] T004 Add failing foundational integration coverage for canonical routing-key usage and trace propagation in `tests/integration/routing_test.go` and `tests/integration/tracing_test.go`.
- [X] T005 [P] Normalize RabbitMQ queue bindings and publisher/subscriber setup for raw and business traffic in `internal/gateway/bridge/uplink.go`, `internal/uplink/service.go`, `internal/downlink/service.go`, and `internal/ws/push/pusher.go`.
- [X] T006 [P] Implement dependency-aware readiness responses and required metrics labels across `cmd/iot-gateway/main.go`, `cmd/iot-uplink/main.go`, `cmd/iot-downlink/main.go`, `cmd/iot-api/main.go`, `cmd/iot-ws/main.go`, and `tests/integration/metrics_test.go`.
- [X] T007 Establish shared audit/logging hooks for authentication failures, uplink processing failures, retry attempts, and final command outcomes in `internal/gateway/service.go`, `internal/uplink/service.go`, `internal/downlink/service.go`, and `internal/api/handler/service.go`.

**Checkpoint**: Canonical routing, shared readiness, migrations, and audit hooks are in place.

---

## Phase 3: User Story 1 - Operate connected devices (Priority: P1) 🎯 MVP

**Goal**: Authenticate devices, track connectivity, normalize upstream traffic, and persist telemetry/event history.

**Independent Test**: Connect a supported simulator or device with valid credentials, publish status and telemetry, then confirm the device becomes visible as online and telemetry/history becomes available without breaking other traffic.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests first, ensure they fail, then implement the story.**

- [X] T008 [P] [US1] Add failing gateway authentication and connectivity integration coverage in `tests/integration/message_flow_test.go` and `tests/integration/routing_test.go`.
- [X] T009 [P] [US1] Add failing uplink normalization, malformed-message, and duplicate-transaction-id coverage in `tests/integration/dji_osd_test.go` and `tests/integration/message_flow_test.go`.
- [X] T010 [P] [US1] Add failing unit coverage for gateway auth/connection and uplink routing/storage in `internal/gateway/mqtt/auth_test.go`, `internal/gateway/connection/manager_test.go`, `internal/uplink/router/router_test.go`, and `internal/uplink/storage/influx_test.go`.

### Implementation for User Story 1

- [X] T011 [US1] Implement device credential lookup and authenticated MQTT admission in `internal/gateway/model/credential.go` and `internal/gateway/mqtt/auth.go`.
- [X] T012 [US1] Implement online/offline connection tracking and stale-device transitions in `internal/gateway/connection/manager.go` and `internal/gateway/service.go`.
- [X] T013 [US1] Implement canonical raw uplink bridging from MQTT to RabbitMQ in `internal/gateway/mqtt/handler.go` and `internal/gateway/bridge/uplink.go`.
- [X] T014 [US1] Implement thing-model-aware uplink validation, duplicate-message idempotency, and normalization in `internal/uplink/processor/processor.go` and `internal/uplink/service.go`.
- [X] T015 [US1] Implement idempotent telemetry/event persistence and canonical downstream routing in `internal/uplink/storage/influx.go` and `internal/uplink/router/router.go`.
- [X] T016 [US1] Wire gateway and uplink runtime readiness to the authenticated ingestion flow in `cmd/iot-gateway/main.go` and `cmd/iot-uplink/main.go`.

**Checkpoint**: Supported devices can authenticate, publish upstream traffic, and produce persisted normalized history.

---

## Phase 4: User Story 2 - Send commands to devices (Priority: P1)

**Goal**: Accept tracked service requests through the API, route them through `iot-downlink`, retry failures, and expose auditable request state.

**Independent Test**: Submit a service request for an online device, receive an immediate tracking identifier, observe the request flow through downlink and gateway, and confirm the final status is queryable.

### Tests for User Story 2 ⚠️

- [X] T017 [P] [US2] Add failing contract coverage for `POST /api/v1/service-requests` and `GET /api/v1/service-requests/{requestID}` in `tests/contract/service_requests_contract_test.go` and `specs/004-core-services-implementation/contracts/openapi.yaml`.
- [X] T018 [P] [US2] Add failing integration coverage for API-to-downlink-to-gateway request flow, timeout handling, late replies after timeout, and retry exhaustion in `tests/integration/dji_e2e_service_test.go` and `tests/integration/message_flow_test.go`.
- [X] T019 [P] [US2] Add failing unit coverage for request persistence, routing, retry behavior, and terminal-state idempotency in `internal/api/handler/service_test.go`, `internal/downlink/router/router_test.go`, `internal/downlink/retry/retry_test.go`, and `internal/downlink/model/service_call_test.go`.

### Implementation for User Story 2

- [X] T020 [US2] Refactor tracked service-request persistence to match the current contract in `internal/downlink/model/service_call.go` and `pkg/models/migrate.go`.
- [X] T021 [US2] Move API command submission from direct dispatcher execution to RabbitMQ-backed downlink orchestration in `internal/api/handler/service.go`, `internal/api/router.go`, and `cmd/iot-api/main.go`.
- [X] T022 [US2] Implement downlink consumer flow that loads pending requests, updates lifecycle state idempotently, and handles replies in `internal/downlink/service.go` and `internal/downlink/dispatcher/dispatcher.go`.
- [X] T023 [US2] Align downlink-to-gateway routing keys and outbound payloads with `pkg/rabbitmq` helpers in `internal/downlink/router/router.go` and `internal/gateway/bridge/downlink.go`.
- [X] T024 [US2] Implement retry, timeout, dead-letter, and audit visibility for tracked requests in `internal/downlink/retry/retry.go`, `internal/downlink/service.go`, and `cmd/iot-downlink/main.go`.
- [X] T025 [US2] Expose service-request status lookup on the contract paths in `internal/api/handler/service.go`, `internal/api/router.go`, and `specs/004-core-services-implementation/contracts/openapi.yaml`.

**Checkpoint**: `iot-api` no longer bypasses `iot-downlink`, and request lifecycle state is fully tracked.

---

## Phase 5: User Story 3 - Manage device inventory and history (Priority: P2)

**Goal**: Support contract-aligned device CRUD and historical telemetry lookup for operations users.

**Independent Test**: Create and update a managed device, retrieve it by serial number, and query historical telemetry for a time range using the documented API contract.

### Tests for User Story 3 ⚠️

- [X] T026 [P] [US3] Add failing contract coverage for device inventory and telemetry history endpoints in `tests/contract/devices_contract_test.go`, `tests/contract/telemetry_contract_test.go`, and `specs/004-core-services-implementation/contracts/openapi.yaml`.
- [X] T027 [P] [US3] Add failing integration coverage for device CRUD and telemetry lookup by serial number in `tests/integration/api_inventory_test.go` and `tests/integration/message_flow_test.go`.
- [X] T028 [P] [US3] Add failing unit coverage for device and telemetry handlers in `internal/api/handler/device_test.go` and `internal/api/handler/telemetry_test.go`.

### Implementation for User Story 3

- [X] T029 [US3] Implement managed-device CRUD by serial number and contract-aligned not-found responses in `internal/api/handler/device.go` and `internal/api/router.go`.
- [X] T030 [US3] Extend managed-device persistence for gateway relationships, thing-model linkage, and status updates in `pkg/models/device.go` and `pkg/models/migrate.go`.
- [X] T031 [US3] Implement contract-aligned telemetry time-range queries for `/api/v1/devices/{deviceSN}/telemetry` in `internal/api/handler/telemetry.go` and `internal/api/router.go`.
- [X] T032 [US3] Wire API readiness to PostgreSQL and InfluxDB availability for inventory and history workflows in `cmd/iot-api/main.go` and `internal/api/router.go`.

**Checkpoint**: Operations users can manage device inventory and query consistent telemetry history through the documented API.

---

## Phase 6: User Story 4 - Receive realtime updates (Priority: P2)

**Goal**: Deliver authorized realtime updates over WebSocket with subscription management and heartbeat cleanup.

**Independent Test**: Open a realtime session, subscribe to supported device topics, generate matching activity, confirm only authorized updates are delivered, and verify stale clients are disconnected.

### Tests for User Story 4 ⚠️

- [X] T033 [P] [US4] Add failing contract coverage for realtime subscription registration in `tests/contract/realtime_contract_test.go` and `specs/004-core-services-implementation/contracts/openapi.yaml`.
- [X] T034 [P] [US4] Add failing integration coverage for subscription delivery and heartbeat cleanup in `tests/integration/ws_realtime_test.go` and `tests/integration/message_flow_test.go`.
- [X] T035 [P] [US4] Add failing unit coverage for topic authorization, fan-out, and heartbeat cleanup in `internal/ws/subscription/manager_test.go`, `internal/ws/push/pusher_test.go`, and `internal/ws/hub/client_test.go`.

### Implementation for User Story 4

- [X] T036 [US4] Implement authorized topic validation and scoped subscription management in `internal/ws/subscription/manager.go` and `internal/ws/service.go`.
- [X] T037 [US4] Implement canonical RabbitMQ bindings and message fan-out for realtime updates in `internal/ws/push/pusher.go` and `cmd/iot-ws/main.go`.
- [X] T038 [US4] Implement heartbeat-driven stale-session cleanup and realtime stats in `internal/ws/hub/client.go`, `internal/ws/hub/hub.go`, and `internal/ws/service.go`.
- [X] T039 [US4] Expose realtime subscription registration through the contract path in `internal/api/handler/realtime.go`, `internal/api/router.go`, and `specs/004-core-services-implementation/contracts/openapi.yaml`.
- [X] T040 [US4] Wire WebSocket readiness to RabbitMQ consumer health and active session state in `cmd/iot-ws/main.go` and `tests/integration/metrics_test.go`.

**Checkpoint**: Authorized clients receive realtime updates without polling, and stale sessions are cleaned up automatically.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and repository-wide quality gates.

- [X] T041 [P] Refresh API documentation and Swagger annotations for implemented contract paths in `cmd/iot-api/main.go`, `internal/api/handler/device.go`, `internal/api/handler/service.go`, `internal/api/handler/telemetry.go`, and `internal/api/handler/realtime.go`.
- [X] T042 [P] Add end-to-end validation for the quickstart flows in `tests/integration/message_flow_test.go`, `tests/integration/dji_e2e_event_test.go`, and `tests/integration/dji_e2e_service_test.go`.
- [X] T043 [P] Add resilience and performance coverage for telemetry latency, device concurrency, realtime fan-out, 99.9% availability, and 30-second recovery objectives in `tests/integration/dji_performance_test.go` and `tests/integration/metrics_test.go`.
- [X] T044 Run `make test`, `make lint`, and `make coverage`, then reconcile failures in `Makefile` and the specific Go files reported by those commands.
- [X] T045 Validate the documented startup and operational flow in `specs/004-core-services-implementation/quickstart.md` against the final service behavior and update any drift in `cmd/iot-gateway/main.go`, `cmd/iot-uplink/main.go`, `cmd/iot-downlink/main.go`, `cmd/iot-api/main.go`, and `cmd/iot-ws/main.go`.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1: Setup** — no dependencies.
- **Phase 2: Foundational** — depends on Phase 1 and blocks story completion.
- **Phase 3: US1** — depends on Phase 2; this is the MVP.
- **Phase 4: US2** — depends on Phase 2 and the canonical routing work from US1.
- **Phase 5: US3** — depends on Phase 2; telemetry history validation benefits from US1 but inventory CRUD is independently testable with seeded data.
- **Phase 6: US4** — depends on Phase 2 and downstream routing from US1.
- **Phase 7: Polish** — depends on all desired stories being complete.

### User Story Dependencies

- **US1 (P1)**: First delivery slice; no dependency on other stories.
- **US2 (P1)**: Requires the shared routing/message foundation and uses the same canonical gateway path validated by US1.
- **US3 (P2)**: Can proceed after Phase 2, but full telemetry-history verification depends on US1 data flow.
- **US4 (P2)**: Requires canonical uplink routing so subscribed updates can be delivered consistently.

### Within Each User Story

- Write contract, integration, and unit tests first and ensure they fail.
- Persistence/model work before service orchestration.
- Service orchestration before API/WebSocket endpoints.
- Readiness, metrics, and audit visibility before story sign-off.

### Parallel Opportunities

- `T002` and `T003` can run in parallel once `T001` is understood.
- `T005`, `T006`, and `T007` can run in parallel during the foundational phase.
- In each story, the listed test tasks can run in parallel.
- US3 can progress in parallel with US2 after Phase 2 if seeded telemetry data is acceptable for interim testing.

---

## Parallel Example: User Story 1

```bash
# Launch failing US1 tests together:
Task: "Add failing gateway authentication and connectivity integration coverage in tests/integration/message_flow_test.go and tests/integration/routing_test.go"
Task: "Add failing uplink normalization, malformed-message, and duplicate-transaction-id coverage in tests/integration/dji_osd_test.go and tests/integration/message_flow_test.go"
Task: "Add failing unit coverage for gateway auth/connection and uplink routing/storage in internal/gateway/mqtt/auth_test.go, internal/gateway/connection/manager_test.go, internal/uplink/router/router_test.go, and internal/uplink/storage/influx_test.go"
```

## Parallel Example: User Story 2

```bash
# Launch failing US2 tests together:
Task: "Add failing contract coverage for POST /api/v1/service-requests and GET /api/v1/service-requests/{requestID} in tests/contract/service_requests_contract_test.go and specs/004-core-services-implementation/contracts/openapi.yaml"
Task: "Add failing integration coverage for API-to-downlink-to-gateway request flow, timeout handling, late replies after timeout, and retry exhaustion in tests/integration/dji_e2e_service_test.go and tests/integration/message_flow_test.go"
Task: "Add failing unit coverage for request persistence, routing, retry behavior, and terminal-state idempotency in internal/api/handler/service_test.go, internal/downlink/router/router_test.go, internal/downlink/retry/retry_test.go, and internal/downlink/model/service_call_test.go"
```

## Parallel Example: User Story 3

```bash
# Launch failing US3 tests together:
Task: "Add failing contract coverage for device inventory and telemetry history endpoints in tests/contract/devices_contract_test.go, tests/contract/telemetry_contract_test.go, and specs/004-core-services-implementation/contracts/openapi.yaml"
Task: "Add failing integration coverage for device CRUD and telemetry lookup by serial number in tests/integration/api_inventory_test.go and tests/integration/message_flow_test.go"
Task: "Add failing unit coverage for device and telemetry handlers in internal/api/handler/device_test.go and internal/api/handler/telemetry_test.go"
```

## Parallel Example: User Story 4

```bash
# Launch failing US4 tests together:
Task: "Add failing contract coverage for realtime subscription registration in tests/contract/realtime_contract_test.go and specs/004-core-services-implementation/contracts/openapi.yaml"
Task: "Add failing integration coverage for subscription delivery and heartbeat cleanup in tests/integration/ws_realtime_test.go and tests/integration/message_flow_test.go"
Task: "Add failing unit coverage for topic authorization, fan-out, and heartbeat cleanup in internal/ws/subscription/manager_test.go, internal/ws/push/pusher_test.go, and internal/ws/hub/client_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup.
2. Complete Phase 2: Foundational.
3. Complete Phase 3: User Story 1.
4. Validate authenticated ingestion, normalized routing, and telemetry persistence before moving on.

### Incremental Delivery

1. Setup + Foundational establish the shared contracts.
2. Deliver **US1** for authenticated device ingress and telemetry.
3. Deliver **US2** for tracked downlink orchestration without API bypass.
4. Deliver **US3** for inventory and telemetry history APIs.
5. Deliver **US4** for realtime subscriptions and push delivery.
6. Finish with cross-cutting validation, performance, and quality gates.

### Parallel Team Strategy

1. Complete Setup + Foundational together.
2. After Phase 2:
   - Developer A: US1/US2 messaging path
   - Developer B: US3 inventory/history APIs
   - Developer C: US4 realtime delivery
3. Rejoin for Phase 7 validation and hardening.

---

## Notes

- `[P]` tasks touch separate files and can be split across contributors.
- Every user story remains independently testable even when it reuses foundational routing or storage work.
- The highest-risk items from the plan are explicitly covered: canonical routing alignment, removal of the `iot-api` downlink bypass, readiness based on real dependencies, and authorization-aware realtime subscriptions.
