# Research: DJI Protocol Implementation

**Feature**: 003-dji-protocol-implementation
**Date**: 2025-02-05

## R1: Import 路径问题分析 ✅ 已解决

### 问题描述

`pkg/adapter/dji/protocol/` 目录下的 23 个文件使用了外部 import 路径：
```go
import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"
```

### 受影响文件

- `protocol/device/commands.go`
- `protocol/device/events.go`
- `protocol/camera/commands.go`
- `protocol/camera/ir_camera.go`
- `protocol/wayline/commands.go`
- `protocol/wayline/events.go`
- `protocol/wayline/requests.go`
- `protocol/drc/commands.go`
- `protocol/drc/events.go`
- `protocol/file/commands.go`
- `protocol/file/events.go`
- `protocol/file/log_events.go`
- `protocol/firmware/commands.go`
- `protocol/firmware/events.go`
- `protocol/live/commands.go`
- `protocol/psdk/commands.go`
- `protocol/psdk/events.go`
- `protocol/safety/commands.go`
- `protocol/config/device_config.go`
- `protocol/config/organization.go`
- `protocol/config/requests.go`
- `protocol/config/storage.go`
- `protocol/aircraft/hms_events.go`

### 决策

**Decision**: ✅ 用户已手动修复 import 路径

**Status**: 已完成 (2025-02-05)

---

## R2: OSD 数据结构集成分析

### 现有结构

**`protocol/aircraft/osd.go`** 定义了完整的 OSD 数据结构：
- `AircraftOSD` - 飞行器 OSD 主结构 (80+ 字段)
- `PayloadInfo` - 负载信息
- `BatteryInfo` / `BatteryDetail` - 电池信息
- `CameraInfo` - 相机信息
- `GimbalInfo` - 云台信息
- `PositionState` - 定位状态
- `ObstacleAvoidance` - 避障状态

**`protocol/aircraft/dock_osd.go`** 定义了机场 OSD：
- `DockOSD` - 机场 OSD 主结构

**`protocol/aircraft/rc_osd.go`** 定义了遥控器 OSD：
- `RCOSD` - 遥控器 OSD 主结构

### 当前 adapter 实现

`pkg/adapter/dji/parser.go` 只做基础 JSON 解析：
```go
func ParseMessage(payload []byte) (*Message, error) {
    var msg Message
    if err := json.Unmarshal(payload, &msg); err != nil {
        return nil, err
    }
    return &msg, nil
}
```

`Message.Data` 是 `json.RawMessage`，未解析具体结构。

### 决策

**Decision**: 创建 OSD 解析器，将 `json.RawMessage` 解析为具体结构

**设计方案**:
```go
// pkg/adapter/dji/integration/osd_parser.go

type OSDParser struct{}

func (p *OSDParser) ParseAircraftOSD(data json.RawMessage) (*aircraft.AircraftOSD, error)
func (p *OSDParser) ParseDockOSD(data json.RawMessage) (*aircraft.DockOSD, error)
func (p *OSDParser) ParseRCOSD(data json.RawMessage) (*aircraft.RCOSD, error)
```

**Rationale**:
- 保持现有 adapter 接口不变
- 在需要时按需解析具体结构
- 支持部分更新（所有字段都是指针）

**Alternatives Rejected**:
- 修改 Message.Data 类型：破坏现有接口
- 在 parser.go 中直接解析：职责不清晰

---

## R3: 服务调用路由机制分析

### 服务方法清单

基于 `protocol/` 目录分析，DJI 协议支持以下服务方法：

**设备控制 (device)**:
- `cover_open` / `cover_close` - 舱盖控制
- `drone_open` / `drone_close` - 飞行器电源
- `charge_open` / `charge_close` - 充电控制
- `device_reboot` - 设备重启
- `debug_mode_open` / `debug_mode_close` - 调试模式
- `battery_maintenance_switch` - 电池维护模式
- `air_conditioner_mode_switch` - 空调模式
- `alarm_state_switch` - 告警状态
- `sdr_workmode_switch` - SDR 工作模式

**航线任务 (wayline)**:
- `flighttask_create` - 创建任务
- `flighttask_prepare` - 准备任务
- `flighttask_execute` - 执行任务
- `flighttask_pause` - 暂停任务
- `flighttask_recovery` - 恢复任务
- `flighttask_undo` - 撤销任务
- `return_home` - 返航
- `return_home_cancel` - 取消返航

**相机控制 (camera)**:
- `camera_mode_switch` - 切换相机模式
- `camera_photo_take` - 拍照
- `camera_recording_start` / `camera_recording_stop` - 录像
- `camera_aim` - 相机瞄准
- `camera_focal_length_set` - 设置焦距
- `gimbal_reset` - 云台复位
- `ir_metering_*` - 红外测温

**实时控制 (drc)**:
- `drc_mode_enter` / `drc_mode_exit` - 进入/退出 DRC 模式
- `drone_control` - 飞行器控制
- `drone_emergency_stop` - 紧急停止
- `heart` - 心跳

**文件管理 (file)**:
- `file_upload_start` / `file_upload_finish` - 文件上传
- `file_upload_list` - 获取文件列表

**固件升级 (firmware)**:
- `ota_create` / `ota_progress` - OTA 升级

**实时视频 (live)**:
- `live_start_push` / `live_stop_push` - 开始/停止推流
- `live_set_quality` - 设置画质
- `live_lens_change` - 切换镜头

### 决策

**Decision**: 实现 ServiceRouter 基于 method 字段路由到具体处理器

**设计方案**:
```go
// pkg/adapter/dji/router/service_router.go

type ServiceHandler func(ctx context.Context, msg *Message) (*Response, error)

type ServiceRouter struct {
    handlers map[string]ServiceHandler
}

func (r *ServiceRouter) Register(method string, handler ServiceHandler)
func (r *ServiceRouter) Route(ctx context.Context, method string, msg *Message) (*Response, error)

// 预注册所有服务方法
func NewServiceRouter() *ServiceRouter {
    r := &ServiceRouter{handlers: make(map[string]ServiceHandler)}
    r.Register("cover_open", handleCoverOpen)
    r.Register("flighttask_prepare", handleFlightTaskPrepare)
    // ... 其他方法
    return r
}
```

**Rationale**:
- 清晰的路由机制
- 易于扩展新方法
- 支持方法级别的错误处理

**Alternatives Rejected**:
- switch-case 路由：不易扩展
- 反射调用：性能开销大

---

## R4: 事件类型与处理机制分析

### 事件类型清单

**设备事件 (device)**:
- `device_exit_homing_notify` - 退出归位通知
- `device_temp_ntfy_need_clear` - 温度通知需清除
- `file_upload_callback` - 文件上传回调

**航线事件 (wayline)**:
- `flighttask_progress` - 任务进度
- `flighttask_ready` - 任务就绪
- `return_home_info` - 返航信息

**HMS 事件 (aircraft)**:
- `hms_notify` - 健康管理系统通知

**DRC 事件 (drc)**:
- `joystick_invalid_notify` - 摇杆无效通知
- `drc_status_notify` - DRC 状态通知

**文件事件 (file)**:
- `highest_priority_upload_flighttask_media` - 高优先级上传
- `file_upload_progress` - 文件上传进度

**固件事件 (firmware)**:
- `ota_progress` - OTA 进度

### 决策

**Decision**: 实现 EventRouter 基于 method 字段路由到具体处理器

**设计方案**:
```go
// pkg/adapter/dji/router/event_router.go

type EventHandler func(ctx context.Context, event *Event) error

type EventRouter struct {
    handlers map[string]EventHandler
}

func (r *EventRouter) Register(method string, handler EventHandler)
func (r *EventRouter) Route(ctx context.Context, method string, event *Event) error
```

**Rationale**:
- 与 ServiceRouter 保持一致的设计
- 支持事件确认机制
- 易于扩展新事件类型

---

## R5: 消息处理器设计

### 决策

**Decision**: 为每种 Topic 类型创建独立的 Handler

**设计方案**:
```go
// pkg/adapter/dji/handler/

// OSD Handler - 处理 thing/product/{sn}/osd
type OSDHandler struct {
    parser *integration.OSDParser
}
func (h *OSDHandler) Handle(ctx context.Context, msg *Message, topic *TopicInfo) (*StandardMessage, error)

// State Handler - 处理 thing/product/{sn}/state
type StateHandler struct{}
func (h *StateHandler) Handle(ctx context.Context, msg *Message, topic *TopicInfo) (*StandardMessage, error)

// Status Handler - 处理 sys/product/{sn}/status
type StatusHandler struct{}
func (h *StatusHandler) Handle(ctx context.Context, msg *Message, topic *TopicInfo) (*StandardMessage, error)

// Service Handler - 处理 thing/product/{sn}/services 和 services_reply
type ServiceHandler struct {
    router *router.ServiceRouter
}
func (h *ServiceHandler) Handle(ctx context.Context, msg *Message, topic *TopicInfo) (*StandardMessage, error)

// Event Handler - 处理 thing/product/{sn}/events
type EventHandler struct {
    router *router.EventRouter
}
func (h *EventHandler) Handle(ctx context.Context, msg *Message, topic *TopicInfo) (*StandardMessage, error)

// Request Handler - 处理 thing/product/{sn}/requests
type RequestHandler struct{}
func (h *RequestHandler) Handle(ctx context.Context, msg *Message, topic *TopicInfo) (*StandardMessage, error)
```

**Rationale**:
- 职责分离，每个 Handler 专注一种消息类型
- 易于测试和维护
- 支持独立扩展

---

## 总结

| 研究项 | 决策 | 状态 |
|--------|------|------|
| R1: Import 路径 | 用户已手动修复 | ✅ 完成 |
| R2: OSD 解析 | 创建 OSDParser 按需解析 | 待实现 |
| R3: 服务路由 | 实现 ServiceRouter | 待实现 |
| R4: 事件路由 | 实现 EventRouter | 待实现 |
| R5: 消息处理器 | 为每种 Topic 创建 Handler | 待实现 |

### 配置常量 (基于 Clarify 结果)

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `SERVICE_CALL_TIMEOUT` | 30s | 服务调用超时时间 |
| `DRC_HEARTBEAT_TIMEOUT` | 3s | DRC 心跳超时时间 |
| `UNKNOWN_DEVICE_POLICY` | discard | 未知设备消息处理策略 |

### 实现顺序

1. ~~**Phase 0**: 修复 import 路径 (R1)~~ ✅ 已完成
2. **Phase 1**: 实现 Handler 框架 (R5)
3. **Phase 2**: 实现 OSD 解析 (R2)
4. **Phase 3**: 实现 ServiceRouter (R3)
5. **Phase 4**: 实现 EventRouter (R4)
