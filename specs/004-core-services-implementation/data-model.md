# Data Model: Core Services Implementation

**Feature**: 004-core-services-implementation
**Date**: 2026-03-23

## Entity Relationship Overview

```text
ThingModel ──────< ManagedDevice >────── DeviceCredential
                     │      │
                     │      └──────< TelemetryRecord
                     │
                     ├──────< DeviceTopology
                     │
                     └──────< ServiceRequest

RealtimeSession ──────< RealtimeSubscription
```

## Entities

### ThingModel

Represents the TSL-based capability definition used to validate and interpret device traffic.

**Fields**:
- `id`: unique identifier
- `product_key`: unique product/capability key
- `product_name`: human-readable product name
- `version`: capability model version
- `tsl_json`: full TSL JSON definition
- `description`: optional description
- `created_at`, `updated_at`, `deleted_at`

**Validation rules**:
- `product_key` is required and unique
- `version` is required
- `tsl_json` is required and must contain properties, services, and events required by the platform constitution

**Relationships**:
- One thing model can be associated with many managed devices

### ManagedDevice

Represents a device known to the platform and tracked across authentication, connectivity, telemetry, and command execution.

**Fields**:
- `id`: unique identifier
- `device_sn`: unique device serial number
- `device_name`: display name
- `device_type`: gateway, aircraft, dock, rc, or other supported platform type
- `vendor`: supported vendor identifier
- `status`: current connectivity/operational status
- `gateway_sn`: optional parent gateway serial number
- `thing_model_id`: optional reference to the capability model
- `last_online_time`: latest observed online time
- `created_at`, `updated_at`, `deleted_at`

**Validation rules**:
- `device_sn` is required and unique
- `device_name` is required for operator-facing inventory management
- `device_type` is required
- `vendor` is required
- `status` defaults to `unknown` until connectivity is observed

**Relationships**:
- Many devices belong to one thing model
- One device can have one active credential record
- One device can have many telemetry records
- One device can have many service requests
- One device can participate in zero or more topology links

### DeviceCredential

Represents the authentication material used to decide whether a device may connect.

**Fields**:
- `id`: unique identifier
- `device_sn`: device serial number
- `username`: credential username
- `password_hash`: stored secret representation
- `enabled`: whether the credential is active
- `created_at`, `updated_at`

**Validation rules**:
- `device_sn` is required and unique within the credential store
- `username` is required
- `password_hash` is required
- disabled credentials must be rejected for new device sessions

**Relationships**:
- One credential record belongs to one managed device

### DeviceTopology

Represents parent-child device relationships for gateway-managed deployments.

**Fields**:
- `id`: unique identifier
- `device_id`: child device reference
- `parent_device_id`: optional parent device reference
- `relation_type`: relationship classification
- `created_at`

**Validation rules**:
- `device_id` is required
- parent and child must not reference the same device
- relation type must be one of the supported topology patterns

**Relationships**:
- A managed device can have zero or one parent link in the active topology
- A managed device can have many child links when acting as a gateway

### TelemetryRecord

Represents a stored upstream measurement, state update, or event for history and downstream use.

**Fields**:
- `measurement`: logical telemetry/event stream name
- `device_sn`: source device serial number
- `gateway_sn`: optional gateway serial number
- `vendor`: source vendor identifier
- `message_type`: property, event, status, or other supported category
- `payload`: normalized business data
- `timestamp`: event time from the source message
- `trace_id`: distributed trace identifier when present

**Validation rules**:
- `device_sn` is required
- `message_type` is required
- `timestamp` is required
- payload must satisfy the registered capability model for supported measurements

**Relationships**:
- Many telemetry records belong to one managed device

### ServiceRequest

Represents a tracked command or property/configuration change request sent toward a device.

**Fields**:
- `id`: unique request identifier
- `device_sn`: target device serial number
- `vendor`: target vendor identifier
- `method`: requested operation
- `params`: request payload
- `call_type`: command, property, or config
- `status`: pending, sent, success, failed, timeout, or retrying
- `tid`: transaction identifier
- `bid`: business identifier
- `retry_count`: current retry count
- `max_retries`: retry limit
- `error`: final or latest error text
- `response`: response payload when available
- `sent_at`, `completed_at`, `created_at`, `updated_at`

**Validation rules**:
- `device_sn`, `vendor`, and `method` are required
- `status` defaults to `pending`
- `max_retries` defaults to 3 when not provided
- `tid` and `bid` should be preserved across retries for traceability

**Relationships**:
- Many service requests belong to one managed device

### RealtimeSession

Represents an active client session receiving live platform updates.

**Fields**:
- `session_id`: unique client/session identifier
- `user_id`: authenticated platform user or client principal
- `device_scope`: optional device-scoping metadata from the request context
- `connected_at`: session start time
- `last_heartbeat_at`: latest successful heartbeat time
- `state`: active or closed

**Validation rules**:
- `session_id` is required
- sessions without valid heartbeat progression must be closed

**Relationships**:
- One realtime session can have many subscriptions

### RealtimeSubscription

Represents an authorized topic subscription attached to a realtime session.

**Fields**:
- `session_id`: owning session identifier
- `topic_pattern`: subscribed topic or pattern
- `created_at`: subscription creation time

**Validation rules**:
- `topic_pattern` is required
- subscriptions must only target supported device topics
- subscriptions must pass authorization checks for the owning client

**Relationships**:
- Many subscriptions belong to one realtime session

## State Transitions

### ManagedDevice connectivity status

```text
unknown -> online
online  -> offline
offline -> online
```

**Rules**:
- A successful authenticated connection can move a device to `online`
- Missing expected traffic or explicit disconnect can move a device to `offline`
- Devices begin in `unknown` until the platform observes activity

### ServiceRequest lifecycle

```text
pending -> sent -> success
              \-> failed -> retrying -> sent
              \-> timeout -> retrying -> sent
```

**Rules**:
- Requests start as `pending`
- A request moves to `sent` once accepted for device delivery
- Final terminal states are `success`, `failed`, and `timeout`
- `retrying` is transitional and increments `retry_count`
- Requests cannot retry after reaching `max_retries`

### RealtimeSession lifecycle

```text
active -> closed
```

**Rules**:
- A session becomes active after a successful WebSocket upgrade
- A missed heartbeat, client disconnect, or explicit shutdown closes the session and removes all subscriptions

## Message and trace envelope requirements

All normalized inter-service messages must preserve:
- `tid`: transaction identifier
- `bid`: business identifier
- `timestamp`: source or processing time
- `device_sn`: source or target device serial number
- `service`: logical producing service name
- `action`: logical business action
- `protocol_meta.vendor`: vendor identifier when available
- W3C trace context headers (`traceparent`, `tracestate`) when propagated through RabbitMQ

## Design notes

- Shared repository models already exist for `ThingModel` and `Device`; planning should extend rather than duplicate them.
- `ServiceRequest` persistence is necessary for operator-facing auditability and retry tracking.
- `TelemetryRecord` is represented operationally in time-series storage rather than a single relational table, but its logical fields must remain stable for API and realtime consumers.
- `RealtimeSession` and `RealtimeSubscription` may remain memory-backed at runtime if that satisfies current scope, but their behavior is part of the feature contract and must be tested.
