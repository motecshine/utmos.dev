package wpml

type DroneModel int

// different typed constants cannot be deduplicated
const (
	DroneM300RTK   DroneModel = 60
	DroneM30       DroneModel = 67
	DroneM3Series  DroneModel = 77
	DroneM350RTK   DroneModel = 89
	DroneM3DSeries DroneModel = 91
	DroneM4Series  DroneModel = 99
	DroneM4DSeries DroneModel = 100
	DroneM400      DroneModel = 103
	DroneDahua     DroneModel = 1006
)

type DroneSubModel int

const (
	DroneSubModel0 DroneSubModel = 0
	DroneSubModel1 DroneSubModel = 1
	DroneSubModel2 DroneSubModel = 2
	DroneSubModel3 DroneSubModel = 3

	DroneSubM30Dual   DroneSubModel = 0
	DroneSubM30Triple DroneSubModel = 1

	DroneSubM3E DroneSubModel = 0
	DroneSubM3T DroneSubModel = 1
	DroneSubM3A DroneSubModel = 3

	DroneSubM3D  DroneSubModel = 0
	DroneSubM3TD DroneSubModel = 1

	DroneSubM4E DroneSubModel = 0
	DroneSubM4T DroneSubModel = 1

	DroneSubM4D  DroneSubModel = 0
	DroneSubM4TD DroneSubModel = 1

	DroneSUBX900 DroneSubModel = 0
)

type PayloadModel int

const (
	PayloadZ30  PayloadModel = 20
	PayloadXT2  PayloadModel = 26
	PayloadXTS  PayloadModel = 41
	PayloadH20  PayloadModel = 42
	PayloadH20T PayloadModel = 43
	PayloadH20N PayloadModel = 61
	PayloadH30  PayloadModel = 82
	PayloadH30T PayloadModel = 83

	PayloadM30Camera  PayloadModel = 52
	PayloadM30TCamera PayloadModel = 53

	PayloadMavic3ECamera PayloadModel = 66
	PayloadMavic3TCamera PayloadModel = 67
	PayloadMavic3ACamera PayloadModel = 129

	PayloadMatrice3DCamera  PayloadModel = 80
	PayloadMatrice3TDCamera PayloadModel = 81

	PayloadMatrice4ECamera PayloadModel = 88
	PayloadMatrice4TCamera PayloadModel = 89

	PayloadMatrice4DCamera  PayloadModel = 98
	PayloadMatrice4TDCamera PayloadModel = 99

	PayloadFPVCamera PayloadModel = 39

	PayloadDockCamera PayloadModel = 165

	PayloadAuxiliaryCamera PayloadModel = 176

	PayloadPSDK PayloadModel = 65534

	PayloadMatrice3D  = PayloadMatrice3DCamera
	PayloadMatrice3TD = PayloadMatrice3TDCamera

	PayloadX900 PayloadModel = 2014
)

type FlightMode string

const (
	FlightModeSafely       FlightMode = "safely"
	FlightModePointToPoint FlightMode = "pointToPoint"
)

type FinishAction string

const (
	FinishActionGoHome            FinishAction = "goHome"
	FinishActionNoAction          FinishAction = "noAction"
	FinishActionAutoLand          FinishAction = "autoLand"
	FinishActionGotoFirstWaypoint FinishAction = "gotoFirstWaypoint"
)

type RCLostAction string

const (
	RCLostActionGoContinue        RCLostAction = "goContinue"
	RCLostActionExecuteLostAction RCLostAction = "executeLostAction"
)

type ExecuteRCLostAction string

const (
	ExecuteRCLostActionGoBack  ExecuteRCLostAction = "goBack"
	ExecuteRCLostActionLanding ExecuteRCLostAction = "landing"
	ExecuteRCLostActionHover   ExecuteRCLostAction = "hover"
)

type HeightMode string

const (
	HeightModeEGM96                 HeightMode = "EGM96"
	HeightModeRelativeToStartPoint  HeightMode = "relativeToStartPoint"
	HeightModeAboveGroundLevel      HeightMode = "aboveGroundLevel"
	HeightModeRealTimeFollowSurface HeightMode = "realTimeFollowSurface"
)

type ExecuteHeightMode string

const (
	ExecuteHeightModeWGS84                 ExecuteHeightMode = "WGS84"
	ExecuteHeightModeRelativeToStartPoint  ExecuteHeightMode = "relativeToStartPoint"
	ExecuteHeightModeRealTimeFollowSurface ExecuteHeightMode = "realTimeFollowSurface"
)

type CoordinateMode string

const (
	CoordinateModeWGS84 CoordinateMode = "WGS84"
)

type PositioningType string

const (
	PositioningTypeGPS     PositioningType = "GPS"
	PositioningTypeRTKBase PositioningType = "RTKBaseStation"
	PositioningTypeQianXun PositioningType = "QianXun"
	PositioningTypeCustom  PositioningType = "Custom"
)

type TemplateType string

const (
	TemplateTypeWaypoint     TemplateType = "waypoint"
	TemplateTypeMapping2D    TemplateType = "mapping2d"
	TemplateTypeMapping3D    TemplateType = "mapping3d"
	TemplateTypeMappingStrip TemplateType = "mappingStrip"
)

// String representations for JSON/testing
//
// different typed constants cannot be deduplicated
const (
	DroneModelMatrice3TD string = "M3TD"
	DroneModelMatrice3E  string = "M3E"
	DroneModelMatrice3T  string = "M3T"
	DroneModelMavic3     string = "Mavic3"
	DroneModelMini3      string = "Mini3"
	DroneModelAir2S      string = "Air2S"
	DroneModelM30        string = "M30"
	DroneModelM300RTK    string = "M300RTK"
	DroneModelM350RTK    string = "M350RTK"
)

const (
	PayloadModelM3TD   string = "M3TD"
	PayloadModelM3E    string = "M3E"
	PayloadModelM3T    string = "M3T"
	PayloadModelMavic3 string = "Mavic3"
	PayloadModelMini3  string = "Mini3"
	PayloadModelAir2S  string = "Air2S"
	PayloadModelM30    string = "M30"
)
