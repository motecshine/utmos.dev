package live

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Live Streaming Commands
// ===============================

// LiveStartPushData represents the start live push data
type LiveStartPushData struct {
	URLType      int    `json:"url_type"`      // URL type: 0=Agora, 1=RTMP, 3=GB28181, 4=WebRTC
	URL          string `json:"url"`           // Live stream URL/parameters
	VideoID      string `json:"video_id"`      // Video stream ID (format: {sn}/{camera_index}/{video_index})
	VideoQuality int    `json:"video_quality"` // Video quality: 0=adaptive, 1=smooth, 2=standard, 3=high, 4=super
}

// LiveStartPushCommand represents the start live push request
type LiveStartPushCommand struct {
	common.Header
	MethodName string            `json:"method"`
	DataValue  LiveStartPushData `json:"data"`
}

// NewLiveStartPushCommand creates a new start live push command
func NewLiveStartPushCommand(data LiveStartPushData) *LiveStartPushCommand {
	return &LiveStartPushCommand{
		Header:     common.NewHeader(),
		MethodName: "live_start_push",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *LiveStartPushCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *LiveStartPushCommand) Data() any {
	return c.DataValue
}

// LiveStopPushData represents the stop live push data
type LiveStopPushData struct {
	VideoID string `json:"video_id"` // Video stream ID
}

// LiveStopPushCommand represents the stop live push request
type LiveStopPushCommand struct {
	common.Header
	MethodName string           `json:"method"`
	DataValue  LiveStopPushData `json:"data"`
}

// NewLiveStopPushCommand creates a new stop live push command
func NewLiveStopPushCommand(data LiveStopPushData) *LiveStopPushCommand {
	return &LiveStopPushCommand{
		Header:     common.NewHeader(),
		MethodName: "live_stop_push",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *LiveStopPushCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *LiveStopPushCommand) Data() any {
	return c.DataValue
}

// LiveSetQualityData represents the set live quality data
type LiveSetQualityData struct {
	VideoID      string `json:"video_id"`      // Video stream ID (format: {sn}/{camera_index}/{video_index})
	VideoQuality int    `json:"video_quality"` // Video quality: 0=adaptive, 1=smooth, 2=standard, 3=high, 4=super
}

// LiveSetQualityCommand represents the set live quality request
type LiveSetQualityCommand struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  LiveSetQualityData `json:"data"`
}

// NewLiveSetQualityCommand creates a new set live quality command
func NewLiveSetQualityCommand(data LiveSetQualityData) *LiveSetQualityCommand {
	return &LiveSetQualityCommand{
		Header:     common.NewHeader(),
		MethodName: "live_set_quality",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *LiveSetQualityCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *LiveSetQualityCommand) Data() any {
	return c.DataValue
}

// LiveLensChangeData represents the live lens change data
type LiveLensChangeData struct {
	// VideoID   string `json:"video_id"`   // Video stream ID (format: {sn}/{camera_index}/{video_index})
	VideoType string `json:"video_type"` // Video type: ir=infrared, normal=default, wide=wide angle, zoom=zoom
}

// LiveLensChangeCommand represents the live lens change request
type LiveLensChangeCommand struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  LiveLensChangeData `json:"data"`
}

// NewLiveLensChangeCommand creates a new live lens change command
func NewLiveLensChangeCommand(data LiveLensChangeData) *LiveLensChangeCommand {
	return &LiveLensChangeCommand{
		Header:     common.NewHeader(),
		MethodName: "live_lens_change",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *LiveLensChangeCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *LiveLensChangeCommand) Data() any {
	return c.DataValue
}

// GetHeader implements Command.GetHeader
func (c *LiveLensChangeCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *LiveSetQualityCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *LiveStartPushCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *LiveStopPushCommand) GetHeader() *common.Header {
	return &c.Header
}
