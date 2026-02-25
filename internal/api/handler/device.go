// Package handler provides HTTP handlers for iot-api
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/utmos/utmos/pkg/models"
)

// Device handles device-related API requests
type Device struct {
	db     *gorm.DB
	logger *logrus.Entry
}

// NewDevice creates a new device handler
func NewDevice(db *gorm.DB, logger *logrus.Entry) *Device {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &Device{
		db:     db,
		logger: logger.WithField("handler", "device"),
	}
}

// CreateDeviceRequest represents the request body for creating a device
type CreateDeviceRequest struct {
	DeviceSN     string  `json:"device_sn" binding:"required"`
	DeviceName   string  `json:"device_name" binding:"required"`
	DeviceType   string  `json:"device_type" binding:"required"`
	Vendor       string  `json:"vendor"`
	GatewaySN    *string `json:"gateway_sn,omitempty"`
	ThingModelID *uint   `json:"thing_model_id,omitempty"`
}

// UpdateDeviceRequest represents the request body for updating a device
type UpdateDeviceRequest struct {
	DeviceName   *string              `json:"device_name,omitempty"`
	DeviceType   *string              `json:"device_type,omitempty"`
	Vendor       *string              `json:"vendor,omitempty"`
	Status       *models.DeviceStatus `json:"status,omitempty"`
	GatewaySN    *string              `json:"gateway_sn,omitempty"`
	ThingModelID *uint                `json:"thing_model_id,omitempty"`
}

// DeviceResponse represents the response for a device
type DeviceResponse struct {
	ID             uint                `json:"id"`
	DeviceSN       string              `json:"device_sn"`
	DeviceName     string              `json:"device_name"`
	DeviceType     string              `json:"device_type"`
	Vendor         string              `json:"vendor"`
	Status         models.DeviceStatus `json:"status"`
	GatewaySN      *string             `json:"gateway_sn,omitempty"`
	ThingModelID   *uint               `json:"thing_model_id,omitempty"`
	LastOnlineTime *string             `json:"last_online_time,omitempty"`
	CreatedAt      string              `json:"created_at"`
	UpdatedAt      string              `json:"updated_at"`
}

// ListDevicesResponse represents the response for listing devices
type ListDevicesResponse struct {
	Devices    []DeviceResponse `json:"devices"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// toDeviceResponse converts a device model to response
func toDeviceResponse(device *models.Device) DeviceResponse {
	resp := DeviceResponse{
		ID:           device.ID,
		DeviceSN:     device.DeviceSN,
		DeviceName:   device.DeviceName,
		DeviceType:   device.DeviceType,
		Vendor:       device.Vendor,
		Status:       device.Status,
		GatewaySN:    device.GatewaySN,
		ThingModelID: device.ThingModelID,
		CreatedAt:    device.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    device.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if device.LastOnlineTime != nil {
		t := device.LastOnlineTime.Format("2006-01-02T15:04:05Z")
		resp.LastOnlineTime = &t
	}
	return resp
}

// Create creates a new device
// @Summary Create a new device
// @Description Create a new device with the provided information
// @Tags devices
// @Accept json
// @Produce json
// @Param device body CreateDeviceRequest true "Device information"
// @Success 201 {object} DeviceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/devices [post]
func (h *Device) Create(c *gin.Context) {
	var req CreateDeviceRequest
	if !bindJSON(c, &req) {
		return
	}

	device := &models.Device{
		DeviceSN:     req.DeviceSN,
		DeviceName:   req.DeviceName,
		DeviceType:   req.DeviceType,
		Vendor:       req.Vendor,
		GatewaySN:    req.GatewaySN,
		ThingModelID: req.ThingModelID,
		Status:       models.DeviceStatusUnknown,
	}

	if device.Vendor == "" {
		device.Vendor = "generic"
	}

	if err := h.db.Create(device).Error; err != nil {
		if isUniqueConstraintError(err) {
			respondError(c, http.StatusConflict, "DEVICE_EXISTS", "Device with this serial number already exists")
			return
		}
		respondInternalError(c, h.logger, err, "Failed to create device", "Failed to create device")
		return
	}

	logWithTrace(h.logger, c.Request.Context()).WithField("device_sn", device.DeviceSN).Info("Device created")
	c.JSON(http.StatusCreated, toDeviceResponse(device))
}

// Get retrieves a device by ID
// @Summary Get a device by ID
// @Description Get device details by ID
// @Tags devices
// @Produce json
// @Param id path int true "Device ID"
// @Success 200 {object} DeviceResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/devices/{id} [get]
func (h *Device) Get(c *gin.Context) {
	id, ok := parseUintID(c, "id")
	if !ok {
		return
	}

	var device models.Device
	if handleDBLookupError(c, h.logger, h.db.First(&device, id).Error,
		"DEVICE_NOT_FOUND", "Device not found",
		"Failed to get device", "Failed to get device") {
		return
	}

	c.JSON(http.StatusOK, toDeviceResponse(&device))
}

// GetBySN retrieves a device by serial number
// @Summary Get a device by serial number
// @Description Get device details by serial number
// @Tags devices
// @Produce json
// @Param sn path string true "Device Serial Number"
// @Success 200 {object} DeviceResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/devices/sn/{sn} [get]
func (h *Device) GetBySN(c *gin.Context) {
	sn, ok := requireStringParam(c, "sn", "INVALID_SN", "Device serial number is required")
	if !ok {
		return
	}

	var device models.Device
	if handleDBLookupError(c, h.logger, h.db.Where("device_sn = ?", sn).First(&device).Error,
		"DEVICE_NOT_FOUND", "Device not found",
		"Failed to get device", "Failed to get device") {
		return
	}

	c.JSON(http.StatusOK, toDeviceResponse(&device))
}

// List lists devices with pagination and filtering
// @Summary List devices
// @Description List devices with pagination and optional filtering
// @Tags devices
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param vendor query string false "Filter by vendor"
// @Param status query string false "Filter by status"
// @Param device_type query string false "Filter by device type"
// @Success 200 {object} ListDevicesResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/devices [get]
func (h *Device) List(c *gin.Context) {
	page, pageSize, offset := parsePagination(c, 20, 100)

	query := h.db.Model(&models.Device{})

	// Apply filters
	if vendor := c.Query("vendor"); vendor != "" {
		query = query.Where("vendor = ?", vendor)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if deviceType := c.Query("device_type"); deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		respondInternalError(c, h.logger, err, "Failed to count devices", "Failed to list devices")
		return
	}

	// Fetch devices
	var devices []models.Device
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&devices).Error; err != nil {
		respondInternalError(c, h.logger, err, "Failed to list devices", "Failed to list devices")
		return
	}

	// Convert to response
	deviceResponses := make([]DeviceResponse, len(devices))
	for i, device := range devices {
		deviceResponses[i] = toDeviceResponse(&device)
	}

	c.JSON(http.StatusOK, ListDevicesResponse{
		Devices:    deviceResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages(total, pageSize),
	})
}

// Update updates a device
// @Summary Update a device
// @Description Update device information
// @Tags devices
// @Accept json
// @Produce json
// @Param id path int true "Device ID"
// @Param device body UpdateDeviceRequest true "Device information to update"
// @Success 200 {object} DeviceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/devices/{id} [put]
func (h *Device) Update(c *gin.Context) {
	id, ok := parseUintID(c, "id")
	if !ok {
		return
	}

	var req UpdateDeviceRequest
	if !bindJSON(c, &req) {
		return
	}

	var device models.Device
	if handleDBLookupError(c, h.logger, h.db.First(&device, id).Error,
		"DEVICE_NOT_FOUND", "Device not found",
		"Failed to get device", "Failed to update device") {
		return
	}

	// Update fields using struct
	if req.DeviceName != nil {
		device.DeviceName = *req.DeviceName
	}
	if req.DeviceType != nil {
		device.DeviceType = *req.DeviceType
	}
	if req.Vendor != nil {
		device.Vendor = *req.Vendor
	}
	if req.Status != nil {
		device.Status = *req.Status
	}
	if req.GatewaySN != nil {
		device.GatewaySN = req.GatewaySN
	}
	if req.ThingModelID != nil {
		device.ThingModelID = req.ThingModelID
	}

	if err := h.db.Save(&device).Error; err != nil {
		respondInternalError(c, h.logger, err, "Failed to update device", "Failed to update device")
		return
	}

	// Reload device
	h.db.First(&device, id)

	logWithTrace(h.logger, c.Request.Context()).WithField("device_id", id).Info("Device updated")
	c.JSON(http.StatusOK, toDeviceResponse(&device))
}

// Delete deletes a device
// @Summary Delete a device
// @Description Delete a device by ID
// @Tags devices
// @Param id path int true "Device ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/devices/{id} [delete]
func (h *Device) Delete(c *gin.Context) {
	id, ok := parseUintID(c, "id")
	if !ok {
		return
	}

	result := h.db.Delete(&models.Device{}, id)
	if result.Error != nil {
		respondInternalError(c, h.logger, result.Error, "Failed to delete device", "Failed to delete device")
		return
	}

	if result.RowsAffected == 0 {
		respondNotFound(c, "DEVICE_NOT_FOUND", "Device not found")
		return
	}

	logWithTrace(h.logger, c.Request.Context()).WithField("device_id", id).Info("Device deleted")
	c.Status(http.StatusNoContent)
}
