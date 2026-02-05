package wpml

import "fmt"

type Placemark struct {
	Point                      *Point                      `xml:"Point,omitempty" json:"point,omitempty"`
	Index                      int                         `xml:"wpml:index" validate:"min=0,max=65535" json:"index"`
	EllipsoidHeight            *float64                    `xml:"wpml:ellipsoidHeight,omitempty" json:"ellipsoid_height,omitempty"`
	Height                     *float64                    `xml:"wpml:height,omitempty" json:"height,omitempty"`
	UseGlobalHeight            *int                        `xml:"wpml:useGlobalHeight,omitempty" json:"use_global_height,omitempty"`
	UseGlobalSpeed             *int                        `xml:"wpml:useGlobalSpeed,omitempty" json:"use_global_speed,omitempty"`
	UseGlobalHeadingParam      *int                        `xml:"wpml:useGlobalHeadingParam,omitempty" json:"use_global_heading_param,omitempty"`
	UseGlobalTurnParam         *int                        `xml:"wpml:useGlobalTurnParam,omitempty" json:"use_global_turn_param,omitempty"`
	GimbalPitchAngle           *float64                    `xml:"wpml:gimbalPitchAngle,omitempty" json:"gimbal_pitch_angle,omitempty"`
	ExecuteHeight              *float64                    `xml:"wpml:executeHeight,omitempty" json:"execute_height,omitempty"`
	WaypointSpeed              *float64                    `xml:"wpml:waypointSpeed,omitempty" json:"waypoint_speed,omitempty"`
	WaypointHeadingParam       *WaypointHeadingParam       `xml:"wpml:waypointHeadingParam,omitempty" json:"waypoint_heading_param,omitempty"`
	WaypointTurnParam          *WaypointTurnParam          `xml:"wpml:waypointTurnParam,omitempty" json:"waypoint_turn_param,omitempty"`
	UseStraightLine            *int                        `xml:"wpml:useStraightLine,omitempty" json:"use_straight_line,omitempty"`
	ActionGroups               []ActionGroup               `xml:"wpml:actionGroup,omitempty" json:"action_groups,omitempty"`
	WaypointGimbalHeadingParam *WaypointGimbalHeadingParam `xml:"wpml:waypointGimbalHeadingParam,omitempty" json:"waypoint_gimbal_heading_param,omitempty"`
	IsRisky                    *int                        `xml:"wpml:isRisky,omitempty" json:"is_risky,omitempty"`
	WaypointWorkType           *int                        `xml:"wpml:waypointWorkType,omitempty" json:"waypoint_work_type,omitempty"`
}

type Point struct {
	Coordinates string `xml:"coordinates" validate:"required" json:"coordinates"`
}

type WaypointHeadingParam struct {
	WaypointHeadingMode        string   `xml:"wpml:waypointHeadingMode" validate:"required" json:"waypoint_heading_mode"`
	WaypointHeadingAngle       *float64 `xml:"wpml:waypointHeadingAngle,omitempty" json:"waypoint_heading_angle,omitempty"`
	WaypointPoiPoint           *string  `xml:"wpml:waypointPoiPoint,omitempty" json:"waypoint_poi_point,omitempty"`
	WaypointHeadingAngleEnable *int     `xml:"wpml:waypointHeadingAngleEnable,omitempty" json:"waypoint_heading_angle_enable,omitempty"`
	WaypointHeadingPathMode    string   `xml:"wpml:waypointHeadingPathMode" validate:"required" json:"waypoint_heading_path_mode"`
	WaypointHeadingPoiIndex    *int     `xml:"wpml:waypointHeadingPoiIndex,omitempty" json:"waypoint_heading_poi_index,omitempty"`
}

type WaypointTurnParam struct {
	WaypointTurnMode        string   `xml:"wpml:waypointTurnMode" validate:"required" json:"waypoint_turn_mode"`
	WaypointTurnDampingDist *float64 `xml:"wpml:waypointTurnDampingDist,omitempty" json:"waypoint_turn_damping_dist,omitempty"`
}

type WaypointGimbalHeadingParam struct {
	WaypointGimbalPitchAngle *float64 `xml:"wpml:waypointGimbalPitchAngle,omitempty" json:"waypoint_gimbal_pitch_angle,omitempty"`
	WaypointGimbalYawAngle   *float64 `xml:"wpml:waypointGimbalYawAngle,omitempty" json:"waypoint_gimbal_yaw_angle,omitempty"`
}

type ActionGroup struct {
	ActionGroupID         int           `xml:"wpml:actionGroupId" validate:"required,min=0,max=65535" json:"action_group_id"`
	ActionGroupStartIndex int           `xml:"wpml:actionGroupStartIndex" validate:"required,min=0,max=65535" json:"action_group_start_index"`
	ActionGroupEndIndex   int           `xml:"wpml:actionGroupEndIndex" validate:"required,min=0,max=65535" json:"action_group_end_index"`
	ActionGroupMode       string        `xml:"wpml:actionGroupMode" validate:"required" json:"action_group_mode"`
	ActionTrigger         ActionTrigger `xml:"wpml:actionTrigger" validate:"required" json:"action_trigger"`
	Actions               []Action      `xml:"wpml:action" validate:"required,dive" json:"actions"`
}

type ActionTrigger struct {
	ActionTriggerType  string   `xml:"wpml:actionTriggerType" validate:"required" json:"action_trigger_type"`
	ActionTriggerParam *float64 `xml:"wpml:actionTriggerParam,omitempty" json:"action_trigger_param,omitempty"`
}

type Action struct {
	ActionID                int                      `xml:"wpml:actionId" validate:"required,min=0,max=65535" json:"action_id"`
	ActionActuatorFunc      string                   `xml:"wpml:actionActuatorFunc" validate:"required" json:"action_actuator_func"`
	ActionActuatorFuncParam *ActionActuatorFuncParam `xml:"wpml:actionActuatorFuncParam,omitempty" json:"action_actuator_func_param,omitempty"`
}

func (a Action) GetActionType() string {
	return a.ActionActuatorFunc
}

type ActionActuatorFuncParam struct {
	GimbalHeadingYawBase      *string            `xml:"wpml:gimbalHeadingYawBase,omitempty" json:"gimbal_heading_yaw_base,omitempty"`
	GimbalRotateMode          *string            `xml:"wpml:gimbalRotateMode,omitempty" json:"gimbal_rotate_mode,omitempty"`
	GimbalPitchRotateEnable   *int               `xml:"wpml:gimbalPitchRotateEnable,omitempty" json:"gimbal_pitch_rotate_enable,omitempty"`
	GimbalPitchRotateAngle    *float64           `xml:"wpml:gimbalPitchRotateAngle,omitempty" json:"gimbal_pitch_rotate_angle,omitempty"`
	GimbalRollRotateEnable    *int               `xml:"wpml:gimbalRollRotateEnable,omitempty" json:"gimbal_roll_rotate_enable,omitempty"`
	GimbalRollRotateAngle     *float64           `xml:"wpml:gimbalRollRotateAngle,omitempty" json:"gimbal_roll_rotate_angle,omitempty"`
	GimbalYawRotateEnable     *int               `xml:"wpml:gimbalYawRotateEnable,omitempty" json:"gimbal_yaw_rotate_enable,omitempty"`
	GimbalYawRotateAngle      *float64           `xml:"wpml:gimbalYawRotateAngle,omitempty" json:"gimbal_yaw_rotate_angle,omitempty"`
	GimbalRotateTimeEnable    *int               `xml:"wpml:gimbalRotateTimeEnable,omitempty" json:"gimbal_rotate_time_enable,omitempty"`
	GimbalRotateTime          *float64           `xml:"wpml:gimbalRotateTime,omitempty" json:"gimbal_rotate_time,omitempty"`
	FocalLength               *float64           `xml:"wpml:focalLength,omitempty" json:"focal_length,omitempty"`
	IsUseFocalFactor          *int               `xml:"wpml:isUseFocalFactor,omitempty" json:"is_use_focal_factor,omitempty"`
	PayloadPositionIndex      *int               `xml:"wpml:payloadPositionIndex,omitempty" json:"payload_position_index,omitempty"`
	FileSuffix                *string            `xml:"wpml:fileSuffix,omitempty" json:"file_suffix,omitempty"`
	PayloadLensIndex          *string            `xml:"wpml:payloadLensIndex,omitempty" json:"payload_lens_index,omitempty"`
	UseGlobalPayloadLensIndex *int               `xml:"wpml:useGlobalPayloadLensIndex,omitempty" json:"use_global_payload_lens_index,omitempty"`
	IsPointFocus              *int               `xml:"wpml:isPointFocus,omitempty" json:"is_point_focus,omitempty"`
	FocusX                    *float64           `xml:"wpml:focusX,omitempty" json:"focus_x,omitempty"`
	FocusY                    *float64           `xml:"wpml:focusY,omitempty" json:"focus_y,omitempty"`
	FocusRegionWidth          *float64           `xml:"wpml:focusRegionWidth,omitempty" json:"focus_region_width,omitempty"`
	FocusRegionHeight         *float64           `xml:"wpml:focusRegionHeight,omitempty" json:"focus_region_height,omitempty"`
	IsInfiniteFocus           *int               `xml:"wpml:isInfiniteFocus,omitempty" json:"is_infinite_focus,omitempty"`
	DirectoryName             *string            `xml:"wpml:directoryName,omitempty" json:"directory_name,omitempty"`
	AircraftHeading           *float64           `xml:"wpml:aircraftHeading,omitempty" json:"aircraft_heading,omitempty"`
	AircraftPathMode          *string            `xml:"wpml:aircraftPathMode,omitempty" json:"aircraft_path_mode,omitempty"`
	HoverTime                 *float64           `xml:"wpml:hoverTime,omitempty" json:"hover_time,omitempty"`
	AccurateFrameValid        *int               `xml:"wpml:accurateFrameValid,omitempty" json:"accurate_frame_valid,omitempty"`
	TargetAngle               *float64           `xml:"wpml:targetAngle,omitempty" json:"target_angle,omitempty"`
	ActionUUID                *string            `xml:"wpml:actionUUID,omitempty" json:"action_uuid,omitempty"`
	ImageWidth                *int               `xml:"wpml:imageWidth,omitempty" json:"image_width,omitempty"`
	ImageHeight               *int               `xml:"wpml:imageHeight,omitempty" json:"image_height,omitempty"`
	AFPos                     *int               `xml:"wpml:AFPos,omitempty" json:"af_pos,omitempty"`
	GimbalPort                *int               `xml:"wpml:gimbalPort,omitempty" json:"gimbal_port,omitempty"`
	OrientedCameraType        *int               `xml:"wpml:orientedCameraType,omitempty" json:"oriented_camera_type,omitempty"`
	OrientedFilePath          *string            `xml:"wpml:orientedFilePath,omitempty" json:"oriented_file_path,omitempty"`
	OrientedFileMD5           *string            `xml:"wpml:orientedFileMD5,omitempty" json:"oriented_file_md5,omitempty"`
	OrientedFileSize          *int               `xml:"wpml:orientedFileSize,omitempty" json:"oriented_file_size,omitempty"`
	OrientedFileSuffix        *string            `xml:"wpml:orientedFileSuffix,omitempty" json:"oriented_file_suffix,omitempty"`
	OrientedCameraApertue     *int               `xml:"wpml:orientedCameraApertue,omitempty" json:"oriented_camera_apertue,omitempty"`
	OrientedCameraLuminance   *int               `xml:"wpml:orientedCameraLuminance,omitempty" json:"oriented_camera_luminance,omitempty"`
	OrientedCameraShutterTime *float64           `xml:"wpml:orientedCameraShutterTime,omitempty" json:"oriented_camera_shutter_time,omitempty"`
	OrientedCameraISO         *int               `xml:"wpml:orientedCameraISO,omitempty" json:"oriented_camera_iso,omitempty"`
	OrientedPhotoMode         *string            `xml:"wpml:orientedPhotoMode,omitempty" json:"oriented_photo_mode,omitempty"`
	PanoShotSubMode           *string            `xml:"wpml:panoShotSubMode,omitempty" json:"pano_shot_sub_mode,omitempty"`
	RecordPointCloudOperate   *string            `xml:"wpml:recordPointCloudOperate,omitempty" json:"record_point_cloud_operate,omitempty"`
	MinShootInterval          *float64           `xml:"wpml:minShootInterval,omitempty" json:"min_shoot_interval,omitempty"`
	CameraFocusType           *string            `xml:"wpml:cameraFocusType,omitempty" json:"camera_focus_type,omitempty"`
	SmartObliqueCycleMode     *string            `xml:"wpml:smartObliqueCycleMode,omitempty" json:"smart_oblique_cycle_mode,omitempty"`
	SmartObliquePoint         *SmartObliquePoint `xml:"wpml:smartObliquePoint,omitempty" json:"smart_oblique_point,omitempty"`
}

func (a ActionActuatorFuncParam) GetActionType() string {
	return "actions"
}

type SmartObliquePoint struct {
	SmartObliqueRunningTime *int     `xml:"wpml:smartObliqueRunningTime,omitempty" json:"smart_oblique_running_time,omitempty"`
	SmartObliqueStayTime    *int     `xml:"wpml:smartObliqueStayTime,omitempty" json:"smart_oblique_stay_time,omitempty"`
	SmartObliqueEulerPitch  *float64 `xml:"wpml:smartObliqueEulerPitch,omitempty" json:"smart_oblique_euler_pitch,omitempty"`
	SmartObliqueEulerRoll   *float64 `xml:"wpml:smartObliqueEulerRoll,omitempty" json:"smart_oblique_euler_roll,omitempty"`
	SmartObliqueEulerYaw    *float64 `xml:"wpml:smartObliqueEulerYaw,omitempty" json:"smart_oblique_euler_yaw,omitempty"`
}

const (
	ActionTypeTakePhoto          = "takePhoto"
	ActionTypeStartRecord        = "startRecord"
	ActionTypeStopRecord         = "stopRecord"
	ActionTypeFocus              = "focus"
	ActionTypeZoom               = "zoom"
	ActionTypeCustomDirName      = "customDirName"
	ActionTypeGimbalRotate       = "gimbalRotate"
	ActionTypeRotateYaw          = "rotateYaw"
	ActionTypeHover              = "hover"
	ActionTypeGimbalEvenlyRotate = "gimbalEvenlyRotate"
	ActionTypeAccurateShoot      = "accurateShoot"
	ActionTypeOrientedShoot      = "orientedShoot"
	ActionTypePanoShot           = "panoShot"
	ActionTypeRecordPointCloud   = "recordPointCloud"
	ActionTypeGimbalAngleLock    = "gimbalAngleLock"
	ActionTypeGimbalAngleUnlock  = "gimbalAngleUnlock"
	ActionTypeStartSmartOblique  = "startSmartOblique"
	ActionTypeStartTimeLapse     = "startTimeLapse"
	ActionTypeStopTimeLapse      = "stopTimeLapse"
	ActionTypeSetFocusType       = "setFocusType"
	ActionTypeTargetDetection    = "targetDetection"
)

const (
	TriggerTypeReachPoint            = "reachPoint"
	TriggerTypePassPoint             = "passPoint"
	TriggerTypeManual                = "manual"
	TriggerTypeBetweenAdjacentPoints = "betweenAdjacentPoints"
	TriggerTypeMultipleTiming        = "multipleTiming"
	TriggerTypeMultipleDistance      = "multipleDistance"
)

const (
	ActionGroupModeSequence = "sequence"
)

const (
	TurnModeCoordinateTurn                           = "coordinateTurn"
	TurnModeToPointAndStopWithDiscontinuityCurvature = "toPointAndStopWithDiscontinuityCurvature"
	TurnModeToPointAndStopWithContinuityCurvature    = "toPointAndStopWithContinuityCurvature"
	TurnModeToPointAndPassWithContinuityCurvature    = "toPointAndPassWithContinuityCurvature"
)

const (
	HeadingModeFollowWayline    = "followWayline"
	HeadingModeManually         = "manually"
	HeadingModeFixed            = "fixed"
	HeadingModeSmoothTransition = "smoothTransition"
	HeadingModeTowardPOI        = "towardPOI"
	HeadingModeFree             = "free"
)

const (
	HeadingPathModeClockwise        = "clockwise"
	HeadingPathModeCounterClockwise = "counterClockwise"
	HeadingPathModeFollowBadArc     = "followBadArc"
)

func NewActionGroup(id int, startIndex, endIndex int) *ActionGroup {
	return &ActionGroup{
		ActionGroupID:         id,
		ActionGroupStartIndex: startIndex,
		ActionGroupEndIndex:   endIndex,
		ActionGroupMode:       ActionGroupModeSequence,
		Actions:               make([]Action, 0),
	}
}

func (ag *ActionGroup) AddAction(action Action) {
	ag.Actions = append(ag.Actions, action)
}

func (ag *ActionGroup) SetTrigger(triggerType string, param *float64) {
	ag.ActionTrigger = ActionTrigger{
		ActionTriggerType:  triggerType,
		ActionTriggerParam: param,
	}
}

func NewWaypoint(longitude, latitude float64, index int) *Placemark {
	coordinates := fmt.Sprintf("%.15f,%.15f", longitude, latitude)
	return &Placemark{
		Point: &Point{
			Coordinates: coordinates,
		},
		Index: index,
	}
}

func (p *Placemark) SetHeight(ellipsoidHeight, height float64) {
	p.EllipsoidHeight = &ellipsoidHeight
	p.Height = &height
}

func (p *Placemark) SetExecuteHeight(height float64) {
	p.ExecuteHeight = &height
}

func (p *Placemark) AddActionGroup(actionGroup ActionGroup) {
	p.ActionGroups = append(p.ActionGroups, actionGroup)
}
