package camera

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// IR Camera Commands
// ===============================

// IRMeteringModeSetData represents the IR metering mode set data
type IRMeteringModeSetData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
	Mode         int    `json:"mode"`          // Metering mode: 0=off, 1=point, 2=area
}

// IRMeteringModeSetCommand represents the IR metering mode set request
type IRMeteringModeSetCommand struct {
	common.Header
	MethodName string                `json:"method"`
	DataValue  IRMeteringModeSetData `json:"data"`
}

// NewIRMeteringModeSetCommand creates a new IR metering mode set request
func NewIRMeteringModeSetCommand(data IRMeteringModeSetData) *IRMeteringModeSetCommand {
	return &IRMeteringModeSetCommand{
		Header:     common.NewHeader(),
		MethodName: "ir_metering_mode_set",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *IRMeteringModeSetCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *IRMeteringModeSetCommand) Data() any { return c.DataValue }

// IRMeteringPointSetData represents the IR metering point set data
type IRMeteringPointSetData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	X            float64 `json:"x"`             // Metering point coordinate x (0-1)
	Y            float64 `json:"y"`             // Metering point coordinate y (0-1)
}

// IRMeteringPointSetCommand represents the IR metering point set request
type IRMeteringPointSetCommand struct {
	common.Header
	MethodName string                 `json:"method"`
	DataValue  IRMeteringPointSetData `json:"data"`
}

// NewIRMeteringPointSetCommand creates a new IR metering point set request
func NewIRMeteringPointSetCommand(data IRMeteringPointSetData) *IRMeteringPointSetCommand {
	return &IRMeteringPointSetCommand{
		Header:     common.NewHeader(),
		MethodName: "ir_metering_point_set",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *IRMeteringPointSetCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *IRMeteringPointSetCommand) Data() any { return c.DataValue }

// IRMeteringAreaSetData represents the IR metering area set data
type IRMeteringAreaSetData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	X            float64 `json:"x"`             // Metering area upper-left corner x coordinate (0-1)
	Y            float64 `json:"y"`             // Metering area upper-left corner y coordinate (0-1)
	Width        float64 `json:"width"`         // Metering area width (0-1)
	Height       float64 `json:"height"`        // Metering area height (0-1)
}

// IRMeteringAreaSetCommand represents the IR metering area set request
type IRMeteringAreaSetCommand struct {
	common.Header
	MethodName string                `json:"method"`
	DataValue  IRMeteringAreaSetData `json:"data"`
}

// NewIRMeteringAreaSetCommand creates a new IR metering area set request
func NewIRMeteringAreaSetCommand(data IRMeteringAreaSetData) *IRMeteringAreaSetCommand {
	return &IRMeteringAreaSetCommand{
		Header:     common.NewHeader(),
		MethodName: "ir_metering_area_set",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *IRMeteringAreaSetCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *IRMeteringAreaSetCommand) Data() any { return c.DataValue }

// GetHeader implements Command.GetHeader
func (c *IRMeteringAreaSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *IRMeteringModeSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *IRMeteringPointSetCommand) GetHeader() *common.Header {
	return &c.Header
}
