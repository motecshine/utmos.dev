package downlink

// ServiceMethods contains common DJI service method names
var ServiceMethods = struct {
	// Flight control
	Takeoff        string
	Land           string
	ReturnHome     string
	FlyToPoint     string
	StopFlyToPoint string

	// Gimbal control
	GimbalRotate string
	GimbalReset  string

	// Camera control
	CameraPhotoTake  string
	CameraVideoStart string
	CameraVideoStop  string
	CameraZoom       string

	// Dock control
	DockDebugModeOpen  string
	DockDebugModeClose string
	DockReboot         string

	// Mission control
	FlightTaskPrepare string
	FlightTaskExecute string
	FlightTaskCancel  string
}{
	Takeoff:            "takeoff",
	Land:               "land",
	ReturnHome:         "return_home",
	FlyToPoint:         "fly_to_point",
	StopFlyToPoint:     "stop_fly_to_point",
	GimbalRotate:       "gimbal_rotate",
	GimbalReset:        "gimbal_reset",
	CameraPhotoTake:    "camera_photo_take",
	CameraVideoStart:   "camera_video_start",
	CameraVideoStop:    "camera_video_stop",
	CameraZoom:         "camera_zoom",
	DockDebugModeOpen:  "dock_debug_mode_open",
	DockDebugModeClose: "dock_debug_mode_close",
	DockReboot:         "dock_reboot",
	FlightTaskPrepare:  "flight_task_prepare",
	FlightTaskExecute:  "flight_task_execute",
	FlightTaskCancel:   "flight_task_cancel",
}

// NewTakeoffCall creates a takeoff service call
func NewTakeoffCall(deviceSN string, height float64) *ServiceCall {
	return NewServiceCall(deviceSN, ServiceMethods.Takeoff, map[string]any{
		"height": height,
	})
}

// NewLandCall creates a land service call
func NewLandCall(deviceSN string) *ServiceCall {
	return NewServiceCall(deviceSN, ServiceMethods.Land, nil)
}

// NewReturnHomeCall creates a return home service call
func NewReturnHomeCall(deviceSN string) *ServiceCall {
	return NewServiceCall(deviceSN, ServiceMethods.ReturnHome, nil)
}

// NewFlyToPointCall creates a fly to point service call
func NewFlyToPointCall(deviceSN string, latitude, longitude, altitude, speed float64) *ServiceCall {
	return NewServiceCall(deviceSN, ServiceMethods.FlyToPoint, map[string]any{
		"latitude":  latitude,
		"longitude": longitude,
		"altitude":  altitude,
		"speed":     speed,
	})
}

// NewGimbalRotateCall creates a gimbal rotate service call
func NewGimbalRotateCall(deviceSN string, pitch, yaw float64) *ServiceCall {
	return NewServiceCall(deviceSN, ServiceMethods.GimbalRotate, map[string]any{
		"pitch": pitch,
		"yaw":   yaw,
	})
}

// NewCameraPhotoCall creates a camera photo take service call
func NewCameraPhotoCall(deviceSN string) *ServiceCall {
	return NewServiceCall(deviceSN, ServiceMethods.CameraPhotoTake, nil)
}
