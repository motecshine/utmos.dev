package integration

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/utmos/utmos/internal/gateway/connection"
	gatewaymodel "github.com/utmos/utmos/internal/gateway/model"
	"github.com/utmos/utmos/internal/gateway/mqtt"
)

func setupGatewayTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&gatewaymodel.DeviceCredential{}))
	return db
}

// TestGatewayAuthenticationFlow tests device credential validation
// through the full authentication lifecycle: credential creation,
// successful authentication, rejection on bad password, disabled-device
// rejection, re-enable, and credential deletion.
func TestGatewayAuthenticationFlow(t *testing.T) {
	db := setupGatewayTestDB(t)
	auth := mqtt.NewAuthenticator(db, nil)
	ctx := context.Background()

	const (
		deviceSN = "INTEG-GW-001"
		username = "gw_user_001"
		password = "s3cret!"
	)

	// Step 1: Create credential
	cred, err := auth.CreateCredential(ctx, deviceSN, username, password)
	require.NoError(t, err)
	require.NotNil(t, cred)
	assert.Equal(t, deviceSN, cred.DeviceSN)
	assert.True(t, cred.Enabled)

	t.Run("authenticate with valid credentials", func(t *testing.T) {
		result, err := auth.Authenticate(ctx, username, password)
		require.NoError(t, err)
		assert.Equal(t, deviceSN, result.DeviceSN)
	})

	t.Run("authenticate by device SN", func(t *testing.T) {
		result, err := auth.AuthenticateByDeviceSN(ctx, deviceSN, password)
		require.NoError(t, err)
		assert.Equal(t, deviceSN, result.DeviceSN)
	})

	t.Run("reject wrong password", func(t *testing.T) {
		result, err := auth.Authenticate(ctx, username, "wrong")
		assert.ErrorIs(t, err, mqtt.ErrInvalidCredentials)
		assert.Nil(t, result)
	})

	t.Run("reject unknown device", func(t *testing.T) {
		result, err := auth.Authenticate(ctx, "no_such_user", password)
		assert.ErrorIs(t, err, mqtt.ErrDeviceNotFound)
		assert.Nil(t, result)
	})

	t.Run("reject disabled device", func(t *testing.T) {
		require.NoError(t, auth.DisableDevice(ctx, deviceSN))

		result, err := auth.Authenticate(ctx, username, password)
		assert.ErrorIs(t, err, mqtt.ErrDeviceDisabled)
		assert.Nil(t, result)

		// Also rejected via device SN path
		result, err = auth.AuthenticateByDeviceSN(ctx, deviceSN, password)
		assert.ErrorIs(t, err, mqtt.ErrDeviceDisabled)
		assert.Nil(t, result)
	})

	t.Run("re-enable allows authentication", func(t *testing.T) {
		require.NoError(t, auth.EnableDevice(ctx, deviceSN))

		result, err := auth.Authenticate(ctx, username, password)
		require.NoError(t, err)
		assert.Equal(t, deviceSN, result.DeviceSN)
	})

	t.Run("delete credential prevents authentication", func(t *testing.T) {
		require.NoError(t, auth.DeleteCredential(ctx, deviceSN))

		result, err := auth.Authenticate(ctx, username, password)
		assert.ErrorIs(t, err, mqtt.ErrDeviceNotFound)
		assert.Nil(t, result)
	})

	t.Run("multiple devices isolated", func(t *testing.T) {
		_, err := auth.CreateCredential(ctx, "DEV-A", "userA", "passA")
		require.NoError(t, err)
		_, err = auth.CreateCredential(ctx, "DEV-B", "userB", "passB")
		require.NoError(t, err)

		resultA, err := auth.Authenticate(ctx, "userA", "passA")
		require.NoError(t, err)
		assert.Equal(t, "DEV-A", resultA.DeviceSN)

		resultB, err := auth.Authenticate(ctx, "userB", "passB")
		require.NoError(t, err)
		assert.Equal(t, "DEV-B", resultB.DeviceSN)

		// Cross-credential must fail
		_, err = auth.Authenticate(ctx, "userA", "passB")
		assert.ErrorIs(t, err, mqtt.ErrInvalidCredentials)
	})
}

// TestDeviceOnlineOfflineTracking tests connection state management
// through the full online/offline lifecycle: connect, verify online,
// disconnect, verify offline, reconnect, stale cleanup, and callback
// notifications.
func TestDeviceOnlineOfflineTracking(t *testing.T) {
	mgr := connection.NewManager(nil)

	const (
		device1 = "CONN-DEV-001"
		device2 = "CONN-DEV-002"
	)

	t.Run("connect marks device online", func(t *testing.T) {
		mgr.Connect(device1, "client-001", "10.0.0.1")

		assert.True(t, mgr.IsOnline(device1))
		assert.Equal(t, 1, mgr.GetOnlineCount())

		state := mgr.GetState(device1)
		require.NotNil(t, state)
		assert.True(t, state.Online)
		assert.Equal(t, "client-001", state.ClientID)
		assert.Equal(t, "10.0.0.1", state.IPAddress)
		assert.NotNil(t, state.ConnectedAt)
	})

	t.Run("disconnect marks device offline", func(t *testing.T) {
		mgr.Disconnect(device1)

		assert.False(t, mgr.IsOnline(device1))
		assert.Equal(t, 0, mgr.GetOnlineCount())

		state := mgr.GetState(device1)
		require.NotNil(t, state)
		assert.False(t, state.Online)
		assert.NotNil(t, state.DisconnectAt)
	})

	t.Run("reconnect updates state", func(t *testing.T) {
		mgr.Connect(device1, "client-001-v2", "10.0.0.2")

		assert.True(t, mgr.IsOnline(device1))
		state := mgr.GetState(device1)
		require.NotNil(t, state)
		assert.Equal(t, "client-001-v2", state.ClientID)
		assert.Equal(t, "10.0.0.2", state.IPAddress)
	})

	t.Run("multiple devices tracked independently", func(t *testing.T) {
		mgr.Connect(device2, "client-002", "10.0.0.3")

		assert.True(t, mgr.IsOnline(device1))
		assert.True(t, mgr.IsOnline(device2))
		assert.Equal(t, 2, mgr.GetOnlineCount())

		onlineDevices := mgr.GetOnlineDevices()
		assert.Len(t, onlineDevices, 2)

		// Disconnect one, other stays online
		mgr.Disconnect(device1)
		assert.False(t, mgr.IsOnline(device1))
		assert.True(t, mgr.IsOnline(device2))
		assert.Equal(t, 1, mgr.GetOnlineCount())
	})

	t.Run("update last seen timestamp", func(t *testing.T) {
		state := mgr.GetState(device2)
		require.NotNil(t, state)
		prevLastSeen := state.LastSeenAt

		time.Sleep(10 * time.Millisecond)
		mgr.UpdateLastSeen(device2)

		state = mgr.GetState(device2)
		require.NotNil(t, state)
		assert.True(t, state.LastSeenAt.After(prevLastSeen))
	})

	t.Run("remove deletes device state", func(t *testing.T) {
		mgr.Remove(device2)
		assert.Nil(t, mgr.GetState(device2))
		assert.False(t, mgr.IsOnline(device2))
		assert.Equal(t, 0, mgr.GetOnlineCount())
	})

	t.Run("callbacks fire on connect and disconnect", func(t *testing.T) {
		callbackMgr := connection.NewManager(nil)
		var connectCount, disconnectCount int32

		callbackMgr.SetOnConnect(func(_ *connection.DeviceState) {
			atomic.AddInt32(&connectCount, 1)
		})
		callbackMgr.SetOnDisconnect(func(_ *connection.DeviceState) {
			atomic.AddInt32(&disconnectCount, 1)
		})

		callbackMgr.Connect("CB-DEV", "cb-client", "10.0.0.99")
		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, int32(1), atomic.LoadInt32(&connectCount))

		callbackMgr.Disconnect("CB-DEV")
		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, int32(1), atomic.LoadInt32(&disconnectCount))
	})

	t.Run("unknown device returns nil state", func(t *testing.T) {
		assert.Nil(t, mgr.GetState("NO-SUCH-DEVICE"))
		assert.False(t, mgr.IsOnline("NO-SUCH-DEVICE"))
	})
}
