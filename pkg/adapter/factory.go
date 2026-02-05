package adapter

import (
	"fmt"
)

// Factory provides methods for creating protocol adapters.
type Factory struct {
	registry *Registry
}

// NewFactory creates a new adapter factory using the global registry.
func NewFactory() *Factory {
	return &Factory{
		registry: globalRegistry,
	}
}

// NewFactoryWithRegistry creates a new adapter factory with a custom registry.
func NewFactoryWithRegistry(registry *Registry) *Factory {
	return &Factory{
		registry: registry,
	}
}

// NewAdapter creates a protocol adapter by vendor name.
// Returns an error if the vendor is not registered.
func (f *Factory) NewAdapter(vendor string) (ProtocolAdapter, error) {
	if vendor == "" {
		return nil, fmt.Errorf("vendor name cannot be empty")
	}

	adapter, err := f.registry.Get(vendor)
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter for vendor %q: %w", vendor, err)
	}

	return adapter, nil
}

// ListVendors returns a list of all available vendor names.
func (f *Factory) ListVendors() []string {
	return f.registry.List()
}

// HasVendor checks if an adapter for the given vendor is registered.
func (f *Factory) HasVendor(vendor string) bool {
	_, err := f.registry.Get(vendor)
	return err == nil
}

// defaultFactory is the default global factory instance.
var defaultFactory = NewFactory()

// NewAdapterByVendor creates a protocol adapter by vendor name using the global factory.
// This is a convenience function for creating adapters without explicitly creating a factory.
func NewAdapterByVendor(vendor string) (ProtocolAdapter, error) {
	return defaultFactory.NewAdapter(vendor)
}

// ListAvailableVendors returns a list of all available vendor names from the global factory.
func ListAvailableVendors() []string {
	return defaultFactory.ListVendors()
}

// IsVendorAvailable checks if an adapter for the given vendor is available in the global factory.
func IsVendorAvailable(vendor string) bool {
	return defaultFactory.HasVendor(vendor)
}
