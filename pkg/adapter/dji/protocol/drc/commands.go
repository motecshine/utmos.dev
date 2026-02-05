package drc

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// DRC (Direct Remote Control) Commands
// ===============================

// FlightAuthorityGrabRequest represents the flight authority grab request
type FlightAuthorityGrabCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  interface{} `json:"data"`
}

// NewFlightAuthorityGrabRequest creates a new flight authority grab request
func NewFlightAuthorityGrabCommand() *FlightAuthorityGrabCommand {
	return &FlightAuthorityGrabCommand{
		Header:     common.NewHeader(),
		MethodName: "flight_authority_grab",
		DataValue:  nil,
	}
}

func (c *FlightAuthorityGrabCommand) Method() string { return c.MethodName }
func (c *FlightAuthorityGrabCommand) Data() any      { return c.DataValue }

// PayloadAuthorityGrabData represents the payload authority grab data
type PayloadAuthorityGrabData struct {
	PayloadIndex string `json:"payload_index"` // Payload index (camera enumeration value)
}

// PayloadAuthorityGrabRequest represents the payload authority grab request
type PayloadAuthorityGrabCommand struct {
	common.Header
	MethodName string                   `json:"method"`
	DataValue  PayloadAuthorityGrabData `json:"data"`
}

// NewPayloadAuthorityGrabRequest creates a new payload authority grab request
func NewPayloadAuthorityGrabCommand(data PayloadAuthorityGrabData) *PayloadAuthorityGrabCommand {
	return &PayloadAuthorityGrabCommand{
		Header:     common.NewHeader(),
		MethodName: "payload_authority_grab",
		DataValue:  data,
	}
}

func (c *PayloadAuthorityGrabCommand) Method() string { return c.MethodName }
func (c *PayloadAuthorityGrabCommand) Data() any      { return c.DataValue }

// MQTTBroker represents MQTT broker connection info
type MQTTBroker struct {
	Address    string `json:"address"`     // Server address, e.g., 192.0.2.1:8883, mqtt.dji.com:8883
	ClientID   string `json:"client_id"`   // MQTT client ID
	Username   string `json:"username"`    // Username for connection
	Password   string `json:"password"`    // Password for authentication
	ExpireTime int    `json:"expire_time"` // Authentication expiration time (seconds)
	EnableTLS  bool   `json:"enable_tls"`  // Whether to enable TLS
}

// DRCModeEnterData represents the DRC mode enter data
type DRCModeEnterData struct {
	MQTTBroker   MQTTBroker `json:"mqtt_broker"`   // MQTT broker connection info
	OSDFrequency int        `json:"osd_frequency"` // OSD frequency (1-30 Hz)
	HSIFrequency int        `json:"hsi_frequency"` // HSI frequency (1-30 Hz)
}

// DRCModeEnterRequest represents the DRC mode enter request
type DRCModeEnterCommand struct {
	common.Header
	MethodName string           `json:"method"`
	DataValue  DRCModeEnterData `json:"data"`
}

// NewDRCModeEnterRequest creates a new DRC mode enter request
func NewDRCModeEnterCommand(data DRCModeEnterData) *DRCModeEnterCommand {
	return &DRCModeEnterCommand{
		Header:     common.NewHeader(),
		MethodName: "drc_mode_enter",
		DataValue:  data,
	}
}

func (c *DRCModeEnterCommand) Method() string { return c.MethodName }
func (c *DRCModeEnterCommand) Data() any      { return c.DataValue }

// DRCModeExitRequest represents the DRC mode exit request
type DRCModeExitCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  interface{} `json:"data"`
}

// NewDRCModeExitRequest creates a new DRC mode exit request
func NewDRCModeExitCommand() *DRCModeExitCommand {
	return &DRCModeExitCommand{
		Header:     common.NewHeader(),
		MethodName: "drc_mode_exit",
		DataValue:  nil,
	}
}

func (c *DRCModeExitCommand) Method() string { return c.MethodName }
func (c *DRCModeExitCommand) Data() any      { return c.DataValue }

// TakeoffToPointData represents the takeoff to point data
type TakeoffToPointData struct {
	TargetLatitude          float64 `json:"target_latitude" binding:"required"`     // Target point latitude (-90 to 90)
	TargetLongitude         float64 `json:"target_longitude" binding:"required"`    // Target point longitude (-180 to 180)
	TargetHeight            float64 `json:"target_height" binding:"required"`       // Target point height (ellipsoidal height, WGS84)
	SecurityTakeoffHeight   float64 `json:"security_takeoff_height" example:"20"`   // Safe takeoff height (relative to dock ALT)
	RthMode                 int     `json:"rth_mode" example:"1"`                   // 【必填】返航模式设置值"0":"智能高度","1":"设定高度"
	RthAltitude             int     `json:"rth_altitude"`                           // Return home altitude (relative to dock ALT)
	RcLostAction            int     `json:"rc_lost_action" example:"2"`             // RC lost action (0=hover, 1=land, 2=RTH)
	ExitWaylineWhenRcLost   int     `json:"exit_wayline_when_rc_lost" example:"1"`  // [Deprecated] Wayline lost action (0=continue, 1=exit)
	CommanderModeLostAction int     `json:"commander_mode_lost_action" example:"1"` // 【必填】Commander mode lost action (0=continue, 1=exit)
	CommanderFlightMode     int     `json:"commander_flight_mode" example:"1"`      // 【必填】指点飞行模式设置值 "0":"智能高度飞行","1":"设定高度飞行"
	CommanderFlightHeight   float64 `json:"commander_flight_height" example:"40"`   // 【必填】Commander flight height (relative to dock ALT)
	FlightID                string  `json:"flight_id" binding:"required"`           // One-key takeoff mission UUID
	MaxSpeed                *int    `json:"max_speed,omitempty" example:"10"`       // Max speed (1-15 m/s, optional)
	//SimulateMission          *SimulateMission `json:"simulate_mission,omitempty"`             // Simulator mission settings (optional)
	FlightSafetyAdvanceCheck int `json:"flight_safety_advance_check"` // 飞行安全预检查 "0":"关闭","1":"开启"
}

// SimulateMission represents simulator mission settings
type SimulateMission struct {
	IsEnable  int     `json:"is_enable" example:"0"` // Enable simulator (0=no, 1=yes)
	Latitude  float64 `json:"latitude"`              // Simulator latitude (-90 to 90)
	Longitude float64 `json:"longitude"`             // Simulator longitude (-180 to 180)
}

// PointTarget represents a geographic target point (used by fly_to_point)
type PointTarget struct {
	Latitude  float64 `json:"latitude"`  // Latitude (-90 to 90 degrees)
	Longitude float64 `json:"longitude"` // Longitude (-180 to 180 degrees)
	Height    float64 `json:"height"`    // Height (meters, relative to takeoff point)
}

// TakeoffToPointRequest represents the takeoff to point request
type TakeoffToPointCommand struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  TakeoffToPointData `json:"data"`
}

// NewTakeoffToPointRequest creates a new takeoff to point request
func NewTakeoffToPointCommand(data TakeoffToPointData) *TakeoffToPointCommand {
	return &TakeoffToPointCommand{
		Header:     common.NewHeader(),
		MethodName: "takeoff_to_point",
		DataValue:  data,
	}
}

func (c *TakeoffToPointCommand) Method() string { return c.MethodName }
func (c *TakeoffToPointCommand) Data() any      { return c.DataValue }

// FlyToPointData represents the fly to point data
type FlyToPointData struct {
	FlyToID  string        `json:"fly_to_id"` // Fly to point ID
	MaxSpeed int           `json:"max_speed"` // Maximum speed (0-15 m/s)
	Points   []PointTarget `json:"points"`    // Target points (supports only 1 point)
}

// FlyToPointRequest represents the fly to point request
type FlyToPointCommand struct {
	common.Header
	MethodName string         `json:"method"`
	DataValue  FlyToPointData `json:"data"`
}

// NewFlyToPointRequest creates a new fly to point request
func NewFlyToPointCommand(data FlyToPointData) *FlyToPointCommand {
	return &FlyToPointCommand{
		Header:     common.NewHeader(),
		MethodName: "fly_to_point",
		DataValue:  data,
	}
}

func (c *FlyToPointCommand) Method() string { return c.MethodName }
func (c *FlyToPointCommand) Data() any      { return c.DataValue }

// FlyToPointStopRequest represents the fly to point stop request
type FlyToPointStopCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  interface{} `json:"data"`
}

// NewFlyToPointStopRequest creates a new fly to point stop request
func NewFlyToPointStopCommand() *FlyToPointStopCommand {
	return &FlyToPointStopCommand{
		Header:     common.NewHeader(),
		MethodName: "fly_to_point_stop",
		DataValue:  nil,
	}
}

func (c *FlyToPointStopCommand) Method() string { return c.MethodName }
func (c *FlyToPointStopCommand) Data() any      { return c.DataValue }

// FlyToPointUpdateData represents the fly to point update data
type FlyToPointUpdateData struct {
	MaxSpeed int           `json:"max_speed"` // Maximum speed (1-15 m/s)
	Points   []PointTarget `json:"points"`    // Updated target points (supports only 1 point)
}

// FlyToPointUpdateRequest represents the fly to point update request
type FlyToPointUpdateCommand struct {
	common.Header
	MethodName string               `json:"method"`
	DataValue  FlyToPointUpdateData `json:"data"`
}

// NewFlyToPointUpdateRequest creates a new fly to point update request
func NewFlyToPointUpdateCommand(data FlyToPointUpdateData) *FlyToPointUpdateCommand {
	return &FlyToPointUpdateCommand{
		Header:     common.NewHeader(),
		MethodName: "fly_to_point_update",
		DataValue:  data,
	}
}

func (c *FlyToPointUpdateCommand) Method() string { return c.MethodName }
func (c *FlyToPointUpdateCommand) Data() any      { return c.DataValue }

// DroneControlData represents the drone control data
type DroneControlData struct {
	Seq int     `json:"seq"` // Command sequence number
	X   float64 `json:"x"`   // Forward/backward speed (-17 to 17 m/s)
	Y   float64 `json:"y"`   // Left/right speed (-17 to 17 m/s)
	H   float64 `json:"h"`   // Up/down speed (-4 to 5 m/s)
	W   float64 `json:"w"`   // Angular velocity (-90 to 90 degrees/s)
}

// DroneControlRequest represents the drone control request
type DroneControlCommand struct {
	common.Header
	MethodName string           `json:"method"`
	DataValue  DroneControlData `json:"data"`
}

// NewDroneControlRequest creates a new drone control request
func NewDroneControlCommand(data DroneControlData) *DroneControlCommand {
	return &DroneControlCommand{
		Header:     common.NewHeader(),
		MethodName: "drone_control",
		DataValue:  data,
	}
}

func (c *DroneControlCommand) Method() string { return c.MethodName }
func (c *DroneControlCommand) Data() any      { return c.DataValue }

// StickControlData represents the stick control data
type StickControlData struct {
	Roll     int `json:"roll"`     // 横滚通道 364-1684 1024为中值（无动作），数值增大表示向右倾斜，减小表示向左倾斜
	Pitch    int `json:"pitch"`    // 俯仰通道 364-1684 1024为中值（无动作），数值增大表示向前俯冲，减小表示向后抬头。
	Throttle int `json:"throttle"` // 升降通道 364-1684 1024为悬停状态，数值增大表示升高，减小表示降低。
	Yaw      int `json:"yaw"`      // 偏航通道 364-1684 1024为中值（无动作），数值增大表示顺时针旋转，减小表示逆时针旋转。
}

// StickControlRequest represents the stick control request
type StickControlCommand struct {
	common.Header
	MethodName string           `json:"method"`
	DataValue  StickControlData `json:"data"`
}

// NewStickControlRequest creates a new stick control request
func NewStickControlCommand(data StickControlData) *StickControlCommand {
	return &StickControlCommand{
		Header:     common.NewHeader(),
		MethodName: "stick_control",
		DataValue:  data,
	}
}

func (c *StickControlCommand) Method() string { return c.MethodName }
func (c *StickControlCommand) Data() any      { return c.DataValue }

// DroneEmergencyStopRequest represents the drone emergency stop request
type DroneEmergencyStopCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  interface{} `json:"data"`
}

// NewDroneEmergencyStopRequest creates a new drone emergency stop request
func NewDroneEmergencyStopCommand() *DroneEmergencyStopCommand {
	return &DroneEmergencyStopCommand{
		Header:     common.NewHeader(),
		MethodName: "drone_emergency_stop",
		DataValue:  nil,
	}
}

func (c *DroneEmergencyStopCommand) Method() string { return c.MethodName }
func (c *DroneEmergencyStopCommand) Data() any      { return c.DataValue }

// HeartBeatData represents the heart beat data
type HeartBeatData struct {
	Seq       int   `json:"seq"`       // Command sequence number
	Timestamp int64 `json:"timestamp"` // Heart beat timestamp (milliseconds)
}

// HeartBeatRequest represents the heart beat request
type HeartBeatCommand struct {
	common.Header
	MethodName string        `json:"method"`
	DataValue  HeartBeatData `json:"data"`
}

// NewHeartBeatRequest creates a new heart beat request
func NewHeartBeatCommand(data HeartBeatData) *HeartBeatCommand {
	return &HeartBeatCommand{
		Header:     common.NewHeader(),
		MethodName: "heart_beat",
		DataValue:  data,
	}
}

func (c *HeartBeatCommand) Method() string { return c.MethodName }
func (c *HeartBeatCommand) Data() any      { return c.DataValue }

// GetHeader implements Command.GetHeader
func (c *DRCModeEnterCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *DRCModeExitCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *DroneControlCommand) GetHeader() *common.Header {
	return &c.Header
}

func (c *StickControlCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *DroneEmergencyStopCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *FlightAuthorityGrabCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *FlyToPointCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *FlyToPointStopCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *FlyToPointUpdateCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *HeartBeatCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *PayloadAuthorityGrabCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *TakeoffToPointCommand) GetHeader() *common.Header {
	return &c.Header
}
