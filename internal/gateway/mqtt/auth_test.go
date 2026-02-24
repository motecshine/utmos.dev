package mqtt

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/utmos/utmos/internal/gateway/model"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.DeviceCredential{})
	require.NoError(t, err)

	return db
}

func TestNewAuthenticator(t *testing.T) {
	db := setupTestDB(t)
	auth := NewAuthenticator(db, nil)

	assert.NotNil(t, auth)
	assert.NotNil(t, auth.db)
	assert.NotNil(t, auth.logger)
}

func TestAuthenticator_CreateCredential(t *testing.T) {
	db := setupTestDB(t)
	auth := NewAuthenticator(db, nil)
	ctx := context.Background()

	credential, err := auth.CreateCredential(ctx, "device-001", "user001", "password123")
	require.NoError(t, err)
	assert.NotNil(t, credential)
	assert.Equal(t, "device-001", credential.DeviceSN)
	assert.Equal(t, "user001", credential.Username)
	assert.True(t, credential.Enabled)
	assert.NotEmpty(t, credential.PasswordHash)

	// Verify password hash
	err = bcrypt.CompareHashAndPassword([]byte(credential.PasswordHash), []byte("password123"))
	assert.NoError(t, err)
}

func TestAuthenticator_Authenticate(t *testing.T) {
	db := setupTestDB(t)
	auth := NewAuthenticator(db, nil)
	ctx := context.Background()

	// Create a credential first
	_, err := auth.CreateCredential(ctx, "device-001", "user001", "password123")
	require.NoError(t, err)

	t.Run("valid credentials", func(t *testing.T) {
		credential, err := auth.Authenticate(ctx, "user001", "password123")
		require.NoError(t, err)
		assert.NotNil(t, credential)
		assert.Equal(t, "device-001", credential.DeviceSN)
	})

	t.Run("invalid password", func(t *testing.T) {
		credential, err := auth.Authenticate(ctx, "user001", "wrongpassword")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Nil(t, credential)
	})

	t.Run("user not found", func(t *testing.T) {
		credential, err := auth.Authenticate(ctx, "nonexistent", "password123")
		assert.ErrorIs(t, err, ErrDeviceNotFound)
		assert.Nil(t, credential)
	})
}

func TestAuthenticator_AuthenticateByDeviceSN(t *testing.T) {
	db := setupTestDB(t)
	auth := NewAuthenticator(db, nil)
	ctx := context.Background()

	// Create a credential first
	_, err := auth.CreateCredential(ctx, "device-001", "user001", "password123")
	require.NoError(t, err)

	t.Run("valid credentials", func(t *testing.T) {
		credential, err := auth.AuthenticateByDeviceSN(ctx, "device-001", "password123")
		require.NoError(t, err)
		assert.NotNil(t, credential)
		assert.Equal(t, "device-001", credential.DeviceSN)
	})

	t.Run("invalid password", func(t *testing.T) {
		credential, err := auth.AuthenticateByDeviceSN(ctx, "device-001", "wrongpassword")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Nil(t, credential)
	})

	t.Run("device not found", func(t *testing.T) {
		credential, err := auth.AuthenticateByDeviceSN(ctx, "nonexistent", "password123")
		assert.ErrorIs(t, err, ErrDeviceNotFound)
		assert.Nil(t, credential)
	})
}

func TestAuthenticator_DisabledDevice(t *testing.T) {
	db := setupTestDB(t)
	auth := NewAuthenticator(db, nil)
	ctx := context.Background()

	// Create and disable a credential
	_, err := auth.CreateCredential(ctx, "device-001", "user001", "password123")
	require.NoError(t, err)

	err = auth.DisableDevice(ctx, "device-001")
	require.NoError(t, err)

	// Try to authenticate
	credential, err := auth.Authenticate(ctx, "user001", "password123")
	assert.ErrorIs(t, err, ErrDeviceDisabled)
	assert.Nil(t, credential)
}

func TestAuthenticator_EnableDevice(t *testing.T) {
	db := setupTestDB(t)
	auth := NewAuthenticator(db, nil)
	ctx := context.Background()

	// Create and disable a credential
	_, err := auth.CreateCredential(ctx, "device-001", "user001", "password123")
	require.NoError(t, err)

	err = auth.DisableDevice(ctx, "device-001")
	require.NoError(t, err)

	// Re-enable
	err = auth.EnableDevice(ctx, "device-001")
	require.NoError(t, err)

	// Should be able to authenticate now
	credential, err := auth.Authenticate(ctx, "user001", "password123")
	require.NoError(t, err)
	assert.NotNil(t, credential)
}

func TestAuthenticator_DeleteCredential(t *testing.T) {
	db := setupTestDB(t)
	auth := NewAuthenticator(db, nil)
	ctx := context.Background()

	// Create a credential
	_, err := auth.CreateCredential(ctx, "device-001", "user001", "password123")
	require.NoError(t, err)

	// Delete it
	err = auth.DeleteCredential(ctx, "device-001")
	require.NoError(t, err)

	// Should not be able to authenticate
	credential, err := auth.Authenticate(ctx, "user001", "password123")
	assert.ErrorIs(t, err, ErrDeviceNotFound)
	assert.Nil(t, credential)
}

func TestAuthenticator_GetCredential(t *testing.T) {
	db := setupTestDB(t)
	auth := NewAuthenticator(db, nil)
	ctx := context.Background()

	// Create a credential
	_, err := auth.CreateCredential(ctx, "device-001", "user001", "password123")
	require.NoError(t, err)

	t.Run("existing device", func(t *testing.T) {
		credential, err := auth.GetCredential(ctx, "device-001")
		require.NoError(t, err)
		assert.NotNil(t, credential)
		assert.Equal(t, "device-001", credential.DeviceSN)
	})

	t.Run("non-existing device", func(t *testing.T) {
		credential, err := auth.GetCredential(ctx, "nonexistent")
		assert.ErrorIs(t, err, ErrDeviceNotFound)
		assert.Nil(t, credential)
	})
}
