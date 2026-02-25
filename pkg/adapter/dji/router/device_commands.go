package router

import (
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/device"
)

// Device command method names.
const (
	MethodCoverOpen                = "cover_open"
	MethodCoverClose               = "cover_close"
	MethodCoverForceClose          = "cover_force_close"
	MethodDroneOpen                = "drone_open"
	MethodDroneClose               = "drone_close"
	MethodChargeOpen               = "charge_open"
	MethodChargeClose              = "charge_close"
	MethodDeviceReboot             = "device_reboot"
	MethodDeviceFormat             = "device_format"
	MethodDroneFormat              = "drone_format"
	MethodDebugModeOpen            = "debug_mode_open"
	MethodDebugModeClose           = "debug_mode_close"
	MethodBatteryMaintenanceSwitch = "battery_maintenance_switch"
	MethodAirConditionerModeSwitch = "air_conditioner_mode_switch"
	MethodAlarmStateSwitch         = "alarm_state_switch"
	MethodSDRWorkmodeSwitch        = "sdr_workmode_switch"
)

// RegisterDeviceCommands registers all device control commands to the router.
// Returns an error if any handler registration fails.
func RegisterDeviceCommands(r *ServiceRouter) error {
	handlers := map[string]ServiceHandlerFunc{
		// Simple commands without data payload
		MethodCoverOpen:       NoDataCommandHandler(MethodCoverOpen),
		MethodCoverClose:      NoDataCommandHandler(MethodCoverClose),
		MethodCoverForceClose: NoDataCommandHandler(MethodCoverForceClose),
		MethodDroneOpen:       NoDataCommandHandler(MethodDroneOpen),
		MethodDroneClose:      NoDataCommandHandler(MethodDroneClose),
		MethodChargeOpen:      NoDataCommandHandler(MethodChargeOpen),
		MethodChargeClose:     NoDataCommandHandler(MethodChargeClose),
		MethodDeviceReboot:    NoDataCommandHandler(MethodDeviceReboot),
		MethodDeviceFormat:    NoDataCommandHandler(MethodDeviceFormat),
		MethodDroneFormat:     NoDataCommandHandler(MethodDroneFormat),
		MethodDebugModeOpen:   NoDataCommandHandler(MethodDebugModeOpen),
		MethodDebugModeClose:  NoDataCommandHandler(MethodDebugModeClose),

		// Commands with data payload - use types from protocol/device package
		MethodBatteryMaintenanceSwitch: SimpleCommandHandler[device.BatteryMaintenanceSwitchData](MethodBatteryMaintenanceSwitch),
		MethodAirConditionerModeSwitch: SimpleCommandHandler[device.AirConditionerModeSwitchData](MethodAirConditionerModeSwitch),
		MethodAlarmStateSwitch:         SimpleCommandHandler[device.AlarmStateSwitchData](MethodAlarmStateSwitch),
		MethodSDRWorkmodeSwitch:        SimpleCommandHandler[device.SDRWorkmodeSwitchData](MethodSDRWorkmodeSwitch),
	}

	return RegisterHandlers(r, handlers)
}
