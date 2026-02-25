package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/downlink/dispatcher"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, time.Second, config.InitialDelay)
	assert.Equal(t, 30*time.Second, config.MaxDelay)
	assert.Equal(t, 2.0, config.Multiplier)
	assert.True(t, config.EnableDeadLetter)
}

func TestNewHandler(t *testing.T) {
	t.Run("with config", func(t *testing.T) {
		config := &Config{
			MaxRetries:   5,
			InitialDelay: 2 * time.Second,
		}
		handler := NewHandler(config, nil)

		require.NotNil(t, handler)
		assert.Equal(t, 5, handler.config.MaxRetries)
	})

	t.Run("without config", func(t *testing.T) {
		handler := NewHandler(nil, nil)

		require.NotNil(t, handler)
		assert.Equal(t, 3, handler.config.MaxRetries)
	})
}

func TestHandler_ScheduleRetry(t *testing.T) {
	config := &Config{
		MaxRetries:       3,
		InitialDelay:     100 * time.Millisecond,
		MaxDelay:         time.Second,
		Multiplier:       2.0,
		EnableDeadLetter: true,
	}
	handler := NewHandler(config, nil)

	t.Run("schedule first retry", func(t *testing.T) {
		call := &dispatcher.ServiceCall{
			ID:         "call-001",
			DeviceSN:   "DEVICE001",
			Vendor:     "dji",
			Method:     "takeoff",
			RetryCount: 0,
		}

		scheduled := handler.ScheduleRetry(call, "connection error")

		assert.True(t, scheduled)
		assert.Equal(t, 1, handler.GetPendingRetries())
	})

	t.Run("max retries exceeded", func(t *testing.T) {
		call := &dispatcher.ServiceCall{
			ID:         "call-002",
			DeviceSN:   "DEVICE002",
			Vendor:     "dji",
			Method:     "land",
			RetryCount: 3, // Already at max
		}

		scheduled := handler.ScheduleRetry(call, "timeout")

		assert.False(t, scheduled)
		assert.Equal(t, 1, handler.GetDeadLetterCount())
	})
}

func TestHandler_CalculateDelay(t *testing.T) {
	config := &Config{
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     time.Second,
		Multiplier:   2.0,
	}
	handler := NewHandler(config, nil)

	testCases := []struct {
		retryCount int
		expected   time.Duration
	}{
		{0, 100 * time.Millisecond},  // 100ms * 2^0 = 100ms
		{1, 200 * time.Millisecond},  // 100ms * 2^1 = 200ms
		{2, 400 * time.Millisecond},  // 100ms * 2^2 = 400ms
		{3, 800 * time.Millisecond},  // 100ms * 2^3 = 800ms
		{4, time.Second},             // 100ms * 2^4 = 1600ms, capped at 1s
		{10, time.Second},            // Capped at max
	}

	for _, tc := range testCases {
		delay := handler.calculateDelay(tc.retryCount)
		assert.Equal(t, tc.expected, delay, "retry count %d", tc.retryCount)
	}
}

func TestHandler_ProcessRetries(t *testing.T) {
	config := &Config{
		MaxRetries:       3,
		InitialDelay:     1 * time.Millisecond, // Very short for testing
		MaxDelay:         10 * time.Millisecond,
		Multiplier:       2.0,
		EnableDeadLetter: true,
	}
	handler := NewHandler(config, nil)

	t.Run("process due retries", func(t *testing.T) {
		call := &dispatcher.ServiceCall{
			ID:         "call-003",
			DeviceSN:   "DEVICE003",
			Vendor:     "dji",
			Method:     "takeoff",
			RetryCount: 0,
		}

		var retriedCall *dispatcher.ServiceCall
		handler.SetOnRetry(func(ctx context.Context, c *dispatcher.ServiceCall) error {
			retriedCall = c
			return nil
		})

		handler.ScheduleRetry(call, "error")

		// Wait for retry to be due
		time.Sleep(5 * time.Millisecond)

		processed := handler.ProcessRetries(context.Background())

		assert.Equal(t, 1, processed)
		assert.NotNil(t, retriedCall)
		assert.Equal(t, "call-003", retriedCall.ID)
		assert.Equal(t, 1, retriedCall.RetryCount)
	})

	t.Run("retry fails and reschedules", func(t *testing.T) {
		handler2 := NewHandler(config, nil)

		call := &dispatcher.ServiceCall{
			ID:         "call-004",
			DeviceSN:   "DEVICE004",
			Vendor:     "dji",
			Method:     "land",
			RetryCount: 0,
		}

		retryCount := 0
		handler2.SetOnRetry(func(ctx context.Context, c *dispatcher.ServiceCall) error {
			retryCount++
			return errors.New("still failing")
		})

		handler2.ScheduleRetry(call, "initial error")
		time.Sleep(5 * time.Millisecond)

		// First retry
		handler2.ProcessRetries(context.Background())
		assert.Equal(t, 1, retryCount)
		assert.Equal(t, 1, handler2.GetPendingRetries()) // Rescheduled

		time.Sleep(5 * time.Millisecond)

		// Second retry
		handler2.ProcessRetries(context.Background())
		assert.Equal(t, 2, retryCount)
		assert.Equal(t, 1, handler2.GetPendingRetries())

		time.Sleep(10 * time.Millisecond)

		// Third retry - should move to dead letter
		handler2.ProcessRetries(context.Background())
		assert.Equal(t, 3, retryCount)
		assert.Equal(t, 0, handler2.GetPendingRetries())
		assert.Equal(t, 1, handler2.GetDeadLetterCount())
	})
}

func TestHandler_DeadLetter(t *testing.T) {
	config := &Config{
		MaxRetries:       1,
		InitialDelay:     time.Millisecond,
		EnableDeadLetter: true,
	}
	handler := NewHandler(config, nil)

	t.Run("dead letter callback", func(t *testing.T) {
		var deadLetterEntry *DeadLetterEntry
		handler.SetOnDeadLetter(func(entry *DeadLetterEntry) {
			deadLetterEntry = entry
		})

		call := &dispatcher.ServiceCall{
			ID:         "call-005",
			DeviceSN:   "DEVICE005",
			Vendor:     "dji",
			Method:     "takeoff",
			RetryCount: 1, // At max
		}

		handler.ScheduleRetry(call, "final error")

		require.NotNil(t, deadLetterEntry)
		assert.Equal(t, "call-005", deadLetterEntry.Call.ID)
		assert.Equal(t, "final error", deadLetterEntry.Error)
	})

	t.Run("get dead letter entries", func(t *testing.T) {
		entries := handler.GetDeadLetterEntries()
		assert.Len(t, entries, 1)
	})

	t.Run("remove from dead letter", func(t *testing.T) {
		removed := handler.RemoveFromDeadLetter("call-005")
		assert.True(t, removed)
		assert.Equal(t, 0, handler.GetDeadLetterCount())
	})

	t.Run("remove non-existent", func(t *testing.T) {
		removed := handler.RemoveFromDeadLetter("non-existent")
		assert.False(t, removed)
	})
}

func TestHandler_RequeueFromDeadLetter(t *testing.T) {
	config := &Config{
		MaxRetries:       1,
		InitialDelay:     time.Millisecond,
		EnableDeadLetter: true,
	}
	handler := NewHandler(config, nil)

	call := &dispatcher.ServiceCall{
		ID:         "call-006",
		DeviceSN:   "DEVICE006",
		Vendor:     "dji",
		Method:     "takeoff",
		RetryCount: 1,
	}

	// Add to dead letter
	handler.ScheduleRetry(call, "error")
	assert.Equal(t, 1, handler.GetDeadLetterCount())
	assert.Equal(t, 0, handler.GetPendingRetries())

	// Requeue
	requeued := handler.RequeueFromDeadLetter("call-006")
	assert.True(t, requeued)
	assert.Equal(t, 0, handler.GetDeadLetterCount())
	assert.Equal(t, 1, handler.GetPendingRetries())

	// Verify call was reset
	handler.mu.RLock()
	retryable := handler.retryQueue[0]
	handler.mu.RUnlock()

	assert.Equal(t, 0, retryable.Call.RetryCount)
	assert.Equal(t, dispatcher.ServiceCallStatusPending, retryable.Call.Status)
}

func TestHandler_ClearDeadLetter(t *testing.T) {
	config := &Config{
		MaxRetries:       0,
		EnableDeadLetter: true,
	}
	handler := NewHandler(config, nil)

	// Add multiple entries
	for i := 0; i < 3; i++ {
		call := &dispatcher.ServiceCall{
			ID:       fmt.Sprintf("call-%c", 'a'+i),
			DeviceSN: "DEVICE",
		}
		handler.ScheduleRetry(call, "error")
	}

	assert.Equal(t, 3, handler.GetDeadLetterCount())

	handler.ClearDeadLetter()
	assert.Equal(t, 0, handler.GetDeadLetterCount())
}

func TestHandler_StartRetryWorker(t *testing.T) {
	config := &Config{
		MaxRetries:   3,
		InitialDelay: time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}
	handler := NewHandler(config, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	retried := make(chan string, 1)
	handler.SetOnRetry(func(ctx context.Context, c *dispatcher.ServiceCall) error {
		retried <- c.ID
		return nil
	})

	call := &dispatcher.ServiceCall{
		ID:         "call-worker",
		DeviceSN:   "DEVICE",
		Vendor:     "dji",
		Method:     "takeoff",
		RetryCount: 0,
	}

	handler.ScheduleRetry(call, "error")
	handler.StartRetryWorker(ctx, 5*time.Millisecond)

	select {
	case id := <-retried:
		assert.Equal(t, "call-worker", id)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("retry worker did not process retry")
	}
}

func TestRetryableCall(t *testing.T) {
	call := &dispatcher.ServiceCall{
		ID:       "call-001",
		DeviceSN: "DEVICE001",
	}

	retryable := &RetryableCall{
		Call:       call,
		RetryCount: 2,
		NextRetry:  time.Now().Add(time.Second),
		LastError:  "connection timeout",
	}

	assert.Equal(t, "call-001", retryable.Call.ID)
	assert.Equal(t, 2, retryable.RetryCount)
	assert.Equal(t, "connection timeout", retryable.LastError)
}

func TestDeadLetterEntry(t *testing.T) {
	call := &dispatcher.ServiceCall{
		ID:       "call-001",
		DeviceSN: "DEVICE001",
	}

	entry := &DeadLetterEntry{
		Call:     call,
		Error:    "max retries exceeded",
		FailedAt: time.Now(),
		Retries:  3,
	}

	assert.Equal(t, "call-001", entry.Call.ID)
	assert.Equal(t, "max retries exceeded", entry.Error)
	assert.Equal(t, 3, entry.Retries)
}
