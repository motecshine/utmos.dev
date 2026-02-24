package router

import (
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/wayline"
)

// Wayline command method names.
const (
	MethodFlighttaskCreate   = "flighttask_create"
	MethodFlighttaskPrepare  = "flighttask_prepare"
	MethodFlighttaskExecute  = "flighttask_execute"
	MethodFlighttaskPause    = "flighttask_pause"
	MethodFlighttaskRecovery = "flighttask_recovery"
	MethodFlighttaskUndo     = "flighttask_undo"
	MethodReturnHome         = "return_home"
	MethodReturnHomeCancel   = "return_home_cancel"
)

// RegisterWaylineCommands registers all wayline commands to the router.
// Returns an error if any handler registration fails.
//
// Registration files share structural pattern but register different command types
func RegisterWaylineCommands(r *ServiceRouter) error {
	handlers := map[string]ServiceHandlerFunc{
		// Flight task commands with data - use types from protocol/wayline package
		MethodFlighttaskCreate:  SimpleCommandHandler[wayline.CreateData](MethodFlighttaskCreate),
		MethodFlighttaskPrepare: SimpleCommandHandler[wayline.PrepareData](MethodFlighttaskPrepare),
		MethodFlighttaskExecute: SimpleCommandHandler[wayline.ExecuteData](MethodFlighttaskExecute),
		MethodFlighttaskUndo:    SimpleCommandHandler[wayline.UndoData](MethodFlighttaskUndo),

		// Simple commands without data payload
		MethodFlighttaskPause:    NoDataCommandHandler(MethodFlighttaskPause),
		MethodFlighttaskRecovery: NoDataCommandHandler(MethodFlighttaskRecovery),
		MethodReturnHome:         NoDataCommandHandler(MethodReturnHome),
		MethodReturnHomeCancel:   NoDataCommandHandler(MethodReturnHomeCancel),
	}

	return RegisterHandlers(r, handlers)
}
