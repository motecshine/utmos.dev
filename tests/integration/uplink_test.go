package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/uplink"
	"github.com/utmos/utmos/internal/uplink/router"
	"github.com/utmos/utmos/internal/uplink/storage"
	"github.com/utmos/utmos/pkg/adapter"
	djiuplink "github.com/utmos/utmos/pkg/adapter/dji/uplink"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// TestUplinkServiceCreation tests uplink service creation and configuration
func TestUplinkServiceCreation(t *testing.T) {
	t.Run("create with default config", func(t *testing.T) {
		svc := uplink.NewService(nil, nil, nil, nil, nil)
		require.NotNil(t, svc)
		defer func() { _ = svc.Stop() }()

		assert.False(t, svc.IsRunning())
		assert.NotNil(t, svc.GetRegistry())
		assert.NotNil(t, svc.GetHandler())
	})

	t.Run("create with custom config", func(t *testing.T) {
		config := &uplink.ServiceConfig{
			QueueName:     "custom.queue",
			RoutingKeys:   []string{"custom.#"},
			Influx:        storage.DefaultConfig(),
			Router:        router.DefaultConfig(),
			EnableStorage: false,
			EnableRouting: false,
		}

		svc := uplink.NewService(config, nil, nil, nil, nil)
		require.NotNil(t, svc)
		defer func() { _ = svc.Stop() }()

		assert.Nil(t, svc.GetStorage())
		assert.Nil(t, svc.GetRouter())
	})
}

// TestUplinkRegistry tests processor registration
func TestUplinkRegistry(t *testing.T) {
	svc := uplink.NewService(nil, nil, nil, nil, nil)
	defer func() { _ = svc.Stop() }()

	registry := svc.GetRegistry()
	require.NotNil(t, registry)

	t.Run("no processors registered by default", func(t *testing.T) {
		vendors := registry.ListVendors()
		assert.Empty(t, vendors, "no processors should be registered by default")
	})

	t.Run("register DJI processor", func(t *testing.T) {
		djiProcessor := djiuplink.NewProcessorAdapter(nil)
		svc.RegisterProcessor(djiProcessor)

		vendors := registry.ListVendors()
		assert.Contains(t, vendors, "dji")
	})

	t.Run("get DJI processor after registration", func(t *testing.T) {
		p, exists := registry.Get("dji")
		assert.True(t, exists)
		assert.Equal(t, "dji", p.GetVendor())
	})

	t.Run("register custom processor", func(t *testing.T) {
		customProcessor := &mockUplinkProcessor{vendor: "custom"}
		svc.RegisterProcessor(customProcessor)

		vendors := registry.ListVendors()
		assert.Contains(t, vendors, "custom")
	})
}

// TestUplinkDJIMessageProcessing tests DJI message processing
func TestUplinkDJIMessageProcessing(t *testing.T) {
	config := &uplink.ServiceConfig{
		QueueName:     "test.queue",
		RoutingKeys:   []string{"test.#"},
		EnableStorage: false,
		EnableRouting: false,
	}
	svc := uplink.NewService(config, nil, nil, nil, nil)
	defer func() { _ = svc.Stop() }()

	// Register DJI processor
	djiProcessor := djiuplink.NewProcessorAdapter(nil)
	svc.RegisterProcessor(djiProcessor)

	t.Run("process OSD property message", func(t *testing.T) {
		data := map[string]interface{}{
			"latitude":        22.5431,
			"longitude":       113.9234,
			"altitude":        100.5,
			"height":          50.0,
			"speed":           10.5,
			"battery_percent": 85.0,
			"flight_mode":     "GPS",
		}
		dataBytes, _ := json.Marshal(data)

		msg := &rabbitmq.StandardMessage{
			TID:       "tid-001",
			BID:       "bid-001",
			DeviceSN:  "DEVICE001",
			Action:    "property.report",
			Timestamp: time.Now().UnixMilli(),
			Data:      dataBytes,
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: "dji",
			},
		}

		err := svc.ProcessMessage(context.Background(), msg)
		assert.NoError(t, err)
	})

	t.Run("process event message", func(t *testing.T) {
		data := map[string]interface{}{
			"method": "fly_to_point_progress",
			"data": map[string]interface{}{
				"progress": 50,
				"status":   "flying",
			},
		}
		dataBytes, _ := json.Marshal(data)

		msg := &rabbitmq.StandardMessage{
			TID:       "tid-002",
			BID:       "bid-002",
			DeviceSN:  "DEVICE002",
			Action:    "event.report",
			Timestamp: time.Now().UnixMilli(),
			Data:      dataBytes,
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: "dji",
				Method: "fly_to_point_progress",
			},
		}

		err := svc.ProcessMessage(context.Background(), msg)
		assert.NoError(t, err)
	})

	t.Run("process device status message", func(t *testing.T) {
		msg := &rabbitmq.StandardMessage{
			TID:       "tid-003",
			BID:       "bid-003",
			DeviceSN:  "DEVICE003",
			Action:    "device.online",
			Timestamp: time.Now().UnixMilli(),
			Data:      json.RawMessage(`{"online": true}`),
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: "dji",
			},
		}

		err := svc.ProcessMessage(context.Background(), msg)
		assert.NoError(t, err)
	})
}

// TestUplinkDJIProcessor tests DJI processor directly
func TestUplinkDJIProcessor(t *testing.T) {
	djiProcessor := djiuplink.NewProcessorAdapter(nil)

	t.Run("can process DJI messages", func(t *testing.T) {
		msg := &rabbitmq.StandardMessage{
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: "dji",
			},
		}
		assert.True(t, djiProcessor.CanProcess(msg))
	})

	t.Run("process property message", func(t *testing.T) {
		data := map[string]interface{}{
			"latitude":  22.5431,
			"longitude": 113.9234,
		}
		dataBytes, _ := json.Marshal(data)

		msg := &rabbitmq.StandardMessage{
			TID:       "tid-001",
			BID:       "bid-001",
			DeviceSN:  "DEVICE001",
			Action:    "property.report",
			Timestamp: time.Now().UnixMilli(),
			Data:      dataBytes,
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: "dji",
			},
		}

		processed, err := djiProcessor.Process(context.Background(), msg)
		require.NoError(t, err)
		require.NotNil(t, processed)

		assert.Equal(t, adapter.MessageTypeProperty, processed.MessageType)
		assert.Equal(t, "DEVICE001", processed.DeviceSN)
		assert.Equal(t, "dji", processed.Vendor)
		assert.Equal(t, 22.5431, processed.Properties["latitude"])
		assert.Equal(t, 113.9234, processed.Properties["longitude"])
	})

	t.Run("flatten nested properties", func(t *testing.T) {
		data := map[string]interface{}{
			"data": map[string]interface{}{
				"position": map[string]interface{}{
					"latitude":  22.5431,
					"longitude": 113.9234,
				},
				"battery": 85,
			},
		}
		dataBytes, _ := json.Marshal(data)

		msg := &rabbitmq.StandardMessage{
			TID:       "tid-002",
			BID:       "bid-002",
			DeviceSN:  "DEVICE002",
			Action:    "property.report",
			Timestamp: time.Now().UnixMilli(),
			Data:      dataBytes,
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: "dji",
			},
		}

		processed, err := djiProcessor.Process(context.Background(), msg)
		require.NoError(t, err)

		// Nested properties should be flattened
		assert.Equal(t, 22.5431, processed.Properties["position.latitude"])
		assert.Equal(t, 113.9234, processed.Properties["position.longitude"])
		assert.Equal(t, float64(85), processed.Properties["battery"])
	})
}

// TestUplinkInfluxStorage tests InfluxDB storage
func TestUplinkInfluxStorage(t *testing.T) {
	influxStorage := storage.NewStorage(nil, nil)
	defer func() { _ = influxStorage.Close() }()

	t.Run("write processed message", func(t *testing.T) {
		msg := &adapter.ProcessedMessage{
			MessageType: adapter.MessageTypeProperty,
			DeviceSN:    "DEVICE001",
			Vendor:      "dji",
			Properties: map[string]interface{}{
				"latitude":  22.5431,
				"longitude": 113.9234,
				"altitude":  100.5,
			},
			Timestamp: time.Now().UnixMilli(),
		}

		err := influxStorage.WriteProcessedMessage(context.Background(), msg)
		assert.NoError(t, err)
	})

	t.Run("write telemetry point", func(t *testing.T) {
		point := storage.NewTelemetryPoint("DEVICE001", "dji", "telemetry").
			AddTag("location", "warehouse").
			AddField("temperature", 25.5).
			AddField("humidity", 60).
			SetTimestamp(time.Now())

		err := influxStorage.WritePoint(point.ToInfluxPoint())
		assert.NoError(t, err)
	})
}

// TestUplinkRouter tests message routing
func TestUplinkRouter(t *testing.T) {
	t.Run("routing key generation", func(t *testing.T) {
		assert.Equal(t, "iot.ws.property", router.RoutingKeyWSProperty)
		assert.Equal(t, "iot.ws.event", router.RoutingKeyWSEvent)
		assert.Equal(t, "iot.ws.status", router.RoutingKeyWSStatus)
		assert.Equal(t, "iot.api.property", router.RoutingKeyAPIProperty)
		assert.Equal(t, "iot.api.event", router.RoutingKeyAPIEvent)
	})

	t.Run("router configuration", func(t *testing.T) {
		config := router.DefaultConfig()
		assert.Equal(t, "iot.topic", config.Exchange)
		assert.True(t, config.EnableWSRouting)
		assert.True(t, config.EnableAPIRouting)
	})

	t.Run("multi-router", func(t *testing.T) {
		multiRouter := router.NewMultiRouter(nil)

		var routedMessages []string
		multiRouter.AddRouter(func(ctx context.Context, msg *adapter.ProcessedMessage) error {
			routedMessages = append(routedMessages, "router1:"+msg.DeviceSN)
			return nil
		})
		multiRouter.AddRouter(func(ctx context.Context, msg *adapter.ProcessedMessage) error {
			routedMessages = append(routedMessages, "router2:"+msg.DeviceSN)
			return nil
		})

		msg := &adapter.ProcessedMessage{
			DeviceSN: "DEVICE001",
		}

		err := multiRouter.Route(context.Background(), msg)
		assert.NoError(t, err)
		assert.Len(t, routedMessages, 2)
		assert.Contains(t, routedMessages, "router1:DEVICE001")
		assert.Contains(t, routedMessages, "router2:DEVICE001")
	})
}

// TestUplinkServiceLifecycle tests service start/stop lifecycle
func TestUplinkServiceLifecycle(t *testing.T) {
	config := &uplink.ServiceConfig{
		QueueName:     "test.queue",
		RoutingKeys:   []string{"test.#"},
		EnableStorage: false,
		EnableRouting: false,
	}
	svc := uplink.NewService(config, nil, nil, nil, nil)

	t.Run("initial state", func(t *testing.T) {
		assert.False(t, svc.IsRunning())
	})

	t.Run("start service", func(t *testing.T) {
		err := svc.Start(context.Background())
		assert.NoError(t, err)
		assert.True(t, svc.IsRunning())
	})

	t.Run("double start prevention", func(t *testing.T) {
		err := svc.Start(context.Background())
		assert.Error(t, err)
	})

	t.Run("stop service", func(t *testing.T) {
		err := svc.Stop()
		assert.NoError(t, err)
		assert.False(t, svc.IsRunning())
	})

	t.Run("double stop is safe", func(t *testing.T) {
		err := svc.Stop()
		assert.NoError(t, err)
	})
}

// TestUplinkServiceStats tests service statistics
func TestUplinkServiceStats(t *testing.T) {
	svc := uplink.NewService(nil, nil, nil, nil, nil)
	defer func() { _ = svc.Stop() }()

	// Register DJI processor
	djiProcessor := djiuplink.NewProcessorAdapter(nil)
	svc.RegisterProcessor(djiProcessor)

	stats := svc.GetStats()
	require.NotNil(t, stats)

	assert.False(t, stats.Running)
	assert.Contains(t, stats.RegisteredVendors, "dji")
	assert.True(t, stats.StorageEnabled)
	assert.True(t, stats.RoutingEnabled)
}

// TestUplinkOSDDataParsing tests OSD data parsing
func TestUplinkOSDDataParsing(t *testing.T) {
	properties := map[string]interface{}{
		"latitude":         22.5431,
		"longitude":        113.9234,
		"altitude":         100.5,
		"height":           50.0,
		"speed":            10.5,
		"heading":          180.0,
		"battery_percent":  85.0,
		"flight_mode":      "GPS",
		"gps_satellites":   12.0,
		"signal_strength":  95.0,
		"home_distance":    500.0,
		"horizontal_speed": 8.5,
		"vertical_speed":   2.0,
	}

	osd, err := djiuplink.ParseOSDData(properties)
	require.NoError(t, err)

	assert.Equal(t, 22.5431, osd.Latitude)
	assert.Equal(t, 113.9234, osd.Longitude)
	assert.Equal(t, 100.5, osd.Altitude)
	assert.Equal(t, 50.0, osd.Height)
	assert.Equal(t, 10.5, osd.Speed)
	assert.Equal(t, 180.0, osd.Heading)
	assert.Equal(t, 85, osd.BatteryPercent)
	assert.Equal(t, "GPS", osd.FlightMode)
	assert.Equal(t, 12, osd.GPSSatellites)
	assert.Equal(t, 95, osd.SignalStrength)
	assert.Equal(t, 500.0, osd.HomeDistance)
	assert.Equal(t, 8.5, osd.HorizontalSpeed)
	assert.Equal(t, 2.0, osd.VerticalSpeed)
}

// TestUplinkMessageTypes tests message type handling
func TestUplinkMessageTypes(t *testing.T) {
	assert.Equal(t, adapter.MessageType("property"), adapter.MessageTypeProperty)
	assert.Equal(t, adapter.MessageType("event"), adapter.MessageTypeEvent)
	assert.Equal(t, adapter.MessageType("service"), adapter.MessageTypeService)
	assert.Equal(t, adapter.MessageType("status"), adapter.MessageTypeStatus)
}

// mockUplinkProcessor is a test processor implementation
type mockUplinkProcessor struct {
	vendor string
}

func (p *mockUplinkProcessor) GetVendor() string {
	return p.vendor
}

func (p *mockUplinkProcessor) CanProcess(msg *rabbitmq.StandardMessage) bool {
	return msg.ProtocolMeta != nil && msg.ProtocolMeta.Vendor == p.vendor
}

func (p *mockUplinkProcessor) Process(ctx context.Context, msg *rabbitmq.StandardMessage) (*adapter.ProcessedMessage, error) {
	return &adapter.ProcessedMessage{
		Original:    msg,
		MessageType: adapter.MessageTypeProperty,
		DeviceSN:    msg.DeviceSN,
		Vendor:      p.vendor,
		Properties:  make(map[string]interface{}),
		Timestamp:   msg.Timestamp,
	}, nil
}
