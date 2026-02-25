package router

// RegisterDRCEvents registers DRC-related events to the event router.
// Returns an error if any handler registration fails.
func RegisterDRCEvents(r *EventRouter) error {
	handlers := map[string]EventHandlerFunc{
		MethodJoystickInvalidNotify: NoDataEventHandler(MethodJoystickInvalidNotify),
		MethodDRCStatusNotify:       NoDataEventHandler(MethodDRCStatusNotify),
	}

	return RegisterEventHandlers(r, handlers)
}
