# Feature Specification: Core Services Implementation

**Feature Branch**: `004-core-services-implementation`
**Created**: 2025-02-05
**Updated**: 2025-02-05
**Status**: Implemented
**Input**: 实现 5 个核心微服务的业务逻辑
**Depends On**: `001-project-setup` (基础设施), `002-protocol-adapter-design` (协议适配器框架), `003-dji-protocol-implementation` (DJI 协议实现)

## Overview

基于 001/002/003 实现的基础设施和协议适配器，完整实现 5 个核心微服务的业务逻辑。本 Spec 覆盖 iot-gateway、iot-uplink、iot-downlink、iot-api、iot-ws 的核心功能实现。

**数据流架构**:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           完整数据流                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────┐    ┌─────────┐    ┌─────────────┐    ┌─────────────┐          │
│  │ Device  │───►│ VerneMQ │───►│ iot-gateway │───►│  RabbitMQ   │          │
│  └─────────┘    └─────────┘    └─────────────┘    └──────┬──────┘          │
│       ▲              │                                    │                 │
│       │              │                                    ▼                 │
│       │              │              ┌─────────────────────────────────┐    │
│       │              │              │         iot-uplink              │    │
│       │              │              │  - 消息解析验证                   │    │
│       │              │              │  - 物模型映射                     │    │
│       │              │              │  - 时序数据写入                   │    │
│       │              │              │  - 消息路由                       │    │
│       │              │              └─────────────┬───────────────────┘    │
│       │              │                            │                         │
│       │              │                            ▼                         │
│       │              │              ┌─────────────────────────────────┐    │
│       │              │              │         iot-ws                  │    │
│       │              │              │  - WebSocket 连接管理            │    │
│       │              │              │  - 实时消息推送                   │    │
│       │              │              └─────────────────────────────────┘    │
│       │              │                                                      │
│  ┌────┴────┐    ┌────┴────┐    ┌─────────────┐    ┌─────────────┐          │
│  │ Device  │◄───│ VerneMQ │◄───│ iot-gateway │◄───│  RabbitMQ   │          │
│  └─────────┘    └─────────┘    └─────────────┘    └──────┬──────┘          │
│                                                          │                  │
│                                                          │                  │
│              ┌─────────────────────────────────┐         │                  │
│              │         iot-downlink            │◄────────┘                  │
│              │  - 消息格式转换                   │                           │
│              │  - 消息确认重试                   │                           │
│              │  - 路由到 gateway               │                           │
│              └─────────────┬───────────────────┘                           │
│                            │                                                │
│                            │                                                │
│              ┌─────────────┴───────────────────┐                           │
│              │         iot-api                 │                           │
│              │  - RESTful API                  │                           │
│              │  - 设备管理                       │                           │
│              │  - 服务调用                       │                           │
│              └─────────────────────────────────┘                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Architecture

### 服务职责划分

| 服务 | 职责 | 数据访问 | 优先级 |
|------|------|----------|--------|
| **iot-gateway** | MQTT↔RabbitMQ 桥接，设备认证 | PostgreSQL (设备认证) | P1 |
| **iot-uplink** | 上行消息处理，物模型映射 | PostgreSQL (读), InfluxDB (写) | P1 |
| **iot-downlink** | 下行消息路由，确认重试 | PostgreSQL | P1 |
| **iot-api** | RESTful API，业务逻辑 | PostgreSQL, InfluxDB | P2 |
| **iot-ws** | WebSocket 实时推送 | 无 (通过 RabbitMQ) | P2 |

### RabbitMQ 消息路由

```
Exchange: iot.topic (topic exchange)

上行消息路由:
  iot-gateway → iot.raw.{vendor}.uplink → iot-uplink
  iot-uplink  → iot.{vendor}.{service}.{action} → iot-ws/iot-api

下行消息路由:
  iot-api     → iot.{vendor}.service.call → iot-downlink
  iot-downlink → iot.raw.{vendor}.downlink → iot-gateway
```

## User Scenarios & Testing

### User Story 1 - MQTT 消息桥接 (Priority: P1)

作为平台运维人员，我需要 iot-gateway 能够将 MQTT 消息转发到 RabbitMQ，以便其他服务处理。

**Why this priority**: Gateway 是整个数据流的入口，是所有其他功能的基础。

**Acceptance Scenarios**:

1. **Given** DJI 设备发布 MQTT 消息到 VerneMQ, **When** iot-gateway 接收到消息, **Then** 能够转换为标准格式并发布到 RabbitMQ `iot.raw.dji.uplink`
2. **Given** iot-downlink 发布下行消息到 RabbitMQ, **When** iot-gateway 接收到消息, **Then** 能够转换为 MQTT 格式并发布到 VerneMQ
3. **Given** 设备连接到 VerneMQ, **When** 设备认证, **Then** iot-gateway 能够验证设备凭证

---

### User Story 2 - 上行消息处理 (Priority: P1)

作为平台运维人员，我需要 iot-uplink 能够处理上行消息，解析物模型并存储时序数据。

**Why this priority**: 上行消息处理是 IoT 平台的核心功能。

**Acceptance Scenarios**:

1. **Given** iot-gateway 发布上行消息到 RabbitMQ, **When** iot-uplink 接收到消息, **Then** 能够解析消息并验证格式
2. **Given** 上行消息包含遥测数据, **When** iot-uplink 处理消息, **Then** 能够写入 InfluxDB 时序数据库
3. **Given** 上行消息需要推送到客户端, **When** iot-uplink 处理完成, **Then** 能够发布到 RabbitMQ 供 iot-ws 消费

---

### User Story 3 - 下行消息路由 (Priority: P1)

作为平台操作员，我需要 iot-downlink 能够将服务调用路由到设备。

**Why this priority**: 下行消息是平台控制设备的核心能力。

**Acceptance Scenarios**:

1. **Given** iot-api 发布服务调用请求, **When** iot-downlink 接收到消息, **Then** 能够路由到正确的 gateway
2. **Given** 服务调用需要确认, **When** 设备未响应, **Then** iot-downlink 能够重试
3. **Given** 服务调用超时, **When** 超过配置时间, **Then** iot-downlink 能够返回超时错误

---

### User Story 4 - RESTful API (Priority: P2)

作为平台用户，我需要通过 HTTP API 管理设备和调用服务。

**Why this priority**: API 是用户与平台交互的主要方式。

**Acceptance Scenarios**:

1. **Given** 用户请求设备列表, **When** 调用 GET /api/v1/devices, **Then** 返回设备列表
2. **Given** 用户请求调用设备服务, **When** 调用 POST /api/v1/devices/{sn}/services/{method}, **Then** 能够下发服务调用
3. **Given** 用户请求设备遥测数据, **When** 调用 GET /api/v1/devices/{sn}/telemetry, **Then** 返回时序数据

---

### User Story 5 - WebSocket 实时推送 (Priority: P2)

作为平台用户，我需要通过 WebSocket 接收实时消息推送。

**Why this priority**: 实时推送是监控场景的核心需求。

**Acceptance Scenarios**:

1. **Given** 客户端建立 WebSocket 连接, **When** 连接成功, **Then** 能够订阅设备消息
2. **Given** 设备上报遥测数据, **When** iot-uplink 处理完成, **Then** iot-ws 能够推送到订阅的客户端
3. **Given** 客户端断开连接, **When** 连接关闭, **Then** 能够清理订阅关系

---

### Edge Cases

- MQTT 连接断开时如何处理？（自动重连，指数退避 1s/2s/4s/8s 最大 30s，消息缓存最大 1000 条）
- RabbitMQ 连接断开时如何处理？（自动重连，消息持久化，确认机制 publisher confirm）
- 消息处理失败时如何处理？（重试 3 次后进入死信队列 `iot.dlx`，告警通知）
- 大量设备同时上线时如何处理？（连接池最大 1000，限流 100 连接/秒）
- WebSocket 连接数过多时如何处理？（单节点最大 10000 连接，超限返回 503）

## Requirements

### Technology Stack Requirements

- **TS-001**: 必须使用 Go 1.22+ 开发
- **TS-002**: 必须使用 Gin Framework 处理 HTTP/WebSocket
- **TS-003**: 必须使用 GORM 访问 PostgreSQL
- **TS-004**: 必须使用 paho.mqtt.golang 连接 VerneMQ
- **TS-005**: 必须复用 001/002/003 实现的基础设施

### Functional Requirements

#### iot-gateway (P1)

- **FR-001**: 必须实现 MQTT 客户端连接 VerneMQ
- **FR-002**: 必须实现设备认证（用户名/密码）
- **FR-003**: 必须实现 MQTT 消息到 RabbitMQ 的转发
- **FR-004**: 必须实现 RabbitMQ 消息到 MQTT 的转发
- **FR-005**: 必须实现设备连接状态管理

#### iot-uplink (P1)

- **FR-006**: 必须实现 RabbitMQ 消息订阅
- **FR-007**: 必须实现消息解析和验证
- **FR-008**: 必须实现物模型映射。将厂商协议数据映射到平台标准物模型(TSL)结构，通过 `product_key` 查询 `ThingModel` 表获取 TSL JSON 定义。厂商物模型属性定义参考各协议文档（如 DJI: `docs/protocol/dji/en/60.api-reference/*/properties.md`），映射维度包括 Properties（属性上报）、Services（服务调用）、Events（事件通知）
- **FR-009**: 必须实现 InfluxDB 时序数据写入
- **FR-010**: 必须实现消息路由到其他服务

#### iot-downlink (P1)

- **FR-011**: 必须实现 RabbitMQ 消息订阅
- **FR-012**: 必须实现消息格式转换
- **FR-013**: 必须实现消息确认和重试机制
- **FR-014**: 必须实现消息路由到 gateway

#### iot-api (P2)

- **FR-015**: 必须实现设备管理 API (CRUD)
- **FR-016**: 必须实现服务调用 API
- **FR-017**: 必须实现遥测数据查询 API
- **FR-018**: 必须实现 OpenAPI 文档

#### iot-ws (P2)

- **FR-019**: 必须实现 WebSocket 连接管理
- **FR-020**: 必须实现消息订阅机制
- **FR-021**: 必须实现实时消息推送
- **FR-022**: 必须实现连接心跳检测

### Non-Functional Requirements

- **NFR-001**: 消息处理延迟 < 100ms (P95)
- **NFR-002**: 支持 1000+ 设备同时在线
- **NFR-003**: 支持 10000+ WebSocket 连接
- **NFR-004**: 服务可用性 > 99.9%（通过故障注入测试验证：模拟 MQTT/RabbitMQ/PostgreSQL/InfluxDB 断连后服务自动恢复，恢复时间 < 30s）

### Test-First Development Requirements

- **TDD-001**: 所有功能开发必须遵循 TDD 原则
- **TDD-002**: 单元测试覆盖率必须 >= 80%
- **TDD-003**: 必须提供集成测试验证完整消息流程
- **TDD-004**: 必须提供端到端测试验证完整数据流

## Success Criteria

### Measurable Outcomes

- **SC-001**: 5 个核心服务全部实现业务逻辑
- **SC-002**: 完整的上行数据流（设备→VerneMQ→gateway→uplink→ws）
- **SC-003**: 完整的下行数据流（api→downlink→gateway→VerneMQ→设备）
- **SC-004**: 单元测试覆盖率 >= 80%
- **SC-005**: 集成测试通过率 100%

## Clarifications

### Session 2025-02-05

- Q: 本 Spec 与 003 的关系？ → A: 003 实现了 DJI 协议适配器，本 Spec 实现核心服务的业务逻辑，两者配合完成完整数据流
- Q: MQTT 客户端库选择？ → A: 使用 paho.mqtt.golang，这是 Eclipse 官方维护的 Go MQTT 客户端
- Q: WebSocket 库选择？ → A: 使用 gorilla/websocket（注：该库已于 2023 年归档为只读，但 API 稳定且社区广泛使用，无需替换。如需迁移可评估 nhooyr.io/websocket）
- Q: FR-002 认证方式选择？ → A: 004 只实现用户名/密码认证，证书认证延后到 005 或后续迭代
