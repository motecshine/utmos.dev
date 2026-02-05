package adapter

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/pkg/rabbitmq"
)

// mockAdapter is a test implementation of ProtocolAdapter.
type mockAdapter struct {
	vendor string
}

func newMockAdapter(vendor string) *mockAdapter {
	return &mockAdapter{vendor: vendor}
}

func (m *mockAdapter) GetVendor() string {
	return m.vendor
}

func (m *mockAdapter) ParseRawMessage(topic string, payload []byte) (*ProtocolMessage, error) {
	return &ProtocolMessage{
		Vendor:      m.vendor,
		Topic:       topic,
		DeviceSN:    "mock-device",
		MessageType: MessageTypeProperty,
		Data:        payload,
	}, nil
}

func (m *mockAdapter) ToStandardMessage(pm *ProtocolMessage) (*rabbitmq.StandardMessage, error) {
	return &rabbitmq.StandardMessage{
		TID:      "mock-tid",
		BID:      "mock-bid",
		Service:  "mock-service",
		Action:   "property.report",
		DeviceSN: pm.DeviceSN,
		Data:     pm.Data,
	}, nil
}

func (m *mockAdapter) FromStandardMessage(sm *rabbitmq.StandardMessage) (*ProtocolMessage, error) {
	return &ProtocolMessage{
		Vendor:      m.vendor,
		DeviceSN:    sm.DeviceSN,
		MessageType: MessageTypeService,
		Data:        sm.Data,
	}, nil
}

func (m *mockAdapter) GetRawPayload(pm *ProtocolMessage) ([]byte, error) {
	return json.Marshal(pm.Data)
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()

	adapter := newMockAdapter("test-vendor")
	registry.Register(adapter)

	// Verify registration
	retrieved, err := registry.Get("test-vendor")
	require.NoError(t, err)
	assert.Equal(t, "test-vendor", retrieved.GetVendor())
}

func TestRegistry_Get_NotFound(t *testing.T) {
	registry := NewRegistry()

	_, err := registry.Get("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRegistry_Get_Found(t *testing.T) {
	registry := NewRegistry()

	adapter := newMockAdapter("dji")
	registry.Register(adapter)

	retrieved, err := registry.Get("dji")
	require.NoError(t, err)
	assert.Equal(t, "dji", retrieved.GetVendor())
}

func TestRegistry_List(t *testing.T) {
	registry := NewRegistry()

	// Empty registry
	vendors := registry.List()
	assert.Empty(t, vendors)

	// Register adapters
	registry.Register(newMockAdapter("dji"))
	registry.Register(newMockAdapter("tuya"))
	registry.Register(newMockAdapter("generic"))

	vendors = registry.List()
	assert.Len(t, vendors, 3)
	assert.Contains(t, vendors, "dji")
	assert.Contains(t, vendors, "tuya")
	assert.Contains(t, vendors, "generic")
}

func TestRegistry_Unregister(t *testing.T) {
	registry := NewRegistry()

	adapter := newMockAdapter("dji")
	registry.Register(adapter)

	// Verify registered
	_, err := registry.Get("dji")
	require.NoError(t, err)

	// Unregister
	registry.Unregister("dji")

	// Verify unregistered
	_, err = registry.Get("dji")
	assert.Error(t, err)
}

func TestRegistry_RegisterOverwrite(t *testing.T) {
	registry := NewRegistry()

	adapter1 := newMockAdapter("dji")
	adapter2 := newMockAdapter("dji")

	registry.Register(adapter1)
	registry.Register(adapter2)

	// Should have only one entry
	vendors := registry.List()
	assert.Len(t, vendors, 1)
}

func TestGlobalRegistry(t *testing.T) {
	// Reset global registry for test
	globalRegistry = NewRegistry()

	adapter := newMockAdapter("global-test")
	Register(adapter)

	retrieved, err := Get("global-test")
	require.NoError(t, err)
	assert.Equal(t, "global-test", retrieved.GetVendor())

	vendors := List()
	assert.Contains(t, vendors, "global-test")

	// Cleanup
	Unregister("global-test")
}
