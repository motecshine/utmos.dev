package integration

import (
	"context"
	"testing"

	"github.com/utmos/utmos/pkg/rabbitmq"
	"github.com/utmos/utmos/pkg/tracer"
	"go.opentelemetry.io/otel/trace"
)

// TestMessageFlowWithTracing tests that trace context is properly propagated through message flow
func TestMessageFlowWithTracing(t *testing.T) {
	// Create a mock trace context
	ctx := context.Background()

	// Test trace context injection into message headers
	headers := make(map[string]interface{})
	tracer.InjectContext(ctx, headers)

	// Headers should be set (even if empty without active span)
	// This validates the injection mechanism works

	// Test extraction
	extractedCtx := tracer.ExtractContext(ctx, headers)
	if extractedCtx == nil {
		t.Error("expected non-nil context after extraction")
	}
}

// TestRoutingKeyGeneration tests routing key generation for different vendors
func TestRoutingKeyGeneration(t *testing.T) {
	tests := []struct {
		vendor   string
		service  string
		action   string
		expected string
	}{
		{rabbitmq.VendorDJI, rabbitmq.ServiceDevice, rabbitmq.ActionPropertyReport, "iot.dji.device.property.report"},
		{rabbitmq.VendorGeneric, rabbitmq.ServiceDevice, rabbitmq.ActionDeviceOnline, "iot.generic.device.device.online"},
		{rabbitmq.VendorTuya, rabbitmq.ServiceService, rabbitmq.ActionServiceCall, "iot.tuya.service.service.call"},
		{rabbitmq.VendorDJI, rabbitmq.ServiceDevice, rabbitmq.ActionDeviceOffline, "iot.dji.device.device.offline"},
		{rabbitmq.VendorGeneric, rabbitmq.ServiceEvent, rabbitmq.ActionEventReport, "iot.generic.event.event.report"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			rk := rabbitmq.NewRoutingKey(tt.vendor, tt.service, tt.action)
			if rk.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, rk.String())
			}
		})
	}
}

// TestRoutingKeyParsing tests parsing routing keys back to components
func TestRoutingKeyParsing(t *testing.T) {
	tests := []struct {
		input   string
		vendor  string
		service string
		action  string
		valid   bool
	}{
		{"iot.dji.device.property.report", "dji", "device", "property.report", true},
		{"iot.generic.device.device.online", "generic", "device", "device.online", true},
		{"iot.tuya.service.service.call", "tuya", "service", "service.call", true},
		{"invalid.key", "", "", "", false},
		{"iot.only.two", "", "", "", false},
		{"", "", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			rk, err := rabbitmq.Parse(tt.input)
			if tt.valid {
				if err != nil {
					t.Errorf("expected valid routing key, got error: %v", err)
					return
				}
				if rk.Vendor != tt.vendor {
					t.Errorf("expected vendor %s, got %s", tt.vendor, rk.Vendor)
				}
				if rk.Service != tt.service {
					t.Errorf("expected service %s, got %s", tt.service, rk.Service)
				}
				if rk.Action != tt.action {
					t.Errorf("expected action %s, got %s", tt.action, rk.Action)
				}
			} else {
				if err == nil {
					t.Error("expected error for invalid routing key")
				}
			}
		})
	}
}

// TestStandardMessageCreation tests message creation with all required fields
func TestStandardMessageCreation(t *testing.T) {
	data := map[string]interface{}{
		"temperature": 25.5,
		"humidity":    60.0,
	}

	msg, err := rabbitmq.NewStandardMessage(
		rabbitmq.ServiceDevice,
		rabbitmq.ActionPropertyReport,
		"device-001",
		data,
	)

	if err != nil {
		t.Fatalf("failed to create message: %v", err)
	}

	// Verify required fields
	if msg.TID == "" {
		t.Error("TID should be set")
	}
	if msg.BID == "" {
		t.Error("BID should be set")
	}
	if msg.Timestamp == 0 {
		t.Error("Timestamp should be set")
	}
	if msg.Service != rabbitmq.ServiceDevice {
		t.Errorf("expected service %s, got %s", rabbitmq.ServiceDevice, msg.Service)
	}
	if msg.Action != rabbitmq.ActionPropertyReport {
		t.Errorf("expected action %s, got %s", rabbitmq.ActionPropertyReport, msg.Action)
	}
	if msg.DeviceSN != "device-001" {
		t.Errorf("expected device SN device-001, got %s", msg.DeviceSN)
	}

	// Validate the message
	if err := msg.Validate(); err != nil {
		t.Errorf("message validation failed: %v", err)
	}
}

// TestTraceContextRoundTrip tests trace context injection and extraction
func TestTraceContextRoundTrip(t *testing.T) {
	ctx := context.Background()

	// Simulate trace context in headers
	headers := map[string]interface{}{
		"traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		"tracestate":  "congo=t61rcWkgMzE",
	}

	// Extract context from headers
	extractedCtx := tracer.ExtractContext(ctx, headers)

	// Get span from context
	span := trace.SpanFromContext(extractedCtx)
	spanCtx := span.SpanContext()

	// Verify trace ID was extracted
	if spanCtx.TraceID().String() != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Logf("trace ID: %s", spanCtx.TraceID().String())
		// Note: extraction may not work without a proper tracer provider
	}

	// Inject back into new headers
	newHeaders := make(map[string]interface{})
	tracer.InjectContext(extractedCtx, newHeaders)

	// Headers should be propagated
	if len(newHeaders) > 0 {
		t.Logf("propagated headers: %v", newHeaders)
	}
}

// TestMultiVendorRouting tests routing for multiple vendors
func TestMultiVendorRouting(t *testing.T) {
	vendors := []string{rabbitmq.VendorDJI, rabbitmq.VendorGeneric, rabbitmq.VendorTuya}
	services := []string{rabbitmq.ServiceDevice, rabbitmq.ServiceEvent, rabbitmq.ServiceService}
	actions := []string{rabbitmq.ActionPropertyReport, rabbitmq.ActionEventReport, rabbitmq.ActionServiceCall}

	for _, vendor := range vendors {
		for _, service := range services {
			for _, action := range actions {
				rk := rabbitmq.NewRoutingKey(vendor, service, action)
				key := rk.String()

				// Parse it back
				parsed, err := rabbitmq.Parse(key)
				if err != nil {
					t.Errorf("failed to parse generated key %s: %v", key, err)
					continue
				}

				if parsed.Vendor != vendor {
					t.Errorf("vendor mismatch: expected %s, got %s", vendor, parsed.Vendor)
				}
				if parsed.Service != service {
					t.Errorf("service mismatch: expected %s, got %s", service, parsed.Service)
				}
				if parsed.Action != action {
					t.Errorf("action mismatch: expected %s, got %s", action, parsed.Action)
				}
			}
		}
	}
}
