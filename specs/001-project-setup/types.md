# Type Definitions: Project Setup

本文档定义了 001-project-setup 功能所需的核心类型和接口，防止实现时的猜测。

## 1. 配置类型

### 1.1 主配置结构

```go
// internal/shared/config/config.go

// Config 应用主配置
type Config struct {
    Server    ServerConfig    `yaml:"server"`
    Database  DatabaseConfig  `yaml:"database"`
    RabbitMQ  RabbitMQConfig  `yaml:"rabbitmq"`
    Tracer    TracerConfig    `yaml:"tracer"`
    Metrics   MetricsConfig   `yaml:"metrics"`
    Logger    LoggerConfig    `yaml:"logger"`
}

// ServerConfig HTTP 服务器配置
type ServerConfig struct {
    Host         string        `yaml:"host" default:"0.0.0.0"`
    Port         int           `yaml:"port" default:"8080"`
    ReadTimeout  time.Duration `yaml:"read_timeout" default:"30s"`
    WriteTimeout time.Duration `yaml:"write_timeout" default:"30s"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
    Postgres PostgresConfig `yaml:"postgres"`
    InfluxDB InfluxDBConfig `yaml:"influxdb"`
}

// PostgresConfig PostgreSQL 配置
type PostgresConfig struct {
    Host            string        `yaml:"host" default:"localhost"`
    Port            int           `yaml:"port" default:"5432"`
    User            string        `yaml:"user"`
    Password        string        `yaml:"password"`
    DBName          string        `yaml:"dbname"`
    SSLMode         string        `yaml:"sslmode" default:"disable"`
    MaxIdleConns    int           `yaml:"max_idle_conns" default:"10"`
    MaxOpenConns    int           `yaml:"max_open_conns" default:"100"`
    ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" default:"1h"`
}

// InfluxDBConfig InfluxDB 配置
type InfluxDBConfig struct {
    URL    string `yaml:"url" default:"http://localhost:8086"`
    Token  string `yaml:"token"`
    Org    string `yaml:"org"`
    Bucket string `yaml:"bucket"`
}
```

### 1.2 RabbitMQ 配置

```go
// RabbitMQConfig RabbitMQ 配置
type RabbitMQConfig struct {
    URL           string        `yaml:"url" default:"amqp://guest:guest@localhost:5672/"`
    ExchangeName  string        `yaml:"exchange_name" default:"iot"`
    ExchangeType  string        `yaml:"exchange_type" default:"topic"`
    PrefetchCount int           `yaml:"prefetch_count" default:"10"`
    Retry         RetryConfig   `yaml:"retry"`
}

// RetryConfig 重试配置
type RetryConfig struct {
    MaxRetries   int           `yaml:"max_retries" default:"10"`
    InitialDelay time.Duration `yaml:"initial_delay" default:"1s"`
    MaxDelay     time.Duration `yaml:"max_delay" default:"30s"`
    Multiplier   float64       `yaml:"multiplier" default:"2.0"`
}
```

### 1.3 Tracer 配置

```go
// TracerConfig 分布式追踪配置
type TracerConfig struct {
    Enabled      bool    `yaml:"enabled" default:"true"`
    Endpoint     string  `yaml:"endpoint" default:"http://localhost:4318/v1/traces"`
    ServiceName  string  `yaml:"service_name"`
    SamplingRate float64 `yaml:"sampling_rate" default:"1.0"` // dev: 1.0, prod: 0.1
    BatchTimeout time.Duration `yaml:"batch_timeout" default:"5s"`
    MaxQueueSize int     `yaml:"max_queue_size" default:"2048"`
}
```

### 1.4 Metrics 配置

```go
// MetricsConfig Prometheus metrics 配置
type MetricsConfig struct {
    Enabled   bool   `yaml:"enabled" default:"true"`
    Path      string `yaml:"path" default:"/metrics"`
    Port      int    `yaml:"port" default:"9090"` // 独立端口或与服务共用
    Namespace string `yaml:"namespace" default:"iot"`
}
```

### 1.5 Logger 配置

```go
// LoggerConfig 日志配置
type LoggerConfig struct {
    Level      string `yaml:"level" default:"info"` // debug, info, warn, error
    Format     string `yaml:"format" default:"json"` // json, text
    Output     string `yaml:"output" default:"stdout"` // stdout, file
    FilePath   string `yaml:"file_path,omitempty"`
}
```

---

## 2. 消息类型

### 2.1 标准消息格式

```go
// pkg/rabbitmq/message.go

// StandardMessage RabbitMQ 标准消息格式
type StandardMessage struct {
    TID       string          `json:"tid"`       // 事务 ID (UUID)
    BID       string          `json:"bid"`       // 业务 ID (UUID)
    Timestamp int64           `json:"timestamp"` // 毫秒级 Unix 时间戳
    Service   string          `json:"service"`   // 发送服务名
    Action    string          `json:"action"`    // 动作标识
    DeviceSN  string          `json:"device_sn"` // 设备序列号
    Data      json.RawMessage `json:"data"`      // 业务数据
}

// MessageHeader RabbitMQ 消息头
type MessageHeader struct {
    Traceparent string `json:"traceparent"` // W3C Trace Context
    Tracestate  string `json:"tracestate"`  // W3C Trace State
    MessageType string `json:"message_type"` // property, event, service
    Vendor      string `json:"vendor,omitempty"` // 厂商标识（可选）
}

// NewStandardMessage 创建标准消息
func NewStandardMessage(service, action, deviceSN string, data interface{}) (*StandardMessage, error)

// Validate 验证消息格式
func (m *StandardMessage) Validate() error
```

### 2.2 Routing Key 类型

```go
// pkg/rabbitmq/routing.go

// RoutingKey RabbitMQ routing key 结构
// 格式: iot.{vendor}.{service}.{action}
type RoutingKey struct {
    Vendor  string // 厂商标识: dji, generic, tuya 等
    Service string // 服务名: gateway, uplink, downlink, api, ws
    Action  string // 动作: property.report, service.call, event.notify 等
}

// String 返回 routing key 字符串
// 示例: iot.dji.gateway.property.report
func (r RoutingKey) String() string

// Parse 解析 routing key 字符串
func Parse(key string) (*RoutingKey, error)

// NewRoutingKey 创建 routing key
func NewRoutingKey(vendor, service, action string) RoutingKey

// 预定义 Action 常量
const (
    ActionPropertyReport  = "property.report"
    ActionPropertySet     = "property.set"
    ActionServiceCall     = "service.call"
    ActionServiceReply    = "service.reply"
    ActionEventNotify     = "event.notify"
    ActionDeviceOnline    = "device.online"
    ActionDeviceOffline   = "device.offline"
)

// 预定义 Vendor 常量
const (
    VendorDJI     = "dji"
    VendorGeneric = "generic"
    VendorTuya    = "tuya"
)
```

---

## 3. 接口定义

### 3.1 RabbitMQ 客户端接口

```go
// pkg/rabbitmq/client.go

// Client RabbitMQ 客户端接口
type Client interface {
    // Connect 连接 RabbitMQ（含指数退避重试）
    Connect(ctx context.Context) error

    // Close 关闭连接
    Close() error

    // IsConnected 检查连接状态
    IsConnected() bool

    // DeclareExchange 声明 Exchange（幂等）
    DeclareExchange(name, exchangeType string) error

    // DeclareQueue 声明 Queue（幂等）
    DeclareQueue(name string, durable bool) error

    // BindQueue 绑定 Queue 到 Exchange
    BindQueue(queueName, routingKey, exchangeName string) error
}

// Publisher 消息发布接口
type Publisher interface {
    // Publish 发布消息（自动注入 W3C Trace Context）
    Publish(ctx context.Context, routingKey string, msg *StandardMessage) error
}

// Subscriber 消息订阅接口
type Subscriber interface {
    // Subscribe 订阅消息
    // handler 返回 error 时消息会被 Nack
    Subscribe(queueName string, handler MessageHandler) error

    // Unsubscribe 取消订阅
    Unsubscribe(queueName string) error
}

// MessageHandler 消息处理函数
// ctx 包含从消息头提取的 trace context
type MessageHandler func(ctx context.Context, msg *StandardMessage) error
```

### 3.2 Tracer 接口

```go
// pkg/tracer/provider.go

// Provider Tracer Provider 接口
type Provider interface {
    // Tracer 获取 Tracer 实例
    Tracer(name string) trace.Tracer

    // Shutdown 关闭 Provider
    Shutdown(ctx context.Context) error
}

// NewProvider 创建 Tracer Provider
func NewProvider(cfg *TracerConfig) (Provider, error)

// pkg/tracer/http.go

// HTTPMiddleware Gin HTTP 追踪中间件
func HTTPMiddleware(tracer trace.Tracer) gin.HandlerFunc

// pkg/tracer/rabbitmq.go

// InjectContext 将 trace context 注入 RabbitMQ 消息头
func InjectContext(ctx context.Context, headers map[string]interface{})

// ExtractContext 从 RabbitMQ 消息头提取 trace context
func ExtractContext(ctx context.Context, headers map[string]interface{}) context.Context
```

### 3.3 Metrics 接口

```go
// pkg/metrics/collector.go

// Collector Metrics 收集器接口
type Collector interface {
    // Registry 获取 Prometheus Registry
    Registry() *prometheus.Registry

    // NewCounter 创建计数器
    // name 格式: iot_{component}_{metric_type}_{unit}
    NewCounter(name, help string, labels []string) *prometheus.CounterVec

    // NewHistogram 创建直方图
    NewHistogram(name, help string, labels []string, buckets []float64) *prometheus.HistogramVec

    // NewGauge 创建仪表盘
    NewGauge(name, help string, labels []string) *prometheus.GaugeVec
}

// 标准标签定义
const (
    LabelService     = "service"      // 服务名
    LabelVendor      = "vendor"       // 厂商
    LabelMessageType = "message_type" // 消息类型
    LabelStatus      = "status"       // 状态 (success, failed)
    LabelMethod      = "method"       // HTTP 方法
    LabelPath        = "path"         // HTTP 路径
    LabelCode        = "code"         // HTTP 状态码
)

// pkg/metrics/handler.go

// Handler 返回 Gin handler for /metrics endpoint
func Handler(collector Collector) gin.HandlerFunc
```

### 3.4 Repository 接口

```go
// pkg/repository/device.go

// DeviceRepository 设备数据访问接口
type DeviceRepository interface {
    // GetByDeviceSN 根据设备序列号查询设备
    GetByDeviceSN(ctx context.Context, deviceSN string) (*models.Device, error)

    // GetVendorByDeviceSN 根据设备序列号查询厂商
    // 用于生成 routing key
    GetVendorByDeviceSN(ctx context.Context, deviceSN string) (string, error)

    // Create 创建设备
    Create(ctx context.Context, device *models.Device) error

    // Update 更新设备
    Update(ctx context.Context, device *models.Device) error

    // UpdateStatus 更新设备状态
    UpdateStatus(ctx context.Context, deviceSN string, status string) error
}
```

---

## 4. 错误类型

```go
// pkg/errors/errors.go

// 错误码定义
type ErrorCode int

const (
    // 通用错误 (1000-1999)
    ErrInternal          ErrorCode = 1000
    ErrInvalidParameter  ErrorCode = 1001
    ErrNotFound          ErrorCode = 1002
    ErrAlreadyExists     ErrorCode = 1003
    ErrUnauthorized      ErrorCode = 1004
    ErrForbidden         ErrorCode = 1005

    // 设备错误 (2000-2999)
    ErrDeviceNotFound    ErrorCode = 2000
    ErrDeviceOffline     ErrorCode = 2001
    ErrDeviceNotReady    ErrorCode = 2002

    // 消息错误 (3000-3999)
    ErrInvalidMessage    ErrorCode = 3000
    ErrInvalidRoutingKey ErrorCode = 3001
    ErrMessageTimeout    ErrorCode = 3002
    ErrTraceContextMissing ErrorCode = 3003

    // 连接错误 (4000-4999)
    ErrRabbitMQConnection ErrorCode = 4000
    ErrDatabaseConnection ErrorCode = 4001
    ErrInfluxDBConnection ErrorCode = 4002
)

// Error 业务错误
type Error struct {
    Code    ErrorCode `json:"code"`
    Message string    `json:"message"`
    Details string    `json:"details,omitempty"`
}

func (e *Error) Error() string

// New 创建错误
func New(code ErrorCode, message string) *Error

// Wrap 包装错误
func Wrap(err error, code ErrorCode, message string) *Error

// Is 判断错误类型
func Is(err error, code ErrorCode) bool
```

---

## 5. 配置文件示例

### 5.1 config.dev.yaml

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

database:
  postgres:
    host: localhost
    port: 5432
    user: postgres
    password: postgres
    dbname: umos_iot_dev
    sslmode: disable
    max_idle_conns: 10
    max_open_conns: 50
    conn_max_lifetime: 1h
  influxdb:
    url: http://localhost:8086
    token: dev-token
    org: umos
    bucket: iot_dev

rabbitmq:
  url: amqp://guest:guest@localhost:5672/
  exchange_name: iot
  exchange_type: topic
  prefetch_count: 10
  retry:
    max_retries: 10
    initial_delay: 1s
    max_delay: 30s
    multiplier: 2.0

tracer:
  enabled: true
  endpoint: http://localhost:4318/v1/traces
  service_name: ${SERVICE_NAME}
  sampling_rate: 1.0  # 开发环境 100% 采样
  batch_timeout: 5s
  max_queue_size: 2048

metrics:
  enabled: true
  path: /metrics
  namespace: iot

logger:
  level: debug
  format: json
  output: stdout
```

### 5.2 config.prod.yaml

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

database:
  postgres:
    host: ${POSTGRES_HOST}
    port: ${POSTGRES_PORT}
    user: ${POSTGRES_USER}
    password: ${POSTGRES_PASSWORD}
    dbname: ${POSTGRES_DB}
    sslmode: require
    max_idle_conns: 20
    max_open_conns: 100
    conn_max_lifetime: 1h
  influxdb:
    url: ${INFLUXDB_URL}
    token: ${INFLUXDB_TOKEN}
    org: ${INFLUXDB_ORG}
    bucket: ${INFLUXDB_BUCKET}

rabbitmq:
  url: ${RABBITMQ_URL}
  exchange_name: iot
  exchange_type: topic
  prefetch_count: 50
  retry:
    max_retries: 10
    initial_delay: 1s
    max_delay: 30s
    multiplier: 2.0

tracer:
  enabled: true
  endpoint: ${TEMPO_ENDPOINT}
  service_name: ${SERVICE_NAME}
  sampling_rate: 0.1  # 生产环境 10% 采样
  batch_timeout: 5s
  max_queue_size: 4096

metrics:
  enabled: true
  path: /metrics
  namespace: iot

logger:
  level: info
  format: json
  output: stdout
```

---

## 6. 死信队列命名规范

```go
// 死信队列命名规范
const (
    // 死信 Exchange
    DeadLetterExchange = "iot.dlx"

    // 死信 Queue 命名: iot.dlq.{original_queue_name}
    // 示例: iot.dlq.uplink.property.report
    DeadLetterQueuePrefix = "iot.dlq."
)

// 死信消息处理策略
// 1. 消息进入死信队列后，记录到 MessageLog 表
// 2. 保留原始 trace_id 和所有 headers
// 3. 添加 x-death 头记录死信原因和时间
// 4. 可通过 Grafana 告警监控死信队列长度
```

---

## Version

**Version**: 1.0.0
**Created**: 2025-02-04
**Last Updated**: 2025-02-04
