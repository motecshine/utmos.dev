<!--
Sync Impact Report:
Version change: 1.4.0 → 1.4.1 (Enhanced code standards with naming conventions and typo prevention)
Modified principles: Language & Code Standards (expanded with detailed naming conventions and typo prevention)
Added sections: None
Removed sections: None
Changes: Added detailed Go naming conventions (package, variable, function, constant, type, interface, file naming). Added mandatory typo prevention requirement - no typos allowed in code, comments, documentation, variable names, function names, file names. Must use spell-checking tools (misspell) in CI/CD. Code review must check for typos.
Templates requiring updates:
  ✅ plan-template.md - Constitution Check section aligns with principles
  ✅ spec-template.md - Requirements section aligns with IoT platform principles
  ✅ tasks-template.md - Task categorization aligns with device abstraction principles
Follow-up TODOs: None
-->

# UMOS IoT Platform Constitution

## Core Principles

### I. Thing Model Driven Architecture (NON-NEGOTIABLE)
所有 IoT 设备必须通过标准化的物模型（Thing Model）进行抽象和描述。物模型使用 TSL (Thing Specification Language) JSON 格式，包含三个维度：Properties（属性）、Services（服务）、Events（事件）。平台必须能够解析、验证和执行物模型定义的能力。新设备接入时，只需提供物模型文件即可自动获得平台支持，无需修改核心代码。

**Rationale**: 物模型驱动架构实现了设备能力的标准化抽象，屏蔽了不同厂商设备的内部实现细节，使得平台能够统一管理各种 IoT 设备。这大大降低了新设备接入的成本，提高了平台的可扩展性和一致性。

### II. Multi-Protocol Support
平台必须同时支持 MQTT 5.0、HTTPS (RESTful API) 和 WebSocket 三种标准通信协议。MQTT 用于设备与平台之间的实时消息传输（属性上报、服务调用、事件推送），HTTPS 用于业务逻辑的 RESTful API 调用，WebSocket 用于服务器向客户端推送实时消息。所有协议必须遵循统一的 Topic/Endpoint 命名规范和消息格式规范。

**Rationale**: 不同场景需要不同的通信协议。MQTT 适合设备端的低带宽、高并发场景，HTTPS 适合业务逻辑的标准化接口，WebSocket 适合实时推送。多协议支持确保了平台能够适应各种业务场景的需求。

### III. Device Abstraction Layer
平台必须提供统一的设备抽象层（Device Abstraction Layer），将不同厂商、不同类型的设备映射到统一的设备接口。设备抽象层必须支持网关设备模式（Gateway Device Pattern），即设备可以通过网关设备连接到平台，网关设备负责管理子设备。设备抽象层必须支持动态设备注册、设备拓扑管理、设备能力发现等功能。

**Rationale**: 设备抽象层实现了设备差异的屏蔽，使得上层业务逻辑无需关心具体设备的实现细节。网关设备模式支持了复杂的设备组网场景。统一的设备接口使得平台能够以一致的方式管理各种设备。

### IV. Extensibility & Plugin Architecture
平台必须采用插件化架构，支持通过插件方式添加新的设备类型、新的功能模块、新的协议适配器。插件必须遵循标准的插件接口规范，能够独立开发、测试、部署。平台必须提供插件注册机制、插件生命周期管理、插件间通信机制。核心平台代码必须保持稳定，新功能通过插件方式扩展。

**Rationale**: 插件化架构确保了平台的核心稳定性，同时提供了灵活的扩展能力。新设备类型、新功能模块可以通过插件方式快速集成，无需修改核心代码。这大大提高了平台的灵活性和可维护性。

### V. Standardized API Design
所有 API 接口必须遵循统一的接口规范，包括统一的 URI 格式、统一的请求/响应格式、统一的错误码规范、统一的认证授权机制。API 版本管理必须采用 URL 路径版本控制（如 `/api/v1/`），支持多版本共存。所有 API 必须提供完整的 OpenAPI/Swagger 文档。

**Rationale**: 标准化的 API 设计降低了开发者的学习成本，提高了 API 的一致性和可维护性。版本管理确保了 API 的向后兼容性，支持平滑升级。完整的 API 文档提高了开发效率。

### VI. Test-First Development (NON-NEGOTIABLE)
所有功能开发必须遵循测试驱动开发（TDD）原则：先编写测试用例（包括单元测试、集成测试、契约测试），确保测试失败，然后实现功能，最后重构。设备接入必须提供模拟设备（Mock Device）用于测试。所有 API 接口必须提供契约测试（Contract Test）确保接口规范的一致性。

**Rationale**: 测试优先确保了代码质量和功能的正确性。模拟设备使得开发者可以在没有真实设备的情况下进行开发和测试。契约测试确保了不同实现之间的接口一致性。

### VII. Observability & Monitoring
平台必须提供完整的可观测性支持，包括结构化日志、分布式追踪、指标监控、告警机制。所有关键操作必须记录日志，日志必须包含请求 ID（Request ID）、设备 ID、时间戳等关键信息。平台必须支持实时监控设备状态、API 性能、系统资源使用情况。异常情况必须能够及时告警。

**Rationale**: 可观测性是生产环境系统的基础要求。结构化日志和分布式追踪帮助快速定位问题。指标监控和告警机制确保系统稳定运行。这对于 IoT 平台尤其重要，因为设备故障可能影响业务连续性。

## Architecture Constraints

### Protocol Standards
- **MQTT**: 必须支持 MQTT 5.0 协议，Topic 命名必须遵循分层结构，消息格式必须遵循统一的 JSON 结构，包含事务标识、业务标识、时间戳、数据载荷等必要字段
- **HTTPS**: 必须支持 RESTful API，URI 格式必须遵循统一的命名规范，请求/响应格式为 JSON，必须支持标准 HTTP 方法（GET、POST、PUT、DELETE）
- **WebSocket**: 必须支持 WebSocket 协议，消息格式为 JSON，必须包含业务代码、版本号、时间戳、数据载荷等必要字段

### Device Model Standards
- 所有设备必须提供物模型文件（TSL JSON 格式）
- 物模型必须包含产品信息、属性定义、服务定义、事件定义等必要组成部分
- 属性必须支持多种推送模式（定时上报、变化上报、被动查询等）
- 服务必须支持同步和异步两种调用模式
- 事件必须支持是否需要回复的标识机制

### Data Model Standards
- 所有消息必须包含事务标识和业务标识，用于消息追踪和业务关联
- 所有时间戳必须使用标准的时间戳格式（推荐毫秒级 Unix 时间戳）
- 所有设备标识必须使用统一的设备标识规范，支持网关设备和子设备的标识
- 所有错误响应必须包含错误码和错误信息，遵循统一的错误处理规范

### Language & Code Standards
- **Go 语言规范**: 所有 Go 代码必须严格遵循 [Uber Go 语言编码规范](https://github.com/uber-go/guide)。该规范涵盖了命名约定、错误处理、并发模式、接口设计、性能优化等方面的最佳实践
- **命名规范** (NON-NEGOTIABLE): 所有代码必须严格遵循 Go 官方命名规范和 Uber Go 规范：
  - **包命名**: 小写字母，简短有意义，不使用下划线或驼峰，如 `iotapi`、`device`、`thingmodel`
  - **变量命名**: 驼峰命名法（camelCase），首字母小写，如 `deviceID`、`thingModel`；私有变量首字母小写，公开变量首字母大写
  - **函数命名**: 驼峰命名法，公开函数首字母大写，私有函数首字母小写，如 `GetDevice()`、`processMessage()`
  - **常量命名**: 全大写字母，使用下划线分隔，如 `MAX_RETRY_COUNT`、`DEFAULT_TIMEOUT`
  - **类型命名**: 驼峰命名法，首字母大写，如 `Device`、`ThingModel`、`MessageHandler`
  - **接口命名**: 通常以 `-er` 结尾或使用描述性名称，如 `Reader`、`Writer`、`DeviceManager`
  - **文件名**: 小写字母，使用下划线分隔，如 `device_manager.go`、`thing_model.go`
- **代码质量** (NON-NEGOTIABLE): 
  - **严禁拼写错误（Typo）**: 所有代码、注释、文档、变量名、函数名、文件名等严禁出现拼写错误。代码审查必须检查拼写错误
  - 必须使用静态代码分析工具（如 `golangci-lint`）进行代码风格检查，确保代码符合规范要求
  - 必须启用拼写检查工具（如 `misspell`、`golangci-lint` 的 `misspell` 检查器）在 CI/CD 流程中自动检测拼写错误
  - 代码审查必须验证是否符合 Uber Go 规范，包括但不限于：包命名、变量命名、函数命名、错误处理方式、并发安全、接口设计原则、拼写正确性等

### Technology Stack Standards (NON-NEGOTIABLE)
平台的基础技术栈选型必须严格遵循以下规定，严禁私自引入其他技术栈组件。如需添加新的技术栈组件，必须经过技术委员会讨论和批准。

**已批准的基础技术栈清单**：
- **编程语言**: Go (Golang) - 所有服务必须使用 Go 语言开发
- **Web 框架**: Gin Framework - HTTP API 服务和 WebSocket 服务必须使用 Gin Framework
- **ORM 框架**: GORM - 数据库访问必须使用 GORM 框架

**技术栈使用原则**：
- **Go 语言**: 所有微服务必须使用 Go 语言开发，版本要求遵循项目统一规定（建议 Go 1.21+）
- **Gin Framework**: 
  - iot-api 服务必须使用 Gin Framework 处理 HTTP RESTful API
  - iot-ws 服务必须使用 Gin Framework 处理 WebSocket 连接
  - 其他服务如需 HTTP 接口（如健康检查），也应使用 Gin Framework
  - 必须遵循 Gin 的最佳实践，包括中间件使用、路由组织、错误处理等
- **GORM**: 
  - 所有数据库访问操作必须使用 GORM 框架
  - PostgreSQL 和 InfluxDB 的数据访问都应通过 GORM（InfluxDB 可能需要适配器）
  - 必须遵循 GORM 的最佳实践，包括模型定义、迁移管理、查询优化、事务处理等
  - 严禁使用原生 SQL 或直接数据库驱动，除非 GORM 无法满足特定需求（需技术委员会批准）

**技术栈变更流程**：
1. 提出新技术栈组件需求，说明使用场景和必要性
2. 技术委员会评估：性能、稳定性、社区支持、学习成本、与现有技术栈的集成、团队技能匹配度
3. 技术委员会投票决定是否批准
4. 批准后更新宪法文档，添加新技术栈组件到已批准清单
5. 更新架构文档、开发文档和示例代码

**Rationale**: 统一的技术栈选型确保了代码的一致性和可维护性，降低了开发者的学习成本，提高了代码复用性。Go 语言的高性能和并发特性适合 IoT 平台的高并发场景。Gin Framework 提供了轻量级、高性能的 HTTP 框架。GORM 提供了类型安全的数据库访问，减少了 SQL 注入风险，提高了开发效率。严格的变更流程避免了技术债务的积累，确保每个技术栈组件的引入都经过充分论证。

### Microservice Architecture (NON-NEGOTIABLE)
平台必须采用微服务架构，拆分为以下5个核心服务，服务之间通过 RabbitMQ 进行异步消息通信。每个服务必须职责单一、边界清晰，严禁跨服务直接调用。

**核心服务清单及职责**：

1. **iot-api** (HTTP API 服务)
   - **职责**: 处理所有 HTTP RESTful API 请求
   - **功能**: 设备管理 API、物模型管理 API、业务逻辑 API、用户认证授权
   - **通信**: 接收 HTTP 请求，通过 RabbitMQ 发布消息到其他服务，接收 RabbitMQ 消息并返回 HTTP 响应
   - **数据访问**: 直接访问 PostgreSQL（业务数据）、InfluxDB（时序数据查询）

2. **iot-ws** (WebSocket 服务)
   - **职责**: 处理所有 WebSocket 连接和实时消息推送
   - **功能**: WebSocket 连接管理、向客户端推送实时消息（设备状态变化、事件通知、服务响应等）、接收客户端通过 WebSocket 发送的消息
   - **通信**: 维护 WebSocket 连接，订阅 RabbitMQ 消息队列接收需要推送的消息，通过 RabbitMQ 发布客户端消息
   - **数据访问**: 不直接访问数据库，所有数据通过 RabbitMQ 消息获取

3. **iot-uplink** (上行消息处理服务)
   - **职责**: 处理从设备到平台的所有上行消息
   - **功能**: 从 RabbitMQ 接收上行消息（属性上报、事件上报、服务响应），消息解析和验证，物模型映射，业务逻辑处理，消息路由到其他服务
   - **通信**: 订阅 RabbitMQ 上行消息队列（从 iot-gateway 接收），处理后发布到 RabbitMQ 相应队列（供其他服务消费）
   - **数据访问**: 读取物模型定义（PostgreSQL），写入时序数据（InfluxDB）
   - **重要约束**: 严禁直接连接 VerneMQ 或任何 MQTT Broker，所有上行消息必须通过 iot-gateway 和 RabbitMQ 接收

4. **iot-downlink** (下行消息处理服务)
   - **职责**: 处理从平台到设备的所有下行消息
   - **功能**: 从 RabbitMQ 接收下行消息（服务调用、属性设置、命令下发），消息格式转换和验证，消息确认和重试，路由到 RabbitMQ（发送给 iot-gateway）
   - **通信**: 订阅 RabbitMQ 下行消息队列（从其他服务接收），发布消息到 RabbitMQ（发送给 iot-gateway）
   - **数据访问**: 读取物模型定义（PostgreSQL），记录消息状态（PostgreSQL）
   - **重要约束**: 严禁直接连接 VerneMQ 或任何 MQTT Broker，所有下行消息必须通过 RabbitMQ 发送给 iot-gateway

5. **iot-gateway** (协议网关服务)
   - **职责**: 协议适配和网关管理，是唯一直接连接 MQTT Broker 的服务
   - **功能**: MQTT 协议适配（连接 VerneMQ），设备认证和授权，设备连接管理，协议转换（如需要支持其他协议），设备拓扑管理，MQTT 消息与 RabbitMQ 消息的双向转换
   - **通信**: 
     - **上行**: 从 VerneMQ 订阅 MQTT Topic 接收设备消息，转换为标准格式后发布到 RabbitMQ 上行队列（供 iot-uplink 消费）
     - **下行**: 从 RabbitMQ 订阅下行消息队列（从 iot-downlink 接收），转换为 MQTT 格式后发布到 VerneMQ MQTT Topic
   - **数据访问**: 读取设备配置和认证信息（PostgreSQL），更新设备连接状态（PostgreSQL）
   - **重要约束**: 这是唯一允许直接连接 VerneMQ 的服务，负责所有 MQTT 协议相关的操作

**服务间通信原则**：
- 所有服务间通信必须通过 RabbitMQ 进行，严禁服务间直接调用（HTTP/gRPC 等）
- **MQTT 协议隔离**: 只有 iot-gateway 服务可以直接连接 VerneMQ，其他服务（iot-uplink、iot-downlink）严禁直接连接 MQTT Broker，必须通过 iot-gateway 和 RabbitMQ 进行消息传递
- RabbitMQ 消息必须包含事务标识（tid）、业务标识（bid）、时间戳等标准字段
- 消息队列命名必须遵循统一规范：`iot.{vendor}.{service}.{action}`（如 `iot.dji.uplink.property.report`），支持多厂商路由
- 服务必须实现消息幂等性处理，支持消息重试和死信队列
- 服务必须实现分布式追踪，在消息中传递追踪上下文

**服务部署原则**：
- 每个服务必须可以独立部署和扩展
- 服务必须实现健康检查接口，支持 Kubernetes/Docker 编排
- 服务必须支持配置外部化，通过环境变量或配置中心管理
- 服务必须实现优雅关闭，确保消息处理完成后再退出

**Rationale**: 微服务架构实现了职责分离，提高了系统的可扩展性和可维护性。通过 RabbitMQ 进行异步通信，实现了服务间的解耦，提高了系统的容错能力。将 MQTT 协议隔离到 iot-gateway 服务，使得其他服务无需关心 MQTT 协议细节，只需处理业务逻辑，进一步提高了系统的可维护性和可测试性。明确的职责划分使得每个服务可以独立开发、测试、部署和扩展。

### Infrastructure & Middleware Standards (NON-NEGOTIABLE)
平台的基础中间件选型必须严格遵循以下规定，严禁私自引入其他中间件。如需添加新的中间件，必须经过技术委员会讨论和批准。

**已批准的基础中间件清单**：
- **消息队列**: VerneMQ (MQTT Broker), RabbitMQ (AMQP/消息队列)
- **数据库**: PostgreSQL (关系型数据库), InfluxDB (时序数据库)
- **可观测性**: Prometheus (指标监控), Loki (日志聚合), Tempo (分布式追踪), Grafana (可视化)

**中间件使用原则**：
- VerneMQ 用于 MQTT 5.0 协议的消息代理，处理设备与平台之间的实时消息传输
- RabbitMQ 用于应用内部的消息队列，处理异步任务和事件驱动架构
- PostgreSQL 用于存储业务数据、设备元数据、物模型定义等结构化数据
- InfluxDB 用于存储设备时序数据（如传感器数据、设备状态历史等）
- Prometheus 用于收集和存储系统指标、应用指标、设备指标
- Loki 用于聚合和查询应用日志、系统日志
- Tempo 用于分布式追踪，追踪请求在微服务间的调用链路
- Grafana 用于统一的可视化展示，包括指标监控、日志查询、追踪分析

**中间件变更流程**：
1. 提出新中间件需求，说明使用场景和必要性
2. 技术委员会评估：性能、稳定性、社区支持、运维成本、与现有中间件的集成
3. 技术委员会投票决定是否批准
4. 批准后更新宪法文档，添加新中间件到已批准清单
5. 更新架构文档和部署文档

**Rationale**: 统一的基础中间件选型确保了技术栈的一致性，降低了运维复杂度，提高了系统的可维护性。严格的变更流程避免了技术债务的积累，确保每个中间件的引入都经过充分论证。

## Development Workflow

### Device Integration Process
1. **物模型定义**: 新设备接入前，必须先定义物模型文件，明确设备的属性、服务、事件
2. **协议适配器开发**: 根据设备使用的通信协议，开发对应的协议适配器
3. **设备抽象层实现**: 实现设备抽象层接口，将设备特定的操作映射到统一接口
4. **测试验证**: 使用模拟设备进行单元测试、集成测试、端到端测试
5. **文档编写**: 编写设备接入文档、API 文档、使用示例

### Code Review Requirements
- 所有代码变更必须经过代码审查
- 代码审查必须验证是否符合宪法原则
- 新设备接入必须验证物模型文件的完整性和正确性
- API 变更必须验证版本管理和向后兼容性
- 性能关键代码必须进行性能测试

### Quality Gates
- 所有代码必须通过静态代码分析（Linting），Go 代码必须通过 `golangci-lint` 检查并符合 Uber Go 规范
- 所有测试必须通过（单元测试覆盖率 ≥ 80%）
- 所有 API 必须提供完整的 OpenAPI 文档
- 所有设备接入必须提供使用示例和测试用例
- 性能关键路径必须满足性能指标要求
- 代码审查必须验证是否符合 Uber Go 编码规范

## Governance

本宪法是 UMOS IoT 平台开发的所有实践的最高指导原则。所有开发活动必须遵循本宪法的原则和约束。任何违反宪法的实践必须经过充分论证和批准。

**Amendment Process**: 宪法修订必须经过以下流程：
1. 提出修订提案，说明修订原因和影响范围
2. 技术委员会审查和讨论
3. 更新宪法文档，更新版本号（遵循语义化版本控制）
4. 更新所有相关模板和文档
5. 通知所有开发人员

**Compliance Review**: 所有 Pull Request 必须验证是否符合宪法原则。代码审查必须包含宪法合规性检查。定期（每季度）进行全面的合规性审查。

**Version Control**: 宪法版本遵循语义化版本控制（MAJOR.MINOR.PATCH）：
- MAJOR: 向后不兼容的原则变更或原则删除
- MINOR: 新增原则或重大原则扩展
- PATCH: 澄清、措辞修正、非语义性改进

**Version**: 1.4.1 | **Ratified**: 2025-01-27 | **Last Amended**: 2025-01-27
