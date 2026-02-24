package wayline

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Wayline Mission Commands (Cloud → Device)
// ===============================

// CreateCommand represents the create flight task command
type CreateCommand struct {
	common.Header
	MethodName string     `json:"method"`
	DataValue  CreateData `json:"data"`
}

// NewCreateCommand creates a new flight task creation command
func NewCreateCommand(data CreateData) *CreateCommand {
	return &CreateCommand{
		Header:     common.NewHeader(),
		MethodName: "flighttask_create",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *CreateCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *CreateCommand) Data() any { return c.DataValue }

// GetHeader returns the event header.
func (c *CreateCommand) GetHeader() *common.Header { return &c.Header }

// PrepareCommand represents the prepare flight task command
type PrepareCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  PrepareData `json:"data"`
}

// NewPrepareCommand creates a new flight task preparation command
func NewPrepareCommand(data PrepareData) *PrepareCommand {
	return &PrepareCommand{
		Header:     common.NewHeader(),
		MethodName: "flighttask_prepare",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *PrepareCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *PrepareCommand) Data() any { return c.DataValue }

// GetHeader returns the event header.
func (c *PrepareCommand) GetHeader() *common.Header { return &c.Header }

// ExecuteCommand represents the execute flight task command
type ExecuteCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  ExecuteData `json:"data"`
}

// NewExecuteCommand creates a new flight task execution command
func NewExecuteCommand(data ExecuteData) *ExecuteCommand {
	return &ExecuteCommand{
		Header:     common.NewHeader(),
		MethodName: "flighttask_execute",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *ExecuteCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *ExecuteCommand) Data() any { return c.DataValue }

// GetHeader returns the event header.
func (c *ExecuteCommand) GetHeader() *common.Header { return &c.Header }

// PauseCommand represents the pause flight task command
type PauseCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewPauseCommand creates a new flight task pause command
func NewPauseCommand() *PauseCommand {
	return &PauseCommand{
		Header:     common.NewHeader(),
		MethodName: "flighttask_pause",
		DataValue:  nil,
	}
}

// Method returns the method name.
func (c *PauseCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *PauseCommand) Data() any { return c.DataValue }

// GetHeader returns the event header.
func (c *PauseCommand) GetHeader() *common.Header { return &c.Header }

// RecoveryCommand represents the recovery flight task command
type RecoveryCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewRecoveryCommand creates a new flight task recovery command
func NewRecoveryCommand() *RecoveryCommand {
	return &RecoveryCommand{
		Header:     common.NewHeader(),
		MethodName: "flighttask_recovery",
		DataValue:  nil,
	}
}

// Method returns the method name.
func (c *RecoveryCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *RecoveryCommand) Data() any { return c.DataValue }

// GetHeader returns the event header.
func (c *RecoveryCommand) GetHeader() *common.Header { return &c.Header }

// UndoCommand represents the undo flight task command
type UndoCommand struct {
	common.Header
	MethodName string   `json:"method"`
	DataValue  UndoData `json:"data"`
}

// NewUndoCommand creates a new flight task undo command
func NewUndoCommand(data UndoData) *UndoCommand {
	return &UndoCommand{
		Header:     common.NewHeader(),
		MethodName: "flighttask_undo",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *UndoCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *UndoCommand) Data() any { return c.DataValue }

// GetHeader returns the event header.
func (c *UndoCommand) GetHeader() *common.Header { return &c.Header }

// ReturnHomeCommand represents the return to home command
type ReturnHomeCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewReturnHomeCommand creates a new return to home command
func NewReturnHomeCommand() *ReturnHomeCommand {
	return &ReturnHomeCommand{
		Header:     common.NewHeader(),
		MethodName: "return_home",
		DataValue:  nil,
	}
}

// Method returns the method name.
func (c *ReturnHomeCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *ReturnHomeCommand) Data() any { return c.DataValue }

// GetHeader returns the event header.
func (c *ReturnHomeCommand) GetHeader() *common.Header { return &c.Header }

// CancelReturnHomeCommand represents the cancel return to home command
type CancelReturnHomeCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewCancelReturnHomeCommand creates a new cancel return to home command
func NewCancelReturnHomeCommand() *CancelReturnHomeCommand {
	return &CancelReturnHomeCommand{
		Header:     common.NewHeader(),
		MethodName: "return_home_cancel",
		DataValue:  nil,
	}
}

// Method returns the method name.
func (c *CancelReturnHomeCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *CancelReturnHomeCommand) Data() any { return c.DataValue }

// GetHeader returns the event header.
func (c *CancelReturnHomeCommand) GetHeader() *common.Header { return &c.Header }

// AbortFlightSetupCommand represents the abort flight setup command (dock1 only)
type AbortFlightSetupCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewAbortFlightSetupCommand creates a new abort flight setup command
func NewAbortFlightSetupCommand() *AbortFlightSetupCommand {
	return &AbortFlightSetupCommand{
		Header:     common.NewHeader(),
		MethodName: "flight_setup_abort",
		DataValue:  nil,
	}
}

// Method returns the method name.
func (c *AbortFlightSetupCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *AbortFlightSetupCommand) Data() any { return c.DataValue }

// GetHeader returns the event header.
func (c *AbortFlightSetupCommand) GetHeader() *common.Header { return &c.Header }
