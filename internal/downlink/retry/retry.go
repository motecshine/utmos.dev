// Package retry provides retry mechanisms for failed service calls
package retry

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/utmos/utmos/internal/downlink/dispatcher"
)

// Retry configuration defaults
const (
	// DefaultMaxRetries is the default maximum number of retry attempts
	DefaultMaxRetries = 3
	// DefaultInitialDelay is the default initial delay before first retry
	DefaultInitialDelay = time.Second
	// DefaultMaxDelay is the default maximum delay between retries
	DefaultMaxDelay = 30 * time.Second
	// DefaultMultiplier is the default exponential backoff multiplier
	DefaultMultiplier = 2.0
)

// Config holds retry configuration
type Config struct {
	// MaxRetries is the maximum number of retry attempts
	MaxRetries int
	// InitialDelay is the initial delay before first retry
	InitialDelay time.Duration
	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration
	// Multiplier is the exponential backoff multiplier
	Multiplier float64
	// EnableDeadLetter enables dead letter queue for failed calls
	EnableDeadLetter bool
}

// DefaultConfig returns default retry configuration
func DefaultConfig() *Config {
	return &Config{
		MaxRetries:       DefaultMaxRetries,
		InitialDelay:     DefaultInitialDelay,
		MaxDelay:         DefaultMaxDelay,
		Multiplier:       DefaultMultiplier,
		EnableDeadLetter: true,
	}
}

// RetryableCall represents a call that can be retried
type RetryableCall struct {
	Call       *dispatcher.ServiceCall
	RetryCount int
	NextRetry  time.Time
	LastError  string
}

// DeadLetterEntry represents a failed call in the dead letter queue
type DeadLetterEntry struct {
	Call      *dispatcher.ServiceCall
	Error     string
	FailedAt  time.Time
	Retries   int
}

// Handler handles retry logic for failed service calls
type Handler struct {
	config      *Config
	logger      *logrus.Entry
	retryQueue  []*RetryableCall
	deadLetter  []*DeadLetterEntry
	mu          sync.RWMutex
	onRetry     func(ctx context.Context, call *dispatcher.ServiceCall) error
	onDeadLetter func(entry *DeadLetterEntry)
}

// NewHandler creates a new retry handler
func NewHandler(config *Config, logger *logrus.Entry) *Handler {
	if config == nil {
		config = DefaultConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &Handler{
		config:     config,
		logger:     logger.WithField("component", "retry-handler"),
		retryQueue: make([]*RetryableCall, 0),
		deadLetter: make([]*DeadLetterEntry, 0),
	}
}

// SetOnRetry sets the callback for retry attempts
func (h *Handler) SetOnRetry(callback func(ctx context.Context, call *dispatcher.ServiceCall) error) {
	h.onRetry = callback
}

// SetOnDeadLetter sets the callback for dead letter entries
func (h *Handler) SetOnDeadLetter(callback func(entry *DeadLetterEntry)) {
	h.onDeadLetter = callback
}

// ScheduleRetry schedules a failed call for retry
func (h *Handler) ScheduleRetry(call *dispatcher.ServiceCall, err string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if max retries exceeded
	if call.RetryCount >= h.config.MaxRetries {
		h.logger.WithFields(logrus.Fields{
			"call_id":     call.ID,
			"device_sn":   call.DeviceSN,
			"retry_count": call.RetryCount,
			"max_retries": h.config.MaxRetries,
		}).Warn("Max retries exceeded, moving to dead letter queue")

		if h.config.EnableDeadLetter {
			h.addToDeadLetter(call, err)
		}
		return false
	}

	// Calculate next retry time with exponential backoff
	delay := h.calculateDelay(call.RetryCount)
	nextRetry := time.Now().Add(delay)

	retryable := &RetryableCall{
		Call:       call,
		RetryCount: call.RetryCount,
		NextRetry:  nextRetry,
		LastError:  err,
	}

	h.retryQueue = append(h.retryQueue, retryable)

	h.logger.WithFields(logrus.Fields{
		"call_id":     call.ID,
		"device_sn":   call.DeviceSN,
		"retry_count": call.RetryCount + 1,
		"next_retry":  nextRetry,
		"delay":       delay,
	}).Debug("Scheduled retry")

	return true
}

// calculateDelay calculates the delay for a retry attempt using exponential backoff
func (h *Handler) calculateDelay(retryCount int) time.Duration {
	delay := float64(h.config.InitialDelay) * math.Pow(h.config.Multiplier, float64(retryCount))
	if delay > float64(h.config.MaxDelay) {
		delay = float64(h.config.MaxDelay)
	}
	return time.Duration(delay)
}

// addToDeadLetter adds a failed call to the dead letter queue
func (h *Handler) addToDeadLetter(call *dispatcher.ServiceCall, err string) {
	entry := &DeadLetterEntry{
		Call:     call,
		Error:    err,
		FailedAt: time.Now(),
		Retries:  call.RetryCount,
	}

	h.deadLetter = append(h.deadLetter, entry)

	h.logger.WithFields(logrus.Fields{
		"call_id":   call.ID,
		"device_sn": call.DeviceSN,
		"error":     err,
	}).Warn("Added to dead letter queue")

	if h.onDeadLetter != nil {
		h.onDeadLetter(entry)
	}
}

// ProcessRetries processes pending retries that are due
func (h *Handler) ProcessRetries(ctx context.Context) int {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now()
	processed := 0
	remaining := make([]*RetryableCall, 0)

	for _, retryable := range h.retryQueue {
		if retryable.NextRetry.After(now) {
			remaining = append(remaining, retryable)
			continue
		}

		// Increment retry count
		retryable.Call.RetryCount++
		retryable.Call.Status = dispatcher.ServiceCallStatusRetrying

		h.logger.WithFields(logrus.Fields{
			"call_id":     retryable.Call.ID,
			"device_sn":   retryable.Call.DeviceSN,
			"retry_count": retryable.Call.RetryCount,
		}).Debug("Processing retry")

		if h.onRetry != nil {
			if err := h.onRetry(ctx, retryable.Call); err != nil {
				h.logger.WithError(err).WithField("call_id", retryable.Call.ID).Error("Retry failed")

				// Check if we should retry again or move to dead letter
				if retryable.Call.RetryCount >= h.config.MaxRetries {
					if h.config.EnableDeadLetter {
						h.addToDeadLetter(retryable.Call, err.Error())
					}
				} else {
					// Schedule another retry
					delay := h.calculateDelay(retryable.Call.RetryCount)
					retryable.NextRetry = time.Now().Add(delay)
					retryable.LastError = err.Error()
					remaining = append(remaining, retryable)
				}
			}
		}
		processed++
	}

	h.retryQueue = remaining
	return processed
}

// GetPendingRetries returns the number of pending retries
func (h *Handler) GetPendingRetries() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.retryQueue)
}

// GetDeadLetterCount returns the number of dead letter entries
func (h *Handler) GetDeadLetterCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.deadLetter)
}

// GetDeadLetterEntries returns all dead letter entries
func (h *Handler) GetDeadLetterEntries() []*DeadLetterEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	entries := make([]*DeadLetterEntry, len(h.deadLetter))
	copy(entries, h.deadLetter)
	return entries
}

// ClearDeadLetter clears the dead letter queue
func (h *Handler) ClearDeadLetter() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.deadLetter = make([]*DeadLetterEntry, 0)
}

// RemoveFromDeadLetter removes a specific entry from the dead letter queue
func (h *Handler) RemoveFromDeadLetter(callID string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i, entry := range h.deadLetter {
		if entry.Call.ID == callID {
			h.deadLetter = append(h.deadLetter[:i], h.deadLetter[i+1:]...)
			return true
		}
	}
	return false
}

// RequeueFromDeadLetter moves an entry from dead letter back to retry queue
func (h *Handler) RequeueFromDeadLetter(callID string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i, entry := range h.deadLetter {
		if entry.Call.ID == callID {
			// Reset retry count and requeue
			entry.Call.RetryCount = 0
			entry.Call.Status = dispatcher.ServiceCallStatusPending

			retryable := &RetryableCall{
				Call:       entry.Call,
				RetryCount: 0,
				NextRetry:  time.Now(),
				LastError:  "",
			}

			h.retryQueue = append(h.retryQueue, retryable)
			h.deadLetter = append(h.deadLetter[:i], h.deadLetter[i+1:]...)

			h.logger.WithField("call_id", callID).Info("Requeued from dead letter")
			return true
		}
	}
	return false
}

// StartRetryWorker starts a background worker that processes retries
func (h *Handler) StartRetryWorker(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				h.logger.Info("Retry worker stopped")
				return
			case <-ticker.C:
				processed := h.ProcessRetries(ctx)
				if processed > 0 {
					h.logger.WithField("processed", processed).Debug("Processed retries")
				}
			}
		}
	}()

	h.logger.WithField("interval", interval).Info("Retry worker started")
}
