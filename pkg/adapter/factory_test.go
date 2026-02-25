package adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/pkg/rabbitmq"
)

// testAdapter is a mock adapter for testing the factory.
type testAdapter struct {
	vendor string
}

func (a *testAdapter) GetVendor() string {
	return a.vendor
}

func (a *testAdapter) ParseRawMessage(topic string, _ []byte) (*ProtocolMessage, error) {
	return &ProtocolMessage{
		Vendor:      a.vendor,
		Topic:       topic,
		DeviceSN:    "test-device",
		MessageType: MessageTypeProperty,
	}, nil
}

func (a *testAdapter) ToStandardMessage(pm *ProtocolMessage) (*rabbitmq.StandardMessage, error) {
	return &rabbitmq.StandardMessage{
		DeviceSN: pm.DeviceSN,
		Action:   "property.report",
	}, nil
}

func (a *testAdapter) FromStandardMessage(sm *rabbitmq.StandardMessage) (*ProtocolMessage, error) {
	return &ProtocolMessage{
		Vendor:      a.vendor,
		DeviceSN:    sm.DeviceSN,
		MessageType: MessageTypeProperty,
	}, nil
}

func (a *testAdapter) GetRawPayload(_ *ProtocolMessage) ([]byte, error) {
	return []byte(`{}`), nil
}

func TestNewFactory(t *testing.T) {
	factory := NewFactory()
	assert.NotNil(t, factory)
	assert.NotNil(t, factory.registry)
}

func TestNewFactoryWithRegistry(t *testing.T) {
	registry := NewRegistry()
	factory := NewFactoryWithRegistry(registry)
	assert.NotNil(t, factory)
	assert.Equal(t, registry, factory.registry)
}

func TestAdapterFactory_NewAdapter(t *testing.T) {
	registry := NewRegistry()
	factory := NewFactoryWithRegistry(registry)

	// Register a test adapter
	testVendor := "test-vendor"
	registry.Register(&testAdapter{vendor: testVendor})

	t.Run("success", func(t *testing.T) {
		adapter, err := factory.NewAdapter(testVendor)
		require.NoError(t, err)
		assert.NotNil(t, adapter)
		assert.Equal(t, testVendor, adapter.GetVendor())
	})

	t.Run("unknown vendor", func(t *testing.T) {
		adapter, err := factory.NewAdapter("unknown-vendor")
		assert.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "unknown-vendor")
		assert.Contains(t, err.Error(), "failed to create adapter")
	})

	t.Run("empty vendor", func(t *testing.T) {
		adapter, err := factory.NewAdapter("")
		assert.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "vendor name cannot be empty")
	})
}

func TestAdapterFactory_ListVendors(t *testing.T) {
	registry := NewRegistry()
	factory := NewFactoryWithRegistry(registry)

	t.Run("empty registry", func(t *testing.T) {
		vendors := factory.ListVendors()
		assert.Empty(t, vendors)
	})

	t.Run("with registered adapters", func(t *testing.T) {
		registry.Register(&testAdapter{vendor: "vendor-a"})
		registry.Register(&testAdapter{vendor: "vendor-b"})

		vendors := factory.ListVendors()
		assert.Len(t, vendors, 2)
		assert.Contains(t, vendors, "vendor-a")
		assert.Contains(t, vendors, "vendor-b")
	})
}

func TestAdapterFactory_HasVendor(t *testing.T) {
	registry := NewRegistry()
	factory := NewFactoryWithRegistry(registry)

	testVendor := "test-vendor"
	registry.Register(&testAdapter{vendor: testVendor})

	t.Run("registered vendor", func(t *testing.T) {
		assert.True(t, factory.HasVendor(testVendor))
	})

	t.Run("unregistered vendor", func(t *testing.T) {
		assert.False(t, factory.HasVendor("unknown-vendor"))
	})
}

func TestNewAdapterByVendor(t *testing.T) {
	// Register a test adapter to the global registry
	testVendor := "global-test-vendor"
	Register(&testAdapter{vendor: testVendor})
	defer Unregister(testVendor)

	t.Run("success", func(t *testing.T) {
		adapter, err := NewAdapterByVendor(testVendor)
		require.NoError(t, err)
		assert.NotNil(t, adapter)
		assert.Equal(t, testVendor, adapter.GetVendor())
	})

	t.Run("unknown vendor", func(t *testing.T) {
		adapter, err := NewAdapterByVendor("nonexistent-vendor")
		assert.Error(t, err)
		assert.Nil(t, adapter)
	})
}

func TestListAvailableVendors(t *testing.T) {
	// Clean up any existing registrations
	for _, v := range List() {
		Unregister(v)
	}

	// Register test adapters
	Register(&testAdapter{vendor: "vendor-x"})
	Register(&testAdapter{vendor: "vendor-y"})
	defer func() {
		Unregister("vendor-x")
		Unregister("vendor-y")
	}()

	vendors := ListAvailableVendors()
	assert.Len(t, vendors, 2)
	assert.Contains(t, vendors, "vendor-x")
	assert.Contains(t, vendors, "vendor-y")
}

func TestIsVendorAvailable(t *testing.T) {
	testVendor := "availability-test-vendor"
	Register(&testAdapter{vendor: testVendor})
	defer Unregister(testVendor)

	assert.True(t, IsVendorAvailable(testVendor))
	assert.False(t, IsVendorAvailable("nonexistent-vendor"))
}
