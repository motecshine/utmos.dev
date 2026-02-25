package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/utmos/utmos/internal/downlink/dispatcher"
	"github.com/utmos/utmos/internal/downlink/model"
)

// Service handles service call API requests
type Service struct {
	db         *gorm.DB
	logger     *logrus.Entry
	dispatcher *dispatcher.DispatchHandler
	repository *model.ServiceCallRepository
}

// NewService creates a new service handler
func NewService(db *gorm.DB, dispatchHandler *dispatcher.DispatchHandler, logger *logrus.Entry) *Service {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	var repo *model.ServiceCallRepository
	if db != nil {
		repo = model.NewServiceCallRepository(db)
	}

	return &Service{
		db:         db,
		logger:     logger.WithField("handler", "service"),
		dispatcher: dispatchHandler,
		repository: repo,
	}
}

// ServiceCallRequest represents the request body for a service call
type ServiceCallRequest struct {
	DeviceSN   string         `json:"device_sn" binding:"required"`
	Vendor     string         `json:"vendor" binding:"required"`
	Method     string         `json:"method" binding:"required"`
	Params     map[string]any `json:"params,omitempty"`
	CallType   string         `json:"call_type,omitempty"`
	MaxRetries int            `json:"max_retries,omitempty"`
}

// ServiceCallResponse represents the response for a service call
type ServiceCallResponse struct {
	ID          string         `json:"id"`
	DeviceSN    string         `json:"device_sn"`
	Vendor      string         `json:"vendor"`
	Method      string         `json:"method"`
	Params      map[string]any `json:"params,omitempty"`
	CallType    string         `json:"call_type"`
	Status      string         `json:"status"`
	TID         string         `json:"tid,omitempty"`
	BID         string         `json:"bid,omitempty"`
	RetryCount  int            `json:"retry_count"`
	MaxRetries  int            `json:"max_retries"`
	Error       string         `json:"error,omitempty"`
	SentAt      *string        `json:"sent_at,omitempty"`
	CompletedAt *string        `json:"completed_at,omitempty"`
	CreatedAt   string         `json:"created_at"`
}

// ListServiceCallsResponse represents the response for listing service calls
type ListServiceCallsResponse struct {
	ServiceCalls []ServiceCallResponse `json:"service_calls"`
	Total        int64                 `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
}

// toServiceCallResponse converts a model to response
func toServiceCallResponse(call *model.ServiceCall) ServiceCallResponse {
	resp := ServiceCallResponse{
		ID:         call.ID,
		DeviceSN:   call.DeviceSN,
		Vendor:     call.Vendor,
		Method:     call.Method,
		CallType:   string(call.CallType),
		Status:     string(call.Status),
		TID:        call.TID,
		BID:        call.BID,
		RetryCount: call.RetryCount,
		MaxRetries: call.MaxRetries,
		Error:      call.Error,
		CreatedAt:  call.CreatedAt.Format(time.RFC3339),
	}

	if call.Params != nil {
		params, _ := call.GetParams()
		resp.Params = params
	}

	if call.SentAt != nil {
		t := call.SentAt.Format(time.RFC3339)
		resp.SentAt = &t
	}
	if call.CompletedAt != nil {
		t := call.CompletedAt.Format(time.RFC3339)
		resp.CompletedAt = &t
	}

	return resp
}

// toDispatcherServiceCall converts request to dispatcher ServiceCall
func toDispatcherServiceCall(req *ServiceCallRequest) *dispatcher.ServiceCall {
	callType := dispatcher.ServiceCallTypeCommand
	switch req.CallType {
	case "property":
		callType = dispatcher.ServiceCallTypeProperty
	case "config":
		callType = dispatcher.ServiceCallTypeConfig
	}

	// Convert params to JSON
	var paramsJSON json.RawMessage
	if req.Params != nil {
		paramsJSON, _ = json.Marshal(req.Params)
	}

	call := dispatcher.NewServiceCall(req.DeviceSN, req.Vendor, req.Method, paramsJSON)
	call.CallType = callType
	if req.MaxRetries > 0 {
		call.MaxRetries = req.MaxRetries
	}

	return call
}

// requireRepository checks that the service call repository is available.
// Returns true if available. On failure it writes a 503 response and returns false.
func (h *Service) requireRepository(c *gin.Context) bool {
	return requireDependency(c, h.repository, "Service call repository")
}

// Call invokes a service call on a device
// @Summary Invoke a service call
// @Description Invoke a service call on a device
// @Tags services
// @Accept json
// @Produce json
// @Param call body ServiceCallRequest true "Service call request"
// @Success 202 {object} ServiceCallResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/services/call [post]
func (h *Service) Call(c *gin.Context) {
	var req ServiceCallRequest
	if !bindJSON(c, &req) {
		return
	}

	// Create dispatcher service call
	dispatcherCall := toDispatcherServiceCall(&req)

	// Dispatch the call
	if h.dispatcher != nil {
		ctx := c.Request.Context()
		result, err := h.dispatcher.Handle(ctx, dispatcherCall)
		if err != nil {
			logWithTrace(h.logger, c.Request.Context()).WithError(err).WithFields(logrus.Fields{
				"device_sn": req.DeviceSN,
				"method":    req.Method,
			}).Error("Failed to dispatch service call")

			// Still return accepted if we can persist the call
			if h.repository != nil {
				modelCall := h.toModelServiceCall(dispatcherCall)
				modelCall.Status = model.ServiceCallStatusFailed
				modelCall.Error = err.Error()
				_ = h.repository.Create(modelCall)

				c.JSON(http.StatusAccepted, toServiceCallResponse(modelCall))
				return
			}

			respondError(c, http.StatusInternalServerError, "DISPATCH_FAILED", "Failed to dispatch service call")
			return
		}

		// Update call with result
		dispatcherCall.TID = result.MessageID
	}

	// Persist the call if repository is available
	if h.repository != nil {
		modelCall := h.toModelServiceCall(dispatcherCall)
		if err := h.repository.Create(modelCall); err != nil {
			logWithTrace(h.logger, c.Request.Context()).WithError(err).Error("Failed to persist service call")
		}

		logWithTrace(h.logger, c.Request.Context()).WithFields(logrus.Fields{
			"call_id":   modelCall.ID,
			"device_sn": req.DeviceSN,
			"method":    req.Method,
		}).Info("Service call dispatched")

		c.JSON(http.StatusAccepted, toServiceCallResponse(modelCall))
		return
	}

	// Return basic response if no repository
	c.JSON(http.StatusAccepted, gin.H{
		"device_sn": req.DeviceSN,
		"method":    req.Method,
		"status":    "sent",
		"tid":       dispatcherCall.TID,
	})
}

// toModelServiceCall converts dispatcher call to model
func (h *Service) toModelServiceCall(call *dispatcher.ServiceCall) *model.ServiceCall {
	modelCall := &model.ServiceCall{
		ID:         call.ID,
		DeviceSN:   call.DeviceSN,
		Vendor:     call.Vendor,
		Method:     call.Method,
		CallType:   model.ServiceCallType(call.CallType),
		Status:     model.ServiceCallStatus(call.Status),
		TID:        call.TID,
		BID:        call.BID,
		RetryCount: call.RetryCount,
		MaxRetries: call.MaxRetries,
		Error:      call.Error,
		CreatedAt:  call.CreatedAt,
		SentAt:     call.SentAt,
	}

	if len(call.Params) > 0 {
		modelCall.SetParamsRaw(call.Params)
	}

	return modelCall
}

// Get retrieves a service call by ID
// @Summary Get a service call by ID
// @Description Get service call details by ID
// @Tags services
// @Produce json
// @Param id path string true "Service Call ID"
// @Success 200 {object} ServiceCallResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/services/calls/{id} [get]
func (h *Service) Get(c *gin.Context) {
	id, ok := requireStringParam(c, "id", "INVALID_ID", "Service call ID is required")
	if !ok {
		return
	}

	if !h.requireRepository(c) {
		return
	}

	call, err := h.repository.FindByID(id)
	if handleDBLookupError(c, h.logger, err,
		"NOT_FOUND", "Service call not found",
		"Failed to get service call", "Failed to get service call") {
		return
	}

	c.JSON(http.StatusOK, toServiceCallResponse(call))
}

// ListByDevice lists service calls for a device
// @Summary List service calls for a device
// @Description List service calls for a specific device
// @Tags services
// @Produce json
// @Param device_sn path string true "Device Serial Number"
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} ListServiceCallsResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/services/calls/device/{device_sn} [get]
func (h *Service) ListByDevice(c *gin.Context) {
	deviceSN, ok := requireStringParam(c, "device_sn", "INVALID_DEVICE_SN", "Device serial number is required")
	if !ok {
		return
	}

	if !h.requireRepository(c) {
		return
	}

	limit := parseLimit(c, 20, 100)

	calls, err := h.repository.FindByDeviceSN(deviceSN, limit)
	if err != nil {
		respondInternalError(c, h.logger, err, "Failed to list service calls", "Failed to list service calls")
		return
	}

	responses := make([]ServiceCallResponse, len(calls))
	for i, call := range calls {
		responses[i] = toServiceCallResponse(&call)
	}

	c.JSON(http.StatusOK, ListServiceCallsResponse{
		ServiceCalls: responses,
		Total:        int64(len(calls)),
		Page:         1,
		PageSize:     limit,
	})
}

// Cancel cancels a pending service call
// @Summary Cancel a service call
// @Description Cancel a pending service call
// @Tags services
// @Param id path string true "Service Call ID"
// @Success 200 {object} ServiceCallResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/services/calls/{id}/cancel [post]
func (h *Service) Cancel(c *gin.Context) {
	id, ok := requireStringParam(c, "id", "INVALID_ID", "Service call ID is required")
	if !ok {
		return
	}

	if !h.requireRepository(c) {
		return
	}

	call, err := h.repository.FindByID(id)
	if handleDBLookupError(c, h.logger, err,
		"NOT_FOUND", "Service call not found",
		"Failed to get service call", "Failed to cancel service call") {
		return
	}

	// Check if call can be cancelled
	if call.IsCompleted() {
		respondBadRequest(c, "CANNOT_CANCEL", "Service call is already completed")
		return
	}

	// Update status
	call.Status = model.ServiceCallStatusCancelled
	now := time.Now()
	call.CompletedAt = &now

	if err := h.repository.Update(call); err != nil {
		respondInternalError(c, h.logger, err, "Failed to cancel service call", "Failed to cancel service call")
		return
	}

	logWithTrace(h.logger, c.Request.Context()).WithField("call_id", id).Info("Service call cancelled")
	c.JSON(http.StatusOK, toServiceCallResponse(call))
}
