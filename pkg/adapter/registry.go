package adapter

import (
	"fmt"
	"sort"
	"sync"
)

// Registry manages protocol adapter registration and lookup.
type Registry struct {
	mu       sync.RWMutex
	adapters map[string]ProtocolAdapter
}

// NewRegistry creates a new adapter registry.
func NewRegistry() *Registry {
	return &Registry{
		adapters: make(map[string]ProtocolAdapter),
	}
}

// Register registers a protocol adapter with the registry.
// If an adapter with the same vendor already exists, it will be overwritten.
func (r *Registry) Register(adapter ProtocolAdapter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adapters[adapter.GetVendor()] = adapter
}

// Unregister removes a protocol adapter from the registry.
func (r *Registry) Unregister(vendor string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.adapters, vendor)
}

// Get retrieves a protocol adapter by vendor name.
func (r *Registry) Get(vendor string) (ProtocolAdapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	adapter, ok := r.adapters[vendor]
	if !ok {
		return nil, fmt.Errorf("adapter for vendor %q not found", vendor)
	}
	return adapter, nil
}

// List returns a list of all registered vendor names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	vendors := make([]string, 0, len(r.adapters))
	for vendor := range r.adapters {
		vendors = append(vendors, vendor)
	}
	sort.Strings(vendors)
	return vendors
}

// globalRegistry is the default global registry instance.
var globalRegistry = NewRegistry()

// Register registers a protocol adapter with the global registry.
func Register(adapter ProtocolAdapter) {
	globalRegistry.Register(adapter)
}

// Unregister removes a protocol adapter from the global registry.
func Unregister(vendor string) {
	globalRegistry.Unregister(vendor)
}

// Get retrieves a protocol adapter from the global registry.
func Get(vendor string) (ProtocolAdapter, error) {
	return globalRegistry.Get(vendor)
}

// List returns all registered vendors from the global registry.
func List() []string {
	return globalRegistry.List()
}
