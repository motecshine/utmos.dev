# Tasks: Core Services Implementation

**Feature**: 004-core-services-implementation
**Generated**: 2025-02-05
**Total Tasks**: 50

## Task Summary

| Phase | Tasks | Status |
|-------|-------|--------|
| Phase 1: iot-gateway | T001-T010 | Complete |
| Phase 2: iot-uplink | T011-T018 | Complete |
| Phase 3: iot-downlink | T019-T026 | Complete |
| Phase 4: iot-api | T027-T036 | Complete |
| Phase 5: iot-ws | T037-T044 | Complete |
| Phase 6: Integration & NFR | T045-T050 | Complete |

---

## Phase 1: iot-gateway (P1)

### T001: 添加 MQTT 客户端依赖
- [x] **Description**: 添加 paho.mqtt.golang 依赖到 go.mod
- **Files**: `go.mod`
- **Acceptance**: `go mod tidy` 成功

### T002: 实现 MQTT 客户端连接管理
- [x] **Description**: 实现 MQTT 客户端连接 VerneMQ，支持自动重连
- **Files**: `internal/gateway/mqtt/client.go`, `internal/gateway/mqtt/client_test.go`
- **Depends**: T001
- **Acceptance**: 单元测试通过，能够连接 VerneMQ

### T003: 实现设备认证
- [x] **Description**: 实现设备用户名/密码认证，查询 PostgreSQL 验证凭证
- **Files**: `internal/gateway/mqtt/auth.go`, `internal/gateway/mqtt/auth_test.go`, `internal/gateway/model/credential.go`
- **Depends**: T002
- **Acceptance**: 单元测试通过，能够验证设备凭证

### T004: 实现 MQTT 消息处理器
- [x] **Description**: 实现 MQTT 消息接收和处理逻辑
- **Files**: `internal/gateway/mqtt/handler.go`, `internal/gateway/mqtt/handler_test.go`
- **Depends**: T002
- **Acceptance**: 单元测试通过，能够接收 MQTT 消息

### T005: 实现上行消息桥接
- [x] **Description**: 实现 MQTT 消息到 RabbitMQ 的转发
- **Files**: `internal/gateway/bridge/uplink.go`, `internal/gateway/bridge/uplink_test.go`
- **Depends**: T004
- **Acceptance**: 单元测试通过，MQTT 消息能够转发到 RabbitMQ

### T006: 实现下行消息桥接
- [x] **Description**: 实现 RabbitMQ 消息到 MQTT 的转发
- **Files**: `internal/gateway/bridge/downlink.go`, `internal/gateway/bridge/downlink_test.go`
- **Depends**: T002
- **Acceptance**: 单元测试通过，RabbitMQ 消息能够转发到 MQTT

### T007: 实现设备连接状态管理
- [x] **Description**: 实现设备在线/离线状态跟踪
- **Files**: `internal/gateway/connection/manager.go`, `internal/gateway/connection/manager_test.go`
- **Depends**: T004
- **Acceptance**: 单元测试通过，能够跟踪设备连接状态

### T008: 实现 Gateway 服务层
- [x] **Description**: 整合 MQTT 客户端、桥接、连接管理到服务层
- **Files**: `internal/gateway/service.go`, `internal/gateway/service_test.go`
- **Depends**: T005, T006, T007
- **Acceptance**: 单元测试通过

### T009: 更新 Gateway main.go
- [x] **Description**: 更新 cmd/iot-gateway/main.go 集成业务逻辑
- **Files**: `cmd/iot-gateway/main.go`
- **Depends**: T008
- **Acceptance**: 服务能够启动并连接 VerneMQ

### T010: Gateway 集成测试
- [x] **Description**: 编写 Gateway 集成测试
- **Files**: `tests/integration/gateway_test.go`
- **Depends**: T009
- **Acceptance**: 集成测试通过

---

## Phase 2: iot-uplink (P1)

### T011: 实现消息处理器接口
- [x] **Description**: 定义消息处理器接口和基础实现
- **Files**: `internal/uplink/processor/processor.go`, `internal/uplink/processor/processor_test.go`
- **Acceptance**: 单元测试通过

### T012: 实现 DJI 消息处理器
- [x] **Description**: 订阅 RabbitMQ 队列 `iot.dji.#` (dji-adapter 输出)，处理标准化后的 DJI 消息。包含物模型映射逻辑：将 DJI 协议数据映射到标准物模型结构 (FR-008)。注意：iot-uplink 不直接调用 pkg/adapter/dji，而是消费 dji-adapter 服务的输出。
- **Files**: `internal/uplink/processor/dji.go`, `internal/uplink/processor/dji_test.go`
- **Depends**: T011
- **Acceptance**: 单元测试通过，能够处理 DJI 消息

### T013: 添加 InfluxDB 客户端依赖
- [x] **Description**: 添加 influxdb-client-go 依赖到 go.mod
- **Files**: `go.mod`
- **Acceptance**: `go mod tidy` 成功

### T014: 实现 InfluxDB 时序数据写入
- [x] **Description**: 实现遥测数据写入 InfluxDB
- **Files**: `internal/uplink/storage/influx.go`, `internal/uplink/storage/influx_test.go`
- **Depends**: T013
- **Acceptance**: 单元测试通过，能够写入 InfluxDB

### T015: 实现消息路由
- [x] **Description**: 实现消息路由到其他服务（iot-ws, iot-api）
- **Files**: `internal/uplink/router/router.go`, `internal/uplink/router/router_test.go`
- **Depends**: T012
- **Acceptance**: 单元测试通过

### T016: 实现 Uplink 服务层
- [x] **Description**: 整合处理器、存储、路由到服务层
- **Files**: `internal/uplink/service.go`, `internal/uplink/service_test.go`
- **Depends**: T012, T014, T015
- **Acceptance**: 单元测试通过

### T017: 更新 Uplink main.go
- [x] **Description**: 更新 cmd/iot-uplink/main.go 集成业务逻辑
- **Files**: `cmd/iot-uplink/main.go`
- **Depends**: T016
- **Acceptance**: 服务能够启动并处理消息

### T018: Uplink 集成测试
- [x] **Description**: 编写 Uplink 集成测试
- **Files**: `tests/integration/uplink_test.go`
- **Depends**: T017
- **Acceptance**: 集成测试通过

---

## Phase 3: iot-downlink (P1)

### T019: 实现消息分发器接口
- [x] **Description**: 定义消息分发器接口和基础实现
- **Files**: `internal/downlink/dispatcher/dispatcher.go`, `internal/downlink/dispatcher/dispatcher_test.go`
- **Acceptance**: 单元测试通过

### T020: 实现 DJI 消息分发器
- [x] **Description**: 发布服务调用到 RabbitMQ 队列 `iot.dji.service.call`，由 dji-adapter 服务转换为 DJI 协议格式。注意：iot-downlink 不直接调用 pkg/adapter/dji，而是发布消息供 dji-adapter 服务消费。
- **Files**: `internal/downlink/dispatcher/dji.go`, `internal/downlink/dispatcher/dji_test.go`
- **Depends**: T019
- **Acceptance**: 单元测试通过

### T021: 实现服务调用记录模型
- [x] **Description**: 实现 ServiceCall 数据模型和数据库迁移
- **Files**: `internal/downlink/model/service_call.go`
- **Acceptance**: 数据库迁移成功

### T022: 实现重试机制
- [x] **Description**: 实现指数退避重试和死信队列
- **Files**: `internal/downlink/retry/retry.go`, `internal/downlink/retry/retry_test.go`
- **Depends**: T021
- **Acceptance**: 单元测试通过

### T023: 实现消息路由到 Gateway
- [x] **Description**: 实现下行消息路由到 iot-gateway
- **Files**: `internal/downlink/router/router.go`, `internal/downlink/router/router_test.go`
- **Depends**: T020
- **Acceptance**: 单元测试通过

### T024: 实现 Downlink 服务层
- [x] **Description**: 整合分发器、重试、路由到服务层
- **Files**: `internal/downlink/service.go`, `internal/downlink/service_test.go`
- **Depends**: T020, T022, T023
- **Acceptance**: 单元测试通过

### T025: 更新 Downlink main.go
- [x] **Description**: 更新 cmd/iot-downlink/main.go 集成业务逻辑
- **Files**: `cmd/iot-downlink/main.go`
- **Depends**: T024
- **Acceptance**: 服务能够启动并处理消息

### T026: Downlink 集成测试
- [x] **Description**: 编写 Downlink 集成测试
- **Files**: `tests/integration/downlink_test.go`
- **Depends**: T025
- **Acceptance**: 集成测试通过

---

## Phase 4: iot-api (P2)

### T027: 实现设备管理 Handler
- [x] **Description**: 实现设备 CRUD API Handler
- **Files**: `internal/api/handler/device.go`, `internal/api/handler/device_test.go`
- **Acceptance**: 单元测试通过

### T028: 实现服务调用 Handler
- [x] **Description**: 实现服务调用 API Handler
- **Files**: `internal/api/handler/service.go`, `internal/api/handler/service_test.go`
- **Depends**: T027
- **Acceptance**: 单元测试通过

### T029: 实现遥测查询 Handler
- [x] **Description**: 实现遥测数据查询 API Handler
- **Files**: `internal/api/handler/telemetry.go`, `internal/api/handler/telemetry_test.go`
- **Depends**: T014
- **Acceptance**: 单元测试通过

### T030: 实现认证中间件
- [x] **Description**: 实现 API 认证中间件
- **Files**: `internal/api/middleware/auth.go`, `internal/api/middleware/auth_test.go`
- **Acceptance**: 单元测试通过

### T031: 实现追踪中间件
- [x] **Description**: 实现分布式追踪中间件
- **Files**: `internal/api/middleware/trace.go`, `internal/api/middleware/trace_test.go`
- **Acceptance**: 单元测试通过

### T032: 实现路由配置
- [x] **Description**: 配置 API 路由
- **Files**: `internal/api/router.go`, `internal/api/router_test.go`
- **Depends**: T027, T028, T029, T030, T031
- **Acceptance**: 单元测试通过

### T033: 更新 API main.go
- [x] **Description**: 更新 cmd/iot-api/main.go 集成业务逻辑
- **Files**: `cmd/iot-api/main.go`
- **Depends**: T032
- **Acceptance**: 服务能够启动并响应 API 请求

### T034: 生成 OpenAPI 文档
- [x] **Description**: 使用 swag 生成 OpenAPI 文档
- **Files**: `docs/swagger.json`, `docs/swagger.yaml`
- **Depends**: T033
- **Acceptance**: 文档生成成功

### T035: API 集成测试
- [x] **Description**: 编写 API 集成测试
- **Files**: `tests/integration/api_test.go`
- **Depends**: T033
- **Acceptance**: 集成测试通过

### T036: API 契约测试
- [x] **Description**: 编写 API 契约测试验证 OpenAPI 规范
- **Files**: `tests/contract/api_contract_test.go`
- **Depends**: T034
- **Acceptance**: 契约测试通过

---

## Phase 5: iot-ws (P2)

### T037: 添加 WebSocket 依赖
- [x] **Description**: 添加 gorilla/websocket 依赖到 go.mod
- **Files**: `go.mod`
- **Acceptance**: `go mod tidy` 成功

### T038: 实现 WebSocket Hub
- [x] **Description**: 实现 WebSocket 连接中心
- **Files**: `internal/ws/hub/hub.go`, `internal/ws/hub/hub_test.go`
- **Depends**: T037
- **Acceptance**: 单元测试通过

### T039: 实现 WebSocket Client
- [x] **Description**: 实现 WebSocket 客户端连接管理
- **Files**: `internal/ws/hub/client.go`, `internal/ws/hub/client_test.go`
- **Depends**: T038
- **Acceptance**: 单元测试通过

### T039-A: 实现 WebSocket 心跳检测
- [x] **Description**: 实现 ping/pong 心跳机制，30s 超时断开无响应连接 (FR-022)
- **Files**: `internal/ws/hub/heartbeat.go`, `internal/ws/hub/heartbeat_test.go`
- **Depends**: T039
- **Acceptance**: 单元测试通过，超时连接能够正确断开

### T040: 实现订阅管理器
- [x] **Description**: 实现消息订阅管理
- **Files**: `internal/ws/subscription/manager.go`, `internal/ws/subscription/manager_test.go`
- **Depends**: T039, T039-A
- **Acceptance**: 单元测试通过

### T041: 实现消息推送
- [x] **Description**: 实现从 RabbitMQ 接收消息并推送到 WebSocket 客户端
- **Files**: `internal/ws/push/pusher.go`, `internal/ws/push/pusher_test.go`
- **Depends**: T040
- **Acceptance**: 单元测试通过

### T042: 实现 WS 服务层
- [x] **Description**: 整合 Hub、订阅、推送到服务层
- **Files**: `internal/ws/service.go`, `internal/ws/service_test.go`
- **Depends**: T038, T040, T041
- **Acceptance**: 单元测试通过

### T043: 更新 WS main.go
- [x] **Description**: 更新 cmd/iot-ws/main.go 集成业务逻辑
- **Files**: `cmd/iot-ws/main.go`
- **Depends**: T042
- **Acceptance**: 服务能够启动并处理 WebSocket 连接

### T044: WS 集成测试
- [x] **Description**: 编写 WebSocket 集成测试
- **Files**: `tests/integration/ws_test.go`
- **Depends**: T043
- **Acceptance**: 集成测试通过

---

## Phase 6: Integration & NFR Testing

### T045: 端到端测试
- [x] **Description**: 编写完整数据流端到端测试
- **Files**: `tests/integration/e2e_test.go`
- **Depends**: T010, T018, T026, T035, T044
- **Acceptance**: 端到端测试通过，验证完整上下行数据流

### T046: 性能测试
- [x] **Description**: 验证消息处理延迟 < 100ms (P95) (NFR-001)
- **Files**: `tests/integration/performance_test.go`
- **Depends**: T045
- **Acceptance**: P95 延迟 < 100ms，测试报告生成

### T047: 设备负载测试
- [x] **Description**: 验证支持 1000+ 设备同时在线 (NFR-002)
- **Files**: `tests/integration/load_device_test.go`
- **Depends**: T010
- **Acceptance**: 1000 设备并发连接测试通过

### T048: WebSocket 负载测试
- [x] **Description**: 验证支持 10000+ WebSocket 连接 (NFR-003)
- **Files**: `tests/integration/load_ws_test.go`
- **Depends**: T044
- **Acceptance**: 10000 WebSocket 连接测试通过

### T049: 覆盖率验证
- [x] **Description**: 验证单元测试覆盖率 >= 80% (TDD-002)
- **Files**: `Makefile` (coverage target)
- **Depends**: T045
- **Acceptance**: `make coverage` 报告显示 >= 80%

### T050: 可用性测试
- [x] **Description**: 验证服务可用性 > 99.9% (NFR-004)，包括故障恢复测试
- **Files**: `tests/integration/availability_test.go`
- **Depends**: T045
- **Acceptance**: 服务故障恢复测试通过

---

## Execution Order

```
Phase 1 (P1): T001 → T002 → [P] T003, T004 → [P] T005, T006, T007 → T008 → T009 → T010
Phase 2 (P1): T011 → T012, T013 → [P] T014, T015 → T016 → T017 → T018
Phase 3 (P1): T019 → T020, T021 → [P] T022, T023 → T024 → T025 → T026
Phase 4 (P2): T027 → [P] T028, T029, T030, T031 → T032 → T033 → [P] T034, T035, T036
Phase 5 (P2): T037 → T038 → T039 → T039-A → T040 → T041 → T042 → T043 → T044
Phase 6:      T045 → [P] T046, T047, T048, T049, T050
```

## Notes

- P1 任务（Phase 1-3）是核心功能，必须优先完成
- P2 任务（Phase 4-5）可以在 P1 完成后并行开发
- 所有任务必须遵循 TDD 原则：先写测试，再实现功能
- 单元测试覆盖率目标 >= 80%
- [P] 标记表示可并行执行的任务

## Edge Case Handling

| Edge Case | Handling Task | Strategy |
|-----------|---------------|----------|
| MQTT 连接断开 | T002 | 自动重连 (paho.mqtt SetAutoReconnect) |
| RabbitMQ 连接断开 | T005, T006 | 自动重连 (amqp091-go reconnect) |
| 消息处理失败 | T022 | 死信队列 + 告警 |
| 大量设备同时上线 | T007 | 连接池 + 限流 |
| WebSocket 连接数过多 | T038 | 连接限制 (max 10000) |
| WebSocket 心跳超时 | T039-A | 30s 超时断开 |
