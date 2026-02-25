package router

import (
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/config"
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/live"
)

// Config command method names.
const (
	MethodConfig           = "config"
	MethodStorageConfigGet = "storage_config_get"
	MethodPhotoStorageSet  = "photo_storage_set"
	MethodVideoStorageSet  = "video_storage_set"
)

// Live command method names.
const (
	MethodLiveStartPush  = "live_start_push"
	MethodLiveStopPush   = "live_stop_push"
	MethodLiveSetQuality = "live_set_quality"
	MethodLiveLensChange = "live_lens_change"
)

// RegisterConfigAndLiveCommands registers all config and live streaming commands to the router.
// Returns an error if any handler registration fails.
func RegisterConfigAndLiveCommands(r *ServiceRouter) error {
	handlers := map[string]ServiceHandlerFunc{
		// Config commands
		MethodConfig:           SimpleCommandHandler[config.Data](MethodConfig),
		MethodStorageConfigGet: SimpleCommandHandler[config.StorageConfigGetData](MethodStorageConfigGet),
		MethodPhotoStorageSet:  SimpleCommandHandler[config.PhotoStorageSetData](MethodPhotoStorageSet),
		MethodVideoStorageSet:  SimpleCommandHandler[config.VideoStorageSetData](MethodVideoStorageSet),
		// Live streaming commands
		MethodLiveStartPush:  SimpleCommandHandler[live.StartPushData](MethodLiveStartPush),
		MethodLiveStopPush:   SimpleCommandHandler[live.StopPushData](MethodLiveStopPush),
		MethodLiveSetQuality: SimpleCommandHandler[live.SetQualityData](MethodLiveSetQuality),
		MethodLiveLensChange: SimpleCommandHandler[live.LensChangeData](MethodLiveLensChange),
	}

	return RegisterHandlers(r, handlers)
}
