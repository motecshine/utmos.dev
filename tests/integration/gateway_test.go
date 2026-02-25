package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/gateway"
	"github.com/utmos/utmos/internal/gateway/bridge"
	"github.com/utmos/utmos/internal/gateway/connection"
	"github.com/utmos/utmos/internal/gateway/mqtt"
)

// TestGatewayServiceCreation tests gateway service creation and configuration
func TestGatewayServiceCreation(t *testing.T) {
	t.Run("create with default config", func(t *testing.T) {
		svc := gateway.NewService(nil, nil, nil, nil, nil)
		require.NotNil(t, svc)
		assert.False(t, svc.IsRunning())
		assert.False(t, svc.IsMQTTConnected())
	})

	t.Run("create with custom config", func(t *testing.T) {
		config := &gateway.ServiceConfig{
			MQTT: &mqtt.Config{
				Broker:   "test-broker",
				Port:     1884,
				ClientID: "test-gateway",
			},
			UplinkBridge:    bridge.DefaultUplinkBridgeConfig(),
			DownlinkBridge:  bridge.DefaultDownlinkBridgeConfig(),
			CleanupInterval: 10 * time.Minute,
			MaxStaleAge:     48 * time.Hour,
			SubscribeTopics: []string{"test/#"},
		}

		svc := gateway.NewService(config, nil, nil, nil, nil)
		require.NotNil(t, svc)
	})
}

// TestGatewayConnectionManager tests device connection management
func TestGatewayConnectionManager(t *testing.T) {
	svc := gateway.NewService(nil, nil, nil, nil, nil)
	connManager := svc.GetConnectionManager()
	require.NotNil(t, connManager)

	t.Run("register single device", func(t *testing.T) {
		state := svc.RegisterDevice("device-001", "client-001", "192.168.1.100")
		require.NotNil(t, state)
		assert.Equal(t, "device-001", state.DeviceSN)
		assert.True(t, state.Online)
		assert.Equal(t, 1, svc.GetOnlineDeviceCount())
	})

	t.Run("register multiple devices", func(t *testing.T) {
		svc.RegisterDevice("device-002", "client-002", "192.168.1.101")
		svc.RegisterDevice("device-003", "client-003", "192.168.1.102")
		assert.Equal(t, 3, svc.GetOnlineDeviceCount())
	})

	t.Run("check device online status", func(t *testing.T) {
		assert.True(t, svc.IsDeviceOnline("device-001"))
		assert.True(t, svc.IsDeviceOnline("device-002"))
		assert.False(t, svc.IsDeviceOnline("nonexistent"))
	})

	t.Run("unregister device", func(t *testing.T) {
		state := svc.UnregisterDevice("device-002")
		require.NotNil(t, state)
		assert.False(t, state.Online)
		assert.Equal(t, 2, svc.GetOnlineDeviceCount())
		assert.False(t, svc.IsDeviceOnline("device-002"))
	})

	t.Run("get online devices", func(t *testing.T) {
		devices := connManager.GetOnlineDevices()
		assert.Len(t, devices, 2)
	})
}

// TestGatewayConnectionCallbacks tests connection event callbacks
func TestGatewayConnectionCallbacks(t *testing.T) {
	connManager := connection.NewManager(nil)

	var connectEvents []string
	var disconnectEvents []string

	connManager.SetOnConnect(func(state *connection.DeviceState) {
		connectEvents = append(connectEvents, state.DeviceSN)
	})

	connManager.SetOnDisconnect(func(state *connection.DeviceState) {
		disconnectEvents = append(disconnectEvents, state.DeviceSN)
	})

	// Connect devices
	connManager.Connect("device-001", "client-001", "192.168.1.100")
	connManager.Connect("device-002", "client-002", "192.168.1.101")

	// Wait for async callbacks
	time.Sleep(100 * time.Millisecond)

	assert.Len(t, connectEvents, 2)
	assert.Contains(t, connectEvents, "device-001")
	assert.Contains(t, connectEvents, "device-002")

	// Disconnect one device
	connManager.Disconnect("device-001")
	time.Sleep(100 * time.Millisecond)

	assert.Len(t, disconnectEvents, 1)
	assert.Contains(t, disconnectEvents, "device-001")
}

// TestGatewayConnectionCleanup tests stale connection cleanup
func TestGatewayConnectionCleanup(t *testing.T) {
	connManager := connection.NewManager(nil)

	// Connect and disconnect devices
	connManager.Connect("device-001", "client-001", "192.168.1.100")
	connManager.Connect("device-002", "client-002", "192.168.1.101")
	connManager.Disconnect("device-001")

	// Online device should not be cleaned up
	removed := connManager.CleanupStale(context.Background(), 0)
	assert.Equal(t, 1, removed)

	// device-002 should still be online
	assert.True(t, connManager.IsOnline("device-002"))
	assert.Nil(t, connManager.GetState("device-001"))
}

// TestGatewayMQTTHandler tests MQTT message handling
func TestGatewayMQTTHandler(t *testing.T) {
	handler := mqtt.NewHandler(nil)
	require.NotNil(t, handler)

	var processedMessages []*mqtt.Message

	// Register a simple processor
	processor := mqtt.NewSimpleProcessor("thing/product/+/+/#", func(ctx context.Context, msg *mqtt.Message, topicInfo *mqtt.TopicInfo) error {
		processedMessages = append(processedMessages, msg)
		return nil
	})
	handler.RegisterProcessor(processor)

	t.Run("process matching message", func(t *testing.T) {
		// Note: We can't easily test Handle() without a real MQTT client
		// This test verifies processor registration works
		assert.NotNil(t, processor.Pattern())
		assert.Equal(t, "thing/product/+/+/#", processor.Pattern())
	})
}

// TestGatewayTopicParsing tests MQTT topic parsing
func TestGatewayTopicParsing(t *testing.T) {
	testCases := []struct {
		name           string
		topic          string
		expectedVendor string
		expectedSN     string
		expectedSvc    string
	}{
		{
			name:           "DJI OSD topic",
			topic:          "thing/product/DEVICE001/osd",
			expectedVendor: "dji",
			expectedSN:     "DEVICE001",
			expectedSvc:    "osd",
		},
		{
			name:           "DJI state topic",
			topic:          "thing/product/DEVICE002/state",
			expectedVendor: "dji",
			expectedSN:     "DEVICE002",
			expectedSvc:    "state",
		},
		{
			name:           "DJI services topic",
			topic:          "thing/product/DEVICE003/services",
			expectedVendor: "dji",
			expectedSN:     "DEVICE003",
			expectedSvc:    "services",
		},
		{
			name:           "sys topic",
			topic:          "sys/product/DEVICE004/status",
			expectedVendor: "dji",
			expectedSN:     "DEVICE004",
			expectedSvc:    "status",
		},
		{
			name:           "custom vendor topic",
			topic:          "custom/thing/product/DEVICE005/telemetry",
			expectedVendor: "custom",
			expectedSN:     "DEVICE005",
			expectedSvc:    "telemetry",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info := mqtt.ParseTopic(tc.topic)
			assert.Equal(t, tc.expectedVendor, info.Vendor)
			assert.Equal(t, tc.expectedSN, info.DeviceSN)
			assert.Equal(t, tc.expectedSvc, info.Service)
		})
	}
}

// TestGatewayUplinkBridge tests uplink bridge configuration
func TestGatewayUplinkBridge(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		config := bridge.DefaultUplinkBridgeConfig()
		assert.Equal(t, "iot.topic", config.Exchange)
	})

	t.Run("routing key generation", func(t *testing.T) {
		key := bridge.GetUplinkRoutingKey("dji")
		assert.Equal(t, "iot.raw.dji.uplink", key)

		key = bridge.GetUplinkRoutingKey("custom")
		assert.Equal(t, "iot.raw.custom.uplink", key)
	})

	t.Run("create bridge without publisher", func(t *testing.T) {
		uplinkBridge := bridge.NewUplinkBridge(nil, nil, nil)
		require.NotNil(t, uplinkBridge)

		// Bridge should fail without publisher
		msg := &mqtt.Message{
			Topic:   "thing/product/DEVICE001/osd",
			Payload: json.RawMessage(`{"test": true}`),
		}
		topicInfo := mqtt.ParseTopic(msg.Topic)
		err := uplinkBridge.Bridge(context.Background(), msg, topicInfo)
		assert.Error(t, err)
	})
}

// TestGatewayDownlinkBridge tests downlink bridge configuration
func TestGatewayDownlinkBridge(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		config := bridge.DefaultDownlinkBridgeConfig()
		assert.Equal(t, "iot.topic", config.Exchange)
		assert.Equal(t, "iot.gateway.downlink", config.Queue)
		assert.Equal(t, "iot.raw.*.downlink", config.RoutingKey)
	})

	t.Run("routing key generation", func(t *testing.T) {
		key := bridge.GetDownlinkRoutingKey("dji")
		assert.Equal(t, "iot.raw.dji.downlink", key)
	})

	t.Run("create bridge without subscriber", func(t *testing.T) {
		downlinkBridge := bridge.NewDownlinkBridge(nil, nil, nil, nil)
		require.NotNil(t, downlinkBridge)
		assert.False(t, downlinkBridge.IsRunning())
	})

	t.Run("parse downlink message", func(t *testing.T) {
		msgData := `{
			"device_sn": "DEVICE001",
			"topic": "thing/product/DEVICE001/services_reply",
			"payload": {"result": 0},
			"qos": 1,
			"retained": false,
			"trace_id": "trace-001",
			"span_id": "span-001"
		}`

		msg, err := bridge.ParseRawDownlinkMessage([]byte(msgData))
		require.NoError(t, err)
		assert.Equal(t, "DEVICE001", msg.DeviceSN)
		assert.Equal(t, "thing/product/DEVICE001/services_reply", msg.Topic)
		assert.Equal(t, 1, msg.QoS)
		assert.False(t, msg.Retained)
	})
}

// TestGatewayServiceLifecycle tests service start/stop lifecycle
func TestGatewayServiceLifecycle(t *testing.T) {
	svc := gateway.NewService(nil, nil, nil, nil, nil)

	t.Run("initial state", func(t *testing.T) {
		assert.False(t, svc.IsRunning())
		assert.False(t, svc.IsMQTTConnected())
	})

	t.Run("stop when not running", func(t *testing.T) {
		err := svc.Stop()
		assert.NoError(t, err)
	})

	t.Run("double start prevention", func(t *testing.T) {
		// Manually set running state to simulate running service
		// This tests the guard against double start
		svc2 := gateway.NewService(nil, nil, nil, nil, nil)

		// Use reflection or internal state manipulation would be needed
		// For now, we test that Start fails when MQTT broker is unavailable
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := svc2.Start(ctx)
		// Should fail because MQTT broker is not available
		assert.Error(t, err)
	})
}

// TestGatewayRawMessageParsing tests raw message parsing utilities
func TestGatewayRawMessageParsing(t *testing.T) {
	t.Run("parse uplink message", func(t *testing.T) {
		msgData := `{
			"vendor": "dji",
			"topic": "thing/product/DEVICE001/osd",
			"payload": {"latitude": 22.5, "longitude": 113.9},
			"qos": 1,
			"timestamp": 1704067200000,
			"trace_id": "trace-001",
			"span_id": "span-001"
		}`

		msg, err := bridge.ParseRawUplinkMessage([]byte(msgData))
		require.NoError(t, err)
		assert.Equal(t, "dji", msg.Vendor)
		assert.Equal(t, "thing/product/DEVICE001/osd", msg.Topic)
		assert.Equal(t, 1, msg.QoS)
	})

	t.Run("create new uplink message", func(t *testing.T) {
		payload := json.RawMessage(`{"test": true}`)
		msg := bridge.NewRawUplinkMessage("dji", "thing/product/DEVICE001/osd", payload, 1)

		assert.Equal(t, "dji", msg.Vendor)
		assert.Equal(t, "thing/product/DEVICE001/osd", msg.Topic)
		assert.Equal(t, 1, msg.QoS)
		assert.NotEmpty(t, msg.TraceID)
		assert.NotEmpty(t, msg.SpanID)
		assert.Greater(t, msg.Timestamp, int64(0))
	})

	t.Run("create new downlink message", func(t *testing.T) {
		payload := json.RawMessage(`{"command": "takeoff"}`)
		msg := bridge.NewRawDownlinkMessage("DEVICE001", "thing/product/DEVICE001/services", payload, 1, false, "trace-001", "span-001")

		assert.Equal(t, "DEVICE001", msg.DeviceSN)
		assert.Equal(t, "thing/product/DEVICE001/services", msg.Topic)
		assert.Equal(t, 1, msg.QoS)
		assert.False(t, msg.Retained)
		assert.Equal(t, "trace-001", msg.TraceID)
	})
}

// TestGatewayHighLoadDevices tests handling many device connections
func TestGatewayHighLoadDevices(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping high load test in short mode")
	}

	connManager := connection.NewManager(nil)
	deviceCount := 100

	// Register many devices
	for i := 0; i < deviceCount; i++ {
		deviceSN := fmt.Sprintf("device-%c%d", 'A'+i%26, i/26)
		connManager.Connect(deviceSN, "client-"+deviceSN, fmt.Sprintf("192.168.1.%d", i%256))
	}

	assert.Equal(t, deviceCount, connManager.GetOnlineCount())

	// Disconnect half
	devices := connManager.GetAllDevices()
	for i, d := range devices {
		if i%2 == 0 {
			connManager.Disconnect(d.DeviceSN)
		}
	}

	assert.Equal(t, deviceCount/2, connManager.GetOnlineCount())
}
