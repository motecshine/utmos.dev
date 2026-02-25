package wpml

// Waylines represents the input schema for converting wayline data into a WPML mission.
type Waylines struct {
	Name                     string             `json:"name" validate:"required,min=1,max=100"`
	Description              string             `json:"description,omitempty" validate:"max=500"`
	DroneModel               DroneModel         `json:"drone_model" validate:"required,drone_model"`
	PayloadModel             PayloadModel       `json:"payload_model" validate:"required,payload_model"`
	PayloadPositionIndex     PayloadPosition    `json:"payload_position_index,omitempty" validate:"payload_position"`
	TemplateType             TemplateType       `json:"template_type" validate:"required"`
	GlobalHeight             float64            `json:"global_height,omitempty" validate:"omitempty,min=5,max=1500"`
	GlobalSpeed              float64            `json:"global_speed,omitempty" validate:"omitempty,min=1,max=15"`
	PhotoSettings            []string           `json:"photo_settings,omitempty" validate:"dive,oneof=wide zoom ir vision"`
	UseLowLightSmart         bool               `json:"use_low_light_smart,omitempty"`
	FinishAction             FinishAction       `json:"finish_action,omitempty"`
	HeightType               HeightMode         `json:"height_type,omitempty" default:"relativeToStartPoint"`
	ClimbMode                string             `json:"climb_mode,omitempty" validate:"omitempty,oneof=vertical inclined"`
	TakeOffSecurityHeight    float64            `json:"take_off_security_height,omitempty" validate:"omitempty,min=20,max=200"`
	GlobalRTHHeight          float64            `json:"global_rth_height,omitempty" validate:"omitempty,min=20,max=1500"`
	AircraftYawMode          string             `json:"aircraft_yaw_mode,omitempty" validate:"omitempty,oneof=followWayline followRoute manual free"`
	GimbalPitchMode          string             `json:"gimbal_pitch_mode,omitempty" validate:"omitempty,oneof=usePointSetting manual free"`
	GlobalTransitionalSpeed  float64            `json:"global_transitional_speed,omitempty" validate:"omitempty,min=1,max=15"`
	TakeOffRefPointLatitude  float64            `json:"take_off_ref_point_latitude,omitempty" validate:"omitempty,min=-90,max=90"`
	TakeOffRefPointLongitude float64            `json:"take_off_ref_point_longitude,omitempty" validate:"omitempty,min=-180,max=180"`
	TakeOffRefPointHeight    float64            `json:"take_off_ref_point_height,omitempty"`
	TakeOffRefPointAGLHeight *float64           `json:"take_off_ref_point_agl_height,omitempty"`
	GlobalWaypointTurnMode   string             `json:"global_waypoint_turn_mode,omitempty" validate:"omitempty,oneof=coordinateTurn toPointAndStopWithDiscontinuityCurvature toPointAndStopWithContinuityCurvature toPointAndPassWithContinuityCurvature"`
	GlobalUseStraightLine    *bool              `json:"global_use_straight_line,omitempty"`
	GlobalTurnDampingDist    float64            `json:"global_turn_damping_dist,omitempty" validate:"omitempty,min=0"`
	Waypoints                []WaylinesWaypoint `json:"waypoints" validate:"required,min=1,dive"`
	SafeHeight               float64            `json:"safe_height,omitempty" validate:"omitempty,min=20,max=500"`
	Distance                 *float64           `json:"distance,omitempty" validate:"omitempty,min=0"`
	Duration                 *float64           `json:"duration,omitempty" validate:"omitempty,min=0"`
}

// WaylinesWaypoint represents a single waypoint in the waylines conversion schema.
type WaylinesWaypoint struct {
	Latitude         float64         `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude        float64         `json:"longitude" validate:"required,min=-180,max=180"`
	Height           float64         `json:"height" validate:"required,min=5,max=500"`
	Speed            float64         `json:"speed,omitempty" validate:"min=1,max=15"`
	TriggerType      string          `json:"trigger_type,omitempty" validate:"oneof=reachPoint passPoint manual betweenAdjacentPoints multipleTiming multipleDistance"`
	TriggerParam     float64         `json:"trigger_param,omitempty" validate:"min=0"`
	WaypointTurnMode string          `json:"waypoint_turn_mode,omitempty" validate:"omitempty,oneof=coordinateTurn toPointAndStopWithDiscontinuityCurvature toPointAndStopWithContinuityCurvature toPointAndPassWithContinuityCurvature"`
	UseStraightLine  *bool           `json:"use_straight_line,omitempty"`
	TurnDampingDist  float64         `json:"turn_damping_dist,omitempty" validate:"min=0"`
	Actions          []ActionRequest `json:"actions,omitempty" validate:"dive"`
}

// Validate validates the Waylines schema using the WPML validator.
func (w *Waylines) Validate() error {
	validator, err := NewValidator()
	if err != nil {
		return err
	}
	return validator.ValidateStruct(w)
}

// ApplyDefaults applies default values to the Waylines schema fields that are not set.
func (w *Waylines) ApplyDefaults() {
	if w.HeightType == "" {
		w.HeightType = HeightModeRelativeToStartPoint
	}
}
