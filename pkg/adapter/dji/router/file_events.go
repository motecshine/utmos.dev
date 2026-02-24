package router

import (
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/file"
)

// RegisterFileEvents registers file-related events to the event router.
// Returns an error if any handler registration fails.
//
// structurally similar to command registration functions but operates on different handler types
func RegisterFileEvents(r *EventRouter) error {
	handlers := map[string]EventHandlerFunc{
		// Use types from protocol/file package
		MethodFileUploadCallback:    SimpleEventHandler[file.FileUploadCallbackData](MethodFileUploadCallback),
		MethodFileUploadProgress:    NoDataEventHandler(MethodFileUploadProgress),
		MethodHighestPriorityUpload: SimpleEventHandler[file.HighestPriorityUploadFlighttaskMediaData](MethodHighestPriorityUpload),
	}

	return RegisterEventHandlers(r, handlers)
}
