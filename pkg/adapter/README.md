# Protocol Adapter Development Guide

This guide explains how to implement a new protocol adapter for the UMOS IoT platform.

## Overview

Protocol adapters convert vendor-specific IoT messages to/from the platform's standard message format. Each vendor (DJI, Tuya, etc.) has its own adapter that handles protocol-specific parsing and conversion.

## Architecture

```
Device → VerneMQ → iot-gateway → RabbitMQ (iot.raw.{vendor}.uplink)
                                        ↓
                               {vendor}-adapter (parse & convert)
                                        ↓
                               RabbitMQ (iot.{vendor}.{service}.{action})
                                        ↓
                                   iot-uplink
```

## Interface Specification

All protocol adapters must implement the `ProtocolAdapter` interface:

```go
type ProtocolAdapter interface {
    // GetVendor returns the vendor identifier (e.g., "dji", "tuya").
    GetVendor() string

    // ParseRawMessage parses raw bytes into a protocol-specific message.
    // The topic parameter is the original MQTT topic.
    ParseRawMessage(topic string, payload []byte) (*ProtocolMessage, error)

    // ToStandardMessage converts a protocol message to a standard message.
    ToStandardMessage(pm *ProtocolMessage) (*rabbitmq.StandardMessage, error)

    // FromStandardMessage converts a standard message to a protocol message.
    // Used for downlink message conversion.
    FromStandardMessage(sm *rabbitmq.StandardMessage) (*ProtocolMessage, error)

    // GetRawPayload returns the raw payload bytes for sending to device.
    // Used for downlink message serialization.
    GetRawPayload(pm *ProtocolMessage) ([]byte, error)
}
```

## Step-by-Step Implementation Guide

### Step 1: Create Package Structure

Create a new package under `pkg/adapter/{vendor}/`:

```
pkg/adapter/{vendor}/
├── adapter.go      # Main adapter implementation
├── types.go        # Protocol-specific types
├── topic.go        # Topic parsing logic
├── parser.go       # Message parsing logic
├── converter.go    # Message conversion logic
├── errors.go       # Custom error definitions
└── *_test.go       # Test files
```

### Step 2: Define Protocol Types

Create `types.go` with protocol-specific structures:

```go
package myvendor

const VendorMyVendor = "myvendor"

// TopicType represents the type of MQTT topic.
type TopicType string

const (
    TopicTypeProperty TopicType = "property"
    TopicTypeEvent    TopicType = "event"
    TopicTypeService  TopicType = "service"
)

// TopicInfo contains parsed topic information.
type TopicInfo struct {
    Type      TopicType
    DeviceSN  string
    Raw       string
}

// VendorMessage represents the vendor's message format.
type VendorMessage struct {
    ID        string          `json:"id"`
    Timestamp int64           `json:"timestamp"`
    Data      json.RawMessage `json:"data"`
}
```

### Step 3: Implement Topic Parser

Create `topic.go` to parse MQTT topics:

```go
package myvendor

import (
    "fmt"
    "strings"
)

// ParseTopic parses a vendor-specific MQTT topic.
func ParseTopic(topic string) (*TopicInfo, error) {
    parts := strings.Split(topic, "/")
    if len(parts) < 3 {
        return nil, fmt.Errorf("invalid topic format: %s", topic)
    }

    return &TopicInfo{
        Type:     TopicType(parts[2]),
        DeviceSN: parts[1],
        Raw:      topic,
    }, nil
}
```

### Step 4: Implement Message Parser

Create `parser.go` to parse message payloads:

```go
package myvendor

import (
    "encoding/json"
    "errors"
)

var ErrEmptyPayload = errors.New("empty payload")

// ParseMessage parses raw bytes into a VendorMessage.
func ParseMessage(payload []byte) (*VendorMessage, error) {
    if len(payload) == 0 {
        return nil, ErrEmptyPayload
    }

    var msg VendorMessage
    if err := json.Unmarshal(payload, &msg); err != nil {
        return nil, err
    }

    return &msg, nil
}
```

### Step 5: Implement Converter

Create `converter.go` for message conversion:

```go
package myvendor

import (
    "github.com/utmos/utmos/pkg/rabbitmq"
)

type Converter struct{}

func NewConverter() *Converter {
    return &Converter{}
}

// ToStandardMessage converts vendor message to standard format.
func (c *Converter) ToStandardMessage(msg *VendorMessage, topic *TopicInfo) (*rabbitmq.StandardMessage, error) {
    action := mapTopicTypeToAction(topic.Type)

    return &rabbitmq.StandardMessage{
        DeviceSN:  topic.DeviceSN,
        TID:       msg.ID,
        Timestamp: msg.Timestamp,
        Action:    action,
        Data:      msg.Data,
        ProtocolMeta: &rabbitmq.ProtocolMeta{
            Vendor:        VendorMyVendor,
            OriginalTopic: topic.Raw,
        },
    }, nil
}

// FromStandardMessage converts standard message to vendor format.
func (c *Converter) FromStandardMessage(sm *rabbitmq.StandardMessage) (*VendorMessage, error) {
    return &VendorMessage{
        ID:        sm.TID,
        Timestamp: sm.Timestamp,
        Data:      sm.Data,
    }, nil
}

func mapTopicTypeToAction(tt TopicType) string {
    switch tt {
    case TopicTypeProperty:
        return "property.report"
    case TopicTypeEvent:
        return "event.report"
    case TopicTypeService:
        return "service.call"
    default:
        return "property.report"
    }
}
```

### Step 6: Implement Adapter

Create `adapter.go` as the main entry point:

```go
package myvendor

import (
    "github.com/utmos/utmos/pkg/adapter"
    "github.com/utmos/utmos/pkg/rabbitmq"
)

type Adapter struct {
    converter *Converter
}

func NewAdapter() *Adapter {
    return &Adapter{
        converter: NewConverter(),
    }
}

func (a *Adapter) GetVendor() string {
    return VendorMyVendor
}

func (a *Adapter) ParseRawMessage(topic string, payload []byte) (*adapter.ProtocolMessage, error) {
    topicInfo, err := ParseTopic(topic)
    if err != nil {
        return nil, err
    }

    msg, err := ParseMessage(payload)
    if err != nil {
        return nil, err
    }

    return &adapter.ProtocolMessage{
        Vendor:      VendorMyVendor,
        Topic:       topic,
        DeviceSN:    topicInfo.DeviceSN,
        MessageType: mapTopicTypeToMessageType(topicInfo.Type),
        TID:         msg.ID,
        Timestamp:   msg.Timestamp,
        Data:        msg.Data,
    }, nil
}

func (a *Adapter) ToStandardMessage(pm *adapter.ProtocolMessage) (*rabbitmq.StandardMessage, error) {
    // Implementation
}

func (a *Adapter) FromStandardMessage(sm *rabbitmq.StandardMessage) (*adapter.ProtocolMessage, error) {
    // Implementation
}

func (a *Adapter) GetRawPayload(pm *adapter.ProtocolMessage) ([]byte, error) {
    // Implementation
}

// Register registers the adapter with the global registry.
func Register() {
    adapter.Register(NewAdapter())
}

// Ensure Adapter implements ProtocolAdapter interface.
var _ adapter.ProtocolAdapter = (*Adapter)(nil)
```

### Step 7: Create Service Main

Create `cmd/{vendor}-adapter/main.go`:

```go
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"

    "github.com/utmos/utmos/internal/shared/config"
    "github.com/utmos/utmos/internal/shared/logger"
    "github.com/utmos/utmos/pkg/adapter"
    myvendor "github.com/utmos/utmos/pkg/adapter/myvendor"
    "github.com/utmos/utmos/pkg/rabbitmq"
)

func main() {
    cfg, _ := config.Load("myvendor-adapter")
    log := logger.New(&cfg.Logger)

    // Register adapter
    myvendor.Register()

    // Get adapter
    vendorAdapter, _ := adapter.Get(myvendor.VendorMyVendor)

    // Connect to RabbitMQ and process messages
    // ...

    // Wait for shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
}
```

## Testing Requirements

### Unit Tests

Each component should have comprehensive unit tests:

```go
// topic_test.go
func TestParseTopic(t *testing.T) {
    tests := []struct {
        name    string
        topic   string
        want    *TopicInfo
        wantErr bool
    }{
        {
            name:  "valid property topic",
            topic: "device/DEVICE001/property",
            want: &TopicInfo{
                Type:     TopicTypeProperty,
                DeviceSN: "DEVICE001",
            },
        },
        {
            name:    "invalid topic",
            topic:   "invalid",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseTopic(tt.topic)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.want.Type, got.Type)
            assert.Equal(t, tt.want.DeviceSN, got.DeviceSN)
        })
    }
}
```

### Integration Tests

Test the complete message flow:

```go
func TestAdapterMessageFlow(t *testing.T) {
    myvendor.Register()
    defer adapter.Unregister(myvendor.VendorMyVendor)

    vendorAdapter, _ := adapter.Get(myvendor.VendorMyVendor)

    // Test uplink flow
    pm, err := vendorAdapter.ParseRawMessage(topic, payload)
    require.NoError(t, err)

    stdMsg, err := vendorAdapter.ToStandardMessage(pm)
    require.NoError(t, err)
    assert.Equal(t, expectedDeviceSN, stdMsg.DeviceSN)
}
```

### Coverage Requirements

- Unit test coverage must be >= 80%
- Run tests with: `go test -cover ./pkg/adapter/myvendor/...`

## Using the Adapter Factory

The adapter factory provides a convenient way to create adapters:

```go
// Using the factory
factory := adapter.NewFactory()
vendorAdapter, err := factory.NewAdapter("myvendor")

// Or use global functions
vendorAdapter, err := adapter.NewAdapterByVendor("myvendor")

// Check available vendors
vendors := adapter.ListAvailableVendors()
available := adapter.IsVendorAvailable("myvendor")
```

## RabbitMQ Queue Conventions

- **Uplink (raw)**: `iot.raw.{vendor}.uplink`
- **Downlink (raw)**: `iot.raw.{vendor}.downlink`
- **Standard messages**: `iot.{vendor}.{service}.{action}`

## Best Practices

1. **Error Handling**: Define custom errors in `errors.go` for better error identification
2. **Validation**: Validate messages early in the parsing process
3. **Logging**: Use structured logging with trace IDs
4. **Metrics**: Export Prometheus metrics for monitoring
5. **Testing**: Write tests first (TDD approach)
6. **Documentation**: Document protocol-specific behaviors

## Reference Implementation

See `pkg/adapter/dji/` for a complete reference implementation of the DJI protocol adapter.
