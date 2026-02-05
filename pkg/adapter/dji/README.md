# DJI Protocol Adapter

DJI Cloud API 协议适配器，实现 DJI 设备与 UMOS 平台的消息转换。

## 架构

```
pkg/adapter/dji/
├── adapter.go          # 主适配器实现
├── topic.go            # MQTT Topic 解析
├── parser.go           # 消息解析
├── converter.go        # 消息转换
├── types.go            # 基础类型定义
├── errors.go           # 错误定义
├── handler/            # 消息处理器
│   ├── handler.go      # Handler 接口
│   ├── registry.go     # Handler 注册表
│   ├── osd_handler.go  # OSD 消息处理
│   ├── state_handler.go    # State 消息处理
│   ├── status_handler.go   # Status 消息处理
│   ├── service_handler.go  # Service 消息处理
│   ├── event_handler.go    # Event 消息处理
│   ├── request_handler.go  # Request 消息处理
│   └── drc_handler.go      # DRC 消息处理
├── router/             # 服务/事件路由
│   ├── service_router.go   # 服务调用路由
│   ├── event_router.go     # 事件路由
│   ├── device_commands.go  # 设备控制命令
│   ├── camera_commands.go  # 相机控制命令
│   ├── wayline_commands.go # 航线任务命令
│   ├── drc_commands.go     # DRC 控制命令
│   └── ...
├── protocol/           # DJI 协议数据结构
│   ├── aircraft/       # 飞行器相关
│   ├── camera/         # 相机相关
│   ├── device/         # 设备相关
│   ├── wayline/        # 航线相关
│   └── ...
├── integration/        # 协议集成
│   └── osd_parser.go   # OSD 数据解析
├── init/               # 初始化
│   └── init.go         # Handler 初始化
├── config/             # 配置
│   └── config.go       # 配置常量
└── wpml/               # WPML 航线文件处理
```

## 使用方法

### 初始化适配器

```go
import (
    dji "github.com/utmos/utmos/pkg/adapter/dji"
    djiinit "github.com/utmos/utmos/pkg/adapter/dji/init"
)

// 创建并初始化适配器
adapter := djiinit.NewInitializedAdapter()

// 或者手动初始化
adapter := dji.NewAdapter()
if err := djiinit.InitializeAdapter(adapter); err != nil {
    log.Fatal(err)
}
```

### 处理消息

```go
// 解析 MQTT Topic
topic, err := adapter.ParseTopic("thing/product/gateway-001/osd")
if err != nil {
    return err
}

// 解析消息
msg, err := adapter.ParseMessage(payload)
if err != nil {
    return err
}

// 处理消息
standardMsg, err := adapter.HandleMessage(ctx, msg, topic)
if err != nil {
    return err
}
```

## Handler 类型

| Handler | Topic | 说明 |
|---------|-------|------|
| OSDHandler | `thing/product/{sn}/osd` | 遥测数据 |
| StateHandler | `thing/product/{sn}/state` | 属性变化 |
| StatusHandler | `sys/product/{sn}/status` | 设备状态 |
| ServiceHandler | `thing/product/{sn}/services` | 服务调用 |
| EventHandler | `thing/product/{sn}/events` | 事件上报 |
| RequestHandler | `thing/product/{sn}/requests` | 设备请求 |
| DRCHandler | `thing/product/{sn}/drc/up` | 实时控制 |

## 配置常量

| 常量 | 值 | 说明 |
|------|-----|------|
| `ServiceCallTimeout` | 30s | 服务调用超时 |
| `DRCHeartbeatTimeout` | 3s | DRC 心跳超时 |
| `UnknownDevicePolicy` | discard | 未知设备处理策略 |

## 可观测性

DJI adapter 集成了完整的可观测性支持：

### Prometheus 指标

```go
import (
    "github.com/utmos/utmos/pkg/adapter/dji/observability"
    "github.com/utmos/utmos/pkg/metrics"
)

// 创建 metrics collector
collector := metrics.NewCollector("iot")
djiMetrics := observability.NewMetrics(collector)

// 记录消息
djiMetrics.RecordMessageReceived("osd", "success")
djiMetrics.RecordProcessingDuration("osd", 0.007)
```

**暴露的指标**:
- `iot_dji_messages_received_total` - 接收消息总数
- `iot_dji_messages_sent_total` - 发送消息总数
- `iot_dji_messages_errors_total` - 错误总数
- `iot_dji_message_processing_duration_seconds` - 处理延迟直方图
- `iot_dji_active_devices` - 活跃设备数

### 结构化日志

```go
import "github.com/utmos/utmos/pkg/adapter/dji/observability"

logger := observability.DefaultLogger()

// 带 trace context 的日志
logger.WithMessage(ctx, "device-001", "osd", "property.report").Info("processing message")

// 输出示例 (JSON):
// {"level":"info","msg":"processing message","trace_id":"abc123","span_id":"def456",
//  "vendor":"dji","device_sn":"device-001","message_type":"osd","method":"property.report"}
```

### 分布式追踪

```go
import "github.com/utmos/utmos/pkg/adapter/dji/observability"

tracer := observability.NewTracer()

// 创建消息处理 span
ctx, span := tracer.StartMessageSpan(ctx, "osd", "property.report", "device-001")
defer span.End()

// 记录错误
if err != nil {
    tracer.RecordError(span, err)
}
```

### Handler Observer (推荐)

```go
import "github.com/utmos/utmos/pkg/adapter/dji/observability"

observer := observability.NewHandlerObserver(metrics, tracer, logger)

// 开始观察
result := observer.StartObserve(ctx, "osd", "property.report", "device-001")

// 处理消息...

// 结束观察 (自动记录 metrics, logs, traces)
result.End()
// 或者出错时
result.EndWithError(err, "parse_error")
```

## 性能

基准测试结果（Apple M1）：

| Handler | 延迟 | 内存分配 |
|---------|------|----------|
| OSD | ~7μs | 5.3KB |
| State | ~1.6μs | 2KB |
| Status | ~1.6μs | 2.2KB |

所有处理延迟远低于 50ms 目标。

## 测试

```bash
# 运行单元测试
go test ./pkg/adapter/dji/...

# 运行基准测试
go test -bench=. ./pkg/adapter/dji/handler/...

# 查看测试覆盖率
go test -cover ./pkg/adapter/dji/...
```
