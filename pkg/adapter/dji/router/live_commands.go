package router

import (
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/live"
)

// Live command method names.
const (
	MethodLiveStartPush  = "live_start_push"
	MethodLiveStopPush   = "live_stop_push"
	MethodLiveSetQuality = "live_set_quality"
	MethodLiveLensChange = "live_lens_change"
)

// RegisterLiveCommands registers all live commands to the router.
// Returns an error if any handler registration fails.
func RegisterLiveCommands(r *ServiceRouter) error {
	handlers := map[string]ServiceHandlerFunc{
		MethodLiveStartPush:  SimpleCommandHandler[live.LiveStartPushData](MethodLiveStartPush),
		MethodLiveStopPush:   SimpleCommandHandler[live.LiveStopPushData](MethodLiveStopPush),
		MethodLiveSetQuality: SimpleCommandHandler[live.LiveSetQualityData](MethodLiveSetQuality),
		MethodLiveLensChange: SimpleCommandHandler[live.LiveLensChangeData](MethodLiveLensChange),
	}

	return RegisterHandlers(r, handlers)
}
