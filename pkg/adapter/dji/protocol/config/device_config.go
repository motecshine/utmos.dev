package config

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Device Configuration Requests (Device → Cloud)
// ===============================

// RequestData represents the data in config request from device
type RequestData struct {
	ConfigType  string `json:"config_type"`  // "json"
	ConfigScope string `json:"config_scope"` // "product"
}

// Note: ResponseData is defined in requests.go

// StorageConfigGetRequestData represents the data in storage_config_get request from device
type StorageConfigGetRequestData struct {
	Module int `json:"module"` // 0: 媒体
}

// StorageConfigGetResponseData represents the response data for storage_config_get
type StorageConfigGetResponseData struct {
	Result int                        `json:"result"`
	Output StorageConfigGetOutputData `json:"output"`
}

// StorageConfigGetOutputData represents the output data for storage config
type StorageConfigGetOutputData struct {
	Bucket          string             `json:"bucket"`
	Credentials     StorageCredentials `json:"credentials"`
	Endpoint        string             `json:"endpoint"`
	Provider        string             `json:"provider"` // "ali", "aws", "minio"
	Region          string             `json:"region"`
	ObjectKeyPrefix string             `json:"object_key_prefix"`
}

// StorageCredentials represents storage credentials (STS token)
type StorageCredentials struct {
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Expire          int    `json:"expire"`         // seconds
	SecurityToken   string `json:"security_token"` // STS token
}

// ===============================
// Flight Task Resource Request (Device → Cloud)
// ===============================

// FlightTaskResourceGetRequestData represents the data in flighttask_resource_get request from device
type FlightTaskResourceGetRequestData struct {
	FlightID string `json:"flight_id"` // 飞行任务ID
}

// FlightTaskResourceGetResponseData represents the response data for flighttask_resource_get
type FlightTaskResourceGetResponseData struct {
	Result int                             `json:"result"`
	Output FlightTaskResourceGetOutputData `json:"output"`
}

// FlightTaskResourceGetResponse represents the flight task resource get response.
type FlightTaskResourceGetResponse struct {
	common.Header
	MethodName string                            `json:"method"`
	DataValue  FlightTaskResourceGetResponseData `json:"output"`
}

// NewFlightTaskResourceGetResponse creates a new FlightTaskResourceGetResponse.
func NewFlightTaskResourceGetResponse(data FlightTaskResourceGetResponseData) *FlightTaskResourceGetResponse {
	return &FlightTaskResourceGetResponse{
		Header:     common.NewHeader(),
		MethodName: "flighttask_resource_get",
		DataValue:  data,
	}
}

// Method implements Command.Method
func (c *FlightTaskResourceGetResponse) Method() string {
	return c.MethodName
}

// Data implements Command.Data
func (c *FlightTaskResourceGetResponse) Data() any {
	return c.DataValue
}

// GetHeader implements Command.GetHeader
func (c *FlightTaskResourceGetResponse) GetHeader() *common.Header {
	return &c.Header
}

// FlightTaskResourceGetOutputData represents the output data for flight task resource
type FlightTaskResourceGetOutputData struct {
	File FlightTaskFile `json:"file"`
}

// FlightTaskFile represents the flight task file information
type FlightTaskFile struct {
	URL         string `json:"url"`         // 航线文件下载URL（预签名URL）
	Fingerprint string `json:"fingerprint"` // 文件MD5签名
}

// ===============================
// Flight Areas Request (Device → Cloud)
// ===============================

// FlightAreasGetRequestData represents the data in flight_areas_get request from device
type FlightAreasGetRequestData struct {
	// DJI文档未明确请求参数，可能为空或包含设备位置信息
}

// FlightAreasGetResponseData represents the response data for flight_areas_get
type FlightAreasGetResponseData struct {
	Result int                      `json:"result"`
	Output FlightAreasGetOutputData `json:"output"`
}

// FlightAreasGetOutputData represents the output data for flight areas
type FlightAreasGetOutputData struct {
	File []FlightAreasFile `json:"files"`
}

// FlightAreasFile represents the flight areas file information
type FlightAreasFile struct {
	Name     string `json:"name"`     // 文件名
	URL      string `json:"url"`      // 飞行区域文件下载URL（预签名URL）
	Checksum string `json:"checksum"` // 文件MD5签名
	Size     int    `json:"size"`     // 文件大小

}
