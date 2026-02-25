package uplink

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/uplink/router"
	"github.com/utmos/utmos/internal/uplink/storage"
	"github.com/utmos/utmos/pkg/adapter"
	djiuplink "github.com/utmos/utmos/pkg/adapter/dji/uplink"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

func TestDefaultServiceConfig(t *testing.T) {
	config := DefaultServiceConfig()

	assert.Equal(t, "iot.uplink.messages", config.QueueName)
	assert.Len(t, config.RoutingKeys, 1) // Changed to wildcard pattern
	assert.Contains(t, config.RoutingKeys, "iot.*.#")
	assert.NotNil(t, config.Influx)
	assert.NotNil(t, config.Router)
	assert.True(t, config.EnableStorage)
	assert.True(t, config.EnableRouting)
}

func TestNewService(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		svc := NewService(nil, nil, nil, nil, nil)
		require.NotNil(t, svc)
		assert.NotNil(t, svc.config)
		assert.NotNil(t, svc.registry)
		assert.NotNil(t, svc.handler)
		assert.NotNil(t, svc.storage)
		assert.NotNil(t, svc.router)

		// Cleanup
		_ = svc.Stop()
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &ServiceConfig{
			QueueName:     "custom.queue",
			RoutingKeys:   []string{"custom.#"},
			Influx:        storage.DefaultConfig(),
			Router:        router.DefaultConfig(),
			EnableStorage: false,
			EnableRouting: false,
		}

		svc := NewService(config, nil, nil, nil, nil)
		require.NotNil(t, svc)
		assert.Nil(t, svc.storage) // Disabled
		assert.Nil(t, svc.router)  // Disabled

		_ = svc.Stop()
	})

	t.Run("processors must be registered by caller", func(t *testing.T) {
		svc := NewService(nil, nil, nil, nil, nil)
		defer func() { _ = svc.Stop() }()

		// No processors registered by default
		vendors := svc.registry.ListVendors()
		assert.Empty(t, vendors)

		// Register DJI processor
		djiProcessor := djiuplink.NewProcessorAdapter(nil)
		svc.RegisterProcessor(djiProcessor)

		vendors = svc.registry.ListVendors()
		assert.Contains(t, vendors, "dji")
	})
}

func TestService_StartStop(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)
	defer func() { _ = svc.Stop() }()

	t.Run("initial state", func(t *testing.T) {
		assert.False(t, svc.IsRunning())
	})

	t.Run("start", func(t *testing.T) {
		err := svc.Start(context.Background())
		assert.NoError(t, err)
		assert.True(t, svc.IsRunning())
	})

	t.Run("double start", func(t *testing.T) {
		err := svc.Start(context.Background())
		assert.Error(t, err)
	})

	t.Run("stop", func(t *testing.T) {
		err := svc.Stop()
		assert.NoError(t, err)
		assert.False(t, svc.IsRunning())
	})

	t.Run("double stop", func(t *testing.T) {
		err := svc.Stop()
		assert.NoError(t, err)
	})
}

func TestService_GetComponents(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)
	defer func() { _ = svc.Stop() }()

	assert.NotNil(t, svc.GetRegistry())
	assert.NotNil(t, svc.GetHandler())
	assert.NotNil(t, svc.GetStorage())
	assert.NotNil(t, svc.GetRouter())
}

func TestService_RegisterProcessor(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)
	defer func() { _ = svc.Stop() }()

	// Create a mock processor
	mockProcessor := &mockTestProcessor{vendor: "custom"}
	svc.RegisterProcessor(mockProcessor)

	vendors := svc.registry.ListVendors()
	assert.Contains(t, vendors, "custom")
}

func TestService_ProcessMessage(t *testing.T) {
	// Disable routing for this test since we don't have a publisher
	config := &ServiceConfig{
		QueueName:     "test.queue",
		RoutingKeys:   []string{"test.#"},
		Influx:        storage.DefaultConfig(),
		EnableStorage: true,
		EnableRouting: false, // Disable routing for test
	}
	svc := NewService(config, nil, nil, nil, nil)
	defer func() { _ = svc.Stop() }()

	// Register DJI processor for tests
	djiProcessor := djiuplink.NewProcessorAdapter(nil)
	svc.RegisterProcessor(djiProcessor)

	t.Run("process DJI property message", func(t *testing.T) {
		data := map[string]any{
			"latitude":  22.5431,
			"longitude": 113.9234,
		}
		dataBytes, _ := json.Marshal(data)

		msg := &rabbitmq.StandardMessage{
			TID:       "tid-001",
			BID:       "bid-001",
			DeviceSN:  "DEVICE001",
			Action:    "property.report",
			Timestamp: 1704067200000,
			Data:      dataBytes,
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: "dji",
			},
		}

		err := svc.ProcessMessage(context.Background(), msg)
		assert.NoError(t, err)
	})

	t.Run("process unknown vendor message", func(t *testing.T) {
		msg := &rabbitmq.StandardMessage{
			TID:       "tid-002",
			BID:       "bid-002",
			DeviceSN:  "DEVICE002",
			Action:    "unknown.action",
			Timestamp: 1704067200000,
			Data:      json.RawMessage(`{}`),
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: "unknown",
			},
		}

		err := svc.ProcessMessage(context.Background(), msg)
		assert.Error(t, err)
	})
}

func TestService_GetStats(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)
	defer func() { _ = svc.Stop() }()

	// Register DJI processor for stats test
	djiProcessor := djiuplink.NewProcessorAdapter(nil)
	svc.RegisterProcessor(djiProcessor)

	stats := svc.GetStats()
	require.NotNil(t, stats)

	assert.False(t, stats.Running)
	assert.Contains(t, stats.RegisteredVendors, "dji")
	assert.True(t, stats.StorageEnabled)
	assert.True(t, stats.RoutingEnabled)

	// Start service
	_ = svc.Start(context.Background())
	stats = svc.GetStats()
	assert.True(t, stats.Running)
}

func TestService_DisabledComponents(t *testing.T) {
	config := &ServiceConfig{
		QueueName:     "test.queue",
		RoutingKeys:   []string{"test.#"},
		EnableStorage: false,
		EnableRouting: false,
	}

	svc := NewService(config, nil, nil, nil, nil)
	defer func() { _ = svc.Stop() }()

	assert.Nil(t, svc.GetStorage())
	assert.Nil(t, svc.GetRouter())

	stats := svc.GetStats()
	assert.False(t, stats.StorageEnabled)
	assert.False(t, stats.RoutingEnabled)
}

// mockTestProcessor is a test processor implementation
type mockTestProcessor struct {
	vendor string
}

func (p *mockTestProcessor) GetVendor() string {
	return p.vendor
}

func (p *mockTestProcessor) CanProcess(msg *rabbitmq.StandardMessage) bool {
	return msg.ProtocolMeta != nil && msg.ProtocolMeta.Vendor == p.vendor
}

func (p *mockTestProcessor) Process(ctx context.Context, msg *rabbitmq.StandardMessage) (*adapter.ProcessedMessage, error) {
	return &adapter.ProcessedMessage{
		Original:    msg,
		MessageType: adapter.MessageTypeProperty,
		DeviceSN:    msg.DeviceSN,
		Vendor:      p.vendor,
		Properties:  make(map[string]any),
		Timestamp:   msg.Timestamp,
	}, nil
}
