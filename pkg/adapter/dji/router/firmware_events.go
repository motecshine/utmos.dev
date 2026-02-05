package router

// RegisterFirmwareEvents registers firmware-related events to the event router.
// Returns an error if any handler registration fails.
func RegisterFirmwareEvents(r *EventRouter) error {
	handlers := map[string]EventHandlerFunc{
		MethodOTAProgress: NoDataEventHandler(MethodOTAProgress),
	}

	return RegisterEventHandlers(r, handlers)
}
