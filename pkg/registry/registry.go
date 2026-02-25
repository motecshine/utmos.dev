// Package registry provides a generic thread-safe registry for vendor-based items
package registry

import (
	"sync"

	"github.com/sirupsen/logrus"
)

// VendorProvider is an interface that items must implement to be stored in the registry
type VendorProvider interface {
	// GetVendor returns the vendor identifier for this item
	GetVendor() string
}

// Registry is a generic thread-safe registry for vendor-based items
type Registry[T VendorProvider] struct {
	items  map[string]T
	mu     sync.RWMutex
	logger *logrus.Entry
	name   string
}

// New creates a new Registry with the given name and logger
func New[T VendorProvider](name string, logger *logrus.Entry) *Registry[T] {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &Registry[T]{
		items:  make(map[string]T),
		logger: logger.WithField("component", name),
		name:   name,
	}
}

// Register registers an item in the registry
func (r *Registry[T]) Register(item T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	vendor := item.GetVendor()
	r.items[vendor] = item
	r.logger.WithField("vendor", vendor).Info("Registered item")
}

// Unregister removes an item from the registry by vendor
func (r *Registry[T]) Unregister(vendor string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.items, vendor)
	r.logger.WithField("vendor", vendor).Info("Unregistered item")
}

// Get returns an item by vendor
func (r *Registry[T]) Get(vendor string) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[vendor]
	return item, exists
}

// ListVendors returns all registered vendor names
func (r *Registry[T]) ListVendors() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	vendors := make([]string, 0, len(r.items))
	for vendor := range r.items {
		vendors = append(vendors, vendor)
	}
	return vendors
}

// Count returns the number of registered items
func (r *Registry[T]) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.items)
}

// Range iterates over all items in the registry
// The callback function should return true to continue iteration, false to stop
func (r *Registry[T]) Range(fn func(vendor string, item T) bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for vendor, item := range r.items {
		if !fn(vendor, item) {
			break
		}
	}
}

// Find returns the first item that matches the predicate
func (r *Registry[T]) Find(predicate func(T) bool) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.items {
		if predicate(item) {
			return item, true
		}
	}

	var zero T
	return zero, false
}

// GetOrFind returns an item by vendor, or finds one using the predicate if not found
func (r *Registry[T]) GetOrFind(vendor string, predicate func(T) bool) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// First try direct vendor lookup
	if vendor != "" {
		if item, exists := r.items[vendor]; exists {
			return item, true
		}
	}

	// Then try predicate
	if predicate != nil {
		for _, item := range r.items {
			if predicate(item) {
				return item, true
			}
		}
	}

	var zero T
	return zero, false
}
