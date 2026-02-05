package drc

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// DRC (Direct Remote Control) Events
// ===============================

// PlannedPathPoint represents a planned trajectory point
type PlannedPathPoint struct {
	Latitude  float64 `json:"latitude"`  // Trajectory point latitude (-90 to 90)
	Longitude float64 `json:"longitude"` // Trajectory point longitude (-180 to 180)
	Height    float64 `json:"height"`    // Trajectory point height (m, ellipsoid height)
}

// FlyToPointProgressData represents the fly to point progress data
type FlyToPointProgressData struct {
	Result            int                `json:"result"`              // Return code (0=success)
	FlyToID           string             `json:"fly_to_id"`           // Fly to target point ID
	Status            string             `json:"status"`              // Status (wayline_progress, wayline_ok, wayline_failed, wayline_cancel)
	WayPointIndex     int                `json:"way_point_index"`     // Current waypoint index
	RemainingDistance float64            `json:"remaining_distance"`  // Remaining task distance (m)
	RemainingTime     float64            `json:"remaining_time"`      // Remaining task time (s)
	PlannedPathPoints []PlannedPathPoint `json:"planned_path_points"` // Planned trajectory point list
}

// FlyToPointProgressEvent represents the fly to point progress event
type FlyToPointProgressEvent struct {
	common.Header
	MethodName string                 `json:"method"`
	DataValue  FlyToPointProgressData `json:"data"`
}

func (e *FlyToPointProgressEvent) Method() string            { return e.MethodName }
func (e *FlyToPointProgressEvent) Data() any                 { return e.DataValue }
func (e *FlyToPointProgressEvent) GetHeader() *common.Header { return &e.Header }

// TakeoffToPointProgressData represents the one-key takeoff progress data
type TakeoffToPointProgressData struct {
	Result            int                `json:"result"`              // Return code (0=success)
	Status            string             `json:"status"`              // Task status (task_ready, wayline_progress, wayline_ok, wayline_failed, wayline_cancel, task_finish)
	FlightID          string             `json:"flight_id"`           // One-key takeoff task UUID
	TrackID           string             `json:"track_id"`            // Track ID
	WayPointIndex     int                `json:"way_point_index"`     // Current waypoint index
	RemainingDistance float64            `json:"remaining_distance"`  // Remaining task distance (m)
	RemainingTime     float64            `json:"remaining_time"`      // Remaining task time (s)
	PlannedPathPoints []PlannedPathPoint `json:"planned_path_points"` // Planned trajectory point list
}

// TakeoffToPointProgressEvent represents the one-key takeoff progress event
type TakeoffToPointProgressEvent struct {
	common.Header
	MethodName string                     `json:"method"`
	DataValue  TakeoffToPointProgressData `json:"data"`
}

func (e *TakeoffToPointProgressEvent) Method() string            { return e.MethodName }
func (e *TakeoffToPointProgressEvent) Data() any                 { return e.DataValue }
func (e *TakeoffToPointProgressEvent) GetHeader() *common.Header { return &e.Header }

// DRCStatusNotifyData represents the DRC link status notification data
type DRCStatusNotifyData struct {
	Result   int `json:"result"`    // Return code (0=success)
	DRCState int `json:"drc_state"` // DRC state (0=disconnected, 1=connecting, 2=connected)
}

// DRCStatusNotifyEvent represents the DRC link status notification event
type DRCStatusNotifyEvent struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  DRCStatusNotifyData `json:"data"`
}

func (e *DRCStatusNotifyEvent) Method() string            { return e.MethodName }
func (e *DRCStatusNotifyEvent) Data() any                 { return e.DataValue }
func (e *DRCStatusNotifyEvent) GetHeader() *common.Header { return &e.Header }

// JoystickInvalidNotifyData represents the joystick control invalid reason data
type JoystickInvalidNotifyData struct {
	Reason int `json:"reason"` // Invalid reason (0=RC disconnected, 1=low battery RTH, 2=low battery landing, 3=near no-fly zone, 4=RC authority grab)
}

// JoystickInvalidNotifyEvent represents the joystick control invalid notification event
type JoystickInvalidNotifyEvent struct {
	common.Header
	MethodName string                    `json:"method"`
	DataValue  JoystickInvalidNotifyData `json:"data"`
}

func (e *JoystickInvalidNotifyEvent) Method() string            { return e.MethodName }
func (e *JoystickInvalidNotifyEvent) Data() any                 { return e.DataValue }
func (e *JoystickInvalidNotifyEvent) GetHeader() *common.Header { return &e.Header }

// POIStatusNotifyData represents the POI (Point of Interest) circling status data
type POIStatusNotifyData struct {
	Result       int     `json:"result"`        // Return code (0=success)
	CircleRadius float64 `json:"circle_radius"` // Circle radius (m)
	CircleSpeed  float64 `json:"circle_speed"`  // Circle speed (m/s)
	Status       int     `json:"status"`        // POI status (0=idle, 1=circling, 2=paused, 3=stopped)
	Reason       int     `json:"reason"`        // Status change reason
}

// POIStatusNotifyEvent represents the POI circling status notification event
type POIStatusNotifyEvent struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  POIStatusNotifyData `json:"data"`
}

func (e *POIStatusNotifyEvent) Method() string            { return e.MethodName }
func (e *POIStatusNotifyEvent) Data() any                 { return e.DataValue }
func (e *POIStatusNotifyEvent) GetHeader() *common.Header { return &e.Header }

// HSIInfoPushData represents the horizontal situation indicator (obstacle avoidance) data
type HSIInfoPushData struct {
	Result          int       `json:"result"`           // Return code (0=success)
	UpDistance      float64   `json:"up_distance"`      // Upward obstacle distance (m)
	DownDistance    float64   `json:"down_distance"`    // Downward obstacle distance (m)
	AroundDistances []float64 `json:"around_distances"` // 360-degree horizontal obstacle distances (m), array length 360
}

// HSIInfoPushEvent represents the obstacle avoidance information push event (high frequency in DRC mode)
type HSIInfoPushEvent struct {
	common.Header
	MethodName string          `json:"method"`
	DataValue  HSIInfoPushData `json:"data"`
}

func (e *HSIInfoPushEvent) Method() string            { return e.MethodName }
func (e *HSIInfoPushEvent) Data() any                 { return e.DataValue }
func (e *HSIInfoPushEvent) GetHeader() *common.Header { return &e.Header }

// CameraPhotoTakeProgressData represents the camera photo taking progress data
type CameraPhotoTakeProgressData struct {
	Result     int    `json:"result"`      // Return code (0=success)
	Status     string `json:"status"`      // Photo taking status (sent, in_progress, ok, failed)
	Percent    int    `json:"percent"`     // Progress percentage (0-100)
	CameraMode int    `json:"camera_mode"` // Camera mode (0=photo, 1=video, 2=playback)
}

// CameraPhotoTakeProgressEvent represents the camera photo taking progress event
type CameraPhotoTakeProgressEvent struct {
	common.Header
	MethodName string                      `json:"method"`
	DataValue  CameraPhotoTakeProgressData `json:"data"`
}

func (e *CameraPhotoTakeProgressEvent) Method() string            { return e.MethodName }
func (e *CameraPhotoTakeProgressEvent) Data() any                 { return e.DataValue }
func (e *CameraPhotoTakeProgressEvent) GetHeader() *common.Header { return &e.Header }
