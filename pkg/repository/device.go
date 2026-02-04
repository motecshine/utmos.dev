// Package repository provides data access layer implementations.
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	pkgerrors "github.com/utmos/utmos/pkg/errors"
	"github.com/utmos/utmos/pkg/models"
)

// DeviceRepository provides device data access.
type DeviceRepository struct {
	db *gorm.DB
}

// NewDeviceRepository creates a new DeviceRepository.
func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// GetByDeviceSN retrieves a device by its serial number.
func (r *DeviceRepository) GetByDeviceSN(ctx context.Context, deviceSN string) (*models.Device, error) {
	var device models.Device
	result := r.db.WithContext(ctx).Where("device_sn = ?", deviceSN).First(&device)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, pkgerrors.New(pkgerrors.ErrDeviceNotFound, "device not found")
		}
		return nil, pkgerrors.Wrap(result.Error, pkgerrors.ErrDatabaseConnection, "failed to query device")
	}
	return &device, nil
}

// GetVendorByDeviceSN retrieves the vendor for a device by its serial number.
func (r *DeviceRepository) GetVendorByDeviceSN(ctx context.Context, deviceSN string) (string, error) {
	var device models.Device
	result := r.db.WithContext(ctx).Select("vendor").Where("device_sn = ?", deviceSN).First(&device)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", pkgerrors.New(pkgerrors.ErrDeviceNotFound, "device not found")
		}
		return "", pkgerrors.Wrap(result.Error, pkgerrors.ErrDatabaseConnection, "failed to query device vendor")
	}
	return device.Vendor, nil
}

// Create creates a new device.
func (r *DeviceRepository) Create(ctx context.Context, device *models.Device) error {
	result := r.db.WithContext(ctx).Create(device)
	if result.Error != nil {
		return pkgerrors.Wrap(result.Error, pkgerrors.ErrDatabaseConnection, "failed to create device")
	}
	return nil
}

// Update updates an existing device.
func (r *DeviceRepository) Update(ctx context.Context, device *models.Device) error {
	result := r.db.WithContext(ctx).Save(device)
	if result.Error != nil {
		return pkgerrors.Wrap(result.Error, pkgerrors.ErrDatabaseConnection, "failed to update device")
	}
	return nil
}

// UpdateStatus updates the status of a device.
func (r *DeviceRepository) UpdateStatus(ctx context.Context, deviceSN string, status models.DeviceStatus) error {
	result := r.db.WithContext(ctx).
		Model(&models.Device{}).
		Where("device_sn = ?", deviceSN).
		Update("status", status)
	if result.Error != nil {
		return pkgerrors.Wrap(result.Error, pkgerrors.ErrDatabaseConnection, "failed to update device status")
	}
	if result.RowsAffected == 0 {
		return pkgerrors.New(pkgerrors.ErrDeviceNotFound, "device not found")
	}
	return nil
}

// List retrieves devices with pagination.
func (r *DeviceRepository) List(ctx context.Context, offset, limit int) ([]models.Device, error) {
	var devices []models.Device
	result := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&devices)
	if result.Error != nil {
		return nil, pkgerrors.Wrap(result.Error, pkgerrors.ErrDatabaseConnection, "failed to list devices")
	}
	return devices, nil
}

// ListByVendor retrieves devices by vendor.
func (r *DeviceRepository) ListByVendor(ctx context.Context, vendor string, offset, limit int) ([]models.Device, error) {
	var devices []models.Device
	result := r.db.WithContext(ctx).Where("vendor = ?", vendor).Offset(offset).Limit(limit).Find(&devices)
	if result.Error != nil {
		return nil, pkgerrors.Wrap(result.Error, pkgerrors.ErrDatabaseConnection, "failed to list devices by vendor")
	}
	return devices, nil
}

// Delete soft-deletes a device.
func (r *DeviceRepository) Delete(ctx context.Context, deviceSN string) error {
	result := r.db.WithContext(ctx).Where("device_sn = ?", deviceSN).Delete(&models.Device{})
	if result.Error != nil {
		return pkgerrors.Wrap(result.Error, pkgerrors.ErrDatabaseConnection, "failed to delete device")
	}
	if result.RowsAffected == 0 {
		return pkgerrors.New(pkgerrors.ErrDeviceNotFound, "device not found")
	}
	return nil
}
