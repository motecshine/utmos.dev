package router

import (
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/wayline"
)

// RegisterWaylineEvents registers wayline-related events to the event router.
// Returns an error if any handler registration fails.
func RegisterWaylineEvents(r *EventRouter) error {
	handlers := map[string]EventHandlerFunc{
		// Use types from protocol/wayline package
		MethodFlighttaskProgress: SimpleEventHandler[wayline.ProgressData](MethodFlighttaskProgress),
		MethodFlighttaskReady:    SimpleEventHandler[wayline.ReadyData](MethodFlighttaskReady),
		MethodReturnHomeInfo:     SimpleEventHandler[wayline.ReturnHomeInfoData](MethodReturnHomeInfo),
	}

	return RegisterEventHandlers(r, handlers)
}
