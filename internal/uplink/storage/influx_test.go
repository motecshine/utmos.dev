package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/pkg/adapter"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "http://localhost:8086", config.URL)
	assert.Equal(t, "", config.Token)
	assert.Equal(t, "utmos", config.Org)
	assert.Equal(t, "iot", config.Bucket)
	assert.Equal(t, 1000, config.BatchSize)
	assert.Equal(t, time.Second, config.FlushInterval)
}

func TestNewStorage(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		storage := NewStorage(nil, nil)
		require.NotNil(t, storage)
		assert.NotNil(t, storage.client)
		assert.NotNil(t, storage.writeAPI)
		_ = storage.Close()
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &Config{
			URL:           "http://custom:8086",
			Token:         "test-token",
			Org:           "test-org",
			Bucket:        "test-bucket",
			BatchSize:     500,
			FlushInterval: 2 * time.Second,
		}

		storage := NewStorage(config, nil)
		require.NotNil(t, storage)
		assert.Equal(t, config, storage.config)
		_ = storage.Close()
	})
}

func TestStorage_WriteProcessedMessage(t *testing.T) {
	storage := NewStorage(nil, nil)
	defer func() { _ = storage.Close() }()

	t.Run("nil message", func(t *testing.T) {
		err := storage.WriteProcessedMessage(context.Background(), nil)
		assert.Error(t, err)
	})

	t.Run("property message", func(t *testing.T) {
		msg := &adapter.ProcessedMessage{
			MessageType: adapter.MessageTypeProperty,
			DeviceSN:    "DEVICE001",
			Vendor:      "dji",
			Properties: map[string]any{
				"latitude":  22.5431,
				"longitude": 113.9234,
				"altitude":  100.5,
			},
			Timestamp: time.Now().UnixMilli(),
		}

		err := storage.WriteProcessedMessage(context.Background(), msg)
		assert.NoError(t, err)
	})

	t.Run("event message", func(t *testing.T) {
		msg := &adapter.ProcessedMessage{
			MessageType: adapter.MessageTypeEvent,
			DeviceSN:    "DEVICE002",
			Vendor:      "dji",
			Events: []adapter.Event{
				{
					Name: "takeoff_complete",
					Params: map[string]any{
						"height": 100,
					},
				},
			},
			Timestamp: time.Now().UnixMilli(),
		}

		err := storage.WriteProcessedMessage(context.Background(), msg)
		assert.NoError(t, err)
	})

	t.Run("empty properties", func(t *testing.T) {
		msg := &adapter.ProcessedMessage{
			MessageType: adapter.MessageTypeProperty,
			DeviceSN:    "DEVICE003",
			Vendor:      "dji",
			Properties:  map[string]any{},
			Timestamp:   time.Now().UnixMilli(),
		}

		err := storage.WriteProcessedMessage(context.Background(), msg)
		assert.NoError(t, err)
	})
}

func TestStorage_WriteTelemetry(t *testing.T) {
	storage := NewStorage(nil, nil)
	defer func() { _ = storage.Close() }()

	t.Run("valid telemetry", func(t *testing.T) {
		fields := map[string]any{
			"temperature": 25.5,
			"humidity":    60,
		}
		tags := map[string]string{
			"location": "warehouse",
		}

		err := storage.WriteTelemetry(
			context.Background(),
			"DEVICE001",
			"dji",
			"environment",
			fields,
			tags,
			time.Now(),
		)
		assert.NoError(t, err)
	})

	t.Run("empty fields", func(t *testing.T) {
		err := storage.WriteTelemetry(
			context.Background(),
			"DEVICE001",
			"dji",
			"environment",
			map[string]any{},
			nil,
			time.Now(),
		)
		assert.Error(t, err)
	})
}

func TestStorage_Close(t *testing.T) {
	storage := NewStorage(nil, nil)

	err := storage.Close()
	assert.NoError(t, err)

	// Double close should not error
	err = storage.Close()
	assert.NoError(t, err)

	// Operations after close should fail
	err = storage.WriteProcessedMessage(context.Background(), &adapter.ProcessedMessage{
		DeviceSN:  "DEVICE001",
		Vendor:    "dji",
		Timestamp: time.Now().UnixMilli(),
	})
	assert.Error(t, err)
}

func TestStorage_GetMeasurement(t *testing.T) {
	storage := NewStorage(nil, nil)
	defer func() { _ = storage.Close() }()

	testCases := []struct {
		messageType adapter.MessageType
		expected    string
	}{
		{adapter.MessageTypeProperty, "telemetry"},
		{adapter.MessageTypeEvent, "events"},
		{adapter.MessageTypeService, "services"},
		{adapter.MessageTypeStatus, "status"},
		{adapter.MessageType("unknown"), "telemetry"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.messageType), func(t *testing.T) {
			msg := &adapter.ProcessedMessage{
				MessageType: tc.messageType,
			}
			result := storage.getMeasurement(msg)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestStorage_ConvertToFields(t *testing.T) {
	storage := NewStorage(nil, nil)
	defer func() { _ = storage.Close() }()

	data := map[string]any{
		"float64":  64.5,
		"float32":  float32(32.5),
		"int":      42,
		"int64":    int64(64),
		"uint":     uint(10),
		"string":   "test",
		"bool":     true,
		"nil":      nil,
		"complex":  []int{1, 2, 3},
	}

	fields := storage.convertToFields(data)

	assert.Equal(t, 64.5, fields["float64"])
	assert.Equal(t, float32(32.5), fields["float32"])
	assert.Equal(t, 42, fields["int"])
	assert.Equal(t, int64(64), fields["int64"])
	assert.Equal(t, uint(10), fields["uint"])
	assert.Equal(t, "test", fields["string"])
	assert.Equal(t, true, fields["bool"])
	assert.NotContains(t, fields, "nil")
	assert.Contains(t, fields, "complex") // converted to string
}

func TestTelemetryPoint(t *testing.T) {
	t.Run("create and build point", func(t *testing.T) {
		point := NewTelemetryPoint("DEVICE001", "dji", "telemetry").
			AddTag("location", "warehouse").
			AddField("temperature", 25.5).
			AddField("humidity", 60).
			SetTimestamp(time.Unix(1704067200, 0))

		assert.Equal(t, "DEVICE001", point.DeviceSN)
		assert.Equal(t, "dji", point.Vendor)
		assert.Equal(t, "telemetry", point.Measurement)
		assert.Equal(t, "warehouse", point.Tags["location"])
		assert.Equal(t, 25.5, point.Fields["temperature"])
		assert.Equal(t, 60, point.Fields["humidity"])
	})

	t.Run("convert to influx point", func(t *testing.T) {
		point := NewTelemetryPoint("DEVICE001", "dji", "telemetry").
			AddField("value", 100)

		influxPoint := point.ToInfluxPoint()
		require.NotNil(t, influxPoint)
	})
}

func TestStorage_Flush(t *testing.T) {
	storage := NewStorage(nil, nil)
	defer func() { _ = storage.Close() }()

	// Write some data
	msg := &adapter.ProcessedMessage{
		MessageType: adapter.MessageTypeProperty,
		DeviceSN:    "DEVICE001",
		Vendor:      "dji",
		Properties: map[string]any{
			"value": 100,
		},
		Timestamp: time.Now().UnixMilli(),
	}
	_ = storage.WriteProcessedMessage(context.Background(), msg)

	// Flush should not panic
	storage.Flush()
}
