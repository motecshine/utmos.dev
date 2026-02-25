package config

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Configuration Management Commands
// ===============================

// Data represents the config data
type Data struct {
	ConfigType  string `json:"config_type"`  // Configuration type (json)
	ConfigScope string `json:"config_scope"` // Configuration scope (product)
}

// ResponseData represents the config response data
type ResponseData struct {
	NTPServerHost string `json:"ntp_server_host"`           // NTP server host
	AppID         string `json:"app_id"`                    // App ID from DJI developer website
	AppKey        string `json:"app_key"`                   // App Key from DJI developer website
	AppLicense    string `json:"app_license"`               // App License from DJI developer website
	NTPServerPort int    `json:"ntp_server_port,omitempty"` // NTP server port (default 123)
}

// Request represents the config request
type Request struct {
	common.Header
	MethodName string     `json:"method"`
	DataValue  Data `json:"data"`
}

// NewRequest creates a new config request
func NewRequest(data Data) *Request {
	return &Request{
		Header:     common.NewHeader(),
		MethodName: "config",
		DataValue:  data,
	}
}

// Method returns the method name.
func (r *Request) Method() string { return r.MethodName }

// Data returns the command/event data.
func (r *Request) Data() any { return r.DataValue }

// UpdateTopoRequest represents the update topology request
type UpdateTopoRequest struct {
	common.Header
	MethodName string         `json:"method"`
	DataValue  UpdateTopoData `json:"data"`
}

// UpdateTopoData represents common topology update data
type UpdateTopoData struct {
	// Gateway device information
	Domain       string `json:"domain"`        // 网关设备的命名空间
	Type         int    `json:"type"`          // 网关设备的产品类型
	SubType      int    `json:"sub_type"`      // 网关子设备的产品子类型
	DeviceSecret string `json:"device_secret"` // 网关设备的密钥
	Nonce        string `json:"nonce"`         // nonce
	Version      string `json:"version"`       // 设备的固件版本号
	// Sub-device information
	SubDevices []SubDevice `json:"sub_devices"` // 子设备列表
}

// SubDevice represents a sub-device in topology
type SubDevice struct {
	SN           string `json:"sn"`            // 子设备序列号（SN）
	Domain       string `json:"domain"`        // 子设备的命名空间
	Type         int    `json:"type"`          // 子设备的产品类型
	SubType      int    `json:"sub_type"`      // 子设备的产品子类型
	Index        string `json:"index"`         // 子设备挂载位置索引，注意是字符串类型
	DeviceSecret string `json:"device_secret"` // 子设备的密钥
	Nonce        string `json:"nonce"`         // nonce
	Version      string `json:"version"`       // 设备的固件版本号
}

// NewUpdateTopoRequest creates a new update topology request
func NewUpdateTopoRequest(data UpdateTopoData) *UpdateTopoRequest {
	return &UpdateTopoRequest{
		Header:     common.NewHeader(),
		MethodName: "update_topo",
		DataValue:  data,
	}
}

// Method returns the method name.
func (r *UpdateTopoRequest) Method() string { return r.MethodName }

// Data returns the command/event data.
func (r *UpdateTopoRequest) Data() any { return r.DataValue }

// GetHeader implements Command.GetHeader
func (r *Request) GetHeader() *common.Header {
	return &r.Header
}

// GetHeader implements Command.GetHeader
func (r *UpdateTopoRequest) GetHeader() *common.Header {
	return &r.Header
}
