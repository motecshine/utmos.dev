package registry

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockItem implements VendorProvider for testing
type mockItem struct {
	vendor string
	value  string
}

func (m *mockItem) GetVendor() string {
	return m.vendor
}

func TestNew(t *testing.T) {
	reg := New[*mockItem]("test-registry", nil)

	require.NotNil(t, reg)
	assert.NotNil(t, reg.items)
	assert.Equal(t, "test-registry", reg.name)
}

func TestRegistry_Register(t *testing.T) {
	reg := New[*mockItem]("test-registry", nil)

	item := &mockItem{vendor: "vendor1", value: "value1"}
	reg.Register(item)

	assert.Equal(t, 1, reg.Count())

	retrieved, exists := reg.Get("vendor1")
	assert.True(t, exists)
	assert.Equal(t, "value1", retrieved.value)
}

func TestRegistry_Unregister(t *testing.T) {
	reg := New[*mockItem]("test-registry", nil)

	item := &mockItem{vendor: "vendor1", value: "value1"}
	reg.Register(item)
	assert.Equal(t, 1, reg.Count())

	reg.Unregister("vendor1")
	assert.Equal(t, 0, reg.Count())

	_, exists := reg.Get("vendor1")
	assert.False(t, exists)
}

func TestRegistry_Get(t *testing.T) {
	reg := New[*mockItem]("test-registry", nil)

	item := &mockItem{vendor: "vendor1", value: "value1"}
	reg.Register(item)

	t.Run("existing item", func(t *testing.T) {
		retrieved, exists := reg.Get("vendor1")
		assert.True(t, exists)
		assert.Equal(t, "value1", retrieved.value)
	})

	t.Run("non-existing item", func(t *testing.T) {
		_, exists := reg.Get("vendor2")
		assert.False(t, exists)
	})
}

func TestRegistry_ListVendors(t *testing.T) {
	reg := New[*mockItem]("test-registry", nil)

	reg.Register(&mockItem{vendor: "vendor1"})
	reg.Register(&mockItem{vendor: "vendor2"})
	reg.Register(&mockItem{vendor: "vendor3"})

	vendors := reg.ListVendors()
	assert.Len(t, vendors, 3)
	assert.Contains(t, vendors, "vendor1")
	assert.Contains(t, vendors, "vendor2")
	assert.Contains(t, vendors, "vendor3")
}

func TestRegistry_Count(t *testing.T) {
	reg := New[*mockItem]("test-registry", nil)

	assert.Equal(t, 0, reg.Count())

	reg.Register(&mockItem{vendor: "vendor1"})
	assert.Equal(t, 1, reg.Count())

	reg.Register(&mockItem{vendor: "vendor2"})
	assert.Equal(t, 2, reg.Count())

	reg.Unregister("vendor1")
	assert.Equal(t, 1, reg.Count())
}

func TestRegistry_Range(t *testing.T) {
	reg := New[*mockItem]("test-registry", nil)

	reg.Register(&mockItem{vendor: "vendor1", value: "v1"})
	reg.Register(&mockItem{vendor: "vendor2", value: "v2"})
	reg.Register(&mockItem{vendor: "vendor3", value: "v3"})

	visited := make(map[string]string)
	reg.Range(func(vendor string, item *mockItem) bool {
		visited[vendor] = item.value
		return true
	})

	assert.Len(t, visited, 3)
	assert.Equal(t, "v1", visited["vendor1"])
	assert.Equal(t, "v2", visited["vendor2"])
	assert.Equal(t, "v3", visited["vendor3"])
}

func TestRegistry_Range_EarlyStop(t *testing.T) {
	reg := New[*mockItem]("test-registry", nil)

	reg.Register(&mockItem{vendor: "vendor1"})
	reg.Register(&mockItem{vendor: "vendor2"})
	reg.Register(&mockItem{vendor: "vendor3"})

	count := 0
	reg.Range(func(vendor string, item *mockItem) bool {
		count++
		return count < 2 // Stop after 2 iterations
	})

	assert.Equal(t, 2, count)
}

func TestRegistry_Find(t *testing.T) {
	reg := New[*mockItem]("test-registry", nil)

	reg.Register(&mockItem{vendor: "vendor1", value: "v1"})
	reg.Register(&mockItem{vendor: "vendor2", value: "v2"})
	reg.Register(&mockItem{vendor: "vendor3", value: "target"})

	t.Run("found", func(t *testing.T) {
		item, found := reg.Find(func(m *mockItem) bool {
			return m.value == "target"
		})
		assert.True(t, found)
		assert.Equal(t, "target", item.value)
	})

	t.Run("not found", func(t *testing.T) {
		_, found := reg.Find(func(m *mockItem) bool {
			return m.value == "nonexistent"
		})
		assert.False(t, found)
	})
}

func TestRegistry_GetOrFind(t *testing.T) {
	reg := New[*mockItem]("test-registry", nil)

	reg.Register(&mockItem{vendor: "vendor1", value: "v1"})
	reg.Register(&mockItem{vendor: "vendor2", value: "v2"})

	t.Run("get by vendor", func(t *testing.T) {
		item, found := reg.GetOrFind("vendor1", nil)
		assert.True(t, found)
		assert.Equal(t, "v1", item.value)
	})

	t.Run("find by predicate", func(t *testing.T) {
		item, found := reg.GetOrFind("", func(m *mockItem) bool {
			return m.value == "v2"
		})
		assert.True(t, found)
		assert.Equal(t, "v2", item.value)
	})

	t.Run("vendor takes precedence", func(t *testing.T) {
		item, found := reg.GetOrFind("vendor1", func(m *mockItem) bool {
			return m.value == "v2"
		})
		assert.True(t, found)
		assert.Equal(t, "v1", item.value) // vendor1, not v2
	})

	t.Run("not found", func(t *testing.T) {
		_, found := reg.GetOrFind("nonexistent", func(m *mockItem) bool {
			return m.value == "nonexistent"
		})
		assert.False(t, found)
	})
}

func TestRegistry_ConcurrentAccess(t *testing.T) {
	reg := New[*mockItem]("test-registry", nil)

	var wg sync.WaitGroup
	itemCount := 100

	// Concurrent registrations
	for i := 0; i < itemCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			reg.Register(&mockItem{
				vendor: fmt.Sprintf("vendor-%c", 'a'+id%26),
				value:  "value",
			})
		}(i)
	}

	wg.Wait()

	// Concurrent reads
	for i := 0; i < itemCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			reg.Get(fmt.Sprintf("vendor-%c", 'a'+id%26))
			reg.ListVendors()
			reg.Count()
		}(i)
	}

	wg.Wait()

	// Should not panic and should have some items
	assert.Greater(t, reg.Count(), 0)
}
