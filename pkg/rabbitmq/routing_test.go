package rabbitmq

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRoutingKey(t *testing.T) {
	tests := []struct {
		name     string
		vendor   string
		service  string
		action   string
		expected string
	}{
		{
			name:     "DJI property report",
			vendor:   "dji",
			service:  "gateway",
			action:   ActionPropertyReport,
			expected: "iot.dji.gateway.property.report",
		},
		{
			name:     "Generic device online",
			vendor:   VendorGeneric,
			service:  "uplink",
			action:   ActionDeviceOnline,
			expected: "iot.generic.uplink.device.online",
		},
		{
			name:     "Tuya service call",
			vendor:   "tuya",
			service:  "downlink",
			action:   ActionServiceCall,
			expected: "iot.tuya.downlink.service.call",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rk := NewRoutingKey(tt.vendor, tt.service, tt.action)
			assert.Equal(t, tt.expected, rk.String())
			assert.Equal(t, tt.vendor, rk.Vendor)
			assert.Equal(t, tt.service, rk.Service)
			assert.Equal(t, tt.action, rk.Action)
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		wantVendor  string
		wantService string
		wantAction  string
		wantErr     bool
	}{
		{
			name:        "valid routing key with single action",
			key:         "iot.dji.gateway.property.report",
			wantVendor:  "dji",
			wantService: "gateway",
			wantAction:  "property.report",
			wantErr:     false,
		},
		{
			name:        "valid routing key with device online",
			key:         "iot.generic.uplink.device.online",
			wantVendor:  "generic",
			wantService: "uplink",
			wantAction:  "device.online",
			wantErr:     false,
		},
		{
			name:        "valid routing key with service reply",
			key:         "iot.tuya.api.service.reply",
			wantVendor:  "tuya",
			wantService: "api",
			wantAction:  "service.reply",
			wantErr:     false,
		},
		{
			name:    "invalid - missing prefix",
			key:     "dji.gateway.property.report",
			wantErr: true,
		},
		{
			name:    "invalid - wrong prefix",
			key:     "mqtt.dji.gateway.property.report",
			wantErr: true,
		},
		{
			name:    "invalid - too few segments",
			key:     "iot.dji.gateway",
			wantErr: true,
		},
		{
			name:    "invalid - empty string",
			key:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rk, err := Parse(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantVendor, rk.Vendor)
			assert.Equal(t, tt.wantService, rk.Service)
			assert.Equal(t, tt.wantAction, rk.Action)
		})
	}
}

func TestRoutingKey_String(t *testing.T) {
	rk := RoutingKey{
		Vendor:  "dji",
		Service: "gateway",
		Action:  "property.report",
	}
	assert.Equal(t, "iot.dji.gateway.property.report", rk.String())
}

func TestRoutingKey_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rk      RoutingKey
		wantErr bool
	}{
		{
			name: "valid routing key",
			rk: RoutingKey{
				Vendor:  "dji",
				Service: "gateway",
				Action:  "property.report",
			},
			wantErr: false,
		},
		{
			name: "empty vendor",
			rk: RoutingKey{
				Vendor:  "",
				Service: "gateway",
				Action:  "property.report",
			},
			wantErr: true,
		},
		{
			name: "empty service",
			rk: RoutingKey{
				Vendor:  "dji",
				Service: "",
				Action:  "property.report",
			},
			wantErr: true,
		},
		{
			name: "empty action",
			rk: RoutingKey{
				Vendor:  "dji",
				Service: "gateway",
				Action:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rk.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewRawRoutingKey(t *testing.T) {
	tests := []struct {
		name      string
		vendor    string
		direction string
		expected  string
	}{
		{
			name:      "DJI uplink",
			vendor:    "dji",
			direction: DirectionUplink,
			expected:  "iot.raw.dji.uplink",
		},
		{
			name:      "DJI downlink",
			vendor:    "dji",
			direction: DirectionDownlink,
			expected:  "iot.raw.dji.downlink",
		},
		{
			name:      "Tuya uplink",
			vendor:    "tuya",
			direction: DirectionUplink,
			expected:  "iot.raw.tuya.uplink",
		},
		{
			name:      "Generic downlink",
			vendor:    VendorGeneric,
			direction: DirectionDownlink,
			expected:  "iot.raw.generic.downlink",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rk := NewRawRoutingKey(tt.vendor, tt.direction)
			assert.Equal(t, tt.expected, rk.String())
			assert.Equal(t, tt.vendor, rk.Vendor)
			assert.Equal(t, tt.direction, rk.Direction)
		})
	}
}

func TestRawRoutingKey_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rk      RawRoutingKey
		wantErr bool
	}{
		{
			name: "valid uplink",
			rk: RawRoutingKey{
				Vendor:    "dji",
				Direction: DirectionUplink,
			},
			wantErr: false,
		},
		{
			name: "valid downlink",
			rk: RawRoutingKey{
				Vendor:    "tuya",
				Direction: DirectionDownlink,
			},
			wantErr: false,
		},
		{
			name: "empty vendor",
			rk: RawRoutingKey{
				Vendor:    "",
				Direction: DirectionUplink,
			},
			wantErr: true,
		},
		{
			name: "invalid direction",
			rk: RawRoutingKey{
				Vendor:    "dji",
				Direction: "invalid",
			},
			wantErr: true,
		},
		{
			name: "empty direction",
			rk: RawRoutingKey{
				Vendor:    "dji",
				Direction: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rk.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildRawBindingPattern(t *testing.T) {
	tests := []struct {
		name      string
		vendor    string
		direction string
		expected  string
	}{
		{
			name:      "specific vendor and direction",
			vendor:    "dji",
			direction: "uplink",
			expected:  "iot.raw.dji.uplink",
		},
		{
			name:      "wildcard vendor",
			vendor:    "",
			direction: "uplink",
			expected:  "iot.raw.*.uplink",
		},
		{
			name:      "wildcard direction",
			vendor:    "dji",
			direction: "",
			expected:  "iot.raw.dji.*",
		},
		{
			name:      "all wildcards",
			vendor:    "",
			direction: "",
			expected:  "iot.raw.*.*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := BuildRawBindingPattern(tt.vendor, tt.direction)
			assert.Equal(t, tt.expected, pattern)
		})
	}
}

func TestBuildBindingPattern(t *testing.T) {
	tests := []struct {
		name     string
		vendor   string
		service  string
		action   string
		expected string
	}{
		{
			name:     "all specified",
			vendor:   "dji",
			service:  "device",
			action:   "property.report",
			expected: "iot.dji.device.property.report",
		},
		{
			name:     "wildcard vendor",
			vendor:   "",
			service:  "device",
			action:   "property.report",
			expected: "iot.*.device.property.report",
		},
		{
			name:     "wildcard service",
			vendor:   "dji",
			service:  "",
			action:   "property.report",
			expected: "iot.dji.*.property.report",
		},
		{
			name:     "wildcard action",
			vendor:   "dji",
			service:  "device",
			action:   "",
			expected: "iot.dji.device.#",
		},
		{
			name:     "all wildcards",
			vendor:   "",
			service:  "",
			action:   "",
			expected: "iot.*.*.#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := BuildBindingPattern(tt.vendor, tt.service, tt.action)
			assert.Equal(t, tt.expected, pattern)
		})
	}
}
