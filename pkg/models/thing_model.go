package models

import (
	"time"

	"gorm.io/gorm"
)

// ThingModel represents a thing model definition
type ThingModel struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProductKey  string         `gorm:"uniqueIndex;type:varchar(100);not null" json:"product_key"`
	ProductName string         `gorm:"type:varchar(200);not null" json:"product_name"`
	Version     string         `gorm:"type:varchar(50);not null" json:"version"`
	TSLJSON     string         `gorm:"type:jsonb;not null" json:"tsl_json"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName returns the table name for ThingModel
func (ThingModel) TableName() string {
	return "thing_models"
}

