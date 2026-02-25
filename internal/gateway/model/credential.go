// Package model contains data models for iot-gateway
package model

import (
	"time"
)

// DeviceCredential represents device authentication credentials
type DeviceCredential struct {
	ID           uint      `gorm:"primaryKey"`
	DeviceSN     string    `gorm:"uniqueIndex;size:64;not null"`
	Username     string    `gorm:"size:64;not null"`
	PasswordHash string    `gorm:"size:256;not null"`
	Enabled      bool      `gorm:"default:true"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for DeviceCredential
func (DeviceCredential) TableName() string {
	return "device_credentials"
}
