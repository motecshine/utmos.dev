# Quickstart: DJI Protocol Implementation

本文档介绍如何测试和使用 DJI 协议实现。

## 1. 环境准备

### 启动基础设施

```bash
# 启动 RabbitMQ, PostgreSQL, InfluxDB
docker-compose up -d rabbitmq postgres influxdb
```

### 启动 dji-adapter 服务

```bash
# 开发模式
make run-dji-adapter

# 或直接运行
go run cmd/dji-adapter/main.go
```

## 2. 测试 OSD 消息处理

### 模拟飞行器 OSD 消息

```bash
# 发送模拟 OSD 消息到 RabbitMQ
go run tests/mocks/send_osd.go
```

**预期结果**:
- dji-adapter 接收 `iot.raw.dji.uplink` 队列消息
- 解析 OSD 数据 (位置、高度、电量等)
- 转换为 StandardMessage
- 发布到 `iot.dji.aircraft.osd` 队列

### 验证消息转换

```bash
# 监听标准消息队列
go run tests/mocks/listen_standard.go
```

**预期输出**:
```json
{
  "tid": "xxx",
  "bid": "xxx",
  "timestamp": 1706000000000,
  "device_sn": "AIRCRAFT-SN-001",
  "service": "dji-adapter",
  "action": "property.report",
  "data": {
    "longitude": 116.397128,
    "latitude": 39.916527,
    "height": 100.5,
    "battery": {
      "capacity_percent": 85
    }
  },
  "protocol_meta": {
    "vendor": "dji",
    "original_topic": "thing/product/DOCK-SN-001/osd"
  }
}
```

## 3. 测试服务调用

### 发送服务调用请求

```bash
# 发送航线任务准备请求
go run tests/mocks/send_service_call.go flighttask_prepare
```

**预期结果**:
- dji-adapter 接收 `iot.dji.service.call` 队列消息
- 转换为 DJI 协议格式
- 发布到 `iot.raw.dji.downlink` 队列

### 模拟设备响应

```bash
# 模拟设备返回服务响应
go run tests/mocks/send_service_reply.go
```

**预期结果**:
- dji-adapter 接收 `iot.raw.dji.uplink` 队列的 services_reply 消息
- 转换为 StandardMessage
- 发布到 `iot.dji.service.reply` 队列

## 4. 测试事件处理

### 模拟设备事件

```bash
# 发送航线任务进度事件
go run tests/mocks/send_event.go flighttask_progress
```

**预期结果**:
- dji-adapter 接收事件消息
- 解析事件类型和数据
- 发布到 `iot.dji.device.event` 队列

## 5. 添加新的服务命令

### 步骤 1: 定义命令数据结构

在 `pkg/adapter/dji/protocol/{module}/commands.go` 中添加:

```go
// NewCustomCommand represents a custom command
type NewCustomCommand struct {
    common.Header
    MethodName string         `json:"method"`
    DataValue  NewCustomData  `json:"data"`
}

func NewNewCustomCommand(data NewCustomData) *NewCustomCommand {
    return &NewCustomCommand{
        Header:     common.NewHeader(),
        MethodName: "new_custom_method",
        DataValue:  data,
    }
}
```

### 步骤 2: 注册到 ServiceRouter

在 `pkg/adapter/dji/router/service_router.go` 中添加:

```go
func init() {
    DefaultRouter.Register("new_custom_method", handleNewCustomMethod)
}

func handleNewCustomMethod(ctx context.Context, msg *Message) (*Response, error) {
    // 处理逻辑
    return &Response{Result: 0}, nil
}
```

### 步骤 3: 添加测试

在 `pkg/adapter/dji/router/service_router_test.go` 中添加:

```go
func TestNewCustomMethod(t *testing.T) {
    router := NewServiceRouter()

    msg := &Message{
        Method: "new_custom_method",
        Data:   json.RawMessage(`{"key": "value"}`),
    }

    resp, err := router.Route(context.Background(), msg.Method, msg)
    require.NoError(t, err)
    assert.Equal(t, 0, resp.Result)
}
```

## 6. 运行测试

```bash
# 运行所有 DJI 适配器测试
go test -v ./pkg/adapter/dji/...

# 运行特定包的测试
go test -v ./pkg/adapter/dji/handler/...
go test -v ./pkg/adapter/dji/router/...

# 运行基准测试
go test -bench=. ./pkg/adapter/dji/handler/...

# 检查测试覆盖率
go test -cover ./pkg/adapter/dji/...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./pkg/adapter/dji/...
go tool cover -html=coverage.out -o coverage.html
```

### 测试示例

#### 测试 OSD Handler

```go
func TestOSDHandler(t *testing.T) {
    handler := handler.NewOSDHandler()

    osdData := json.RawMessage(`{
        "mode_code": 0,
        "longitude": 113.943,
        "latitude": 22.577,
        "height": 100.5
    }`)

    msg := &dji.Message{
        TID:       "test-tid",
        BID:       "test-bid",
        Timestamp: time.Now().UnixMilli(),
        Data:      osdData,
    }

    topic := &dji.TopicInfo{
        Type:      dji.TopicTypeOSD,
        DeviceSN:  "test-device",
        GatewaySN: "test-gateway",
    }

    result, err := handler.Handle(context.Background(), msg, topic)
    require.NoError(t, err)
    assert.Equal(t, "test-device", result.DeviceSN)
}
```

#### 测试 Service Router

```go
func TestServiceRouter(t *testing.T) {
    r := router.NewServiceRouter()
    router.RegisterDeviceCommands(r)

    req := &router.ServiceRequest{
        Method: "cover_open",
        Data:   nil,
    }

    resp, err := r.RouteService(context.Background(), req)
    require.NoError(t, err)
    assert.Equal(t, 0, resp.Result)
}
```

## 7. 常见问题

### Q: 消息未被处理？

检查:
1. RabbitMQ 连接是否正常
2. 队列绑定是否正确
3. 消息格式是否符合 DJI 协议

### Q: 服务调用超时？

检查:
1. 设备是否在线
2. 网络连接是否正常
3. 超时配置是否合理

### Q: 如何调试消息流？

1. 启用 DEBUG 日志级别
2. 使用 RabbitMQ 管理界面查看队列
3. 使用 Tempo 查看分布式追踪
