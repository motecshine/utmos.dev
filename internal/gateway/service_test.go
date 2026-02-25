package gateway

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/gateway/bridge"
	"github.com/utmos/utmos/internal/gateway/mqtt"
)

func TestDefaultServiceConfig(t *testing.T) {
	config := DefaultServiceConfig()

	require.NotNil(t, config)
	assert.NotNil(t, config.MQTT)
	assert.NotNil(t, config.UplinkBridge)
	assert.NotNil(t, config.DownlinkBridge)
	assert.Equal(t, 5*time.Minute, config.CleanupInterval)
	assert.Equal(t, 24*time.Hour, config.MaxStaleAge)
	assert.Len(t, config.SubscribeTopics, 2)
}

func TestNewService(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		svc := NewService(nil, nil, nil, nil, nil)
		require.NotNil(t, svc)
		assert.NotNil(t, svc.config)
		assert.NotNil(t, svc.mqttClient)
		assert.NotNil(t, svc.mqttHandler)
		assert.NotNil(t, svc.connManager)
		assert.NotNil(t, svc.uplinkBridge)
		assert.NotNil(t, svc.downlinkBridge)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &ServiceConfig{
			MQTT: &mqtt.Config{
				Broker:   "custom-broker",
				Port:     1884,
				ClientID: "custom-client",
			},
			UplinkBridge:    bridge.DefaultUplinkBridgeConfig(),
			DownlinkBridge:  bridge.DefaultDownlinkBridgeConfig(),
			CleanupInterval: 10 * time.Minute,
			MaxStaleAge:     48 * time.Hour,
			SubscribeTopics: []string{"custom/#"},
		}

		svc := NewService(config, nil, nil, nil, nil)
		require.NotNil(t, svc)
		assert.Equal(t, "custom-broker", svc.config.MQTT.Broker)
		assert.Equal(t, 1884, svc.config.MQTT.Port)
		assert.Equal(t, 10*time.Minute, svc.config.CleanupInterval)
	})
}

func TestService_IsRunning(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)
	assert.False(t, svc.IsRunning())
}

func TestService_IsMQTTConnected(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)
	// Not connected initially
	assert.False(t, svc.IsMQTTConnected())
}

func TestService_GetOnlineDeviceCount(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)
	assert.Equal(t, 0, svc.GetOnlineDeviceCount())
}

func TestService_GetConnectionManager(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)
	assert.NotNil(t, svc.GetConnectionManager())
}

func TestService_GetMQTTClient(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)
	assert.NotNil(t, svc.GetMQTTClient())
}

func TestService_GetMQTTHandler(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)
	assert.NotNil(t, svc.GetMQTTHandler())
}

func TestService_RegisterDevice(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)

	state := svc.RegisterDevice("device-001", "client-001", "192.168.1.100")

	require.NotNil(t, state)
	assert.Equal(t, "device-001", state.DeviceSN)
	assert.Equal(t, "client-001", state.ClientID)
	assert.Equal(t, "192.168.1.100", state.IPAddress)
	assert.True(t, state.Online)
	assert.Equal(t, 1, svc.GetOnlineDeviceCount())
}

func TestService_UnregisterDevice(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)

	// Register first
	svc.RegisterDevice("device-001", "client-001", "192.168.1.100")
	assert.Equal(t, 1, svc.GetOnlineDeviceCount())

	// Unregister
	state := svc.UnregisterDevice("device-001")
	require.NotNil(t, state)
	assert.False(t, state.Online)
	assert.Equal(t, 0, svc.GetOnlineDeviceCount())
}

func TestService_IsDeviceOnline(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)

	// Initially not online
	assert.False(t, svc.IsDeviceOnline("device-001"))

	// Register
	svc.RegisterDevice("device-001", "client-001", "192.168.1.100")
	assert.True(t, svc.IsDeviceOnline("device-001"))

	// Unregister
	svc.UnregisterDevice("device-001")
	assert.False(t, svc.IsDeviceOnline("device-001"))
}

func TestService_Stop_NotRunning(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)

	// Should not error when stopping a non-running service
	err := svc.Stop()
	assert.NoError(t, err)
}

func TestService_Start_AlreadyRunning(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)

	// Manually set running state
	svc.mu.Lock()
	svc.running = true
	svc.mu.Unlock()

	// Should error when starting an already running service
	err := svc.Start(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already running")
}

func TestService_MultipleDevices(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)

	// Register multiple devices
	svc.RegisterDevice("device-001", "client-001", "192.168.1.100")
	svc.RegisterDevice("device-002", "client-002", "192.168.1.101")
	svc.RegisterDevice("device-003", "client-003", "192.168.1.102")

	assert.Equal(t, 3, svc.GetOnlineDeviceCount())

	// Unregister one
	svc.UnregisterDevice("device-002")
	assert.Equal(t, 2, svc.GetOnlineDeviceCount())

	// Check individual status
	assert.True(t, svc.IsDeviceOnline("device-001"))
	assert.False(t, svc.IsDeviceOnline("device-002"))
	assert.True(t, svc.IsDeviceOnline("device-003"))
}
