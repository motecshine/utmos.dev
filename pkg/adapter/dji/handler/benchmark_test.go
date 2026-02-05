package handler

import (
	"context"
	"encoding/json"
	"testing"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
)

// BenchmarkOSDHandler_Handle benchmarks OSD message processing.
func BenchmarkOSDHandler_Handle(b *testing.B) {
	handler := NewOSDHandler()

	// Sample aircraft OSD data
	osdData := json.RawMessage(`{
		"mode_code": 0,
		"longitude": 113.943,
		"latitude": 22.577,
		"height": 100.5,
		"elevation": 50.2,
		"horizontal_speed": 5.0,
		"vertical_speed": 1.0,
		"attitude_pitch": 0.5,
		"attitude_roll": 0.2,
		"attitude_head": 180.0,
		"battery": {
			"capacity_percent": 85,
			"remain_flight_time": 1200
		}
	}`)

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 1234567890123,
		Method:    "osd",
		Data:      osdData,
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeOSD,
		DeviceSN:  "test-device",
		GatewaySN: "test-gateway",
		Raw:       "thing/product/test-gateway/osd",
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.Handle(ctx, msg, topic)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkStateHandler_Handle benchmarks State message processing.
func BenchmarkStateHandler_Handle(b *testing.B) {
	handler := NewStateHandler()

	stateData := json.RawMessage(`{
		"mode_code": 1,
		"cover_state": 0,
		"drone_in_dock": 1
	}`)

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 1234567890123,
		Method:    "state",
		Data:      stateData,
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeState,
		DeviceSN:  "test-device",
		GatewaySN: "test-gateway",
		Raw:       "thing/product/test-gateway/state",
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.Handle(ctx, msg, topic)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkStatusHandler_Handle benchmarks Status message processing.
func BenchmarkStatusHandler_Handle(b *testing.B) {
	handler := NewStatusHandler()

	statusData := json.RawMessage(`{
		"online": true,
		"sub_devices": [
			{"sn": "aircraft-001", "type": 60}
		]
	}`)

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 1234567890123,
		Method:    "status",
		Data:      statusData,
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeStatus,
		DeviceSN:  "test-device",
		GatewaySN: "test-gateway",
		Raw:       "sys/product/test-gateway/status",
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.Handle(ctx, msg, topic)
		if err != nil {
			b.Fatal(err)
		}
	}
}
