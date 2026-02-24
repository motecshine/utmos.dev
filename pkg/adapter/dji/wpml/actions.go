package wpml

// PayloadPosition represents the position index of a payload on the drone.
type PayloadPosition int

// Payload position constants for identifying gimbal mount positions.
const (
	// PayloadPosition0 is the primary gimbal position (position 0).
	PayloadPosition0 PayloadPosition = 0
	// PayloadPosition1 is the secondary gimbal position (position 1).
	PayloadPosition1 PayloadPosition = 1
	// PayloadPosition2 is the tertiary gimbal position (position 2).
	PayloadPosition2 PayloadPosition = 2
	// PayloadPosition7 is the auxiliary gimbal position (position 7).
	PayloadPosition7 PayloadPosition = 7
)

// ActionInterface is the interface that all WPML action types must implement.
type ActionInterface interface {
	GetActionType() string
}

// TakePhotoAction represents an action to capture a photo with the specified payload.
type TakePhotoAction struct {
	PayloadPositionIndex      PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	FileSuffix                string          `json:"file_suffix,omitempty"`
	UseGlobalPayloadLensIndex bool            `json:"use_global_payload_lens_index"`
	PayloadLensIndex          *string         `json:"payload_lens_index,omitempty"`
}

// GetActionType returns the action type identifier for TakePhotoAction.
func (a *TakePhotoAction) GetActionType() string {
	return ActionTypeTakePhoto
}

// StartRecordAction represents an action to start video recording with the specified payload.
type StartRecordAction struct {
	PayloadPositionIndex      PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	FileSuffix                string          `json:"file_suffix,omitempty"`
	UseGlobalPayloadLensIndex bool            `json:"use_global_payload_lens_index"`
	PayloadLensIndex          *string         `json:"payload_lens_index,omitempty"`
}

// GetActionType returns the action type identifier for StartRecordAction.
func (a *StartRecordAction) GetActionType() string {
	return ActionTypeStartRecord
}

// StopRecordAction represents an action to stop video recording with the specified payload.
type StopRecordAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	PayloadLensIndex     *string         `json:"payload_lens_index,omitempty"`
}

// GetActionType returns the action type identifier for StopRecordAction.
func (a *StopRecordAction) GetActionType() string {
	return ActionTypeStopRecord
}

// FocusAction represents an action to adjust the camera focus at a specified point.
type FocusAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	IsPointFocus         bool            `json:"is_point_focus"`
	FocusX               float64         `json:"focus_x" validate:"required,min=0,max=1"`
	FocusY               float64         `json:"focus_y" validate:"required,min=0,max=1"`
	IsInfiniteFocus      bool            `json:"is_infinite_focus"`
	FocusRegionWidth     *float64        `json:"focus_region_width,omitempty" validate:"omitempty,min=0,max=1"`
	FocusRegionHeight    *float64        `json:"focus_region_height,omitempty" validate:"omitempty,min=0,max=1"`
}

// GetActionType returns the action type identifier for FocusAction.
func (a *FocusAction) GetActionType() string {
	return ActionTypeFocus
}

// ZoomAction represents an action to adjust the camera zoom to a specified focal length.
type ZoomAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	FocalLength          float64         `json:"focal_length" validate:"required,gt=0"`
}

// GetActionType returns the action type identifier for ZoomAction.
func (a *ZoomAction) GetActionType() string {
	return ActionTypeZoom
}

// CustomDirNameAction represents an action to set a custom directory name for file storage.
type CustomDirNameAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	DirectoryName        string          `json:"directory_name" validate:"required"`
}

// GetActionType returns the action type identifier for CustomDirNameAction.
func (a *CustomDirNameAction) GetActionType() string {
	return ActionTypeCustomDirName
}

// GimbalRotateAction represents an action to rotate the gimbal to specified angles.
type GimbalRotateAction struct {
	PayloadPositionIndex    PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	GimbalHeadingYawBase    string          `json:"gimbal_heading_yaw_base" validate:"required"`
	GimbalRotateMode        string          `json:"gimbal_rotate_mode" validate:"required"`
	GimbalPitchRotateEnable bool            `json:"gimbal_pitch_rotate_enable"`
	GimbalPitchRotateAngle  float64         `json:"gimbal_pitch_rotate_angle"`
	GimbalRollRotateEnable  bool            `json:"gimbal_roll_rotate_enable"`
	GimbalRollRotateAngle   float64         `json:"gimbal_roll_rotate_angle"`
	GimbalYawRotateEnable   bool            `json:"gimbal_yaw_rotate_enable"`
	GimbalYawRotateAngle    float64         `json:"gimbal_yaw_rotate_angle"`
	GimbalRotateTimeEnable  bool            `json:"gimbal_rotate_time_enable"`
	GimbalRotateTime        float64         `json:"gimbal_rotate_time" validate:"min=0"`
}

// GetActionType returns the action type identifier for GimbalRotateAction.
func (a *GimbalRotateAction) GetActionType() string {
	return ActionTypeGimbalRotate
}

// RotateYawAction represents an action to rotate the aircraft yaw to a specified heading.
type RotateYawAction struct {
	AircraftHeading  float64 `json:"aircraft_heading" validate:"min=-180,max=180"`
	AircraftPathMode *string `json:"aircraft_path_mode,omitempty"`
}

// GetActionType returns the action type identifier for RotateYawAction.
func (a *RotateYawAction) GetActionType() string {
	return ActionTypeRotateYaw
}

// HoverAction represents an action to hover the drone at the current position for a specified duration.
type HoverAction struct {
	HoverTime float64 `json:"hover_time" validate:"required,gt=0"`
}

// GetActionType returns the action type identifier for HoverAction.
func (a *HoverAction) GetActionType() string {
	return ActionTypeHover
}

// GimbalEvenlyRotateAction represents an action to evenly rotate the gimbal pitch between waypoints.
type GimbalEvenlyRotateAction struct {
	GimbalPitchRotateAngle float64         `json:"gimbal_pitch_rotate_angle"`
	PayloadPositionIndex   PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

// GetActionType returns the action type identifier for GimbalEvenlyRotateAction.
func (a *GimbalEvenlyRotateAction) GetActionType() string {
	return ActionTypeGimbalEvenlyRotate
}

// OrientedShootAction represents an action to capture an oriented photograph with precise camera positioning.
type OrientedShootAction struct {
	GimbalPitchRotateAngle    float64         `json:"gimbal_pitch_rotate_angle"`
	GimbalYawRotateAngle      float64         `json:"gimbal_yaw_rotate_angle"`
	FocusX                    int             `json:"focus_x" validate:"required,min=0,max=960"`
	FocusY                    int             `json:"focus_y" validate:"required,min=0,max=720"`
	FocusRegionWidth          int             `json:"focus_region_width" validate:"required,gt=0,max=960"`
	FocusRegionHeight         int             `json:"focus_region_height" validate:"required,gt=0,max=720"`
	FocalLength               float64         `json:"focal_length" validate:"required,gt=0"`
	AircraftHeading           float64         `json:"aircraft_heading" validate:"required,min=-180,max=180"`
	AccurateFrameValid        bool            `json:"accurate_frame_valid"`
	PayloadPositionIndex      PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	UseGlobalPayloadLensIndex bool            `json:"use_global_payload_lens_index"`
	TargetAngle               float64         `json:"target_angle" validate:"required,min=0,max=360"`
	ActionUUID                string          `json:"action_uuid" validate:"required"`
	ImageWidth                int             `json:"image_width" validate:"required,gt=0"`
	ImageHeight               int             `json:"image_height" validate:"required,gt=0"`
	AFPos                     int             `json:"af_pos" validate:"min=0"`
	GimbalPort                int             `json:"gimbal_port" validate:"min=0"`
	OrientedCameraType        int             `json:"oriented_camera_type" validate:"min=0"`
	OrientedFilePath          string          `json:"oriented_file_path" validate:"required"`
	OrientedFileMD5           string          `json:"oriented_file_md5" validate:"required"`
	OrientedFileSize          int             `json:"oriented_file_size" validate:"required,gt=0"`
	OrientedFileSuffix        string          `json:"oriented_file_suffix,omitempty"`
	OrientedCameraApertue     int             `json:"oriented_camera_apertue" validate:"min=0"`
	OrientedCameraLuminance   int             `json:"oriented_camera_luminance" validate:"min=0"`
	OrientedCameraShutterTime float64         `json:"oriented_camera_shutter_time" validate:"required,gt=0"`
	OrientedCameraISO         int             `json:"oriented_camera_iso" validate:"min=0"`
	OrientedPhotoMode         string          `json:"oriented_photo_mode" validate:"required"`
	PayloadLensIndex          *string         `json:"payload_lens_index,omitempty"`
}

// GetActionType returns the action type identifier for OrientedShootAction.
func (a *OrientedShootAction) GetActionType() string {
	return ActionTypeOrientedShoot
}

// PanoShotAction represents an action to capture a panoramic photograph.
type PanoShotAction struct {
	PayloadPositionIndex      PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	UseGlobalPayloadLensIndex bool            `json:"use_global_payload_lens_index"`
	PanoShotSubMode           string          `json:"pano_shot_sub_mode" validate:"required"`
	PayloadLensIndex          *string         `json:"payload_lens_index,omitempty"`
}

// GetActionType returns the action type identifier for PanoShotAction.
func (a *PanoShotAction) GetActionType() string {
	return ActionTypePanoShot
}

// RecordPointCloudAction represents an action to start or stop point cloud recording.
type RecordPointCloudAction struct {
	PayloadPositionIndex    PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	RecordPointCloudOperate string          `json:"record_point_cloud_operate" validate:"required"`
}

// GetActionType returns the action type identifier for RecordPointCloudAction.
func (a *RecordPointCloudAction) GetActionType() string {
	return ActionTypeRecordPointCloud
}

// AccurateShootAction represents an action to capture an accurate photograph with precise positioning parameters.
type AccurateShootAction struct {
	GimbalPitchRotateAngle    float64         `json:"gimbal_pitch_rotate_angle"`
	GimbalYawRotateAngle      float64         `json:"gimbal_yaw_rotate_angle"`
	FocusX                    int             `json:"focus_x" validate:"required,min=0,max=960"`
	FocusY                    int             `json:"focus_y" validate:"required,min=0,max=720"`
	FocusRegionWidth          int             `json:"focus_region_width" validate:"required,gt=0,max=960"`
	FocusRegionHeight         int             `json:"focus_region_height" validate:"required,gt=0,max=720"`
	FocalLength               float64         `json:"focal_length" validate:"required,gt=0"`
	AircraftHeading           float64         `json:"aircraft_heading" validate:"required,min=-180,max=180"`
	AccurateFrameValid        bool            `json:"accurate_frame_valid"`
	PayloadPositionIndex      PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	UseGlobalPayloadLensIndex bool            `json:"use_global_payload_lens_index"`
	TargetAngle               float64         `json:"target_angle" validate:"required,min=0,max=360"`
	ImageWidth                int             `json:"image_width" validate:"required,gt=0"`
	ImageHeight               int             `json:"image_height" validate:"required,gt=0"`
	AFPos                     int             `json:"af_pos" validate:"min=0"`
	GimbalPort                int             `json:"gimbal_port" validate:"min=0"`
	AccurateCameraType        int             `json:"accurate_camera_type" validate:"min=0"`
	AccurateFilePath          string          `json:"accurate_file_path" validate:"required"`
	AccurateFileMD5           string          `json:"accurate_file_md5" validate:"required"`
	AccurateFileSize          int             `json:"accurate_file_size" validate:"required,gt=0"`
	AccurateFileSuffix        string          `json:"accurate_file_suffix,omitempty"`
	AccurateCameraApertue     int             `json:"accurate_camera_apertue" validate:"min=0"`
	AccurateCameraLuminance   int             `json:"accurate_camera_luminance" validate:"min=0"`
	AccurateCameraShutterTime float64         `json:"accurate_camera_shutter_time" validate:"required,gt=0"`
	AccurateCameraISO         int             `json:"accurate_camera_iso" validate:"min=0"`
	PayloadLensIndex          *string         `json:"payload_lens_index,omitempty"`
}

// GetActionType returns the action type identifier for AccurateShootAction.
func (a *AccurateShootAction) GetActionType() string {
	return ActionTypeAccurateShoot
}

// GimbalAngleLockAction represents an action to lock the gimbal at its current angle.
type GimbalAngleLockAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

// GetActionType returns the action type identifier for GimbalAngleLockAction.
func (a *GimbalAngleLockAction) GetActionType() string {
	return ActionTypeGimbalAngleLock
}

// GimbalAngleUnlockAction represents an action to unlock the gimbal angle.
type GimbalAngleUnlockAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

// GetActionType returns the action type identifier for GimbalAngleUnlockAction.
func (a *GimbalAngleUnlockAction) GetActionType() string {
	return ActionTypeGimbalAngleUnlock
}

// StartSmartObliqueAction represents an action to start smart oblique photography.
type StartSmartObliqueAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

// GetActionType returns the action type identifier for StartSmartObliqueAction.
func (a *StartSmartObliqueAction) GetActionType() string {
	return ActionTypeStartSmartOblique
}

// StartTimeLapseAction represents an action to start time-lapse recording.
type StartTimeLapseAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

// GetActionType returns the action type identifier for StartTimeLapseAction.
func (a *StartTimeLapseAction) GetActionType() string {
	return ActionTypeStartTimeLapse
}

// StopTimeLapseAction represents an action to stop time-lapse recording.
type StopTimeLapseAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

// GetActionType returns the action type identifier for StopTimeLapseAction.
func (a *StopTimeLapseAction) GetActionType() string {
	return ActionTypeStopTimeLapse
}

// SetFocusTypeAction represents an action to set the camera focus type.
type SetFocusTypeAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

// GetActionType returns the action type identifier for SetFocusTypeAction.
func (a *SetFocusTypeAction) GetActionType() string {
	return ActionTypeSetFocusType
}

// TargetDetectionAction represents an action to perform target detection using the payload.
type TargetDetectionAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

// GetActionType returns the action type identifier for TargetDetectionAction.
func (a *TargetDetectionAction) GetActionType() string {
	return ActionTypeTargetDetection
}

// CreateActionFromType creates and returns a new action instance based on the given action type string.
// Returns nil if the action type is not recognized.
func CreateActionFromType(actionType string) ActionInterface {
	switch actionType {
	case ActionTypeTakePhoto:
		return &TakePhotoAction{}
	case ActionTypeStartRecord:
		return &StartRecordAction{}
	case ActionTypeStopRecord:
		return &StopRecordAction{}
	case ActionTypeFocus:
		return &FocusAction{}
	case ActionTypeZoom:
		return &ZoomAction{}
	case ActionTypeCustomDirName:
		return &CustomDirNameAction{}
	case ActionTypeGimbalRotate:
		return &GimbalRotateAction{}
	case ActionTypeRotateYaw:
		return &RotateYawAction{}
	case ActionTypeHover:
		return &HoverAction{}
	case ActionTypeGimbalEvenlyRotate:
		return &GimbalEvenlyRotateAction{}
	case ActionTypeOrientedShoot:
		return &OrientedShootAction{}
	case ActionTypePanoShot:
		return &PanoShotAction{}
	case ActionTypeRecordPointCloud:
		return &RecordPointCloudAction{}
	case ActionTypeAccurateShoot:
		return &AccurateShootAction{}
	case ActionTypeGimbalAngleLock:
		return &GimbalAngleLockAction{}
	case ActionTypeGimbalAngleUnlock:
		return &GimbalAngleUnlockAction{}
	case ActionTypeStartSmartOblique:
		return &StartSmartObliqueAction{}
	case ActionTypeStartTimeLapse:
		return &StartTimeLapseAction{}
	case ActionTypeStopTimeLapse:
		return &StopTimeLapseAction{}
	case ActionTypeSetFocusType:
		return &SetFocusTypeAction{}
	case ActionTypeTargetDetection:
		return &TargetDetectionAction{}
	default:
		return nil
	}
}
