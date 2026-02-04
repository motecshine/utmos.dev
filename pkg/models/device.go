// Package models provides GORM data models for UMOS IoT platform.
package models

import (
	"time"

	"gorm.io/gorm"
)

// DeviceStatus represents the status of a device.
type DeviceStatus string

const (
	DeviceStatusOnline  DeviceStatus = "online"
	DeviceStatusOffline DeviceStatus = "offline"
	DeviceStatusUnknown DeviceStatus = "unknown"
)

// Device represents an IoT device in the system.
type Device struct {
	ThingModel     *ThingModel    `gorm:"foreignKey:ThingModelID" json:"thing_model,omitempty"`
	GatewaySN      *string        `gorm:"index;size:100" json:"gateway_sn,omitempty"`
	LastOnlineTime *time.Time     `json:"last_online_time,omitempty"`
	ThingModelID   *uint          `gorm:"index" json:"thing_model_id,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeviceSN       string         `gorm:"uniqueIndex;size:100;not null" json:"device_sn"`
	DeviceName     string         `gorm:"size:200;not null" json:"device_name"`
	DeviceType     string         `gorm:"size:50;not null" json:"device_type"`
	Vendor         string         `gorm:"index;size:50;not null;default:'generic'" json:"vendor"`
	Status         DeviceStatus   `gorm:"size:20;default:'unknown'" json:"status"`
	ID             uint           `gorm:"primaryKey" json:"id"`
}

// TableName returns the table name for the Device model.
func (Device) TableName() string {
	return "devices"
}
