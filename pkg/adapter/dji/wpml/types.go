package wpml

// DroneModel represents a DJI drone model identifier.
type DroneModel int

// Supported drone model constants.
const (
	// DroneM300RTK is the DJI Matrice 300 RTK drone model.
	DroneM300RTK DroneModel = 60
	// DroneM30 is the DJI Matrice 30 drone model.
	DroneM30 DroneModel = 67
	// DroneM3Series is the DJI Mavic 3 series drone model.
	DroneM3Series DroneModel = 77
	// DroneM350RTK is the DJI Matrice 350 RTK drone model.
	DroneM350RTK DroneModel = 89
	// DroneM3DSeries is the DJI Matrice 3D series drone model.
	DroneM3DSeries DroneModel = 91
	// DroneM4Series is the DJI Matrice 4 series drone model.
	DroneM4Series DroneModel = 99
	// DroneM4DSeries is the DJI Matrice 4D series drone model.
	DroneM4DSeries DroneModel = 100
	// DroneM400 is the DJI Matrice 400 drone model.
	DroneM400 DroneModel = 103
	// DroneDahua is the Dahua drone model.
	DroneDahua DroneModel = 1006
)

// DroneSubModel represents a DJI drone sub-model variant.
type DroneSubModel int

// Drone sub-model constants for differentiating variants within a model series.
const (
	// DroneSubModel0 is the default drone sub-model variant (index 0).
	DroneSubModel0 DroneSubModel = 0
	// DroneSubModel1 is drone sub-model variant index 1.
	DroneSubModel1 DroneSubModel = 1
	// DroneSubModel2 is drone sub-model variant index 2.
	DroneSubModel2 DroneSubModel = 2
	// DroneSubModel3 is drone sub-model variant index 3.
	DroneSubModel3 DroneSubModel = 3

	// DroneSubM30Dual is the M30 dual-sensor sub-model.
	DroneSubM30Dual DroneSubModel = 0
	// DroneSubM30Triple is the M30 triple-sensor sub-model.
	DroneSubM30Triple DroneSubModel = 1

	// DroneSubM3E is the Mavic 3 Enterprise sub-model.
	DroneSubM3E DroneSubModel = 0
	// DroneSubM3T is the Mavic 3 Thermal sub-model.
	DroneSubM3T DroneSubModel = 1
	// DroneSubM3A is the Mavic 3 Advanced sub-model.
	DroneSubM3A DroneSubModel = 3

	// DroneSubM3D is the Matrice 3D sub-model.
	DroneSubM3D DroneSubModel = 0
	// DroneSubM3TD is the Matrice 3TD (thermal-dual) sub-model.
	DroneSubM3TD DroneSubModel = 1

	// DroneSubM4E is the Matrice 4 Enterprise sub-model.
	DroneSubM4E DroneSubModel = 0
	// DroneSubM4T is the Matrice 4 Thermal sub-model.
	DroneSubM4T DroneSubModel = 1

	// DroneSubM4D is the Matrice 4D sub-model.
	DroneSubM4D DroneSubModel = 0
	// DroneSubM4TD is the Matrice 4TD (thermal-dual) sub-model.
	DroneSubM4TD DroneSubModel = 1

	// DroneSUBX900 is the X900 drone sub-model.
	DroneSUBX900 DroneSubModel = 0
)

// PayloadModel represents a DJI payload (camera/sensor) model identifier.
type PayloadModel int

// Supported payload model constants.
const (
	// PayloadZ30 is the Zenmuse Z30 payload.
	PayloadZ30 PayloadModel = 20
	// PayloadXT2 is the Zenmuse XT2 thermal payload.
	PayloadXT2 PayloadModel = 26
	// PayloadXTS is the Zenmuse XTS thermal payload.
	PayloadXTS PayloadModel = 41
	// PayloadH20 is the Zenmuse H20 payload.
	PayloadH20 PayloadModel = 42
	// PayloadH20T is the Zenmuse H20T thermal payload.
	PayloadH20T PayloadModel = 43
	// PayloadH20N is the Zenmuse H20N night-vision payload.
	PayloadH20N PayloadModel = 61
	// PayloadH30 is the Zenmuse H30 payload.
	PayloadH30 PayloadModel = 82
	// PayloadH30T is the Zenmuse H30T thermal payload.
	PayloadH30T PayloadModel = 83

	// PayloadM30Camera is the Matrice 30 built-in camera payload.
	PayloadM30Camera PayloadModel = 52
	// PayloadM30TCamera is the Matrice 30T built-in thermal camera payload.
	PayloadM30TCamera PayloadModel = 53

	// PayloadMavic3ECamera is the Mavic 3 Enterprise camera payload.
	PayloadMavic3ECamera PayloadModel = 66
	// PayloadMavic3TCamera is the Mavic 3 Thermal camera payload.
	PayloadMavic3TCamera PayloadModel = 67
	// PayloadMavic3ACamera is the Mavic 3 Advanced camera payload.
	PayloadMavic3ACamera PayloadModel = 129

	// PayloadMatrice3DCamera is the Matrice 3D camera payload.
	PayloadMatrice3DCamera PayloadModel = 80
	// PayloadMatrice3TDCamera is the Matrice 3TD thermal-dual camera payload.
	PayloadMatrice3TDCamera PayloadModel = 81

	// PayloadMatrice4ECamera is the Matrice 4 Enterprise camera payload.
	PayloadMatrice4ECamera PayloadModel = 88
	// PayloadMatrice4TCamera is the Matrice 4 Thermal camera payload.
	PayloadMatrice4TCamera PayloadModel = 89

	// PayloadMatrice4DCamera is the Matrice 4D camera payload.
	PayloadMatrice4DCamera PayloadModel = 98
	// PayloadMatrice4TDCamera is the Matrice 4TD thermal-dual camera payload.
	PayloadMatrice4TDCamera PayloadModel = 99

	// PayloadFPVCamera is the FPV camera payload.
	PayloadFPVCamera PayloadModel = 39

	// PayloadDockCamera is the Dock station camera payload.
	PayloadDockCamera PayloadModel = 165

	// PayloadAuxiliaryCamera is the auxiliary camera payload.
	PayloadAuxiliaryCamera PayloadModel = 176

	// PayloadPSDK is the Payload SDK (third-party) payload.
	PayloadPSDK PayloadModel = 65534

	// PayloadMatrice3D is an alias for PayloadMatrice3DCamera.
	PayloadMatrice3D = PayloadMatrice3DCamera
	// PayloadMatrice3TD is an alias for PayloadMatrice3TDCamera.
	PayloadMatrice3TD = PayloadMatrice3TDCamera

	// PayloadX900 is the X900 payload.
	PayloadX900 PayloadModel = 2014
)

// FlightMode represents the flight mode for flying to a wayline.
type FlightMode string

// Supported flight mode constants.
const (
	// FlightModeSafely is the safe flight mode that avoids obstacles.
	FlightModeSafely FlightMode = "safely"
	// FlightModePointToPoint is the direct point-to-point flight mode.
	FlightModePointToPoint FlightMode = "pointToPoint"
)

// FinishAction represents the action to perform when the mission finishes.
type FinishAction string

// Supported finish action constants.
const (
	// FinishActionGoHome returns the drone to its home point after mission completion.
	FinishActionGoHome FinishAction = "goHome"
	// FinishActionNoAction keeps the drone hovering at the last waypoint.
	FinishActionNoAction FinishAction = "noAction"
	// FinishActionAutoLand lands the drone at the last waypoint.
	FinishActionAutoLand FinishAction = "autoLand"
	// FinishActionGotoFirstWaypoint flies the drone back to the first waypoint.
	FinishActionGotoFirstWaypoint FinishAction = "gotoFirstWaypoint"
)

// RCLostAction represents the action to perform when the remote controller signal is lost.
type RCLostAction string

// Supported RC lost action constants.
const (
	// RCLostActionGoContinue continues the mission when RC signal is lost.
	RCLostActionGoContinue RCLostAction = "goContinue"
	// RCLostActionExecuteLostAction executes the configured lost action when RC signal is lost.
	RCLostActionExecuteLostAction RCLostAction = "executeLostAction"
)

// ExecuteRCLostAction represents the specific action to execute when the RC signal is lost.
type ExecuteRCLostAction string

// Supported execute RC lost action constants.
const (
	// ExecuteRCLostActionGoBack returns the drone to the home point.
	ExecuteRCLostActionGoBack ExecuteRCLostAction = "goBack"
	// ExecuteRCLostActionLanding lands the drone at its current position.
	ExecuteRCLostActionLanding ExecuteRCLostAction = "landing"
	// ExecuteRCLostActionHover hovers the drone at its current position.
	ExecuteRCLostActionHover ExecuteRCLostAction = "hover"
)

// HeightMode represents the reference frame for waypoint heights.
type HeightMode string

// Supported height mode constants.
const (
	// HeightModeEGM96 uses the EGM96 geoid model for height reference.
	HeightModeEGM96 HeightMode = "EGM96"
	// HeightModeRelativeToStartPoint uses height relative to the takeoff point.
	HeightModeRelativeToStartPoint HeightMode = "relativeToStartPoint"
	// HeightModeAboveGroundLevel uses height above the ground level.
	HeightModeAboveGroundLevel HeightMode = "aboveGroundLevel"
	// HeightModeRealTimeFollowSurface follows the terrain surface in real time.
	HeightModeRealTimeFollowSurface HeightMode = "realTimeFollowSurface"
)

// ExecuteHeightMode represents the height mode used during wayline execution.
type ExecuteHeightMode string

// Supported execute height mode constants.
const (
	// ExecuteHeightModeWGS84 uses the WGS84 ellipsoidal height for execution.
	ExecuteHeightModeWGS84 ExecuteHeightMode = "WGS84"
	// ExecuteHeightModeRelativeToStartPoint uses height relative to the start point for execution.
	ExecuteHeightModeRelativeToStartPoint ExecuteHeightMode = "relativeToStartPoint"
	// ExecuteHeightModeRealTimeFollowSurface follows the terrain surface in real time during execution.
	ExecuteHeightModeRealTimeFollowSurface ExecuteHeightMode = "realTimeFollowSurface"
)

// CoordinateMode represents the coordinate reference system used for waypoints.
type CoordinateMode string

// Supported coordinate mode constants.
const (
	// CoordinateModeWGS84 is the WGS84 coordinate reference system.
	CoordinateModeWGS84 CoordinateMode = "WGS84"
)

// PositioningType represents the positioning method used for waypoints.
type PositioningType string

// Supported positioning type constants.
const (
	// PositioningTypeGPS uses standard GPS positioning.
	PositioningTypeGPS PositioningType = "GPS"
	// PositioningTypeRTKBase uses RTK base station positioning for higher accuracy.
	PositioningTypeRTKBase PositioningType = "RTKBaseStation"
	// PositioningTypeQianXun uses QianXun network RTK positioning.
	PositioningTypeQianXun PositioningType = "QianXun"
	// PositioningTypeCustom uses a custom positioning source.
	PositioningTypeCustom PositioningType = "Custom"
)

// TemplateType represents the type of WPML mission template.
type TemplateType string

// Supported template type constants.
const (
	// TemplateTypeWaypoint is a waypoint-based mission template.
	TemplateTypeWaypoint TemplateType = "waypoint"
	// TemplateTypeMapping2D is a 2D mapping mission template.
	TemplateTypeMapping2D TemplateType = "mapping2d"
	// TemplateTypeMapping3D is a 3D mapping mission template.
	TemplateTypeMapping3D TemplateType = "mapping3d"
	// TemplateTypeMappingStrip is a strip mapping mission template.
	TemplateTypeMappingStrip TemplateType = "mappingStrip"
)

// String representations for JSON/testing
//
// different typed constants cannot be deduplicated

// Drone model string constants for JSON serialization and testing.
const (
	// DroneModelMatrice3TD is the string identifier for the Matrice 3TD drone.
	DroneModelMatrice3TD string = "M3TD"
	// DroneModelMatrice3E is the string identifier for the Mavic 3 Enterprise drone.
	DroneModelMatrice3E string = "M3E"
	// DroneModelMatrice3T is the string identifier for the Mavic 3 Thermal drone.
	DroneModelMatrice3T string = "M3T"
	// DroneModelMavic3 is the string identifier for the Mavic 3 drone.
	DroneModelMavic3 string = "Mavic3"
	// DroneModelMini3 is the string identifier for the Mini 3 drone.
	DroneModelMini3 string = "Mini3"
	// DroneModelAir2S is the string identifier for the Air 2S drone.
	DroneModelAir2S string = "Air2S"
	// DroneModelM30 is the string identifier for the Matrice 30 drone.
	DroneModelM30 string = "M30"
	// DroneModelM300RTK is the string identifier for the Matrice 300 RTK drone.
	DroneModelM300RTK string = "M300RTK"
	// DroneModelM350RTK is the string identifier for the Matrice 350 RTK drone.
	DroneModelM350RTK string = "M350RTK"
)

// Payload model string constants for JSON serialization and testing.
const (
	// PayloadModelM3TD is the string identifier for the M3TD payload.
	PayloadModelM3TD string = "M3TD"
	// PayloadModelM3E is the string identifier for the M3E payload.
	PayloadModelM3E string = "M3E"
	// PayloadModelM3T is the string identifier for the M3T payload.
	PayloadModelM3T string = "M3T"
	// PayloadModelMavic3 is the string identifier for the Mavic 3 payload.
	PayloadModelMavic3 string = "Mavic3"
	// PayloadModelMini3 is the string identifier for the Mini 3 payload.
	PayloadModelMini3 string = "Mini3"
	// PayloadModelAir2S is the string identifier for the Air 2S payload.
	PayloadModelAir2S string = "Air2S"
	// PayloadModelM30 is the string identifier for the M30 payload.
	PayloadModelM30 string = "M30"
)
