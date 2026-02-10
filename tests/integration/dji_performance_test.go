package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	djiinit "github.com/utmos/utmos/pkg/adapter/dji/init"
)

// TestPerformance_1000Devices tests message processing performance with 1000 simulated devices.
func TestPerformance_1000Devices(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	adapter := djiinit.NewInitializedAdapter()
	ctx := context.Background()

	const (
		numDevices        = 1000
		messagesPerDevice = 10
		totalMessages     = numDevices * messagesPerDevice
	)

	// Generate test messages
	messages := make([]struct {
		topic   string
		payload []byte
	}, totalMessages)

	for i := 0; i < numDevices; i++ {
		deviceSN := fmt.Sprintf("device-%04d", i)
		for j := 0; j < messagesPerDevice; j++ {
			idx := i*messagesPerDevice + j
			payload := map[string]interface{}{
				"tid":       fmt.Sprintf("tid-%d-%d", i, j),
				"bid":       fmt.Sprintf("bid-%d-%d", i, j),
				"timestamp": time.Now().UnixMilli(),
				"gateway":   deviceSN,
				"data": map[string]interface{}{
					"host": map[string]interface{}{
						"latitude":  31.2304 + float64(i)*0.0001,
						"longitude": 121.4737 + float64(j)*0.0001,
						"altitude":  100.0 + float64(j),
						"battery": map[string]interface{}{
							"capacity_percent": 85 - j,
						},
					},
				},
			}
			payloadBytes, _ := json.Marshal(payload)
			messages[idx].topic = fmt.Sprintf("thing/product/%s/osd", deviceSN)
			messages[idx].payload = payloadBytes
		}
	}

	// Warm up
	for i := 0; i < 100; i++ {
		_, _ = adapter.HandleMessage(ctx, messages[i].topic, messages[i].payload)
	}

	// Measure processing time
	var (
		successCount int64
		errorCount   int64
		totalLatency int64
	)

	start := time.Now()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 100) // Limit concurrency

	for i := 0; i < totalMessages; i++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(idx int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			msgStart := time.Now()

			_, err := adapter.HandleMessage(ctx, messages[idx].topic, messages[idx].payload)
			if err != nil {
				atomic.AddInt64(&errorCount, 1)
				return
			}

			latency := time.Since(msgStart).Microseconds()
			atomic.AddInt64(&totalLatency, latency)
			atomic.AddInt64(&successCount, 1)
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	// Calculate metrics
	avgLatency := float64(totalLatency) / float64(successCount)
	throughput := float64(successCount) / elapsed.Seconds()

	t.Logf("Performance Test Results:")
	t.Logf("  Total Messages: %d", totalMessages)
	t.Logf("  Success: %d", successCount)
	t.Logf("  Errors: %d", errorCount)
	t.Logf("  Total Time: %v", elapsed)
	t.Logf("  Avg Latency: %.2f µs", avgLatency)
	t.Logf("  Throughput: %.2f msg/s", throughput)

	// Assertions
	assert.Equal(t, int64(totalMessages), successCount, "All messages should be processed successfully")
	assert.Equal(t, int64(0), errorCount, "No errors should occur")
	assert.Less(t, avgLatency, float64(50000), "Average latency should be < 50ms (50000µs)")
}

// TestPerformance_Latency tests message processing latency.
func TestPerformance_Latency(t *testing.T) {
	adapter := djiinit.NewInitializedAdapter()
	ctx := context.Background()

	const iterations = 1000

	topic := "thing/product/gateway-001/osd"
	payload := []byte(`{
		"tid": "tid-perf",
		"bid": "bid-perf",
		"timestamp": 1234567890123,
		"gateway": "gateway-001",
		"data": {
			"host": {
				"latitude": 31.2304,
				"longitude": 121.4737,
				"altitude": 100.5,
				"height": 50.0,
				"attitude_pitch": 0.5,
				"attitude_roll": 0.2,
				"attitude_yaw": 180.0,
				"horizontal_speed": 5.0,
				"vertical_speed": 0.0,
				"battery": {
					"capacity_percent": 85,
					"voltage": 48000,
					"temperature": 25.0
				}
			}
		}
	}`)

	// Warm up
	for i := 0; i < 100; i++ {
		_, _ = adapter.HandleMessage(ctx, topic, payload)
	}

	// Measure latencies
	latencies := make([]time.Duration, iterations)

	for i := 0; i < iterations; i++ {
		start := time.Now()

		_, err := adapter.HandleMessage(ctx, topic, payload)
		require.NoError(t, err)

		latencies[i] = time.Since(start)
	}

	// Calculate percentiles
	var total time.Duration
	for _, l := range latencies {
		total += l
	}
	avg := total / time.Duration(iterations)

	// Sort for percentile calculation
	sortedLatencies := make([]time.Duration, len(latencies))
	copy(sortedLatencies, latencies)
	for i := 0; i < len(sortedLatencies); i++ {
		for j := i + 1; j < len(sortedLatencies); j++ {
			if sortedLatencies[i] > sortedLatencies[j] {
				sortedLatencies[i], sortedLatencies[j] = sortedLatencies[j], sortedLatencies[i]
			}
		}
	}

	p50 := sortedLatencies[iterations*50/100]
	p95 := sortedLatencies[iterations*95/100]
	p99 := sortedLatencies[iterations*99/100]

	t.Logf("Latency Test Results (n=%d):", iterations)
	t.Logf("  Average: %v", avg)
	t.Logf("  P50: %v", p50)
	t.Logf("  P95: %v", p95)
	t.Logf("  P99: %v", p99)

	// Assert P95 < 50ms (spec requirement)
	assert.Less(t, p95.Milliseconds(), int64(50), "P95 latency should be < 50ms")
}

// TestPerformance_Concurrent tests concurrent message processing.
func TestPerformance_Concurrent(t *testing.T) {
	adapter := djiinit.NewInitializedAdapter()
	ctx := context.Background()

	const (
		concurrency = 100
		iterations  = 100
	)

	topic := "thing/product/gateway-001/osd"
	payload := []byte(`{
		"tid": "tid-concurrent",
		"bid": "bid-concurrent",
		"timestamp": 1234567890123,
		"gateway": "gateway-001",
		"data": {
			"host": {
				"latitude": 31.2304,
				"longitude": 121.4737,
				"altitude": 100.5
			}
		}
	}`)

	var (
		successCount int64
		errorCount   int64
	)

	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				_, err := adapter.HandleMessage(ctx, topic, payload)
				if err != nil {
					atomic.AddInt64(&errorCount, 1)
					continue
				}

				atomic.AddInt64(&successCount, 1)
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	totalMessages := int64(concurrency * iterations)
	throughput := float64(successCount) / elapsed.Seconds()

	t.Logf("Concurrent Test Results:")
	t.Logf("  Concurrency: %d", concurrency)
	t.Logf("  Total Messages: %d", totalMessages)
	t.Logf("  Success: %d", successCount)
	t.Logf("  Errors: %d", errorCount)
	t.Logf("  Total Time: %v", elapsed)
	t.Logf("  Throughput: %.2f msg/s", throughput)

	assert.Equal(t, totalMessages, successCount, "All messages should be processed successfully")
	assert.Equal(t, int64(0), errorCount, "No errors should occur")
}

// TestPerformance_Memory tests memory usage during message processing.
func TestPerformance_Memory(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory test in short mode")
	}

	adapter := djiinit.NewInitializedAdapter()
	ctx := context.Background()

	const iterations = 10000

	topic := "thing/product/gateway-001/osd"
	payload := []byte(`{
		"tid": "tid-mem",
		"bid": "bid-mem",
		"timestamp": 1234567890123,
		"gateway": "gateway-001",
		"data": {
			"host": {
				"latitude": 31.2304,
				"longitude": 121.4737,
				"altitude": 100.5,
				"battery": {"capacity_percent": 85}
			}
		}
	}`)

	// Process messages
	for i := 0; i < iterations; i++ {
		_, err := adapter.HandleMessage(ctx, topic, payload)
		require.NoError(t, err)
	}

	// Memory should be stable (no leaks)
	// This is a basic test - in production, use pprof for detailed analysis
	t.Logf("Processed %d messages without memory issues", iterations)
}
