package common

import (
	"encoding/json"
	"strconv"
)

// parseFlexJSON is a generic helper that parses JSON data as either a number or a string.
// It handles both numeric and string representations, as well as empty strings (returning 0).
func parseFlexJSON[T ~int64 | ~int](data []byte, bitSize int) (T, error) {
	var num T
	if err := json.Unmarshal(data, &num); err == nil {
		return num, nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return 0, err
	}

	if str == "" {
		return 0, nil
	}

	n, err := strconv.ParseInt(str, 10, bitSize)
	if err != nil {
		return 0, err
	}
	return T(n), nil
}

// FlexInt64 is a flexible int64 that can unmarshal from both string and number JSON values
// DJI API sometimes returns timestamps as strings (e.g., "1704067200") instead of numbers
// It also handles empty strings by returning 0
type FlexInt64 int64

// UnmarshalJSON implements json.Unmarshaler
func (f *FlexInt64) UnmarshalJSON(data []byte) error {
	v, err := parseFlexJSON[int64](data, 64)
	if err != nil {
		return err
	}
	*f = FlexInt64(v)
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
	v, err := parseFlexJSON[int](data, 32)
	if err != nil {
		return err
	}
	*f = FlexInt(v)
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
	// DeviceTypeAircraft is the aircraft device type.
	DeviceTypeAircraft DeviceType = 0
	// DeviceTypeDock is the dock device type.
	DeviceTypeDock DeviceType = 3
	// DeviceTypeRC is the remote controller device type.
	DeviceTypeRC DeviceType = 4
)

// StorageType represents storage location
type StorageType int

const (
	// StorageTypeAircraft is the aircraft storage type.
	StorageTypeAircraft StorageType = 0
	// StorageTypePayload is the payload storage type.
	StorageTypePayload StorageType = 1
)

// PlannedPathPoint represents a planned trajectory point
type PlannedPathPoint struct {
	Latitude  float64 `json:"latitude"`  // Trajectory point latitude (-90 to 90)
	Longitude float64 `json:"longitude"` // Trajectory point longitude (-180 to 180)
	Height    float64 `json:"height"`    // Trajectory point height (m, ellipsoid height)
}
