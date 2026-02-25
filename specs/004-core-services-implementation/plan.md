# Implementation Plan: Core Services Implementation

**Branch**: `004-core-services-implementation` | **Date**: 2025-02-05 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-core-services-implementation/spec.md`

## Summary

实现 5 个核心微服务（iot-gateway、iot-uplink、iot-downlink、iot-api、iot-ws）的业务逻辑，完成完整的上下行数据流。基于 001/002/003 已实现的基础设施和协议适配器框架。

## Technical Context

**Language/Version**: Go 1.22+
**Primary Dependencies**: Gin Framework, GORM, paho.mqtt.golang, gorilla/websocket, amqp091-go
**Storage**: PostgreSQL (业务数据), InfluxDB (时序数据)
**Testing**: go test, testify (assert/require/mock)
**Target Platform**: Linux server (Docker/Kubernetes)
**Project Type**: Microservices
**Performance Goals**: 消息处理延迟 < 100ms (P95), 1000+ 设备同时在线
**Constraints**: 只有 iot-gateway 可连接 VerneMQ，服务间通过 RabbitMQ 通信
**Scale/Scope**: 5 个核心服务，完整上下行数据流

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| Thing Model Driven Architecture | ✅ PASS | 使用物模型进行设备抽象 |
| Multi-Protocol Support | ✅ PASS | MQTT + HTTP + WebSocket |
| Device Abstraction Layer | ✅ PASS | 通过 dji-adapter 实现 |
| Extensibility & Plugin Architecture | ✅ PASS | 协议适配器可插拔 |
| Standardized API Design | ✅ PASS | RESTful API + OpenAPI |
| Test-First Development | ✅ PASS | TDD 原则 |
| Observability & Monitoring | ✅ PASS | 复用 001 可观测性基础设施 |
| Technology Stack Standards | ✅ PASS | Go + Gin + GORM |
| Microservice Architecture | ✅ PASS | 5 个核心服务 + RabbitMQ |
| Infrastructure Standards | ✅ PASS | VerneMQ + RabbitMQ + PostgreSQL + InfluxDB |

## Project Structure

### Documentation (this feature)

```text
specs/004-core-services-implementation/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (API contracts)
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
cmd/
├── iot-gateway/main.go      # Gateway 服务入口 (已有骨架)
├── iot-uplink/main.go       # Uplink 服务入口 (已有骨架)
├── iot-downlink/main.go     # Downlink 服务入口 (已有骨架)
├── iot-api/main.go          # API 服务入口 (已有骨架)
└── iot-ws/main.go           # WebSocket 服务入口 (已有骨架)

internal/
├── gateway/                 # iot-gateway 业务逻辑 (新增)
│   ├── mqtt/               # MQTT 客户端
│   │   ├── client.go       # MQTT 连接管理
│   │   ├── handler.go      # 消息处理
│   │   └── auth.go         # 设备认证
│   ├── bridge/             # MQTT↔RabbitMQ 桥接
│   │   ├── uplink.go       # 上行桥接
│   │   └── downlink.go     # 下行桥接
│   └── service.go          # 服务层
│
├── uplink/                  # iot-uplink 业务逻辑 (新增)
│   ├── processor/          # 消息处理器
│   │   ├── processor.go    # 处理器接口
│   │   └── dji.go          # DJI 消息处理
│   ├── storage/            # 数据存储
│   │   └── influx.go       # InfluxDB 写入
│   ├── router/             # 消息路由
│   │   └── router.go       # 路由逻辑
│   └── service.go          # 服务层
│
├── downlink/                # iot-downlink 业务逻辑 (新增)
│   ├── dispatcher/         # 消息分发
│   │   ├── dispatcher.go   # 分发器接口
│   │   └── dji.go          # DJI 消息分发
│   ├── retry/              # 重试机制
│   │   └── retry.go        # 重试逻辑
│   └── service.go          # 服务层
│
├── api/                     # iot-api 业务逻辑 (新增)
│   ├── handler/            # HTTP 处理器
│   │   ├── device.go       # 设备管理
│   │   ├── service.go      # 服务调用
│   │   └── telemetry.go    # 遥测查询
│   ├── middleware/         # 中间件
│   │   ├── auth.go         # 认证
│   │   └── trace.go        # 追踪
│   └── router.go           # 路由配置
│
├── ws/                      # iot-ws 业务逻辑 (新增)
│   ├── hub/                # 连接管理
│   │   ├── hub.go          # 连接中心
│   │   └── client.go       # 客户端连接
│   ├── subscription/       # 订阅管理
│   │   └── manager.go      # 订阅管理器
│   └── service.go          # 服务层
│
└── shared/                  # 共享代码 (已有)
    ├── config/
    ├── logger/
    ├── database/
    └── server/

pkg/
├── adapter/dji/             # DJI 协议适配器 (已有)
├── rabbitmq/                # RabbitMQ 客户端 (已有)
├── metrics/                 # Prometheus 指标 (已有)
├── tracer/                  # 分布式追踪 (已有)
└── models/                  # 数据模型 (已有)

tests/
├── integration/             # 集成测试 (新增)
│   ├── gateway_test.go
│   ├── uplink_test.go
│   ├── downlink_test.go
│   └── e2e_test.go
└── mocks/                   # Mock 数据 (已有部分)
```

**Structure Decision**: 采用微服务架构，每个服务在 `internal/` 下有独立的业务逻辑目录，共享代码在 `internal/shared/` 和 `pkg/`。

## Implementation Phases

### Phase 1: iot-gateway (P1)

**目标**: 实现 MQTT↔RabbitMQ 桥接

**任务**:
1. 实现 MQTT 客户端连接 VerneMQ
2. 实现设备认证（用户名/密码）
3. 实现上行消息桥接（MQTT → RabbitMQ）
4. 实现下行消息桥接（RabbitMQ → MQTT）
5. 实现设备连接状态管理

**依赖**: paho.mqtt.golang

### Phase 2: iot-uplink (P1)

**目标**: 实现上行消息处理

**任务**:
1. 实现 RabbitMQ 消息订阅
2. 集成 dji-adapter 进行消息解析
3. 实现 InfluxDB 时序数据写入
4. 实现消息路由到 iot-ws

**依赖**: Phase 1 完成, dji-adapter

### Phase 3: iot-downlink (P1)

**目标**: 实现下行消息路由

**任务**:
1. 实现 RabbitMQ 消息订阅
2. 集成 dji-adapter 进行消息转换
3. 实现消息确认和重试机制
4. 实现消息路由到 iot-gateway

**依赖**: Phase 1 完成, dji-adapter

### Phase 4: iot-api (P2)

**目标**: 实现 RESTful API

**任务**:
1. 实现设备管理 API (CRUD)
2. 实现服务调用 API
3. 实现遥测数据查询 API
4. 生成 OpenAPI 文档

**依赖**: Phase 2, Phase 3 完成

### Phase 5: iot-ws (P2)

**目标**: 实现 WebSocket 实时推送

**任务**:
1. 实现 WebSocket 连接管理
2. 实现消息订阅机制
3. 实现实时消息推送
4. 实现连接心跳检测

**依赖**: Phase 2 完成

### Phase 6: Integration Testing

**目标**: 验证完整数据流

**任务**:
1. 编写集成测试
2. 编写端到端测试
3. 性能测试

**依赖**: Phase 1-5 完成

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违规项。

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| MQTT 连接不稳定 | 高 | 实现自动重连和消息缓存 |
| RabbitMQ 消息丢失 | 高 | 使用持久化队列和确认机制 |
| InfluxDB 写入性能 | 中 | 批量写入，异步处理 |
| WebSocket 连接数过多 | 中 | 连接限制，负载均衡 |

## Dependencies

### External Libraries (需要添加)

```go
// go.mod additions
require (
    github.com/eclipse/paho.mqtt.golang v1.4.3  // MQTT 客户端
    github.com/gorilla/websocket v1.5.1         // WebSocket
    github.com/influxdata/influxdb-client-go/v2 v2.13.0  // InfluxDB 客户端
)
```

### Internal Dependencies

- `pkg/adapter/dji` - DJI 协议适配器 (003 已实现)
- `pkg/rabbitmq` - RabbitMQ 客户端 (001 已实现)
- `pkg/metrics` - Prometheus 指标 (001 已实现)
- `pkg/tracer` - 分布式追踪 (001 已实现)
- `internal/shared` - 共享基础设施 (001 已实现)
