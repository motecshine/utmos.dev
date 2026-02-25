# UMOS IoT Platform

统一的多厂商 IoT 设备管理平台，支持多种通信协议和设备类型。

## 项目概述

UMOS (Unified Multi-vendor IoT Operating System) 是一个基于微服务架构的 IoT 平台，旨在提供统一的设备管理、数据采集、消息路由和可观测性能力。

## 核心特性

- **多协议支持**: MQTT 5.0、HTTPS RESTful API、WebSocket
- **多厂商支持**: 通过统一的物模型抽象支持不同厂商的设备
- **分布式追踪**: 基于 OpenTelemetry 的完整追踪链路
- **统一 Metrics**: Prometheus 格式的指标监控
- **消息路由**: 基于 RabbitMQ 的多厂商消息路由机制

## 技术栈

- **语言**: Go 1.22+
- **Web 框架**: Gin Framework
- **ORM**: GORM
- **日志**: logrus
- **消息队列**: RabbitMQ
- **MQTT Broker**: VerneMQ
- **数据库**: PostgreSQL, InfluxDB
- **可观测性**: Prometheus, Loki, Tempo, Grafana

## 项目结构

```
umos/
├── cmd/                    # 服务入口
│   ├── iot-api/           # HTTP API 服务
│   ├── iot-ws/            # WebSocket 服务
│   ├── iot-uplink/        # 上行消息处理服务
│   ├── iot-downlink/      # 下行消息处理服务
│   └── iot-gateway/       # MQTT 网关服务
├── internal/              # 内部实现
│   └── shared/           # 共享代码
├── pkg/                   # 公共包
│   ├── models/           # 数据模型
│   ├── metrics/          # Metrics 包
│   ├── tracer/           # 追踪包
│   └── rabbitmq/         # RabbitMQ 客户端
├── api/                   # API 定义
├── deployments/          # 部署配置
└── tests/                # 测试代码
```

## 快速开始

### 前置要求

- Go 1.22+
- Docker & Docker Compose
- Make (可选)

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd umos
```

2. **安装依赖**
```bash
go mod download
```

3. **启动基础设施**
```bash
docker-compose up -d
```

4. **运行服务**
```bash
# 运行单个服务
make run-api

# 或运行所有服务
make build
```

## 开发指南

### 代码规范

本项目严格遵循 [Uber Go 编码规范](https://github.com/uber-go/guide)，所有代码必须通过以下检查：

- `golangci-lint` 代码风格检查
- `misspell` 拼写检查
- `go vet` 静态分析

### 运行测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage
```

### 代码检查

```bash
# 运行 linter
make lint

# 格式化代码
make fmt
```

## 部署

### Docker 部署

```bash
# 构建所有服务镜像
make docker-build

# 使用 docker-compose 启动
make docker-compose-up
```

### Kubernetes 部署

Kubernetes manifests 位于 `deployments/kubernetes/` 目录。

## 文档

- [架构文档](docs/architecture/)
- [API 文档](api/v1/)
- [开发指南](docs/development/)

## 贡献

欢迎贡献代码！请确保：

1. 代码遵循 Uber Go 编码规范
2. 所有测试通过
3. 通过 linter 检查
4. 无拼写错误

## 许可证

[待定]

