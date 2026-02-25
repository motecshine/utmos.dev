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

	"github.com/utmos/utmos/internal/gateway/connection"
	"github.com/utmos/utmos/internal/uplink/processor"
	djiuplink "github.com/utmos/utmos/pkg/adapter/dji/uplink"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// TestLoadDeviceConnections tests support for 1000+ simultaneous device connections (NFR-002)
func TestLoadDeviceConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	t.Run("1000 device connections", func(t *testing.T) {
		manager := connection.NewManager(nil)

		deviceCount := 1000
		var wg sync.WaitGroup
		var successCount int64

		start := time.Now()

		// Simulate 1000 device connections
		for i := 0; i < deviceCount; i++ {
			wg.Add(1)
			go func(deviceID int) {
				defer wg.Done()

				deviceSN := fmt.Sprintf("LOAD-DEVICE-%04d", deviceID)
				clientID := fmt.Sprintf("client-%04d", deviceID)

				// Register device connection
				manager.Connect(deviceSN, clientID, "127.0.0.1")
				atomic.AddInt64(&successCount, 1)
			}(i)
		}

		wg.Wait()
		elapsed := time.Since(start)

		t.Logf("Registered %d devices in %v", successCount, elapsed)
		t.Logf("Registration rate: %.2f devices/second", float64(successCount)/elapsed.Seconds())

		// Verify all devices are connected
		assert.Equal(t, int64(deviceCount), successCount)
		assert.Equal(t, deviceCount, manager.GetOnlineCount())

		// Verify we can get device status
		for i := 0; i < 10; i++ {
			deviceSN := fmt.Sprintf("LOAD-DEVICE-%04d", i)
			assert.True(t, manager.IsOnline(deviceSN))
		}
	})

	t.Run("device connection churn", func(t *testing.T) {
		manager := connection.NewManager(nil)

		deviceCount := 500
		iterations := 3
		var totalOps int64

		start := time.Now()

		for iter := 0; iter < iterations; iter++ {
			var wg sync.WaitGroup

			// Connect devices
			for i := 0; i < deviceCount; i++ {
				wg.Add(1)
				go func(deviceID int) {
					defer wg.Done()
					deviceSN := fmt.Sprintf("CHURN-DEVICE-%04d", deviceID)
					clientID := fmt.Sprintf("client-%04d-%d", deviceID, iter)
					manager.Connect(deviceSN, clientID, "127.0.0.1")
					atomic.AddInt64(&totalOps, 1)
				}(i)
			}
			wg.Wait()

			// Disconnect half the devices
			for i := 0; i < deviceCount/2; i++ {
				wg.Add(1)
				go func(deviceID int) {
					defer wg.Done()
					deviceSN := fmt.Sprintf("CHURN-DEVICE-%04d", deviceID)
					manager.Disconnect(deviceSN)
					atomic.AddInt64(&totalOps, 1)
				}(i)
			}
			wg.Wait()
		}

		elapsed := time.Since(start)
		t.Logf("Completed %d operations in %v", totalOps, elapsed)
		t.Logf("Operation rate: %.2f ops/second", float64(totalOps)/elapsed.Seconds())

		// Should handle high churn rate
		assert.Greater(t, float64(totalOps)/elapsed.Seconds(), 1000.0)
	})
}

// TestLoadDeviceMessages tests message processing under device load
func TestLoadDeviceMessages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	t.Run("messages from 1000 devices", func(t *testing.T) {
		registry := processor.NewRegistry(nil)
		djiProcessor := djiuplink.NewProcessorAdapter(nil)
		registry.Register(djiProcessor)
		handler := processor.NewMessageHandler(registry, nil)

		deviceCount := 1000
		messagesPerDevice := 10
		totalMessages := deviceCount * messagesPerDevice

		var wg sync.WaitGroup
		var processedCount int64
		var errorCount int64

		start := time.Now()

		for d := 0; d < deviceCount; d++ {
			wg.Add(1)
			go func(deviceID int) {
				defer wg.Done()

				deviceSN := fmt.Sprintf("MSG-DEVICE-%04d", deviceID)

				for m := 0; m < messagesPerDevice; m++ {
					msg := &rabbitmq.StandardMessage{
						TID:       fmt.Sprintf("tid-%d-%d", deviceID, m),
						BID:       fmt.Sprintf("bid-%d-%d", deviceID, m),
						Service:   "iot-gateway",
						Action:    "telemetry.report",
						DeviceSN:  deviceSN,
						Timestamp: time.Now().UnixMilli(),
						Data:      json.RawMessage(`{"latitude": 39.9042, "longitude": 116.4074}`),
						ProtocolMeta: &rabbitmq.ProtocolMeta{
							Vendor: "dji",
						},
					}

					err := handler.Handle(context.Background(), msg)
					if err != nil {
						atomic.AddInt64(&errorCount, 1)
					} else {
						atomic.AddInt64(&processedCount, 1)
					}
				}
			}(d)
		}

		wg.Wait()
		elapsed := time.Since(start)

		t.Logf("Processed %d messages from %d devices in %v", processedCount, deviceCount, elapsed)
		t.Logf("Throughput: %.2f messages/second", float64(processedCount)/elapsed.Seconds())
		t.Logf("Errors: %d", errorCount)

		// All messages should be processed
		assert.Equal(t, int64(totalMessages), processedCount)
		assert.Equal(t, int64(0), errorCount)

		// Should process at least 5000 messages per second
		assert.Greater(t, float64(processedCount)/elapsed.Seconds(), 5000.0)
	})
}

// TestLoadDeviceStatusUpdates tests device status update performance
func TestLoadDeviceStatusUpdates(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	manager := connection.NewManager(nil)

	deviceCount := 1000

	// Register all devices first
	for i := 0; i < deviceCount; i++ {
		deviceSN := fmt.Sprintf("STATUS-DEVICE-%04d", i)
		clientID := fmt.Sprintf("client-%04d", i)
		manager.Connect(deviceSN, clientID, "127.0.0.1")
	}

	t.Run("concurrent status queries", func(t *testing.T) {
		queryCount := 10000
		var wg sync.WaitGroup
		var successCount int64

		start := time.Now()

		for i := 0; i < queryCount; i++ {
			wg.Add(1)
			go func(queryID int) {
				defer wg.Done()
				deviceSN := fmt.Sprintf("STATUS-DEVICE-%04d", queryID%deviceCount)
				if manager.IsOnline(deviceSN) {
					atomic.AddInt64(&successCount, 1)
				}
			}(i)
		}

		wg.Wait()
		elapsed := time.Since(start)

		t.Logf("Completed %d status queries in %v", queryCount, elapsed)
		t.Logf("Query rate: %.2f queries/second", float64(queryCount)/elapsed.Seconds())

		assert.Equal(t, int64(queryCount), successCount)
		// Should handle at least 10000 queries per second
		assert.Greater(t, float64(queryCount)/elapsed.Seconds(), 10000.0)
	})

	t.Run("concurrent status updates", func(t *testing.T) {
		updateCount := 5000
		var wg sync.WaitGroup

		start := time.Now()

		for i := 0; i < updateCount; i++ {
			wg.Add(1)
			go func(updateID int) {
				defer wg.Done()
				deviceSN := fmt.Sprintf("STATUS-DEVICE-%04d", updateID%deviceCount)
				// Simulate heartbeat update
				manager.UpdateLastSeen(deviceSN)
			}(i)
		}

		wg.Wait()
		elapsed := time.Since(start)

		t.Logf("Completed %d status updates in %v", updateCount, elapsed)
		t.Logf("Update rate: %.2f updates/second", float64(updateCount)/elapsed.Seconds())

		// Should handle at least 5000 updates per second
		assert.Greater(t, float64(updateCount)/elapsed.Seconds(), 5000.0)
	})
}

// TestLoadDeviceReconnection tests device reconnection handling
func TestLoadDeviceReconnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	manager := connection.NewManager(nil)

	deviceCount := 100
	reconnectCount := 10

	var wg sync.WaitGroup
	var totalReconnects int64

	start := time.Now()

	for d := 0; d < deviceCount; d++ {
		wg.Add(1)
		go func(deviceID int) {
			defer wg.Done()

			deviceSN := fmt.Sprintf("RECONNECT-DEVICE-%04d", deviceID)

			for r := 0; r < reconnectCount; r++ {
				clientID := fmt.Sprintf("client-%04d-%d", deviceID, r)

				// Connect
				manager.Connect(deviceSN, clientID, "127.0.0.1")

				// Small delay
				time.Sleep(time.Millisecond)

				// Disconnect
				manager.Disconnect(deviceSN)

				atomic.AddInt64(&totalReconnects, 1)
			}
		}(d)
	}

	wg.Wait()
	elapsed := time.Since(start)

	expectedReconnects := int64(deviceCount * reconnectCount)
	t.Logf("Completed %d reconnections in %v", totalReconnects, elapsed)
	t.Logf("Reconnection rate: %.2f reconnects/second", float64(totalReconnects)/elapsed.Seconds())

	assert.Equal(t, expectedReconnects, totalReconnects)
}

// BenchmarkDeviceConnection benchmarks device connection operations
func BenchmarkDeviceConnection(b *testing.B) {
	manager := connection.NewManager(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deviceSN := fmt.Sprintf("BENCH-DEVICE-%d", i)
		clientID := fmt.Sprintf("client-%d", i)
		manager.Connect(deviceSN, clientID, "127.0.0.1")
	}
}

// BenchmarkDeviceStatusQuery benchmarks device status queries
func BenchmarkDeviceStatusQuery(b *testing.B) {
	manager := connection.NewManager(nil)

	// Pre-register devices
	for i := 0; i < 1000; i++ {
		deviceSN := fmt.Sprintf("BENCH-DEVICE-%d", i)
		clientID := fmt.Sprintf("client-%d", i)
		manager.Connect(deviceSN, clientID, "127.0.0.1")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deviceSN := fmt.Sprintf("BENCH-DEVICE-%d", i%1000)
		_ = manager.IsOnline(deviceSN)
	}
}
