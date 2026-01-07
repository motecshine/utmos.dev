# Implementation Plan: Project Setup

**Branch**: `001-project-setup` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-project-setup/spec.md`

## Summary

设置 UMOS IoT 平台的基础项目结构，包括5个微服务的项目骨架、开发环境配置、CI/CD 流程、代码质量检查工具等。**核心实现目标**：实现分布式追踪（Trace）基础设施（OpenTelemetry + Tempo），定义多厂商 RabbitMQ routing key 规范（`iot.{vendor}.{service}.{action}`），建立服务间消息路由机制（Topic Exchange + 动态 Queue），实现统一的 metrics 包（Prometheus + 统一封装），确保所有服务能够通过标准化的消息路由、追踪和指标监控机制串联起来。确保项目遵循宪法规定的技术栈（Go + Gin + GORM）、代码规范（Uber Go）、微服务架构和中间件选型。

## Technical Context

**Language/Version**: Go 1.22  
**Primary Dependencies**: 
- Gin Framework (HTTP/WebSocket)
- GORM (ORM)
- logrus (结构化日志库)
- RabbitMQ Client (消息队列)
- VerneMQ Client (MQTT, 仅 iot-gateway)
- PostgreSQL Driver (通过 GORM)
- InfluxDB Client (时序数据库)
- OpenTelemetry SDK (分布式追踪)
- Prometheus Go Client Library (指标监控)

**Storage**: 
- PostgreSQL (业务数据、设备元数据、物模型定义)
- InfluxDB (时序数据)

**Testing**: 
- Go testing package (单元测试)
- Testify (测试断言库)
- Mockery (Mock 生成工具，用于模拟设备)

**Target Platform**: Linux (Docker/Kubernetes 部署)  
**Project Type**: Microservices (5个独立服务)  
**Performance Goals**: 
- API 响应时间 P95 < 200ms
- 支持 10k+ 并发设备连接
- 消息处理吞吐量 > 10k msg/s

**Constraints**: 
- 必须遵循 Uber Go 编码规范
- 严禁拼写错误（Typo）
- 服务间通信必须通过 RabbitMQ
- 只有 iot-gateway 可以连接 VerneMQ
- RabbitMQ routing key 必须遵循 `iot.{vendor}.{service}.{action}` 格式
- 所有消息必须包含 W3C Trace Context（traceparent、tracestate）
- Metrics 命名必须遵循 `iot_{component}_{metric_type}_{unit}` 格式

**Scale/Scope**: 
- 5个微服务
- 支持多租户架构
- 可水平扩展

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### ✅ 技术栈合规性
- **Go 语言**: ✅ 使用 Go 1.21+
- **Gin Framework**: ✅ iot-api 和 iot-ws 使用 Gin
- **GORM**: ✅ 所有数据库访问使用 GORM

### ✅ 代码规范合规性
- **Uber Go 规范**: ✅ 必须遵循
- **命名规范**: ✅ 包、变量、函数、常量、类型、接口、文件名必须遵循规范
- **Typo 检查**: ✅ 必须启用 misspell 检查

### ✅ 微服务架构合规性
- **5个核心服务**: ✅ iot-api, iot-ws, iot-uplink, iot-downlink, iot-gateway
- **RabbitMQ 通信**: ✅ 服务间必须通过 RabbitMQ
- **MQTT 隔离**: ✅ 只有 iot-gateway 连接 VerneMQ

### ✅ 中间件合规性
- **消息队列**: ✅ VerneMQ, RabbitMQ
- **数据库**: ✅ PostgreSQL, InfluxDB
- **可观测性**: ✅ Prometheus, Loki, Tempo, Grafana

### ⚠️ 需要确认的事项
- Go 版本具体选择（1.21, 1.22, 1.23?）→ ✅ 已确认：Go 1.22
- 各服务的独立仓库还是 monorepo？→ ✅ 已确认：Monorepo
- CI/CD 平台选择（GitHub Actions, GitLab CI, Jenkins?）→ ✅ 已确认：GitHub Actions
- 配置管理方案（环境变量、配置中心？）→ ✅ 已确认：YAML 配置文件

## Project Structure

### Documentation (this feature)

```text
specs/001-project-setup/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (基础数据模型)
├── quickstart.md         # Phase 1 output (快速开始指南)
└── contracts/           # Phase 1 output (API 契约，如果有)
```

### Source Code (repository root)

```text
umos/
├── cmd/                          # 各服务的入口
│   ├── iot-api/
│   │   └── main.go
│   ├── iot-ws/
│   │   └── main.go
│   ├── iot-uplink/
│   │   └── main.go
│   ├── iot-downlink/
│   │   └── main.go
│   └── iot-gateway/
│       └── main.go
│
├── internal/                     # 内部包（不对外暴露）
│   ├── api/                      # iot-api 内部实现
│   │   ├── handler/              # HTTP handlers
│   │   ├── middleware/           # Gin 中间件
│   │   └── router/               # 路由配置
│   ├── ws/                       # iot-ws 内部实现
│   │   ├── connection/           # WebSocket 连接管理
│   │   └── handler/              # WebSocket handlers
│   ├── uplink/                   # iot-uplink 内部实现
│   │   ├── processor/            # 消息处理器
│   │   └── transformer/          # 消息转换器
│   ├── downlink/                 # iot-downlink 内部实现
│   │   ├── processor/            # 消息处理器
│   │   └── transformer/          # 消息转换器
│   ├── gateway/                  # iot-gateway 内部实现
│   │   ├── mqtt/                 # MQTT 客户端
│   │   ├── auth/                 # 设备认证
│   │   └── converter/            # MQTT ↔ RabbitMQ 转换
│   └── shared/                   # 共享代码
│       ├── config/               # 配置管理（YAML）
│       └── logger/               # 日志
│
├── pkg/                          # 可对外暴露的公共包
│   ├── models/                   # 数据模型（GORM models）
│   │   ├── device.go
│   │   ├── thing_model.go
│   │   └── message.go
│   ├── repository/               # 数据访问层（GORM）
│   │   ├── device.go
│   │   └── thing_model.go
│   ├── metrics/                 # 统一 metrics 包（Prometheus 封装）
│   │   ├── collector.go         # Metrics 收集器
│   │   ├── middleware.go        # 中间件 metrics（RabbitMQ、PostgreSQL、InfluxDB）
│   │   ├── business.go          # 业务 metrics API
│   │   └── handler.go           # HTTP handler for /metrics endpoint
│   ├── tracer/                  # 分布式追踪（OpenTelemetry）
│   │   ├── provider.go          # Tracer Provider
│   │   ├── http.go              # HTTP 追踪中间件
│   │   └── rabbitmq.go          # RabbitMQ 消息追踪
│   ├── rabbitmq/                # RabbitMQ 客户端封装
│   │   ├── client.go            # RabbitMQ 客户端
│   │   ├── routing.go           # Routing key 生成和解析
│   │   ├── exchange.go          # Exchange 和 Queue 管理
│   │   ├── message.go           # 标准消息格式
│   │   ├── publisher.go        # 消息发布
│   │   └── subscriber.go       # 消息订阅
│   └── errors/                   # 错误定义
│
├── api/                          # API 定义（OpenAPI/Swagger）
│   └── v1/
│       └── openapi.yaml
│
├── scripts/                     # 脚本
│   ├── build.sh                  # 构建脚本
│   ├── test.sh                   # 测试脚本
│   └── lint.sh                   # 代码检查脚本
│
├── deployments/                  # 部署配置
│   ├── docker/                   # Dockerfile
│   │   ├── iot-api.Dockerfile
│   │   ├── iot-ws.Dockerfile
│   │   ├── iot-uplink.Dockerfile
│   │   ├── iot-downlink.Dockerfile
│   │   └── iot-gateway.Dockerfile
│   └── kubernetes/               # K8s manifests
│       ├── iot-api/
│       ├── iot-ws/
│       ├── iot-uplink/
│       ├── iot-downlink/
│       └── iot-gateway/
│
├── tests/                        # 测试
│   ├── unit/                     # 单元测试
│   ├── integration/              # 集成测试
│   ├── contract/                 # 契约测试
│   └── mocks/                    # Mock 文件
│
├── docs/                         # 文档
│   ├── architecture/             # 架构文档
│   ├── api/                      # API 文档
│   └── development/              # 开发文档
│
├── .github/                      # GitHub Actions
│   └── workflows/
│       ├── ci.yml                # CI 流程
│       └── cd.yml                # CD 流程
│
├── .golangci.yml                 # golangci-lint 配置
├── go.mod                        # Go 模块定义
├── go.sum                        # Go 依赖校验
├── Makefile                      # Make 命令
├── docker-compose.yml            # 本地开发环境
└── README.md                     # 项目说明
```

**Structure Decision**: 采用 monorepo 结构，5个微服务共享公共代码（pkg/、internal/shared/），但各自独立部署。这种结构便于代码复用、统一管理和版本控制。

## Implementation Phases

### Phase 0: 基础设施搭建
1. **项目初始化**
   - 创建 monorepo 目录结构
   - 初始化 Go modules（go.mod）
   - 配置 golangci-lint 和 misspell
   - 设置 GitHub Actions CI/CD

2. **统一 Metrics 包实现** (`pkg/metrics`)
   - 基于 Prometheus Go Client Library 创建统一封装
   - 实现中间件 metrics 自动收集（RabbitMQ、PostgreSQL、InfluxDB）
   - 提供业务 metrics API（Counter、Histogram、Gauge）
   - 实现 `/metrics` HTTP 端点暴露
   - 定义统一标签规范（service、vendor、message_type、status）
   - 实现命名规范 `iot_{component}_{metric_type}_{unit}`

3. **分布式追踪基础设施** (`pkg/tracer`)
   - 集成 OpenTelemetry SDK
   - 实现 W3C Trace Context 传播（HTTP 和 RabbitMQ）
   - 创建 Tracer Provider 和配置
   - 实现 HTTP 追踪中间件（Gin）
   - 实现 RabbitMQ 消息追踪（消息头传递 traceparent、tracestate）
   - 配置 Tempo 导出器

4. **RabbitMQ 客户端封装** (`pkg/rabbitmq`)
   - 实现 RabbitMQ 客户端封装
   - 实现 Topic Exchange（`iot`）和 Queue 管理
   - 实现 routing key 生成和解析（`iot.{vendor}.{service}.{action}`）
   - 实现消息格式标准化（包含 W3C Trace Context）
   - 实现消息发布和订阅封装

### Phase 1: 服务骨架创建
1. **5个微服务入口** (`cmd/`)
   - 创建各服务的 main.go
   - 实现服务启动和优雅关闭
   - 集成配置管理、日志、metrics、tracer
   - 实现健康检查端点（`/health`、`/ready`）

2. **共享基础设施** (`internal/shared/`)
   - 配置管理（YAML 配置文件，支持多环境）
   - 结构化日志（使用 logrus，JSON 格式，包含 trace_id）
   - 错误处理统一封装

3. **数据模型** (`pkg/models/`)
   - Device、ThingModel、DeviceProperty、DeviceEvent、MessageLog
   - GORM 模型定义和 AutoMigrate 迁移

### Phase 2: 数据库和中间件集成
1. **PostgreSQL 集成**
   - GORM 连接和配置
   - GORM AutoMigrate 数据库迁移（不使用 SQL 脚本）
   - Repository 层实现

2. **InfluxDB 集成**
   - InfluxDB 客户端配置
   - 时序数据写入封装

3. **RabbitMQ 集成**
   - Exchange 和 Queue 初始化
   - 消息发布和订阅实现
   - 消息格式验证

### Phase 3: 测试和文档
1. **单元测试**
   - 各包单元测试（覆盖率 ≥ 80%）
   - Mock 设备生成

2. **集成测试**
   - 服务间消息流转测试
   - 分布式追踪端到端测试
   - Metrics 收集和暴露测试

3. **文档**
   - API 文档（OpenAPI/Swagger）
   - 架构文档更新
   - 开发指南

## Key Implementation Details

### 1. 统一 Metrics 包设计

**包结构**: `pkg/metrics/`
- `collector.go`: Metrics 收集器，管理 Prometheus Registry
- `middleware.go`: 中间件 metrics 自动收集
  - RabbitMQ: 连接数、消息数、延迟
  - PostgreSQL: 连接池、查询延迟、错误数
  - InfluxDB: 写入延迟、错误数
- `business.go`: 业务 metrics API
  - `NewCounter(name, labels)`: 创建计数器
  - `NewHistogram(name, labels)`: 创建直方图
  - `NewGauge(name, labels)`: 创建仪表盘

**标签规范**:
- 必需标签: `service`（服务名）、`vendor`（厂商）、`message_type`（消息类型）、`status`（状态）
- 可选标签: 业务自定义标签

**命名规范**: `iot_{component}_{metric_type}_{unit}`
- 示例: `iot_rabbitmq_connection_total`、`iot_postgres_query_duration_seconds`、`iot_message_processed_total`

### 2. 分布式追踪设计

**包结构**: `pkg/tracer/`
- `provider.go`: OpenTelemetry Tracer Provider 配置
- `http.go`: Gin HTTP 追踪中间件
  - 从 HTTP 请求头提取 W3C Trace Context
  - 创建新的 span
  - 在响应头中返回 trace_id
- `rabbitmq.go`: RabbitMQ 消息追踪
  - 在消息头中注入 W3C Trace Context（traceparent、tracestate）
  - 从消息头提取 Trace Context
  - 创建新的 span

**W3C Trace Context 传递**:
- HTTP: 通过 `traceparent` 和 `tracestate` 请求头
- RabbitMQ: 通过消息的 `headers` 字段（`traceparent`、`tracestate`）

### 3. RabbitMQ Routing Key 设计

**格式**: `iot.{vendor}.{service}.{action}`

**示例**:
- `iot.dji.gateway.uplink.property` - DJI 设备属性上报（gateway → uplink）
- `iot.dji.uplink.property.report` - DJI 属性上报（处理后）
- `iot.dji.downlink.service.call` - DJI 服务调用（下行）

**Vendor 获取**:
- 从消息中提取 device_sn（topic 或 payload）
- 查询数据库获取 vendor
- 用于生成 routing key

**Exchange 和 Queue**:
- Exchange: `iot` (Topic Exchange)
- Queue: 每个服务创建自己的 Queue，通过 routing key 绑定

### 4. 消息格式设计

**标准消息格式**:
```json
{
  "tid": "transaction-uuid",
  "bid": "business-uuid",
  "timestamp": 1234567890123,
  "service": "iot-api",
  "action": "device.query",
  "device_sn": "device-serial-number",
  "data": {
    // 具体业务数据
  }
}
```

**RabbitMQ 消息头**:
- `traceparent`: W3C Trace Context
- `tracestate`: W3C Trace State
- `message_type`: 消息类型（property、event、service）
- `vendor`: 厂商标识（可选，用于路由）

## Dependencies

### 核心依赖
- `github.com/gin-gonic/gin`: HTTP/WebSocket 框架
- `gorm.io/gorm`: ORM 框架
- `github.com/streadway/amqp`: RabbitMQ 客户端
- `github.com/prometheus/client_golang`: Prometheus 客户端
- `go.opentelemetry.io/otel`: OpenTelemetry SDK
- `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp`: Tempo 导出器

### 开发依赖
- `github.com/golangci/golangci-lint`: 代码检查
- `github.com/client9/misspell`: 拼写检查
- `github.com/stretchr/testify`: 测试断言
- `github.com/vektra/mockery`: Mock 生成

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违反宪法的设计决策。

