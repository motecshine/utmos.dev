# Quickstart: Core Services Implementation

**Feature**: 004-core-services-implementation
**Date**: 2025-02-05

## Prerequisites

- Go 1.22+
- Docker & Docker Compose
- 已完成 001/002/003 的基础设施部署

## Quick Start

### 1. 启动基础设施

```bash
# 启动所有依赖服务
docker-compose up -d

# 验证服务状态
docker-compose ps
```

预期输出:
```
NAME                STATUS
postgres            running (healthy)
influxdb            running (healthy)
rabbitmq            running (healthy)
vernemq             running (healthy)
```

### 2. 配置环境变量

```bash
# 复制示例配置
cp configs/iot-gateway.example.yaml configs/iot-gateway.yaml
cp configs/iot-uplink.example.yaml configs/iot-uplink.yaml
cp configs/iot-downlink.example.yaml configs/iot-downlink.yaml
cp configs/iot-api.example.yaml configs/iot-api.yaml
cp configs/iot-ws.example.yaml configs/iot-ws.yaml

# 或使用环境变量
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=umos
export POSTGRES_PASSWORD=umos
export POSTGRES_DB=umos

export RABBITMQ_HOST=localhost
export RABBITMQ_PORT=5672
export RABBITMQ_USER=guest
export RABBITMQ_PASSWORD=guest

export INFLUXDB_URL=http://localhost:8086
export INFLUXDB_TOKEN=your-token
export INFLUXDB_ORG=umos
export INFLUXDB_BUCKET=iot

export VERNEMQ_HOST=localhost
export VERNEMQ_PORT=1883
```

### 3. 运行服务

```bash
# 终端 1: 启动 iot-gateway
make run-gateway

# 终端 2: 启动 iot-uplink
make run-uplink

# 终端 3: 启动 iot-downlink
make run-downlink

# 终端 4: 启动 iot-api
make run-api

# 终端 5: 启动 iot-ws
make run-ws
```

### 4. 验证服务健康

```bash
# 检查各服务健康状态
curl http://localhost:8080/health  # iot-api
curl http://localhost:8081/health  # iot-ws
curl http://localhost:8082/health  # iot-gateway
curl http://localhost:8083/health  # iot-uplink
curl http://localhost:8084/health  # iot-downlink
```

## Integration Scenarios

### Scenario 1: 设备上线

**步骤**:

1. 设备连接 VerneMQ (MQTT)
2. iot-gateway 验证设备凭证
3. 设备发布 status 消息
4. iot-uplink 处理上线事件
5. iot-ws 推送到订阅客户端

**测试命令**:

```bash
# 使用 mosquitto_pub 模拟设备
mosquitto_pub -h localhost -p 1883 \
  -u "device-001" -P "secret" \
  -t "sys/product/device-001/status" \
  -m '{"tid":"uuid-1","bid":"uuid-2","timestamp":1234567890,"data":{"online":true}}'
```

**预期结果**:
- iot-gateway 日志显示设备连接
- iot-uplink 日志显示消息处理
- WebSocket 客户端收到上线通知

### Scenario 2: 遥测数据上报

**步骤**:

1. 设备发布 OSD 消息到 VerneMQ
2. iot-gateway 转发到 RabbitMQ
3. iot-uplink 解析并写入 InfluxDB
4. iot-uplink 路由到 iot-ws
5. iot-ws 推送到订阅客户端

**测试命令**:

```bash
# 模拟 DJI 飞行器 OSD 数据
mosquitto_pub -h localhost -p 1883 \
  -u "device-001" -P "secret" \
  -t "thing/product/device-001/osd" \
  -m '{
    "tid":"uuid-1",
    "bid":"uuid-2",
    "timestamp":1234567890,
    "data":{
      "host":{
        "latitude":31.2304,
        "longitude":121.4737,
        "altitude":100.5,
        "battery":{"capacity_percent":85}
      }
    }
  }'
```

**验证 InfluxDB**:

```bash
# 查询遥测数据
influx query 'from(bucket:"iot") |> range(start:-1h) |> filter(fn:(r) => r._measurement == "dji_aircraft_osd")'
```

### Scenario 3: 服务调用

**步骤**:

1. 客户端调用 iot-api 服务调用接口
2. iot-api 发布消息到 RabbitMQ
3. iot-downlink 处理并路由到 iot-gateway
4. iot-gateway 发布 MQTT 消息到设备
5. 设备响应，原路返回

**测试命令**:

```bash
# 调用设备服务
curl -X POST http://localhost:8080/api/v1/devices/device-001/services/flighttask_prepare \
  -H "Content-Type: application/json" \
  -d '{
    "params": {
      "file_id": "wayline-001",
      "task_type": 0
    }
  }'
```

**预期响应**:

```json
{
  "tid": "uuid-generated",
  "bid": "uuid-generated",
  "status": "pending",
  "message": "Service call initiated"
}
```

### Scenario 4: WebSocket 实时推送

**步骤**:

1. 客户端建立 WebSocket 连接
2. 客户端订阅设备主题
3. 设备上报数据
4. 客户端收到实时推送

**测试代码**:

```javascript
// WebSocket 客户端示例
const ws = new WebSocket('ws://localhost:8081/ws');

ws.onopen = () => {
  // 订阅设备 OSD 数据
  ws.send(JSON.stringify({
    type: 'subscribe',
    topics: ['device.osd.device-001']
  }));
};

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  console.log('Received:', msg);
};
```

## API Reference

### Device Management

```bash
# 获取设备列表
GET /api/v1/devices

# 获取设备详情
GET /api/v1/devices/{sn}

# 创建设备
POST /api/v1/devices
{
  "sn": "device-001",
  "name": "Test Device",
  "type": "aircraft",
  "vendor": "dji"
}

# 更新设备
PUT /api/v1/devices/{sn}
{
  "name": "Updated Name"
}

# 删除设备
DELETE /api/v1/devices/{sn}
```

### Service Call

```bash
# 调用设备服务
POST /api/v1/devices/{sn}/services/{method}
{
  "params": {...}
}

# 查询服务调用状态
GET /api/v1/services/{tid}
```

### Telemetry Query

```bash
# 查询遥测数据
GET /api/v1/devices/{sn}/telemetry?start=-1h&end=now

# 查询最新遥测
GET /api/v1/devices/{sn}/telemetry/latest
```

## Troubleshooting

### 问题: MQTT 连接失败

**症状**: 设备无法连接 VerneMQ

**解决方案**:
1. 检查 VerneMQ 是否运行: `docker-compose ps vernemq`
2. 检查设备凭证是否正确
3. 检查 iot-gateway 日志

### 问题: 消息未到达 iot-uplink

**症状**: 设备发送消息但 iot-uplink 未收到

**解决方案**:
1. 检查 RabbitMQ 连接: `curl http://localhost:15672/api/queues`
2. 检查 iot-gateway 日志是否有转发记录
3. 检查 RabbitMQ exchange 和 queue 绑定

### 问题: InfluxDB 写入失败

**症状**: 遥测数据未写入 InfluxDB

**解决方案**:
1. 检查 InfluxDB 连接: `curl http://localhost:8086/health`
2. 检查 token 和 bucket 配置
3. 检查 iot-uplink 日志

### 问题: WebSocket 连接断开

**症状**: WebSocket 连接频繁断开

**解决方案**:
1. 检查心跳配置
2. 检查网络稳定性
3. 检查 iot-ws 日志

## Performance Tuning

### RabbitMQ

```yaml
# 增加预取数量
rabbitmq:
  prefetch_count: 100

# 启用消息持久化
rabbitmq:
  durable: true
```

### InfluxDB

```yaml
# 批量写入配置
influxdb:
  batch_size: 1000
  flush_interval: 1s
```

### WebSocket

```yaml
# 连接限制
websocket:
  max_connections: 10000
  read_buffer_size: 1024
  write_buffer_size: 1024
```
