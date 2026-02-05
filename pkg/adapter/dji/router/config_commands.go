package router

import (
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/config"
)

// Config command method names.
const (
	MethodConfig           = "config"
	MethodStorageConfigGet = "storage_config_get"
	MethodPhotoStorageSet  = "photo_storage_set"
	MethodVideoStorageSet  = "video_storage_set"
)

// RegisterConfigCommands registers all config commands to the router.
// Returns an error if any handler registration fails.
func RegisterConfigCommands(r *ServiceRouter) error {
	handlers := map[string]ServiceHandlerFunc{
		// Config request - use types from protocol/config package
		MethodConfig:           SimpleCommandHandler[config.ConfigData](MethodConfig),
		MethodStorageConfigGet: SimpleCommandHandler[config.StorageConfigGetData](MethodStorageConfigGet),
		MethodPhotoStorageSet:  SimpleCommandHandler[config.PhotoStorageSetData](MethodPhotoStorageSet),
		MethodVideoStorageSet:  SimpleCommandHandler[config.VideoStorageSetData](MethodVideoStorageSet),
	}

	return RegisterHandlers(r, handlers)
}
