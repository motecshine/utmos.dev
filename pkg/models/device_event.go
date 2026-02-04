package models

import (
	"time"

	"gorm.io/datatypes"
)

// DeviceEvent represents a device event record.
type DeviceEvent struct {
	Device    *Device        `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	EventData datatypes.JSON `gorm:"type:jsonb;not null" json:"event_data"`
	Timestamp time.Time      `gorm:"index:idx_device_event_timestamp;not null" json:"timestamp"`
	CreatedAt time.Time      `json:"created_at"`
	EventKey  string         `gorm:"index;size:100;not null" json:"event_key"`
	DeviceID  uint           `gorm:"index:idx_device_event_timestamp;not null" json:"device_id"`
	ID        uint           `gorm:"primaryKey" json:"id"`
}

// TableName returns the table name for the DeviceEvent model.
func (DeviceEvent) TableName() string {
	return "device_events"
}
