# Feature Specification: DJI Protocol Implementation

**Feature Branch**: `003-dji-protocol-implementation`
**Created**: 2025-02-05
**Updated**: 2025-02-05
**Status**: Draft
**Input**: User description: "完整实现 DJI vendor protocol uplink downlink"
**Depends On**: `002-protocol-adapter-design` (协议适配器框架)

## Overview

基于 002 实现的协议适配器框架，完整实现 DJI Cloud API 协议的上下行消息处理。本 Spec 覆盖 DJI 协议的所有核心功能模块，包括设备管理、遥测数据、服务调用、事件处理、航线管理、媒体管理、实时控制等。

**数据流架构**:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           上行数据流 (Uplink)                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  DJI Device ──MQTT──► VerneMQ ──► iot-gateway ──RMQ──► dji-adapter         │
│                                                                             │
│  MQTT Topics (设备发布):                                                     │
│  ├── thing/product/{gateway_sn}/osd          → 遥测数据定时上报              │
│  ├── thing/product/{gateway_sn}/state        → 属性变化上报                  │
│  ├── thing/product/{gateway_sn}/events       → 事件上报                      │
│  ├── thing/product/{gateway_sn}/services_reply → 服务调用响应                │
│  ├── thing/product/{gateway_sn}/requests     → 设备主动请求                  │
│  └── sys/product/{gateway_sn}/status         → 设备上下线状态                │
│                                                                             │
│  dji-adapter 处理:                                                          │
│  ├── 解析 DJI 协议格式 (JSON)                                                │
│  ├── 提取设备标识 (gateway_sn, device_sn)                                    │
│  ├── 转换为 StandardMessage                                                 │
│  └── 发布到 RabbitMQ: iot.dji.{service}.{action}                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                           下行数据流 (Downlink)                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  iot-api ──RMQ──► iot-downlink ──RMQ──► dji-adapter ──RMQ──► iot-gateway   │
│                                                                             │
│  dji-adapter 处理:                                                          │
│  ├── 订阅 RabbitMQ: iot.dji.service.call                                    │
│  ├── 解析 StandardMessage                                                   │
│  ├── 转换为 DJI 协议格式                                                     │
│  └── 发布到 RabbitMQ: iot.raw.dji.downlink                                  │
│                                                                             │
│  MQTT Topics (平台发布):                                                     │
│  ├── thing/product/{gateway_sn}/services     → 服务调用                      │
│  ├── thing/product/{gateway_sn}/events_reply → 事件响应                      │
│  ├── thing/product/{gateway_sn}/requests_reply → 请求响应                    │
│  ├── thing/product/{gateway_sn}/property/set → 属性设置                      │
│  └── thing/product/{gateway_sn}/drc/down     → 实时控制指令                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Architecture

### DJI Cloud API 物模型 (Thing Model)

DJI 协议基于物模型 (TSL - Thing Specification Language) 设计，包含三大核心概念：

| 概念 | 说明 | 数据流向 |
|------|------|----------|
| **Properties** | 设备属性，如电量、位置、状态 | 设备 → 平台 (上报)，平台 → 设备 (设置) |
| **Services** | 服务调用，如起飞、降落、拍照 | 平台 → 设备 (调用)，设备 → 平台 (响应) |
| **Events** | 事件通知，如低电量告警、任务完成 | 设备 → 平台 (上报)，平台 → 设备 (确认) |

### 协议模块划分

基于 `pkg/adapter/dji/protocol` 目录结构，DJI 协议分为以下模块：

| 模块 | 目录 | 功能 | 优先级 |
|------|------|------|--------|
| **device** | `protocol/device` | 设备管理、拓扑、固件 | P1 |
| **aircraft** | `protocol/aircraft` | 飞行器遥测、控制 | P1 |
| **common** | `protocol/common` | 公共数据结构、错误码 | P1 |
| **config** | `protocol/config` | 配置管理 | P2 |
| **camera** | `protocol/camera` | 相机控制、拍照、录像 | P2 |
| **wayline** | `protocol/wayline` | 航线管理、任务执行 | P2 |
| **file** | `protocol/file` | 文件上传下载 | P3 |
| **live** | `protocol/live` | 实时视频流 | P3 |
| **drc** | `protocol/drc` | 实时控制 (DRC) | P3 |
| **firmware** | `protocol/firmware` | 固件升级 | P3 |
| **safety** | `protocol/safety` | 安全相关 | P3 |
| **psdk** | `protocol/psdk` | PSDK 负载 | P4 |

### 消息格式

**DJI 原始消息格式** (MQTT Payload):
```json
{
  "tid": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "bid": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "timestamp": 1234567890123,
  "gateway": "gateway_sn",
  "data": {
    // 协议特定数据
  },
  "method": "service_method_name"
}
```

**标准消息格式** (StandardMessage):
```json
{
  "tid": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "bid": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "timestamp": 1234567890123,
  "device_sn": "device_serial_number",
  "service": "device|aircraft|camera|wayline|...",
  "action": "property.report|service.call|event|...",
  "data": {},
  "protocol_meta": {
    "vendor": "dji",
    "original_topic": "thing/product/{sn}/osd",
    "qos": 1,
    "method": "original_method_name"
  }
}
```

## User Scenarios & Testing

### User Story 1 - 设备遥测数据上报 (Priority: P1)

作为平台运维人员，我需要实时接收 DJI 设备的遥测数据（OSD），以便监控设备状态和位置。

**Why this priority**: 遥测数据是 IoT 平台的核心功能，是所有其他功能的基础。

**Independent Test**: 可以通过模拟 DJI 设备发送 OSD 消息，验证 dji-adapter 能够正确解析并转换为标准消息。

**Acceptance Scenarios**:

1. **Given** DJI 机场/遥控器发布 OSD 消息到 `thing/product/{gateway_sn}/osd`, **When** dji-adapter 接收到该消息, **Then** 能够解析飞行器位置、高度、速度、电量等数据，转换为 StandardMessage 并发布到 `iot.dji.aircraft.osd`
2. **Given** OSD 消息包含嵌套设备数据（如机场+飞行器）, **When** dji-adapter 解析消息, **Then** 能够正确提取各设备的遥测数据并分别处理
3. **Given** OSD 消息频率为 0.5Hz (每2秒一次), **When** 持续接收消息, **Then** 处理延迟 < 50ms (P95)

---

### User Story 2 - 设备属性变化上报 (Priority: P1)

作为平台运维人员，我需要接收设备属性变化通知，以便及时响应设备状态变更。

**Why this priority**: 属性变化是事件驱动架构的核心，用于触发业务逻辑。

**Independent Test**: 可以通过模拟设备属性变化消息，验证 dji-adapter 能够正确解析并转换。

**Acceptance Scenarios**:

1. **Given** DJI 设备属性发生变化（如飞行模式切换）, **When** 设备发布 state 消息到 `thing/product/{gateway_sn}/state`, **Then** dji-adapter 能够解析变化的属性并发布到 `iot.dji.device.state`
2. **Given** 多个属性同时变化, **When** dji-adapter 接收消息, **Then** 能够正确解析所有变化的属性

---

### User Story 3 - 设备上下线状态管理 (Priority: P1)

作为平台运维人员，我需要知道设备的在线/离线状态，以便管理设备连接。

**Why this priority**: 设备状态是设备管理的基础功能。

**Independent Test**: 可以通过模拟设备上线/下线消息，验证状态更新。

**Acceptance Scenarios**:

1. **Given** DJI 设备连接到 VerneMQ, **When** 设备发布 status 消息到 `sys/product/{gateway_sn}/status`, **Then** dji-adapter 能够解析设备拓扑（机场、飞行器、遥控器）并更新在线状态
2. **Given** 设备断开连接, **When** VerneMQ 检测到断开, **Then** 平台能够收到离线通知并更新状态

---

### User Story 4 - 服务调用下发 (Priority: P1)

作为平台操作员，我需要向 DJI 设备下发控制指令（如起飞、降落、返航），以便远程控制设备。

**Why this priority**: 服务调用是平台控制设备的核心能力。

**Independent Test**: 可以通过 API 发送服务调用请求，验证 dji-adapter 能够正确转换并下发。

**Acceptance Scenarios**:

1. **Given** 平台发送服务调用请求（如 `flighttask_prepare`）, **When** dji-adapter 接收到 StandardMessage, **Then** 能够转换为 DJI 协议格式并发布到 `thing/product/{gateway_sn}/services`
2. **Given** 设备返回服务调用响应, **When** dji-adapter 接收到 `services_reply`, **Then** 能够解析响应结果并发布到 `iot.dji.service.reply`
3. **Given** 服务调用超时, **When** 超过配置的超时时间, **Then** 平台能够收到超时错误

---

### User Story 5 - 事件上报与确认 (Priority: P1)

作为平台运维人员，我需要接收设备事件通知（如告警、任务完成），并能够确认事件。

**Why this priority**: 事件是设备主动通知平台的重要机制。

**Independent Test**: 可以通过模拟设备事件消息，验证事件处理流程。

**Acceptance Scenarios**:

1. **Given** DJI 设备发生事件（如低电量告警）, **When** 设备发布 events 消息, **Then** dji-adapter 能够解析事件类型和数据，发布到 `iot.dji.device.event`
2. **Given** 平台需要确认事件, **When** 发送事件确认, **Then** dji-adapter 能够转换并发布到 `thing/product/{gateway_sn}/events_reply`

---

### User Story 6 - 航线任务管理 (Priority: P2)

作为平台操作员，我需要上传航线文件并执行航线任务，以便实现自动化飞行。

**Why this priority**: 航线任务是 DJI 机场的核心业务功能。

**Independent Test**: 可以通过上传航线文件并执行任务，验证完整流程。

**Acceptance Scenarios**:

1. **Given** 平台上传 WPML 航线文件, **When** 调用航线上传服务, **Then** dji-adapter 能够处理文件上传流程
2. **Given** 平台下发航线任务, **When** 调用 `flighttask_prepare` 和 `flighttask_execute`, **Then** 设备能够执行航线任务
3. **Given** 航线任务执行中, **When** 设备上报任务进度, **Then** 平台能够接收进度更新

---

### User Story 7 - 相机控制 (Priority: P2)

作为平台操作员，我需要控制设备相机（拍照、录像、变焦），以便获取现场图像。

**Why this priority**: 相机控制是无人机应用的常见需求。

**Independent Test**: 可以通过发送相机控制指令，验证相机响应。

**Acceptance Scenarios**:

1. **Given** 平台发送拍照指令, **When** dji-adapter 处理请求, **Then** 设备能够执行拍照并返回结果
2. **Given** 平台发送录像开始/停止指令, **When** dji-adapter 处理请求, **Then** 设备能够控制录像状态

---

### User Story 8 - 实时控制 DRC (Priority: P3)

作为平台操作员，我需要实时控制飞行器（虚拟摇杆），以便进行精细操控。

**Why this priority**: DRC 是高级功能，需要低延迟通信。

**Independent Test**: 可以通过发送 DRC 指令，验证实时控制响应。

**Acceptance Scenarios**:

1. **Given** 平台建立 DRC 连接, **When** 发送虚拟摇杆指令, **Then** 飞行器能够实时响应
2. **Given** DRC 连接中断, **When** 超过心跳超时, **Then** 飞行器能够安全悬停

---

### User Story 9 - 媒体文件管理 (Priority: P3)

作为平台运维人员，我需要管理设备上的媒体文件（照片、视频），以便下载和存储。

**Why this priority**: 媒体管理是数据采集的重要环节。

**Independent Test**: 可以通过获取文件列表并下载文件，验证媒体管理功能。

**Acceptance Scenarios**:

1. **Given** 设备上有媒体文件, **When** 平台请求文件列表, **Then** 能够获取文件元数据
2. **Given** 平台请求下载文件, **When** 设备上传文件, **Then** 平台能够接收并存储文件

---

### Edge Cases

- 消息格式不符合 DJI 协议规范时如何处理？（记录错误日志，发送到死信队列，不影响其他消息处理）
- 设备 SN 无法识别时如何处理？（**记录日志并丢弃**，安全优先，防止未授权设备接入）
- 服务调用响应超时如何处理？（**30 秒超时**，返回超时错误，支持配置重试策略）
- 大量设备同时上报数据时如何保证性能？（支持水平扩展，消息批处理）
- 协议版本不兼容时如何处理？（支持协议版本检测，向后兼容）
- DRC 心跳超时如何处理？（**3 秒超时**，触发飞行器安全悬停）

## Requirements

### Technology Stack Requirements

- **TS-001**: 必须使用 Go 1.22+ 开发
- **TS-002**: 必须复用 002 实现的协议适配器框架 (`pkg/adapter`)
- **TS-003**: 必须复用 001 实现的基础设施 (`pkg/rabbitmq`, `pkg/tracer`, `pkg/metrics`)
- **TS-004**: 必须遵循 Uber Go 语言编码规范
- **TS-005**: 必须使用 `golangci-lint` 进行代码检查

### Functional Requirements

#### 核心协议支持 (P1)

- **FR-001**: 系统必须支持 DJI OSD 消息解析，包括飞行器、机场、遥控器的遥测数据
- **FR-002**: 系统必须支持 DJI State 消息解析，处理属性变化通知
- **FR-003**: 系统必须支持 DJI Status 消息解析，处理设备上下线状态
- **FR-004**: 系统必须支持 DJI Services 消息的双向转换（调用和响应）
- **FR-005**: 系统必须支持 DJI Events 消息的双向转换（上报和确认）
- **FR-006**: 系统必须支持 DJI Requests 消息处理（设备主动请求）

#### 业务功能支持 (P2)

- **FR-007**: 系统必须支持航线任务管理相关服务（prepare, execute, pause, resume, cancel）
- **FR-008**: 系统必须支持相机控制相关服务（拍照、录像、变焦、云台控制）
- **FR-009**: 系统必须支持设备配置管理（读取、设置设备参数）

#### 高级功能支持 (P3)

- **FR-010**: 系统必须支持 DRC 实时控制协议
- **FR-011**: 系统必须支持媒体文件管理（文件列表、上传、下载）
- **FR-012**: 系统必须支持固件升级流程
- **FR-013**: 系统必须支持实时视频流管理

#### 协议适配要求

- **FR-014**: 所有 DJI 协议消息必须转换为 StandardMessage 格式
- **FR-015**: StandardMessage 必须包含完整的 protocol_meta 信息
- **FR-016**: 必须支持设备拓扑解析（机场-飞行器-负载的层级关系）
- **FR-017**: 必须支持 DJI 错误码到平台错误码的映射

#### 可观测性要求

- **FR-018**: 必须记录所有消息处理的结构化日志（包含 trace_id, device_sn, message_type）
- **FR-019**: 必须暴露 Prometheus 指标（消息数量、处理延迟、错误率）
- **FR-020**: 必须支持 W3C Trace Context 传播

### Test-First Development Requirements

- **TDD-001**: 所有功能开发必须遵循 TDD 原则
- **TDD-002**: 必须提供 DJI 协议消息的 Mock 数据用于测试
- **TDD-003**: 单元测试覆盖率必须 >= 80%
- **TDD-004**: 必须提供集成测试验证完整消息流程

### Key Entities

- **DJIMessage**: DJI 原始消息结构，包含 tid, bid, timestamp, gateway, data, method
- **OSDData**: 遥测数据结构，包含飞行器、机场、遥控器的状态数据
- **StateData**: 属性变化数据结构
- **ServiceRequest/Response**: 服务调用请求和响应结构
- **EventData**: 事件数据结构
- **DeviceTopology**: 设备拓扑结构（机场-飞行器-负载）
- **WaylineTask**: 航线任务结构
- **DRCCommand**: 实时控制指令结构

## Success Criteria

### Measurable Outcomes

- **SC-001**: 支持 DJI Cloud API 核心 Topic（osd, state, status, services, events, requests）的完整解析
- **SC-002**: 消息处理延迟 < 50ms (P95)
- **SC-003**: 消息处理成功率 > 99.9%
- **SC-004**: 单元测试覆盖率 >= 80%
- **SC-005**: 支持至少 1000 台设备同时在线的消息处理能力
- **SC-006**: 所有服务调用支持超时和重试机制

## Clarifications

### Session 2025-02-05

- Q: 本 Spec 与 002 的关系？ → A: 002 实现了协议适配器框架和 DJI 协议层适配示例，本 Spec 在此基础上完整实现所有 DJI 协议功能
- Q: 是否需要实现所有 DJI 协议功能？ → A: 按优先级分阶段实现，P1 为核心功能必须实现，P2/P3 可根据实际需求调整
- Q: WPML 航线文件如何处理？ → A: 航线文件解析使用 `pkg/adapter/dji/wpml` 模块，本 Spec 负责航线任务的服务调用流程

### Session 2025-02-05 (Clarify)

- Q: 服务调用超时时间应该设置为多少？ → A: **30 秒**，适合大多数服务调用场景，平衡响应速度和可靠性
- Q: 未知设备 SN 的消息如何处理？ → A: **记录日志并丢弃**，安全优先，防止未授权设备接入
- Q: DRC 实时控制心跳超时时间？ → A: **3 秒**，快速检测连接中断，及时触发安全悬停

### Configuration Constants

基于上述澄清，定义以下配置常量：

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `SERVICE_CALL_TIMEOUT` | 30s | 服务调用超时时间 |
| `DRC_HEARTBEAT_TIMEOUT` | 3s | DRC 心跳超时时间 |
| `UNKNOWN_DEVICE_POLICY` | discard | 未知设备消息处理策略 (discard/forward/dlq) |
