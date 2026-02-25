// Package model provides data models for iot-downlink
package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// ServiceCallStatus represents the status of a service call
type ServiceCallStatus string

const (
	// ServiceCallStatusPending indicates the service call is pending.
	ServiceCallStatusPending ServiceCallStatus = "pending"
	// ServiceCallStatusSent indicates the service call has been sent.
	ServiceCallStatusSent ServiceCallStatus = "sent"
	// ServiceCallStatusSuccess indicates the service call succeeded.
	ServiceCallStatusSuccess ServiceCallStatus = "success"
	// ServiceCallStatusFailed indicates the service call failed.
	ServiceCallStatusFailed ServiceCallStatus = "failed"
	// ServiceCallStatusTimeout indicates the service call timed out.
	ServiceCallStatusTimeout ServiceCallStatus = "timeout"
	// ServiceCallStatusRetrying indicates the service call is being retried.
	ServiceCallStatusRetrying ServiceCallStatus = "retrying"
	// ServiceCallStatusCancelled indicates the service call was cancelled.
	ServiceCallStatusCancelled ServiceCallStatus = "cancelled"
)

// ServiceCallType represents the type of service call
type ServiceCallType string

const (
	// ServiceCallTypeCommand is the command service call type.
	ServiceCallTypeCommand ServiceCallType = "command"
	// ServiceCallTypeProperty is the property service call type.
	ServiceCallTypeProperty ServiceCallType = "property"
	// ServiceCallTypeConfig is the config service call type.
	ServiceCallTypeConfig ServiceCallType = "config"
)

// ServiceCall represents a service call record in the database
type ServiceCall struct {
	ID          string            `gorm:"primaryKey;type:varchar(36)" json:"id"`
	DeviceSN    string            `gorm:"type:varchar(64);index;not null" json:"device_sn"`
	Vendor      string            `gorm:"type:varchar(32);index;not null" json:"vendor"`
	Method      string            `gorm:"type:varchar(64);not null" json:"method"`
	Params      json.RawMessage   `gorm:"type:jsonb" json:"params"`
	CallType    ServiceCallType   `gorm:"type:varchar(32);not null;default:'command'" json:"call_type"`
	Status      ServiceCallStatus `gorm:"type:varchar(32);index;not null;default:'pending'" json:"status"`
	TID         string            `gorm:"type:varchar(36);index" json:"tid"`
	BID         string            `gorm:"type:varchar(36)" json:"bid"`
	RetryCount  int               `gorm:"default:0" json:"retry_count"`
	MaxRetries  int               `gorm:"default:3" json:"max_retries"`
	Error       string            `gorm:"type:text" json:"error,omitempty"`
	Response    json.RawMessage   `gorm:"type:jsonb" json:"response,omitempty"`
	SentAt      *time.Time        `gorm:"index" json:"sent_at,omitempty"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	CreatedAt   time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName returns the table name for ServiceCall
func (ServiceCall) TableName() string {
	return "service_calls"
}

// BeforeCreate sets default values before creating
func (s *ServiceCall) BeforeCreate(tx *gorm.DB) error {
	if s.Status == "" {
		s.Status = ServiceCallStatusPending
	}
	if s.CallType == "" {
		s.CallType = ServiceCallTypeCommand
	}
	if s.MaxRetries == 0 {
		s.MaxRetries = 3
	}
	return nil
}

// marshalToJSON marshals a map into json.RawMessage
func marshalToJSON(data map[string]any) (json.RawMessage, error) {
	return json.Marshal(data)
}

// unmarshalFromJSON unmarshals json.RawMessage into a map
func unmarshalFromJSON(data json.RawMessage) (map[string]any, error) {
	if data == nil {
		return nil, nil
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// setJSONField marshals data to JSON and assigns it to the given field
func (s *ServiceCall) setJSONField(field *json.RawMessage, data map[string]any) error {
	result, err := marshalToJSON(data)
	if err != nil {
		return err
	}
	*field = result
	return nil
}

// SetParams sets the params from a map
func (s *ServiceCall) SetParams(params map[string]any) error {
	return s.setJSONField(&s.Params, params)
}

// SetParamsRaw sets the params from raw JSON
func (s *ServiceCall) SetParamsRaw(params json.RawMessage) {
	s.Params = params
}

// GetParams returns the params as a map
func (s *ServiceCall) GetParams() (map[string]any, error) {
	return unmarshalFromJSON(s.Params)
}

// SetResponse sets the response from a map
func (s *ServiceCall) SetResponse(response map[string]any) error {
	return s.setJSONField(&s.Response, response)
}

// GetResponse returns the response as a map
func (s *ServiceCall) GetResponse() (map[string]any, error) {
	return unmarshalFromJSON(s.Response)
}

// MarkSent marks the service call as sent
func (s *ServiceCall) MarkSent() {
	now := time.Now()
	s.SentAt = &now
	s.Status = ServiceCallStatusSent
}

// MarkSuccess marks the service call as successful
func (s *ServiceCall) MarkSuccess(response map[string]any) error {
	now := time.Now()
	s.CompletedAt = &now
	s.Status = ServiceCallStatusSuccess
	if response != nil {
		return s.SetResponse(response)
	}
	return nil
}

// MarkFailed marks the service call as failed
func (s *ServiceCall) MarkFailed(err string) {
	now := time.Now()
	s.CompletedAt = &now
	s.Status = ServiceCallStatusFailed
	s.Error = err
}

// MarkTimeout marks the service call as timed out
func (s *ServiceCall) MarkTimeout() {
	now := time.Now()
	s.CompletedAt = &now
	s.Status = ServiceCallStatusTimeout
}

// MarkRetrying marks the service call as retrying
func (s *ServiceCall) MarkRetrying() {
	s.Status = ServiceCallStatusRetrying
	s.RetryCount++
}

// CanRetry checks if the service call can be retried
func (s *ServiceCall) CanRetry() bool {
	return s.RetryCount < s.MaxRetries &&
		(s.Status == ServiceCallStatusFailed || s.Status == ServiceCallStatusTimeout)
}

// IsPending checks if the service call is pending
func (s *ServiceCall) IsPending() bool {
	return s.Status == ServiceCallStatusPending
}

// IsCompleted checks if the service call is completed
func (s *ServiceCall) IsCompleted() bool {
	return s.Status == ServiceCallStatusSuccess ||
		s.Status == ServiceCallStatusFailed ||
		s.Status == ServiceCallStatusTimeout ||
		s.Status == ServiceCallStatusCancelled
}

// ServiceCallRepository provides database operations for ServiceCall
type ServiceCallRepository struct {
	db *gorm.DB
}

// NewServiceCallRepository creates a new repository
func NewServiceCallRepository(db *gorm.DB) *ServiceCallRepository {
	return &ServiceCallRepository{db: db}
}

// Create creates a new service call
func (r *ServiceCallRepository) Create(call *ServiceCall) error {
	return r.db.Create(call).Error
}

// Update updates a service call
func (r *ServiceCallRepository) Update(call *ServiceCall) error {
	return r.db.Save(call).Error
}

// findByField finds a single service call by a given field and value
func (r *ServiceCallRepository) findByField(field, value string) (*ServiceCall, error) {
	var call ServiceCall
	if err := r.db.First(&call, field+" = ?", value).Error; err != nil {
		return nil, err
	}
	return &call, nil
}

// FindByID finds a service call by ID
func (r *ServiceCallRepository) FindByID(id string) (*ServiceCall, error) {
	return r.findByField("id", id)
}

// FindByTID finds a service call by TID
func (r *ServiceCallRepository) FindByTID(tid string) (*ServiceCall, error) {
	return r.findByField("tid", tid)
}

// findWithQuery executes a query with an optional limit and returns matching service calls
func (r *ServiceCallRepository) findWithQuery(query *gorm.DB, limit int) ([]ServiceCall, error) {
	var calls []ServiceCall
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&calls).Error; err != nil {
		return nil, err
	}
	return calls, nil
}

// FindByDeviceSN finds service calls by device serial number
func (r *ServiceCallRepository) FindByDeviceSN(deviceSN string, limit int) ([]ServiceCall, error) {
	query := r.db.Where("device_sn = ?", deviceSN).Order("created_at DESC")
	return r.findWithQuery(query, limit)
}

// FindPending finds pending service calls
func (r *ServiceCallRepository) FindPending(limit int) ([]ServiceCall, error) {
	query := r.db.Where("status = ?", ServiceCallStatusPending).Order("created_at ASC")
	return r.findWithQuery(query, limit)
}

// FindRetryable finds service calls that can be retried
func (r *ServiceCallRepository) FindRetryable(limit int) ([]ServiceCall, error) {
	query := r.db.Where("status IN ? AND retry_count < max_retries",
		[]ServiceCallStatus{ServiceCallStatusFailed, ServiceCallStatusTimeout}).
		Order("created_at ASC")
	return r.findWithQuery(query, limit)
}

// UpdateStatus updates the status of a service call
func (r *ServiceCallRepository) UpdateStatus(id string, status ServiceCallStatus) error {
	return r.db.Model(&ServiceCall{}).Where("id = ?", id).Update("status", status).Error
}

// Delete deletes a service call
func (r *ServiceCallRepository) Delete(id string) error {
	return r.db.Delete(&ServiceCall{}, "id = ?", id).Error
}

// AutoMigrate runs database migrations
func (r *ServiceCallRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&ServiceCall{})
}
