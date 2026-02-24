package config

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Storage Control Commands
// ===============================

// PhotoStorageSetData represents the photo storage set data
type PhotoStorageSetData struct {
	PayloadIndex         string   `json:"payload_index"`          // Camera enumeration value
	PhotoStorageSettings []string `json:"photo_storage_settings"` // Photo storage types: current, wide, zoom, ir (multiple selection)
}

// PhotoStorageSetCommand represents the photo storage set request
type PhotoStorageSetCommand struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  PhotoStorageSetData `json:"data"`
}

// NewPhotoStorageSetCommand creates a new photo storage set request
func NewPhotoStorageSetCommand(data PhotoStorageSetData) *PhotoStorageSetCommand {
	return &PhotoStorageSetCommand{
		Header:     common.NewHeader(),
		MethodName: "photo_storage_set",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *PhotoStorageSetCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *PhotoStorageSetCommand) Data() any { return c.DataValue }

// VideoStorageSetData represents the video storage set data
type VideoStorageSetData struct {
	PayloadIndex         string   `json:"payload_index"`          // Camera enumeration value
	VideoStorageSettings []string `json:"video_storage_settings"` // Video storage types: current, wide, zoom, ir (multiple selection)
}

// VideoStorageSetCommand represents the video storage set request
type VideoStorageSetCommand struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  VideoStorageSetData `json:"data"`
}

// NewVideoStorageSetCommand creates a new video storage set request
func NewVideoStorageSetCommand(data VideoStorageSetData) *VideoStorageSetCommand {
	return &VideoStorageSetCommand{
		Header:     common.NewHeader(),
		MethodName: "video_storage_set",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *VideoStorageSetCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *VideoStorageSetCommand) Data() any { return c.DataValue }

// StorageConfigGetData represents the storage config get data
type StorageConfigGetData struct {
	Module int `json:"module"` // Module enum: 0=media, 1=psdk ui resource
}

// StorageConfigGetRequest represents the storage config get request
type StorageConfigGetRequest struct {
	common.Header
	MethodName string               `json:"method"`
	DataValue  StorageConfigGetData `json:"data"`
}

// NewStorageConfigGetRequest creates a new storage config get request
func NewStorageConfigGetRequest(data StorageConfigGetData) *StorageConfigGetRequest {
	return &StorageConfigGetRequest{
		Header:     common.NewHeader(),
		MethodName: "storage_config_get",
		DataValue:  data,
	}
}

// Method returns the method name.
func (r *StorageConfigGetRequest) Method() string { return r.MethodName }

// Data returns the command/event data.
func (r *StorageConfigGetRequest) Data() any { return r.DataValue }

// GetHeader implements Command.GetHeader
func (c *PhotoStorageSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (r *StorageConfigGetRequest) GetHeader() *common.Header {
	return &r.Header
}

// GetHeader implements Command.GetHeader
func (c *VideoStorageSetCommand) GetHeader() *common.Header {
	return &c.Header
}
