package router

import (
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/file"
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/firmware"
)

// File command method names.
const (
	MethodFileUploadStart  = "fileupload_start"
	MethodFileUploadFinish = "fileupload_update"
	MethodFileUploadList   = "fileupload_list"
)

// Firmware command method names.
const (
	MethodOTACreate = "ota_create"
)

// RegisterFileAndFirmwareCommands registers all file upload and firmware commands to the router.
// Returns an error if any handler registration fails.
func RegisterFileAndFirmwareCommands(r *ServiceRouter) error {
	handlers := map[string]ServiceHandlerFunc{
		// File upload commands
		MethodFileUploadStart:  SimpleCommandHandler[file.FileUploadStartData](MethodFileUploadStart),
		MethodFileUploadFinish: SimpleCommandHandler[file.FileUploadUpdateData](MethodFileUploadFinish),
		MethodFileUploadList:   SimpleCommandHandler[file.FileUploadListData](MethodFileUploadList),
		// Firmware commands
		MethodOTACreate: SimpleCommandHandler[firmware.OTACreateData](MethodOTACreate),
	}

	return RegisterHandlers(r, handlers)
}
