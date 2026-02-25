# Implementation Plan: DJI Protocol Implementation

**Branch**: `003-dji-protocol-implementation` | **Date**: 2025-02-05 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-dji-protocol-implementation/spec.md`

## Summary

Complete implementation of the DJI Cloud API protocol adapter for the UMOS IoT platform. The adapter translates all DJI uplink and downlink MQTT message types (OSD, State, Status, Services, Events, Requests, DRC) into the platform's StandardMessage format via a handler/router architecture. Built on top of the protocol adapter framework established in `002-protocol-adapter-design`, the implementation covers 12 protocol modules (aircraft, camera, config, common, device, drc, file, firmware, live, psdk, safety, wayline) with 50+ registered service/event methods, organized across 78 tasks in 13 phases.

## Technical Context

**Language/Version**: Go 1.22+ (module declares go 1.25.5)
**Primary Dependencies**: Gin v1.11.0 (HTTP/health), GORM v1.31.1 (PostgreSQL ORM), amqp091-go v1.10.0 (RabbitMQ), logrus v1.9.4 (structured logging), OpenTelemetry v1.40.0 (distributed tracing), Prometheus client_golang v1.23.2 (metrics), google/uuid v1.6.0 (transaction IDs), testify v1.11.1 (test assertions), nbio/xml (WPML parsing)
**Storage**: PostgreSQL (business data, device metadata, thing model definitions via GORM), InfluxDB v2 (time-series telemetry data via influxdb-client-go)
**Testing**: `go test` with `testify/assert` and `testify/require`; TDD workflow (tests written first, must fail before implementation); race detection enabled; target coverage >= 80%
**Target Platform**: Linux server (Docker containers, Kubernetes orchestration)
**Project Type**: Microservices (6 services: iot-api, iot-ws, iot-uplink, iot-downlink, iot-gateway, dji-adapter)
**Performance Goals**: <50ms P95 message processing latency, 1000 concurrent devices, >99.9% message processing success rate
**Constraints**: MQTT isolation (only iot-gateway connects to VerneMQ), all inter-service communication via RabbitMQ (no HTTP/gRPC between services), RabbitMQ routing keys follow `iot.{vendor}.{service}.{action}` format, all messages must include `device_sn`, `tid`, `bid`, `timestamp`, and W3C Trace Context headers
**Scale/Scope**: 12 protocol modules, 50+ service/event methods, 78 tasks across 13 phases, 9 user stories covering P1-P3 priorities

## Constitution Check

*GATE: Verified against UMOS IoT Platform Constitution v1.4.1*

| Principle | Status | Evidence |
|-----------|--------|----------|
| **I. Thing Model Driven (NON-NEGOTIABLE)** | Pass | TSL-based protocol modules in `pkg/adapter/dji/protocol/` with Properties (OSD/State), Services (50+ methods), and Events (core/wayline/drc/file/firmware). Protocol types defined as Go structs matching DJI Cloud API TSL definitions. |
| **II. Multi-Protocol Support** | Pass | MQTT topics parsed via `ParseTopic()` (thing/product/{sn}/osd, sys/product/{sn}/status, etc.), HTTP health/metrics endpoints via Gin, WebSocket service in `internal/ws/`. |
| **III. Device Abstraction Layer** | Pass | Handler/Router pattern abstracts DJI-specific protocol into `MessageHandler` interface. `HandlerRegistry` dispatches by `TopicType`. Gateway device pattern supported via `TopicInfo.GatewaySN`/`DeviceSN` separation. |
| **IV. Extensibility & Plugin Architecture** | Pass | Plugin-style handler registration via `Registry.Register(handler)` and `ServiceRouter.RegisterServiceHandler(method, fn)`. New commands added by implementing `HandlerFunc` and calling registration functions (e.g., `RegisterDeviceCommands`, `RegisterWaylineCommands`). |
| **V. Standardized API Design** | Pass | RESTful API via Gin in `internal/api/`, health/readiness/metrics endpoints at `/health`, `/ready`, `/metrics`. Unified `StandardMessage` format with `ProtocolMeta` for all inter-service messages. |
| **VI. Test-First Development (NON-NEGOTIABLE)** | Pass | TDD workflow enforced: 78 tasks organized with tests-first phases. Every handler, router, and protocol module has corresponding `*_test.go` files. Integration tests in `tests/integration/`. Benchmark tests in `handler/benchmark_test.go`. |
| **VII. Observability & Monitoring** | Pass | `pkg/adapter/dji/observability/` provides Prometheus metrics (`Metrics`), structured logging (`Logger` with trace context, vendor, device_sn fields), OpenTelemetry tracing (`Tracer` with span attributes for messaging.system, dji.method, dji.device_sn), and unified `HandlerObserver` for handler instrumentation. |

No constitution violations detected.

## Project Structure

### Documentation (this feature)

```text
specs/003-dji-protocol-implementation/
├── plan.md              # This file
├── spec.md              # Feature specification (DJI Cloud API protocol requirements)
└── tasks.md             # 78 tasks across 13 phases (all completed)
```

### Source Code (repository root)

```text
cmd/
├── dji-adapter/          # DJI adapter service entry point (main.go)
│                         # - RabbitMQ consumer for uplink/downlink queues
│                         # - Gin HTTP server for /health, /ready, /metrics
│                         # - Prometheus metrics (messages_processed, processing_duration, parse_errors)
│                         # - Graceful shutdown with signal handling
├── iot-api/              # HTTP RESTful API service
├── iot-downlink/         # Downlink message processor
├── iot-gateway/          # MQTT gateway (VerneMQ bridge)
├── iot-uplink/           # Uplink message processor
└── iot-ws/               # WebSocket real-time push service

pkg/adapter/dji/
├── adapter.go            # Core Adapter implementing ProtocolAdapter interface
│                         # - HandleMessage: topic parsing -> handler dispatch -> fallback converter
│                         # - ParseRawMessage, ToStandardMessage, FromStandardMessage, GetRawPayload
│                         # - MessageHandler and HandlerRegistry interfaces
├── converter.go          # Bidirectional DJI <-> StandardMessage conversion
│                         # - MapTopicTypeToAction, MapActionToTopicType mappings
├── parser.go             # DJI message JSON parsing (ParseMessage)
├── topic.go              # MQTT topic parsing (thing/product/{sn}/osd, sys/product/{sn}/status)
│                         # - ParseTopic, BuildTopic, TopicInfo struct
├── types.go              # Core types: Message, TopicType, Direction, action constants
├── errors.go             # Sentinel errors (ErrMissingTID, ErrEmptyTopic, ErrInvalidTopic, etc.)
├── handler/              # Message handlers by topic type
│   ├── handler.go        # Handler interface definition
│   ├── base_handler.go   # Shared handler logic (MessageConfig, BuildStandardMessage, HandleMessage)
│   ├── registry.go       # Thread-safe HandlerRegistry (map[TopicType]Handler with RWMutex)
│   ├── osd_handler.go    # OSD handler with OSDParser integration (aircraft/dock/rc field extraction)
│   ├── state_handler.go  # State handler for property change messages
│   ├── status_handler.go # Status handler for device online/offline topology
│   ├── service_handler.go# Service handler delegating to ServiceRouter
│   ├── routed_handlers.go# EventHandler and RequestHandler (baseHandler + router delegation)
│   ├── drc_handler.go    # DRC real-time control handler (drc/up, drc/down topics)
│   └── benchmark_test.go # Performance benchmarks for handler processing
├── router/               # Method-level routing for services and events
│   ├── router.go         # Generic handlerRegistry[T] with thread-safe register/get/list/has
│   ├── handler.go        # SimpleCommandHandler[T], NoDataCommandHandler, SimpleEventHandler[T]
│   │                     # RegisterHandlers, RegisterEventHandlers batch registration helpers
│   ├── service_router.go # ServiceRouter: method -> ServiceHandlerFunc dispatch
│   ├── event_router.go   # EventRouter: method -> EventHandlerFunc dispatch (with need_reply)
│   ├── device_commands.go    # 16 device commands (cover, drone, charge, reboot, format, debug, etc.)
│   ├── camera_commands.go    # 9 camera commands (mode, photo, recording, aim, focal, gimbal, IR)
│   ├── wayline_commands.go   # 8 wayline commands (create, prepare, execute, pause, recovery, undo, return)
│   ├── drc_commands.go       # 5 DRC commands (mode_enter/exit, drone_control, emergency_stop, heart)
│   ├── file_firmware_commands.go # File commands (upload start/finish/list) + firmware (ota_create)
│   ├── config_live_commands.go   # Config commands + live stream commands (start/stop/quality/lens)
│   ├── core_events.go       # Core events (device_exit_homing, temp_notify, file_upload_callback, hms)
│   ├── wayline_events.go    # Wayline events (flighttask_progress, flighttask_ready, return_home_info)
│   ├── drc_events.go        # DRC events (joystick_invalid_notify, drc_status_notify)
│   ├── file_events.go       # File events (highest_priority_upload, file_upload_progress)
│   └── firmware_events.go   # Firmware events (ota_progress)
├── protocol/             # DJI Cloud API protocol data types (TSL Go struct definitions)
│   ├── aircraft/         # Aircraft OSD, Dock OSD, RC OSD, HMS events
│   ├── camera/           # Camera commands, IR camera types
│   ├── common/           # Common header, command, types (shared across modules)
│   ├── config/           # Device config, organization, requests, storage
│   ├── device/           # Device commands, events
│   ├── drc/              # DRC commands, events
│   ├── file/             # File commands, events, log events
│   ├── firmware/         # Firmware commands, events
│   ├── live/             # Live stream commands
│   ├── psdk/             # PSDK payload commands, events
│   ├── safety/           # Safety commands
│   └── wayline/          # Wayline commands, events, requests, types
├── integration/          # OSD parser with auto-detection (aircraft/dock/rc)
│   └── osd_parser.go     # ParseOSD, ParsedOSD with type detection heuristics
├── observability/        # Observability instrumentation
│   ├── metrics.go        # Prometheus metrics (received/sent/errors counters, processing histogram, active gauge)
│   ├── logger.go         # Structured logger with vendor/device_sn/trace context fields
│   ├── tracer.go         # OpenTelemetry tracer (message spans, service call spans, event spans)
│   └── handler_observer.go # Unified handler observer (StartObserve -> End/EndWithError)
├── config/               # Configuration constants
│   └── config.go         # SERVICE_CALL_TIMEOUT=30s, DRC_HEARTBEAT_TIMEOUT=3s, UNKNOWN_DEVICE_POLICY=discard
├── wpml/                 # WPML wayline file handling
│   ├── types.go          # WPML XML types
│   ├── structures.go     # Wayline structure definitions
│   ├── template.go       # Template handling
│   ├── placemark.go      # Placemark/waypoint types
│   ├── actions.go        # Action definitions
│   ├── action_types.go   # Action type enums
│   ├── converter.go      # WPML <-> internal conversion
│   ├── converter_schema.go # Schema conversion helpers
│   ├── serializer.go     # XML serialization
│   ├── validator.go      # WPML validation logic
│   ├── kmz.go            # KMZ file read/write
│   └── errors.go         # WPML-specific errors
├── uplink/               # Uplink processing
│   ├── adapter.go        # Uplink adapter
│   ├── processor.go      # Message processor
│   └── osd.go            # OSD-specific uplink processing
├── downlink/             # Downlink processing
│   ├── adapter.go        # Downlink adapter
│   ├── dispatcher.go     # Service call dispatcher (routing key generation, StandardMessage creation)
│   └── methods.go        # Downlink method definitions
└── init/                 # Handler initialization
    └── init.go           # InitializeAdapter: registers all service commands, events, and handlers
                          # - registerServiceCommands: device, camera, wayline, drc, file/firmware, config/live
                          # - registerEvents: core, wayline, drc, file, firmware
                          # - registerHandlers: OSD, State, Status, Service, Event, Request, DRC

internal/
├── api/                  # HTTP API service internals
│   ├── router.go         # Gin router setup with middleware
│   ├── handler/          # API handlers (device, service, telemetry)
│   └── middleware/       # Auth and trace middleware
├── downlink/             # Downlink service internals
│   ├── service.go        # Downlink service orchestration
│   ├── dispatcher/       # Message dispatcher
│   ├── router/           # Downlink routing
│   ├── retry/            # Retry logic with backoff
│   └── model/            # Service call model
├── uplink/               # Uplink service internals
│   ├── service.go        # Uplink service orchestration
│   ├── processor/        # Message processor
│   ├── router/           # Uplink routing
│   └── storage/          # InfluxDB time-series storage
├── gateway/              # Gateway service internals
│   ├── service.go        # Gateway service orchestration
│   ├── bridge/           # MQTT <-> RabbitMQ bridge (uplink, downlink, parse)
│   ├── connection/       # Device connection manager
│   ├── mqtt/             # MQTT client, auth, handler
│   └── model/            # Credential model
├── ws/                   # WebSocket service internals
│   ├── service.go        # WebSocket service orchestration
│   ├── hub/              # WebSocket hub and client management
│   ├── push/             # Message push logic
│   └── subscription/     # Subscription manager
└── shared/               # Shared infrastructure
    ├── config/           # Configuration loading
    ├── database/         # PostgreSQL connection setup
    └── server/           # Graceful shutdown helper

pkg/                      # Shared public packages
├── adapter/              # Protocol adapter framework (ProtocolAdapter interface, global registry)
├── config/               # Package-level configuration types
├── errors/               # Shared error types
├── logger/               # Structured logger with trace context support
├── metrics/              # Prometheus metrics collector (NewCounter, NewHistogram, NewGauge)
├── models/               # Shared domain models
├── rabbitmq/             # RabbitMQ client, publisher, StandardMessage, ProtocolMeta, routing keys
├── registry/             # Generic registry utilities
├── repository/           # Repository pattern interfaces
└── tracer/               # OpenTelemetry tracer setup

tests/
└── integration/          # Integration and E2E tests
    ├── api_test.go
    ├── availability_test.go
    ├── downlink_test.go
    ├── e2e_test.go
    ├── gateway_test.go
    ├── load_device_test.go
    ├── load_ws_test.go
    ├── performance_test.go
    ├── uplink_test.go
    └── ws_test.go
```

**Structure Decision**: Go microservices layout with `cmd/` for service entry points, `pkg/` for reusable packages (including the DJI adapter as a vendor-specific protocol plugin), and `internal/` for service-specific logic. The DJI adapter (`pkg/adapter/dji/`) follows a layered architecture: protocol types define the data model, handlers process messages by topic type, routers dispatch by method name, and the init package wires everything together. This structure enables independent testing of each layer and straightforward extension for new DJI methods or new vendor adapters.

## Complexity Tracking

No constitution violations. All architecture decisions align with the UMOS IoT Platform Constitution v1.4.1 principles. The 6th service (dji-adapter) extends the 5-service microservice constraint defined in the constitution, but this is an expected outcome of the plugin architecture principle (Principle IV) -- vendor-specific protocol adapters are designed to be deployed as independent services that communicate exclusively via RabbitMQ.
