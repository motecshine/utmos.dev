package router

import (
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/file"
)

// File command method names.
const (
	MethodFileUploadStart  = "fileupload_start"
	MethodFileUploadFinish = "fileupload_update"
	MethodFileUploadList   = "fileupload_list"
)

// RegisterFileCommands registers all file commands to the router.
// Returns an error if any handler registration fails.
func RegisterFileCommands(r *ServiceRouter) error {
	handlers := map[string]ServiceHandlerFunc{
		// File upload commands - use types from protocol/file package
		MethodFileUploadStart:  SimpleCommandHandler[file.FileUploadStartData](MethodFileUploadStart),
		MethodFileUploadFinish: SimpleCommandHandler[file.FileUploadUpdateData](MethodFileUploadFinish),
		MethodFileUploadList:   SimpleCommandHandler[file.FileUploadListData](MethodFileUploadList),
	}

	return RegisterHandlers(r, handlers)
}
