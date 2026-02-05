package router

// RegisterWaylineEvents registers wayline-related events to the event router.
// Returns an error if any handler registration fails.
func RegisterWaylineEvents(r *EventRouter) error {
	handlers := map[string]EventHandlerFunc{
		MethodFlighttaskProgress: SimpleEventHandler[FlighttaskProgressData](MethodFlighttaskProgress),
		MethodFlighttaskReady:    SimpleEventHandler[FlighttaskReadyData](MethodFlighttaskReady),
		MethodReturnHomeInfo:     SimpleEventHandler[ReturnHomeInfoData](MethodReturnHomeInfo),
	}

	return RegisterEventHandlers(r, handlers)
}
