package wayline

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Wayline Mission Requests (Device → Cloud)
// ===============================

// ResourceGetRequest represents the get flight task resource request
type ResourceGetRequest struct {
	common.Header
	MethodName string          `json:"method"`
	DataValue  ResourceGetData `json:"data"`
}

// NewResourceGetRequest creates a new flight task resource get request
func NewResourceGetRequest(data ResourceGetData) *ResourceGetRequest {
	return &ResourceGetRequest{
		Header:     common.NewHeader(),
		MethodName: "flighttask_resource_get",
		DataValue:  data,
	}
}

func (r *ResourceGetRequest) Method() string            { return r.MethodName }
func (r *ResourceGetRequest) Data() any                 { return r.DataValue }
func (r *ResourceGetRequest) GetHeader() *common.Header { return &r.Header }

// ===============================
// Request Data Types
// ===============================

// ResourceGetData represents the get flight task resource data
type ResourceGetData struct {
	FlighttaskID string `json:"flighttask_id"` // Flight task ID
}
