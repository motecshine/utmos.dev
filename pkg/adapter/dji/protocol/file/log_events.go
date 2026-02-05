package file

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Log File Upload Events
// 文件上传进度通知
// ===============================

// FileUploadProgress represents the file upload progress information
type FileUploadProgress struct {
	Progress    int    `json:"progress"`     // 进度值 (0-100)
	CurrentStep int    `json:"current_step"` // 当前步骤
	TotalStep   int    `json:"total_step"`   // 总步骤数
	FinishTime  int64  `json:"finish_time"`  // 上传完成时间 (ms timestamp)
	UploadRate  int    `json:"upload_rate"`  // 上传速率 (bytes/s)
	Result      int    `json:"result"`       // 返回码 (0=成功)
	Status      string `json:"status"`       // 上传状态 (ok, failed, etc.)
}

// FileUploadItem represents a single file upload item
type FileUploadItem struct {
	Module      string             `json:"module"`      // 所属设备类型: "0"=飞行器, "3"=机场
	Size        int                `json:"size"`        // 文件大小 (bytes)
	DeviceSN    string             `json:"device_sn"`   // 设备序列号 (SN)
	Key         string             `json:"key"`         // 对象存储桶 Key
	Fingerprint string             `json:"fingerprint"` // 文件指纹 (MD5)
	Progress    FileUploadProgress `json:"progress"`    // 进度信息
}

// FileUploadProgressExt represents the file upload extended information
type FileUploadProgressExt struct {
	Files []FileUploadItem `json:"files"` // 文件列表
}

// FileUploadProgressOutput represents the file upload output
type FileUploadProgressOutput struct {
	Ext    FileUploadProgressExt `json:"ext"`    // 扩展信息
	Status string                `json:"status"` // 状态 (ok, failed, etc.)
}

// FileUploadProgressData represents the file upload progress data
type FileUploadProgressData struct {
	Result int                      `json:"result"` // Return code (0=success)
	Output FileUploadProgressOutput `json:"output"` // Output data
}

// FileUploadProgressEvent represents the log file upload progress event
type FileUploadProgressEvent struct {
	common.Header
	MethodName string                 `json:"method"`
	DataValue  FileUploadProgressData `json:"data"`
}

func (e *FileUploadProgressEvent) Method() string            { return e.MethodName }
func (e *FileUploadProgressEvent) Data() any                 { return e.DataValue }
func (e *FileUploadProgressEvent) GetHeader() *common.Header { return &e.Header }
