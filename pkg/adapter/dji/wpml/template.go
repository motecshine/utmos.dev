package wpml

import (
	"time"

	"github.com/nbio/xml"
)

type WPMLMission struct {
	Template  *TemplateDocument
	Waylines  *WaylinesDocument
	Resources map[string][]byte
}

type TemplateDocument struct {
	XMLName  xml.Name                `xml:"kml" json:"xml_name"`
	XMLNS    string                  `xml:"xmlns,attr" json:"xmlns"`
	WPMLNS   string                  `xml:"xmlns:wpml,attr" json:"wpml_ns"`
	Document TemplateDocumentContent `xml:"Document" json:"document"`
}

type TemplateDocumentContent struct {
	Author     string `xml:"wpml:author,omitempty" json:"author,omitempty"`
	CreateTime int64  `xml:"wpml:createTime,omitempty" json:"create_time,omitempty"`
	UpdateTime int64  `xml:"wpml:updateTime,omitempty" json:"update_time,omitempty"`

	MissionConfig MissionConfig `xml:"wpml:missionConfig" validate:"required" json:"mission_config"`

	Folders []TemplateFolder `xml:"Folder" validate:"required,dive" json:"folders"`
}

type WaylinesDocument struct {
	XMLName  xml.Name                `xml:"kml" json:"xml_name"`
	XMLNS    string                  `xml:"xmlns,attr" json:"xmlns"`
	WPMLNS   string                  `xml:"xmlns:wpml,attr" json:"wpml_ns"`
	Document WaylinesDocumentContent `xml:"Document" json:"document"`
}

type WaylinesDocumentContent struct {
	Folders []WaylineFolder `xml:"Folder" validate:"required,dive" json:"folders"`

	MissionConfig WaylinesMissionConfig `xml:"wpml:missionConfig" validate:"required" json:"mission_config"`
}

type WaylinesMissionConfig struct {
	FlyToWaylineMode          FlightMode           `xml:"wpml:flyToWaylineMode" validate:"required" json:"fly_to_wayline_mode"`
	FinishAction              FinishAction         `xml:"wpml:finishAction" validate:"required" json:"finish_action"`
	ExitOnRCLost              RCLostAction         `xml:"wpml:exitOnRCLost" validate:"required" json:"exit_on_rc_lost"`
	ExecuteRCLostAction       *ExecuteRCLostAction `xml:"wpml:executeRCLostAction,omitempty" json:"execute_rc_lost_action,omitempty"`
	TakeOffSecurityHeight     float64              `xml:"wpml:takeOffSecurityHeight" validate:"required,min=1.2,max=1500" json:"take_off_security_height"`
	GlobalTransitionalSpeed   float64              `xml:"wpml:globalTransitionalSpeed" validate:"required,min=1,max=15" json:"global_transitional_speed"`
	GlobalRTHHeight           float64              `xml:"wpml:globalRTHHeight" validate:"required,min=2,max=1500" json:"global_rth_height"`
	DroneInfo                 DroneInfo            `xml:"wpml:droneInfo" validate:"required" json:"drone_info"`
	WaylineAvoidLimitAreaMode *int                 `xml:"wpml:waylineAvoidLimitAreaMode,omitempty" json:"wayline_avoid_limit_area_mode,omitempty"`
	PayloadInfo               PayloadInfo          `xml:"wpml:payloadInfo" validate:"required" json:"payload_info"`
	AutoRerouteInfo           *AutoRerouteInfo     `xml:"wpml:autoRerouteInfo,omitempty" json:"auto_reroute_info,omitempty"`
}

type TemplateFolder struct {
	TemplateType    TemplateType `xml:"wpml:templateType" validate:"required" json:"template_type"`
	TemplateID      int          `xml:"wpml:templateId" validate:"min=0,max=65535" json:"template_id"`
	AutoFlightSpeed float64      `xml:"wpml:autoFlightSpeed" validate:"required,min=1,max=15" json:"auto_flight_speed"`

	WaylineCoordinateSysParam *WaylineCoordinateSysParam `xml:"wpml:waylineCoordinateSysParam,omitempty" json:"wayline_coordinate_sys_param,omitempty"`
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

	Polygon    *Polygon    `xml:"Polygon,omitempty" json:"polygon,omitempty"`
	LineString *LineString `xml:"LineString,omitempty" json:"line_string,omitempty"`
	Placemarks []Placemark `xml:"Placemark,omitempty" json:"placemarks,omitempty"`
}

type WaylineFolder struct {
	TemplateID        int               `xml:"wpml:templateId" validate:"min=0,max=65535" json:"template_id"`
	WaylineID         int               `xml:"wpml:waylineId" validate:"min=0,max=65535" json:"wayline_id"`
	AutoFlightSpeed   float64           `xml:"wpml:autoFlightSpeed" validate:"required,min=1,max=15" json:"auto_flight_speed"`
	ExecuteHeightMode ExecuteHeightMode `xml:"wpml:executeHeightMode" validate:"required" json:"execute_height_mode"`
	Distance          *float64          `xml:"wpml:distance,omitempty" json:"distance,omitempty"`
	Duration          *float64          `xml:"wpml:duration,omitempty" json:"duration,omitempty"`

	Placemarks []Placemark `xml:"Placemark" validate:"required,dive" json:"placemarks"`

	StartActionGroup *ActionGroup `xml:"wpml:startActionGroup,omitempty" json:"start_action_group,omitempty"`
}

func NewWPMLMission() *WPMLMission {
	now := time.Now().UnixMilli()

	return &WPMLMission{
		Template: &TemplateDocument{
			XMLNS:  "http://www.opengis.net/kml/2.2",
			WPMLNS: "http://www.dji.com/wpmz/1.0.6",
			Document: TemplateDocumentContent{
				CreateTime: now,
				UpdateTime: now,
			},
		},
		Waylines: &WaylinesDocument{
			XMLNS:    "http://www.opengis.net/kml/2.2",
			WPMLNS:   "http://www.dji.com/wpmz/1.0.6",
			Document: WaylinesDocumentContent{},
		},
		Resources: make(map[string][]byte),
	}
}

func (m *WPMLMission) SetAuthor(author string) {
	if m.Template != nil {
		m.Template.Document.Author = author
	}
}

func (m *WPMLMission) UpdateTimestamp() {
	now := time.Now().UnixMilli()
	if m.Template != nil {
		m.Template.Document.UpdateTime = now
	}
}

func (m *WPMLMission) SetMissionConfig(config MissionConfig) {
	if m.Template != nil {
		m.Template.Document.MissionConfig = config
	}

	if m.Waylines != nil {
		waylineAvoidMode := 0

		globalRTHHeight := 100.0
		if config.GlobalRTHHeight != nil {
			globalRTHHeight = *config.GlobalRTHHeight
		}

		m.Waylines.Document.MissionConfig = WaylinesMissionConfig{
			FlyToWaylineMode:          config.FlyToWaylineMode,
			FinishAction:              config.FinishAction,
			ExitOnRCLost:              config.ExitOnRCLost,
			ExecuteRCLostAction:       config.ExecuteRCLostAction,
			TakeOffSecurityHeight:     config.TakeOffSecurityHeight,
			GlobalTransitionalSpeed:   config.GlobalTransitionalSpeed,
			GlobalRTHHeight:           globalRTHHeight,
			DroneInfo:                 config.DroneInfo,
			WaylineAvoidLimitAreaMode: &waylineAvoidMode,
			PayloadInfo:               config.PayloadInfo,
			AutoRerouteInfo:           config.AutoRerouteInfo,
		}
	}
}

func (m *WPMLMission) AddResource(filename string, data []byte) {
	if m.Resources == nil {
		m.Resources = make(map[string][]byte)
	}
	m.Resources[filename] = data
}
