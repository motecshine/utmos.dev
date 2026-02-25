package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/uplink/processor"
	"github.com/utmos/utmos/internal/ws"
	"github.com/utmos/utmos/internal/ws/hub"
	djiuplink "github.com/utmos/utmos/pkg/adapter/dji/uplink"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// TestPerformanceMessageLatency tests message processing latency (NFR-001)
// Target: P95 latency < 100ms
func TestPerformanceMessageLatency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	t.Run("uplink message processing latency", func(t *testing.T) {
		registry := processor.NewRegistry(nil)
		djiProcessor := djiuplink.NewProcessorAdapter(nil)
		registry.Register(djiProcessor)
		handler := processor.NewMessageHandler(registry, nil)

		messageCount := 1000
		latencies := make([]time.Duration, messageCount)

		for i := 0; i < messageCount; i++ {
			msg := &rabbitmq.StandardMessage{
				TID:       fmt.Sprintf("perf-tid-%d", i),
				BID:       fmt.Sprintf("perf-bid-%d", i),
				Service:   "iot-gateway",
				Action:    "telemetry.report",
				DeviceSN:  "PERF-DRONE-001",
				Timestamp: time.Now().UnixMilli(),
				Data:      json.RawMessage(`{"latitude": 39.9042, "longitude": 116.4074, "altitude": 100.5}`),
				ProtocolMeta: &rabbitmq.ProtocolMeta{
					Vendor: "dji",
				},
			}

			start := time.Now()
			_ = handler.Handle(context.Background(), msg)
			latencies[i] = time.Since(start)
		}

		// Calculate P95 latency
		p95 := calculateP95(latencies)
		t.Logf("P95 latency: %v", p95)

		// NFR-001: P95 latency < 100ms
		assert.Less(t, p95, 100*time.Millisecond, "P95 latency should be less than 100ms")
	})

	t.Run("websocket push latency", func(t *testing.T) {
		wsSvc := ws.NewService(nil, nil, nil)
		err := wsSvc.Start(context.Background())
		require.NoError(t, err)
		defer func() { _ = wsSvc.Stop() }()

		// Subscribe mock clients
		for i := 0; i < 100; i++ {
			clientID := fmt.Sprintf("perf-client-%d", i)
			wsSvc.SubscriptionManager().Subscribe(clientID, "perf.topic")
		}

		messageCount := 1000
		latencies := make([]time.Duration, messageCount)

		for i := 0; i < messageCount; i++ {
			msg := &hub.Message{
				Type:  hub.MessageTypeEvent,
				Event: "perf.topic",
				Data:  map[string]interface{}{"index": i},
			}

			start := time.Now()
			wsSvc.Pusher().PushToTopic("perf.topic", msg)
			latencies[i] = time.Since(start)
		}

		// Wait for all messages to be processed
		time.Sleep(500 * time.Millisecond)

		// Calculate P95 latency
		p95 := calculateP95(latencies)
		t.Logf("WebSocket push P95 latency: %v", p95)

		// P95 latency should be reasonable for push operations
		assert.Less(t, p95, 10*time.Millisecond, "Push P95 latency should be less than 10ms")
	})
}

// TestPerformanceThroughput tests message throughput
func TestPerformanceThroughput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	t.Run("message processing throughput", func(t *testing.T) {
		registry := processor.NewRegistry(nil)
		djiProcessor := djiuplink.NewProcessorAdapter(nil)
		registry.Register(djiProcessor)
		handler := processor.NewMessageHandler(registry, nil)

		messageCount := 10000
		start := time.Now()

		for i := 0; i < messageCount; i++ {
			msg := &rabbitmq.StandardMessage{
				TID:       "throughput-tid",
				BID:       "throughput-bid",
				Service:   "iot-gateway",
				Action:    "telemetry.report",
				DeviceSN:  "THROUGHPUT-DRONE-001",
				Timestamp: time.Now().UnixMilli(),
				Data:      json.RawMessage(`{"latitude": 39.9042, "longitude": 116.4074}`),
				ProtocolMeta: &rabbitmq.ProtocolMeta{
					Vendor: "dji",
				},
			}
			_ = handler.Handle(context.Background(), msg)
		}

		elapsed := time.Since(start)
		throughput := float64(messageCount) / elapsed.Seconds()

		t.Logf("Processed %d messages in %v", messageCount, elapsed)
		t.Logf("Throughput: %.2f messages/second", throughput)

		// Should process at least 1000 messages per second
		assert.Greater(t, throughput, 1000.0, "Throughput should be at least 1000 msg/s")
	})

	t.Run("concurrent message processing", func(t *testing.T) {
		registry := processor.NewRegistry(nil)
		djiProcessor := djiuplink.NewProcessorAdapter(nil)
		registry.Register(djiProcessor)
		handler := processor.NewMessageHandler(registry, nil)

		messageCount := 10000
		workerCount := 10
		messagesPerWorker := messageCount / workerCount

		var wg sync.WaitGroup
		start := time.Now()

		for w := 0; w < workerCount; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				for i := 0; i < messagesPerWorker; i++ {
					msg := &rabbitmq.StandardMessage{
						TID:       "concurrent-tid",
						BID:       "concurrent-bid",
						Service:   "iot-gateway",
						Action:    "telemetry.report",
						DeviceSN:  "CONCURRENT-DRONE-001",
						Timestamp: time.Now().UnixMilli(),
						Data:      json.RawMessage(fmt.Sprintf(`{"worker": %d}`, workerID)),
						ProtocolMeta: &rabbitmq.ProtocolMeta{
							Vendor: "dji",
						},
					}
					_ = handler.Handle(context.Background(), msg)
				}
			}(w)
		}

		wg.Wait()
		elapsed := time.Since(start)
		throughput := float64(messageCount) / elapsed.Seconds()

		t.Logf("Concurrent processing: %d messages in %v", messageCount, elapsed)
		t.Logf("Concurrent throughput: %.2f messages/second", throughput)

		// Concurrent processing should be faster
		assert.Greater(t, throughput, 5000.0, "Concurrent throughput should be at least 5000 msg/s")
	})
}

// TestPerformanceMemory tests memory usage under load
func TestPerformanceMemory(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	t.Run("websocket service memory", func(t *testing.T) {
		wsSvc := ws.NewService(nil, nil, nil)
		err := wsSvc.Start(context.Background())
		require.NoError(t, err)
		defer func() { _ = wsSvc.Stop() }()

		// Add many subscriptions
		for i := 0; i < 1000; i++ {
			clientID := fmt.Sprintf("mem-client-%d", i)
			for j := 0; j < 10; j++ {
				topic := fmt.Sprintf("mem.topic.%d", j)
				wsSvc.SubscriptionManager().Subscribe(clientID, topic)
			}
		}

		stats := wsSvc.GetStats()
		t.Logf("Active topics: %d", stats.ActiveTopics)

		// Should handle 10000 subscriptions
		assert.Equal(t, 10, stats.ActiveTopics)
	})
}

// calculateP95 calculates the 95th percentile latency
func calculateP95(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	// Sort latencies
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Calculate P95 index
	p95Index := int(float64(len(sorted)) * 0.95)
	if p95Index >= len(sorted) {
		p95Index = len(sorted) - 1
	}

	return sorted[p95Index]
}

// BenchmarkMessageProcessing benchmarks message processing
func BenchmarkMessageProcessing(b *testing.B) {
	registry := processor.NewRegistry(nil)
	djiProcessor := djiuplink.NewProcessorAdapter(nil)
	registry.Register(djiProcessor)
	handler := processor.NewMessageHandler(registry, nil)

	msg := &rabbitmq.StandardMessage{
		TID:       "bench-tid",
		BID:       "bench-bid",
		Service:   "iot-gateway",
		Action:    "telemetry.report",
		DeviceSN:  "BENCH-DRONE-001",
		Timestamp: time.Now().UnixMilli(),
		Data:      json.RawMessage(`{"latitude": 39.9042, "longitude": 116.4074}`),
		ProtocolMeta: &rabbitmq.ProtocolMeta{
			Vendor: "dji",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler.Handle(context.Background(), msg)
	}
}

// BenchmarkWebSocketPush benchmarks WebSocket push operations
func BenchmarkWebSocketPush(b *testing.B) {
	wsSvc := ws.NewService(nil, nil, nil)
	_ = wsSvc.Start(context.Background())
	defer func() { _ = wsSvc.Stop() }()

	// Subscribe clients
	for i := 0; i < 100; i++ {
		wsSvc.SubscriptionManager().Subscribe(fmt.Sprintf("bench-client-%d", i), "bench.topic")
	}

	msg := &hub.Message{
		Type:  hub.MessageTypeEvent,
		Event: "bench.topic",
		Data:  map[string]interface{}{"test": "data"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wsSvc.Pusher().PushToTopic("bench.topic", msg)
	}
}
