# Tasks: Project Setup with Distributed Tracing and Multi-Vendor RabbitMQ Routing

**Input**: Design documents from `/specs/001-project-setup/`
**Prerequisites**: plan.md, spec.md, data-model.md, types.md, research.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story. Each phase follows TDD principle: tests first, then implementation.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create monorepo directory structure per plan.md in repository root
- [X] T002 Initialize Go module with `go mod init github.com/utmos/utmos` in repository root
- [X] T003 [P] Add core dependencies to go.mod (Gin, GORM, RabbitMQ client, OpenTelemetry, Prometheus, logrus)
- [X] T004 [P] Configure golangci-lint with misspell checker in .golangci.yml
- [X] T005 [P] Create Makefile with build, test, lint, run commands
- [X] T006 [P] Setup GitHub Actions CI workflow in .github/workflows/ci.yml
- [X] T007 [P] Create docker-compose.yml for local development (PostgreSQL, RabbitMQ, InfluxDB, Tempo, Grafana)
- [X] T008 [P] Create README.md with project overview and setup instructions
- [X] T009 [P] Create config.dev.yaml and config.prod.yaml per types.md in configs/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### 2.1 Shared Infrastructure

- [X] T010 [P] Create internal/shared/config/config.go with Config struct per types.md Section 1
- [X] T011 [P] Create internal/shared/config/loader.go with Load() function for YAML config (multi-environment support)
- [X] T012 [P] Create internal/shared/logger/logger.go with logrus wrapper (JSON format, trace_id support)
- [X] T013 [P] Create pkg/errors/errors.go with ErrorCode enum and Error struct per types.md Section 4

### 2.2 Data Models (GORM)

- [X] T014 [P] Create pkg/models/device.go with Device model per data-model.md
- [X] T015 [P] Create pkg/models/thing_model.go with ThingModel model per data-model.md
- [X] T016 [P] Create pkg/models/device_property.go with DeviceProperty model per data-model.md
- [X] T017 [P] Create pkg/models/device_event.go with DeviceEvent model per data-model.md
- [X] T018 [P] Create pkg/models/message_log.go with MessageLog model per data-model.md
- [X] T019 Create pkg/models/migrate.go with AutoMigrate() function for all models

### 2.3 Database Connection

- [X] T020 Create internal/shared/database/postgres.go with GORM connection and pool configuration

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - 分布式追踪基础设施 (Priority: P1) 🎯 MVP

**Goal**: 实现分布式追踪基础设施，能够在所有微服务之间追踪消息流转

**Independent Test**: 发送一条设备消息，验证 trace_id 能够在所有服务之间传递，并在 Tempo 中查询到完整的调用链路

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T021 [P] [US1] Create pkg/tracer/provider_test.go with unit tests for TracerProvider
- [X] T022 [P] [US1] Create pkg/tracer/http_test.go with unit tests for HTTP middleware
- [X] T023 [P] [US1] Create pkg/tracer/rabbitmq_test.go with unit tests for context injection/extraction

### Implementation for User Story 1

- [X] T024 [US1] Create pkg/tracer/provider.go with OpenTelemetry TracerProvider per types.md Section 3.2 (depends on T021)
  - NewProvider(cfg *TracerConfig) function
  - Tempo OTLP HTTP exporter configuration
  - Sampling rate configuration (dev: 100%, prod: 10%)
  - Graceful shutdown support

- [X] T025 [US1] Create pkg/tracer/http.go with Gin HTTP tracing middleware per types.md Section 3.2 (depends on T022)
  - HTTPMiddleware(tracer trace.Tracer) gin.HandlerFunc
  - Extract W3C Trace Context from request headers (traceparent, tracestate)
  - Create span for each HTTP request
  - Add trace_id to response headers

- [X] T026 [US1] Create pkg/tracer/rabbitmq.go with RabbitMQ message tracing per types.md Section 3.2 (depends on T023)
  - InjectContext(ctx, headers) function
  - ExtractContext(ctx, headers) function
  - W3C Trace Context propagation in message headers

- [X] T027 [US1] Update internal/shared/logger/logger.go to include trace_id and span_id in log entries
  - WithTrace(ctx) function to extract trace context
  - Automatic trace_id injection in all log calls

**Checkpoint**: User Story 1 complete - trace_id can be passed through HTTP and RabbitMQ

---

## Phase 4: User Story 2 - 多厂商 RabbitMQ Routing Key 定义 (Priority: P1)

**Goal**: 定义多厂商 RabbitMQ routing key 规范，支持不同厂商的消息路由

**Independent Test**: 模拟不同厂商（DJI、通用 MQTT、Tuya）的消息，验证 routing key 能够正确生成和解析

### Tests for User Story 2

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T028 [P] [US2] Create pkg/rabbitmq/routing_test.go with unit tests for routing key generation and parsing
- [X] T029 [P] [US2] Create pkg/rabbitmq/message_test.go with unit tests for StandardMessage validation
- [X] T030 [P] [US2] Create pkg/rabbitmq/client_test.go with unit tests for RabbitMQ client
- [X] T031 [P] [US2] Create pkg/repository/device_test.go with unit tests for GetVendorByDeviceSN

### Implementation for User Story 2

- [X] T032 [US2] Create pkg/rabbitmq/routing.go with RoutingKey struct and functions per types.md Section 2.2 (depends on T028)
  - RoutingKey struct with Vendor, Service, Action fields
  - NewRoutingKey(vendor, service, action) function
  - Parse(key string) function
  - String() method returning `iot.{vendor}.{service}.{action}`
  - Predefined constants: VendorDJI, VendorGeneric, VendorTuya, ActionPropertyReport, etc.

- [X] T033 [US2] Create pkg/rabbitmq/message.go with StandardMessage struct per types.md Section 2.1 (depends on T029)
  - StandardMessage struct with TID, BID, Timestamp, Service, Action, DeviceSN, Data
  - MessageHeader struct with Traceparent, Tracestate, MessageType, Vendor
  - NewStandardMessage() function
  - Validate() method

- [X] T034 [US2] Create pkg/rabbitmq/client.go with RabbitMQ client per types.md Section 3.1 (depends on T030)
  - Client interface implementation
  - Connect() with exponential backoff retry (1s→2s→4s→8s...max 30s, 10 retries)
  - DeclareExchange(), DeclareQueue(), BindQueue() functions
  - IsConnected(), Close() functions

- [X] T035 [US2] Create pkg/rabbitmq/exchange.go with Exchange and Queue management
  - DeclareTopicExchange("iot") function
  - DeclareQueueWithDLQ() function (dead letter queue support)
  - BindQueueToExchange() function

- [X] T036 [US2] Create pkg/repository/device.go with DeviceRepository per types.md Section 3.4 (depends on T031)
  - GetByDeviceSN(ctx, deviceSN) function
  - GetVendorByDeviceSN(ctx, deviceSN) function for routing key generation
  - Create(), Update(), UpdateStatus() functions

**Checkpoint**: User Story 2 complete - routing keys can be generated and parsed correctly for multi-vendor support

---

## Phase 5: User Story 3 - 统一 Metrics 包 (Priority: P1)

**Goal**: 实现统一的 metrics 包来处理框架基础中间件和业务代码的 metrics

**Independent Test**: 查看 Prometheus metrics 端点，验证框架中间件和业务 metrics 都能正确暴露

### Tests for User Story 3

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T037 [P] [US3] Create pkg/metrics/collector_test.go with unit tests for Collector
- [X] T038 [P] [US3] Create pkg/metrics/business_test.go with unit tests for Counter, Histogram, Gauge APIs
- [X] T039 [P] [US3] Create pkg/metrics/middleware_test.go with unit tests for middleware metrics

### Implementation for User Story 3

- [X] T040 [US3] Create pkg/metrics/collector.go with Prometheus Registry management per types.md Section 3.3 (depends on T037)
  - Collector interface implementation
  - Registry() function
  - NewCollector(namespace string) function
  - Standard label constants: LabelService, LabelVendor, LabelMessageType, LabelStatus

- [X] T041 [US3] Create pkg/metrics/business.go with Counter, Histogram, Gauge APIs per types.md Section 3.3 (depends on T038)
  - NewCounter(name, help, labels) function
  - NewHistogram(name, help, labels, buckets) function
  - NewGauge(name, help, labels) function
  - Naming convention enforcement: `iot_{component}_{metric_type}_{unit}`

- [X] T042 [US3] Create pkg/metrics/middleware.go with middleware metrics collection (depends on T039)
  - RabbitMQ metrics: iot_rabbitmq_connection_total, iot_rabbitmq_message_total, iot_rabbitmq_message_duration_seconds
  - PostgreSQL metrics: iot_postgres_connection_pool_size, iot_postgres_query_duration_seconds, iot_postgres_error_total
  - InfluxDB metrics: iot_influxdb_write_duration_seconds, iot_influxdb_error_total

- [X] T043 [US3] Create pkg/metrics/handler.go with HTTP handler for /metrics endpoint
  - Handler(collector Collector) gin.HandlerFunc
  - Prometheus exposition format

**Checkpoint**: User Story 3 complete - metrics can be collected from middleware and business code, exposed at /metrics endpoint

---

## Phase 6: User Story 4 - 服务间消息串联机制 (Priority: P1)

**Goal**: 建立服务间消息串联机制，确保所有服务能够通过 RabbitMQ 和分布式追踪串联起来

**Independent Test**: 端到端测试，验证一条消息从设备到客户端响应的完整流程

**Dependencies**: US1 (tracing), US2 (routing), US3 (metrics)

### Tests for User Story 4

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T044 [P] [US4] Create pkg/rabbitmq/publisher_test.go with unit tests for message publishing
- [X] T045 [P] [US4] Create pkg/rabbitmq/subscriber_test.go with unit tests for message subscription
- [X] T046 [P] [US4] Create tests/integration/message_flow_test.go with integration test for message flow

### Implementation for User Story 4

- [X] T047 [US4] Create pkg/rabbitmq/publisher.go with message publishing per types.md Section 3.1 (depends on T044)
  - Publisher interface implementation
  - Publish(ctx, routingKey, msg) function
  - Automatic W3C Trace Context injection via pkg/tracer/rabbitmq.go

- [X] T048 [US4] Create pkg/rabbitmq/subscriber.go with message subscription per types.md Section 3.1 (depends on T045)
  - Subscriber interface implementation
  - Subscribe(queueName, handler) function
  - Manual Ack mode with Nack on error (dead letter queue)
  - Automatic W3C Trace Context extraction

### Service Skeletons

- [X] T049 [US4] Create cmd/iot-api/main.go with service skeleton
  - Config loading, logger, tracer, metrics initialization
  - Gin router with /health, /ready, /metrics endpoints
  - RabbitMQ publisher integration
  - Graceful shutdown

- [X] T050 [US4] Create cmd/iot-ws/main.go with service skeleton
  - Config loading, logger, tracer, metrics initialization
  - WebSocket connection management
  - RabbitMQ subscriber integration
  - Graceful shutdown

- [X] T051 [US4] Create cmd/iot-uplink/main.go with service skeleton
  - Config loading, logger, tracer, metrics initialization
  - RabbitMQ subscriber (from gateway) and publisher (to api/ws)
  - Graceful shutdown

- [X] T052 [US4] Create cmd/iot-downlink/main.go with service skeleton
  - Config loading, logger, tracer, metrics initialization
  - RabbitMQ subscriber (from api) and publisher (to gateway)
  - Graceful shutdown

- [X] T053 [US4] Create cmd/iot-gateway/main.go with service skeleton
  - Config loading, logger, tracer, metrics initialization
  - MQTT client connection to VerneMQ
  - RabbitMQ publisher (to uplink) and subscriber (from downlink)
  - MQTT ↔ RabbitMQ message conversion
  - Graceful shutdown

- [X] T054 [US4] Create internal/shared/server/graceful.go with graceful shutdown helper
  - WaitForShutdown(ctx, timeout) function
  - Signal handling (SIGINT, SIGTERM)
  - Resource cleanup order

**Checkpoint**: User Story 4 complete - all services can communicate via RabbitMQ with distributed tracing

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

### Integration Tests

- [X] T055 [P] Create tests/integration/tracing_test.go with end-to-end tracing test
- [X] T056 [P] Create tests/integration/routing_test.go with multi-vendor routing test
- [X] T057 [P] Create tests/integration/metrics_test.go with metrics collection test

### Documentation

- [X] T058 [P] Update api/v1/openapi.yaml with service endpoints
- [X] T059 [P] Create docs/architecture/overview.md with architecture documentation
- [X] T060 [P] Create docs/development/getting-started.md with development guide

### Deployment

- [X] T061 [P] Create deployments/docker/iot-api.Dockerfile
- [X] T062 [P] Create deployments/docker/iot-ws.Dockerfile
- [X] T063 [P] Create deployments/docker/iot-uplink.Dockerfile
- [X] T064 [P] Create deployments/docker/iot-downlink.Dockerfile
- [X] T065 [P] Create deployments/docker/iot-gateway.Dockerfile

### Validation

- [ ] T066 Run quickstart.md validation to ensure all setup steps work correctly
- [X] T067 Run golangci-lint to verify code quality and no typos
- [ ] T068 Verify test coverage >= 80% for all packages

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Foundational phase completion
  - US1, US2, US3 can proceed in parallel after Phase 2
  - US4 depends on US1 (tracing) and US2 (routing) completion
- **Polish (Phase 7)**: Depends on all user stories being complete

### User Story Dependencies

```
Phase 2 (Foundational)
         │
         ├──────────────┬──────────────┐
         │              │              │
         ▼              ▼              ▼
    US1 (Tracing)  US2 (Routing)  US3 (Metrics)
         │              │              │
         └──────┬───────┘              │
                │                      │
                ▼                      │
           US4 (Message Chain) ◄───────┘
                │
                ▼
         Phase 7 (Polish)
```

### Within Each User Story (TDD Order)

1. **Tests First**: Write tests, ensure they FAIL
2. **Implementation**: Implement to make tests pass
3. **Refactor**: Clean up code while keeping tests green

### Parallel Opportunities

- **Phase 1**: All tasks marked [P] can run in parallel
- **Phase 2**: All tasks in 2.1, 2.2, 2.3 can run in parallel within each section
- **Phase 3-5**: US1, US2, US3 can run in parallel after Phase 2
- **Phase 6**: US4 must wait for US1 and US2
- **Phase 7**: All tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests first (TDD):
Task: "Create pkg/tracer/provider_test.go"
Task: "Create pkg/tracer/http_test.go"
Task: "Create pkg/tracer/rabbitmq_test.go"

# Then implement (after tests fail):
Task: "Create pkg/tracer/provider.go" (after provider_test.go)
Task: "Create pkg/tracer/http.go" (after http_test.go)
Task: "Create pkg/tracer/rabbitmq.go" (after rabbitmq_test.go)
```

---

## Parallel Example: User Story 2

```bash
# Launch all tests first (TDD):
Task: "Create pkg/rabbitmq/routing_test.go"
Task: "Create pkg/rabbitmq/message_test.go"
Task: "Create pkg/rabbitmq/client_test.go"
Task: "Create pkg/repository/device_test.go"

# Then implement (after tests fail):
Task: "Create pkg/rabbitmq/routing.go"
Task: "Create pkg/rabbitmq/message.go"
Task: "Create pkg/rabbitmq/client.go"
Task: "Create pkg/repository/device.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Distributed Tracing)
4. **STOP and VALIDATE**: Test trace_id propagation independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 (Tracing) → Test independently → Deploy/Demo
3. Add User Story 2 (Routing) → Test independently → Deploy/Demo
4. Add User Story 3 (Metrics) → Test independently → Deploy/Demo
5. Add User Story 4 (Message Chaining) → Requires US1+US2 → Deploy/Demo
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Tracing)
   - Developer B: User Story 2 (Routing)
   - Developer C: User Story 3 (Metrics)
3. After US1 and US2 complete:
   - Developer A + B: User Story 4 (Message Chaining)
4. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- **TDD Required**: Write tests FIRST, ensure they FAIL, then implement
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- All type definitions reference types.md to prevent implementation guessing
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence

---

## Task Summary

**Total Tasks**: 68

**Tasks per Phase**:
- Phase 1 (Setup): 9 tasks
- Phase 2 (Foundational): 11 tasks
- Phase 3 (US1 - Tracing): 7 tasks (3 tests + 4 implementation)
- Phase 4 (US2 - Routing): 9 tasks (4 tests + 5 implementation)
- Phase 5 (US3 - Metrics): 7 tasks (3 tests + 4 implementation)
- Phase 6 (US4 - Message Chaining): 11 tasks (3 tests + 8 implementation)
- Phase 7 (Polish): 14 tasks

**Parallel Opportunities**:
- Phase 1: 7 parallel tasks
- Phase 2: 10 parallel tasks
- Phase 3-5: Can run in parallel after Phase 2 (US1, US2, US3 independent)
- Phase 7: 11 parallel tasks

**Suggested MVP Scope**: Phase 1 + Phase 2 + Phase 3 (User Story 1 - Distributed Tracing)

**Key Improvements over Previous Version**:
1. ✅ TDD Compliant: Tests before implementation in each user story
2. ✅ No File Conflicts: Each task operates on a single file
3. ✅ Type References: All tasks reference types.md for implementation details
4. ✅ Clear Dependencies: US4 explicitly depends on US1 and US2
5. ✅ Merged Conflicting Tasks: Previous T036-T038 and T039-T041 merged into single tasks

