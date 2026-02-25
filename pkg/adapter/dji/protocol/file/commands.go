package file

import (
	"encoding/json"
	"fmt"

	"github.com/utmos/utmos/pkg/adapter/dji/protocol/common"
)

// ===============================
// File Management Commands
// ===============================

// UploadFlighttaskMediaPrioritizeData represents the upload flight task media prioritize data
type UploadFlighttaskMediaPrioritizeData struct {
	FlightID string `json:"flight_id"` // Flight task ID
}

// UploadFlighttaskMediaPrioritizeCommand represents the upload flight task media prioritize request
type UploadFlighttaskMediaPrioritizeCommand struct {
	common.Header
	MethodName string                              `json:"method"`
	DataValue  UploadFlighttaskMediaPrioritizeData `json:"data"`
}

// NewUploadFlighttaskMediaPrioritizeCommand creates a new upload flight task media prioritize request
func NewUploadFlighttaskMediaPrioritizeCommand(data UploadFlighttaskMediaPrioritizeData) *UploadFlighttaskMediaPrioritizeCommand {
	return &UploadFlighttaskMediaPrioritizeCommand{
		Header:     common.NewHeader(),
		MethodName: "upload_flighttask_media_prioritize",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *UploadFlighttaskMediaPrioritizeCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *UploadFlighttaskMediaPrioritizeCommand) Data() any { return c.DataValue }

// ===============================
// Log File Management Commands (Services)
// 日志文件管理命令
// ===============================

// UploadListData represents the file upload list request data
// 获取设备可上传的文件列表
type UploadListData struct {
	ModuleList []string `json:"module_list"` // 文件所属过滤列表: "0"=飞行器, "3"=机场
}

// UploadListCommand represents the file upload list request
type UploadListCommand struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  UploadListData `json:"data"`
}

// NewUploadListCommand creates a new file upload list request
func NewUploadListCommand(data UploadListData) *UploadListCommand {
	return &UploadListCommand{
		Header:     common.NewHeader(),
		MethodName: "fileupload_list",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *UploadListCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *UploadListCommand) Data() any { return c.DataValue }

// Credentials represents the cloud storage credentials
type Credentials struct {
	AccessKeyID     string `json:"access_key_id"`     // 访问密钥 ID
	AccessKeySecret string `json:"access_key_secret"` // 秘密访问密钥
	Expire          int64  `json:"expire"`            // 访问密钥过期时间 (秒)
	SecurityToken   string `json:"security_token"`    // 会话凭证
}

// UploadStartFile represents a file in the upload start request
type UploadStartFile struct {
	List      json.RawMessage `json:"list"`       // 日志列表
	Module    string          `json:"module"`     // 日志所属模块: "0"=飞行器, "3"=机场
	ObjectKey string          `json:"object_key"` // 文件在对象存储桶的 Key
}

// UploadStartParams represents the params in the upload start request
type UploadStartParams struct {
	Files []UploadStartFile `json:"files"`
}

// UploadStartData represents the file upload start data
// 发起日志文件上传
type UploadStartData struct {
	Bucket      string                `json:"bucket"`      // 对象存储桶名称
	Region      string                `json:"region"`      // 数据中心所在的地域
	Credentials Credentials           `json:"credentials"` // 凭证信息
	Endpoint    string                `json:"endpoint"`    // 对外服务的访问域名
	Provider    string                `json:"provider"`    // 云厂商枚举值: "ali"=阿里云, "aws"=亚马逊云, "minio"=minio
	Params      UploadStartParams `json:"params"`
}

// UploadStartCommand represents the file upload start request
type UploadStartCommand struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  UploadStartData `json:"data"`
}

// NewUploadStartCommand creates a new file upload start request
func NewUploadStartCommand(data UploadStartData) *UploadStartCommand {
	return &UploadStartCommand{
		Header:     common.NewHeader(),
		MethodName: "fileupload_start",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *UploadStartCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *UploadStartCommand) Data() any { return c.DataValue }

// UploadUpdateData represents the file upload update data
// 上传状态更新
type UploadUpdateData struct {
	Status     string   `json:"status"`      // 上传状态: "cancel"=取消
	ModuleList []string `json:"module_list"` // 日志所属模块列表: "0"=飞行器, "3"=机场
}

// UploadUpdateCommand represents the file upload update request
type UploadUpdateCommand struct {
	common.Header
	MethodName string               `json:"method"`
	DataValue  UploadUpdateData `json:"data"`
}

// NewUploadUpdateCommand creates a new file upload update request
func NewUploadUpdateCommand(data UploadUpdateData) *UploadUpdateCommand {
	return &UploadUpdateCommand{
		Header:     common.NewHeader(),
		MethodName: "fileupload_update",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *UploadUpdateCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *UploadUpdateCommand) Data() any { return c.DataValue }

// ===============================
// Service Reply Structures (services_reply)
// 服务响应结构体
// ===============================

// UploadListItem represents a file item in the file list
// 支持两种格式:
// 1. DJI 官方文档格式: boot_index, start_time, end_time, size
// 2. 实际设备返回格式: file_name, file_size, path, start_time, end_time
type UploadListItem struct {
	// DJI 官方文档格式字段
	BootIndex int `json:"boot_index,omitempty"` // 文件索引

	// 实际设备返回格式字段
	FileName string `json:"file_name,omitempty"` // 文件名
	Path     string `json:"path,omitempty"`      // 文件路径

	// 通用字段 (两种格式都有)
	StartTime common.FlexInt64 `json:"start_time"`          // 日志开始时间 (秒), 可能为空字符串
	EndTime   common.FlexInt64 `json:"end_time"`            // 日志结束时间 (秒)
	Size      int              `json:"size,omitempty"`      // 日志文件大小 (bytes) - 官方文档格式
	FileSize  int              `json:"file_size,omitempty"` // 日志文件大小 (bytes) - 实际设备格式
}

// GetFileSize 获取文件大小 (兼容两种格式)
func (f *UploadListItem) GetFileSize() int {
	if f.FileSize > 0 {
		return f.FileSize
	}
	return f.Size
}

// GetFileName 获取文件名 (兼容两种格式)
func (f *UploadListItem) GetFileName() string {
	if f.FileName != "" {
		return f.FileName
	}
	return fmt.Sprintf("log_%d", f.BootIndex)
}

// UploadListFile represents a file group in the file list reply
type UploadListFile struct {
	DeviceSN string           `json:"device_sn"` // 设备序列号 (SN)
	Result   int              `json:"result"`    // 返回码 (非0代表错误)
	Module   common.FlexInt   `json:"module"`    // 所属设备类型: 0=飞行器, 3=机场 (可能是数字或字符串)
	List     []UploadListItem `json:"list"`      // 文件索引列表
}

// UploadListOutput represents the output field in the reply
type UploadListOutput struct {
	Status string `json:"status"` // 状态: "ok"
}

// UploadListReplyData represents the file upload list reply data
type UploadListReplyData struct {
	Result    int                `json:"result"`              // 返回码 (非0代表错误)
	ResultMsg string             `json:"resultMsg,omitempty"` // 返回消息
	Files     []UploadListFile   `json:"files"`               // 文件列表
	Output    *UploadListOutput  `json:"output,omitempty"`    // 输出状态
}

// UploadStartReplyData represents the file upload start reply data
type UploadStartReplyData struct {
	Result int `json:"result"` // 返回码 (非0代表错误)
}

// UploadUpdateReplyData represents the file upload update reply data
type UploadUpdateReplyData struct {
	Result int `json:"result"` // 返回码 (非0代表错误)
}

// GetHeader implements Command.GetHeader
func (c *UploadListCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *UploadStartCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *UploadUpdateCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *UploadFlighttaskMediaPrioritizeCommand) GetHeader() *common.Header {
	return &c.Header
}
