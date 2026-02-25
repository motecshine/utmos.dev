// Package connection provides device connection state management
package connection

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// DeviceState represents the connection state of a device
type DeviceState struct {
	DeviceSN     string
	Online       bool
	ConnectedAt  *time.Time
	LastSeenAt   time.Time
	DisconnectAt *time.Time
	ClientID     string
	IPAddress    string
}

// Manager manages device connection states
type Manager struct {
	devices map[string]*DeviceState
	mu      sync.RWMutex
	logger  *logrus.Entry

	// Callbacks
	onConnect    func(state *DeviceState)
	onDisconnect func(state *DeviceState)
}

// NewManager creates a new connection manager
func NewManager(logger *logrus.Entry) *Manager {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &Manager{
		devices: make(map[string]*DeviceState),
		logger:  logger.WithField("component", "connection-manager"),
	}
}

// SetOnConnect sets the callback for device connection events
func (m *Manager) SetOnConnect(callback func(state *DeviceState)) {
	m.onConnect = callback
}

// SetOnDisconnect sets the callback for device disconnection events
func (m *Manager) SetOnDisconnect(callback func(state *DeviceState)) {
	m.onDisconnect = callback
}

// Connect marks a device as connected
func (m *Manager) Connect(deviceSN, clientID, ipAddress string) *DeviceState {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	state := &DeviceState{
		DeviceSN:    deviceSN,
		Online:      true,
		ConnectedAt: &now,
		LastSeenAt:  now,
		ClientID:    clientID,
		IPAddress:   ipAddress,
	}

	m.devices[deviceSN] = state

	m.logger.WithFields(logrus.Fields{
		"device_sn":  deviceSN,
		"client_id":  clientID,
		"ip_address": ipAddress,
	}).Info("Device connected")

	if m.onConnect != nil {
		go m.onConnect(state)
	}

	return state
}

// Disconnect marks a device as disconnected
func (m *Manager) Disconnect(deviceSN string) *DeviceState {
	m.mu.Lock()
	defer m.mu.Unlock()

	state, exists := m.devices[deviceSN]
	if !exists {
		return nil
	}

	now := time.Now()
	state.Online = false
	state.DisconnectAt = &now

	m.logger.WithFields(logrus.Fields{
		"device_sn": deviceSN,
		"client_id": state.ClientID,
	}).Info("Device disconnected")

	if m.onDisconnect != nil {
		go m.onDisconnect(state)
	}

	return state
}

// UpdateLastSeen updates the last seen timestamp for a device
func (m *Manager) UpdateLastSeen(deviceSN string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if state, exists := m.devices[deviceSN]; exists {
		state.LastSeenAt = time.Now()
	}
}

// GetState returns the connection state for a device
func (m *Manager) GetState(deviceSN string) *DeviceState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if state, exists := m.devices[deviceSN]; exists {
		// Return a copy to prevent race conditions
		stateCopy := *state
		return &stateCopy
	}
	return nil
}

// IsOnline checks if a device is online
func (m *Manager) IsOnline(deviceSN string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if state, exists := m.devices[deviceSN]; exists {
		return state.Online
	}
	return false
}

// GetOnlineDevices returns a list of all online devices
func (m *Manager) GetOnlineDevices() []*DeviceState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var online []*DeviceState
	for _, state := range m.devices {
		if state.Online {
			stateCopy := *state
			online = append(online, &stateCopy)
		}
	}
	return online
}

// GetOnlineCount returns the number of online devices
func (m *Manager) GetOnlineCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, state := range m.devices {
		if state.Online {
			count++
		}
	}
	return count
}

// GetAllDevices returns all device states
func (m *Manager) GetAllDevices() []*DeviceState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	devices := make([]*DeviceState, 0, len(m.devices))
	for _, state := range m.devices {
		stateCopy := *state
		devices = append(devices, &stateCopy)
	}
	return devices
}

// Remove removes a device from the manager
func (m *Manager) Remove(deviceSN string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.devices, deviceSN)
}

// CleanupStale removes devices that haven't been seen for the specified duration
func (m *Manager) CleanupStale(ctx context.Context, maxAge time.Duration) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	threshold := time.Now().Add(-maxAge)
	removed := 0

	for deviceSN, state := range m.devices {
		if !state.Online && state.LastSeenAt.Before(threshold) {
			delete(m.devices, deviceSN)
			removed++
			m.logger.WithField("device_sn", deviceSN).Debug("Removed stale device state")
		}
	}

	return removed
}

// StartCleanupRoutine starts a background routine to clean up stale device states
func (m *Manager) StartCleanupRoutine(ctx context.Context, interval, maxAge time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				removed := m.CleanupStale(ctx, maxAge)
				if removed > 0 {
					m.logger.WithField("removed", removed).Debug("Cleaned up stale device states")
				}
			}
		}
	}()
}
