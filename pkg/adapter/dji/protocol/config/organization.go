package config

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Organization Management Commands
// ===============================

// AirportBindStatusRequest represents the airport bind status request
type AirportBindStatusRequest struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  interface{} `json:"data"`
}

// NewAirportBindStatusRequest creates a new airport bind status request
func NewAirportBindStatusRequest() *AirportBindStatusRequest {
	return &AirportBindStatusRequest{
		Header:     common.NewHeader(),
		MethodName: "airport_bind_status",
		DataValue:  nil,
	}
}

func (r *AirportBindStatusRequest) Method() string { return r.MethodName }
func (r *AirportBindStatusRequest) Data() any      { return r.DataValue }

// AirportOrganizationGetRequest represents the airport organization get request
type AirportOrganizationGetRequest struct {
	common.Header
	MethodName string      `json:"method"`
	DataValue  interface{} `json:"data"`
}

// NewAirportOrganizationGetRequest creates a new airport organization get request
func NewAirportOrganizationGetRequest() *AirportOrganizationGetRequest {
	return &AirportOrganizationGetRequest{
		Header:     common.NewHeader(),
		MethodName: "airport_organization_get",
		DataValue:  nil,
	}
}

func (r *AirportOrganizationGetRequest) Method() string { return r.MethodName }
func (r *AirportOrganizationGetRequest) Data() any      { return r.DataValue }

// AirportOrganizationBindData represents the airport organization bind data
type AirportOrganizationBindData struct {
	OrganizationID string `json:"organization_id"` // Organization ID to bind
	AccessKey      string `json:"access_key"`      // Access key for binding
}

// AirportOrganizationBindRequest represents the airport organization bind request
type AirportOrganizationBindRequest struct {
	common.Header
	MethodName string                      `json:"method"`
	DataValue  AirportOrganizationBindData `json:"data"`
}

// NewAirportOrganizationBindRequest creates a new airport organization bind request
func NewAirportOrganizationBindRequest(data AirportOrganizationBindData) *AirportOrganizationBindRequest {
	return &AirportOrganizationBindRequest{
		Header:     common.NewHeader(),
		MethodName: "airport_organization_bind",
		DataValue:  data,
	}
}

func (r *AirportOrganizationBindRequest) Method() string { return r.MethodName }
func (r *AirportOrganizationBindRequest) Data() any      { return r.DataValue }

// GetHeader implements Command.GetHeader
func (r *AirportBindStatusRequest) GetHeader() *common.Header {
	return &r.Header
}

// GetHeader implements Command.GetHeader
func (r *AirportOrganizationBindRequest) GetHeader() *common.Header {
	return &r.Header
}

// GetHeader implements Command.GetHeader
func (r *AirportOrganizationGetRequest) GetHeader() *common.Header {
	return &r.Header
}

// ===============================
// Device → Cloud Request Data Structures
// ===============================

// AirportBindStatusRequestData represents the data in airport_bind_status request from device
type AirportBindStatusRequestData struct {
	Devices []struct {
		SN string `json:"sn"`
	} `json:"devices"`
}

// AirportBindStatusResponseData represents the response data for airport_bind_status
type AirportBindStatusResponseData struct {
	Result int `json:"result"`
	Output struct {
		BindStatus []AirportBindStatusItem `json:"bind_status"`
	} `json:"output"`
}

// AirportBindStatusItem represents a single device's bind status
type AirportBindStatusItem struct {
	SN                       string `json:"sn"`
	IsDeviceBindOrganization bool   `json:"is_device_bind_organization"`
	OrganizationID           string `json:"organization_id"`
	OrganizationName         string `json:"organization_name"`
	DeviceCallsign           string `json:"device_callsign"`
}

// AirportOrganizationGetRequestData represents the data in airport_organization_get request from device
type AirportOrganizationGetRequestData struct {
	DeviceBindingCode string `json:"device_binding_code"`
	OrganizationID    string `json:"organization_id"`
}

// AirportOrganizationGetResponseData represents the response data for airport_organization_get
type AirportOrganizationGetResponseData struct {
	Result int `json:"result"`
	Output struct {
		OrganizationName string `json:"organization_name"`
	} `json:"output"`
}

// AirportOrganizationBindRequestData represents the data in airport_organization_bind request from device
type AirportOrganizationBindRequestData struct {
	BindDevices []AirportBindDeviceItem `json:"bind_devices"`
}

// AirportBindDeviceItem represents a single device's bind parameters
type AirportBindDeviceItem struct {
	DeviceBindingCode string `json:"device_binding_code"`
	OrganizationID    string `json:"organization_id"`
	DeviceCallsign    string `json:"device_callsign"`
	SN                string `json:"sn"`
	DeviceModelKey    string `json:"device_model_key"`
}

// AirportOrganizationBindResponseData represents the response data for airport_organization_bind
type AirportOrganizationBindResponseData struct {
	Result int `json:"result"`
	Output struct {
		ErrInfos []AirportBindErrorInfo `json:"err_infos"`
	} `json:"output"`
}

// AirportBindErrorInfo represents error information for a single device
type AirportBindErrorInfo struct {
	SN      string `json:"sn"`
	ErrCode int    `json:"err_code"`
}
