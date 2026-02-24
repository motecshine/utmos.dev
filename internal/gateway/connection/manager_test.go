package connection

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	manager := NewManager(nil)
	require.NotNil(t, manager)
	assert.NotNil(t, manager.devices)
	assert.NotNil(t, manager.logger)
}

func TestManager_Connect(t *testing.T) {
	manager := NewManager(nil)

	state := manager.Connect("device-001", "client-001", "192.168.1.100")

	require.NotNil(t, state)
	assert.Equal(t, "device-001", state.DeviceSN)
	assert.Equal(t, "client-001", state.ClientID)
	assert.Equal(t, "192.168.1.100", state.IPAddress)
	assert.True(t, state.Online)
	assert.NotNil(t, state.ConnectedAt)
}

func TestManager_Disconnect(t *testing.T) {
	manager := NewManager(nil)

	// Connect first
	manager.Connect("device-001", "client-001", "192.168.1.100")

	// Then disconnect
	state := manager.Disconnect("device-001")

	require.NotNil(t, state)
	assert.False(t, state.Online)
	assert.NotNil(t, state.DisconnectAt)
}

func TestManager_Disconnect_NotFound(t *testing.T) {
	manager := NewManager(nil)

	state := manager.Disconnect("nonexistent")
	assert.Nil(t, state)
}

func TestManager_GetState(t *testing.T) {
	manager := NewManager(nil)

	t.Run("existing device", func(t *testing.T) {
		manager.Connect("device-001", "client-001", "192.168.1.100")
		state := manager.GetState("device-001")
		require.NotNil(t, state)
		assert.Equal(t, "device-001", state.DeviceSN)
	})

	t.Run("non-existing device", func(t *testing.T) {
		state := manager.GetState("nonexistent")
		assert.Nil(t, state)
	})
}

func TestManager_IsOnline(t *testing.T) {
	manager := NewManager(nil)

	// Initially not online
	assert.False(t, manager.IsOnline("device-001"))

	// Connect
	manager.Connect("device-001", "client-001", "192.168.1.100")
	assert.True(t, manager.IsOnline("device-001"))

	// Disconnect
	manager.Disconnect("device-001")
	assert.False(t, manager.IsOnline("device-001"))
}

func TestManager_UpdateLastSeen(t *testing.T) {
	manager := NewManager(nil)

	manager.Connect("device-001", "client-001", "192.168.1.100")
	initialState := manager.GetState("device-001")
	initialLastSeen := initialState.LastSeenAt

	time.Sleep(10 * time.Millisecond)
	manager.UpdateLastSeen("device-001")

	updatedState := manager.GetState("device-001")
	assert.True(t, updatedState.LastSeenAt.After(initialLastSeen))
}

func TestManager_GetOnlineDevices(t *testing.T) {
	manager := NewManager(nil)

	// Connect multiple devices
	manager.Connect("device-001", "client-001", "192.168.1.100")
	manager.Connect("device-002", "client-002", "192.168.1.101")
	manager.Connect("device-003", "client-003", "192.168.1.102")

	// Disconnect one
	manager.Disconnect("device-002")

	online := manager.GetOnlineDevices()
	assert.Len(t, online, 2)
}

func TestManager_GetOnlineCount(t *testing.T) {
	manager := NewManager(nil)

	assert.Equal(t, 0, manager.GetOnlineCount())

	manager.Connect("device-001", "client-001", "192.168.1.100")
	assert.Equal(t, 1, manager.GetOnlineCount())

	manager.Connect("device-002", "client-002", "192.168.1.101")
	assert.Equal(t, 2, manager.GetOnlineCount())

	manager.Disconnect("device-001")
	assert.Equal(t, 1, manager.GetOnlineCount())
}

func TestManager_GetAllDevices(t *testing.T) {
	manager := NewManager(nil)

	manager.Connect("device-001", "client-001", "192.168.1.100")
	manager.Connect("device-002", "client-002", "192.168.1.101")

	devices := manager.GetAllDevices()
	assert.Len(t, devices, 2)
}

func TestManager_Remove(t *testing.T) {
	manager := NewManager(nil)

	manager.Connect("device-001", "client-001", "192.168.1.100")
	assert.NotNil(t, manager.GetState("device-001"))

	manager.Remove("device-001")
	assert.Nil(t, manager.GetState("device-001"))
}

func TestManager_Callbacks(t *testing.T) {
	manager := NewManager(nil)

	var connectCalled int32
	var disconnectCalled int32

	manager.SetOnConnect(func(state *DeviceState) {
		atomic.AddInt32(&connectCalled, 1)
	})

	manager.SetOnDisconnect(func(state *DeviceState) {
		atomic.AddInt32(&disconnectCalled, 1)
	})

	manager.Connect("device-001", "client-001", "192.168.1.100")
	time.Sleep(50 * time.Millisecond) // Wait for goroutine
	assert.Equal(t, int32(1), atomic.LoadInt32(&connectCalled))

	manager.Disconnect("device-001")
	time.Sleep(50 * time.Millisecond) // Wait for goroutine
	assert.Equal(t, int32(1), atomic.LoadInt32(&disconnectCalled))
}

func TestManager_CleanupStale(t *testing.T) {
	manager := NewManager(nil)

	// Connect and immediately disconnect
	manager.Connect("device-001", "client-001", "192.168.1.100")
	manager.Disconnect("device-001")

	// Should not be removed yet (too recent)
	removed := manager.CleanupStale(context.Background(), 1*time.Hour)
	assert.Equal(t, 0, removed)
	assert.NotNil(t, manager.GetState("device-001"))

	// Should be removed with very short max age
	removed = manager.CleanupStale(context.Background(), 0)
	assert.Equal(t, 1, removed)
	assert.Nil(t, manager.GetState("device-001"))
}

func TestManager_CleanupStale_OnlineNotRemoved(t *testing.T) {
	manager := NewManager(nil)

	// Connect device (stays online)
	manager.Connect("device-001", "client-001", "192.168.1.100")

	// Should not be removed even with 0 max age (because it's online)
	removed := manager.CleanupStale(context.Background(), 0)
	assert.Equal(t, 0, removed)
	assert.NotNil(t, manager.GetState("device-001"))
}
