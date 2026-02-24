package device

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// CoverOpenCommand represents the open cover command
type CoverOpenCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewCoverOpenCommand creates a new cover open command
func NewCoverOpenCommand() *CoverOpenCommand {
	return &CoverOpenCommand{
		Header:     common.NewHeader(),
		MethodName: "cover_open",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *CoverOpenCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CoverOpenCommand) Data() any {
	return c.DataValue
}

// CoverCloseCommand represents the close cover command
type CoverCloseCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewCoverCloseCommand creates a new cover close command
func NewCoverCloseCommand() *CoverCloseCommand {
	return &CoverCloseCommand{
		Header:     common.NewHeader(),
		MethodName: "cover_close",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *CoverCloseCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CoverCloseCommand) Data() any {
	return c.DataValue
}

// CoverForceCloseCommand represents the force close cover command
type CoverForceCloseCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewCoverForceCloseCommand creates a new force close cover command
func NewCoverForceCloseCommand() *CoverForceCloseCommand {
	return &CoverForceCloseCommand{
		Header:     common.NewHeader(),
		MethodName: "cover_force_close",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *CoverForceCloseCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CoverForceCloseCommand) Data() any {
	return c.DataValue
}

// DroneOpenCommand represents the drone power on command
type DroneOpenCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewDroneOpenCommand creates a new drone power on command
func NewDroneOpenCommand() *DroneOpenCommand {
	return &DroneOpenCommand{
		Header:     common.NewHeader(),
		MethodName: "drone_open",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *DroneOpenCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *DroneOpenCommand) Data() any {
	return c.DataValue
}

// DroneCloseCommand represents the drone power off command
type DroneCloseCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewDroneCloseCommand creates a new drone power off command
func NewDroneCloseCommand() *DroneCloseCommand {
	return &DroneCloseCommand{
		Header:     common.NewHeader(),
		MethodName: "drone_close",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *DroneCloseCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *DroneCloseCommand) Data() any {
	return c.DataValue
}

// ChargeOpenCommand represents the start charging command
type ChargeOpenCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewChargeOpenCommand creates a new start charging command
func NewChargeOpenCommand() *ChargeOpenCommand {
	return &ChargeOpenCommand{
		Header:     common.NewHeader(),
		MethodName: "charge_open",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *ChargeOpenCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *ChargeOpenCommand) Data() any {
	return c.DataValue
}

// ChargeCloseCommand represents the stop charging command
type ChargeCloseCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewChargeCloseCommand creates a new stop charging command
func NewChargeCloseCommand() *ChargeCloseCommand {
	return &ChargeCloseCommand{
		Header:     common.NewHeader(),
		MethodName: "charge_close",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *ChargeCloseCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *ChargeCloseCommand) Data() any {
	return c.DataValue
}

// DeviceRebootCommand represents the device reboot command
type DeviceRebootCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewDeviceRebootCommand creates a new device reboot command
func NewDeviceRebootCommand() *DeviceRebootCommand {
	return &DeviceRebootCommand{
		Header:     common.NewHeader(),
		MethodName: "device_reboot",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *DeviceRebootCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *DeviceRebootCommand) Data() any {
	return c.DataValue
}

// DeviceFormatCommand represents the device format command
type DeviceFormatCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewDeviceFormatCommand creates a new device format command
func NewDeviceFormatCommand() *DeviceFormatCommand {
	return &DeviceFormatCommand{
		Header:     common.NewHeader(),
		MethodName: "device_format",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *DeviceFormatCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *DeviceFormatCommand) Data() any {
	return c.DataValue
}

// DroneFormatCommand represents the drone format command
type DroneFormatCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewDroneFormatCommand creates a new drone format command
func NewDroneFormatCommand() *DroneFormatCommand {
	return &DroneFormatCommand{
		Header:     common.NewHeader(),
		MethodName: "drone_format",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *DroneFormatCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *DroneFormatCommand) Data() any {
	return c.DataValue
}

// PutterOpenCommand represents the pusher open command (dock1 only)
type PutterOpenCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewPutterOpenCommand creates a new pusher open command
func NewPutterOpenCommand() *PutterOpenCommand {
	return &PutterOpenCommand{
		Header:     common.NewHeader(),
		MethodName: "putter_open",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *PutterOpenCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *PutterOpenCommand) Data() any {
	return c.DataValue
}

// PutterCloseCommand represents the pusher close command (dock1 only)
type PutterCloseCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewPutterCloseCommand creates a new pusher close command
func NewPutterCloseCommand() *PutterCloseCommand {
	return &PutterCloseCommand{
		Header:     common.NewHeader(),
		MethodName: "putter_close",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *PutterCloseCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *PutterCloseCommand) Data() any {
	return c.DataValue
}

// DebugModeOpenCommand represents the debug mode enable command
type DebugModeOpenCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewDebugModeOpenCommand creates a new debug mode enable command
func NewDebugModeOpenCommand() *DebugModeOpenCommand {
	return &DebugModeOpenCommand{
		Header:     common.NewHeader(),
		MethodName: "debug_mode_open",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *DebugModeOpenCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *DebugModeOpenCommand) Data() any {
	return c.DataValue
}

// DebugModeCloseCommand represents the debug mode disable command
type DebugModeCloseCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewDebugModeCloseCommand creates a new debug mode disable command
func NewDebugModeCloseCommand() *DebugModeCloseCommand {
	return &DebugModeCloseCommand{
		Header:     common.NewHeader(),
		MethodName: "debug_mode_close",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *DebugModeCloseCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *DebugModeCloseCommand) Data() any {
	return c.DataValue
}

// SupplementLightOpenCommand represents the supplement light enable command
type SupplementLightOpenCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewSupplementLightOpenCommand creates a new supplement light enable command
func NewSupplementLightOpenCommand() *SupplementLightOpenCommand {
	return &SupplementLightOpenCommand{
		Header:     common.NewHeader(),
		MethodName: "supplement_light_open",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *SupplementLightOpenCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *SupplementLightOpenCommand) Data() any {
	return c.DataValue
}

// SupplementLightCloseCommand represents the supplement light disable command
type SupplementLightCloseCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewSupplementLightCloseCommand creates a new supplement light disable command
func NewSupplementLightCloseCommand() *SupplementLightCloseCommand {
	return &SupplementLightCloseCommand{
		Header:     common.NewHeader(),
		MethodName: "supplement_light_close",
		DataValue:  nil,
	}
}

// Method implements Command.Method
func (c *SupplementLightCloseCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *SupplementLightCloseCommand) Data() any {
	return c.DataValue
}

// ===============================
// Complex Device Control Commands (data_type: "object")
// ===============================

// Battery maintenance mode constants
const (
	BatteryMaintenanceModeOff = 0 // Turn off battery maintenance mode
	BatteryMaintenanceModeOn  = 1 // Turn on battery maintenance mode
)

// BatteryMaintenanceSwitchData represents the battery maintenance switch data
type BatteryMaintenanceSwitchData struct {
	Action int `json:"action"` // Battery maintenance mode: 0=off, 1=on
}

// BatteryMaintenanceSwitchCommand represents the battery maintenance mode switch command
type BatteryMaintenanceSwitchCommand struct {
	common.Header
	MethodName string                       `json:"method"`
	DataValue  BatteryMaintenanceSwitchData `json:"data"`
}

// NewBatteryMaintenanceSwitchCommand creates a new battery maintenance mode switch command
func NewBatteryMaintenanceSwitchCommand(data BatteryMaintenanceSwitchData) *BatteryMaintenanceSwitchCommand {
	return &BatteryMaintenanceSwitchCommand{
		Header:     common.NewHeader(),
		MethodName: "battery_maintenance_switch",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *BatteryMaintenanceSwitchCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *BatteryMaintenanceSwitchCommand) Data() any {
	return c.DataValue
}

// AirConditioner mode constants
const (
	AirConditionerModeIdle       = 0 // Idle mode (turn off cooling, heating, or dehumidification)
	AirConditionerModeCooling    = 1 // Cooling mode
	AirConditionerModeHeating    = 2 // Heating mode
	AirConditionerModeDehumidify = 3 // Dehumidification mode (auto cooling or heating dehumidification)
)

// AirConditionerModeSwitchData represents the air conditioner mode switch data
type AirConditionerModeSwitchData struct {
	Action int `json:"action"` // Air conditioner mode: 0=idle, 1=cooling, 2=heating, 3=dehumidify
}

// AirConditionerModeSwitchCommand represents the air conditioner mode switch command
type AirConditionerModeSwitchCommand struct {
	common.Header
	MethodName string                       `json:"method"`
	DataValue  AirConditionerModeSwitchData `json:"data"`
}

// NewAirConditionerModeSwitchCommand creates a new air conditioner mode switch command
func NewAirConditionerModeSwitchCommand(data AirConditionerModeSwitchData) *AirConditionerModeSwitchCommand {
	return &AirConditionerModeSwitchCommand{
		Header:     common.NewHeader(),
		MethodName: "air_conditioner_mode_switch",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *AirConditionerModeSwitchCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *AirConditionerModeSwitchCommand) Data() any {
	return c.DataValue
}

// AlarmStateSwitchData represents the alarm state switch data
type AlarmStateSwitchData struct {
	Action int `json:"action"` // Alarm state: 0=off, 1=on
}

// AlarmStateSwitchCommand represents the alarm state switch command
type AlarmStateSwitchCommand struct {
	common.Header
	MethodName string               `json:"method"`
	DataValue  AlarmStateSwitchData `json:"data"`
}

// NewAlarmStateSwitchCommand creates a new alarm state switch command
func NewAlarmStateSwitchCommand(data AlarmStateSwitchData) *AlarmStateSwitchCommand {
	return &AlarmStateSwitchCommand{
		Header:     common.NewHeader(),
		MethodName: "alarm_state_switch",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *AlarmStateSwitchCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *AlarmStateSwitchCommand) Data() any {
	return c.DataValue
}

// Battery store mode constants
const (
	BatteryStoreModePlan    = 1 // Plan mode
	BatteryStoreModeStandby = 2 // Standby mode
)

// BatteryStoreModeSwitchData represents the battery store mode switch data
type BatteryStoreModeSwitchData struct {
	Action int `json:"action"` // Battery store mode: 1=plan, 2=standby
}

// BatteryStoreModeSwitch represents the battery storage mode switch command
type BatteryStoreModeSwitch struct {
	common.Header
	MethodName string                     `json:"method"`
	DataValue  BatteryStoreModeSwitchData `json:"data"`
}

// NewBatteryStoreModeSwitch creates a new battery storage mode switch command
func NewBatteryStoreModeSwitch(data BatteryStoreModeSwitchData) *BatteryStoreModeSwitch {
	return &BatteryStoreModeSwitch{
		Header:     common.NewHeader(),
		MethodName: "battery_store_mode_switch",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *BatteryStoreModeSwitch) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *BatteryStoreModeSwitch) Data() any {
	return c.DataValue
}

// SDR work mode constants
const (
	SDRWorkmodeSdrOnly    = 0 // SDR only mode
	SDRWorkmode4GEnhanced = 1 // 4G enhanced mode (SDR + 4G)
)

// SDRWorkmodeSwitchData represents the SDR work mode switch data
type SDRWorkmodeSwitchData struct {
	LinkWorkmode int `json:"link_workmode"` // Link work mode: 0=SDR only, 1=4G enhanced
}

// SDRWorkmodeSwitch represents the SDR work mode switch command
type SDRWorkmodeSwitch struct {
	common.Header
	MethodName string                `json:"method"`
	DataValue  SDRWorkmodeSwitchData `json:"data"`
}

// NewSDRWorkmodeSwitch creates a new SDR work mode switch command
func NewSDRWorkmodeSwitch(data SDRWorkmodeSwitchData) *SDRWorkmodeSwitch {
	return &SDRWorkmodeSwitch{
		Header:     common.NewHeader(),
		MethodName: "sdr_workmode_switch",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *SDRWorkmodeSwitch) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *SDRWorkmodeSwitch) Data() any {
	return c.DataValue
}

// GetHeader implements Command.GetHeader
func (c *AirConditionerModeSwitchCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *AlarmStateSwitchCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *BatteryMaintenanceSwitchCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *ChargeCloseCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *ChargeOpenCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CoverCloseCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CoverForceCloseCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CoverOpenCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *DebugModeCloseCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *DebugModeOpenCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *DeviceFormatCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *DeviceRebootCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *DroneCloseCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *DroneFormatCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *DroneOpenCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *PutterCloseCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *PutterOpenCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *SupplementLightCloseCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *SupplementLightOpenCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *BatteryStoreModeSwitch) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *SDRWorkmodeSwitch) GetHeader() *common.Header {
	return &c.Header
}
