package models

import (
	"github.com/utmos/utmos/internal/downlink/model"
	"gorm.io/gorm"
)

// AutoMigrate runs GORM auto migration for all models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&ThingModel{},
		&Device{},
		&DeviceProperty{},
		&DeviceEvent{},
		&MessageLog{},
		&model.ServiceCall{},
	)
}
