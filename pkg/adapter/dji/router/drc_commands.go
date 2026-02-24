package router

import (
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/drc"
)

// DRC command method names.
const (
	MethodDRCModeEnter       = "drc_mode_enter"
	MethodDRCModeExit        = "drc_mode_exit"
	MethodDroneControl       = "drone_control"
	MethodDroneEmergencyStop = "drone_emergency_stop"
	MethodHeart              = "heart"
)

// RegisterDRCCommands registers all DRC commands to the router.
// Returns an error if any handler registration fails.
//
// Registration files share structural pattern but register different command types
func RegisterDRCCommands(r *ServiceRouter) error {
	handlers := map[string]ServiceHandlerFunc{
		// DRC mode commands - use types from protocol/drc package
		MethodDRCModeEnter: SimpleCommandHandler[drc.DRCModeEnterData](MethodDRCModeEnter),
		MethodDRCModeExit:  NoDataCommandHandler(MethodDRCModeExit),

		// Control commands
		MethodDroneControl:       SimpleCommandHandler[drc.DroneControlData](MethodDroneControl),
		MethodDroneEmergencyStop: NoDataCommandHandler(MethodDroneEmergencyStop),

		// Heartbeat
		MethodHeart: SimpleCommandHandler[drc.HeartBeatData](MethodHeart),
	}

	return RegisterHandlers(r, handlers)
}
