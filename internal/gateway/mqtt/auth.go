// Package mqtt provides MQTT authentication functionality
package mqtt

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/utmos/utmos/internal/gateway/model"
)

var (
	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrDeviceDisabled is returned when device is disabled
	ErrDeviceDisabled = errors.New("device is disabled")
	// ErrDeviceNotFound is returned when device is not found
	ErrDeviceNotFound = errors.New("device not found")
)

// Authenticator handles device authentication
type Authenticator struct {
	db     *gorm.DB
	logger *logrus.Entry
}

// NewAuthenticator creates a new Authenticator
func NewAuthenticator(db *gorm.DB, logger *logrus.Entry) *Authenticator {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &Authenticator{
		db:     db,
		logger: logger.WithField("component", "authenticator"),
	}
}

// findCredential looks up a DeviceCredential by a single column condition.
// It returns ErrDeviceNotFound if no record matches, or wraps other DB errors.
func (a *Authenticator) findCredential(ctx context.Context, field, value string) (*model.DeviceCredential, error) {
	var credential model.DeviceCredential

	result := a.db.WithContext(ctx).Where(field+" = ?", value).First(&credential)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			a.logger.WithField(field, value).Debug("Device not found")
			return nil, ErrDeviceNotFound
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	return &credential, nil
}

// verifyPassword compares a bcrypt hash with a plaintext password.
// It returns ErrInvalidCredentials when the password does not match.
func (a *Authenticator) verifyPassword(credential *model.DeviceCredential, password, logField, logValue string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(credential.PasswordHash), []byte(password)); err != nil {
		a.logger.WithField(logField, logValue).Debug("Invalid password")
		return ErrInvalidCredentials
	}
	return nil
}

// Authenticate validates device credentials
func (a *Authenticator) Authenticate(ctx context.Context, username, password string) (*model.DeviceCredential, error) {
	credential, err := a.findCredential(ctx, "username", username)
	if err != nil {
		return nil, err
	}

	if !credential.Enabled {
		a.logger.WithField("device_sn", credential.DeviceSN).Debug("Device is disabled")
		return nil, ErrDeviceDisabled
	}

	if err := a.verifyPassword(credential, password, "username", username); err != nil {
		return nil, err
	}

	a.logger.WithFields(logrus.Fields{
		"device_sn": credential.DeviceSN,
		"username":  username,
	}).Debug("Device authenticated successfully")

	return credential, nil
}

// AuthenticateByDeviceSN validates device credentials by device SN
func (a *Authenticator) AuthenticateByDeviceSN(ctx context.Context, deviceSN, password string) (*model.DeviceCredential, error) {
	credential, err := a.findCredential(ctx, "device_sn", deviceSN)
	if err != nil {
		return nil, err
	}

	if !credential.Enabled {
		a.logger.WithField("device_sn", deviceSN).Debug("Device is disabled")
		return nil, ErrDeviceDisabled
	}

	if err := a.verifyPassword(credential, password, "device_sn", deviceSN); err != nil {
		return nil, err
	}

	a.logger.WithField("device_sn", deviceSN).Debug("Device authenticated successfully")

	return credential, nil
}

// CreateCredential creates a new device credential
func (a *Authenticator) CreateCredential(ctx context.Context, deviceSN, username, password string) (*model.DeviceCredential, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	credential := &model.DeviceCredential{
		DeviceSN:     deviceSN,
		Username:     username,
		PasswordHash: string(hashedPassword),
		Enabled:      true,
	}

	result := a.db.WithContext(ctx).Create(credential)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create credential: %w", result.Error)
	}

	a.logger.WithFields(logrus.Fields{
		"device_sn": deviceSN,
		"username":  username,
	}).Info("Device credential created")

	return credential, nil
}

// UpdateCredential updates device credential
func (a *Authenticator) UpdateCredential(ctx context.Context, deviceSN string, updates map[string]any) error {
	result := a.db.WithContext(ctx).Model(&model.DeviceCredential{}).
		Where("device_sn = ?", deviceSN).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update credential: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}

	return nil
}

// setDeviceEnabled sets the enabled flag for a device.
func (a *Authenticator) setDeviceEnabled(ctx context.Context, deviceSN string, enabled bool) error {
	return a.UpdateCredential(ctx, deviceSN, map[string]any{"enabled": enabled})
}

// EnableDevice enables a device
func (a *Authenticator) EnableDevice(ctx context.Context, deviceSN string) error {
	return a.setDeviceEnabled(ctx, deviceSN, true)
}

// DisableDevice disables a device
func (a *Authenticator) DisableDevice(ctx context.Context, deviceSN string) error {
	return a.setDeviceEnabled(ctx, deviceSN, false)
}

// DeleteCredential deletes device credential
func (a *Authenticator) DeleteCredential(ctx context.Context, deviceSN string) error {
	result := a.db.WithContext(ctx).Where("device_sn = ?", deviceSN).Delete(&model.DeviceCredential{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete credential: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}

	a.logger.WithField("device_sn", deviceSN).Info("Device credential deleted")
	return nil
}

// GetCredential retrieves device credential by device SN
func (a *Authenticator) GetCredential(ctx context.Context, deviceSN string) (*model.DeviceCredential, error) {
	return a.findCredential(ctx, "device_sn", deviceSN)
}
