# Implementation Plan: Protocol Adapter Design

**Branch**: `002-protocol-adapter-design` | **Date**: 2025-02-05 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-protocol-adapter-design/spec.md`

## Summary

实现协议适配器框架和 DJI 协议适配器服务（`dji-adapter`），支持将 DJI 设备的原始消息转换为平台标准消息格式。协议适配器从 RabbitMQ 订阅原始消息（由 iot-gateway 发布），解析协议特定格式，转换为标准消息后发布到 RabbitMQ 供 iot-uplink 消费。同时支持下行消息的反向转换。

**核心实现目标**：
1. 定义协议适配器接口规范（`pkg/adapter/`）
2. 实现 DJI 协议解析器（`pkg/adapter/dji/`）
3. 创建 `dji-adapter` 服务（`cmd/dji-adapter/`）
4. 扩展 001 的 `StandardMessage` 增加 `protocol_meta` 字段
5. 定义原始消息队列规范（`iot.raw.{vendor}.{uplink|downlink}`）

## Technical Context

**Language/Version**: Go 1.22 (与 001 保持一致)
**Primary Dependencies**:
- 复用 001: Gin Framework, GORM, logrus, RabbitMQ Client, OpenTelemetry, Prometheus
- 新增: 无新依赖，完全复用 001 基础设施

**Storage**:
- PostgreSQL (读取设备配置、物模型定义)
- 复用 001 的数据库连接层

**Testing**:
- Go testing package (单元测试)
- Testify (测试断言)
- 模拟消息测试（Mock RabbitMQ 消息）

**Target Platform**: Linux (Docker/Kubernetes 部署)
**Project Type**: Microservices (扩展服务)
**Performance Goals**:
- 消息转换延迟 P95 < 50ms
- 消息处理成功率 > 99.9%

**Constraints**:
- 必须遵循 Uber Go 编码规范
- 严禁拼写错误（Typo）
- 协议适配器不直接连接 VerneMQ，只从 RabbitMQ 订阅消息
- 必须复用 001 的基础设施（tracer、metrics、logger、config）

**Scale/Scope**:
- 1 个新服务 (`dji-adapter`)
- 1 个新 pkg 包 (`pkg/adapter/`)
- 扩展现有 `pkg/rabbitmq/message.go`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### ✅ 技术栈合规性
- **Go 语言**: ✅ 使用 Go 1.22
- **Gin Framework**: ✅ dji-adapter 健康检查使用 Gin
- **GORM**: ✅ 读取设备配置使用 GORM

### ✅ 代码规范合规性
- **Uber Go 规范**: ✅ 必须遵循
- **命名规范**: ✅ 包、变量、函数、常量、类型、接口、文件名必须遵循规范
- **Typo 检查**: ✅ 必须启用 misspell 检查

### ✅ 微服务架构合规性
- **扩展服务**: ✅ dji-adapter 是扩展服务，独立于5个核心服务
- **RabbitMQ 通信**: ✅ 从 RabbitMQ 订阅原始消息
- **MQTT 隔离**: ✅ 不直接连接 VerneMQ，通过 iot-gateway 桥接

### ✅ 物模型驱动架构
- **统一物模型**: ✅ 所有协议共享统一的物模型定义
- **协议映射**: ✅ dji-adapter 负责将 DJI 格式映射到标准物模型

### ✅ 可观测性支持
- **分布式追踪**: ✅ 复用 001 的 pkg/tracer
- **指标监控**: ✅ 复用 001 的 pkg/metrics
- **结构化日志**: ✅ 复用 001 的 internal/shared/logger

### ✅ 测试优先开发
- **TDD 原则**: ✅ 先写测试，确保失败，再实现
- **模拟消息**: ✅ 提供模拟 DJI 原始消息用于测试
- **契约测试**: ✅ 协议适配器接口必须有契约测试

## Project Structure

### Documentation (this feature)

```text
specs/002-protocol-adapter-design/
├── plan.md              # This file
├── spec.md              # Feature specification
├── research.md          # Phase 0 output (DJI 协议研究)
├── data-model.md        # Phase 1 output (协议消息模型)
├── quickstart.md        # Phase 1 output (快速开始指南)
├── contracts/           # Phase 1 output (协议适配器接口契约)
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
umos/
├── cmd/
│   ├── iot-api/              # [001] 已实现
│   ├── iot-ws/               # [001] 已实现
│   ├── iot-uplink/           # [001] 已实现
│   ├── iot-downlink/         # [001] 已实现
│   ├── iot-gateway/          # [001] 已实现
│   └── dji-adapter/          # [002] 新增 - DJI 协议适配器服务
│       └── main.go
│
├── internal/
│   └── shared/               # [001] 已实现，复用
│
├── pkg/
│   ├── adapter/              # [002] 新增 - 协议适配器框架
│   │   ├── interface.go      # 协议适配器接口定义
│   │   ├── registry.go       # 协议适配器注册表
│   │   ├── raw_message.go    # 原始消息定义
│   │   └── dji/              # DJI 协议适配器实现
│   │       ├── adapter.go    # DJI 适配器实现
│   │       ├── parser.go     # DJI 消息解析器
│   │       ├── converter.go  # DJI ↔ 标准消息转换器
│   │       ├── topic.go      # DJI Topic 解析
│   │       └── types.go      # DJI 协议类型定义
│   │
│   ├── rabbitmq/             # [001] 已实现，需扩展
│   │   └── message.go        # 扩展 StandardMessage 增加 protocol_meta
│   │
│   ├── tracer/               # [001] 已实现，复用
│   ├── metrics/              # [001] 已实现，复用
│   ├── models/               # [001] 已实现，复用
│   ├── repository/           # [001] 已实现，复用
│   └── errors/               # [001] 已实现，复用
│
├── tests/
│   ├── unit/                 # [001] 已有结构
│   ├── integration/          # [001] 已有结构
│   │   └── dji_adapter_test.go  # [002] 新增
│   └── mocks/
│       └── dji_messages.go   # [002] 新增 - DJI 模拟消息
│
└── deployments/
    └── docker/
        └── dji-adapter.Dockerfile  # [002] 新增
```

**Structure Decision**: 采用 `pkg/adapter/` 包结构，将协议适配器框架作为公共包，各厂商适配器作为子包（如 `pkg/adapter/dji/`）。这种结构便于后续扩展其他厂商适配器（如 `pkg/adapter/tuya/`）。

## Implementation Phases

### Phase 0: 研究和准备

1. **DJI 协议研究**
   - 分析 `docs/dji` 目录下的 DJI Cloud API 文档
   - 理解 DJI MQTT Topic 结构和消息格式
   - 识别核心 Topic：osd、state、services、events、status
   - 记录协议特定的元数据（如 method、need_reply 等）

2. **001 基础设施评估**
   - 确认 001 实现的可复用组件
   - 识别需要扩展的部分（StandardMessage、RoutingKey）

### Phase 1: 框架设计

1. **协议适配器接口定义** (`pkg/adapter/interface.go`)
   - `ProtocolAdapter` 接口
   - `ParseRawMessage(raw []byte) (*ProtocolMessage, error)`
   - `ToStandardMessage(pm *ProtocolMessage) (*StandardMessage, error)`
   - `FromStandardMessage(sm *StandardMessage) ([]byte, error)`
   - `GetVendor() string`

2. **原始消息定义** (`pkg/adapter/raw_message.go`)
   - `RawMessage` 结构，包含原始 payload 和元数据
   - 原始 MQTT Topic、QoS、时间戳等

3. **扩展 StandardMessage** (`pkg/rabbitmq/message.go`)
   - 新增 `ProtocolMeta` 字段
   - `Vendor`、`OriginalTopic`、`QoS` 等元数据

4. **RabbitMQ 队列规范**
   - 原始消息上行：`iot.raw.{vendor}.uplink`
   - 原始消息下行：`iot.raw.{vendor}.downlink`
   - 标准消息：`iot.{vendor}.{service}.{action}`（复用 001）

### Phase 2: DJI 适配器实现

1. **DJI Topic 解析器** (`pkg/adapter/dji/topic.go`)
   - 解析 `thing/product/{device_sn}/osd` 等 Topic
   - 提取 device_sn、gateway_sn、消息类型

2. **DJI 消息解析器** (`pkg/adapter/dji/parser.go`)
   - 解析 DJI JSON 消息格式
   - 提取 tid、bid、timestamp、method、data 等字段

3. **DJI 消息转换器** (`pkg/adapter/dji/converter.go`)
   - DJI 消息 → 标准消息（上行）
   - 标准消息 → DJI 消息（下行）

4. **DJI 适配器服务** (`cmd/dji-adapter/main.go`)
   - 订阅 `iot.raw.dji.uplink` 队列
   - 订阅 `iot.dji.*.downlink` 队列（下行）
   - 解析、转换、发布消息

### Phase 3: 测试和集成

1. **单元测试**
   - Topic 解析测试
   - 消息解析测试
   - 消息转换测试
   - 覆盖率 ≥ 80%

2. **集成测试**
   - 端到端消息流转测试
   - 模拟 iot-gateway 发送原始消息
   - 验证 dji-adapter 正确转换并发布标准消息

3. **文档和部署**
   - Dockerfile
   - 配置文件示例
   - 开发文档

## Key Implementation Details

### 1. 协议适配器接口设计

```go
// pkg/adapter/interface.go
package adapter

type ProtocolAdapter interface {
    // GetVendor returns the vendor identifier (e.g., "dji", "tuya")
    GetVendor() string

    // ParseRawMessage parses raw bytes into a protocol-specific message
    ParseRawMessage(topic string, payload []byte) (*ProtocolMessage, error)

    // ToStandardMessage converts protocol message to standard message
    ToStandardMessage(pm *ProtocolMessage) (*rabbitmq.StandardMessage, error)

    // FromStandardMessage converts standard message to protocol format
    FromStandardMessage(sm *rabbitmq.StandardMessage) (*ProtocolMessage, error)

    // GetRawPayload returns the raw payload for sending to device
    GetRawPayload(pm *ProtocolMessage) ([]byte, error)
}
```

### 2. 原始消息结构

```go
// pkg/adapter/raw_message.go
package adapter

type RawMessage struct {
    Vendor      string            `json:"vendor"`
    Topic       string            `json:"topic"`
    Payload     []byte            `json:"payload"`
    QoS         int               `json:"qos"`
    Timestamp   int64             `json:"timestamp"`
    Headers     map[string]string `json:"headers"`
}
```

### 3. 扩展 StandardMessage

```go
// pkg/rabbitmq/message.go (扩展)
type ProtocolMeta struct {
    Vendor        string `json:"vendor"`
    OriginalTopic string `json:"original_topic,omitempty"`
    QoS           *int   `json:"qos,omitempty"`
    Method        string `json:"method,omitempty"`
    NeedReply     *bool  `json:"need_reply,omitempty"`
}

type StandardMessage struct {
    // ... 现有字段 ...
    ProtocolMeta *ProtocolMeta `json:"protocol_meta,omitempty"`
}
```

### 4. DJI Topic 解析

```go
// pkg/adapter/dji/topic.go
package dji

// DJI Topic patterns:
// thing/product/{device_sn}/osd       - 属性定时上报
// thing/product/{device_sn}/state     - 属性变化上报
// thing/product/{gateway_sn}/services - 服务调用
// thing/product/{gateway_sn}/events   - 事件上报
// sys/product/{gateway_sn}/status     - 设备状态

type TopicInfo struct {
    Type      TopicType // osd, state, services, events, status
    DeviceSN  string
    GatewaySN string
    Direction Direction // uplink, downlink
}

func ParseTopic(topic string) (*TopicInfo, error)
```

### 5. 消息流转

```
上行流程:
设备 → VerneMQ → iot-gateway → RabbitMQ(iot.raw.dji.uplink)
                                        ↓
                               dji-adapter (解析转换)
                                        ↓
                               RabbitMQ(iot.dji.device.property.report)
                                        ↓
                                   iot-uplink

下行流程:
iot-api → RabbitMQ → iot-downlink → RabbitMQ(iot.dji.service.call)
                                            ↓
                                   dji-adapter (反向转换)
                                            ↓
                                   RabbitMQ(iot.raw.dji.downlink)
                                            ↓
                                   iot-gateway → VerneMQ → 设备
```

## Dependencies

### 复用 001 的依赖
- `github.com/gin-gonic/gin`: HTTP 框架（健康检查）
- `gorm.io/gorm`: ORM 框架（读取配置）
- `github.com/streadway/amqp`: RabbitMQ 客户端
- `github.com/prometheus/client_golang`: Prometheus 客户端
- `go.opentelemetry.io/otel`: OpenTelemetry SDK
- `github.com/sirupsen/logrus`: 日志库

### 新增依赖
- 无新增依赖

## 001 需要的扩展（前置工作）

在开始 002 实现前，需要对 001 进行以下小范围扩展：

| 文件 | 修改内容 | 估计行数 |
|------|----------|----------|
| `pkg/rabbitmq/message.go` | 添加 `ProtocolMeta` 结构和字段 | +30 行 |
| `pkg/rabbitmq/routing.go` | 添加原始消息队列常量 | +10 行 |
| `cmd/iot-gateway/main.go` | 添加发布原始消息到 RabbitMQ 的逻辑 | +50 行 |

**总计**: ~90 行扩展，向后兼容

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违反宪法的设计决策。

## Next Steps

1. 运行 `/speckit.tasks` 生成详细任务列表
2. 先完成 001 的扩展工作（~90 行）
3. 按任务列表实现 002 功能
