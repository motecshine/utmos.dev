# Quick Start Guide: UMOS IoT Platform Setup

本文档提供 UMOS IoT 平台的快速开始指南。

## 前置要求

- Go 1.22+
- Docker & Docker Compose
- Make (可选，但推荐)
- Git

## 1. 克隆项目

```bash
git clone <repository-url>
cd umos
```

## 2. 安装依赖

```bash
# 安装 Go 依赖
go mod download

# 安装开发工具
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/vektra/mockery/v2@latest
```

## 3. 启动基础设施

使用 Docker Compose 启动所有中间件：

```bash
docker-compose up -d
```

这将启动：
- PostgreSQL
- InfluxDB
- RabbitMQ
- VerneMQ
- Prometheus
- Loki
- Tempo
- Grafana

## 4. 配置环境变量

复制环境变量模板：

```bash
cp .env.example .env
```

编辑 `.env` 文件，配置数据库连接、消息队列等：

```env
# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=umos
POSTGRES_PASSWORD=umos123
POSTGRES_DB=umos_iot

# InfluxDB
INFLUXDB_URL=http://localhost:8086
INFLUXDB_TOKEN=your-token
INFLUXDB_ORG=umos
INFLUXDB_BUCKET=iot_data

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# VerneMQ
VERNEMQ_HOST=localhost
VERNEMQ_PORT=1883

# 服务配置
IOT_API_PORT=8080
IOT_WS_PORT=8081
```

## 5. 初始化数据库

```bash
# 运行数据库迁移
make migrate-up

# 或手动运行
go run cmd/migrate/main.go up
```

## 6. 运行代码检查

```bash
# 运行 golangci-lint
make lint

# 检查拼写错误
make spell-check
```

## 7. 运行测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage
```

## 8. 构建服务

```bash
# 构建所有服务
make build

# 构建单个服务
make build-api
make build-ws
make build-uplink
make build-downlink
make build-gateway
```

## 9. 运行服务

### 开发模式（热重载）

```bash
# 运行单个服务
make run-api
make run-ws
make run-uplink
make run-downlink
make run-gateway

# 或使用 air 进行热重载
air -c .air.toml
```

### 生产模式

```bash
# 运行所有服务
make run-all

# 或使用 Docker Compose
docker-compose -f docker-compose.services.yml up
```

## 10. 验证服务

### 健康检查

```bash
# iot-api 健康检查
curl http://localhost:8080/health

# iot-ws 健康检查
curl http://localhost:8081/health
```

### 查看日志

```bash
# 查看服务日志
docker-compose logs -f iot-api

# 或直接查看应用日志
tail -f logs/iot-api.log
```

### 访问 Grafana

打开浏览器访问: http://localhost:3000

默认用户名/密码: admin/admin

## 11. 开发工作流

### 创建新功能

```bash
# 创建新的 feature branch
git checkout -b 002-feature-name

# 开发代码
# ...

# 运行测试和检查
make test
make lint

# 提交代码
git add .
git commit -m "feat: add new feature"
```

### 代码审查检查清单

- [ ] 代码通过 `make lint`
- [ ] 无拼写错误（`make spell-check`）
- [ ] 测试通过（`make test`）
- [ ] 测试覆盖率 ≥ 80%
- [ ] 遵循 Uber Go 规范
- [ ] 遵循命名规范
- [ ] 添加必要的注释和文档

## 12. 常见问题

### 端口冲突

如果端口被占用，修改 `.env` 文件中的端口配置。

### 数据库连接失败

检查 PostgreSQL 是否启动：
```bash
docker-compose ps postgres
```

### RabbitMQ 连接失败

检查 RabbitMQ 是否启动：
```bash
docker-compose ps rabbitmq
```

访问 RabbitMQ 管理界面: http://localhost:15672
默认用户名/密码: guest/guest

## 13. 下一步

- 阅读 [架构文档](../../docs/architecture/microservice-architecture.md)
- 阅读 [API 文档](../../docs/api/)
- 查看 [开发指南](../../docs/development/)

## 有用的命令

```bash
# 查看所有 Make 命令
make help

# 清理构建产物
make clean

# 格式化代码
make fmt

# 运行所有检查
make check

# 查看服务状态
make status
```

