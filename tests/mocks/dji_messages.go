// Package mocks provides mock DJI messages for testing.
package mocks

import (
	"encoding/json"
	"fmt"
	"time"
)

// mustMarshal marshals v to JSON and panics on error.
// This is safe to use in test mocks where marshaling should never fail.
func mustMarshal(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("mocks: failed to marshal JSON: %v", err))
	}
	return data
}

// DJIMockMessage represents a mock DJI message structure.
type DJIMockMessage struct {
	TID       string          `json:"tid"`
	BID       string          `json:"bid"`
	Timestamp int64           `json:"timestamp"`
	Method    string          `json:"method,omitempty"`
	NeedReply *int            `json:"need_reply,omitempty"`
	Data      json.RawMessage `json:"data"`
}

// OSDData represents OSD (On-Screen Display) telemetry data.
type OSDData struct {
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Altitude        float64 `json:"altitude"`
	Height          float64 `json:"height"`
	Speed           float64 `json:"speed"`
	Heading         float64 `json:"heading"`
	BatteryPercent  int     `json:"battery_percent"`
	FlightMode      string  `json:"flight_mode"`
	GPSSatellites   int     `json:"gps_satellites"`
	HomeDistance    float64 `json:"home_distance"`
	VerticalSpeed   float64 `json:"vertical_speed"`
	HorizontalSpeed float64 `json:"horizontal_speed"`
}

// StateData represents device state change data.
type StateData struct {
	FirmwareVersion string `json:"firmware_version"`
	SerialNumber    string `json:"serial_number"`
	DeviceModel     string `json:"device_model"`
	Online          bool   `json:"online"`
}

// EventData represents an event notification.
type EventData struct {
	EventType string          `json:"event_type"`
	Progress  int             `json:"progress,omitempty"`
	Result    int             `json:"result,omitempty"`
	Message   string          `json:"message,omitempty"`
	Extra     json.RawMessage `json:"extra,omitempty"`
}

// ServiceRequestData represents a service call request.
type ServiceRequestData struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

// ServiceReplyData represents a service call reply.
type ServiceReplyData struct {
	Result int             `json:"result"`
	Output json.RawMessage `json:"output,omitempty"`
}

// NewOSDMessage creates a mock OSD message.
func NewOSDMessage(deviceSN string) (topic string, payload []byte) {
	topic = "thing/product/" + deviceSN + "/osd"

	osdData := OSDData{
		Latitude:        39.9042,
		Longitude:       116.4074,
		Altitude:        100.5,
		Height:          50.0,
		Speed:           15.5,
		Heading:         180.0,
		BatteryPercent:  85,
		FlightMode:      "GPS",
		GPSSatellites:   18,
		HomeDistance:    250.0,
		VerticalSpeed:   2.0,
		HorizontalSpeed: 15.0,
	}

	msg := DJIMockMessage{
		TID:       "tid-osd-" + deviceSN,
		BID:       "bid-osd-" + deviceSN,
		Timestamp: time.Now().UnixMilli(),
		Data:      mustMarshal(osdData),
	}

	payload = mustMarshal(msg)
	return topic, payload
}

// NewStateMessage creates a mock state message.
func NewStateMessage(deviceSN string) (topic string, payload []byte) {
	topic = "thing/product/" + deviceSN + "/state"

	stateData := StateData{
		FirmwareVersion: "v01.00.0500",
		SerialNumber:    deviceSN,
		DeviceModel:     "Mavic 3 Enterprise",
		Online:          true,
	}

	msg := DJIMockMessage{
		TID:       "tid-state-" + deviceSN,
		BID:       "bid-state-" + deviceSN,
		Timestamp: time.Now().UnixMilli(),
		Data:      mustMarshal(stateData),
	}

	payload = mustMarshal(msg)
	return topic, payload
}

// NewEventMessage creates a mock event message.
func NewEventMessage(gatewaySN, eventType string, progress int) (topic string, payload []byte) {
	topic = "thing/product/" + gatewaySN + "/events"

	eventData := EventData{
		EventType: eventType,
		Progress:  progress,
		Result:    0,
		Message:   "Event in progress",
	}

	needReply := 0
	msg := DJIMockMessage{
		TID:       "tid-event-" + gatewaySN,
		BID:       "bid-event-" + gatewaySN,
		Timestamp: time.Now().UnixMilli(),
		Method:    eventType,
		NeedReply: &needReply,
		Data:      mustMarshal(eventData),
	}

	payload = mustMarshal(msg)
	return topic, payload
}

// NewServicesRequestMessage creates a mock services request message.
func NewServicesRequestMessage(gatewaySN, method string, params any) (topic string, payload []byte) {
	topic = "thing/product/" + gatewaySN + "/services"

	requestData := ServiceRequestData{
		Method: method,
		Params: mustMarshal(params),
	}

	needReply := 1
	msg := DJIMockMessage{
		TID:       "tid-svc-" + gatewaySN,
		BID:       "bid-svc-" + gatewaySN,
		Timestamp: time.Now().UnixMilli(),
		Method:    method,
		NeedReply: &needReply,
		Data:      mustMarshal(requestData),
	}

	payload = mustMarshal(msg)
	return topic, payload
}

// NewServicesReplyMessage creates a mock services reply message.
func NewServicesReplyMessage(gatewaySN, method string, result int, output any) (topic string, payload []byte) {
	topic = "thing/product/" + gatewaySN + "/services_reply"

	replyData := ServiceReplyData{
		Result: result,
		Output: mustMarshal(output),
	}

	msg := DJIMockMessage{
		TID:       "tid-reply-" + gatewaySN,
		BID:       "bid-reply-" + gatewaySN,
		Timestamp: time.Now().UnixMilli(),
		Method:    method,
		Data:      mustMarshal(replyData),
	}

	payload = mustMarshal(msg)
	return topic, payload
}

// NewStatusMessage creates a mock device status message.
func NewStatusMessage(gatewaySN string, online bool) (topic string, payload []byte) {
	topic = "sys/product/" + gatewaySN + "/status"

	statusData := map[string]any{
		"status":    "online",
		"timestamp": time.Now().UnixMilli(),
	}
	if !online {
		statusData["status"] = "offline"
	}

	msg := DJIMockMessage{
		TID:       "tid-status-" + gatewaySN,
		BID:       "bid-status-" + gatewaySN,
		Timestamp: time.Now().UnixMilli(),
		Data:      mustMarshal(statusData),
	}

	payload = mustMarshal(msg)
	return topic, payload
}

// SampleMessages returns a collection of sample DJI messages for testing.
func SampleMessages() map[string]struct {
	Topic   string
	Payload []byte
} {
	samples := make(map[string]struct {
		Topic   string
		Payload []byte
	})

	// OSD message
	osdTopic, osdPayload := NewOSDMessage("1ZNBH1D00C00FK")
	samples["osd"] = struct {
		Topic   string
		Payload []byte
	}{osdTopic, osdPayload}

	// State message
	stateTopic, statePayload := NewStateMessage("1ZNBH1D00C00FK")
	samples["state"] = struct {
		Topic   string
		Payload []byte
	}{stateTopic, statePayload}

	// Event message
	eventTopic, eventPayload := NewEventMessage("DOCK001", "fly_to_point_progress", 50)
	samples["event"] = struct {
		Topic   string
		Payload []byte
	}{eventTopic, eventPayload}

	// Services request
	svcTopic, svcPayload := NewServicesRequestMessage("DOCK001", "takeoff", map[string]any{
		"height": 50,
	})
	samples["services_request"] = struct {
		Topic   string
		Payload []byte
	}{svcTopic, svcPayload}

	// Services reply
	replyTopic, replyPayload := NewServicesReplyMessage("DOCK001", "takeoff", 0, map[string]any{
		"status": "success",
	})
	samples["services_reply"] = struct {
		Topic   string
		Payload []byte
	}{replyTopic, replyPayload}

	// Status message
	statusTopic, statusPayload := NewStatusMessage("DOCK001", true)
	samples["status"] = struct {
		Topic   string
		Payload []byte
	}{statusTopic, statusPayload}

	return samples
}
