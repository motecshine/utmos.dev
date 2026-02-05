package common

import (
	"encoding/json"
	"strconv"
)

// FlexInt64 is a flexible int64 that can unmarshal from both string and number JSON values
// DJI API sometimes returns timestamps as strings (e.g., "1704067200") instead of numbers
// It also handles empty strings by returning 0
type FlexInt64 int64

// UnmarshalJSON implements json.Unmarshaler
func (f *FlexInt64) UnmarshalJSON(data []byte) error {
	// First try to unmarshal as a number
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		*f = FlexInt64(num)
		return nil
	}

	// Then try to unmarshal as a string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Handle empty string
	if str == "" {
		*f = 0
		return nil
	}

	// Parse the string as int64
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*f = FlexInt64(num)
	return nil
}

// MarshalJSON implements json.Marshaler
func (f FlexInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(f))
}

// Int64 returns the underlying int64 value
func (f FlexInt64) Int64() int64 {
	return int64(f)
}

// FlexInt is a flexible int that can unmarshal from both string and number JSON values
// Similar to FlexInt64 but for int type
type FlexInt int

// UnmarshalJSON implements json.Unmarshaler
func (f *FlexInt) UnmarshalJSON(data []byte) error {
	// First try to unmarshal as a number
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		*f = FlexInt(num)
		return nil
	}

	// Then try to unmarshal as a string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Handle empty string
	if str == "" {
		*f = 0
		return nil
	}

	// Parse the string as int
	num64, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return err
	}
	*f = FlexInt(num64)
	return nil
}

// MarshalJSON implements json.Marshaler
func (f FlexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(f))
}

// Int returns the underlying int value
func (f FlexInt) Int() int {
	return int(f)
}

// FileInfo represents file information for file upload/download
type FileInfo struct {
	Name     string `json:"name"`     // File name
	URL      string `json:"url"`      // File URL
	Checksum string `json:"checksum"` // File checksum (MD5 or SHA256)
	Size     int    `json:"size"`     // File size in bytes
}

// Position represents GPS coordinate
type Position struct {
	Latitude  float64 `json:"latitude"`  // Latitude (-90 to 90)
	Longitude float64 `json:"longitude"` // Longitude (-180 to 180)
	Height    float64 `json:"height"`    // Height above sea level (meters)
}

// DeviceType represents device type enumeration
type DeviceType int

const (
	DeviceTypeAircraft DeviceType = 0 // Aircraft
	DeviceTypeDock     DeviceType = 3 // Dock
	DeviceTypeRC       DeviceType = 4 // Remote Controller
)

// StorageType represents storage location
type StorageType int

const (
	StorageTypeAircraft StorageType = 0 // Aircraft storage
	StorageTypePayload  StorageType = 1 // Payload storage
)

// PlannedPathPoint represents a planned trajectory point
type PlannedPathPoint struct {
	Latitude  float64 `json:"latitude"`  // Trajectory point latitude (-90 to 90)
	Longitude float64 `json:"longitude"` // Trajectory point longitude (-180 to 180)
	Height    float64 `json:"height"`    // Trajectory point height (m, ellipsoid height)
}
