package models

import (
	"time"

	"gorm.io/datatypes"
)

// DeviceProperty represents a device property value.
type DeviceProperty struct {
	Device        *Device        `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	PropertyValue datatypes.JSON `gorm:"type:jsonb;not null" json:"property_value"`
	UpdatedAt     time.Time      `json:"updated_at"`
	PropertyKey   string         `gorm:"uniqueIndex:idx_device_property;size:100;not null" json:"property_key"`
	DeviceID      uint           `gorm:"uniqueIndex:idx_device_property;not null" json:"device_id"`
	ID            uint           `gorm:"primaryKey" json:"id"`
}

// TableName returns the table name for the DeviceProperty model.
func (DeviceProperty) TableName() string {
	return "device_properties"
}
