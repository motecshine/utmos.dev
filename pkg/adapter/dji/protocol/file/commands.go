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

// FileUploadListData represents the file upload list request data
// 获取设备可上传的文件列表
type FileUploadListData struct {
	ModuleList []string `json:"module_list"` // 文件所属过滤列表: "0"=飞行器, "3"=机场
}

// FileUploadListCommand represents the file upload list request
type FileUploadListCommand struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  FileUploadListData `json:"data"`
}

// NewFileUploadListCommand creates a new file upload list request
func NewFileUploadListCommand(data FileUploadListData) *FileUploadListCommand {
	return &FileUploadListCommand{
		Header:     common.NewHeader(),
		MethodName: "fileupload_list",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *FileUploadListCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *FileUploadListCommand) Data() any { return c.DataValue }

// Credentials represents the cloud storage credentials
type Credentials struct {
	AccessKeyID     string `json:"access_key_id"`     // 访问密钥 ID
	AccessKeySecret string `json:"access_key_secret"` // 秘密访问密钥
	Expire          int64  `json:"expire"`            // 访问密钥过期时间 (秒)
	SecurityToken   string `json:"security_token"`    // 会话凭证
}

// FileUploadStartFile represents a file in the upload start request
type FileUploadStartFile struct {
	List      json.RawMessage `json:"list"`       // 日志列表
	Module    string          `json:"module"`     // 日志所属模块: "0"=飞行器, "3"=机场
	ObjectKey string          `json:"object_key"` // 文件在对象存储桶的 Key
}

// FileUploadStartParams represents the params in the upload start request
type FileUploadStartParams struct {
	Files []FileUploadStartFile `json:"files"`
}

// FileUploadStartData represents the file upload start data
// 发起日志文件上传
type FileUploadStartData struct {
	Bucket      string                `json:"bucket"`      // 对象存储桶名称
	Region      string                `json:"region"`      // 数据中心所在的地域
	Credentials Credentials           `json:"credentials"` // 凭证信息
	Endpoint    string                `json:"endpoint"`    // 对外服务的访问域名
	Provider    string                `json:"provider"`    // 云厂商枚举值: "ali"=阿里云, "aws"=亚马逊云, "minio"=minio
	Params      FileUploadStartParams `json:"params"`
}

// FileUploadStartCommand represents the file upload start request
type FileUploadStartCommand struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  FileUploadStartData `json:"data"`
}

// NewFileUploadStartCommand creates a new file upload start request
func NewFileUploadStartCommand(data FileUploadStartData) *FileUploadStartCommand {
	return &FileUploadStartCommand{
		Header:     common.NewHeader(),
		MethodName: "fileupload_start",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *FileUploadStartCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *FileUploadStartCommand) Data() any { return c.DataValue }

// FileUploadUpdateData represents the file upload update data
// 上传状态更新
type FileUploadUpdateData struct {
	Status     string   `json:"status"`      // 上传状态: "cancel"=取消
	ModuleList []string `json:"module_list"` // 日志所属模块列表: "0"=飞行器, "3"=机场
}

// FileUploadUpdateCommand represents the file upload update request
type FileUploadUpdateCommand struct {
	common.Header
	MethodName string               `json:"method"`
	DataValue  FileUploadUpdateData `json:"data"`
}

// NewFileUploadUpdateCommand creates a new file upload update request
func NewFileUploadUpdateCommand(data FileUploadUpdateData) *FileUploadUpdateCommand {
	return &FileUploadUpdateCommand{
		Header:     common.NewHeader(),
		MethodName: "fileupload_update",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *FileUploadUpdateCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *FileUploadUpdateCommand) Data() any { return c.DataValue }

// ===============================
// Service Reply Structures (services_reply)
// 服务响应结构体
// ===============================

// FileUploadListItem represents a file item in the file list
// 支持两种格式:
// 1. DJI 官方文档格式: boot_index, start_time, end_time, size
// 2. 实际设备返回格式: file_name, file_size, path, start_time, end_time
type FileUploadListItem struct {
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
func (f *FileUploadListItem) GetFileSize() int {
	if f.FileSize > 0 {
		return f.FileSize
	}
	return f.Size
}

// GetFileName 获取文件名 (兼容两种格式)
func (f *FileUploadListItem) GetFileName() string {
	if f.FileName != "" {
		return f.FileName
	}
	return fmt.Sprintf("log_%d", f.BootIndex)
}

// FileUploadListFile represents a file group in the file list reply
type FileUploadListFile struct {
	DeviceSN string               `json:"device_sn"` // 设备序列号 (SN)
	Result   int                  `json:"result"`    // 返回码 (非0代表错误)
	Module   common.FlexInt       `json:"module"`    // 所属设备类型: 0=飞行器, 3=机场 (可能是数字或字符串)
	List     []FileUploadListItem `json:"list"`      // 文件索引列表
}

// FileUploadListOutput represents the output field in the reply
type FileUploadListOutput struct {
	Status string `json:"status"` // 状态: "ok"
}

// FileUploadListReplyData represents the file upload list reply data
type FileUploadListReplyData struct {
	Result    int                   `json:"result"`              // 返回码 (非0代表错误)
	ResultMsg string                `json:"resultMsg,omitempty"` // 返回消息
	Files     []FileUploadListFile  `json:"files"`               // 文件列表
	Output    *FileUploadListOutput `json:"output,omitempty"`    // 输出状态
}

// FileUploadStartReplyData represents the file upload start reply data
type FileUploadStartReplyData struct {
	Result int `json:"result"` // 返回码 (非0代表错误)
}

// FileUploadUpdateReplyData represents the file upload update reply data
type FileUploadUpdateReplyData struct {
	Result int `json:"result"` // 返回码 (非0代表错误)
}

// GetHeader implements Command.GetHeader
func (c *FileUploadListCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *FileUploadStartCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *FileUploadUpdateCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *UploadFlighttaskMediaPrioritizeCommand) GetHeader() *common.Header {
	return &c.Header
}
