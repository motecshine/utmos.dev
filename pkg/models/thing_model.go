package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ThingModel represents a device thing model (TSL JSON).
type ThingModel struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProductKey  string         `gorm:"uniqueIndex;size:100;not null" json:"product_key"`
	ProductName string         `gorm:"size:200;not null" json:"product_name"`
	Version     string         `gorm:"size:50;not null" json:"version"`
	TSLJSON     datatypes.JSON `gorm:"type:jsonb;not null" json:"tsl_json"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName returns the table name for the ThingModel model.
func (ThingModel) TableName() string {
	return "thing_models"
}
