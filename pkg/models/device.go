package models

import (
	"time"

	"gorm.io/gorm"
)

// DeviceStatus represents device status enum
type DeviceStatus string

const (
	// DeviceStatusOnline indicates device is online
	DeviceStatusOnline DeviceStatus = "online"
	// DeviceStatusOffline indicates device is offline
	DeviceStatusOffline DeviceStatus = "offline"
	// DeviceStatusUnknown indicates device status is unknown
	DeviceStatusUnknown DeviceStatus = "unknown"
)

// Device represents a device entity
type Device struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	DeviceSN      string         `gorm:"uniqueIndex;type:varchar(100);not null" json:"device_sn"`
	DeviceName    string         `gorm:"type:varchar(200);not null" json:"device_name"`
	DeviceType    string         `gorm:"type:varchar(50);not null" json:"device_type"`
	GatewaySN     *string        `gorm:"index;type:varchar(100)" json:"gateway_sn,omitempty"`
	ThingModelID  uint           `gorm:"index;not null" json:"thing_model_id"`
	Vendor        string         `gorm:"type:varchar(50);not null" json:"vendor"`
	Status        DeviceStatus  `gorm:"type:varchar(20);default:'unknown'" json:"status"`
	LastOnlineTime *time.Time    `json:"last_online_time,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName returns the table name for Device
func (Device) TableName() string {
	return "devices"
}

