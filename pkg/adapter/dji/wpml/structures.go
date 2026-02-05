package wpml

import (
	"time"

	"github.com/nbio/xml"
)

type Document struct {
	XMLName  xml.Name        `xml:"kml" json:"xml_name"`
	XMLNS    string          `xml:"xmlns,attr" json:"xmlns"`
	WPMLNS   string          `xml:"xmlns:wpml,attr" json:"wpml_ns"`
	Document DocumentContent `xml:"Document" json:"document"`
}

type DocumentContent struct {
	Author        string         `xml:"wpml:author,omitempty" json:"author,omitempty"`
	CreateTime    int64          `xml:"wpml:createTime,omitempty" json:"create_time,omitempty"`
	UpdateTime    int64          `xml:"wpml:updateTime,omitempty" json:"update_time,omitempty"`
	MissionConfig *MissionConfig `xml:"wpml:missionConfig,omitempty" json:"mission_config,omitempty"`
	Folders       []Folder       `xml:"Folder,omitempty" json:"folders,omitempty"`
}

type MissionConfig struct {
	FlyToWaylineMode          FlightMode           `xml:"wpml:flyToWaylineMode" validate:"required" json:"fly_to_wayline_mode"`
	FinishAction              FinishAction         `xml:"wpml:finishAction" validate:"required" json:"finish_action"`
	ExitOnRCLost              RCLostAction         `xml:"wpml:exitOnRCLost" validate:"required" json:"exit_on_rc_lost"`
	ExecuteRCLostAction       *ExecuteRCLostAction `xml:"wpml:executeRCLostAction,omitempty" json:"execute_rc_lost_action,omitempty"`
	TakeOffSecurityHeight     float64              `xml:"wpml:takeOffSecurityHeight" validate:"required,min=1.2,max=1500" json:"take_off_security_height"`
	TakeOffRefPoint           *string              `xml:"wpml:takeOffRefPoint,omitempty" json:"take_off_ref_point,omitempty"`
	TakeOffRefPointAGLHeight  *float64             `xml:"wpml:takeOffRefPointAGLHeight,omitempty" json:"take_off_ref_point_agl_height,omitempty"`
	GlobalTransitionalSpeed   float64              `xml:"wpml:globalTransitionalSpeed" validate:"required,min=1,max=15" json:"global_transitional_speed"`
	GlobalRTHHeight           *float64             `xml:"wpml:globalRTHHeight,omitempty" json:"global_rth_height,omitempty"`
	DroneInfo                 DroneInfo            `xml:"wpml:droneInfo" validate:"required" json:"drone_info"`
	WaylineAvoidLimitAreaMode *int                 `xml:"wpml:waylineAvoidLimitAreaMode,omitempty" json:"wayline_avoid_limit_area_mode,omitempty"`
	PayloadInfo               PayloadInfo          `xml:"wpml:payloadInfo" validate:"required" json:"payload_info"`
	AutoRerouteInfo           *AutoRerouteInfo     `xml:"wpml:autoRerouteInfo,omitempty" json:"auto_reroute_info,omitempty"`
}

type DroneInfo struct {
	DroneEnumValue    int `xml:"wpml:droneEnumValue" validate:"required" json:"drone_enum_value"`
	DroneSubEnumValue int `xml:"wpml:droneSubEnumValue" validate:"gte=0" json:"drone_sub_enum_value"`
}

type PayloadInfo struct {
	PayloadEnumValue     int  `xml:"wpml:payloadEnumValue" validate:"required" json:"payload_enum_value"`
	PayloadSubEnumValue  *int `xml:"wpml:payloadSubEnumValue,omitempty" json:"payload_sub_enum_value,omitempty"`
	PayloadPositionIndex int  `xml:"wpml:payloadPositionIndex" validate:"gte=0" json:"payload_position_index"`
}

type AutoRerouteInfo struct {
	MissionAutoRerouteMode      int `xml:"wpml:missionAutoRerouteMode" validate:"min=0,max=1" json:"mission_auto_reroute_mode"`
	TransitionalAutoRerouteMode int `xml:"wpml:transitionalAutoRerouteMode" validate:"min=0,max=1" json:"transitional_auto_reroute_mode"`
}

type Folder struct {
	TemplateType              *TemplateType              `xml:"wpml:templateType,omitempty" json:"template_type,omitempty"`
	TemplateID                int                        `xml:"wpml:templateId" validate:"min=0,max=65535" json:"template_id"`
	WaylineCoordinateSysParam *WaylineCoordinateSysParam `xml:"wpml:waylineCoordinateSysParam,omitempty" json:"wayline_coordinate_sys_param,omitempty"`
	AutoFlightSpeed           float64                    `xml:"wpml:autoFlightSpeed" validate:"required,min=1,max=15" json:"auto_flight_speed"`
	PayloadParam              *PayloadParam              `xml:"wpml:payloadParam,omitempty" json:"payload_param,omitempty"`

	GlobalWaypointTurnMode     *string                     `xml:"wpml:globalWaypointTurnMode,omitempty" json:"global_waypoint_turn_mode,omitempty"`
	GlobalUseStraightLine      *int                        `xml:"wpml:globalUseStraightLine,omitempty" json:"global_use_straight_line,omitempty"`
	GimbalPitchMode            *string                     `xml:"wpml:gimbalPitchMode,omitempty" json:"gimbal_pitch_mode,omitempty"`
	GlobalHeight               *float64                    `xml:"wpml:globalHeight,omitempty" json:"global_height,omitempty"`
	GlobalWaypointHeadingParam *GlobalWaypointHeadingParam `xml:"wpml:globalWaypointHeadingParam,omitempty" json:"global_waypoint_heading_param,omitempty"`

	CaliFlightEnable        *int                 `xml:"wpml:caliFlightEnable,omitempty" json:"cali_flight_enable,omitempty"`
	ElevationOptimizeEnable *int                 `xml:"wpml:elevationOptimizeEnable,omitempty" json:"elevation_optimize_enable,omitempty"`
	SmartObliqueEnable      *int                 `xml:"wpml:smartObliqueEnable,omitempty" json:"smart_oblique_enable,omitempty"`
	SmartObliqueGimbalPitch *int                 `xml:"wpml:smartObliqueGimbalPitch,omitempty" json:"smart_oblique_gimbal_pitch,omitempty"`
	ShootType               *string              `xml:"wpml:shootType,omitempty" json:"shoot_type,omitempty"`
	Direction               *int                 `xml:"wpml:direction,omitempty" json:"direction,omitempty"`
	Margin                  *int                 `xml:"wpml:margin,omitempty" json:"margin,omitempty"`
	Overlap                 *Overlap             `xml:"wpml:overlap,omitempty" json:"overlap,omitempty"`
	EllipsoidHeight         *float64             `xml:"wpml:ellipsoidHeight,omitempty" json:"ellipsoid_height,omitempty"`
	Height                  *float64             `xml:"wpml:height,omitempty" json:"height,omitempty"`
	FacadeWaylineEnable     *int                 `xml:"wpml:facadeWaylineEnable,omitempty" json:"facade_wayline_enable,omitempty"`
	MappingHeadingParam     *MappingHeadingParam `xml:"wpml:mappingHeadingParam,omitempty" json:"mapping_heading_param,omitempty"`
	GimbalPitchAngle        *float64             `xml:"wpml:gimbalPitchAngle,omitempty" json:"gimbal_pitch_angle,omitempty"`

	InclinedGimbalPitch *int     `xml:"wpml:inclinedGimbalPitch,omitempty" json:"inclined_gimbal_pitch,omitempty"`
	InclinedFlightSpeed *float64 `xml:"wpml:inclinedFlightSpeed,omitempty" json:"inclined_flight_speed,omitempty"`

	SingleLineEnable         *int     `xml:"wpml:singleLineEnable,omitempty" json:"single_line_enable,omitempty"`
	CuttingDistance          *float64 `xml:"wpml:cuttingDistance,omitempty" json:"cutting_distance,omitempty"`
	BoundaryOptimEnable      *int     `xml:"wpml:boundaryOptimEnable,omitempty" json:"boundary_optim_enable,omitempty"`
	LeftExtend               *int     `xml:"wpml:leftExtend,omitempty" json:"left_extend,omitempty"`
	RightExtend              *int     `xml:"wpml:rightExtend,omitempty" json:"right_extend,omitempty"`
	IncludeCenterEnable      *int     `xml:"wpml:includeCenterEnable,omitempty" json:"include_center_enable,omitempty"`
	StripUseTemplateAltitude *int     `xml:"wpml:stripUseTemplateAltitude,omitempty" json:"strip_use_template_altitude,omitempty"`

	ExecuteHeightMode *ExecuteHeightMode `xml:"wpml:executeHeightMode,omitempty" json:"execute_height_mode,omitempty"`
	WaylineID         *int               `xml:"wpml:waylineId,omitempty" json:"wayline_id,omitempty"`
	Distance          *float64           `xml:"wpml:distance,omitempty" json:"distance,omitempty"`
	Duration          *float64           `xml:"wpml:duration,omitempty" json:"duration,omitempty"`

	Polygon    *Polygon    `xml:"Polygon,omitempty" json:"polygon,omitempty"`
	LineString *LineString `xml:"LineString,omitempty" json:"line_string,omitempty"`
	Placemarks []Placemark `xml:"Placemark,omitempty" json:"placemarks,omitempty"`

	StartActionGroup *ActionGroup `xml:"wpml:startActionGroup,omitempty" json:"start_action_group,omitempty"`
}

type WaylineCoordinateSysParam struct {
	CoordinateMode          CoordinateMode   `xml:"wpml:coordinateMode" validate:"required" json:"coordinate_mode"`
	HeightMode              HeightMode       `xml:"wpml:heightMode" validate:"required" json:"height_mode"`
	PositioningType         *PositioningType `xml:"wpml:positioningType,omitempty" json:"positioning_type,omitempty"`
	GlobalShootHeight       *float64         `xml:"wpml:globalShootHeight,omitempty" json:"global_shoot_height,omitempty"`
	SurfaceFollowModeEnable *int             `xml:"wpml:surfaceFollowModeEnable,omitempty" json:"surface_follow_mode_enable,omitempty"`
	SurfaceRelativeHeight   *float64         `xml:"wpml:surfaceRelativeHeight,omitempty" json:"surface_relative_height,omitempty"`
}

type PayloadParam struct {
	PayloadPositionIndex int     `xml:"wpml:payloadPositionIndex" validate:"required" json:"payload_position_index"`
	FocusMode            *string `xml:"wpml:focusMode,omitempty" json:"focus_mode,omitempty"`
	MeteringMode         *string `xml:"wpml:meteringMode,omitempty" json:"metering_mode,omitempty"`
	DewarpingEnable      *int    `xml:"wpml:dewarpingEnable,omitempty" json:"dewarping_enable,omitempty"`
	ReturnMode           *string `xml:"wpml:returnMode,omitempty" json:"return_mode,omitempty"`
	SamplingRate         *int    `xml:"wpml:samplingRate,omitempty" json:"sampling_rate,omitempty"`
	ScanningMode         *string `xml:"wpml:scanningMode,omitempty" json:"scanning_mode,omitempty"`
	ModelColoringEnable  *int    `xml:"wpml:modelColoringEnable,omitempty" json:"model_coloring_enable,omitempty"`
	ImageFormat          string  `xml:"wpml:imageFormat" validate:"required" json:"image_format"`
}

type GlobalWaypointHeadingParam struct {
	WaypointHeadingMode     string   `xml:"wpml:waypointHeadingMode" validate:"required" json:"waypoint_heading_mode"`
	WaypointHeadingAngle    *float64 `xml:"wpml:waypointHeadingAngle,omitempty" json:"waypoint_heading_angle,omitempty"`
	WaypointPoiPoint        *string  `xml:"wpml:waypointPoiPoint,omitempty" json:"waypoint_poi_point,omitempty"`
	WaypointHeadingPathMode string   `xml:"wpml:waypointHeadingPathMode" validate:"required" json:"waypoint_heading_path_mode"`
	WaypointHeadingPoiIndex *int     `xml:"wpml:waypointHeadingPoiIndex,omitempty" json:"waypoint_heading_poi_index,omitempty"`
}

type Overlap struct {
	OrthoLidarOverlapH     *int `xml:"wpml:orthoLidarOverlapH,omitempty" json:"ortho_lidar_overlap_h,omitempty"`
	OrthoLidarOverlapW     *int `xml:"wpml:orthoLidarOverlapW,omitempty" json:"ortho_lidar_overlap_w,omitempty"`
	OrthoCameraOverlapH    *int `xml:"wpml:orthoCameraOverlapH,omitempty" json:"ortho_camera_overlap_h,omitempty"`
	OrthoCameraOverlapW    *int `xml:"wpml:orthoCameraOverlapW,omitempty" json:"ortho_camera_overlap_w,omitempty"`
	InclinedLidarOverlapH  *int `xml:"wpml:inclinedLidarOverlapH,omitempty" json:"inclined_lidar_overlap_h,omitempty"`
	InclinedLidarOverlapW  *int `xml:"wpml:inclinedLidarOverlapW,omitempty" json:"inclined_lidar_overlap_w,omitempty"`
	InclinedCameraOverlapH *int `xml:"wpml:inclinedCameraOverlapH,omitempty" json:"inclined_camera_overlap_h,omitempty"`
	InclinedCameraOverlapW *int `xml:"wpml:inclinedCameraOverlapW,omitempty" json:"inclined_camera_overlap_w,omitempty"`
}

type MappingHeadingParam struct {
	MappingHeadingMode  string `xml:"wpml:mappingHeadingMode" validate:"required" json:"mapping_heading_mode"`
	MappingHeadingAngle *int   `xml:"wpml:mappingHeadingAngle,omitempty" json:"mapping_heading_angle,omitempty"`
}

type Polygon struct {
	OuterBoundaryIs OuterBoundaryIs `xml:"outerBoundaryIs" json:"outer_boundary_is"`
}

type OuterBoundaryIs struct {
	LinearRing LinearRing `xml:"LinearRing" json:"linear_ring"`
}

type LinearRing struct {
	Coordinates string `xml:"coordinates" json:"coordinates"`
}

type LineString struct {
	Coordinates string `xml:"coordinates" json:"coordinates"`
}

func NewDocument() *Document {
	return &Document{
		XMLNS:  "http://www.opengis.net/kml/2.2",
		WPMLNS: "http://www.dji.com/wpmz/1.0.6",
		Document: DocumentContent{
			CreateTime: time.Now().UnixMilli(),
			UpdateTime: time.Now().UnixMilli(),
		},
	}
}

func (d *Document) SetAuthor(author string) {
	d.Document.Author = author
}

func (d *Document) UpdateTimestamp() {
	d.Document.UpdateTime = time.Now().UnixMilli()
}
