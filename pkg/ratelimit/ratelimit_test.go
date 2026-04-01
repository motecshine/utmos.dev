package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenLimiter_Allow(t *testing.T) {
	// Create a limiter with capacity 5 and refill rate 2 per second
	limiter := NewTokenLimiter(5, 2)

	// Should allow up to 5 tokens immediately
	for i := 0; i < 5; i++ {
		assert.True(t, limiter.Allow(1), "Should allow token %d", i+1)
	}

	// 6th should fail as bucket is exhausted
	assert.False(t, limiter.Allow(1), "Should not allow 6th token when bucket is empty")
}

func TestTokenLimiter_Refill(t *testing.T) {
	// Create a limiter with capacity 2 and refill rate 10 per second
	limiter := NewTokenLimiter(2, 10)

	// Exhaust the bucket
	assert.True(t, limiter.Allow(1))
	assert.True(t, limiter.Allow(1))
	assert.False(t, limiter.Allow(1))

	// Wait for refill (100ms * 10 tokens/s = 1 token)
	time.Sleep(150 * time.Millisecond)

	// Should have refilled at least 1 token
	assert.True(t, limiter.Allow(1), "Should allow after refill")
}

func TestTokenLimiter_Capacity(t *testing.T) {
	limiter := NewTokenLimiter(10, 5)
	assert.Equal(t, float64(10), limiter.Capacity())
	assert.Equal(t, float64(5), limiter.RefillRate())
}

func TestTokenLimiter_BucketLevel(t *testing.T) {
	limiter := NewTokenLimiter(5, 10)

	// Initial level should be at capacity
	assert.Equal(t, float64(5), limiter.BucketLevel())

	// Use some tokens
	limiter.Allow(2)
	// Allow small floating point tolerance
	level := limiter.BucketLevel()
	assert.InDelta(t, float64(3), level, 0.001, "Bucket level should be approximately 3")
}

func TestConnectionLimiter_Allow(t *testing.T) {
	// Global: 3 per second with burst 3, Per-client: 2 per second with burst 2
	cl := NewConnectionLimiter(3, 3, 2, 2)

	// First 3 connections should be allowed (global burst)
	assert.True(t, cl.Allow("client1"), "First connection should be allowed")
	assert.True(t, cl.Allow("client2"), "Second connection should be allowed")
	assert.True(t, cl.Allow("client3"), "Third connection should be allowed")

	// 4th should fail global limit
	assert.False(t, cl.Allow("client4"), "4th connection should be blocked by global limit")
}

func TestConnectionLimiter_PerClientLimit(t *testing.T) {
	// Global: 10 per second with burst 10, Per-client: 1 per second with burst 2
	cl := NewConnectionLimiter(10, 10, 1, 2)

	// client1 can connect (first burst of 2)
	assert.True(t, cl.Allow("client1"), "First connection should be allowed")
	assert.True(t, cl.Allow("client1"), "Second connection should be allowed (burst)")

	// client2 can also connect
	assert.True(t, cl.Allow("client2"), "client2 first connection should be allowed")

	// client1 is now limited by per-client rate (needs to wait for refill)
	// But since burst is 2, we need to wait for refill
	time.Sleep(1100 * time.Millisecond) // Wait for 1+ token to refill

	assert.True(t, cl.Allow("client1"), "client1 should be allowed after refill")
}

func TestConnectionLimiter_RemoveClient(t *testing.T) {
	cl := NewConnectionLimiter(10, 10, 1, 2)

	assert.Equal(t, 0, cl.ClientCount())

	cl.Allow("client1")
	cl.Allow("client1")

	assert.Equal(t, 1, cl.ClientCount())

	cl.RemoveClient("client1")

	assert.Equal(t, 0, cl.ClientCount())
}
