package file

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// File Management Events
// ===============================

// HighestPriorityUploadFlighttaskMediaData represents the media file upload priority data
type HighestPriorityUploadFlighttaskMediaData struct {
	FlightID string `json:"flight_id"` // Flight task ID with highest upload priority
}

// HighestPriorityUploadFlighttaskMediaEvent represents the media file upload priority event
type HighestPriorityUploadFlighttaskMediaEvent struct {
	common.Header
	MethodName string                                   `json:"method"`
	DataValue  HighestPriorityUploadFlighttaskMediaData `json:"data"`
}

// Method returns the method name.
func (e *HighestPriorityUploadFlighttaskMediaEvent) Method() string { return e.MethodName }

// Data returns the command/event data.
func (e *HighestPriorityUploadFlighttaskMediaEvent) Data() any { return e.DataValue }

// GetHeader returns the event header.
func (e *HighestPriorityUploadFlighttaskMediaEvent) GetHeader() *common.Header {
	return &e.Header
}

// ShootPosition represents the shooting position
type ShootPosition struct {
	Lat float64 `json:"lat"` // Shooting position latitude
	Lng float64 `json:"lng"` // Shooting position longitude
}

// MediaMetadata represents the media file metadata
type MediaMetadata struct {
	GimbalYawDegree  float64       `json:"gimbal_yaw_degree"` // Gimbal yaw angle
	AbsoluteAltitude float64       `json:"absolute_altitude"` // Shooting absolute altitude
	RelativeAltitude float64       `json:"relative_altitude"` // Shooting relative altitude
	CreateTime       string        `json:"created_time"`      // Media shooting time (ISO8601 format)
	ShootPosition    ShootPosition `json:"shoot_position"`    // Shooting position
}

// FileExt represents the file extension information
type FileExt struct {
	FlightID        string `json:"flight_id"`         // Flight task ID
	DroneModelKey   string `json:"drone_model_key"`   // Drone product enum value
	PayloadModelKey string `json:"payload_model_key"` // Payload product enum value
	IsOriginal      bool   `json:"is_original"`       // Whether it is an original image
}

// FileInfo represents the file information
type FileInfo struct {
	ObjectKey      string        `json:"object_key"`        // File key in object storage bucket
	Path           string        `json:"path"`              // File business path
	Name           string        `json:"name"`              // File name
	Ext            FileExt       `json:"ext"`               // File extension information
	Metadata       MediaMetadata `json:"metadata"`          // Media metadata
	CloudToCloudID string        `json:"cloud_to_cloud_id"` // Cloud-to-cloud storage bucket ID
}

// FlightTaskInfo represents the flight task information
type FlightTaskInfo struct {
	UploadedFileCount int `json:"uploaded_file_count"` // Current uploaded media count for this flight
	ExpectedFileCount int `json:"expected_file_count"` // Total media count for this flight
	FlightType        int `json:"flight_type"`         // Flight type (0=wayline task, 1=one-key takeoff task)
}

// FileUploadCallbackData represents the file upload result data
type FileUploadCallbackData struct {
	File       FileInfo       `json:"file"`        // File information
	FlightTask FlightTaskInfo `json:"flight_task"` // Flight task information
}

// FileUploadCallbackEvent represents the file upload result event
type FileUploadCallbackEvent struct {
	common.Header
	MethodName string                 `json:"method"`
	DataValue  FileUploadCallbackData `json:"data"`
}

// Method returns the method name.
func (e *FileUploadCallbackEvent) Method() string { return e.MethodName }

// Data returns the command/event data.
func (e *FileUploadCallbackEvent) Data() any { return e.DataValue }

// GetHeader returns the event header.
func (e *FileUploadCallbackEvent) GetHeader() *common.Header { return &e.Header }
