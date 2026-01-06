<!--
Sync Impact Report:
- Version: 1.1.0 → 1.1.1
- Modified principles: I (明确要求遵循 Uber Go Style Guide)
- Added sections: None
- Removed sections: None
- Templates requiring updates: ✅ plan/spec/tasks templates already aligned; no further action
- Follow-up TODOs: None
-->
# IoT Platform Constitution

**Purpose**: Establish non-negotiable principles for the IoT platform feature (downlink/uplink/WS) covering quality, security, observability, and governance.

## Core Principles

### I. Code Quality & Test Gates (NON-NEGOTIABLE)
All changes MUST pass gofmt, go vet, staticcheck, and go test. New code MUST meet coverage: ≥80% overall, ≥90% on critical paths. Coding style MUST comply with the Uber Go Style Guide (formatting, naming, imports, error handling). No unused/dead code or unchecked errors. Contract/integration tests are REQUIRED when touching external/broker/storage APIs. PRs MUST include test evidence and peer review.

### II. Observability & Traceability
All services MUST emit structured logs (logrus) with trace/tenant/device fields, metrics to Prometheus, traces via W3C traceparent/OpenTelemetry to Tempo. Alerts MUST cover base + business SLOs and rate-limit events. Trace propagation is mandatory across MQTT/HTTP/WS/RabbitMQ.

### III. Security & Tenancy
MUST enforce device mTLS (optional username/password) and OIDC/JWT for operators/clients. Multi-tenant isolation is mandatory for data, auth, and quotas. Least privilege and audit logging MUST be applied to command/config actions.

### IV. Performance & Resilience
Downlink p95 ingest-to-persist <500ms; WS p95 fan-out <200ms; command delivery ≥99% for online devices. Rate limiting/backpressure per tenant/device is REQUIRED with clear rejection and telemetry. Replay window bounded (default 15 minutes). DR targets: RPO ≤15m, RTO ≤60m.

### V. Documentation & Traceability
Spec/plan/tasks MUST stay consistent with tech choices (Gin, logrus, VerneMQ, RabbitMQ, InfluxDB, PostgreSQL). Requirements and acceptance criteria MUST be measurable and traceable to tasks. Checklists act as requirement unit tests and MUST be maintained; deviations require documented rationale.

## Additional Constraints
Use single Go module `github.com/utmos/utmos`, Go 1.23+. Adhere to Uber Go Style Guide. Observability stack: Prometheus/Loki/Tempo/Grafana. API framework: Gin (REST) + gorilla/websocket.

## Workflow & Quality Gates
- Tests first for external contracts/integration changes.
- CI gates: fmt, vet, staticcheck, go test with coverage thresholds; contract/integration tests run where applicable.
- No user story work starts before foundational phase completion (config/auth/observability/clients/tests).
- Multi-tenant rate-limit/alert rules MUST be defined and validated before release.

## Governance
This constitution supersedes other practices. Amendments require version bump (SemVer): MAJOR for breaking/removal, MINOR for new principles/sections, PATCH for clarifications. Ratified items apply to all PR reviews; reviewers MUST check compliance with Core Principles and Workflow gates.

**Version**: 1.1.1 | **Ratified**: 2025-01-06 | **Last Amended**: 2025-01-06
