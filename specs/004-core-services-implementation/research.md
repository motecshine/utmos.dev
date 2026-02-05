# Research: Core Services Implementation

**Feature**: 004-core-services-implementation
**Date**: 2025-02-05

## Research Topics

### 1. MQTT Client Library Selection

**Decision**: paho.mqtt.golang

**Rationale**:
- Eclipse 官方维护的 Go MQTT 客户端
- 支持 MQTT 3.1.1 和 5.0
- 成熟稳定，社区活跃
- 支持自动重连、消息持久化

**Alternatives Considered**:
- `emqx/mqtt-go`: 功能较新，但社区较小
- `surgemq/surgemq`: 已停止维护

**Usage Pattern**:
```go
import mqtt "github.com/eclipse/paho.mqtt.golang"

opts := mqtt.NewClientOptions()
opts.AddBroker("tcp://vernemq:1883")
opts.SetClientID("iot-gateway")
opts.SetUsername("gateway")
opts.SetPassword("secret")
opts.SetAutoReconnect(true)
opts.SetOnConnectHandler(onConnect)
opts.SetConnectionLostHandler(onConnectionLost)

client := mqtt.NewClient(opts)
if token := client.Connect(); token.Wait() && token.Error() != nil {
    return token.Error()
}
```

### 2. WebSocket Library Selection

**Decision**: gorilla/websocket

**Rationale**:
- Go 生态最成熟的 WebSocket 库
- 高性能，低内存占用
- 支持压缩、子协议
- 与 Gin 框架集成良好

**Alternatives Considered**:
- `nhooyr/websocket`: 更现代的 API，但社区较小
- `gobwas/ws`: 更底层，需要更多手动管理

**Usage Pattern**:
```go
import "github.com/gorilla/websocket"

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // 生产环境需要验证
    },
}

func wsHandler(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    defer conn.Close()
    // 处理连接
}
```

### 3. InfluxDB Client Selection

**Decision**: influxdb-client-go/v2

**Rationale**:
- InfluxDB 官方 Go 客户端
- 支持 InfluxDB 2.x API
- 支持批量写入、异步写入
- 支持 Flux 查询语言

**Usage Pattern**:
```go
import influxdb2 "github.com/influxdata/influxdb-client-go/v2"

client := influxdb2.NewClient("http://influxdb:8086", "token")
writeAPI := client.WriteAPIBlocking("org", "bucket")

p := influxdb2.NewPoint("telemetry",
    map[string]string{"device_sn": "device-001"},
    map[string]interface{}{"temperature": 25.5},
    time.Now())

writeAPI.WritePoint(context.Background(), p)
```

### 4. Message Routing Pattern

**Decision**: Topic-based routing with RabbitMQ

**Rationale**:
- 使用 RabbitMQ topic exchange 实现灵活路由
- 路由键格式: `iot.{vendor}.{service}.{action}`
- 支持通配符订阅: `iot.dji.#`, `iot.*.uplink.*`

**Routing Table**:

| Source | Routing Key | Consumer |
|--------|-------------|----------|
| iot-gateway | `iot.raw.dji.uplink` | iot-uplink |
| iot-uplink | `iot.dji.aircraft.osd` | iot-ws |
| iot-uplink | `iot.dji.device.state` | iot-ws |
| iot-uplink | `iot.dji.device.event` | iot-ws, iot-api |
| iot-api | `iot.dji.service.call` | iot-downlink |
| iot-downlink | `iot.raw.dji.downlink` | iot-gateway |

### 5. Device Authentication Strategy

**Decision**: Username/Password + Device Registry

**Rationale**:
- 简单可靠，易于实现
- 设备凭证存储在 PostgreSQL
- 支持后续扩展到证书认证

**Flow**:
1. 设备连接 VerneMQ，提供用户名/密码
2. VerneMQ 调用 iot-gateway 的认证 webhook
3. iot-gateway 查询 PostgreSQL 验证凭证
4. 返回认证结果

**Database Schema**:
```sql
CREATE TABLE device_credentials (
    id SERIAL PRIMARY KEY,
    device_sn VARCHAR(64) UNIQUE NOT NULL,
    username VARCHAR(64) NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 6. WebSocket Connection Management

**Decision**: Hub pattern with subscription manager

**Rationale**:
- 集中管理所有 WebSocket 连接
- 支持按设备/主题订阅
- 高效的消息广播

**Architecture**:
```
Hub
├── clients map[*Client]bool     # 所有连接
├── subscriptions map[string][]*Client  # 主题订阅
├── register chan *Client        # 注册通道
├── unregister chan *Client      # 注销通道
└── broadcast chan Message       # 广播通道
```

### 7. Retry Mechanism

**Decision**: Exponential backoff with dead letter queue

**Rationale**:
- 指数退避避免雪崩
- 死信队列保存失败消息
- 支持手动重试

**Configuration**:
```go
type RetryConfig struct {
    MaxRetries     int           // 最大重试次数: 3
    InitialDelay   time.Duration // 初始延迟: 1s
    MaxDelay       time.Duration // 最大延迟: 30s
    Multiplier     float64       // 退避系数: 2.0
    DeadLetterQueue string       // 死信队列名
}
```

### 8. InfluxDB Data Model

**Decision**: Measurement per message type

**Rationale**:
- 每种消息类型一个 measurement
- 设备 SN 作为 tag
- 遥测数据作为 field

**Schema**:
```
Measurement: dji_aircraft_osd
Tags:
  - device_sn
  - gateway_sn
Fields:
  - latitude (float)
  - longitude (float)
  - altitude (float)
  - speed (float)
  - battery (int)
  - ...
Timestamp: message timestamp
```

## Conclusions

1. **MQTT**: 使用 paho.mqtt.golang，支持自动重连
2. **WebSocket**: 使用 gorilla/websocket + Hub 模式
3. **InfluxDB**: 使用官方客户端，批量写入
4. **Routing**: RabbitMQ topic exchange，标准路由键格式
5. **Auth**: 用户名/密码 + PostgreSQL 设备注册表
6. **Retry**: 指数退避 + 死信队列
