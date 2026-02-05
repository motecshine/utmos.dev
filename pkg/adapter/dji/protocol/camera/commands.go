package camera

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Camera Control Commands
// ===============================

// CameraModeSwitchData represents the camera mode switch data
type CameraModeSwitchData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
	CameraMode   int    `json:"camera_mode"`   // Camera mode: 0=photo, 1=video
}

// CameraModeSwitchRequest represents the camera mode switch request
type CameraModeSwitchCommand struct {
	common.Header
	MethodName string               `json:"method"`
	DataValue  CameraModeSwitchData `json:"data"`
}

// NewCameraModeSwitchRequest creates a new camera mode switch request
func NewCameraModeSwitchCommand(data CameraModeSwitchData) *CameraModeSwitchCommand {
	return &CameraModeSwitchCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_mode_switch",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraModeSwitchCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraModeSwitchCommand) Data() any {
	return c.DataValue
}

// GetHeader implements Command.GetHeader
func (c *CameraModeSwitchCommand) GetHeader() *common.Header {
	return &c.Header
}

// CameraPhotoTakeData represents the camera photo take data
type CameraPhotoTakeData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
}

// CameraPhotoTakeRequest represents the camera photo take request
type CameraPhotoTakeCommand struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  CameraPhotoTakeData `json:"data"`
}

// NewCameraPhotoTakeRequest creates a new camera photo take request
func NewCameraPhotoTakeCommand(data CameraPhotoTakeData) *CameraPhotoTakeCommand {
	return &CameraPhotoTakeCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_photo_take",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraPhotoTakeCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraPhotoTakeCommand) Data() any {
	return c.DataValue
}

// CameraPhotoStopData represents the camera photo stop data
type CameraPhotoStopData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
}

// CameraPhotoStopRequest represents the camera photo stop request
type CameraPhotoStopCommand struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  CameraPhotoStopData `json:"data"`
}

// NewCameraPhotoStopRequest creates a new camera photo stop request
func NewCameraPhotoStopCommand(data CameraPhotoStopData) *CameraPhotoStopCommand {
	return &CameraPhotoStopCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_photo_stop",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraPhotoStopCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraPhotoStopCommand) Data() any {
	return c.DataValue
}

// CameraRecordingStartData represents the camera recording start data
type CameraRecordingStartData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
}

// CameraRecordingStartRequest represents the camera recording start request
type CameraRecordingStartCommand struct {
	common.Header
	MethodName string                   `json:"method"`
	DataValue  CameraRecordingStartData `json:"data"`
}

// NewCameraRecordingStartRequest creates a new camera recording start request
func NewCameraRecordingStartCommand(data CameraRecordingStartData) *CameraRecordingStartCommand {
	return &CameraRecordingStartCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_recording_start",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraRecordingStartCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraRecordingStartCommand) Data() any {
	return c.DataValue
}

// CameraRecordingStopData represents the camera recording stop data
type CameraRecordingStopData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
}

// CameraRecordingStopRequest represents the camera recording stop request
type CameraRecordingStopCommand struct {
	common.Header
	MethodName string                  `json:"method"`
	DataValue  CameraRecordingStopData `json:"data"`
}

// NewCameraRecordingStopRequest creates a new camera recording stop request
func NewCameraRecordingStopCommand(data CameraRecordingStopData) *CameraRecordingStopCommand {
	return &CameraRecordingStopCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_recording_stop",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraRecordingStopCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraRecordingStopCommand) Data() any {
	return c.DataValue
}

// CameraScreenDragData represents the camera screen drag data
type CameraScreenDragData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	Locked       bool    `json:"locked"`        // Whether aircraft head and gimbal relationship is locked
	PitchSpeed   float64 `json:"pitch_speed"`   // Gimbal pitch speed (rad/s)
	YawSpeed     float64 `json:"yaw_speed"`     // Gimbal yaw speed (rad/s, only effective when not locked)
}

// CameraScreenDragRequest represents the camera screen drag request
type CameraScreenDragCommand struct {
	common.Header
	MethodName string               `json:"method"`
	DataValue  CameraScreenDragData `json:"data"`
}

// NewCameraScreenDragCommand creates a new camera screen drag request
func NewCameraScreenDragCommand(data CameraScreenDragData) *CameraScreenDragCommand {
	return &CameraScreenDragCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_screen_drag",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraScreenDragCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraScreenDragCommand) Data() any {
	return c.DataValue
}

// CameraAimData represents the camera aim data
type CameraAimData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	CameraType   string  `json:"camera_type"`   // Camera type: ir, wide, zoom
	Locked       bool    `json:"locked"`        // Whether aircraft head and gimbal relationship is locked
	X            float64 `json:"x"`             // Target coordinate x (0-1)
	Y            float64 `json:"y"`             // Target coordinate y (0-1)
}

// CameraAimCommand represents the camera aim request
type CameraAimCommand struct {
	common.Header
	MethodName string        `json:"method"`
	DataValue  CameraAimData `json:"data"`
}

// NewCameraAimCommand creates a new camera aim request
func NewCameraAimCommand(data CameraAimData) *CameraAimCommand {
	return &CameraAimCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_aim",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraAimCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraAimCommand) Data() any {
	return c.DataValue
}

// CameraFocalLengthSetData represents the camera focal length set data
type CameraFocalLengthSetData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	CameraType   string  `json:"camera_type"`   // Camera type: ir, wide, zoom
	ZoomFactor   float64 `json:"zoom_factor"`   // Zoom factor (2-200 for visible light, 2-20 for IR)
}

// CameraFocalLengthSetCommand represents the camera focal length set request
type CameraFocalLengthSetCommand struct {
	common.Header
	MethodName string                   `json:"method"`
	DataValue  CameraFocalLengthSetData `json:"data"`
}

// NewCameraFocalLengthSetCommand creates a new camera focal length set request
func NewCameraFocalLengthSetCommand(data CameraFocalLengthSetData) *CameraFocalLengthSetCommand {
	return &CameraFocalLengthSetCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_focal_length_set",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraFocalLengthSetCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraFocalLengthSetCommand) Data() any {
	return c.DataValue
}

// CameraFrameZoomData represents the camera frame zoom data
type CameraFrameZoomData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	CameraType   string  `json:"camera_type"`   // Camera type: ir, wide, zoom
	Locked       bool    `json:"locked"`        // Whether aircraft head and gimbal relationship is locked
	X            float64 `json:"x"`             // Target coordinate x (0-1)
	Y            float64 `json:"y"`             // Target coordinate y (0-1)
	Width        float64 `json:"width"`         // Frame width (0-1)
	Height       float64 `json:"height"`        // Frame height (0-1)
}

// CameraFrameZoomRequest represents the camera frame zoom request
type CameraFrameZoomCommand struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  CameraFrameZoomData `json:"data"`
}

// CameraFrameZoomRequest creates a new camera frame zoom request
func NewCameraFrameZoomCommand(data CameraFrameZoomData) *CameraFrameZoomCommand {
	return &CameraFrameZoomCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_frame_zoom",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraFrameZoomCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraFrameZoomCommand) Data() any {
	return c.DataValue
}

// CameraLookAtData represents the camera look at data
type CameraLookAtData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	Locked       bool    `json:"locked"`        // Whether aircraft head and gimbal relationship is locked
	Latitude     float64 `json:"latitude"`      // Target point latitude (-90 to 90 degrees)
	Longitude    float64 `json:"longitude"`     // Target point longitude (-180 to 180 degrees)
	Height       float64 `json:"height"`        // Target point height (meters, relative to takeoff point)
}

// CameraLookAtCommand represents the camera look at request
type CameraLookAtCommand struct {
	common.Header
	MethodName string           `json:"method"`
	DataValue  CameraLookAtData `json:"data"`
}

// NewCameraLookAtCommand creates a new camera look at request
func NewCameraLookAtCommand(data CameraLookAtData) *CameraLookAtCommand {
	return &CameraLookAtCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_look_at",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraLookAtCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraLookAtCommand) Data() any {
	return c.DataValue
}

// CameraScreenSplitData represents the camera screen split data
type CameraScreenSplitData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
	Enable       bool   `json:"enable"`        // Whether to enable split screen mode
}

// CameraScreenSplitCommand represents the camera screen split request
type CameraScreenSplitCommand struct {
	common.Header
	MethodName string                `json:"method"`
	DataValue  CameraScreenSplitData `json:"data"`
}

// NewCameraScreenSplitCommand creates a new camera screen split request
func NewCameraScreenSplitCommand(data CameraScreenSplitData) *CameraScreenSplitCommand {
	return &CameraScreenSplitCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_screen_split",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraScreenSplitCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraScreenSplitCommand) Data() any {
	return c.DataValue
}

// CameraExposureModeSetData represents the camera exposure mode set data
type CameraExposureModeSetData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
	CameraType   string `json:"camera_type"`   // Camera type: wide, zoom
	ExposureMode int    `json:"exposure_mode"` // Exposure mode: 1=auto, 2=shutter_priority, 3=aperture_priority, 4=manual
}

// CameraExposureModeSetCommand represents the camera exposure mode set request
type CameraExposureModeSetCommand struct {
	common.Header
	MethodName string                    `json:"method"`
	DataValue  CameraExposureModeSetData `json:"data"`
}

// NewCameraExposureModeSetCommand creates a new camera exposure mode set request
func NewCameraExposureModeSetCommand(data CameraExposureModeSetData) *CameraExposureModeSetCommand {
	return &CameraExposureModeSetCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_exposure_mode_set",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraExposureModeSetCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraExposureModeSetCommand) Data() any {
	return c.DataValue
}

// CameraExposureSetData represents the camera exposure set data
type CameraExposureSetData struct {
	PayloadIndex  string `json:"payload_index"`  // Camera enumeration value
	CameraType    string `json:"camera_type"`    // Camera type: wide, zoom
	ExposureValue string `json:"exposure_value"` // Exposure value: 1=-5.0EV, 2=-4.7EV, ..., 16=0EV, ..., 31=5.0EV, 255=FIXED
}

// CameraExposureSetCommand represents the camera exposure set request
type CameraExposureSetCommand struct {
	common.Header
	MethodName string                `json:"method"`
	DataValue  CameraExposureSetData `json:"data"`
}

// NewCameraExposureSetCommand creates a new camera exposure set request
func NewCameraExposureSetCommand(data CameraExposureSetData) *CameraExposureSetCommand {
	return &CameraExposureSetCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_exposure_set",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraExposureSetCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraExposureSetCommand) Data() any {
	return c.DataValue
}

// CameraFocusModeSetData represents the camera focus mode set data
type CameraFocusModeSetData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
	CameraType   string `json:"camera_type"`   // Camera type: wide, zoom (M30 series only supports zoom)
	FocusMode    int    `json:"focus_mode"`    // Focus mode: 0=MF (manual), 1=AFS (auto single), 2=AFC (auto continuous)
}

// CameraFocusModeSetCommand represents the camera focus mode set request
type CameraFocusModeSetCommand struct {
	common.Header
	MethodName string                 `json:"method"`
	DataValue  CameraFocusModeSetData `json:"data"`
}

// NewCameraFocusModeSetCommand creates a new camera focus mode set request
func NewCameraFocusModeSetCommand(data CameraFocusModeSetData) *CameraFocusModeSetCommand {
	return &CameraFocusModeSetCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_focus_mode_set",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraFocusModeSetCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraFocusModeSetCommand) Data() any {
	return c.DataValue
}

// CameraFocusValueSetData represents the camera focus value set data
type CameraFocusValueSetData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
	CameraType   string `json:"camera_type"`   // Camera type: wide, zoom (M30 series only supports zoom)
	FocusValue   int    `json:"focus_value"`   // Focus value (range from zoom_min_focus_value to zoom_max_focus_value in OSD)
}

// CameraFocusValueSetCommand represents the camera focus value set request
type CameraFocusValueSetCommand struct {
	common.Header
	MethodName string                  `json:"method"`
	DataValue  CameraFocusValueSetData `json:"data"`
}

// NewCameraFocusValueSetCommand creates a new camera focus value set request
func NewCameraFocusValueSetCommand(data CameraFocusValueSetData) *CameraFocusValueSetCommand {
	return &CameraFocusValueSetCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_focus_value_set",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraFocusValueSetCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraFocusValueSetCommand) Data() any {
	return c.DataValue
}

// CameraPointFocusActionData represents the camera point focus action data
type CameraPointFocusActionData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	CameraType   string  `json:"camera_type"`   // Camera type: wide, zoom (M30 series only supports zoom)
	X            float64 `json:"x"`             // Focus point coordinate x (0-1)
	Y            float64 `json:"y"`             // Focus point coordinate y (0-1)
}

// CameraPointFocusActionCommand represents the camera point focus action request
type CameraPointFocusActionCommand struct {
	common.Header
	MethodName string                     `json:"method"`
	DataValue  CameraPointFocusActionData `json:"data"`
}

// NewCameraPointFocusActionCommand creates a new camera point focus action request
func NewCameraPointFocusActionCommand(data CameraPointFocusActionData) *CameraPointFocusActionCommand {
	return &CameraPointFocusActionCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_point_focus_action",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *CameraPointFocusActionCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *CameraPointFocusActionCommand) Data() any {
	return c.DataValue
}

// GimbalResetData represents the gimbal reset data
type GimbalResetData struct {
	PayloadIndex string `json:"payload_index"` // Payload index (camera enumeration value)
	ResetMode    int    `json:"reset_mode"`    // Reset mode: 0=center, 1=down, 2=yaw center, 3=pitch down
}

// GimbalResetCommand represents the gimbal reset request
type GimbalResetCommand struct {
	common.Header
	MethodName string          `json:"method"`
	DataValue  GimbalResetData `json:"data"`
}

// NewGimbalResetCommand creates a new gimbal reset request
func NewGimbalResetCommand(data GimbalResetData) *GimbalResetCommand {
	return &GimbalResetCommand{
		Header:     common.NewHeader(),
		MethodName: "gimbal_reset",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *GimbalResetCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *GimbalResetCommand) Data() any {
	return c.DataValue
}

// GetHeader implements Command.GetHeader
func (c *CameraAimCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraExposureModeSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraExposureSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraFocalLengthSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraFrameZoomCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraFocusModeSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraFocusValueSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraLookAtCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraPhotoStopCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraPhotoTakeCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraPointFocusActionCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraRecordingStartCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraRecordingStopCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraScreenDragCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CameraScreenSplitCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *GimbalResetCommand) GetHeader() *common.Header {
	return &c.Header
}
