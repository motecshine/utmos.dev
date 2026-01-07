# Tasks: Project Setup with Distributed Tracing and Multi-Vendor RabbitMQ Routing

**Input**: Design documents from `/specs/001-project-setup/`
**Prerequisites**: plan.md, spec.md, data-model.md, research.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create monorepo directory structure per plan.md in repository root
- [X] T002 Initialize Go module with `go mod init` in repository root
- [X] T003 [P] Create go.mod and add core dependencies (Gin, GORM, RabbitMQ client, OpenTelemetry, Prometheus)
- [X] T004 [P] Configure golangci-lint in .golangci.yml
- [X] T005 [P] Configure misspell for typo checking in .golangci.yml
- [X] T006 [P] Create Makefile with build, test, lint commands
- [X] T007 [P] Setup GitHub Actions CI workflow in .github/workflows/ci.yml
- [X] T008 [P] Create docker-compose.yml for local development environment
- [X] T009 [P] Create README.md with project overview and setup instructions

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T010 [P] Create internal/shared/config package for YAML configuration management in internal/shared/config/config.go (support multi-environment: dev, staging, prod)
- [X] T011 [P] Create internal/shared/logger package for structured logging using logrus (github.com/sirupsen/logrus) in internal/shared/logger/logger.go
- [X] T012 [P] Create pkg/errors package for error definitions in pkg/errors/errors.go
- [X] T013 Create pkg/models/device.go with Device model (GORM)
- [X] T014 Create pkg/models/thing_model.go with ThingModel model (GORM)
- [X] T015 Create pkg/models/device_property.go with DeviceProperty model (GORM)
- [X] T016 Create pkg/models/device_event.go with DeviceEvent model (GORM)
- [X] T017 Create pkg/models/message_log.go with MessageLog model (GORM)
- [X] T018 Create database migration using GORM AutoMigrate (not SQL scripts) in pkg/models/ with migration initialization function
- [X] T019 Setup PostgreSQL connection configuration in internal/shared/config/database.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - 分布式追踪基础设施 (Priority: P1) 🎯 MVP

**Goal**: 实现分布式追踪基础设施，能够在所有微服务之间追踪消息流转

**Independent Test**: 发送一条设备消息，验证 trace_id 能够在所有服务（iot-gateway → iot-uplink → iot-api/iot-ws）之间传递，并在 Tempo 中查询到完整的调用链路

### Implementation for User Story 1

- [ ] T020 [US1] Create pkg/tracer/provider.go with OpenTelemetry Tracer Provider configuration
- [ ] T021 [US1] Create pkg/tracer/http.go with Gin HTTP tracing middleware (extract W3C Trace Context from headers)
- [ ] T022 [US1] Create pkg/tracer/rabbitmq.go with RabbitMQ message tracing (inject/extract W3C Trace Context in message headers)
- [ ] T023 [US1] Configure Tempo exporter in pkg/tracer/provider.go
- [ ] T024 [US1] Integrate HTTP tracing middleware into Gin router setup
- [ ] T025 [US1] Update internal/shared/logger to include trace_id in log entries
- [ ] T026 [US1] Integrate tracer with metrics package to include trace_id and span_id in metrics labels

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently - trace_id can be passed through HTTP and RabbitMQ

---

## Phase 4: User Story 2 - 多厂商 RabbitMQ Routing Key 定义 (Priority: P1)

**Goal**: 定义多厂商 RabbitMQ routing key 规范，支持不同厂商的消息路由

**Independent Test**: 模拟不同厂商（如 DJI、通用 MQTT）的消息，验证 routing key 能够正确路由到对应的处理服务

### Implementation for User Story 2

- [ ] T027 [US2] Create pkg/rabbitmq/routing.go with routing key generation function (format: `iot.{vendor}.{service}.{action}`)
- [ ] T028 [US2] Create pkg/rabbitmq/routing.go with routing key parsing function
- [ ] T029 [US2] Create pkg/rabbitmq/exchange.go with Topic Exchange (`iot`) management
- [ ] T030 [US2] Create pkg/rabbitmq/exchange.go with Queue creation and binding functions
- [ ] T031 [US2] Create pkg/repository/device.go with GetVendorByDeviceSN function to query vendor from database
- [ ] T032 [US2] Create pkg/rabbitmq/message.go with standard message format (includes device_sn, tid, bid, timestamp)
- [ ] T033 [US2] Create pkg/rabbitmq/client.go with RabbitMQ client wrapper
- [ ] T034 [US2] Implement routing key generation logic that queries vendor from database using device_sn

**Checkpoint**: At this point, User Story 2 should be fully functional - routing keys can be generated and parsed correctly for multi-vendor support

---

## Phase 5: User Story 3 - 统一 Metrics 包 (Priority: P1)

**Goal**: 实现统一的 metrics 包来处理框架基础中间件和业务代码的 metrics

**Independent Test**: 查看 Prometheus metrics 端点，验证框架中间件（RabbitMQ 连接数、PostgreSQL 查询延迟等）和业务 metrics（消息处理数量、错误率等）都能正确暴露

### Implementation for User Story 3

- [ ] T035 [US3] Create pkg/metrics/collector.go with Prometheus Registry management
- [ ] T036 [US3] Create pkg/metrics/middleware.go with RabbitMQ metrics collection (connection count, message count, latency)
- [ ] T037 [US3] Create pkg/metrics/middleware.go with PostgreSQL metrics collection (connection pool, query latency, error count)
- [ ] T038 [US3] Create pkg/metrics/middleware.go with InfluxDB metrics collection (write latency, error count)
- [ ] T039 [US3] Create pkg/metrics/business.go with Counter API (NewCounter function)
- [ ] T040 [US3] Create pkg/metrics/business.go with Histogram API (NewHistogram function)
- [ ] T041 [US3] Create pkg/metrics/business.go with Gauge API (NewGauge function)
- [ ] T042 [US3] Implement unified label specification (service, vendor, message_type, status, trace_id, span_id) in pkg/metrics/collector.go
- [ ] T043 [US3] Implement naming convention `iot_{component}_{metric_type}_{unit}` in pkg/metrics/collector.go
- [ ] T044 [US3] Create HTTP handler for `/metrics` endpoint in pkg/metrics/handler.go
- [ ] T045 [US3] Integrate `/metrics` endpoint into Gin router or create standalone HTTP server

**Checkpoint**: At this point, User Story 3 should be fully functional - metrics can be collected from middleware and business code, exposed at `/metrics` endpoint

---

## Phase 6: User Story 4 - 服务间消息串联机制 (Priority: P1)

**Goal**: 建立服务间消息串联机制，确保所有服务能够通过 RabbitMQ 和分布式追踪串联起来

**Independent Test**: 端到端测试，验证一条消息从设备到客户端响应的完整流程，所有服务都能正确处理和传递消息

### Implementation for User Story 4

- [ ] T046 [US4] Create pkg/rabbitmq/publisher.go with message publishing function (includes W3C Trace Context injection)
- [ ] T047 [US4] Create pkg/rabbitmq/subscriber.go with message subscription function (includes W3C Trace Context extraction)
- [ ] T048 [US4] Integrate W3C Trace Context injection into RabbitMQ message publishing in pkg/rabbitmq/publisher.go
- [ ] T049 [US4] Integrate W3C Trace Context extraction from RabbitMQ message headers in pkg/rabbitmq/subscriber.go
- [ ] T050 [US4] Create message validation function in pkg/rabbitmq/message.go (validate standard message format)
- [ ] T051 [US4] Create cmd/iot-api/main.go with service skeleton and RabbitMQ integration
- [ ] T052 [US4] Create cmd/iot-ws/main.go with service skeleton and RabbitMQ integration
- [ ] T053 [US4] Create cmd/iot-uplink/main.go with service skeleton and RabbitMQ integration
- [ ] T054 [US4] Create cmd/iot-downlink/main.go with service skeleton and RabbitMQ integration
- [ ] T055 [US4] Create cmd/iot-gateway/main.go with service skeleton, MQTT and RabbitMQ integration
- [ ] T056 [US4] Integrate tracing middleware and trace_id extraction in all service entry points (cmd/*/main.go)
- [ ] T057 [US4] Implement health check endpoints (`/health`, `/ready`) in all services
- [ ] T058 [US4] Implement graceful shutdown in all service main.go files

**Checkpoint**: At this point, User Story 4 should be fully functional - all services can communicate via RabbitMQ with distributed tracing

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T059 [P] Create integration tests for distributed tracing in tests/integration/test_tracing.go
- [ ] T060 [P] Create integration tests for RabbitMQ routing in tests/integration/test_routing.go
- [ ] T061 [P] Create integration tests for metrics collection in tests/integration/test_metrics.go
- [ ] T062 [P] Create unit tests for pkg/tracer package (coverage ≥ 80%)
- [ ] T063 [P] Create unit tests for pkg/rabbitmq package (coverage ≥ 80%)
- [ ] T064 [P] Create unit tests for pkg/metrics package (coverage ≥ 80%)
- [ ] T065 [P] Update API documentation in api/v1/openapi.yaml
- [ ] T066 [P] Update architecture documentation in docs/architecture/
- [ ] T067 [P] Create development guide in docs/development/
- [ ] T068 [P] Add Dockerfiles for all services in deployments/docker/
- [ ] T069 [P] Add Kubernetes manifests for all services in deployments/kubernetes/
- [ ] T070 Run quickstart.md validation to ensure all setup steps work correctly
- [ ] T071 Code cleanup and refactoring across all packages
- [ ] T072 Performance optimization for message processing and tracing

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (US1 → US2 → US3 → US4)
- **Polish (Phase 7)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - May use database models from Phase 2
- **User Story 3 (P1)**: Can start after Foundational (Phase 2) - Independent metrics implementation
- **User Story 4 (P1)**: Depends on US1 (tracing) and US2 (routing) - Integrates all components

### Within Each User Story

- Core infrastructure before integration
- Package creation before service integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, US1, US2, US3 can start in parallel
- US4 should start after US1 and US2 are complete (depends on tracing and routing)
- All tests in Polish phase marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tracer components in parallel:
Task: "Create pkg/tracer/provider.go with OpenTelemetry Tracer Provider configuration"
Task: "Create pkg/tracer/http.go with Gin HTTP tracing middleware"
Task: "Create pkg/tracer/rabbitmq.go with RabbitMQ message tracing"
```

---

## Parallel Example: User Story 2

```bash
# Launch routing and exchange components in parallel:
Task: "Create pkg/rabbitmq/routing.go with routing key generation function"
Task: "Create pkg/rabbitmq/exchange.go with Topic Exchange management"
Task: "Create pkg/rabbitmq/message.go with standard message format"
```

---

## Parallel Example: User Story 3

```bash
# Launch metrics components in parallel:
Task: "Create pkg/metrics/middleware.go with RabbitMQ metrics collection"
Task: "Create pkg/metrics/middleware.go with PostgreSQL metrics collection"
Task: "Create pkg/metrics/business.go with Counter API"
Task: "Create pkg/metrics/business.go with Histogram API"
Task: "Create pkg/metrics/business.go with Gauge API"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Distributed Tracing)
4. **STOP and VALIDATE**: Test User Story 1 independently - verify trace_id propagation
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 (Tracing) → Test independently → Deploy/Demo
3. Add User Story 2 (Routing) → Test independently → Deploy/Demo
4. Add User Story 3 (Metrics) → Test independently → Deploy/Demo
5. Add User Story 4 (Message Chaining) → Test independently → Deploy/Demo
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
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- All tasks follow strict checklist format with Task ID, parallel markers, story labels, and file paths

---

## Task Summary

**Total Tasks**: 72

**Tasks per Phase**:
- Phase 1 (Setup): 9 tasks
- Phase 2 (Foundational): 10 tasks
- Phase 3 (US1 - Tracing): 7 tasks
- Phase 4 (US2 - Routing): 8 tasks
- Phase 5 (US3 - Metrics): 11 tasks
- Phase 6 (US4 - Message Chaining): 13 tasks
- Phase 7 (Polish): 14 tasks

**Parallel Opportunities**: 
- Phase 1: 6 parallel tasks
- Phase 2: 8 parallel tasks
- Phase 3-5: Can run in parallel after Phase 2
- Phase 7: 13 parallel tasks

**Suggested MVP Scope**: Phase 1 + Phase 2 + Phase 3 (User Story 1 - Distributed Tracing)

