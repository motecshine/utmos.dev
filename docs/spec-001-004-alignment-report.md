# Spec 001-004 与代码实现对齐分析报告

## 1. 报告目的

本文档用于分析 `specs/001-project-setup` 至 `specs/004-core-services-implementation` 与当前仓库实现之间的对齐程度，判断：

- spec 目标是否已经在代码中落地
- 运行时功能是否完整
- 测试通过是否代表真实链路已经打通
- 当前实现与 spec 之间的主要差异点是什么

分析基于当前仓库代码、测试、启动入口和契约文档完成。当前仓库执行 `go test ./...` 结果为通过，但这不等同于所有 spec 所要求的运行时闭环已经成立。

---

## 2. 总体结论

整体结论是：**001-004 与代码为部分对齐，但整体不能判定为功能完整**。

当前仓库具有以下特点：

- 基础设施层已经成型，包括 tracing、logging、metrics、RabbitMQ、GORM migration、核心服务骨架。
- 大量单元测试、集成测试、契约测试已经存在，说明接口和模块边界有较强覆盖。
- 但若按 spec 的“端到端运行时链路”要求来判断，仍存在多个关键断点。
- 其中最严重的问题集中在：
  - `iot-gateway -> dji-adapter` 上行消息格式不兼容
  - `iot-downlink -> dji-adapter -> iot-gateway` 下行链路与 spec 不一致
  - `cmd/dji-adapter` 未启用完整 DJI handler/router 初始化
  - 命令状态机、late response、幂等、审计链路未真正闭环
  - 实时订阅鉴权与 spec 004 合约不一致
  - readiness 依赖检查不满足 spec 004

建议结论：

- `001`: 基本能力大体落地，但仍有闭环缺口
- `002`: 适配器框架已建立，但运行时落地不完整
- `003`: 协议代码资产很多，但“完整 DJI 协议实现”未真正上线
- `004`: 核心服务可运行，但关键业务要求未全部闭环

---

## 3. 分项结论概览

| Spec | 主题 | 对齐度 | 结论 |
|------|------|--------|------|
| 001 | 基础设施、Trace、Routing、Metrics | 中等偏高 | 有实现，但未完全闭环 |
| 002 | 协议适配器框架 | 中等 | 框架已成型，运行时仍偏骨架 |
| 003 | DJI 全量协议实现 | 偏低 | 代码很多，但实际运行链路未完整启用 |
| 004 | 核心服务实现 | 中等偏低 | 骨架齐全，但关键验收项不完整 |

---

## 4. 关键发现

### 4.1 最关键断点一：`gateway -> dji-adapter` 上行链路不兼容

spec 002/003 的设计是：

- `iot-gateway` 将原始厂商消息发布到 `iot.raw.{vendor}.uplink`
- `dji-adapter` 订阅原始 DJI 上行消息
- `dji-adapter` 解析原始 topic/payload，转换后再发布标准消息

当前代码中：

- `iot-gateway` 在 [internal/gateway/bridge/uplink.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/gateway/bridge/uplink.go#L83) 使用 `iot.raw.{vendor}.uplink` 作为 routing key，这一点与 spec 一致。
- 但它发布的消息体不是“DJI 原始 payload + AMQP header 元信息”，而是一个 `StandardMessage`，其中 `original topic` 被塞进了 `data` 和 `protocol_meta`。
- `dji-adapter` 在 [cmd/dji-adapter/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/dji-adapter/main.go#L234) 处理上行消息时，却要求从 `msg.Headers["original_topic"]` 提取 topic，并把 `msg.Body` 直接当作 DJI 原始 payload 传给 `ParseRawMessage`。

这意味着：

- `gateway` 发出的结构，`dji-adapter` 运行时无法按预期解析。
- 即使 routing key 对了，真正的消息内容模型仍然不兼容。
- 所以 `iot-gateway -> dji-adapter` 的核心运行闭环实际上没有打通。

相关证据：

- [internal/gateway/bridge/uplink.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/gateway/bridge/uplink.go#L83)
- [pkg/rabbitmq/message.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/rabbitmq/message.go#L21)
- [cmd/dji-adapter/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/dji-adapter/main.go#L234)

### 4.2 最关键断点二：下行链路与 spec 设计不一致

spec 003 的下行设计是：

`iot-api -> iot-downlink -> dji-adapter -> iot-gateway -> MQTT device`

当前代码中存在两套互相冲突的路径：

1. `iot-downlink` 自己有一套路由到 `iot.gateway.*` 的逻辑  
   证据：
   - [internal/downlink/router/router.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/downlink/router/router.go#L19)
   - [internal/downlink/service.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/downlink/service.go#L236)

2. `iot-gateway` 的 `DownlinkBridge` 直接把 RabbitMQ 中某些标准消息转成 MQTT 下发  
   证据：
   - [internal/gateway/bridge/downlink.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/gateway/bridge/downlink.go#L177)

3. `pkg/adapter/dji/downlink` 又定义了一套 DJI dispatcher，但其 routing key 格式并不满足 `iot.{vendor}.{service}.{action}`  
   证据：
   - [pkg/adapter/dji/downlink/dispatcher.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/adapter/dji/downlink/dispatcher.go#L16)

4. `cmd/dji-adapter` 订阅下行时绑定的是 `iot.dji.*.service.#`，与 DJI dispatcher 实际发出的 `iot.dji.service.call` / `iot.dji.property.set` 风格并不一致  
   证据：
   - [cmd/dji-adapter/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/dji-adapter/main.go#L321)

结果是：

- spec 设计要求的“协议转换必须经 dji-adapter”没有真正统一落实。
- 实际运行时更像是：`iot-downlink` 可以直接绕过 `dji-adapter` 把消息打给 `iot-gateway`。
- 这使 002/003 的 adapter 角色在核心链路中被削弱甚至绕过。

### 4.3 最关键断点三：`cmd/dji-adapter` 没有启用完整协议处理器

仓库中 DJI 相关代码很多：

- `pkg/adapter/dji/router/*`
- `pkg/adapter/dji/handler/*`
- `pkg/adapter/dji/protocol/*`
- `pkg/adapter/dji/init/init.go`

从代码结构上看，完整版本应该是：

- 初始化 `ServiceRouter`
- 初始化 `EventRouter`
- 注册 DJI 服务命令、事件、DRC、wayline、camera、file、firmware、live 等 handler
- 将这些 handler 挂到 `dji.Adapter` 上

完整初始化逻辑在：

- [pkg/adapter/dji/init/init.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/adapter/dji/init/init.go#L80)

但 `cmd/dji-adapter` 实际只做了：

- `dji.Register()`
- `adapter.Get(dji.VendorDJI)`

即只注册了 `NewAdapter()` 的基础版本，没有调用初始化逻辑。证据：

- [cmd/dji-adapter/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/dji-adapter/main.go#L69)
- [pkg/adapter/dji/adapter.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/adapter/dji/adapter.go#L21)

这会导致：

- 运行时没有启用完整 handler registry。
- `HandleMessage()` 无法利用 service/event/request/drc 等专门处理器。
- 很多协议能力只存在于代码与测试中，没有真正进入生产入口。

### 4.4 最关键断点四：DJI topic 支持范围小于 spec 003

spec 003 明确覆盖：

- `osd`
- `state`
- `events`
- `services_reply`
- `requests`
- `requests_reply`
- `status`
- `drc/up`
- `drc/down`
- 以及 wayline、camera、file、firmware、live 等对应方法域

但运行时 topic 解析器 [pkg/adapter/dji/topic.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/adapter/dji/topic.go#L36) 实际只支持：

- `osd`
- `state`
- `services`
- `services_reply`
- `events`
- `status`
- `status_reply`

不支持：

- `requests`
- `requests_reply`
- `events_reply`
- `drc/up`
- `drc/down`

而这些 topic 常量本身其实已定义在 [pkg/adapter/dji/types.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/adapter/dji/types.go#L34)。

这说明：

- 协议范围在“类型定义层”比“运行时解析层”更完整。
- 003 所要求的“完整实现”在入口级别并未实现。

### 4.5 Trace 在 `dji-adapter` 内未继续标准传播

001 和 003 都要求：

- HTTP 和 RabbitMQ 中使用 W3C Trace Context
- 服务间持续传播 `traceparent` / `tracestate`

仓库统一 RabbitMQ 封装是支持这点的：

- `Publisher.Publish()` 自动注入 trace header  
  [pkg/rabbitmq/publisher.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/rabbitmq/publisher.go#L27)
- `Subscriber.handleDelivery()` 自动提取 trace header  
  [pkg/rabbitmq/subscriber.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/rabbitmq/subscriber.go#L72)

但 `cmd/dji-adapter` 没有复用这个统一发布路径，而是自己直接调用 `channel.PublishWithContext(context.Background(), ...)`：

- 上行发布没有带 trace header  
  [cmd/dji-adapter/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/dji-adapter/main.go#L289)
- 下行发布也只写了 `original_topic`、`device_sn`  
  [cmd/dji-adapter/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/dji-adapter/main.go#L417)

结果：

- trace 进入 `dji-adapter` 后不能保证继续跨服务传播。
- 001 关于服务间 trace 连贯性的目标未完全达到。

### 4.6 命令生命周期和响应闭环未真正完成

spec 004 要求命令生命周期必须支持：

- `pending`
- `sent`
- `retrying`
- `success`
- `failed`
- `timeout`
- late response 审计记录且不能 reopen terminal state

当前仓库中：

- `ServiceCall` model 定义了这些状态和状态迁移方法  
  [internal/downlink/model/service_call.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/downlink/model/service_call.go#L13)
- `FindByTID()` 也存在  
  [internal/downlink/model/service_call.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/downlink/model/service_call.go#L225)

但是我没有在代码中找到：

- 设备响应被消费并映射回 `ServiceCall`
- `FindByTID()` 被实际使用
- `MarkSuccess()` / `MarkTimeout()` / `MarkFailed()` / `MarkRetrying()` 在真实链路中被调用
- late response 被记录审计但不变更终态的实现

进一步看 `iot-downlink`：

- `consumeMessages()` 仍是空实现，仅等待上下文取消  
  [internal/downlink/service.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/downlink/service.go#L226)

所以当前状态是：

- 命令状态模型存在
- 契约测试和结构存在
- 但真正设备响应驱动的生命周期闭环未落地

### 4.7 幂等与审计落地不足

spec 004 要求：

- duplicate upstream message / duplicate service response 按 `tid + bid` 幂等
- late response 要审计
- 认证失败、处理失败、重试、终态都要形成审计记录

当前代码中：

- `MessageLog` model 存在  
  [pkg/models/message_log.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/models/message_log.go#L31)
- AutoMigrate 也包含该表  
  [pkg/models/migrate.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/models/migrate.go#L9)

但我没有找到：

- 对 `message_logs` 的实际写入路径
- 基于 `tid + bid` 的统一去重逻辑
- duplicate reply / late reply 的审计处理逻辑

因此这里属于“模型已建，业务未闭环”。

### 4.8 实时订阅接口与鉴权不符合 spec 004

spec 004 合约明确包含：

- `/api/v1/realtime/subscriptions`

证据：

- [specs/004-core-services-implementation/contracts/openapi.yaml](/Users/laputalaputa/Desktop/code/utmos.dev/specs/004-core-services-implementation/contracts/openapi.yaml#L431)

但 API 实际路由中只有：

- `/api/v1/devices`
- `/api/v1/services`
- `/api/v1/telemetry`

没有 realtime subscription HTTP 接口。证据：

- [internal/api/router.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/api/router.go#L174)

同时，WebSocket 订阅逻辑中：

- 连接时只从 query 参数提取 `device_sn` 和 `user_id`
- 收到 `subscribe` 消息后直接订阅 topic
- 没有 per-topic 授权校验
- 没有“有效 topic 激活、非法 topic 返回显式错误”的 ack 语义

证据：

- [internal/ws/service.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/ws/service.go#L191)
- [internal/ws/service.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/ws/service.go#L242)

因此：

- spec 004 FR-013 / FR-020 只实现了“订阅能力”，没有实现“授权控制和合约接口”。

### 4.9 Readiness 不满足 spec 004 依赖要求

spec 004 对 readiness 的定义是严格的：

- `iot-gateway`: RabbitMQ + VerneMQ
- `iot-uplink`: RabbitMQ + database
- `iot-downlink`: RabbitMQ + database
- `iot-api`: RabbitMQ + PostgreSQL + InfluxDB
- `iot-ws`: RabbitMQ

当前实现中：

- `iot-gateway` readiness 检查 RabbitMQ + MQTT，基本接近 spec  
  [cmd/iot-gateway/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/iot-gateway/main.go#L96)
- `iot-uplink` readiness 只检查 RabbitMQ 和 service running，没有检查 InfluxDB  
  [cmd/iot-uplink/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/iot-uplink/main.go#L111)
- `iot-downlink` 根本没有初始化数据库，readiness 也只检查 RabbitMQ + service  
  [cmd/iot-downlink/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/iot-downlink/main.go#L47)
- `iot-ws` readiness 只看 service running，没有检查 RabbitMQ  
  [cmd/iot-ws/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/iot-ws/main.go#L100)
- `iot-api` 检查数据库和 telemetry handler，但没有把 RabbitMQ 连接性算入 readiness  
  [internal/api/router.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/api/router.go#L123)

这部分与 spec 004 存在明确偏差。

### 4.10 MQTT 认证能力存在，但未看见完整 admission 集成

spec 004 FR-001 要求平台必须认证设备后才允许设备通信。

当前仓库中：

- `internal/gateway/mqtt/auth.go` 实现了设备凭证认证
- 集成测试也覆盖了认证生命周期

证据：

- [internal/gateway/mqtt/auth.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/gateway/mqtt/auth.go#L32)
- [tests/integration/gateway_auth_test.go](/Users/laputalaputa/Desktop/code/utmos.dev/tests/integration/gateway_auth_test.go#L29)

但从运行入口和 gateway service 看：

- 我没有看到它被真正接到 MQTT broker admission path 上
- 没看到 VerneMQ auth webhook / plugin / ACL 集成代码

因此可以判断为：

- 认证模块能力存在
- 但“平台运行时强制 admission”是否真正生效，当前代码证据不足

### 4.11 001 的 routing key 规范没有完全统一

001 要求统一格式：

- `iot.{vendor}.{service}.{action}`

仓库里确实存在该规范实现：

- [pkg/rabbitmq/routing.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/rabbitmq/routing.go#L47)
- [tests/integration/message_flow_test.go](/Users/laputalaputa/Desktop/code/utmos.dev/tests/integration/message_flow_test.go#L34)

但运行时内部仍存在大量非规范 key：

- `iot.ws.property`
- `iot.ws.event`
- `iot.api.property`
- `iot.gateway.command`
- `iot.gateway.downlink`

证据：

- [internal/uplink/router/router.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/uplink/router/router.go#L18)
- [internal/downlink/router/router.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/downlink/router/router.go#L19)

所以当前状态更像是：

- “外部规范”已定义
- “内部服务间 key”仍有历史风格并存
- 未真正统一为 spec 001 规定的一套模型

### 4.12 001 的 `device_sn -> vendor` 查询能力存在，但运行时并未统一使用

001 要求厂商应通过 `sn + vendor` 建模，并通过 `device_sn` 查询 vendor。

当前仓库中：

- `GetVendorByDeviceSN()` 已实现  
  [pkg/repository/device.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/repository/device.go#L37)

但 gateway topic 解析时仍直接根据 topic 格式硬编码 vendor：

- `thing` / `sys` 默认归类为 `dji`
- 否则从 topic 前缀拿 vendor

证据：

- [internal/gateway/mqtt/handler.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/gateway/mqtt/handler.go#L147)

这说明：

- repository 能力已存在
- 运行时链路还没完全按 001 设计切换到“查库获取 vendor”

### 4.13 Metrics 框架已实现，但 middleware 级自动采集未真正接入服务

001 对统一 metrics 包的要求较高，包括：

- RabbitMQ / PostgreSQL / InfluxDB 自动采集
- 统一命名和标签规范
- `/metrics` 暴露

当前仓库中：

- `pkg/metrics` 基础能力完整
- 命名规范基本符合 `iot_{component}_{metric}_{unit}`
- 服务入口都暴露了 `/metrics`

证据：

- [pkg/metrics/collector.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/metrics/collector.go#L1)
- [pkg/metrics/middleware.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/metrics/middleware.go#L1)
- [internal/api/router.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/api/router.go#L167)

但是：

- 我没有看到 `NewMiddlewareMetrics()` 在服务启动入口中被真正接入
- middleware 指标更多还是“定义了”，不是“自动接到 RabbitMQ/Postgres/Influx 客户端生命周期”

证据：

- `NewMiddlewareMetrics` 只在 metrics 包和测试中出现，没有形成服务级接线

因此这里属于：

- 基础框架到位
- 自动化中间件采集落地不足

### 4.14 004 的 capability / thing model 校验没有真正实现

spec 004 FR-004 要求：

- 上行消息在存储或转发前，需要根据 capability model / thing model 做校验

当前仓库中：

- `ThingModel` 数据模型存在  
  [pkg/models/thing_model.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/models/thing_model.go#L10)
- `Device` 也能关联 `ThingModelID`  
  [pkg/models/device.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/models/device.go#L20)

但我没有找到：

- uplink 处理前基于 thing model 的 schema 验证
- property/event/service 与 capability model 的约束验证代码

这说明：

- 数据模型存在
- spec 要求的业务校验并未落地

### 4.15 Gateway 限流要求未实现

spec 004 FR-018 要求：

- gateway 对 MQTT 连接施加 token-bucket rate limiting
- 超限返回 CONNACK 错误

我在 `internal/gateway` 和测试中没有找到：

- token bucket
- rate limiter
- CONNACK reject
- reconnect surge control

因此这部分可以判定为未实现。

---

## 5. 各 Spec 详细评估

## 5.1 Spec 001: Project Setup with Trace / Routing / Metrics

### 已对齐部分

- 已实现 OpenTelemetry provider  
  [pkg/tracer/provider.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/tracer/provider.go#L20)
- 已实现 HTTP trace middleware  
  [internal/api/middleware/trace.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/api/middleware/trace.go#L1)
- 已实现 RabbitMQ trace inject/extract  
  [pkg/tracer/rabbitmq.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/tracer/rabbitmq.go#L54)
- 已实现统一 logger，并支持 trace_id / span_id 打日志  
  [pkg/logger/logger.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/logger/logger.go#L63)
- 已实现统一 metrics collector  
  [pkg/metrics/collector.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/metrics/collector.go#L17)
- 已实现统一配置加载  
  [internal/shared/config/loader.go](/Users/laputalaputa/Desktop/code/utmos.dev/internal/shared/config/loader.go#L1)
- 已实现 GORM AutoMigrate  
  [pkg/models/migrate.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/models/migrate.go#L9)

### 部分对齐部分

- routing key 规范实现了，但内部服务间未完全统一
- `/metrics` 已广泛存在，但 middleware 自动采集落地不足
- `device_sn -> vendor` 能力存在，但运行时链路未统一使用
- trace 跨服务基础具备，但 `dji-adapter` 自定义发布绕过统一封装，导致传播链断裂

### 未完全达成项

- `iot.{vendor}.{service}.{action}` 未成为全局唯一规范
- 所有 RabbitMQ 消息的 trace header 传播未严格保证
- 服务间完整 trace 链路在 adapter 节点处不稳定

### 评估结论

`001` 可以认为基础设施已经大部分具备，但还不能视为“完全对齐”。更准确地说，**001 已完成 70%-80%，剩余问题主要集中在运行时统一性而不是基础设施缺失**。

---

## 5.2 Spec 002: Protocol Adapter Design

### 已对齐部分

- 已有通用 `ProtocolAdapter` 接口  
  [pkg/adapter/interface.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/adapter/interface.go#L56)
- 已有 adapter registry  
  [pkg/adapter/registry.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/adapter/registry.go#L9)
- 已有 raw routing key 定义  
  [pkg/rabbitmq/routing.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/rabbitmq/routing.go#L140)
- 已有独立 `cmd/dji-adapter` 入口  
  [cmd/dji-adapter/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/dji-adapter/main.go#L1)
- 已有 protocol meta 结构  
  [pkg/rabbitmq/message.go](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/rabbitmq/message.go#L11)

### 部分对齐部分

- 适配器框架已存在，但 gateway/adaptor 的 raw message 契约不一致
- 独立服务已存在，但运行时没有完整启用 handler/router
- 标准消息框架存在，但物模型映射与 capability 验证未形成明确业务闭环

### 未完全达成项

- 适配器服务与原始消息队列之间的真实契约未打通
- 共享物模型只停留在数据模型层，没有转成运行时校验流程
- 作为独立服务的 dji-adapter 仍偏“示例骨架”，未达到规范中的完整适配器职责

### 评估结论

`002` 的**架构方向是对的，接口层是成立的，但运行时执行面只完成了一半左右**。

---

## 5.3 Spec 003: DJI Protocol Implementation

### 已对齐部分

- DJI 协议目录结构较完整
- OSD / State / Status / Service / Event / Wayline / Camera / DRC / Firmware / File / Live / WPML 等均有代码资产
- 大量单测和集成测试存在，说明协议对象建模比较充分

关键代码示例：

- [pkg/adapter/dji/protocol](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/adapter/dji/protocol/common/header.go)
- [pkg/adapter/dji/router](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/adapter/dji/router/service_router.go)
- [pkg/adapter/dji/handler](/Users/laputalaputa/Desktop/code/utmos.dev/pkg/adapter/dji/handler/service_handler.go)

### 部分对齐部分

- 代码层有很多完整模块，但运行时入口未启用
- 服务命令和事件 router 都有，但没有接入 `cmd/dji-adapter`
- `requests` / `requests_reply` / `drc` 等类型定义存在，但运行时 topic parser 不支持

### 未完全达成项

- 003 要求“完整实现 DJI vendor protocol uplink downlink”，当前实际未达成
- 完整 handler registry 未接入运行时
- 完整 topic 支持未接入运行时
- 完整上下行运行闭环未成立

### 评估结论

`003` **最容易产生“代码很多所以已经完成”的错觉**。事实上，当前更接近：

- 协议代码资产丰富
- 测试资产丰富
- 但生产运行入口仍是简化版本

因此 `003` 不应判断为已完成。

---

## 5.4 Spec 004: Core Services Implementation

### 已对齐部分

- API 层已实现设备 CRUD
- Service call API 已存在
- Telemetry 查询接口已存在
- WS 服务和订阅管理已存在
- Gateway、Uplink、Downlink、API、WS 五个核心服务入口都存在
- readiness / health / metrics 基础端点都基本具备

关键代码：

- [cmd/iot-api/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/iot-api/main.go#L1)
- [cmd/iot-gateway/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/iot-gateway/main.go#L1)
- [cmd/iot-uplink/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/iot-uplink/main.go#L1)
- [cmd/iot-downlink/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/iot-downlink/main.go#L1)
- [cmd/iot-ws/main.go](/Users/laputalaputa/Desktop/code/utmos.dev/cmd/iot-ws/main.go#L1)

### 部分对齐部分

- 设备认证模块存在，但未看到 admission 级集成
- 设备在线状态管理存在，但与 MQTT 真正连接层的结合证据有限
- service call 模型存在，但终态生命周期未闭环
- realtime 已实现，但未实现 spec 004 的授权控制

### 未完全达成项

- 命令状态机的完整闭环
- late response 审计且不 reopen terminal state
- duplicate / idempotent 统一处理
- realtime subscription contract endpoint
- per-topic authorization
- gateway token-bucket rate limiting
- readiness 依赖校验严格满足 spec
- capability model 校验

### 评估结论

`004` 更像是“核心服务骨架 + 局部业务实现已完成”，但距离 spec 定义的完整成品仍有明显差距。

---

## 6. 测试现状与风险说明

当前仓库执行：

```bash
go test ./...
```

结果通过。

这说明：

- 各模块内部设计大体一致
- 类型、基础行为、部分契约已经稳定

但这不代表：

- 服务启动后的消息契约一定互相兼容
- 真正的上行原始消息和 adapter 之间已经联通
- 真正的下行闭环已经按 spec 成立
- 命令状态和响应消费链已经打通

当前测试的主要风险是：

- 很多测试证明“模块可以工作”
- 但没有完全证明“生产入口按 spec 串起来也工作”

换句话说，**当前测试通过并不能推导出 001-004 功能已完整落地**。

---

## 7. 差异清单

以下差异按严重程度排序。

### P0 差异

1. `gateway -> dji-adapter` 上行消息契约不兼容
2. 下行链路与 spec 设计不一致，存在绕过 `dji-adapter` 的路径
3. `cmd/dji-adapter` 未启用完整 DJI router/handler 初始化
4. DJI dispatcher routing key 与 adapter 消费 binding 不匹配

### P1 差异

5. 命令状态机未闭环，响应消费链未落地
6. late response / duplicate response / idempotent 逻辑缺失
7. realtime subscription API 合约未实现
8. realtime per-topic authorization 未实现
9. readiness 未按 spec 004 依赖定义落地
10. trace 在 `dji-adapter` 内未继续使用统一传播机制

### P2 差异

11. capability / thing model 校验未接入 uplink
12. middleware 级 metrics 自动采集未真正接线
13. `device_sn -> vendor` 查询能力存在但未统一用于运行时路由
14. 内部 routing key 仍有非规范风格并存

### P3 差异

15. gateway token-bucket rate limiting 未实现
16. 审计日志模型存在但业务写入不足

---

## 8. 建议修复顺序

建议按以下顺序修复，避免先修外围再修主链路。

### 第一阶段：先修主链路

1. 统一 raw uplink / raw downlink 消息契约
2. 让 `iot-gateway` 与 `dji-adapter` 使用同一套原始消息结构
3. 明确是否坚持 `iot-api -> iot-downlink -> dji-adapter -> iot-gateway` 为唯一合法下行链路
4. 删除或改造绕过 `dji-adapter` 的直达 gateway 路由

### 第二阶段：启用完整 DJI adapter 运行时

1. 在 `cmd/dji-adapter` 中调用完整初始化逻辑
2. 让 service/event/request/drc handler 真正进入生产入口
3. 扩展 topic parser，补齐 `requests`、`requests_reply`、`events_reply`、`drc/up`、`drc/down`

### 第三阶段：闭合命令生命周期

1. 实现 downlink response consume
2. 基于 `tid/bid` 做 response correlation
3. 实现 `pending -> sent -> success/failed/timeout`
4. 实现 late response audit 和 duplicate idempotency

### 第四阶段：补齐 004 非功能和安全要求

1. realtime subscription API
2. per-topic authorization
3. readiness 严格依赖检查
4. gateway rate limiting
5. audit/message log 真正写入

---

## 9. 最终判断

如果问题是：

“分析 001-004 spec 与代码功能是否对齐，功能是否完整，找出差异点”

那么最终判断是：

- **对齐：部分对齐**
- **完整性：不完整**
- **主要问题：不是代码量不够，而是运行时链路没有真正按 spec 串起来**

最核心的一句话总结是：

**当前仓库已经具备较完整的模块、模型和测试资产，但 002-004 中最关键的 adapter 驱动链路与 command lifecycle 闭环尚未真正落地，因此不能认定 specs 001-004 已完整实现。**

