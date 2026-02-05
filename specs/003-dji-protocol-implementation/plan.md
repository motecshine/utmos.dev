# Implementation Plan: DJI Protocol Implementation

**Branch**: `003-dji-protocol-implementation` | **Date**: 2025-02-05 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-dji-protocol-implementation/spec.md`
**Depends On**: `001-project-setup`, `002-protocol-adapter-design`

## Summary

基于 002 实现的协议适配器框架，完整实现 DJI Cloud API 协议的上下行消息处理。主要工作是将 `pkg/adapter/dji/protocol/` 目录下已定义的数据结构集成到现有的 adapter 中，实现具体的消息解析、路由和转换逻辑。

**核心目标**:
1. 集成 OSD/State/Status 消息的完整解析
2. 实现服务调用的双向转换和路由
3. 实现事件处理的双向转换
4. 支持航线任务、相机控制等业务功能

## Technical Context

**Language/Version**: Go 1.22+
**Primary Dependencies**:
- 复用 001: `pkg/tracer`, `pkg/metrics`, `pkg/rabbitmq`, `internal/shared/logger`, `internal/shared/config`
- 复用 002: `pkg/adapter` (ProtocolAdapter 接口、Registry、Factory)
- 已有: `pkg/adapter/dji/protocol/*` (DJI 协议数据结构)
- 已有: `pkg/adapter/dji/wpml` (WPML 航线解析)

**Storage**:
- PostgreSQL (设备状态、任务记录)
- InfluxDB (OSD 时序数据)
- RabbitMQ (消息队列)

**Testing**: `go test` with TDD, coverage >= 80%

**Target Platform**: Linux server (Docker/Kubernetes)

**Project Type**: Microservice (dji-adapter)

**Performance Goals**:
- 消息处理延迟 < 50ms (P95)
- 支持 1000+ 设备同时在线
- 消息处理成功率 > 99.9%

**Constraints**:
- 必须复用 001/002 基础设施
- 必须遵循宪法定义的微服务架构
- 只有 iot-gateway 可连接 VerneMQ

**Scale/Scope**:
- 12 个协议模块 (aircraft, camera, device, wayline, drc, file, firmware, live, psdk, safety, config, common)
- 50+ 服务命令类型
- 20+ 事件类型

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Thing Model Driven Architecture | ✅ PASS | DJI 协议基于物模型 (Properties, Services, Events) |
| II. Multi-Protocol Support | ✅ PASS | 通过 dji-adapter 适配 DJI MQTT 协议 |
| III. Device Abstraction Layer | ✅ PASS | 使用 ProtocolAdapter 接口抽象 |
| IV. Extensibility & Plugin Architecture | ✅ PASS | 基于 002 的适配器框架扩展 |
| V. Standardized API Design | ✅ PASS | 使用 StandardMessage 统一消息格式 |
| VI. Test-First Development | ✅ PASS | TDD 原则，覆盖率 >= 80% |
| VII. Observability & Monitoring | ✅ PASS | 复用 pkg/tracer, pkg/metrics |
| Language & Code Standards | ✅ PASS | Uber Go 规范，golangci-lint |
| Technology Stack Standards | ✅ PASS | Go + Gin + GORM |
| Microservice Architecture | ✅ PASS | dji-adapter 作为独立服务 |
| Infrastructure Standards | ✅ PASS | RabbitMQ + PostgreSQL + InfluxDB |

**Gate Result**: ✅ ALL PASS - 可以继续 Phase 0

## Project Structure

### Documentation (this feature)

```text
specs/003-dji-protocol-implementation/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (service method definitions)
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
# 已有代码 (001/002 实现)
pkg/
├── adapter/
│   ├── interface.go           # ProtocolAdapter 接口 ✅
│   ├── registry.go            # 适配器注册表 ✅
│   ├── factory.go             # 适配器工厂 ✅
│   └── dji/
│       ├── adapter.go         # DJI 适配器 ✅ (需扩展)
│       ├── topic.go           # Topic 解析 ✅
│       ├── parser.go          # 消息解析 ✅ (需扩展)
│       ├── converter.go       # 消息转换 ✅ (需扩展)
│       ├── types.go           # 基础类型 ✅
│       ├── errors.go          # 错误定义 ✅
│       ├── protocol/          # 协议数据结构 ✅ (需集成)
│       │   ├── aircraft/      # 飞行器 OSD
│       │   ├── camera/        # 相机控制
│       │   ├── common/        # 公共结构
│       │   ├── config/        # 配置管理
│       │   ├── device/        # 设备控制
│       │   ├── drc/           # 实时控制
│       │   ├── file/          # 文件管理
│       │   ├── firmware/      # 固件升级
│       │   ├── live/          # 实时视频
│       │   ├── psdk/          # PSDK 负载
│       │   ├── safety/        # 安全相关
│       │   └── wayline/       # 航线任务
│       └── wpml/              # WPML 航线解析 ✅
├── rabbitmq/                  # RabbitMQ 客户端 ✅
├── tracer/                    # 分布式追踪 ✅
├── metrics/                   # Prometheus 指标 ✅
└── models/                    # GORM 模型 ✅

cmd/
└── dji-adapter/
    └── main.go                # 服务入口 ✅ (需扩展)

# 003 新增/修改
pkg/adapter/dji/
├── handler/                   # 新增: 消息处理器
│   ├── osd_handler.go         # OSD 消息处理
│   ├── state_handler.go       # State 消息处理
│   ├── status_handler.go      # Status 消息处理
│   ├── service_handler.go     # Service 消息处理
│   ├── event_handler.go       # Event 消息处理
│   └── request_handler.go     # Request 消息处理
├── router/                    # 新增: 服务路由
│   ├── service_router.go      # 服务调用路由
│   └── event_router.go        # 事件路由
└── integration/               # 新增: 协议集成
    ├── osd_parser.go          # OSD 数据解析
    ├── command_builder.go     # 命令构建器
    └── event_parser.go        # 事件解析

tests/
├── integration/
│   └── dji_adapter_test.go    # 集成测试 ✅ (需扩展)
└── mocks/
    └── dji_messages.go        # Mock 消息 ✅ (需扩展)
```

**Structure Decision**: 基于现有 002 结构扩展，新增 `handler/`, `router/`, `integration/` 子目录组织新代码。

## Complexity Tracking

> 无宪法违规，无需记录

## Phase 0: Research Tasks

### R1: 修复 protocol 目录 import 路径问题

**问题**: `pkg/adapter/dji/protocol/device/commands.go` 等文件使用外部 import 路径
```go
import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"
```
**需要**: 改为本地路径 `github.com/utmos/utmos/pkg/adapter/dji/protocol/common`

### R2: 分析 OSD 数据结构与 adapter 集成方式

**问题**: `protocol/aircraft/osd.go` 定义了 OSD 结构，但 adapter 未使用
**需要**: 研究如何在 `parser.go` 中集成 OSD 解析

### R3: 分析服务调用路由机制

**问题**: 需要根据 `method` 字段路由到具体命令处理器
**需要**: 设计 ServiceRouter 将 method 映射到 Command 类型

### R4: 分析事件类型与处理机制

**问题**: `protocol/*/events.go` 定义了事件类型，需要集成
**需要**: 设计 EventRouter 处理不同事件类型

## Phase 1: Design Artifacts

### D1: data-model.md
- DJI 消息类型枚举
- OSD 数据结构映射
- 服务命令注册表
- 事件类型注册表

### D2: contracts/
- `service_methods.yaml` - 所有服务方法定义
- `event_types.yaml` - 所有事件类型定义
- `osd_schema.yaml` - OSD 数据 schema

### D3: quickstart.md
- 如何测试 OSD 消息处理
- 如何测试服务调用
- 如何添加新的服务命令

## Implementation Phases

详细任务列表见 [tasks.md](./tasks.md)

### Phase 2: Foundational (P0 - Blocking)
- Handler/Router 接口定义
- 配置常量 (超时、策略)
- HandlerRegistry 实现

### Phase 3-6: 核心协议 (P1 - MVP)

| Phase | 功能 | 任务数 | 工作量 |
|-------|------|--------|--------|
| 3 | OSD 完整解析 | 7 | 2 天 |
| 4 | State/Status 处理 | 5 | 1 天 |
| 5 | Services 路由 (16 方法) | 7 | 2 天 |
| 6 | Events 路由 (4 事件) | 6 | 1 天 |

**P1 总计**: 25 任务, 6 天

### Phase 7-9: 业务功能 (P2)

| Phase | 功能 | 任务数 | 工作量 |
|-------|------|--------|--------|
| 7 | Wayline 航线任务 (8 方法 + 3 事件) | 5 | 1 天 |
| 8 | Camera 相机控制 (7 方法) | 3 | 0.5 天 |
| 9 | Config 配置管理 (2 方法) | 4 | 0.5 天 |

**P2 总计**: 12 任务, 2 天

### Phase 10-12: 高级功能 (P3)

| Phase | 功能 | 任务数 | 工作量 |
|-------|------|--------|--------|
| 10 | DRC 实时控制 (5 方法 + 2 事件) | 7 | 1.5 天 |
| 11 | File 文件管理 (3 方法 + 3 事件) | 4 | 0.5 天 |
| 12 | Firmware/Live (5 方法 + 1 事件) | 5 | 0.5 天 |

**P3 总计**: 16 任务, 2.5 天

### Phase 13: Polish & Integration
- 集成测试 (4 个)
- 性能测试 (1000 设备)
- 文档更新
- **预计工作量**: 1 天

## Summary

| 优先级 | Phases | 任务数 | 工作量 | 累计 |
|--------|--------|--------|--------|------|
| Setup | 1 | 4 | - | - |
| P0 | 2 | 7 | 1 天 | 1 天 |
| P1 | 3-6 | 25 | 6 天 | 7 天 |
| P2 | 7-9 | 12 | 2 天 | 9 天 |
| P3 | 10-12 | 16 | 2.5 天 | 11.5 天 |
| Polish | 13 | 9 | 1 天 | 12.5 天 |

**总计**: 73 任务 (含 4 个已完成的 Setup 任务), 约 12.5 天

## Milestones

1. **Milestone 1**: OSD 可用 (Phase 0-1) → 3 天
2. **Milestone 2**: 核心协议完成 (Phase 2-4) → 7 天
3. **Milestone 3**: 业务功能完成 (Phase 5-7) → 9 天
4. **Milestone 4**: 高级功能完成 (Phase 8-10) → 11.5 天
5. **Milestone 5**: 生产就绪 (Phase 11) → 12.5 天

## Next Steps

1. ✅ 运行 Phase 0 研究任务 (research.md 已完成)
2. ✅ 生成 data-model.md 和 contracts/ (已完成)
3. ✅ 生成 tasks.md (已完成)
4. 开始 Phase 0 实现
