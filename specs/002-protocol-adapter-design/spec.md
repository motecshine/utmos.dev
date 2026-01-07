# Feature Specification: IoT Protocol Adapter Design

**Feature Branch**: `002-protocol-adapter-design`  
**Created**: 2025-01-27  
**Status**: Draft  
**Input**: User description: "我们需要讨论下怎么设计对接不同的 iot 协议"

## Overview

设计一个可扩展的协议适配器架构，支持对接不同的 IoT 协议，实现协议消息与平台内部标准消息格式的转换。协议适配器作为独立的微服务实现，与 iot-gateway 通过 RabbitMQ 进行通信。所有协议共享统一的物模型定义，协议适配器负责将协议特定的消息格式映射到标准物模型格式，并在标准消息中保留协议元数据信息。

**架构定位**: 协议适配器是独立于宪法定义的5个核心服务（iot-api、iot-ws、iot-uplink、iot-downlink、iot-gateway）的扩展微服务。每个协议适配器服务实例负责特定协议的适配（如 `dji-mqtt-adapter`、`dji-https-adapter`、`dji-ws-adapter`）。协议适配器与 iot-gateway 通过 RabbitMQ 进行消息通信，不直接连接 MQTT Broker（VerneMQ）。这种设计实现了协议适配逻辑与核心网关服务的解耦，提高了系统的可扩展性和可维护性。

**本 Spec 的范围**: 完成协议适配器框架和接口定义，实现 DJI 协议（基于 `docs/dji` 文档）的协议层适配作为示例。协议层适配包括：消息格式转换、连接管理、协议元数据提取，不涉及业务功能（如航线管理、媒体管理、实时流媒体等）。其他协议（如通用 MQTT、HTTPS、WebSocket）的适配器将在后续 Spec 中实现。

## User Scenarios & Testing

### User Story 1 - 支持 DJI 协议层适配 (Priority: P1)

作为平台开发者，我需要实现 DJI 协议的协议层适配，能够接收和转换 DJI 协议消息，以便验证协议适配器框架的可行性。

**Why this priority**: DJI 协议是平台需要支持的第一个具体 IoT 协议，基于 `docs/dji` 文档实现，作为协议适配器框架的验证和示例。本 Spec 只关注协议层适配，不涉及业务逻辑。

**Independent Test**: 可以通过模拟 DJI 设备（如 DJI Dock、DJI Pilot 2）发送协议消息，验证协议适配器能够正确接收、解析、转换并路由到 RabbitMQ。

**Acceptance Scenarios**:

1. **Given** DJI 设备通过 MQTT 连接到平台, **When** 设备上报属性数据（osd/state Topic）, **Then** 协议适配器能够接收并转换为标准消息格式，发布到 RabbitMQ 队列 `iot.adapter.dji.mqtt.uplink`
2. **Given** 平台需要下发服务调用给 DJI 设备, **When** 通过 RabbitMQ 队列 `iot.adapter.dji.mqtt.downlink` 接收标准消息, **Then** 协议适配器能够转换为 DJI MQTT 格式并发送到设备
3. **Given** DJI 设备通过 HTTPS 调用 API, **When** 设备发送 HTTP 请求, **Then** 协议适配器能够接收并转换为标准消息格式，发布到 RabbitMQ 队列 `iot.adapter.dji.https.uplink`
4. **Given** DJI 设备通过 WebSocket 连接, **When** 平台通过 RabbitMQ 队列 `iot.adapter.dji.ws.downlink` 推送消息, **Then** 协议适配器能够转换为 WebSocket 格式并推送给设备

---

### User Story 2 - 协议适配器框架可扩展性 (Priority: P2)

作为平台开发者，我需要协议适配器框架具有良好的可扩展性，以便后续能够轻松添加其他协议（如通用 MQTT、CoAP、LoRaWAN 等）。

**Why this priority**: 框架的可扩展性决定了未来添加新协议的成本和复杂度。

**Independent Test**: 可以通过实现一个新的协议适配器接口，验证框架是否支持启动时加载和注册。

**Acceptance Scenarios**:

1. **Given** 协议适配器框架已实现, **When** 开发者实现新的协议适配器接口并在服务启动时注册, **Then** 能够成功注册并加载
2. **Given** 新协议适配器已注册, **When** 设备使用该协议连接, **Then** 框架能够正确路由到对应的适配器

---

### User Story 3 - 通过插件架构添加新协议支持 (Priority: P3)

作为平台开发者，我需要协议适配器框架采用插件化架构，支持通过插件方式添加新的 IoT 协议支持，以便扩展平台能力而无需修改核心代码。

**Why this priority**: 插件化架构提高了系统的可扩展性，符合宪法要求。插件化架构使得新协议适配器可以作为独立的插件模块开发、测试和部署。

**Independent Test**: 可以通过实现协议适配器插件接口，验证新协议插件能够在服务启动时成功注册和使用。

**Acceptance Scenarios**:

1. **Given** 新协议适配器插件已开发, **When** 插件在服务启动时注册到平台, **Then** 平台能够识别并使用该协议
2. **Given** 协议适配器插件已加载, **When** 设备使用该协议连接, **Then** 平台能够正确处理消息

---

### Edge Cases

- 协议连接断开时如何处理未完成的消息？
- 协议消息格式不符合规范时如何处理？
- 多个设备使用不同协议但相同设备标识时如何区分？
- 协议适配器插件加载失败时如何降级处理？

## Requirements

### Technology Stack Requirements

- **TS-001**: 协议适配器服务必须使用 Go (Golang) 语言开发，版本要求 Go 1.21+
- **TS-002**: 如果协议适配器需要提供 HTTP API（如健康检查、管理接口），必须使用 Gin Framework
- **TS-003**: 如果协议适配器需要访问数据库（如读取设备配置、记录消息状态），必须使用 GORM 框架
- **TS-004**: 所有代码必须严格遵循 [Uber Go 语言编码规范](https://github.com/uber-go/guide)，包括命名规范、错误处理、并发模式等
- **TS-005**: 所有代码、注释、文档、变量名、函数名、文件名等严禁出现拼写错误，必须使用 `golangci-lint` 和 `misspell` 进行代码检查

### Functional Requirements

- **FR-001**: 系统必须实现协议适配器框架，包括接口定义、消息转换机制、RabbitMQ 集成等基础设施
- **FR-002**: 系统必须实现 DJI 协议适配器，支持 DJI Cloud API 的 MQTT、HTTPS、WebSocket 三种协议的协议层适配（基于 `docs/dji` 文档）
- **FR-003**: 系统必须支持 DJI 协议消息格式的解析和转换，包括 Topic 解析、消息结构转换、协议元数据提取
- **FR-016**: 系统必须支持 DJI MQTT 核心 Topic：`thing/product/{device_sn}/osd`（属性定时上报）、`thing/product/{device_sn}/state`（属性变化上报）、`thing/product/{gateway_sn}/services`（服务调用）、`thing/product/{gateway_sn}/events`（事件上报）、`sys/product/{gateway_sn}/status`（设备状态）
- **FR-017**: 系统必须支持 DJI HTTPS 核心 API：设备认证、基础查询接口
- **FR-018**: 系统必须支持 DJI WebSocket 核心功能：连接建立、消息推送
- **FR-019**: [DEFERRED] 其他 DJI Topic 和 API（如 DRC、Wayline、Media Management 等）将在后续 Spec 中实现
- **FR-014**: 系统必须支持 DJI 协议的物模型映射框架，能够将 DJI 设备属性、服务、事件映射到标准物模型格式（映射规则可配置）
- **FR-015**: [OUT OF SCOPE] 业务功能（如航线管理、媒体管理、实时流媒体、设备管理业务逻辑等）不在本 Spec 范围内，将在后续 Spec 中实现
- **FR-004**: 系统必须提供协议适配器接口定义，支持在服务启动时加载协议适配器（不支持运行时动态加载和热更新）
- **FR-020**: [DEFERRED] 运行时动态加载和热更新功能将在后续版本中实现
- **FR-005**: 系统必须将不同协议的消息转换为统一的标准消息格式
- **FR-006**: 系统必须支持协议消息的路由和分发到正确的处理服务，RabbitMQ 队列命名必须遵循统一规范：`iot.adapter.{protocol}.{transport}.{direction}`（如 `iot.adapter.dji.mqtt.uplink`、`iot.adapter.dji.https.downlink`）
- **FR-007**: 系统必须支持设备认证和授权（协议层面）
- **FR-008**: 系统必须提供完整的可观测性支持，包括：
  - 结构化日志记录：所有协议消息处理必须记录结构化日志，包含请求 ID（Request ID）、设备 ID、时间戳、协议类型、消息类型等关键信息
  - 分布式追踪：所有消息处理必须支持分布式追踪，在消息中传递追踪上下文（Trace ID、Span ID），与 Tempo 集成
  - 指标监控：必须暴露 Prometheus 指标，包括消息处理数量、处理延迟、错误率、连接数等
  - 告警机制：关键错误和异常情况必须能够及时告警
- **FR-009**: 协议适配器必须作为独立的微服务实现，每个协议适配器服务实例负责特定协议的适配（如 `dji-mqtt-adapter`、`dji-https-adapter`、`dji-ws-adapter`）。协议适配器与 iot-gateway 通过 RabbitMQ 进行通信，不直接连接 MQTT Broker（VerneMQ）。每个协议适配器服务必须可以独立部署和扩展
- **FR-010**: 所有协议必须共享统一的物模型定义，协议适配器必须负责将协议特定的消息格式映射到标准物模型格式
- **FR-011**: 标准消息必须包含协议元数据字段，协议适配器必须提取协议特定特性（如 MQTT QoS、WebSocket 子协议、HTTPS 请求方法等）并转换为元数据
- **FR-012**: 协议适配器不支持协议间转换，只负责协议消息与标准消息格式的双向转换（协议 → 标准消息，标准消息 → 协议）
- **FR-013**: 协议适配器必须支持环境变量 + 配置文件（YAML）的配置管理方式，敏感信息（如密码、密钥）使用环境变量，非敏感配置使用 YAML 文件
- **FR-021**: 如果协议适配器提供 HTTP API 接口（如健康检查、管理接口），必须遵循标准化 API 设计规范：统一的 URI 格式、统一的请求/响应格式（JSON）、统一的错误码规范、URL 路径版本控制（如 `/api/v1/`）、完整的 OpenAPI/Swagger 文档

### Test-First Development Requirements

- **TDD-001**: 所有功能开发必须遵循测试驱动开发（TDD）原则：先编写测试用例（包括单元测试、集成测试），确保测试失败，然后实现功能，最后重构
- **TDD-002**: 协议适配器必须提供模拟设备（Mock Device）用于测试，支持模拟 DJI 设备发送协议消息
- **TDD-003**: 所有协议适配器接口必须提供契约测试（Contract Test），确保接口规范的一致性
- **TDD-004**: 单元测试覆盖率必须 ≥ 80%

### Key Entities

- **ProtocolAdapterService**: 协议适配器微服务，独立的服务实例，负责特定协议的适配（如 mqtt-adapter、https-adapter、websocket-adapter）
- **ProtocolAdapter**: 协议适配器接口，定义协议消息的接收、发送、转换方法
- **ProtocolMessage**: 协议消息，包含协议特定的消息格式和元数据（如 MQTT QoS、WebSocket 子协议等）
- **StandardMessage**: 标准消息，平台内部统一的消息格式，包含协议元数据字段
- **ThingModel**: 物模型定义，所有协议共享的统一物模型，不因协议而改变
- **MessageConverter**: 消息转换器，负责协议消息与标准消息的双向转换（协议 → 标准消息，标准消息 → 协议）
- **ProtocolRegistry**: 协议注册表，管理已注册的协议适配器服务

## Success Criteria

### Measurable Outcomes

- **SC-001**: 协议适配器框架能够支持 DJI 协议的完整接入（MQTT、HTTPS、WebSocket）
- **SC-005**: 框架接口清晰，后续添加新协议适配器的开发时间 < 2 天
- **SC-002**: 协议消息转换延迟 < 50ms (P95)
- **SC-003**: 协议适配器插件在服务启动时加载时间 < 5 秒
- **SC-004**: 协议消息处理成功率 > 99.9%

## Clarifications

### Session 2025-01-27

- Q: 协议适配器应该作为 iot-gateway 服务的一部分，还是作为独立的可插拔组件？ → A: 作为独立的微服务，与 iot-gateway 通过 RabbitMQ 通信
- Q: 协议适配器如何与物模型关联？是否需要为每个协议定义独立的物模型映射？ → A: 所有协议共享统一的物模型定义，协议适配器负责将协议消息映射到物模型格式
- Q: 协议适配器如何处理协议特定的特性（如 MQTT QoS、WebSocket 子协议等）？ → A: 在标准消息中保留协议元数据字段，协议适配器负责提取和转换
- Q: 协议适配器是否需要支持协议转换（如 MQTT 转 HTTPS）？ → A: 不支持协议间转换，只负责协议消息与标准消息的双向转换
- Q: 协议适配器的配置管理方式是什么？ → A: 环境变量 + 配置文件（YAML），敏感信息用环境变量
- Q: spec2 的完成范围是什么？"完成框架"具体指什么？ → A: 完成协议适配器框架和接口定义，实现 DJI 协议（基于 docs/dji 文档）的协议层适配（消息转换、连接管理），不涉及业务功能
- Q: DJI 协议适配器需要实现哪些具体功能？是否需要实现 docs/dji 文档中的所有功能？ → A: 只实现协议层适配（消息转换、连接管理），不涉及业务功能（如航线管理、媒体管理等）
- Q: DJI 协议适配器需要支持 docs/dji 文档中的哪些具体 Topic 和 API？ → A: 支持核心 Topic 和 API（属性上报 osd/state、服务调用 services、事件上报 events），其他 Topic 后续扩展
- Q: 协议适配器框架是否需要支持动态加载和热更新？ → A: 启动时加载协议适配器，不支持运行时动态加载和热更新
- Q: 协议适配器服务的部署方式是什么？每个协议一个独立的服务实例（如 dji-mqtt-adapter、dji-https-adapter、dji-ws-adapter），还是多个协议共享一个服务实例（如 dji-adapter 支持 MQTT/HTTPS/WebSocket）？ → A: 每个协议一个独立的服务实例（如 dji-mqtt-adapter、dji-https-adapter、dji-ws-adapter），每个服务实例可以独立部署和扩展

