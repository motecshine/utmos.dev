// Package ratelimit provides token bucket rate limiting implementations.
package ratelimit

import (
	"sync"
	"time"
)

// TokenLimiter implements a token bucket rate limiter.
type TokenLimiter struct {
	mu          sync.Mutex
	bucket      float64
	capacity    float64
	refillRate  float64 // tokens per second
	lastRefill  time.Time
}

// NewTokenLimiter creates a new token bucket limiter.
// capacity is the maximum number of tokens in the bucket.
// refillRate is the number of tokens added per second.
func NewTokenLimiter(capacity float64, refillRate float64) *TokenLimiter {
	now := time.Now()
	return &TokenLimiter{
		bucket:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: now,
	}
}

// Allow reports whether n tokens can be consumed immediately.
func (l *TokenLimiter) Allow(n float64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.refill()

	if l.bucket >= n {
		l.bucket -= n
		return true
	}
	return false
}

// Wait blocks until n tokens become available.
func (l *TokenLimiter) Wait(n float64) bool {
	l.mu.Lock()

	// Refill first
	l.refill()

	for l.bucket < n {
		// Calculate wait time for needed tokens
		needed := n - l.bucket
		waitTime := time.Duration(needed/l.refillRate * float64(time.Second))

		// Release lock while waiting
		l.mu.Unlock()
		time.Sleep(waitTime)
		l.mu.Lock()
		l.refill()
	}

	l.bucket -= n
	l.mu.Unlock()
	return true
}

// Capacity returns the bucket capacity.
func (l *TokenLimiter) Capacity() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.capacity
}

// RefillRate returns the refill rate in tokens per second.
func (l *TokenLimiter) RefillRate() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.refillRate
}

// refill adds tokens based on elapsed time since last refill.
func (l *TokenLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(l.lastRefill).Seconds()
	tokensToAdd := elapsed * l.refillRate

	l.bucket += tokensToAdd
	if l.bucket > l.capacity {
		l.bucket = l.capacity
	}

	l.lastRefill = now
}

// BucketLevel returns the current number of tokens in the bucket.
func (l *TokenLimiter) BucketLevel() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refill()
	return l.bucket
}

// ConnectionLimiter provides rate limiting for MQTT connections.
type ConnectionLimiter struct {
	limiter      *TokenLimiter
	clientLimits map[string]*TokenLimiter // per-client limiters
	mu           sync.RWMutex
	defaultRate  float64
	defaultBurst float64
}

// NewConnectionLimiter creates a connection limiter with global and per-client limits.
// globalRate is the maximum connections per second globally.
// globalBurst is the maximum burst of connections allowed.
// perClientRate is the maximum connections per second per client.
// perClientBurst is the maximum burst per client.
func NewConnectionLimiter(globalRate, globalBurst, perClientRate, perClientBurst float64) *ConnectionLimiter {
	return &ConnectionLimiter{
		limiter:      NewTokenLimiter(globalBurst, globalRate),
		clientLimits: make(map[string]*TokenLimiter),
		defaultRate:  perClientRate,
		defaultBurst: perClientBurst,
	}
}

// Allow checks if a connection from the given client ID should be allowed.
func (cl *ConnectionLimiter) Allow(clientID string) bool {
	// Check global limit first
	if !cl.limiter.Allow(1) {
		return false
	}

	// Check per-client limit
	cl.mu.RLock()
	clientLimiter, exists := cl.clientLimits[clientID]
	cl.mu.RUnlock()

	if !exists {
		cl.mu.Lock()
		// Double-check after acquiring write lock
		if clientLimiter, exists = cl.clientLimits[clientID]; !exists {
			clientLimiter = NewTokenLimiter(cl.defaultBurst, cl.defaultRate)
			cl.clientLimits[clientID] = clientLimiter
		}
		cl.mu.Unlock()
	}

	return clientLimiter.Allow(1)
}

// RemoveClient removes a client's rate limiter.
func (cl *ConnectionLimiter) RemoveClient(clientID string) {
	cl.mu.Lock()
	delete(cl.clientLimits, clientID)
	cl.mu.Unlock()
}

// ClientCount returns the number of tracked client limiters.
func (cl *ConnectionLimiter) ClientCount() int {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	return len(cl.clientLimits)
}
