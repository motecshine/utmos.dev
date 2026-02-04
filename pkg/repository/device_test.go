package repository

import (
	"context"
	"testing"

	"github.com/utmos/utmos/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Auto migrate the Device model
	if err := db.AutoMigrate(&models.Device{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestDeviceRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDeviceRepository(db)
	ctx := context.Background()

	device := &models.Device{
		DeviceSN:     "test-device-001",
		DeviceName:   "Test Device",
		DeviceType:   "sensor",
		Vendor:       "dji",
		Status:       models.DeviceStatusOnline,
		ThingModelID: nil,
	}

	err := repo.Create(ctx, device)
	if err != nil {
		t.Fatalf("failed to create device: %v", err)
	}

	if device.ID == 0 {
		t.Error("expected device ID to be set")
	}
}

func TestDeviceRepository_GetByDeviceSN(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDeviceRepository(db)
	ctx := context.Background()

	// Create a device first
	device := &models.Device{
		DeviceSN:   "test-device-002",
		DeviceName: "Test Device 2",
		DeviceType: "gateway",
		Vendor:     "generic",
		Status:     models.DeviceStatusOnline,
	}
	if err := repo.Create(ctx, device); err != nil {
		t.Fatalf("failed to create device: %v", err)
	}

	// Get by device SN
	found, err := repo.GetByDeviceSN(ctx, "test-device-002")
	if err != nil {
		t.Fatalf("failed to get device: %v", err)
	}

	if found.DeviceSN != "test-device-002" {
		t.Errorf("expected device SN test-device-002, got %s", found.DeviceSN)
	}

	if found.Vendor != "generic" {
		t.Errorf("expected vendor generic, got %s", found.Vendor)
	}
}

func TestDeviceRepository_GetByDeviceSN_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDeviceRepository(db)
	ctx := context.Background()

	_, err := repo.GetByDeviceSN(ctx, "non-existent-device")
	if err == nil {
		t.Error("expected error for non-existent device")
	}
}

func TestDeviceRepository_GetVendorByDeviceSN(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDeviceRepository(db)
	ctx := context.Background()

	// Create devices with different vendors
	devices := []*models.Device{
		{DeviceSN: "dji-001", DeviceName: "DJI Drone", DeviceType: "drone", Vendor: "dji", Status: models.DeviceStatusOnline},
		{DeviceSN: "tuya-001", DeviceName: "Tuya Sensor", DeviceType: "sensor", Vendor: "tuya", Status: models.DeviceStatusOnline},
		{DeviceSN: "generic-001", DeviceName: "Generic Device", DeviceType: "device", Vendor: "generic", Status: models.DeviceStatusOnline},
	}

	for _, d := range devices {
		if err := repo.Create(ctx, d); err != nil {
			t.Fatalf("failed to create device: %v", err)
		}
	}

	tests := []struct {
		deviceSN       string
		expectedVendor string
	}{
		{"dji-001", "dji"},
		{"tuya-001", "tuya"},
		{"generic-001", "generic"},
	}

	for _, tt := range tests {
		t.Run(tt.deviceSN, func(t *testing.T) {
			vendor, err := repo.GetVendorByDeviceSN(ctx, tt.deviceSN)
			if err != nil {
				t.Fatalf("failed to get vendor: %v", err)
			}
			if vendor != tt.expectedVendor {
				t.Errorf("expected vendor %s, got %s", tt.expectedVendor, vendor)
			}
		})
	}
}

func TestDeviceRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDeviceRepository(db)
	ctx := context.Background()

	// Create a device
	device := &models.Device{
		DeviceSN:   "update-test-001",
		DeviceName: "Original Name",
		DeviceType: "sensor",
		Vendor:     "dji",
		Status:     models.DeviceStatusOffline,
	}
	if err := repo.Create(ctx, device); err != nil {
		t.Fatalf("failed to create device: %v", err)
	}

	// Update the device
	device.DeviceName = "Updated Name"
	device.Status = models.DeviceStatusOnline
	if err := repo.Update(ctx, device); err != nil {
		t.Fatalf("failed to update device: %v", err)
	}

	// Verify the update
	found, err := repo.GetByDeviceSN(ctx, "update-test-001")
	if err != nil {
		t.Fatalf("failed to get device: %v", err)
	}

	if found.DeviceName != "Updated Name" {
		t.Errorf("expected name Updated Name, got %s", found.DeviceName)
	}
	if found.Status != models.DeviceStatusOnline {
		t.Errorf("expected status online, got %s", found.Status)
	}
}

func TestDeviceRepository_UpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDeviceRepository(db)
	ctx := context.Background()

	// Create a device
	device := &models.Device{
		DeviceSN:   "status-test-001",
		DeviceName: "Status Test",
		DeviceType: "sensor",
		Vendor:     "generic",
		Status:     models.DeviceStatusOffline,
	}
	if err := repo.Create(ctx, device); err != nil {
		t.Fatalf("failed to create device: %v", err)
	}

	// Update status
	if err := repo.UpdateStatus(ctx, "status-test-001", models.DeviceStatusOnline); err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	// Verify
	found, err := repo.GetByDeviceSN(ctx, "status-test-001")
	if err != nil {
		t.Fatalf("failed to get device: %v", err)
	}

	if found.Status != models.DeviceStatusOnline {
		t.Errorf("expected status online, got %s", found.Status)
	}
}

func TestDeviceRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDeviceRepository(db)
	ctx := context.Background()

	// Create a device
	device := &models.Device{
		DeviceSN:   "delete-test-001",
		DeviceName: "Delete Test",
		DeviceType: "sensor",
		Vendor:     "dji",
		Status:     models.DeviceStatusOnline,
	}
	if err := repo.Create(ctx, device); err != nil {
		t.Fatalf("failed to create device: %v", err)
	}

	// Delete the device
	if err := repo.Delete(ctx, "delete-test-001"); err != nil {
		t.Fatalf("failed to delete device: %v", err)
	}

	// Verify deletion
	_, err := repo.GetByDeviceSN(ctx, "delete-test-001")
	if err == nil {
		t.Error("expected error for deleted device")
	}
}

func TestDeviceRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDeviceRepository(db)
	ctx := context.Background()

	// Create multiple devices
	for i := 0; i < 5; i++ {
		device := &models.Device{
			DeviceSN:   "list-test-" + string(rune('a'+i)),
			DeviceName: "List Test Device",
			DeviceType: "sensor",
			Vendor:     "generic",
			Status:     models.DeviceStatusOnline,
		}
		if err := repo.Create(ctx, device); err != nil {
			t.Fatalf("failed to create device: %v", err)
		}
	}

	// List with pagination
	devices, err := repo.List(ctx, 0, 3)
	if err != nil {
		t.Fatalf("failed to list devices: %v", err)
	}

	if len(devices) != 3 {
		t.Errorf("expected 3 devices, got %d", len(devices))
	}
}

func TestDeviceRepository_ListByVendor(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDeviceRepository(db)
	ctx := context.Background()

	// Create devices with different vendors
	devices := []*models.Device{
		{DeviceSN: "vendor-dji-1", DeviceName: "DJI 1", DeviceType: "drone", Vendor: "dji", Status: models.DeviceStatusOnline},
		{DeviceSN: "vendor-dji-2", DeviceName: "DJI 2", DeviceType: "drone", Vendor: "dji", Status: models.DeviceStatusOnline},
		{DeviceSN: "vendor-tuya-1", DeviceName: "Tuya 1", DeviceType: "sensor", Vendor: "tuya", Status: models.DeviceStatusOnline},
	}

	for _, d := range devices {
		if err := repo.Create(ctx, d); err != nil {
			t.Fatalf("failed to create device: %v", err)
		}
	}

	// List by vendor
	djiDevices, err := repo.ListByVendor(ctx, "dji", 0, 100)
	if err != nil {
		t.Fatalf("failed to list by vendor: %v", err)
	}

	if len(djiDevices) != 2 {
		t.Errorf("expected 2 DJI devices, got %d", len(djiDevices))
	}
}

