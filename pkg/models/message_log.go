package models

import (
	"time"

	"gorm.io/datatypes"
)

// MessageDirection represents the direction of a message.
type MessageDirection string

const (
	// MessageDirectionUplink is the uplink message direction.
	MessageDirectionUplink MessageDirection = "uplink"
	// MessageDirectionDownlink is the downlink message direction.
	MessageDirectionDownlink MessageDirection = "downlink"
)

// MessageStatus represents the status of a message.
type MessageStatus string

const (
	// MessageStatusSuccess indicates the message was processed successfully.
	MessageStatusSuccess MessageStatus = "success"
	// MessageStatusFailed indicates the message processing failed.
	MessageStatusFailed MessageStatus = "failed"
	// MessageStatusPending indicates the message is pending processing.
	MessageStatusPending MessageStatus = "pending"
)

// MessageLog represents a message log record for debugging and tracing.
type MessageLog struct {
	ErrorMessage *string          `gorm:"type:text" json:"error_message,omitempty"`
	MessageData  datatypes.JSON   `gorm:"type:jsonb;not null" json:"message_data"`
	CreatedAt    time.Time        `gorm:"index:idx_message_log_device_created" json:"created_at"`
	TID          string           `gorm:"index;size:100" json:"tid"`
	BID          string           `gorm:"index;size:100" json:"bid"`
	Service      string           `gorm:"size:50;not null" json:"service"`
	MessageType  string           `gorm:"size:100;not null" json:"message_type"`
	DeviceSN     string           `gorm:"index:idx_message_log_device_created;size:100" json:"device_sn"`
	Direction    MessageDirection `gorm:"size:20;not null" json:"direction"`
	Status       MessageStatus    `gorm:"size:20;default:'pending'" json:"status"`
	ID           uint             `gorm:"primaryKey" json:"id"`
}

// TableName returns the table name for the MessageLog model.
func (MessageLog) TableName() string {
	return "message_logs"
}
