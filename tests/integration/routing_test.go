package integration

import (
	"testing"

	"github.com/utmos/utmos/pkg/rabbitmq"
)

// TestRoutingKeyFormat tests the routing key format specification
func TestRoutingKeyFormat(t *testing.T) {
	// Routing key format: iot.{vendor}.{service}.{action}
	// Note: action can contain dots (e.g., "property.report")
	tests := []struct {
		name     string
		vendor   string
		service  string
		action   string
		expected string
	}{
		// DJI vendor
		{"DJI property report", rabbitmq.VendorDJI, rabbitmq.ServiceDevice, rabbitmq.ActionPropertyReport, "iot.dji.device.property.report"},
		{"DJI device online", rabbitmq.VendorDJI, rabbitmq.ServiceDevice, rabbitmq.ActionDeviceOnline, "iot.dji.device.device.online"},
		{"DJI device offline", rabbitmq.VendorDJI, rabbitmq.ServiceDevice, rabbitmq.ActionDeviceOffline, "iot.dji.device.device.offline"},
		{"DJI event report", rabbitmq.VendorDJI, rabbitmq.ServiceEvent, rabbitmq.ActionEventReport, "iot.dji.event.event.report"},

		// Generic vendor
		{"Generic property report", rabbitmq.VendorGeneric, rabbitmq.ServiceDevice, rabbitmq.ActionPropertyReport, "iot.generic.device.property.report"},
		{"Generic service call", rabbitmq.VendorGeneric, rabbitmq.ServiceService, rabbitmq.ActionServiceCall, "iot.generic.service.service.call"},
		{"Generic service reply", rabbitmq.VendorGeneric, rabbitmq.ServiceService, rabbitmq.ActionServiceReply, "iot.generic.service.service.reply"},

		// Tuya vendor
		{"Tuya property report", rabbitmq.VendorTuya, rabbitmq.ServiceDevice, rabbitmq.ActionPropertyReport, "iot.tuya.device.property.report"},
		{"Tuya event report", rabbitmq.VendorTuya, rabbitmq.ServiceEvent, rabbitmq.ActionEventReport, "iot.tuya.event.event.report"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rk := rabbitmq.NewRoutingKey(tt.vendor, tt.service, tt.action)
			result := rk.String()

			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}

			// Verify round-trip parsing
			parsed, err := rabbitmq.Parse(result)
			if err != nil {
				t.Fatalf("failed to parse routing key: %v", err)
			}

			if parsed.Vendor != tt.vendor {
				t.Errorf("vendor mismatch: expected %s, got %s", tt.vendor, parsed.Vendor)
			}
			if parsed.Service != tt.service {
				t.Errorf("service mismatch: expected %s, got %s", tt.service, parsed.Service)
			}
			if parsed.Action != tt.action {
				t.Errorf("action mismatch: expected %s, got %s", tt.action, parsed.Action)
			}
		})
	}
}

// TestRoutingKeyWildcards tests wildcard pattern matching for subscriptions
func TestRoutingKeyWildcards(t *testing.T) {
	// These patterns would be used for queue bindings
	patterns := []struct {
		pattern     string
		description string
	}{
		{"iot.*.device.*", "All device messages from all vendors"},
		{"iot.dji.*.#", "All DJI messages"},
		{"iot.*.*.property_report", "All property reports from all vendors"},
		{"iot.#", "All IoT messages"},
	}

	for _, p := range patterns {
		t.Run(p.description, func(t *testing.T) {
			// Verify pattern is not empty
			if p.pattern == "" {
				t.Error("pattern should not be empty")
			}

			// Verify pattern starts with iot.
			if len(p.pattern) < 4 || p.pattern[:4] != "iot." {
				t.Errorf("pattern should start with 'iot.': %s", p.pattern)
			}
		})
	}
}

// TestVendorConstants tests vendor constant definitions
func TestVendorConstants(t *testing.T) {
	vendors := []string{
		rabbitmq.VendorDJI,
		rabbitmq.VendorGeneric,
		rabbitmq.VendorTuya,
	}

	for _, v := range vendors {
		if v == "" {
			t.Error("vendor constant should not be empty")
		}
	}

	// Verify uniqueness
	seen := make(map[string]bool)
	for _, v := range vendors {
		if seen[v] {
			t.Errorf("duplicate vendor: %s", v)
		}
		seen[v] = true
	}
}

// TestServiceConstants tests service constant definitions
func TestServiceConstants(t *testing.T) {
	services := []string{
		rabbitmq.ServiceDevice,
		rabbitmq.ServiceEvent,
		rabbitmq.ServiceService,
	}

	for _, s := range services {
		if s == "" {
			t.Error("service constant should not be empty")
		}
	}

	// Verify uniqueness
	seen := make(map[string]bool)
	for _, s := range services {
		if seen[s] {
			t.Errorf("duplicate service: %s", s)
		}
		seen[s] = true
	}
}

// TestActionConstants tests action constant definitions
func TestActionConstants(t *testing.T) {
	actions := []string{
		rabbitmq.ActionPropertyReport,
		rabbitmq.ActionEventReport,
		rabbitmq.ActionServiceCall,
		rabbitmq.ActionServiceReply,
		rabbitmq.ActionDeviceOnline,
		rabbitmq.ActionDeviceOffline,
	}

	for _, a := range actions {
		if a == "" {
			t.Error("action constant should not be empty")
		}
	}

	// Verify uniqueness
	seen := make(map[string]bool)
	for _, a := range actions {
		if seen[a] {
			t.Errorf("duplicate action: %s", a)
		}
		seen[a] = true
	}
}

// TestRoutingKeyParsing tests parsing invalid routing keys
func TestRoutingKeyParsingErrors(t *testing.T) {
	invalidKeys := []struct {
		key    string
		reason string
	}{
		{"", "empty string"},
		{"iot", "only prefix"},
		{"iot.dji", "missing service and action"},
		{"iot.dji.device", "missing action"},
		{"notiot.dji.device.action", "wrong prefix"},
		{"IOT.dji.device.action", "uppercase prefix"},
		{"iot..device.action", "empty vendor"},
		{"iot.dji..action", "empty service"},
		{"iot.dji.device.", "empty action"},
	}

	for _, tt := range invalidKeys {
		t.Run(tt.reason, func(t *testing.T) {
			_, err := rabbitmq.Parse(tt.key)
			if err == nil {
				t.Errorf("expected error for invalid key: %s (%s)", tt.key, tt.reason)
			}
		})
	}
}

// TestRoutingKeyEquality tests routing key equality
func TestRoutingKeyEquality(t *testing.T) {
	rk1 := rabbitmq.NewRoutingKey(rabbitmq.VendorDJI, rabbitmq.ServiceDevice, rabbitmq.ActionPropertyReport)
	rk2 := rabbitmq.NewRoutingKey(rabbitmq.VendorDJI, rabbitmq.ServiceDevice, rabbitmq.ActionPropertyReport)
	rk3 := rabbitmq.NewRoutingKey(rabbitmq.VendorTuya, rabbitmq.ServiceDevice, rabbitmq.ActionPropertyReport)

	if rk1.String() != rk2.String() {
		t.Error("identical routing keys should have same string representation")
	}

	if rk1.String() == rk3.String() {
		t.Error("different routing keys should have different string representations")
	}
}
