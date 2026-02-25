# Research: Project Setup

## Research Questions

### 1. Go 版本选择

**Decision**: Go 1.22

**Rationale**: 
- Go 1.22 是当前稳定版本，提供了良好的性能和稳定性
- 支持最新的语言特性和标准库改进
- 与 Gin Framework 和 GORM 兼容性良好
- 社区支持广泛

**Alternatives considered**:
- Go 1.21: 稳定但功能较旧
- Go 1.23: 可能包含实验性特性，稳定性待验证

### 2. 仓库结构选择

**Decision**: Monorepo（单一仓库）

**Rationale**:
- 5个微服务共享大量公共代码（models、repository、shared utilities）
- 便于统一版本管理和依赖管理
- 简化代码审查和跨服务重构
- 便于统一 CI/CD 配置
- 符合微服务架构但共享代码库的最佳实践

**Alternatives considered**:
- 多仓库（每个服务独立仓库）: 增加管理复杂度，不利于代码复用

### 3. CI/CD 平台选择

**Decision**: GitHub Actions

**Rationale**:
- 与 GitHub 集成良好
- 配置简单，易于维护
- 支持矩阵构建（5个服务并行构建）
- 丰富的 Actions 市场插件
- 免费额度满足中小团队需求

**Alternatives considered**:
- GitLab CI: 功能强大但配置复杂
- Jenkins: 需要自建服务器，维护成本高

### 4. 配置管理方案

**Decision**: 环境变量 + 配置文件（YAML）

**Rationale**:
- 环境变量用于敏感信息（密码、密钥）
- YAML 配置文件用于非敏感配置（端口、超时时间等）
- 支持多环境（dev、staging、prod）
- 符合 12-Factor App 原则
- 未来可扩展为配置中心（如 Consul、etcd）

**Alternatives considered**:
- 纯环境变量: 配置项过多时难以管理
- 配置中心: 初期复杂度高，后续可迁移

### 5. 日志方案

**Decision**: 结构化日志（JSON 格式）+ Loki

**Rationale**:
- 结构化日志便于解析和查询
- JSON 格式与 Loki 集成良好
- 包含 trace_id、service、level、message、timestamp 等标准字段
- 符合可观测性原则

### 6. 指标监控方案

**Decision**: Prometheus + Grafana

**Rationale**:
- Prometheus 是标准的指标监控工具
- 与 Go 应用集成简单（prometheus/client_golang）
- Grafana 提供丰富的可视化能力
- 符合宪法规定的中间件选型

### 7. 分布式追踪方案

**Decision**: OpenTelemetry + Tempo

**Rationale**:
- OpenTelemetry 是行业标准
- Go 有良好的 OpenTelemetry SDK 支持
- Tempo 与 Prometheus/Loki 集成良好
- 符合宪法规定的中间件选型

### 8. 代码质量检查工具

**Decision**: golangci-lint + misspell

**Rationale**:
- golangci-lint 是 Go 社区最流行的 linter
- 支持 Uber Go 规范检查
- 集成 misspell 检查器，防止拼写错误
- 可配置性强，支持自定义规则
- 符合宪法要求

### 9. 测试框架选择

**Decision**: Go testing + Testify + Mockery

**Rationale**:
- Go testing 是标准库，无需额外依赖
- Testify 提供丰富的断言和测试工具
- Mockery 自动生成 Mock 对象，符合 TDD 原则
- 与 Go 生态集成良好

### 10. 数据库迁移方案

**Decision**: GORM Migrator + 手动迁移脚本

**Rationale**:
- GORM Migrator 提供基本的迁移能力
- 复杂迁移使用手动 SQL 脚本
- 版本控制迁移文件
- 符合 GORM 最佳实践

## Technology Stack Summary

| 组件 | 技术选型 | 版本 | 用途 |
|------|---------|------|------|
| 编程语言 | Go | 1.22 | 所有服务 |
| Web 框架 | Gin | Latest | HTTP/WebSocket |
| ORM | GORM | Latest | 数据库访问 |
| 消息队列 | RabbitMQ | Latest | 服务间通信 |
| MQTT Broker | VerneMQ | Latest | 设备连接 |
| 关系数据库 | PostgreSQL | 14+ | 业务数据 |
| 时序数据库 | InfluxDB | 2.x | 时序数据 |
| 指标监控 | Prometheus | Latest | 指标收集 |
| 日志聚合 | Loki | Latest | 日志查询 |
| 分布式追踪 | Tempo | Latest | 链路追踪 |
| 可视化 | Grafana | Latest | 统一展示 |
| 代码检查 | golangci-lint | Latest | 代码质量 |
| 拼写检查 | misspell | Latest | Typo 检查 |
| 测试框架 | Testify | Latest | 测试断言 |
| Mock 工具 | Mockery | Latest | Mock 生成 |

## Best Practices

### Go 项目结构
- 遵循 Go 标准项目布局（Standard Go Project Layout）
- cmd/ 目录存放各服务入口
- internal/ 目录存放内部包（不对外暴露）
- pkg/ 目录存放可对外暴露的公共包

### Gin Framework
- 使用中间件处理认证、日志、追踪
- 路由按功能模块组织
- 统一错误处理和响应格式

### GORM
- 模型定义在 pkg/models/
- Repository 模式封装数据访问
- 使用事务处理复杂操作
- 避免 N+1 查询问题

### RabbitMQ
- 统一的客户端封装（internal/shared/rabbitmq）
- 消息格式标准化（包含 tid、bid、timestamp）
- 实现消息幂等性
- 支持消息重试和死信队列

### 可观测性
- 结构化日志（JSON 格式）
- 统一的指标命名规范
- 分布式追踪上下文传递
- 健康检查端点（/health, /ready）

