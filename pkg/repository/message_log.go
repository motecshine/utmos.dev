// Package repository provides data access layer implementations.
package repository

import (
	"context"

	"gorm.io/gorm"

	pkgerrors "github.com/utmos/utmos/pkg/errors"
	"github.com/utmos/utmos/pkg/models"
)

// MessageLogRepository provides message log data access.
type MessageLogRepository struct {
	db *gorm.DB
}

// NewMessageLogRepository creates a new MessageLogRepository.
func NewMessageLogRepository(db *gorm.DB) *MessageLogRepository {
	return &MessageLogRepository{db: db}
}

// Create creates a new message log entry.
func (r *MessageLogRepository) Create(ctx context.Context, log *models.MessageLog) error {
	result := r.db.WithContext(ctx).Create(log)
	if result.Error != nil {
		return pkgerrors.Wrap(result.Error, pkgerrors.ErrDatabaseConnection, "failed to create message log")
	}
	return nil
}

// CreateAuthAttempt creates a message log entry for an authentication attempt.
func (r *MessageLogRepository) CreateAuthAttempt(ctx context.Context, username, deviceSN, service string, success bool, errorMsg string) error {
	status := models.MessageStatusSuccess
	var errMsg *string
	if !success {
		status = models.MessageStatusFailed
		errMsg = &errorMsg
	}

	log := &models.MessageLog{
		DeviceSN:     deviceSN,
		Service:      service,
		MessageType:  "auth_attempt",
		Direction:    models.MessageDirectionUplink,
		Status:       status,
		ErrorMessage: errMsg,
		TID:          username, // Using username as TID for auth attempts
	}

	// Store additional data as JSON
	log.MessageData = []byte(`{"username":"` + username + `"}`)

	return r.Create(ctx, log)
}

// CreateMessageProcessingFailure creates a message log entry for a processing failure.
func (r *MessageLogRepository) CreateMessageProcessingFailure(ctx context.Context, deviceSN, tid, bid, service, messageType string, direction models.MessageDirection, errorMsg string) error {
	return r.Create(ctx, &models.MessageLog{
		DeviceSN:     deviceSN,
		TID:          tid,
		BID:          bid,
		Service:      service,
		MessageType:  messageType,
		Direction:    direction,
		Status:       models.MessageStatusFailed,
		ErrorMessage: &errorMsg,
	})
}

// CreateIdempotencyRecord creates a message log entry for an idempotency check.
func (r *MessageLogRepository) CreateIdempotencyRecord(ctx context.Context, deviceSN, tid, bid, service, messageType string, direction models.MessageDirection) error {
	return r.Create(ctx, &models.MessageLog{
		DeviceSN:    deviceSN,
		TID:         tid,
		BID:         bid,
		Service:     service,
		MessageType: messageType,
		Direction:   direction,
		Status:      models.MessageStatusSuccess,
	})
}

// CreateRetryRecord creates a message log entry for a retry attempt.
func (r *MessageLogRepository) CreateRetryRecord(ctx context.Context, deviceSN, tid, bid, service, messageType string, direction models.MessageDirection, attemptNumber int, errorMsg string) error {
	return r.Create(ctx, &models.MessageLog{
		DeviceSN:     deviceSN,
		TID:          tid,
		BID:          bid,
		Service:      service,
		MessageType:  messageType,
		Direction:    direction,
		Status:       models.MessageStatusPending,
		ErrorMessage: &errorMsg,
	})
}

// CreateTerminalStateRecord creates a message log entry for a terminal state (success/failed/timeout).
func (r *MessageLogRepository) CreateTerminalStateRecord(ctx context.Context, deviceSN, tid, bid, service, messageType string, direction models.MessageDirection, status models.MessageStatus, errorMsg string) error {
	return r.Create(ctx, &models.MessageLog{
		DeviceSN:     deviceSN,
		TID:          tid,
		BID:          bid,
		Service:      service,
		MessageType:  messageType,
		Direction:    direction,
		Status:       status,
		ErrorMessage: &errorMsg,
	})
}

// CreateLateResponseRecord creates a message log entry for a late response after terminal state.
func (r *MessageLogRepository) CreateLateResponseRecord(ctx context.Context, deviceSN, tid, bid, service, messageType string, direction models.MessageDirection, errorMsg string) error {
	return r.Create(ctx, &models.MessageLog{
		DeviceSN:     deviceSN,
		TID:          tid,
		BID:          bid,
		Service:      service,
		MessageType:  messageType,
		Direction:    direction,
		Status:       models.MessageStatusFailed,
		ErrorMessage: &errorMsg,
	})
}
