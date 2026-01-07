# Feature Specification: Project Setup with Distributed Tracing and Multi-Vendor RabbitMQ Routing

**Feature Branch**: `001-project-setup`  
**Created**: 2025-01-27  
**Status**: Draft  
**Input**: User description: "ok 我们详细讨论001 spec ， 我们在 001 需要 完成 trace， 多厂商 rabbitmq routingkey定义， 来串联起所有的service"

## Overview

设置 UMOS IoT 平台的基础项目结构，包括5个微服务的项目骨架、开发环境配置、CI/CD 流程、代码质量检查工具等。**本 Spec 的核心目标**：实现分布式追踪（Trace）基础设施，定义多厂商 RabbitMQ routing key 规范，建立服务间消息路由机制，确保所有服务能够通过标准化的消息路由和追踪机制串联起来。

## User Scenarios & Testing

### User Story 1 - 分布式追踪基础设施 (Priority: P1)

作为平台开发者，我需要实现分布式追踪基础设施，能够在所有微服务之间追踪消息流转，以便快速定位问题和分析系统性能。

**Why this priority**: 分布式追踪是微服务架构中可观测性的核心组件，是串联所有服务的关键基础设施。没有追踪能力，无法有效调试和监控跨服务调用。

**Independent Test**: 可以通过发送一条设备消息，验证 trace_id 能够在所有服务（iot-gateway → iot-uplink → iot-api/iot-ws）之间传递，并在 Tempo 中查询到完整的调用链路。

**Acceptance Scenarios**:

1. **Given** 设备通过 MQTT 发送属性上报消息, **When** 消息经过 iot-gateway → iot-uplink → iot-api, **Then** 所有服务日志中包含相同的 trace_id，且能在 Tempo 中查询到完整调用链路
2. **Given** 客户端通过 HTTP API 调用服务, **When** 请求经过 iot-api → iot-downlink → iot-gateway, **Then** 所有服务日志中包含相同的 trace_id，且能在 Tempo 中查询到完整调用链路
3. **Given** 消息在 RabbitMQ 中流转, **When** 消息被多个服务消费, **Then** trace_id 在消息头中传递，每个服务都能正确提取和记录

---

### User Story 2 - 多厂商 RabbitMQ Routing Key 定义 (Priority: P1)

作为平台开发者，我需要定义多厂商 RabbitMQ routing key 规范，支持不同厂商（如 DJI、其他 IoT 厂商）的消息路由，以便实现统一的消息路由机制。

**Why this priority**: 多厂商支持是平台的核心能力，routing key 定义是消息路由的基础，必须在项目初始化时建立。

**Independent Test**: 可以通过模拟不同厂商（如 DJI、通用 MQTT）的消息，验证 routing key 能够正确路由到对应的处理服务。

**Acceptance Scenarios**:

1. **Given** DJI 设备发送属性上报消息, **When** 消息通过 RabbitMQ 路由, **Then** routing key 能够标识厂商（dji）和消息类型（property.report），正确路由到处理服务
2. **Given** 不同厂商的设备发送相同类型的消息, **When** 消息通过 RabbitMQ 路由, **Then** routing key 能够区分厂商，路由到对应的处理逻辑
3. **Given** 消息需要跨服务流转, **When** 使用标准化的 routing key, **Then** 所有服务能够正确识别和路由消息

---

### User Story 3 - 统一 Metrics 包 (Priority: P1)

作为平台开发者，我需要统一的 metrics 包来处理框架基础中间件（RabbitMQ、PostgreSQL、InfluxDB 等）的 metrics，同时也应该在业务代码中使用，以便统一监控和告警。

**Why this priority**: 统一的 metrics 包是可观测性的核心组件，能够统一管理框架中间件和业务 metrics，便于监控和告警配置。

**Independent Test**: 可以通过查看 Prometheus 指标端点，验证框架中间件（RabbitMQ 连接数、PostgreSQL 查询延迟等）和业务 metrics（消息处理数量、错误率等）都能正确暴露。

**Acceptance Scenarios**:

1. **Given** 服务启动, **When** 访问 Prometheus metrics 端点, **Then** 能够看到框架中间件的 metrics（RabbitMQ、PostgreSQL、InfluxDB 等）
2. **Given** 业务代码使用 metrics 包, **When** 记录业务 metrics, **Then** 能够统一暴露到 Prometheus，与框架 metrics 使用相同的格式和标签规范
3. **Given** 多个服务使用 metrics 包, **When** 查看 Grafana 仪表板, **Then** 能够统一展示所有服务的框架和业务 metrics

---

### User Story 4 - 服务间消息串联机制 (Priority: P1)

作为平台开发者，我需要建立服务间消息串联机制，确保所有服务能够通过 RabbitMQ 和分布式追踪串联起来，形成完整的消息流转链路。

**Why this priority**: 服务间消息串联是微服务架构的基础，必须确保消息能够在所有服务之间正确流转，且能够被追踪。

**Independent Test**: 可以通过端到端测试，验证一条消息从设备到客户端响应的完整流程，所有服务都能正确处理和传递消息。

**Acceptance Scenarios**:

1. **Given** 设备发送消息, **When** 消息经过所有服务处理, **Then** 消息能够正确流转，trace_id 贯穿整个链路
2. **Given** 服务间消息传递, **When** 使用标准化的消息格式和 routing key, **Then** 所有服务能够正确解析和路由消息
3. **Given** 消息处理失败, **When** 发生错误, **Then** 错误信息能够通过 trace_id 关联到完整的调用链路

---

### Edge Cases

- 消息丢失时如何通过 trace_id 追踪？
- 多个厂商使用相同的设备标识时如何区分？（通过数据库中的 sn + vendor 映射区分）
- trace_id 在消息重试时如何保持一致性？
- routing key 不匹配时如何处理？

## Requirements

### Functional Requirements

- **FR-001**: 系统必须实现分布式追踪基础设施，集成 OpenTelemetry SDK 和 Tempo，使用 W3C Trace Context 标准在 HTTP 请求头和 RabbitMQ 消息头中传递 trace context
- **FR-002**: 系统必须定义多厂商 RabbitMQ routing key 规范，格式为 `iot.{vendor}.{service}.{action}`，支持厂商标识（vendor）和消息类型（message_type）的路由
- **FR-003**: 系统必须实现服务间消息串联机制，使用 RabbitMQ Topic Exchange（`iot`）和动态 Queue 绑定，确保消息能够在所有服务之间正确流转
- **FR-004**: 所有 RabbitMQ 消息必须在消息头（headers）中包含 W3C Trace Context（`traceparent`、`tracestate`），支持分布式追踪
- **FR-005**: 所有服务必须能够从 HTTP 请求头或 RabbitMQ 消息头中提取 W3C Trace Context，创建新的 span，并在日志和指标中记录 trace_id 和 span_id
- **FR-006**: routing key 必须遵循统一规范 `iot.{vendor}.{service}.{action}`，支持多厂商扩展
- **FR-007**: 消息格式必须包含设备序列号（device_sn），设备上云流程中将 sn + vendor 写入数据库，通过 sn 查询数据库获取 vendor 用于 routing key 生成
- **FR-008**: 系统必须实现统一的 metrics 包，基于 Prometheus Go Client Library（`github.com/prometheus/client_golang`），创建统一封装包（如 `pkg/metrics`），提供中间件集成（RabbitMQ、PostgreSQL、InfluxDB）自动收集连接数、查询延迟等 metrics，同时支持手动注册自定义 metrics
- **FR-009**: metrics 包必须支持业务代码使用，提供统一的 API 用于记录业务 metrics（计数器 Counter、直方图 Histogram、仪表盘 Gauge 等），与框架 metrics 使用相同的标签规范
- **FR-010**: 所有 metrics 必须暴露 Prometheus 格式，使用统一的标签规范（必需标签：service、vendor、message_type、status 等），同时支持业务代码添加自定义标签，支持 Grafana 可视化
- **FR-011**: 所有服务必须在统一 HTTP 端点 `/metrics` 暴露 Prometheus 格式的 metrics，使用 Gin 中间件或独立 HTTP 服务器实现
- **FR-012**: Metrics 命名必须遵循统一规范 `iot_{component}_{metric_type}_{unit}`，例如：`iot_rabbitmq_connection_total`、`iot_postgres_query_duration_seconds`、`iot_message_processed_total`
- **FR-013**: 系统必须实现统一的配置管理，使用 YAML 配置文件，支持多环境配置（dev、staging、prod），所有服务必须通过统一的配置包加载配置
- **FR-014**: 系统必须使用 logrus（`github.com/sirupsen/logrus`）作为结构化日志库，所有服务必须通过统一的 logger 包使用 logrus，支持 JSON 格式输出和 trace_id 记录
- **FR-015**: 系统必须使用 GORM 的 AutoMigrate 功能进行数据库迁移，不得使用手动 SQL 脚本，所有数据模型变更必须通过 GORM 迁移实现

### Key Entities

- **TraceContext**: 追踪上下文，包含 trace_id、span_id、parent_span_id 等
- **RoutingKey**: RabbitMQ routing key，格式为 `iot.{vendor}.{service}.{action}`，例如：`iot.dji.gateway.uplink.property`、`iot.dji.uplink.property.report`
- **TopicExchange**: RabbitMQ Topic Exchange，名称为 `iot`，用于服务间消息路由
- **Queue**: RabbitMQ Queue，每个服务创建自己的 Queue，通过 routing key 绑定到 `iot` Topic Exchange
- **MessageHeader**: 消息头，包含 W3C Trace Context（traceparent、tracestate）、message_type 等元数据
- **Vendor**: 厂商标识，如 "dji"、"generic" 等，存储在数据库中，通过 device_sn 查询获取
- **DeviceSN**: 设备序列号，每条消息必须包含（在 topic 或 payload 中），用于查询设备信息和 vendor
- **MetricsPackage**: 统一的 metrics 封装包，基于 Prometheus Go Client Library，提供统一的 API 和标签规范，支持框架中间件和业务代码使用

## Success Criteria

### Measurable Outcomes

- **SC-001**: 分布式追踪基础设施能够追踪所有服务间的消息流转，trace_id 传递成功率 > 99.9%
- **SC-002**: 多厂商 routing key 规范能够支持至少 3 个厂商（DJI、通用 MQTT、其他）
- **SC-003**: 服务间消息串联机制能够支持端到端消息流转，消息丢失率 < 0.1%
- **SC-004**: 在 Tempo 中能够查询到完整的调用链路，链路完整性 > 95%
- **SC-005**: 统一的 metrics 包能够自动收集框架中间件 metrics（RabbitMQ、PostgreSQL、InfluxDB），并支持业务代码使用，所有 metrics 在 `/metrics` 端点暴露

## Clarifications

### Session 2025-01-27

- Q: RabbitMQ routing key 的格式规范是什么？ → A: `iot.{vendor}.{service}.{action}` 格式，例如：`iot.dji.gateway.uplink.property`、`iot.dji.uplink.property.report`
- Q: 分布式追踪的实现方式是什么？ → A: OpenTelemetry SDK + W3C Trace Context 传播，在 HTTP 请求头和 RabbitMQ 消息头中传递 trace context
- Q: 厂商标识（vendor）如何确定？ → A: 设备上云流程中将 sn + vendor 写入数据库，每条消息都包含 sn（在 topic 或 payload 中），通过 sn 查询数据库获取 vendor，不需要在消息体中包含 vendor 字段
- Q: RabbitMQ Exchange 和 Queue 的架构设计是什么？ → A: Topic Exchange + 动态 Queue 绑定，使用 `iot` Topic Exchange，每个服务创建自己的 Queue，通过 routing key 绑定
- Q: W3C Trace Context 在 RabbitMQ 消息头中的传递方式是什么？ → A: 使用 RabbitMQ 消息头（headers）传递，在消息的 headers 字段中包含 `traceparent` 和 `tracestate`，与 HTTP 头传递方式一致
- Q: Metrics 包的技术选型和实现方式是什么？ → A: Prometheus Go Client Library + 统一封装包，使用 `github.com/prometheus/client_golang`，创建统一的封装包（如 `pkg/metrics`），提供统一的 API 和标签规范
- Q: 框架中间件 Metrics 的自动收集方式是什么？ → A: 统一封装包自动收集 + 手动注册，metrics 包提供中间件集成（RabbitMQ、PostgreSQL、InfluxDB），自动收集连接数、查询延迟等，同时支持手动注册自定义 metrics
- Q: Metrics 标签规范的定义方式是什么？ → A: 统一标签规范 + 可选业务标签，定义必需标签（service、vendor、message_type、status 等），同时支持业务代码添加自定义标签
- Q: Metrics 暴露端点的方式是什么？ → A: 统一 HTTP 端点 `/metrics`，所有服务在 `/metrics` 端点暴露 Prometheus 格式的 metrics，使用 Gin 中间件或独立 HTTP 服务器
- Q: Metrics 命名规范是什么？ → A: 统一命名规范 `iot_{component}_{metric_type}_{unit}`，例如：`iot_rabbitmq_connection_total`、`iot_postgres_query_duration_seconds`、`iot_message_processed_total`
- Q: 配置管理方式是什么？ → A: 使用 YAML 配置文件，支持多环境配置（dev、staging、prod），所有服务通过统一的配置包加载配置
- Q: 日志库的技术选型是什么？ → A: 使用 logrus（`github.com/sirupsen/logrus`）作为结构化日志库，所有服务通过统一的 logger 包使用 logrus
- Q: 数据库迁移的实现方式是什么？ → A: 使用 GORM 的 AutoMigrate 功能进行数据库迁移，不得使用手动 SQL 脚本，所有数据模型变更通过 GORM 迁移实现
- Q: Go 模块包名是什么？ → A: 使用 `github.com/utmos/utmos` 作为 Go 模块包名

