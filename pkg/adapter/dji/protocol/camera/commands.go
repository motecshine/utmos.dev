package camera

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Camera Control Commands
// ===============================

// ModeSwitchData represents the camera mode switch data
type ModeSwitchData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
	CameraMode   int    `json:"camera_mode"`   // Camera mode: 0=photo, 1=video
}

// ModeSwitchCommand represents the camera mode switch request
type ModeSwitchCommand struct {
	common.Header
	MethodName string               `json:"method"`
	DataValue  ModeSwitchData `json:"data"`
}

// NewModeSwitchCommand creates a new camera mode switch request
func NewModeSwitchCommand(data ModeSwitchData) *ModeSwitchCommand {
	return &ModeSwitchCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_mode_switch",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *ModeSwitchCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *ModeSwitchCommand) Data() any {
	return c.DataValue
}

// GetHeader implements Command.GetHeader
func (c *ModeSwitchCommand) GetHeader() *common.Header {
	return &c.Header
}

// PhotoTakeData represents the camera photo take data
type PhotoTakeData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
}

// PhotoTakeCommand represents the camera photo take request
type PhotoTakeCommand struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  PhotoTakeData `json:"data"`
}

// NewPhotoTakeCommand creates a new camera photo take request
func NewPhotoTakeCommand(data PhotoTakeData) *PhotoTakeCommand {
	return &PhotoTakeCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_photo_take",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *PhotoTakeCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *PhotoTakeCommand) Data() any {
	return c.DataValue
}

// PhotoStopData represents the camera photo stop data
type PhotoStopData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
}

// PhotoStopCommand represents the camera photo stop request
type PhotoStopCommand struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  PhotoStopData `json:"data"`
}

// NewPhotoStopCommand creates a new camera photo stop request
func NewPhotoStopCommand(data PhotoStopData) *PhotoStopCommand {
	return &PhotoStopCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_photo_stop",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *PhotoStopCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *PhotoStopCommand) Data() any {
	return c.DataValue
}

// RecordingStartData represents the camera recording start data
type RecordingStartData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
}

// RecordingStartCommand represents the camera recording start request
type RecordingStartCommand struct {
	common.Header
	MethodName string                   `json:"method"`
	DataValue  RecordingStartData `json:"data"`
}

// NewRecordingStartCommand creates a new camera recording start request
func NewRecordingStartCommand(data RecordingStartData) *RecordingStartCommand {
	return &RecordingStartCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_recording_start",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *RecordingStartCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *RecordingStartCommand) Data() any {
	return c.DataValue
}

// RecordingStopData represents the camera recording stop data
type RecordingStopData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
}

// RecordingStopCommand represents the camera recording stop request
type RecordingStopCommand struct {
	common.Header
	MethodName string                  `json:"method"`
	DataValue  RecordingStopData `json:"data"`
}

// NewRecordingStopCommand creates a new camera recording stop request
func NewRecordingStopCommand(data RecordingStopData) *RecordingStopCommand {
	return &RecordingStopCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_recording_stop",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *RecordingStopCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *RecordingStopCommand) Data() any {
	return c.DataValue
}

// ScreenDragData represents the camera screen drag data
type ScreenDragData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	Locked       bool    `json:"locked"`        // Whether aircraft head and gimbal relationship is locked
	PitchSpeed   float64 `json:"pitch_speed"`   // Gimbal pitch speed (rad/s)
	YawSpeed     float64 `json:"yaw_speed"`     // Gimbal yaw speed (rad/s, only effective when not locked)
}

// ScreenDragCommand represents the camera screen drag request
type ScreenDragCommand struct {
	common.Header
	MethodName string               `json:"method"`
	DataValue  ScreenDragData `json:"data"`
}

// NewScreenDragCommand creates a new camera screen drag request
func NewScreenDragCommand(data ScreenDragData) *ScreenDragCommand {
	return &ScreenDragCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_screen_drag",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *ScreenDragCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *ScreenDragCommand) Data() any {
	return c.DataValue
}

// AimData represents the camera aim data
type AimData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	CameraType   string  `json:"camera_type"`   // Camera type: ir, wide, zoom
	Locked       bool    `json:"locked"`        // Whether aircraft head and gimbal relationship is locked
	X            float64 `json:"x"`             // Target coordinate x (0-1)
	Y            float64 `json:"y"`             // Target coordinate y (0-1)
}

// AimCommand represents the camera aim request
type AimCommand struct {
	common.Header
	MethodName string        `json:"method"`
	DataValue  AimData `json:"data"`
}

// NewAimCommand creates a new camera aim request
func NewAimCommand(data AimData) *AimCommand {
	return &AimCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_aim",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *AimCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *AimCommand) Data() any {
	return c.DataValue
}

// FocalLengthSetData represents the camera focal length set data
type FocalLengthSetData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	CameraType   string  `json:"camera_type"`   // Camera type: ir, wide, zoom
	ZoomFactor   float64 `json:"zoom_factor"`   // Zoom factor (2-200 for visible light, 2-20 for IR)
}

// FocalLengthSetCommand represents the camera focal length set request
type FocalLengthSetCommand struct {
	common.Header
	MethodName string                   `json:"method"`
	DataValue  FocalLengthSetData `json:"data"`
}

// NewFocalLengthSetCommand creates a new camera focal length set request
func NewFocalLengthSetCommand(data FocalLengthSetData) *FocalLengthSetCommand {
	return &FocalLengthSetCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_focal_length_set",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *FocalLengthSetCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *FocalLengthSetCommand) Data() any {
	return c.DataValue
}

// FrameZoomData represents the camera frame zoom data
type FrameZoomData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	CameraType   string  `json:"camera_type"`   // Camera type: ir, wide, zoom
	Locked       bool    `json:"locked"`        // Whether aircraft head and gimbal relationship is locked
	X            float64 `json:"x"`             // Target coordinate x (0-1)
	Y            float64 `json:"y"`             // Target coordinate y (0-1)
	Width        float64 `json:"width"`         // Frame width (0-1)
	Height       float64 `json:"height"`        // Frame height (0-1)
}

// FrameZoomCommand represents the camera frame zoom request
type FrameZoomCommand struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  FrameZoomData `json:"data"`
}

// NewFrameZoomCommand creates a new camera frame zoom request
func NewFrameZoomCommand(data FrameZoomData) *FrameZoomCommand {
	return &FrameZoomCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_frame_zoom",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *FrameZoomCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *FrameZoomCommand) Data() any {
	return c.DataValue
}

// LookAtData represents the camera look at data
type LookAtData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	Locked       bool    `json:"locked"`        // Whether aircraft head and gimbal relationship is locked
	Latitude     float64 `json:"latitude"`      // Target point latitude (-90 to 90 degrees)
	Longitude    float64 `json:"longitude"`     // Target point longitude (-180 to 180 degrees)
	Height       float64 `json:"height"`        // Target point height (meters, relative to takeoff point)
}

// LookAtCommand represents the camera look at request
type LookAtCommand struct {
	common.Header
	MethodName string           `json:"method"`
	DataValue  LookAtData `json:"data"`
}

// NewLookAtCommand creates a new camera look at request
func NewLookAtCommand(data LookAtData) *LookAtCommand {
	return &LookAtCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_look_at",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *LookAtCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *LookAtCommand) Data() any {
	return c.DataValue
}

// ScreenSplitData represents the camera screen split data
type ScreenSplitData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
	Enable       bool   `json:"enable"`        // Whether to enable split screen mode
}

// ScreenSplitCommand represents the camera screen split request
type ScreenSplitCommand struct {
	common.Header
	MethodName string                `json:"method"`
	DataValue  ScreenSplitData `json:"data"`
}

// NewScreenSplitCommand creates a new camera screen split request
func NewScreenSplitCommand(data ScreenSplitData) *ScreenSplitCommand {
	return &ScreenSplitCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_screen_split",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *ScreenSplitCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *ScreenSplitCommand) Data() any {
	return c.DataValue
}

// ExposureModeSetData represents the camera exposure mode set data
type ExposureModeSetData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
	CameraType   string `json:"camera_type"`   // Camera type: wide, zoom
	ExposureMode int    `json:"exposure_mode"` // Exposure mode: 1=auto, 2=shutter_priority, 3=aperture_priority, 4=manual
}

// ExposureModeSetCommand represents the camera exposure mode set request
type ExposureModeSetCommand struct {
	common.Header
	MethodName string                    `json:"method"`
	DataValue  ExposureModeSetData `json:"data"`
}

// NewExposureModeSetCommand creates a new camera exposure mode set request
func NewExposureModeSetCommand(data ExposureModeSetData) *ExposureModeSetCommand {
	return &ExposureModeSetCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_exposure_mode_set",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *ExposureModeSetCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *ExposureModeSetCommand) Data() any {
	return c.DataValue
}

// ExposureSetData represents the camera exposure set data
type ExposureSetData struct {
	PayloadIndex  string `json:"payload_index"`  // Camera enumeration value
	CameraType    string `json:"camera_type"`    // Camera type: wide, zoom
	ExposureValue string `json:"exposure_value"` // Exposure value: 1=-5.0EV, 2=-4.7EV, ..., 16=0EV, ..., 31=5.0EV, 255=FIXED
}

// ExposureSetCommand represents the camera exposure set request
type ExposureSetCommand struct {
	common.Header
	MethodName string                `json:"method"`
	DataValue  ExposureSetData `json:"data"`
}

// NewExposureSetCommand creates a new camera exposure set request
func NewExposureSetCommand(data ExposureSetData) *ExposureSetCommand {
	return &ExposureSetCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_exposure_set",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *ExposureSetCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *ExposureSetCommand) Data() any {
	return c.DataValue
}

// FocusModeSetData represents the camera focus mode set data
type FocusModeSetData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
	CameraType   string `json:"camera_type"`   // Camera type: wide, zoom (M30 series only supports zoom)
	FocusMode    int    `json:"focus_mode"`    // Focus mode: 0=MF (manual), 1=AFS (auto single), 2=AFC (auto continuous)
}

// FocusModeSetCommand represents the camera focus mode set request
type FocusModeSetCommand struct {
	common.Header
	MethodName string                 `json:"method"`
	DataValue  FocusModeSetData `json:"data"`
}

// NewFocusModeSetCommand creates a new camera focus mode set request
func NewFocusModeSetCommand(data FocusModeSetData) *FocusModeSetCommand {
	return &FocusModeSetCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_focus_mode_set",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *FocusModeSetCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *FocusModeSetCommand) Data() any {
	return c.DataValue
}

// FocusValueSetData represents the camera focus value set data
type FocusValueSetData struct {
	PayloadIndex string `json:"payload_index"` // Camera enumeration value
	CameraType   string `json:"camera_type"`   // Camera type: wide, zoom (M30 series only supports zoom)
	FocusValue   int    `json:"focus_value"`   // Focus value (range from zoom_min_focus_value to zoom_max_focus_value in OSD)
}

// FocusValueSetCommand represents the camera focus value set request
type FocusValueSetCommand struct {
	common.Header
	MethodName string                  `json:"method"`
	DataValue  FocusValueSetData `json:"data"`
}

// NewFocusValueSetCommand creates a new camera focus value set request
func NewFocusValueSetCommand(data FocusValueSetData) *FocusValueSetCommand {
	return &FocusValueSetCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_focus_value_set",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *FocusValueSetCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *FocusValueSetCommand) Data() any {
	return c.DataValue
}

// PointFocusActionData represents the camera point focus action data
type PointFocusActionData struct {
	PayloadIndex string  `json:"payload_index"` // Camera enumeration value
	CameraType   string  `json:"camera_type"`   // Camera type: wide, zoom (M30 series only supports zoom)
	X            float64 `json:"x"`             // Focus point coordinate x (0-1)
	Y            float64 `json:"y"`             // Focus point coordinate y (0-1)
}

// PointFocusActionCommand represents the camera point focus action request
type PointFocusActionCommand struct {
	common.Header
	MethodName string                     `json:"method"`
	DataValue  PointFocusActionData `json:"data"`
}

// NewPointFocusActionCommand creates a new camera point focus action request
func NewPointFocusActionCommand(data PointFocusActionData) *PointFocusActionCommand {
	return &PointFocusActionCommand{
		Header:     common.NewHeader(),
		MethodName: "camera_point_focus_action",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *PointFocusActionCommand) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *PointFocusActionCommand) Data() any {
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
func (c *AimCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *ExposureModeSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *ExposureSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *FocalLengthSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *FrameZoomCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *FocusModeSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *FocusValueSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *LookAtCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *PhotoStopCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *PhotoTakeCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *PointFocusActionCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *RecordingStartCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *RecordingStopCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *ScreenDragCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *ScreenSplitCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *GimbalResetCommand) GetHeader() *common.Header {
	return &c.Header
}
