# Feature Specification: IoT Protocol Adapter Design

**Feature Branch**: `002-protocol-adapter-design`
**Created**: 2025-01-27
**Updated**: 2025-02-05
**Status**: Draft
**Input**: User description: "我们需要讨论下怎么设计对接不同的 iot 协议"

## Overview

设计一个可扩展的协议适配器架构，支持对接不同的 IoT 协议（如 DJI、Tuya、通用 MQTT 等），实现协议消息与平台内部标准消息格式的转换。协议适配器作为独立的微服务实现，从 RabbitMQ 订阅原始消息并发布标准消息。所有协议共享统一的物模型定义，协议适配器负责将协议特定的消息格式映射到标准物模型格式，并在标准消息中保留协议元数据信息。

**本 Spec 的范围**: 完成协议适配器框架和接口定义，实现 DJI 协议（基于 `docs/dji` 文档）的协议层适配作为示例。协议层适配包括：消息格式转换、协议元数据提取，不涉及业务功能（如航线管理、媒体管理、实时流媒体等）。其他协议（如 Tuya、通用 MQTT）的适配器将在后续 Spec 中实现。

## Architecture

### 系统架构图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              用户层                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│   Web 管理后台 ──HTTP──► iot-api (REST API，供 Web 调用)                     │
│   Web 管理后台 ◄──WS───► iot-ws  (WebSocket，向 Web 推送实时消息)             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                              设备层                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│   DJI Drone ─────┐                                                          │
│   DJI Dock ──────┼───MQTT───► VerneMQ (MQTT Broker)                        │
│   Tuya Device ───┤                  ▲                                       │
│   Generic ───────┘                  │ MQTT (订阅/发布)                       │
│                                     │                                       │
│                              iot-gateway                                    │
│                         (MQTT Client, 唯一连接 VerneMQ)                      │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼ RabbitMQ (原始消息)
                         Queue: iot.raw.{vendor}.uplink
                                    │
           ┌────────────────────────┼────────────────────────┐
           ▼                        ▼                        ▼
┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│   dji-adapter    │    │  tuya-adapter    │    │ generic-adapter  │
│  (DJI 协议解析)   │    │  (Tuya 协议解析)  │    │   (直接转发)      │
│  原始消息→标准消息 │    │  原始消息→标准消息 │    │  原始消息→标准消息 │
└────────┬─────────┘    └────────┬─────────┘    └────────┬─────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 ▼ RabbitMQ (标准消息)
                      Queue: iot.{vendor}.{service}.{action}
                                 │
                                 ▼
                            iot-uplink
                      (标准消息 → 业务处理 → 存储)
                                 │
                                 ▼ RabbitMQ
         ┌───────────────────────┴───────────────────────┐
         ▼                                               ▼
      iot-ws                                          iot-api
   (推送给 Web 前端)                               (Web 可查询)
         │                                               │
         └───────────────────────┬───────────────────────┘
                                 ▼
                           Web 管理后台
```

### 下行流程（用户控制设备）

```
Web 管理后台 ──HTTP──► iot-api ──RMQ──► iot-downlink ──RMQ──► dji-adapter
                                                                   │
                                                          (标准消息→DJI格式)
                                                                   │
                                                                   ▼ RMQ
                                                            iot-gateway
                                                                   │
                                                                   ▼ MQTT
                                                              VerneMQ → 设备
```

### 架构定位

协议适配器是独立于宪法定义的5个核心服务（iot-api、iot-ws、iot-uplink、iot-downlink、iot-gateway）的**扩展微服务**。

**核心服务职责**（宪法定义）:

| 服务 | 面向 | 职责 |
|------|------|------|
| **iot-gateway** | 设备 | 唯一连接 VerneMQ 的 MQTT Client，将 MQTT 原始消息转发到 RabbitMQ |
| **iot-api** | Web | REST API，供 Web 管理后台调用 |
| **iot-ws** | Web | WebSocket，向 Web 管理后台推送实时消息 |
| **iot-uplink** | 内部 | 处理上行标准消息，执行业务逻辑（存储、事件触发等） |
| **iot-downlink** | 内部 | 处理下行请求，生成标准消息发送给设备 |

**扩展服务职责**（本 Spec 定义）:

| 服务 | 面向 | 职责 |
|------|------|------|
| **dji-adapter** | 内部 | 订阅 DJI 原始消息，解析 DJI 协议，发布标准消息；订阅下行标准消息，转换为 DJI 协议格式 |
| **tuya-adapter** | 内部 | 订阅 Tuya 原始消息，解析 Tuya 协议，发布标准消息（后续 Spec） |
| **generic-adapter** | 内部 | 订阅通用原始消息，直接转发为标准消息（后续 Spec） |

**关键设计决策**:

1. **每个厂商一个适配器服务**：`dji-adapter` 处理所有 DJI 协议消息（无论来自 MQTT/HTTPS/WebSocket），不是按传输协议拆分
2. **适配器不直接连接设备**：适配器从 RabbitMQ 订阅原始消息，不直接连接 VerneMQ 或接收 HTTP/WebSocket
3. **双向转换**：适配器负责上行（原始→标准）和下行（标准→原始）的双向消息转换

## User Scenarios & Testing

### User Story 1 - 支持 DJI 协议层适配 (Priority: P1)

作为平台开发者，我需要实现 DJI 协议的协议层适配，能够解析和转换 DJI 协议消息，以便验证协议适配器框架的可行性。

**Why this priority**: DJI 协议是平台需要支持的第一个具体 IoT 协议，基于 `docs/dji` 文档实现，作为协议适配器框架的验证和示例。本 Spec 只关注协议层适配，不涉及业务逻辑。

**Independent Test**: 可以通过模拟 iot-gateway 发送 DJI 原始消息到 RabbitMQ，验证 dji-adapter 能够正确解析、转换并发布标准消息。

**Acceptance Scenarios**:

1. **Given** iot-gateway 将 DJI 设备的 MQTT 原始消息发布到 RabbitMQ 队列 `iot.raw.dji.uplink`, **When** dji-adapter 订阅并接收该消息, **Then** dji-adapter 能够解析 DJI 协议格式（osd/state/services/events），转换为标准消息格式，发布到 RabbitMQ 队列 `iot.dji.device.property.report`
2. **Given** iot-downlink 将下行指令发布到 RabbitMQ 队列 `iot.dji.service.call`, **When** dji-adapter 订阅并接收该标准消息, **Then** dji-adapter 能够转换为 DJI MQTT 格式，发布到 RabbitMQ 队列 `iot.raw.dji.downlink`，由 iot-gateway 发送到设备
3. **Given** DJI 设备上报事件消息, **When** dji-adapter 接收到原始事件消息, **Then** dji-adapter 能够解析事件类型和数据，转换为标准事件消息格式

---

### User Story 2 - 协议适配器框架可扩展性 (Priority: P2)

作为平台开发者，我需要协议适配器框架具有良好的可扩展性，以便后续能够轻松添加其他协议适配器（如 tuya-adapter、generic-adapter 等）。

**Why this priority**: 框架的可扩展性决定了未来添加新协议的成本和复杂度。

**Independent Test**: 可以通过实现一个新的协议适配器接口，验证框架是否支持快速开发和部署。

**Acceptance Scenarios**:

1. **Given** 协议适配器框架已实现, **When** 开发者基于框架开发新的协议适配器（如 tuya-adapter）, **Then** 只需实现协议解析和转换逻辑，框架负责 RabbitMQ 订阅/发布、追踪、指标等基础设施
2. **Given** 新协议适配器已开发, **When** 部署为独立服务并配置 RabbitMQ 队列绑定, **Then** 能够自动接收对应厂商的原始消息并处理

---

### Edge Cases

- 原始消息格式不符合预期协议规范时如何处理？（记录错误日志，发送到死信队列）
- 协议适配器处理消息失败时如何重试？（RabbitMQ Nack + 死信队列）
- 多个设备使用相同设备标识但来自不同厂商时如何区分？（通过 vendor 字段区分）
- 协议适配器服务重启时如何保证消息不丢失？（RabbitMQ 持久化 + 手动 Ack）

## Requirements

### Technology Stack Requirements

- **TS-001**: 协议适配器服务必须使用 Go (Golang) 语言开发，版本要求 Go 1.22+
- **TS-002**: 如果协议适配器需要提供 HTTP API（如健康检查、管理接口），必须使用 Gin Framework
- **TS-003**: 如果协议适配器需要访问数据库（如读取设备配置、记录消息状态），必须使用 GORM 框架
- **TS-004**: 所有代码必须严格遵循 [Uber Go 语言编码规范](https://github.com/uber-go/guide)，包括命名规范、错误处理、并发模式等
- **TS-005**: 所有代码、注释、文档、变量名、函数名、文件名等严禁出现拼写错误，必须使用 `golangci-lint` 和 `misspell` 进行代码检查

### Functional Requirements

- **FR-001**: 系统必须实现协议适配器框架，包括接口定义、消息转换机制、RabbitMQ 集成等基础设施
- **FR-002**: 系统必须实现 `dji-adapter` 服务，支持 DJI Cloud API 协议的解析和转换（基于 `docs/dji` 文档）
- **FR-003**: 系统必须支持 DJI 协议消息格式的解析和转换，包括 Topic 解析、消息结构转换、协议元数据提取
- **FR-004**: 系统必须支持 DJI MQTT 核心 Topic：
  - `thing/product/{device_sn}/osd`（属性定时上报）
  - `thing/product/{device_sn}/state`（属性变化上报）
  - `thing/product/{gateway_sn}/services`（服务调用）
  - `thing/product/{gateway_sn}/events`（事件上报）
  - `sys/product/{gateway_sn}/status`（设备状态）
- **FR-005**: [DEFERRED] 其他 DJI Topic（如 DRC、Wayline、Media Management 等）将在后续 Spec 中实现
- **FR-006**: [OUT OF SCOPE] 业务功能（如航线管理、媒体管理、实时流媒体、设备管理业务逻辑等）不在本 Spec 范围内
- **FR-007**: 系统必须提供协议适配器接口定义，包括：
  - `ParseRawMessage(raw []byte) (*ProtocolMessage, error)` - 解析原始消息
  - `ToStandardMessage(pm *ProtocolMessage) (*StandardMessage, error)` - 转换为标准消息
  - `FromStandardMessage(sm *StandardMessage) ([]byte, error)` - 从标准消息转换为原始格式
- **FR-008**: 系统必须将不同协议的消息转换为统一的标准消息格式（复用 001 定义的 `StandardMessage`）
- **FR-009**: RabbitMQ 队列命名规范：
  - 原始消息（上行）: `iot.raw.{vendor}.uplink`（如 `iot.raw.dji.uplink`）
  - 原始消息（下行）: `iot.raw.{vendor}.downlink`（如 `iot.raw.dji.downlink`）
  - 标准消息: `iot.{vendor}.{service}.{action}`（复用 001 定义的格式）
- **FR-010**: 系统必须提供完整的可观测性支持（复用 001 的 pkg/tracer、pkg/metrics、internal/shared/logger）：
  - 结构化日志记录：包含 trace_id、device_id、vendor、message_type 等
  - 分布式追踪：W3C Trace Context 传播
  - 指标监控：消息处理数量、处理延迟、错误率
- **FR-011**: 协议适配器必须作为独立的微服务实现（`cmd/dji-adapter/main.go`），可独立部署和扩展
- **FR-012**: 所有协议必须共享统一的物模型定义，协议适配器必须负责将协议特定的消息格式映射到标准物模型格式
- **FR-013**: 标准消息必须包含协议元数据字段（扩展 001 的 `StandardMessage`），包括：
  - `protocol_meta.vendor` - 厂商标识（dji、tuya、generic）
  - `protocol_meta.original_topic` - 原始 MQTT Topic
  - `protocol_meta.qos` - MQTT QoS 级别
- **FR-014**: 协议适配器必须支持环境变量 + 配置文件（YAML）的配置管理方式（复用 001 的配置框架）

### Test-First Development Requirements

- **TDD-001**: 所有功能开发必须遵循测试驱动开发（TDD）原则
- **TDD-002**: 协议适配器必须提供模拟消息用于测试，支持模拟 iot-gateway 发送 DJI 原始消息
- **TDD-003**: 所有协议适配器接口必须提供契约测试（Contract Test）
- **TDD-004**: 单元测试覆盖率必须 ≥ 80%

### Key Entities

- **ProtocolAdapter**: 协议适配器接口，定义协议消息的解析、转换方法
- **ProtocolMessage**: 协议消息，包含协议特定的消息格式和元数据
- **StandardMessage**: 标准消息，平台内部统一的消息格式（复用 001 定义，需扩展 protocol_meta 字段）
- **RawMessage**: 原始消息，iot-gateway 从 VerneMQ 接收的未处理消息
- **MessageConverter**: 消息转换器，负责协议消息与标准消息的双向转换
- **ThingModel**: 物模型定义，所有协议共享的统一物模型（复用 001 定义）

## Success Criteria

### Measurable Outcomes

- **SC-001**: dji-adapter 能够正确解析 DJI 核心 Topic（osd、state、services、events、status）的消息
- **SC-002**: 协议消息转换延迟 < 50ms (P95)
- **SC-003**: 协议消息处理成功率 > 99.9%
- **SC-004**: 框架接口清晰，后续添加新协议适配器（如 tuya-adapter）的开发时间 < 2 天
- **SC-005**: 单元测试覆盖率 ≥ 80%

## Clarifications

### Session 2025-01-27

- Q: 协议适配器应该作为 iot-gateway 服务的一部分，还是作为独立的可插拔组件？ → A: 作为独立的微服务，从 RabbitMQ 订阅原始消息
- Q: 协议适配器如何与物模型关联？ → A: 所有协议共享统一的物模型定义，协议适配器负责将协议消息映射到物模型格式
- Q: 协议适配器如何处理协议特定的特性（如 MQTT QoS）？ → A: 在标准消息中保留 protocol_meta 字段
- Q: 协议适配器是否需要支持协议转换（如 MQTT 转 HTTPS）？ → A: 不支持协议间转换，只负责协议消息与标准消息的双向转换
- Q: 协议适配器的配置管理方式是什么？ → A: 复用 001 的配置框架（环境变量 + YAML）
- Q: spec2 的完成范围是什么？ → A: 完成协议适配器框架和 dji-adapter 服务，只实现协议层适配，不涉及业务功能

### Session 2025-02-05 (架构澄清)

- Q: dji-mqtt-adapter、dji-https-adapter、dji-ws-adapter 三个服务是否违反设计？ → A: 是的，原设计有误。正确设计是**每个厂商一个适配器服务**（如 dji-adapter），而不是按传输协议拆分
- Q: iot-api 和 iot-ws 的定位？ → A: iot-api 是给 Web 管理后台调用的 REST API；iot-ws 是给 Web 管理后台的 WebSocket 长连接，用于实时推送。它们不直接面向设备
- Q: 设备如何连接平台？ → A: 设备直连 VerneMQ（MQTT Broker），iot-gateway 作为 MQTT Client 订阅设备消息并转发到 RabbitMQ
- Q: 协议适配器从哪里获取消息？ → A: 从 RabbitMQ 订阅原始消息（`iot.raw.{vendor}.uplink`），不直接连接设备或 VerneMQ

