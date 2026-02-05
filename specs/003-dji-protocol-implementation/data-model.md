# Data Model: DJI Protocol Implementation

**Feature**: 003-dji-protocol-implementation
**Date**: 2025-02-05

## 1. 消息类型枚举

### Topic 类型

```go
// TopicType 定义 DJI MQTT Topic 类型
type TopicType string

const (
    TopicTypeOSD           TopicType = "osd"            // 遥测数据定时上报
    TopicTypeState         TopicType = "state"          // 属性变化上报
    TopicTypeEvents        TopicType = "events"         // 事件上报
    TopicTypeEventsReply   TopicType = "events_reply"   // 事件响应
    TopicTypeServices      TopicType = "services"       // 服务调用
    TopicTypeServicesReply TopicType = "services_reply" // 服务响应
    TopicTypeRequests      TopicType = "requests"       // 设备请求
    TopicTypeRequestsReply TopicType = "requests_reply" // 请求响应
    TopicTypeStatus        TopicType = "status"         // 设备状态
    TopicTypeStatusReply   TopicType = "status_reply"   // 状态响应
    TopicTypePropertySet   TopicType = "property/set"   // 属性设置
    TopicTypeDRCUp         TopicType = "drc/up"         // DRC 上行
    TopicTypeDRCDown       TopicType = "drc/down"       // DRC 下行
)
```

### 消息方向

```go
// Direction 定义消息方向
type Direction string

const (
    DirectionUplink   Direction = "uplink"   // 设备 → 平台
    DirectionDownlink Direction = "downlink" // 平台 → 设备
)
```

### 设备类型

```go
// DeviceType 定义设备类型
type DeviceType int

const (
    DeviceTypeUnknown  DeviceType = 0
    DeviceTypeAircraft DeviceType = 60  // 飞行器
    DeviceTypeDock     DeviceType = 165 // 机场
    DeviceTypeDock2    DeviceType = 167 // 机场2
    DeviceTypeRC       DeviceType = 56  // 遥控器
    DeviceTypeRCPlus   DeviceType = 144 // 遥控器 Plus
)
```

## 2. OSD 数据结构

### 飞行器 OSD (AircraftOSD)

| 字段 | 类型 | 说明 |
|------|------|------|
| `mode_code` | *int | 飞行器状态 (0-20) |
| `longitude` | *float64 | 经度 |
| `latitude` | *float64 | 纬度 |
| `height` | *float64 | 绝对高度 (m) |
| `elevation` | *float64 | 相对起飞点高度 (m) |
| `horizontal_speed` | *float64 | 水平速度 (m/s) |
| `vertical_speed` | *float64 | 垂直速度 (m/s) |
| `attitude_pitch` | *float64 | 俯仰角 (度) |
| `attitude_roll` | *float64 | 横滚角 (度) |
| `attitude_head` | *float64 | 航向角 (度) |
| `battery` | *BatteryInfo | 电池信息 |
| `cameras` | []CameraInfo | 相机信息 |
| `payloads` | []PayloadInfo | 负载信息 |

### 机场 OSD (DockOSD)

| 字段 | 类型 | 说明 |
|------|------|------|
| `mode_code` | *int | 机场状态 |
| `cover_state` | *int | 舱盖状态 (0=关闭, 1=打开, 2=半开) |
| `putter_state` | *int | 推杆状态 |
| `charge_state` | *int | 充电状态 |
| `drone_in_dock` | *int | 飞行器是否在舱内 |
| `network_state` | *NetworkState | 网络状态 |
| `storage` | *Storage | 存储信息 |
| `environment_temperature` | *float64 | 环境温度 |
| `environment_humidity` | *int | 环境湿度 |

### 遥控器 OSD (RCOSD)

| 字段 | 类型 | 说明 |
|------|------|------|
| `mode_code` | *int | 遥控器状态 |
| `capacity_percent` | *int | 电量百分比 |
| `longitude` | *float64 | 经度 |
| `latitude` | *float64 | 纬度 |

## 3. 服务命令注册表

### 设备控制命令

| Method | 模块 | 说明 | 数据类型 |
|--------|------|------|----------|
| `cover_open` | device | 打开舱盖 | nil |
| `cover_close` | device | 关闭舱盖 | nil |
| `cover_force_close` | device | 强制关闭舱盖 | nil |
| `drone_open` | device | 飞行器开机 | nil |
| `drone_close` | device | 飞行器关机 | nil |
| `charge_open` | device | 开始充电 | nil |
| `charge_close` | device | 停止充电 | nil |
| `device_reboot` | device | 设备重启 | nil |
| `device_format` | device | 格式化设备 | nil |
| `drone_format` | device | 格式化飞行器 | nil |
| `debug_mode_open` | device | 开启调试模式 | nil |
| `debug_mode_close` | device | 关闭调试模式 | nil |
| `battery_maintenance_switch` | device | 电池维护模式 | BatteryMaintenanceSwitchData |
| `air_conditioner_mode_switch` | device | 空调模式切换 | AirConditionerModeSwitchData |
| `alarm_state_switch` | device | 告警状态切换 | AlarmStateSwitchData |
| `sdr_workmode_switch` | device | SDR 工作模式 | SDRWorkmodeSwitchData |

### 航线任务命令

| Method | 模块 | 说明 | 数据类型 |
|--------|------|------|----------|
| `flighttask_create` | wayline | 创建任务 | CreateData |
| `flighttask_prepare` | wayline | 准备任务 | PrepareData |
| `flighttask_execute` | wayline | 执行任务 | ExecuteData |
| `flighttask_pause` | wayline | 暂停任务 | nil |
| `flighttask_recovery` | wayline | 恢复任务 | nil |
| `flighttask_undo` | wayline | 撤销任务 | UndoData |
| `return_home` | wayline | 返航 | nil |
| `return_home_cancel` | wayline | 取消返航 | nil |

### 相机控制命令

| Method | 模块 | 说明 | 数据类型 |
|--------|------|------|----------|
| `camera_mode_switch` | camera | 切换相机模式 | CameraModeSwitchData |
| `camera_photo_take` | camera | 拍照 | CameraPhotoTakeData |
| `camera_recording_start` | camera | 开始录像 | CameraRecordingData |
| `camera_recording_stop` | camera | 停止录像 | CameraRecordingData |
| `camera_aim` | camera | 相机瞄准 | CameraAimData |
| `camera_focal_length_set` | camera | 设置焦距 | CameraFocalLengthData |
| `gimbal_reset` | camera | 云台复位 | GimbalResetData |

### DRC 控制命令

| Method | 模块 | 说明 | 数据类型 |
|--------|------|------|----------|
| `drc_mode_enter` | drc | 进入 DRC 模式 | DRCModeEnterData |
| `drc_mode_exit` | drc | 退出 DRC 模式 | nil |
| `drone_control` | drc | 飞行器控制 | DroneControlData |
| `drone_emergency_stop` | drc | 紧急停止 | nil |
| `heart` | drc | 心跳 | HeartData |

### 文件管理命令

| Method | 模块 | 说明 | 数据类型 |
|--------|------|------|----------|
| `file_upload_start` | file | 开始上传 | FileUploadStartData |
| `file_upload_finish` | file | 完成上传 | FileUploadFinishData |
| `file_upload_list` | file | 获取文件列表 | FileUploadListData |

### 固件升级命令

| Method | 模块 | 说明 | 数据类型 |
|--------|------|------|----------|
| `ota_create` | firmware | 创建 OTA 任务 | OTACreateData |

### 实时视频命令

| Method | 模块 | 说明 | 数据类型 |
|--------|------|------|----------|
| `live_start_push` | live | 开始推流 | LiveStartPushData |
| `live_stop_push` | live | 停止推流 | LiveStopPushData |
| `live_set_quality` | live | 设置画质 | LiveSetQualityData |
| `live_lens_change` | live | 切换镜头 | LiveLensChangeData |

## 4. 事件类型注册表

### 设备事件

| Method | 模块 | 说明 | 需要回复 |
|--------|------|------|----------|
| `device_exit_homing_notify` | device | 退出归位通知 | 否 |
| `device_temp_ntfy_need_clear` | device | 温度通知需清除 | 否 |
| `file_upload_callback` | device | 文件上传回调 | 是 |

### 航线事件

| Method | 模块 | 说明 | 需要回复 |
|--------|------|------|----------|
| `flighttask_progress` | wayline | 任务进度 | 否 |
| `flighttask_ready` | wayline | 任务就绪 | 是 |
| `return_home_info` | wayline | 返航信息 | 否 |

### HMS 事件

| Method | 模块 | 说明 | 需要回复 |
|--------|------|------|----------|
| `hms` | aircraft | 健康管理系统通知 | 否 |

### DRC 事件

| Method | 模块 | 说明 | 需要回复 |
|--------|------|------|----------|
| `joystick_invalid_notify` | drc | 摇杆无效通知 | 否 |
| `drc_status_notify` | drc | DRC 状态通知 | 否 |

### 文件事件

| Method | 模块 | 说明 | 需要回复 |
|--------|------|------|----------|
| `highest_priority_upload_flighttask_media` | file | 高优先级上传 | 是 |
| `file_upload_progress` | file | 文件上传进度 | 否 |

### 固件事件

| Method | 模块 | 说明 | 需要回复 |
|--------|------|------|----------|
| `ota_progress` | firmware | OTA 进度 | 否 |

## 5. 错误码映射

### DJI 错误码 → 平台错误码

| DJI Code | 平台 Code | 说明 |
|----------|-----------|------|
| 0 | 0 | 成功 |
| 314000 | 1001 | 参数错误 |
| 314001 | 1002 | 设备离线 |
| 314002 | 1003 | 设备忙 |
| 314003 | 1004 | 任务冲突 |
| 316001 | 2001 | 航线文件错误 |
| 316002 | 2002 | 航线任务不存在 |
| 317001 | 3001 | 相机错误 |
| 319001 | 4001 | DRC 连接失败 |

## 6. 实体关系图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              DJI Protocol Entities                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────┐         ┌─────────────┐         ┌─────────────┐           │
│  │   Message   │────────►│  TopicInfo  │────────►│  TopicType  │           │
│  │             │         │             │         │             │           │
│  │ tid         │         │ type        │         │ osd         │           │
│  │ bid         │         │ device_sn   │         │ state       │           │
│  │ timestamp   │         │ gateway_sn  │         │ events      │           │
│  │ method      │         │ direction   │         │ services    │           │
│  │ data        │         │ raw         │         │ status      │           │
│  └──────┬──────┘         └─────────────┘         └─────────────┘           │
│         │                                                                   │
│         │ data (json.RawMessage)                                           │
│         │                                                                   │
│         ▼                                                                   │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         Parsed Data Types                            │   │
│  ├─────────────────────────────────────────────────────────────────────┤   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ AircraftOSD │  │   DockOSD   │  │    RCOSD    │  │  EventData  │ │   │
│  │  │             │  │             │  │             │  │             │ │   │
│  │  │ mode_code   │  │ mode_code   │  │ mode_code   │  │ method      │ │   │
│  │  │ longitude   │  │ cover_state │  │ capacity    │  │ data        │ │   │
│  │  │ latitude    │  │ charge_state│  │ longitude   │  │ need_reply  │ │   │
│  │  │ height      │  │ drone_in    │  │ latitude    │  │             │ │   │
│  │  │ battery     │  │ storage     │  │             │  │             │ │   │
│  │  │ cameras     │  │ network     │  │             │  │             │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │   │
│  │  │ServiceReq   │  │ServiceReply │  │  StatusData │                  │   │
│  │  │             │  │             │  │             │                  │   │
│  │  │ method      │  │ result      │  │ online      │                  │   │
│  │  │ data        │  │ output      │  │ topology    │                  │   │
│  │  │             │  │             │  │             │                  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘                  │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 7. 状态机

### 飞行器状态 (mode_code)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Aircraft Mode Code State Machine                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│    ┌─────┐                                                                  │
│    │  0  │ Standby (待机)                                                   │
│    └──┬──┘                                                                  │
│       │ drone_open                                                          │
│       ▼                                                                     │
│    ┌─────┐                                                                  │
│    │  1  │ Preparing (准备中)                                               │
│    └──┬──┘                                                                  │
│       │ ready                                                               │
│       ▼                                                                     │
│    ┌─────┐     flighttask_execute    ┌─────┐                               │
│    │  2  │ ─────────────────────────►│  3  │ Manual Flight (手动飞行)       │
│    │Ready│                           │     │                                │
│    └──┬──┘                           └──┬──┘                                │
│       │                                 │                                   │
│       │ flighttask_execute              │ return_home                       │
│       ▼                                 ▼                                   │
│    ┌─────┐                           ┌─────┐                               │
│    │  4  │ Auto Flight (自动飞行)    │  5  │ Return Home (返航)             │
│    └──┬──┘                           └──┬──┘                                │
│       │                                 │                                   │
│       │ complete / pause                │ arrived                           │
│       ▼                                 ▼                                   │
│    ┌─────┐                           ┌─────┐                               │
│    │  6  │ Landing (降落中)          │  7  │ Landed (已降落)                │
│    └─────┘                           └─────┘                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 机场状态 (mode_code)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Dock Mode Code State Machine                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│    ┌─────┐                                                                  │
│    │  0  │ Idle (空闲)                                                      │
│    └──┬──┘                                                                  │
│       │ flighttask_prepare                                                  │
│       ▼                                                                     │
│    ┌─────┐                                                                  │
│    │  1  │ Preparing (准备中) - 开舱盖、推出飞行器                           │
│    └──┬──┘                                                                  │
│       │ ready                                                               │
│       ▼                                                                     │
│    ┌─────┐                                                                  │
│    │  2  │ Ready (就绪) - 等待起飞                                          │
│    └──┬──┘                                                                  │
│       │ flighttask_execute                                                  │
│       ▼                                                                     │
│    ┌─────┐                                                                  │
│    │  3  │ Working (工作中) - 飞行器执行任务                                 │
│    └──┬──┘                                                                  │
│       │ return_home / complete                                              │
│       ▼                                                                     │
│    ┌─────┐                                                                  │
│    │  4  │ Recovering (回收中) - 飞行器降落、收回、关舱盖                    │
│    └──┬──┘                                                                  │
│       │ complete                                                            │
│       ▼                                                                     │
│    ┌─────┐                                                                  │
│    │  0  │ Idle (空闲)                                                      │
│    └─────┘                                                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```
