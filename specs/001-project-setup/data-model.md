# Data Model: Project Setup

本文档定义了 UMOS IoT 平台的基础数据模型，这些模型将在项目初始化时创建。

## 核心实体

### 1. Device (设备)

**用途**: 存储设备基本信息

**字段**:
- `id` (uint, Primary Key): 设备ID
- `device_sn` (string, Unique): 设备序列号
- `device_name` (string): 设备名称
- `device_type` (string): 设备类型
- `vendor` (string, Index): 厂商标识（dji, generic, tuya 等），用于 RabbitMQ routing key 生成
- `gateway_sn` (string, Nullable): 网关设备序列号（如果是子设备）
- `thing_model_id` (uint, Foreign Key): 关联的物模型ID
- `status` (enum): 设备状态（online, offline, unknown）
- `last_online_time` (timestamp, Nullable): 最后在线时间
- `created_at` (timestamp): 创建时间
- `updated_at` (timestamp): 更新时间

**索引**:
- `device_sn` (Unique Index)
- `vendor` (Index)
- `gateway_sn` (Index)
- `thing_model_id` (Index)

### 2. ThingModel (物模型)

**用途**: 存储物模型定义（TSL JSON）

**字段**:
- `id` (uint, Primary Key): 物模型ID
- `product_key` (string, Unique): 产品标识
- `product_name` (string): 产品名称
- `version` (string): 物模型版本
- `tsl_json` (jsonb): 物模型定义（TSL JSON格式）
- `description` (text, Nullable): 描述
- `created_at` (timestamp): 创建时间
- `updated_at` (timestamp): 更新时间

**索引**:
- `product_key` (Unique Index)

### 3. DeviceProperty (设备属性)

**用途**: 存储设备属性值（当前值）

**字段**:
- `id` (uint, Primary Key): 记录ID
- `device_id` (uint, Foreign Key): 设备ID
- `property_key` (string): 属性标识
- `property_value` (jsonb): 属性值（JSON格式）
- `updated_at` (timestamp): 更新时间

**索引**:
- `device_id` + `property_key` (Composite Unique Index)

### 4. DeviceEvent (设备事件)

**用途**: 存储设备事件记录

**字段**:
- `id` (uint, Primary Key): 事件ID
- `device_id` (uint, Foreign Key): 设备ID
- `event_key` (string): 事件标识
- `event_data` (jsonb): 事件数据（JSON格式）
- `timestamp` (timestamp): 事件时间
- `created_at` (timestamp): 创建时间

**索引**:
- `device_id` + `timestamp` (Composite Index)
- `event_key` (Index)

### 5. MessageLog (消息日志)

**用途**: 记录消息处理日志（用于调试和追踪）

**字段**:
- `id` (uint, Primary Key): 日志ID
- `tid` (string, Index): 事务ID
- `bid` (string, Index): 业务ID
- `service` (string): 服务名称
- `direction` (enum): 消息方向（uplink, downlink）
- `message_type` (string): 消息类型
- `device_sn` (string, Index): 设备序列号
- `message_data` (jsonb): 消息数据（JSON格式）
- `status` (enum): 处理状态（success, failed, pending）
- `error_message` (text, Nullable): 错误信息
- `created_at` (timestamp): 创建时间

**索引**:
- `tid` (Index)
- `bid` (Index)
- `device_sn` + `created_at` (Composite Index)

## 关系说明

```
ThingModel (1) ──< (N) Device
Device (1) ──< (N) DeviceProperty
Device (1) ──< (N) DeviceEvent
Device (1) ──< (N) MessageLog
```

## 数据库初始化

**重要说明**: 根据 FR-015，系统必须使用 GORM 的 AutoMigrate 功能进行数据库迁移，不得使用手动 SQL 脚本。以下 SQL 脚本仅作为数据库结构的参考文档，实际迁移将通过 GORM AutoMigrate 实现。

### PostgreSQL Schema (参考)

```sql
-- 创建数据库
CREATE DATABASE umos_iot;

-- 创建扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";  -- 用于全文搜索

-- 物模型表
CREATE TABLE thing_models (
    id SERIAL PRIMARY KEY,
    product_key VARCHAR(100) UNIQUE NOT NULL,
    product_name VARCHAR(200) NOT NULL,
    version VARCHAR(50) NOT NULL,
    tsl_json JSONB NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_thing_models_product_key ON thing_models(product_key);

-- 设备表
CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    device_sn VARCHAR(100) UNIQUE NOT NULL,
    device_name VARCHAR(200) NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    vendor VARCHAR(50) NOT NULL DEFAULT 'generic',
    gateway_sn VARCHAR(100),
    thing_model_id INTEGER REFERENCES thing_models(id),
    status VARCHAR(20) DEFAULT 'unknown',
    last_online_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_devices_device_sn ON devices(device_sn);
CREATE INDEX idx_devices_vendor ON devices(vendor);
CREATE INDEX idx_devices_gateway_sn ON devices(gateway_sn);
CREATE INDEX idx_devices_thing_model_id ON devices(thing_model_id);

-- 设备属性表
CREATE TABLE device_properties (
    id SERIAL PRIMARY KEY,
    device_id INTEGER REFERENCES devices(id) ON DELETE CASCADE,
    property_key VARCHAR(100) NOT NULL,
    property_value JSONB NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(device_id, property_key)
);

CREATE INDEX idx_device_properties_device_id ON device_properties(device_id);

-- 设备事件表
CREATE TABLE device_events (
    id SERIAL PRIMARY KEY,
    device_id INTEGER REFERENCES devices(id) ON DELETE CASCADE,
    event_key VARCHAR(100) NOT NULL,
    event_data JSONB NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_device_events_device_id_timestamp ON device_events(device_id, timestamp DESC);
CREATE INDEX idx_device_events_event_key ON device_events(event_key);

-- 消息日志表
CREATE TABLE message_logs (
    id SERIAL PRIMARY KEY,
    tid VARCHAR(100),
    bid VARCHAR(100),
    service VARCHAR(50) NOT NULL,
    direction VARCHAR(20) NOT NULL,
    message_type VARCHAR(100) NOT NULL,
    device_sn VARCHAR(100),
    message_data JSONB NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_message_logs_tid ON message_logs(tid);
CREATE INDEX idx_message_logs_bid ON message_logs(bid);
CREATE INDEX idx_message_logs_device_sn_created_at ON message_logs(device_sn, created_at DESC);
```

### GORM Models

模型定义将在 `pkg/models/` 目录下创建，使用 GORM 标签进行映射。所有数据库迁移将通过 GORM 的 `AutoMigrate` 功能实现（参见 FR-015 和 T018）。

## 时序数据（InfluxDB）

时序数据存储在 InfluxDB 中，包括：
- 设备属性历史值
- 设备性能指标
- 系统监控指标

具体的 measurement 和 tag 设计将在后续功能开发中定义。

