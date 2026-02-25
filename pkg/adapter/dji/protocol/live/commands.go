package live

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Live Streaming Commands
// ===============================

// StartPushData represents the start live push data
type StartPushData struct {
	URLType      int    `json:"url_type"`      // URL type: 0=Agora, 1=RTMP, 3=GB28181, 4=WebRTC
	URL          string `json:"url"`           // Live stream URL/parameters
	VideoID      string `json:"video_id"`      // Video stream ID (format: {sn}/{camera_index}/{video_index})
	VideoQuality int    `json:"video_quality"` // Video quality: 0=adaptive, 1=smooth, 2=standard, 3=high, 4=super
}

// StartPushCommand represents the start live push request
type StartPushCommand struct {
	common.Header
	MethodName string            `json:"method"`
	DataValue  StartPushData `json:"data"`
}

// NewStartPushCommand creates a new start live push command
func NewStartPushCommand(data StartPushData) *StartPushCommand {
	return &StartPushCommand{
		Header:     common.NewHeader(),
		MethodName: "live_start_push",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *StartPushCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *StartPushCommand) Data() any {
	return c.DataValue
}

// StopPushData represents the stop live push data
type StopPushData struct {
	VideoID string `json:"video_id"` // Video stream ID
}

// StopPushCommand represents the stop live push request
type StopPushCommand struct {
	common.Header
	MethodName string           `json:"method"`
	DataValue  StopPushData `json:"data"`
}

// NewStopPushCommand creates a new stop live push command
func NewStopPushCommand(data StopPushData) *StopPushCommand {
	return &StopPushCommand{
		Header:     common.NewHeader(),
		MethodName: "live_stop_push",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *StopPushCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *StopPushCommand) Data() any {
	return c.DataValue
}

// SetQualityData represents the set live quality data
type SetQualityData struct {
	VideoID      string `json:"video_id"`      // Video stream ID (format: {sn}/{camera_index}/{video_index})
	VideoQuality int    `json:"video_quality"` // Video quality: 0=adaptive, 1=smooth, 2=standard, 3=high, 4=super
}

// SetQualityCommand represents the set live quality request
type SetQualityCommand struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  SetQualityData `json:"data"`
}

// NewSetQualityCommand creates a new set live quality command
func NewSetQualityCommand(data SetQualityData) *SetQualityCommand {
	return &SetQualityCommand{
		Header:     common.NewHeader(),
		MethodName: "live_set_quality",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *SetQualityCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *SetQualityCommand) Data() any {
	return c.DataValue
}

// LensChangeData represents the live lens change data
type LensChangeData struct {
	// VideoID   string `json:"video_id"`   // Video stream ID (format: {sn}/{camera_index}/{video_index})
	VideoType string `json:"video_type"` // Video type: ir=infrared, normal=default, wide=wide angle, zoom=zoom
}

// LensChangeCommand represents the live lens change request
type LensChangeCommand struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  LensChangeData `json:"data"`
}

// NewLensChangeCommand creates a new live lens change command
func NewLensChangeCommand(data LensChangeData) *LensChangeCommand {
	return &LensChangeCommand{
		Header:     common.NewHeader(),
		MethodName: "live_lens_change",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *LensChangeCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *LensChangeCommand) Data() any {
	return c.DataValue
}

// GetHeader implements Command.GetHeader
func (c *LensChangeCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *SetQualityCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *StartPushCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *StopPushCommand) GetHeader() *common.Header {
	return &c.Header
}
