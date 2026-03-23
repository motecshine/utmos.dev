# Feature Specification: Core Services Implementation

**Feature Branch**: `004-core-services-implementation`
**Created**: 2025-02-05
**Updated**: 2026-03-23
**Status**: Draft
**Input**: User description: "implement core services"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Operate connected devices (Priority: P1)

As a platform operator, I want supported devices to connect to the platform, report their status, and publish telemetry so I can monitor live operations from a single system.

**Why this priority**: Without reliable device onboarding and upstream data handling, the platform cannot deliver its core value.

**Independent Test**: Can be fully tested by connecting a supported device or simulator, sending valid status and telemetry messages, and confirming that current device state and recent telemetry become visible to operators.

**Acceptance Scenarios**:

1. **Given** a supported device with valid credentials, **When** it connects and publishes status data, **Then** the platform records the device as online and makes the updated status available to operators.
2. **Given** a connected device publishes valid telemetry, **When** the platform processes the message, **Then** the telemetry is stored and made available for historical lookup and downstream consumption.
3. **Given** a device publishes a malformed or unsupported message, **When** the platform receives it, **Then** the message is rejected safely, the failure is traceable, and other device traffic continues to flow.

---

### User Story 2 - Send commands to devices (Priority: P1)

As a platform operator, I want to send service requests to a device and track whether each request succeeds, fails, or times out so I can control field equipment with confidence.

**Why this priority**: Command execution is the primary control loop for managed devices and is required for meaningful platform operations.

**Independent Test**: Can be fully tested by submitting a service request for a supported device, observing immediate acknowledgement, and verifying that the final request outcome is recorded and visible.

**Acceptance Scenarios**:

1. **Given** an authorized client submits a valid service request for an online device, **When** the platform accepts the request, **Then** it returns a tracking identifier and records the request as pending.
2. **Given** a device acknowledges and completes a service request, **When** the platform receives the response, **Then** the final outcome is recorded and made available to the requesting client.
3. **Given** a device does not respond within the allowed window, **When** retry attempts are exhausted, **Then** the platform marks the request as timed out or failed and preserves the failure details for review.
4. **Given** a service request has already reached a terminal timeout outcome, **When** a late device response arrives with the same tracking identifiers, **Then** the platform records the late response for audit purposes without reopening or overwriting the terminal state.

---

### User Story 3 - Manage device inventory and history (Priority: P2)

As an operations user, I want to manage device records and review historical telemetry so I can keep the fleet inventory accurate and investigate device behavior over time.

**Why this priority**: Operational teams need authoritative device records and historical context, but the platform still provides baseline value before this workflow is added.

**Independent Test**: Can be fully tested by creating and updating a device record, retrieving device details, and querying recent telemetry history for a supported device.

**Acceptance Scenarios**:

1. **Given** an authorized client creates or updates a device record, **When** the request is validated, **Then** the platform stores the change and returns the current device data.
2. **Given** telemetry has been collected for a device, **When** an authorized client requests historical data for a time range, **Then** the platform returns the matching records in a consistent format.
3. **Given** a client requests a device that does not exist, **When** the platform processes the request, **Then** it returns a clear not-found response without affecting other records.

---

### User Story 4 - Receive realtime updates (Priority: P2)

As a monitoring client, I want to subscribe to realtime device updates so I can react to status changes, telemetry, and events without polling.

**Why this priority**: Realtime visibility improves operator responsiveness, but it builds on the core ingestion and routing capabilities delivered earlier.

**Independent Test**: Can be fully tested by opening a realtime session, subscribing to a device topic, generating device activity, and confirming only matching updates are delivered.

**Acceptance Scenarios**:

1. **Given** an authorized client opens a realtime session, **When** it subscribes to supported device topics, **Then** the platform confirms the subscription and begins delivering matching updates.
2. **Given** subscribed devices publish new status, telemetry, or event data, **When** the platform processes those updates, **Then** subscribed clients receive the updates in near real time.
3. **Given** a realtime client becomes unresponsive, **When** the heartbeat window is missed, **Then** the platform closes the stale session and releases its subscriptions.

### Edge Cases

- Duplicate upstream messages or duplicate command responses reusing the same `tid` and `bid` are handled idempotently per FR-011 — no duplicate records or duplicate state transitions are created.
- Late command responses arriving after a terminal timeout outcome are recorded for audit purposes without reopening or overwriting the terminal state (per US2-Scenario-4).
- When many devices reconnect simultaneously after a broker or network interruption, the gateway applies a token-bucket rate limit on incoming connections; excess MQTT connects receive a CONNACK error and rely on standard MQTT client reconnect backoff.
- When a realtime client subscribes to device topics it is not authorized to observe, the platform validates each topic individually, returns accepted topics in the subscription ack, and returns rejected topics with an authorization error. Valid topics in the same request still activate.
- When a storage or messaging dependency is temporarily unavailable while commands or telemetry are in flight, services nack and requeue messages up to a configurable retry ceiling, then route exhausted messages to a dead-letter queue for operator review.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The platform MUST authenticate supported devices before allowing them to exchange operational messages.
- **FR-002**: The platform MUST maintain current connectivity status for each managed device and update that status when devices connect, disconnect, or stop sending expected heartbeats.
- **FR-003**: The platform MUST accept upstream device messages and normalize them into a standard platform message model that preserves transaction identity, business identity, event time, device identity, vendor context, and trace context.
- **FR-004**: The platform MUST validate supported upstream messages against the registered device capability model before storing or forwarding business data.
- **FR-005**: The platform MUST store supported telemetry and event data so authorized clients can review device history.
- **FR-006**: The platform MUST route validated upstream updates to the appropriate downstream consumers, including realtime delivery channels and business-facing interfaces.
- **FR-007**: The platform MUST allow authorized clients to create, view, update, and remove managed device records.
- **FR-008**: The platform MUST allow authorized clients to submit service requests for supported devices and immediately receive a tracking identifier for each accepted request.
- **FR-009**: The platform MUST track each device service request through the authoritative lifecycle states `pending`, `sent`, `retrying`, `success`, `failed`, and `timeout`.
- **FR-010**: The platform MUST retry unacknowledged or failed service deliveries according to a defined retry policy and preserve the final failure outcome for operator review.
- **FR-011**: The platform MUST process duplicate upstream messages and duplicate service responses idempotently when they repeat the same `tid` and `bid`, so duplicate deliveries do not create duplicate telemetry/history records or duplicate terminal service-request transitions.
- **FR-012**: The platform MUST deliver downlink device requests only through the designated gateway path rather than direct service-to-device calls.
- **FR-013**: The platform MUST allow authorized clients to establish realtime sessions, manage subscriptions to supported device topics, and receive matching updates.
- **FR-014**: The platform MUST detect stale realtime sessions and close them after missed heartbeat expectations.
- **FR-015**: The platform MUST provide integration-ready interface documentation for device management, telemetry retrieval, service requests, and realtime subscriptions.
- **FR-016**: The platform MUST expose service health and operational readiness information for each core service so operators can determine whether the end-to-end flow is available.
- **FR-017**: The platform MUST record auditable operational events for authentication attempts, message processing failures, retries, and final command outcomes.

### Non-Functional Requirements

- **NFR-001**: At least 95% of valid device telemetry updates MUST become available for subscribed clients or downstream consumers within 1 second during representative acceptance testing.
- **NFR-002**: The platform MUST support at least 1,000 simultaneously connected devices during controlled load testing without dropping authenticated sessions.
- **NFR-003**: The platform MUST support at least 10,000 simultaneous realtime client connections per node during controlled load testing.
- **NFR-004**: Core platform services MUST recover from transient messaging or storage dependency interruptions within 30 seconds without manual intervention.
- **NFR-005**: The end-to-end operational flow MUST demonstrate at least 99.9% availability during controlled resilience testing windows.
- **NFR-006**: Release readiness MUST require passing unit, integration, end-to-end, and contract tests, with unit test coverage of at least 80%.

### Key Entities *(include if feature involves data)*

- **Managed Device**: A supported field asset known to the platform, including its identity, type, vendor, status, connectivity state, and relationship to any parent gateway.
- **Device Credential**: The authentication record that determines whether a device is allowed to connect and exchange messages with the platform.
- **Capability Model**: The normalized description of a device's properties, services, and events that the platform uses to validate and interpret device traffic.
- **Telemetry Record**: A timestamped snapshot of reported device measurements or state used for monitoring, history, and downstream processing.
- **Service Request**: A tracked command issued to a device, including its target device, requested operation, status history, retries, and final outcome.
- **Realtime Subscription**: A client-managed registration describing which device topics should be delivered over an active realtime session.

### Assumptions

- Existing platform foundations already provide supported device protocol translation and baseline infrastructure for core services to build upon.
- The initial scope targets supported vendor/device families already modeled by the platform's capability definitions.
- Clients interacting with device management, command, and realtime features already use the platform's standard authorization model.
- This feature covers the five core runtime services and their end-to-end behavior, not new device family onboarding.

## Clarifications

### Session 2026-03-23

- Q: When a storage or messaging dependency is temporarily unavailable during in-flight processing, what should the service do? → A: Nack + requeue with max retries → dead-letter queue for operator review.
- Q: When many devices reconnect simultaneously after a broker/network interruption, how should the gateway handle the surge? → A: Token-bucket rate limit; excess MQTT connects get CONNACK error and use standard client backoff.
- Q: When a realtime client subscribes to device topics it is not authorized to observe, how should the platform respond? → A: Per-topic accept/reject; valid topics activate, rejected topics return authorization error.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Operators can connect a supported device and observe its online status in the platform within 60 seconds of successful authentication.
- **SC-002**: During acceptance testing, at least 95% of valid telemetry reports become queryable and eligible for realtime delivery within 1 second of ingestion.
- **SC-003**: For 100% of accepted service requests, the platform returns a tracking identifier immediately and records a final outcome within the configured execution window.
- **SC-004**: All four primary user stories can be demonstrated end to end without manual data correction or out-of-band service intervention.
- **SC-005**: Controlled verification demonstrates support for 1,000 device connections and 10,000 realtime client connections while maintaining the required recovery and availability thresholds.
