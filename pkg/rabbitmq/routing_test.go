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
			vendor:   VendorDJI,
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
			vendor:   VendorTuya,
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
