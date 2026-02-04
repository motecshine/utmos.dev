package models

import (
	"time"

	"gorm.io/datatypes"
)

// MessageDirection represents the direction of a message.
type MessageDirection string

const (
	MessageDirectionUplink   MessageDirection = "uplink"
	MessageDirectionDownlink MessageDirection = "downlink"
)

// MessageStatus represents the status of a message.
type MessageStatus string

const (
	MessageStatusSuccess MessageStatus = "success"
	MessageStatusFailed  MessageStatus = "failed"
	MessageStatusPending MessageStatus = "pending"
)

// MessageLog represents a message log record for debugging and tracing.
type MessageLog struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	TID          string           `gorm:"index;size:100" json:"tid"`
	BID          string           `gorm:"index;size:100" json:"bid"`
	Service      string           `gorm:"size:50;not null" json:"service"`
	Direction    MessageDirection `gorm:"size:20;not null" json:"direction"`
	MessageType  string           `gorm:"size:100;not null" json:"message_type"`
	DeviceSN     string           `gorm:"index:idx_message_log_device_created;size:100" json:"device_sn"`
	MessageData  datatypes.JSON   `gorm:"type:jsonb;not null" json:"message_data"`
	Status       MessageStatus    `gorm:"size:20;default:'pending'" json:"status"`
	ErrorMessage *string          `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time        `gorm:"index:idx_message_log_device_created" json:"created_at"`
}

// TableName returns the table name for the MessageLog model.
func (MessageLog) TableName() string {
	return "message_logs"
}
