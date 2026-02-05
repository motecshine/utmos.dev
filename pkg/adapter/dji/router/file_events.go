package router

// RegisterFileEvents registers file-related events to the event router.
// Returns an error if any handler registration fails.
func RegisterFileEvents(r *EventRouter) error {
	handlers := map[string]EventHandlerFunc{
		MethodFileUploadCallback:    SimpleEventHandler[FileUploadCallbackData](MethodFileUploadCallback),
		MethodFileUploadProgress:    NoDataEventHandler(MethodFileUploadProgress),
		MethodHighestPriorityUpload: NoDataEventHandler(MethodHighestPriorityUpload),
	}

	return RegisterEventHandlers(r, handlers)
}
