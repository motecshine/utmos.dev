package router

import (
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/firmware"
)

// Firmware command method names.
const (
	MethodOTACreate = "ota_create"
)

// RegisterFirmwareCommands registers all firmware commands to the router.
// Returns an error if any handler registration fails.
func RegisterFirmwareCommands(r *ServiceRouter) error {
	handlers := map[string]ServiceHandlerFunc{
		MethodOTACreate: SimpleCommandHandler[firmware.OTACreateData](MethodOTACreate),
	}

	return RegisterHandlers(r, handlers)
}
