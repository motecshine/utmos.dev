package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/utmos/utmos/internal/api"
	"github.com/utmos/utmos/internal/downlink/model"
	"github.com/utmos/utmos/internal/gateway/connection"
	"github.com/utmos/utmos/internal/uplink/processor"
	"github.com/utmos/utmos/internal/ws"
	djiuplink "github.com/utmos/utmos/pkg/adapter/dji/uplink"
	"github.com/utmos/utmos/pkg/models"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// TestAvailabilityServiceRecovery tests service recovery after failures (NFR-004)
// Target: Service availability > 99.9%
func TestAvailabilityServiceRecovery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping availability test in short mode")
	}

	t.Run("websocket service recovery", func(t *testing.T) {
		wsSvc := ws.NewService(nil, nil, nil)

		// Start service
		err := wsSvc.Start(context.Background())
		require.NoError(t, err)
		assert.True(t, wsSvc.IsRunning())

		// Stop service
		err = wsSvc.Stop()
		require.NoError(t, err)
		assert.False(t, wsSvc.IsRunning())

		// Restart service
		err = wsSvc.Start(context.Background())
		require.NoError(t, err)
		assert.True(t, wsSvc.IsRunning())

		// Verify service is functional
		wsSvc.SubscriptionManager().Subscribe("test-client", "test.topic")
		assert.True(t, wsSvc.SubscriptionManager().IsSubscribed("test-client", "test.topic"))

		_ = wsSvc.Stop()
	})

	t.Run("connection manager recovery", func(t *testing.T) {
		manager := connection.NewManager(nil)

		// Register some devices
		for i := 0; i < 100; i++ {
			manager.Connect(fmt.Sprintf("device-%d", i), fmt.Sprintf("client-%d", i), "127.0.0.1")
		}

		assert.Equal(t, 100, manager.GetOnlineCount())

		// Disconnect all devices
		for i := 0; i < 100; i++ {
			manager.Disconnect(fmt.Sprintf("device-%d", i))
		}

		assert.Equal(t, 0, manager.GetOnlineCount())

		// Should be able to register new devices
		manager.Connect("new-device", "new-client", "127.0.0.1")
		assert.Equal(t, 1, manager.GetOnlineCount())
	})

	t.Run("processor registry recovery", func(t *testing.T) {
		registry := processor.NewRegistry(nil)
		djiProcessor := djiuplink.NewProcessorAdapter(nil)
		registry.Register(djiProcessor)
		handler := processor.NewMessageHandler(registry, nil)

		// Process messages
		for i := 0; i < 100; i++ {
			msg := &rabbitmq.StandardMessage{
				TID:       "recovery-tid",
				BID:       "recovery-bid",
				Service:   "iot-gateway",
				Action:    "telemetry.report",
				DeviceSN:  "RECOVERY-DEVICE",
				Timestamp: time.Now().UnixMilli(),
				Data:      json.RawMessage(`{"test": "data"}`),
				ProtocolMeta: &rabbitmq.ProtocolMeta{
					Vendor: "dji",
				},
			}
			_ = handler.Handle(context.Background(), msg)
		}

		// Registry should still be functional
		_, found := registry.Get("dji")
		assert.True(t, found)
	})
}

// TestAvailabilityUnderLoad tests service availability under sustained load
func TestAvailabilityUnderLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping availability test in short mode")
	}

	t.Run("api availability under load", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		require.NoError(t, err)
		err = db.AutoMigrate(&models.Device{}, &model.ServiceCall{})
		require.NoError(t, err)

		config := &api.Config{
			EnableAuth:  false,
			EnableTrace: false,
		}
		router := api.NewRouter(config, db, nil, nil, nil)

		server := httptest.NewServer(router)
		defer server.Close()

		duration := 5 * time.Second
		requestCount := 0
		successCount := 0
		var mu sync.Mutex

		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		var wg sync.WaitGroup
		workerCount := 10

		for w := 0; w < workerCount; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				client := &http.Client{Timeout: time.Second}

				for {
					select {
					case <-ctx.Done():
						return
					default:
						resp, err := client.Get(server.URL + "/health")
						mu.Lock()
						requestCount++
						if err == nil && resp.StatusCode == http.StatusOK {
							successCount++
						}
						if resp != nil {
							_ = resp.Body.Close()
						}
						mu.Unlock()
					}
				}
			}()
		}

		wg.Wait()

		availability := float64(successCount) / float64(requestCount) * 100
		t.Logf("Requests: %d, Successes: %d, Availability: %.2f%%", requestCount, successCount, availability)

		// NFR-004: Availability > 99.9%
		assert.Greater(t, availability, 99.9, "Availability should be > 99.9%")
	})

	t.Run("websocket availability under load", func(t *testing.T) {
		wsSvc := ws.NewService(nil, nil, nil)
		err := wsSvc.Start(context.Background())
		require.NoError(t, err)
		defer func() { _ = wsSvc.Stop() }()

		duration := 5 * time.Second
		var operationCount int64
		var successCount int64

		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		var wg sync.WaitGroup
		workerCount := 10

		for w := 0; w < workerCount; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				for {
					select {
					case <-ctx.Done():
						return
					default:
						clientID := fmt.Sprintf("avail-client-%d", workerID)
						topic := "avail.topic"

						// Subscribe
						wsSvc.SubscriptionManager().Subscribe(clientID, topic)
						atomic.AddInt64(&operationCount, 1)

						if wsSvc.SubscriptionManager().IsSubscribed(clientID, topic) {
							atomic.AddInt64(&successCount, 1)
						}

						// Unsubscribe
						wsSvc.SubscriptionManager().Unsubscribe(clientID, topic)
						atomic.AddInt64(&operationCount, 1)

						if !wsSvc.SubscriptionManager().IsSubscribed(clientID, topic) {
							atomic.AddInt64(&successCount, 1)
						}
					}
				}
			}(w)
		}

		wg.Wait()

		availability := float64(successCount) / float64(operationCount) * 100
		t.Logf("Operations: %d, Successes: %d, Availability: %.2f%%", operationCount, successCount, availability)

		assert.Greater(t, availability, 99.9, "Availability should be > 99.9%")
	})
}

// TestAvailabilityGracefulDegradation tests graceful degradation under stress
func TestAvailabilityGracefulDegradation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping availability test in short mode")
	}

	t.Run("message processing under memory pressure", func(t *testing.T) {
		registry := processor.NewRegistry(nil)
		djiProcessor := djiuplink.NewProcessorAdapter(nil)
		registry.Register(djiProcessor)
		handler := processor.NewMessageHandler(registry, nil)

		// Process many messages with large payloads
		messageCount := 10000
		var successCount int64
		var errorCount int64

		var wg sync.WaitGroup
		workerCount := 10
		messagesPerWorker := messageCount / workerCount

		for w := 0; w < workerCount; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for i := 0; i < messagesPerWorker; i++ {
					// Create message with larger payload
					largeData := make(map[string]interface{})
					for j := 0; j < 100; j++ {
						largeData[fmt.Sprintf("field_%d", j)] = fmt.Sprintf("value_%d", j)
					}
					dataBytes, _ := json.Marshal(largeData)

					msg := &rabbitmq.StandardMessage{
						TID:       "stress-tid",
						BID:       "stress-bid",
						Service:   "iot-gateway",
						Action:    "telemetry.report",
						DeviceSN:  "STRESS-DEVICE",
						Timestamp: time.Now().UnixMilli(),
						Data:      dataBytes,
						ProtocolMeta: &rabbitmq.ProtocolMeta{
							Vendor: "dji",
						},
					}

					err := handler.Handle(context.Background(), msg)
					if err != nil {
						atomic.AddInt64(&errorCount, 1)
					} else {
						atomic.AddInt64(&successCount, 1)
					}
				}
			}()
		}

		wg.Wait()

		successRate := float64(successCount) / float64(messageCount) * 100
		t.Logf("Messages: %d, Successes: %d, Errors: %d, Success Rate: %.2f%%",
			messageCount, successCount, errorCount, successRate)

		// Should maintain high success rate even under stress
		assert.Greater(t, successRate, 99.0, "Success rate should be > 99%")
	})
}

// TestAvailabilityHealthChecks tests health check endpoints
func TestAvailabilityHealthChecks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping availability test in short mode")
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&models.Device{}, &model.ServiceCall{})
	require.NoError(t, err)

	config := &api.Config{
		EnableAuth:  false,
		EnableTrace: false,
	}
	router := api.NewRouter(config, db, nil, nil, nil)

	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("health endpoint always available", func(t *testing.T) {
		checkCount := 100
		successCount := 0

		for i := 0; i < checkCount; i++ {
			resp, err := http.Get(server.URL + "/health")
			if err == nil && resp.StatusCode == http.StatusOK {
				successCount++
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
		}

		assert.Equal(t, checkCount, successCount, "All health checks should succeed")
	})

	t.Run("ready endpoint always available", func(t *testing.T) {
		checkCount := 100
		successCount := 0

		for i := 0; i < checkCount; i++ {
			resp, err := http.Get(server.URL + "/ready")
			if err == nil && resp.StatusCode == http.StatusOK {
				successCount++
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
		}

		assert.Equal(t, checkCount, successCount, "All ready checks should succeed")
	})
}

// TestAvailabilityConcurrentOperations tests availability during concurrent operations
func TestAvailabilityConcurrentOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping availability test in short mode")
	}

	manager := connection.NewManager(nil)

	t.Run("concurrent register and query", func(t *testing.T) {
		var wg sync.WaitGroup
		var registerSuccess int64
		var querySuccess int64

		duration := 3 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		// Register workers
		for w := 0; w < 5; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				counter := 0
				for {
					select {
					case <-ctx.Done():
						return
					default:
						deviceSN := fmt.Sprintf("concurrent-device-%d-%d", workerID, counter%100)
						clientID := fmt.Sprintf("client-%d", counter)
						manager.Connect(deviceSN, clientID, "127.0.0.1")
						atomic.AddInt64(&registerSuccess, 1)
						counter++
					}
				}
			}(w)
		}

		// Query workers
		for w := 0; w < 5; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				counter := 0
				for {
					select {
					case <-ctx.Done():
						return
					default:
						deviceSN := fmt.Sprintf("concurrent-device-%d-%d", workerID%5, counter%100)
						_ = manager.IsOnline(deviceSN)
						atomic.AddInt64(&querySuccess, 1)
						counter++
					}
				}
			}(w)
		}

		wg.Wait()

		t.Logf("Registers: %d, Queries: %d", registerSuccess, querySuccess)

		// Both operations should complete successfully
		assert.Greater(t, registerSuccess, int64(0))
		assert.Greater(t, querySuccess, int64(0))
	})
}

// TestAvailabilityErrorRecovery tests recovery from errors
func TestAvailabilityErrorRecovery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping availability test in short mode")
	}

	t.Run("processor handles invalid messages", func(t *testing.T) {
		registry := processor.NewRegistry(nil)
		djiProcessor := djiuplink.NewProcessorAdapter(nil)
		registry.Register(djiProcessor)
		handler := processor.NewMessageHandler(registry, nil)

		// Mix of valid and invalid messages
		validCount := 0
		invalidCount := 0

		for i := 0; i < 100; i++ {
			var msg *rabbitmq.StandardMessage

			if i%10 == 0 {
				// Invalid message (missing vendor)
				msg = &rabbitmq.StandardMessage{
					TID:       "error-tid",
					BID:       "error-bid",
					Service:   "iot-gateway",
					Action:    "telemetry.report",
					DeviceSN:  "ERROR-DEVICE",
					Timestamp: time.Now().UnixMilli(),
					Data:      json.RawMessage(`{"test": "data"}`),
					// No ProtocolMeta
				}
				invalidCount++
			} else {
				// Valid message
				msg = &rabbitmq.StandardMessage{
					TID:       "valid-tid",
					BID:       "valid-bid",
					Service:   "iot-gateway",
					Action:    "telemetry.report",
					DeviceSN:  "VALID-DEVICE",
					Timestamp: time.Now().UnixMilli(),
					Data:      json.RawMessage(`{"test": "data"}`),
					ProtocolMeta: &rabbitmq.ProtocolMeta{
						Vendor: "dji",
					},
				}
				validCount++
			}

			// Should not panic
			_ = handler.Handle(context.Background(), msg)
		}

		t.Logf("Valid: %d, Invalid: %d", validCount, invalidCount)

		// Registry should still be functional after processing invalid messages
		_, found := registry.Get("dji")
		assert.True(t, found)
	})
}
