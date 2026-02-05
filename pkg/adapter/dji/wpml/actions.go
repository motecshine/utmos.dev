package wpml

type PayloadPosition int

const (
	PayloadPosition0 PayloadPosition = 0
	PayloadPosition1 PayloadPosition = 1
	PayloadPosition2 PayloadPosition = 2
	PayloadPosition7 PayloadPosition = 7
)

type ActionInterface interface {
	GetActionType() string
}

type TakePhotoAction struct {
	PayloadPositionIndex      PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	FileSuffix                string          `json:"file_suffix,omitempty"`
	UseGlobalPayloadLensIndex bool            `json:"use_global_payload_lens_index"`
	PayloadLensIndex          *string         `json:"payload_lens_index,omitempty"`
}

func (a *TakePhotoAction) GetActionType() string {
	return ActionTypeTakePhoto
}

type StartRecordAction struct {
	PayloadPositionIndex      PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	FileSuffix                string          `json:"file_suffix,omitempty"`
	UseGlobalPayloadLensIndex bool            `json:"use_global_payload_lens_index"`
	PayloadLensIndex          *string         `json:"payload_lens_index,omitempty"`
}

func (a *StartRecordAction) GetActionType() string {
	return ActionTypeStartRecord
}

type StopRecordAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	PayloadLensIndex     *string         `json:"payload_lens_index,omitempty"`
}

func (a *StopRecordAction) GetActionType() string {
	return ActionTypeStopRecord
}

type FocusAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	IsPointFocus         bool            `json:"is_point_focus"`
	FocusX               float64         `json:"focus_x" validate:"required,min=0,max=1"`
	FocusY               float64         `json:"focus_y" validate:"required,min=0,max=1"`
	IsInfiniteFocus      bool            `json:"is_infinite_focus"`
	FocusRegionWidth     *float64        `json:"focus_region_width,omitempty" validate:"omitempty,min=0,max=1"`
	FocusRegionHeight    *float64        `json:"focus_region_height,omitempty" validate:"omitempty,min=0,max=1"`
}

func (a *FocusAction) GetActionType() string {
	return ActionTypeFocus
}

type ZoomAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	FocalLength          float64         `json:"focal_length" validate:"required,gt=0"`
}

func (a *ZoomAction) GetActionType() string {
	return ActionTypeZoom
}

type CustomDirNameAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	DirectoryName        string          `json:"directory_name" validate:"required"`
}

func (a *CustomDirNameAction) GetActionType() string {
	return ActionTypeCustomDirName
}

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

func (a *GimbalRotateAction) GetActionType() string {
	return ActionTypeGimbalRotate
}

type RotateYawAction struct {
	AircraftHeading  float64 `json:"aircraft_heading" validate:"min=-180,max=180"`
	AircraftPathMode *string `json:"aircraft_path_mode,omitempty"`
}

func (a *RotateYawAction) GetActionType() string {
	return ActionTypeRotateYaw
}

type HoverAction struct {
	HoverTime float64 `json:"hover_time" validate:"required,gt=0"`
}

func (a *HoverAction) GetActionType() string {
	return ActionTypeHover
}

type GimbalEvenlyRotateAction struct {
	GimbalPitchRotateAngle float64         `json:"gimbal_pitch_rotate_angle"`
	PayloadPositionIndex   PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

func (a *GimbalEvenlyRotateAction) GetActionType() string {
	return ActionTypeGimbalEvenlyRotate
}

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

func (a *OrientedShootAction) GetActionType() string {
	return ActionTypeOrientedShoot
}

type PanoShotAction struct {
	PayloadPositionIndex      PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	UseGlobalPayloadLensIndex bool            `json:"use_global_payload_lens_index"`
	PanoShotSubMode           string          `json:"pano_shot_sub_mode" validate:"required"`
	PayloadLensIndex          *string         `json:"payload_lens_index,omitempty"`
}

func (a *PanoShotAction) GetActionType() string {
	return ActionTypePanoShot
}

type RecordPointCloudAction struct {
	PayloadPositionIndex    PayloadPosition `json:"payload_position_index" validate:"payload_position"`
	RecordPointCloudOperate string          `json:"record_point_cloud_operate" validate:"required"`
}

func (a *RecordPointCloudAction) GetActionType() string {
	return ActionTypeRecordPointCloud
}

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

func (a *AccurateShootAction) GetActionType() string {
	return ActionTypeAccurateShoot
}

type GimbalAngleLockAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

func (a *GimbalAngleLockAction) GetActionType() string {
	return ActionTypeGimbalAngleLock
}

type GimbalAngleUnlockAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

func (a *GimbalAngleUnlockAction) GetActionType() string {
	return ActionTypeGimbalAngleUnlock
}

type StartSmartObliqueAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

func (a *StartSmartObliqueAction) GetActionType() string {
	return ActionTypeStartSmartOblique
}

type StartTimeLapseAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

func (a *StartTimeLapseAction) GetActionType() string {
	return ActionTypeStartTimeLapse
}

type StopTimeLapseAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

func (a *StopTimeLapseAction) GetActionType() string {
	return ActionTypeStopTimeLapse
}

type SetFocusTypeAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

func (a *SetFocusTypeAction) GetActionType() string {
	return ActionTypeSetFocusType
}

type TargetDetectionAction struct {
	PayloadPositionIndex PayloadPosition `json:"payload_position_index" validate:"payload_position"`
}

func (a *TargetDetectionAction) GetActionType() string {
	return ActionTypeTargetDetection
}

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
