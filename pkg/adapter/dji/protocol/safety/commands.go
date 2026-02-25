package safety

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Flight Safety Commands
// ===============================

// FlightAreaFile represents the flight area file information
type FlightAreaFile struct {
	Name     string `json:"name"`     // File name
	URL      string `json:"url"`      // File URL
	Checksum string `json:"checksum"` // File SHA256 signature
	Size     int    `json:"size"`     // File size in bytes
}

// FlightAreasUpdateCommand represents the flight areas update request (data is null)
type FlightAreasUpdateCommand struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  any `json:"data"`
}

// NewFlightAreasUpdateCommand creates a new flight areas update request
func NewFlightAreasUpdateCommand() *FlightAreasUpdateCommand {
	return &FlightAreasUpdateCommand{
		Header:     common.NewHeader(),
		MethodName: "flight_areas_update",
		DataValue:  nil,
	}
}

// Method returns the method name.
func (c *FlightAreasUpdateCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *FlightAreasUpdateCommand) Data() any { return c.DataValue }

// FlightAreasGetData represents the flight areas get data
type FlightAreasGetData struct {
	AreaType *string `json:"area_type,omitempty"` // Area type filter (optional)
}

// FlightAreasGetRequest represents the flight areas get request
type FlightAreasGetRequest struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  FlightAreasGetData `json:"data"`
}

// NewFlightAreasGetRequest creates a new flight areas get request
func NewFlightAreasGetRequest(data FlightAreasGetData) *FlightAreasGetRequest {
	return &FlightAreasGetRequest{
		Header:     common.NewHeader(),
		MethodName: "flight_areas_get",
		DataValue:  data,
	}
}

// Method returns the method name.
func (r *FlightAreasGetRequest) Method() string { return r.MethodName }

// Data returns the command/event data.
func (r *FlightAreasGetRequest) Data() any { return r.DataValue }

// UnlockLicenseSwitchData represents the unlock license switch data
type UnlockLicenseSwitchData struct {
	LicenseID int  `json:"license_id"` // License ID (unique identifier)
	Enable    bool `json:"enable"`     // Enable/disable unlock license
}

// UnlockLicenseSwitchCommand represents the unlock license switch request
type UnlockLicenseSwitchCommand struct {
	common.Header
	MethodName string                  `json:"method"`
	DataValue  UnlockLicenseSwitchData `json:"data"`
}

// NewUnlockLicenseSwitchCommand creates a new unlock license switch request
func NewUnlockLicenseSwitchCommand(data UnlockLicenseSwitchData) *UnlockLicenseSwitchCommand {
	return &UnlockLicenseSwitchCommand{
		Header:     common.NewHeader(),
		MethodName: "unlock_license_switch",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *UnlockLicenseSwitchCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *UnlockLicenseSwitchCommand) Data() any { return c.DataValue }

// LicenseFileInfo represents the license file information
type LicenseFileInfo struct {
	URL         string `json:"url"`         // File URL
	Fingerprint string `json:"fingerprint"` // File MD5 signature
}

// UnlockLicenseUpdateData represents the unlock license update data
type UnlockLicenseUpdateData struct {
	File *LicenseFileInfo `json:"file,omitempty"` // Offline license file (optional, use Flysafe server if omitted)
}

// UnlockLicenseUpdateCommand represents the unlock license update request
type UnlockLicenseUpdateCommand struct {
	common.Header
	MethodName string                  `json:"method"`
	DataValue  UnlockLicenseUpdateData `json:"data"`
}

// NewUnlockLicenseUpdateCommand creates a new unlock license update request
func NewUnlockLicenseUpdateCommand(data UnlockLicenseUpdateData) *UnlockLicenseUpdateCommand {
	return &UnlockLicenseUpdateCommand{
		Header:     common.NewHeader(),
		MethodName: "unlock_license_update",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *UnlockLicenseUpdateCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *UnlockLicenseUpdateCommand) Data() any { return c.DataValue }

// LicenseCommonFields represents the common fields of unlock license
type LicenseCommonFields struct {
	LicenseID int    `json:"license_id"` // License unique identifier
	Name      string `json:"name"`       // License name
	Type      int    `json:"type"`       // License type (0-6)
	GroupID   int    `json:"group_id"`   // License group unique identifier
	UserID    string `json:"user_id"`    // License owner user account
	DeviceSN  string `json:"device_sn"`  // Device serial number bound to license
	BeginTime int64  `json:"begin_time"` // Valid start timestamp (seconds)
	EndTime   int64  `json:"end_time"`   // Valid end timestamp (seconds)
	UserOnly  bool   `json:"user_only"`  // Whether to verify user account
	Enabled   bool   `json:"enabled"`    // Whether enabled
}

// AreaUnlockInfo represents the area unlock license information (type=0)
type AreaUnlockInfo struct {
	AreaIDs []int `json:"area_ids"` // Area ID collection
}

// CircleUnlockInfo represents the circle unlock license information (type=1)
type CircleUnlockInfo struct {
	Radius    int     `json:"radius"`    // Circle radius (meters)
	Latitude  float64 `json:"latitude"`  // Circle center latitude (-90 to 90)
	Longitude float64 `json:"longitude"` // Circle center longitude (-180 to 180)
	Height    int     `json:"height"`    // Height limit (meters, 0-65535)
}

// CountryUnlockInfo represents the country/region unlock license information (type=2)
type CountryUnlockInfo struct {
	CountryNumber int `json:"country_number"` // Country/region digital code (ISO 3166-1)
	Height        int `json:"height"`         // Height limit (meters, 0-65535)
}

// HeightUnlockInfo represents the height unlock license information (type=3)
type HeightUnlockInfo struct {
	Height int `json:"height"` // Height limit (meters, 0-65535)
}

// PolygonPoint represents a polygon vertex
type PolygonPoint struct {
	Latitude  float64 `json:"latitude"`  // Latitude (-90 to 90)
	Longitude float64 `json:"longitude"` // Longitude (-180 to 180)
}

// PolygonUnlockInfo represents the polygon unlock license information (type=4)
type PolygonUnlockInfo struct {
	Points []PolygonPoint `json:"points"` // Polygon vertex GPS coordinate collection
}

// PowerUnlockInfo represents the power unlock license information (type=5)
type PowerUnlockInfo struct {
	// Empty struct
}

// RIDUnlockInfo represents the RID unlock license information (type=6)
type RIDUnlockInfo struct {
	Level int `json:"level"` // RID type (1=EU RID, 2=China RID)
}

// UnlockLicense represents the unlock license with all possible types
type UnlockLicense struct {
	CommonFields  LicenseCommonFields `json:"common_fields"`            // Common license information
	AreaUnlock    *AreaUnlockInfo     `json:"area_unlock,omitempty"`    // Area unlock (type=0)
	CircleUnlock  *CircleUnlockInfo   `json:"circle_unlock,omitempty"`  // Circle unlock (type=1)
	CountryUnlock *CountryUnlockInfo  `json:"country_unlock,omitempty"` // Country/region unlock (type=2)
	HeightUnlock  *HeightUnlockInfo   `json:"height_unlock,omitempty"`  // Height unlock (type=3)
	PolygonUnlock *PolygonUnlockInfo  `json:"polygon_unlock,omitempty"` // Polygon unlock (type=4)
	PowerUnlock   *PowerUnlockInfo    `json:"power_unlock,omitempty"`   // Power unlock (type=5)
	RIDUnlock     *RIDUnlockInfo      `json:"rid_unlock,omitempty"`     // RID unlock (type=6)
}

// UnlockLicenseListData represents the unlock license list request data
type UnlockLicenseListData struct {
	DeviceModelDomain int `json:"device_model_domain"` // License location (0=aircraft, 3=dock)
}

// UnlockLicenseListResponse represents the unlock license list response data
type UnlockLicenseListResponse struct {
	Result            int             `json:"result"`              // Result code (non-zero indicates error)
	DeviceModelDomain int             `json:"device_model_domain"` // License location
	Consistence       bool            `json:"consistence"`         // License consistency
	Licenses          []UnlockLicense `json:"licenses"`            // License list
}

// FlightAreasGetResponse represents the flight areas get response
type FlightAreasGetResponse struct {
	Files []FlightAreaFile `json:"files"` // Flight area file list
}

// UnlockLicenseListCommand represents the unlock license list request
type UnlockLicenseListCommand struct {
	common.Header
	MethodName string                `json:"method"`
	DataValue  UnlockLicenseListData `json:"data"`
}

// NewUnlockLicenseListCommand creates a new unlock license list request
func NewUnlockLicenseListCommand(data UnlockLicenseListData) *UnlockLicenseListCommand {
	return &UnlockLicenseListCommand{
		Header:     common.NewHeader(),
		MethodName: "unlock_license_list",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *UnlockLicenseListCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *UnlockLicenseListCommand) Data() any { return c.DataValue }

// GetHeader implements Command.GetHeader
func (r *FlightAreasGetRequest) GetHeader() *common.Header {
	return &r.Header
}

// GetHeader implements Command.GetHeader
func (c *FlightAreasUpdateCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *UnlockLicenseListCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *UnlockLicenseSwitchCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *UnlockLicenseUpdateCommand) GetHeader() *common.Header {
	return &c.Header
}
