# Data Model: Core Services Implementation

**Feature**: 004-core-services-implementation
**Date**: 2025-02-05

## Entity Relationship Diagram

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│ DeviceCredential│       │     Device      │       │  DeviceTopology │
├─────────────────┤       ├─────────────────┤       ├─────────────────┤
│ id              │       │ id              │       │ id              │
│ device_sn (UK)  │──────►│ sn (UK)         │◄──────│ device_id (FK)  │
│ username        │       │ name            │       │ parent_id (FK)  │
│ password_hash   │       │ type            │       │ relation_type   │
│ enabled         │       │ vendor          │       │ created_at      │
│ created_at      │       │ model           │       └─────────────────┘
│ updated_at      │       │ status          │
└─────────────────┘       │ online          │       ┌─────────────────┐
                          │ last_online_at  │       │  ServiceCall    │
                          │ created_at      │       ├─────────────────┤
                          │ updated_at      │       │ id              │
                          └────────┬────────┘       │ tid (UK)        │
                                   │                │ bid             │
                                   │                │ device_sn       │
                                   ▼                │ method          │
                          ┌─────────────────┐       │ params (JSON)   │
                          │  ThingModel     │       │ status          │
                          ├─────────────────┤       │ result (JSON)   │
                          │ id              │       │ error_code      │
                          │ device_id (FK)  │       │ error_msg       │
                          │ version         │       │ created_at      │
                          │ properties      │       │ completed_at    │
                          │ services        │       └─────────────────┘
                          │ events          │
                          │ created_at      │       ┌─────────────────┐
                          │ updated_at      │       │  WSConnection   │
                          └─────────────────┘       ├─────────────────┤
                                                    │ id              │
                                                    │ client_id (UK)  │
                                                    │ user_id         │
                                                    │ subscriptions   │
                                                    │ connected_at    │
                                                    │ last_ping_at    │
                                                    └─────────────────┘
```

## PostgreSQL Entities

### DeviceCredential

设备认证凭证，用于 MQTT 连接认证。

```go
// internal/gateway/model/credential.go
type DeviceCredential struct {
    ID           uint      `gorm:"primaryKey"`
    DeviceSN     string    `gorm:"uniqueIndex;size:64;not null"`
    Username     string    `gorm:"size:64;not null"`
    PasswordHash string    `gorm:"size:256;not null"`
    Enabled      bool      `gorm:"default:true"`
    CreatedAt    time.Time `gorm:"autoCreateTime"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (DeviceCredential) TableName() string {
    return "device_credentials"
}
```

### Device

设备基本信息。

```go
// pkg/models/device.go (已存在，需扩展)
type Device struct {
    ID           uint      `gorm:"primaryKey"`
    SN           string    `gorm:"uniqueIndex;size:64;not null"`
    Name         string    `gorm:"size:128"`
    Type         string    `gorm:"size:32;not null"` // gateway, aircraft, dock, rc
    Vendor       string    `gorm:"size:32;not null"` // dji, etc.
    Model        string    `gorm:"size:64"`
    Status       string    `gorm:"size:32;default:'inactive'"` // active, inactive, maintenance
    Online       bool      `gorm:"default:false"`
    LastOnlineAt *time.Time
    CreatedAt    time.Time `gorm:"autoCreateTime"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
```

### DeviceTopology

设备拓扑关系（网关-子设备）。

```go
// pkg/models/topology.go
type DeviceTopology struct {
    ID           uint      `gorm:"primaryKey"`
    DeviceID     uint      `gorm:"not null;index"`
    Device       Device    `gorm:"foreignKey:DeviceID"`
    ParentID     *uint     `gorm:"index"`
    Parent       *Device   `gorm:"foreignKey:ParentID"`
    RelationType string    `gorm:"size:32;not null"` // gateway_aircraft, dock_aircraft
    CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (DeviceTopology) TableName() string {
    return "device_topologies"
}
```

### ServiceCall

服务调用记录，用于追踪和重试。

```go
// internal/downlink/model/service_call.go
type ServiceCall struct {
    ID          uint            `gorm:"primaryKey"`
    TID         string          `gorm:"uniqueIndex;size:64;not null"` // Transaction ID
    BID         string          `gorm:"size:64;not null"`             // Business ID
    DeviceSN    string          `gorm:"index;size:64;not null"`
    Method      string          `gorm:"size:128;not null"`
    Params      datatypes.JSON  `gorm:"type:jsonb"`
    Status      string          `gorm:"size:32;default:'pending'"` // pending, sent, success, failed, timeout
    Result      datatypes.JSON  `gorm:"type:jsonb"`
    ErrorCode   *int
    ErrorMsg    *string         `gorm:"size:512"`
    RetryCount  int             `gorm:"default:0"`
    CreatedAt   time.Time       `gorm:"autoCreateTime"`
    CompletedAt *time.Time
}

func (ServiceCall) TableName() string {
    return "service_calls"
}
```

### WSConnection (内存结构)

WebSocket 连接信息，存储在内存中。

```go
// internal/ws/model/connection.go
type WSConnection struct {
    ID            string
    ClientID      string
    UserID        string
    Subscriptions []string  // 订阅的主题列表
    ConnectedAt   time.Time
    LastPingAt    time.Time
    Conn          *websocket.Conn
}
```

## InfluxDB Measurements

### dji_aircraft_osd

飞行器遥测数据。

```
Measurement: dji_aircraft_osd
Tags:
  - device_sn: string      # 飞行器 SN
  - gateway_sn: string     # 网关 SN
  - vendor: string         # 厂商 (dji)
Fields:
  - latitude: float        # 纬度
  - longitude: float       # 经度
  - altitude: float        # 高度 (m)
  - height: float          # 相对高度 (m)
  - speed_x: float         # X 轴速度 (m/s)
  - speed_y: float         # Y 轴速度 (m/s)
  - speed_z: float         # Z 轴速度 (m/s)
  - attitude_pitch: float  # 俯仰角 (deg)
  - attitude_roll: float   # 横滚角 (deg)
  - attitude_yaw: float    # 偏航角 (deg)
  - battery_percent: int   # 电池电量 (%)
  - flight_mode: int       # 飞行模式
  - gear: int              # 起落架状态
Timestamp: message timestamp (nanoseconds)
```

### dji_dock_osd

机场遥测数据。

```
Measurement: dji_dock_osd
Tags:
  - device_sn: string      # 机场 SN
  - vendor: string         # 厂商 (dji)
Fields:
  - network_state: int     # 网络状态
  - drone_in_dock: bool    # 飞行器是否在舱内
  - drone_charge_state: int # 充电状态
  - cover_state: int       # 舱盖状态
  - putter_state: int      # 推杆状态
  - supplement_light_state: int # 补光灯状态
  - temperature: float     # 温度 (°C)
  - humidity: float        # 湿度 (%)
  - rainfall: int          # 降雨量
  - wind_speed: float      # 风速 (m/s)
Timestamp: message timestamp (nanoseconds)
```

### dji_device_event

设备事件记录。

```
Measurement: dji_device_event
Tags:
  - device_sn: string      # 设备 SN
  - gateway_sn: string     # 网关 SN
  - vendor: string         # 厂商 (dji)
  - event_type: string     # 事件类型
Fields:
  - event_data: string     # 事件数据 (JSON)
  - need_reply: bool       # 是否需要回复
Timestamp: message timestamp (nanoseconds)
```

## RabbitMQ Message Formats

### Raw Uplink Message

从 iot-gateway 发送到 iot-uplink 的原始消息。

```go
// pkg/rabbitmq/message.go (扩展)
type RawUplinkMessage struct {
    Vendor    string          `json:"vendor"`     // dji
    Topic     string          `json:"topic"`      // MQTT topic
    Payload   json.RawMessage `json:"payload"`    // 原始 payload
    QoS       int             `json:"qos"`
    Timestamp int64           `json:"timestamp"`
    TraceID   string          `json:"trace_id"`
    SpanID    string          `json:"span_id"`
}
```

### Standard Message

标准化消息格式，用于服务间通信。

```go
// pkg/rabbitmq/message.go (已存在)
type StandardMessage struct {
    TID          string                 `json:"tid"`
    BID          string                 `json:"bid"`
    Timestamp    int64                  `json:"timestamp"`
    DeviceSN     string                 `json:"device_sn"`
    GatewaySN    string                 `json:"gateway_sn,omitempty"`
    Service      string                 `json:"service"`
    Action       string                 `json:"action"`
    Data         map[string]interface{} `json:"data"`
    ProtocolMeta ProtocolMeta           `json:"protocol_meta"`
}

type ProtocolMeta struct {
    Vendor        string `json:"vendor"`
    OriginalTopic string `json:"original_topic"`
    QoS           int    `json:"qos"`
    Method        string `json:"method,omitempty"`
}
```

### Service Call Request

服务调用请求消息。

```go
// pkg/rabbitmq/message.go (扩展)
type ServiceCallRequest struct {
    TID       string                 `json:"tid"`
    BID       string                 `json:"bid"`
    DeviceSN  string                 `json:"device_sn"`
    Method    string                 `json:"method"`
    Params    map[string]interface{} `json:"params"`
    Timeout   int                    `json:"timeout"` // seconds
    TraceID   string                 `json:"trace_id"`
    SpanID    string                 `json:"span_id"`
}
```

### Service Call Response

服务调用响应消息。

```go
// pkg/rabbitmq/message.go (扩展)
type ServiceCallResponse struct {
    TID       string                 `json:"tid"`
    BID       string                 `json:"bid"`
    DeviceSN  string                 `json:"device_sn"`
    Method    string                 `json:"method"`
    Result    int                    `json:"result"` // 0 = success
    Output    map[string]interface{} `json:"output,omitempty"`
    ErrorCode *int                   `json:"error_code,omitempty"`
    ErrorMsg  *string                `json:"error_msg,omitempty"`
    TraceID   string                 `json:"trace_id"`
    SpanID    string                 `json:"span_id"`
}
```

## WebSocket Message Formats

### Client Subscribe

客户端订阅消息。

```go
// internal/ws/message/subscribe.go
type SubscribeMessage struct {
    Type    string   `json:"type"`    // "subscribe"
    Topics  []string `json:"topics"`  // ["device.osd.device-001", "device.event.*"]
}
```

### Server Push

服务器推送消息。

```go
// internal/ws/message/push.go
type PushMessage struct {
    Type      string                 `json:"type"`      // "message"
    Topic     string                 `json:"topic"`     // "device.osd.device-001"
    DeviceSN  string                 `json:"device_sn"`
    Timestamp int64                  `json:"timestamp"`
    Data      map[string]interface{} `json:"data"`
}
```

## Validation Rules

### DeviceCredential

- `device_sn`: 必填，唯一，最大 64 字符
- `username`: 必填，最大 64 字符
- `password_hash`: 必填，bcrypt 哈希

### Device

- `sn`: 必填，唯一，最大 64 字符
- `type`: 必填，枚举值 (gateway, aircraft, dock, rc)
- `vendor`: 必填，枚举值 (dji)

### ServiceCall

- `tid`: 必填，唯一，UUID 格式
- `bid`: 必填，UUID 格式
- `device_sn`: 必填，必须存在于 devices 表
- `method`: 必填，最大 128 字符
- `status`: 枚举值 (pending, sent, success, failed, timeout)

## State Transitions

### ServiceCall Status

```
pending ──► sent ──► success
              │
              └──► failed
              │
              └──► timeout
```

### Device Online Status

```
offline ──► online (收到消息)
   ▲           │
   │           │
   └───────────┘ (心跳超时)
```
