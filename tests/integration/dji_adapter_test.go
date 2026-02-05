package integration

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/rabbitmq"
	"github.com/utmos/utmos/tests/mocks"
)

// TestDJIAdapterUplinkFlow tests the complete uplink message flow:
// Raw DJI message -> DJI Adapter -> Standard Message
func TestDJIAdapterUplinkFlow(t *testing.T) {
	// Register DJI adapter
	dji.Register()
	defer adapter.Unregister(dji.VendorDJI)

	// Get adapter from registry
	djiAdapter, err := adapter.Get(dji.VendorDJI)
	require.NoError(t, err)

	t.Run("OSD message flow", func(t *testing.T) {
		// Simulate raw message from iot-gateway
		topic, payload := mocks.NewOSDMessage("1ZNBH1D00C00FK")

		// Step 1: Parse raw message
		pm, err := djiAdapter.ParseRawMessage(topic, payload)
		require.NoError(t, err)
		assert.Equal(t, dji.VendorDJI, pm.Vendor)
		assert.Equal(t, "1ZNBH1D00C00FK", pm.DeviceSN)
		assert.Equal(t, adapter.MessageTypeProperty, pm.MessageType)

		// Step 2: Convert to standard message
		stdMsg, err := djiAdapter.ToStandardMessage(pm)
		require.NoError(t, err)
		assert.Equal(t, "1ZNBH1D00C00FK", stdMsg.DeviceSN)
		assert.Equal(t, "property.report", stdMsg.Action)
		assert.NotEmpty(t, stdMsg.TID)
		assert.NotNil(t, stdMsg.ProtocolMeta)
		assert.Equal(t, dji.VendorDJI, stdMsg.ProtocolMeta.Vendor)

		// Verify data is preserved
		assert.NotNil(t, stdMsg.Data)
	})

	t.Run("State message flow", func(t *testing.T) {
		topic, payload := mocks.NewStateMessage("1ZNBH1D00C00FK")

		pm, err := djiAdapter.ParseRawMessage(topic, payload)
		require.NoError(t, err)
		assert.Equal(t, adapter.MessageTypeProperty, pm.MessageType)

		stdMsg, err := djiAdapter.ToStandardMessage(pm)
		require.NoError(t, err)
		assert.Equal(t, "property.report", stdMsg.Action)
	})

	t.Run("Event message flow", func(t *testing.T) {
		topic, payload := mocks.NewEventMessage("DOCK001", "fly_to_point_progress", 50)

		pm, err := djiAdapter.ParseRawMessage(topic, payload)
		require.NoError(t, err)
		assert.Equal(t, adapter.MessageTypeEvent, pm.MessageType)
		assert.Equal(t, "fly_to_point_progress", pm.Method)

		stdMsg, err := djiAdapter.ToStandardMessage(pm)
		require.NoError(t, err)
		assert.Equal(t, "event.report", stdMsg.Action)
	})

	t.Run("Services request message flow", func(t *testing.T) {
		topic, payload := mocks.NewServicesRequestMessage("DOCK001", "takeoff", map[string]any{
			"height": 50,
		})

		pm, err := djiAdapter.ParseRawMessage(topic, payload)
		require.NoError(t, err)
		assert.Equal(t, adapter.MessageTypeService, pm.MessageType)
		assert.Equal(t, "takeoff", pm.Method)

		stdMsg, err := djiAdapter.ToStandardMessage(pm)
		require.NoError(t, err)
		assert.Equal(t, "service.call", stdMsg.Action)
	})

	t.Run("Status message flow", func(t *testing.T) {
		topic, payload := mocks.NewStatusMessage("DOCK001", true)

		pm, err := djiAdapter.ParseRawMessage(topic, payload)
		require.NoError(t, err)
		assert.Equal(t, adapter.MessageTypeStatus, pm.MessageType)

		stdMsg, err := djiAdapter.ToStandardMessage(pm)
		require.NoError(t, err)
		assert.Equal(t, "device.online", stdMsg.Action)
	})
}

// TestDJIAdapterDownlinkFlow tests the complete downlink message flow:
// Standard Message -> DJI Adapter -> Raw DJI message
func TestDJIAdapterDownlinkFlow(t *testing.T) {
	// Register DJI adapter
	dji.Register()
	defer adapter.Unregister(dji.VendorDJI)

	djiAdapter, err := adapter.Get(dji.VendorDJI)
	require.NoError(t, err)

	t.Run("Service call downlink", func(t *testing.T) {
		// Create a standard message for service call
		stdMsg := &rabbitmq.StandardMessage{
			DeviceSN:  "1ZNBH1D00C00FK",
			TID:       "tid-downlink-001",
			BID:       "bid-downlink-001",
			Timestamp: time.Now().UnixMilli(),
			Action:    "service.call",
			Data:      json.RawMessage(`{"method":"takeoff","params":{"height":50}}`),
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: dji.VendorDJI,
			},
		}

		// Step 1: Convert from standard message
		pm, err := djiAdapter.FromStandardMessage(stdMsg)
		require.NoError(t, err)
		assert.Equal(t, dji.VendorDJI, pm.Vendor)
		assert.Equal(t, "1ZNBH1D00C00FK", pm.DeviceSN)
		assert.Equal(t, adapter.MessageTypeService, pm.MessageType)

		// Step 2: Get raw payload
		rawPayload, err := djiAdapter.GetRawPayload(pm)
		require.NoError(t, err)
		assert.NotEmpty(t, rawPayload)

		// Verify the raw payload is valid JSON
		var djiMsg map[string]any
		err = json.Unmarshal(rawPayload, &djiMsg)
		require.NoError(t, err)
		assert.Contains(t, djiMsg, "tid")
		assert.Contains(t, djiMsg, "bid")
		assert.Contains(t, djiMsg, "timestamp")
	})

	t.Run("Property set downlink", func(t *testing.T) {
		stdMsg := &rabbitmq.StandardMessage{
			DeviceSN:  "1ZNBH1D00C00FK",
			TID:       "tid-prop-001",
			BID:       "bid-prop-001",
			Timestamp: time.Now().UnixMilli(),
			Action:    "property.set",
			Data:      json.RawMessage(`{"brightness":80}`),
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: dji.VendorDJI,
			},
		}

		pm, err := djiAdapter.FromStandardMessage(stdMsg)
		require.NoError(t, err)
		assert.Equal(t, adapter.MessageTypeProperty, pm.MessageType)

		rawPayload, err := djiAdapter.GetRawPayload(pm)
		require.NoError(t, err)
		assert.NotEmpty(t, rawPayload)
	})
}

// TestDJIAdapterRoundTrip tests that messages can be converted back and forth
// without losing essential information.
func TestDJIAdapterRoundTrip(t *testing.T) {
	dji.Register()
	defer adapter.Unregister(dji.VendorDJI)

	djiAdapter, err := adapter.Get(dji.VendorDJI)
	require.NoError(t, err)

	t.Run("OSD round trip", func(t *testing.T) {
		// Original raw message
		topic, payload := mocks.NewOSDMessage("DEVICE001")

		// Uplink: Raw -> Protocol -> Standard
		pm1, err := djiAdapter.ParseRawMessage(topic, payload)
		require.NoError(t, err)

		stdMsg, err := djiAdapter.ToStandardMessage(pm1)
		require.NoError(t, err)

		// Verify essential fields are preserved
		assert.Equal(t, pm1.DeviceSN, stdMsg.DeviceSN)
		assert.Equal(t, pm1.TID, stdMsg.TID)
		assert.Equal(t, pm1.BID, stdMsg.BID)
	})

	t.Run("Service call round trip", func(t *testing.T) {
		// Original standard message
		originalStdMsg := &rabbitmq.StandardMessage{
			DeviceSN:  "DEVICE002",
			TID:       "tid-round-001",
			BID:       "bid-round-001",
			Timestamp: time.Now().UnixMilli(),
			Action:    "service.call",
			Data:      json.RawMessage(`{"command":"start"}`),
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: dji.VendorDJI,
			},
		}

		// Downlink: Standard -> Protocol -> Raw
		pm, err := djiAdapter.FromStandardMessage(originalStdMsg)
		require.NoError(t, err)

		rawPayload, err := djiAdapter.GetRawPayload(pm)
		require.NoError(t, err)

		// Verify the raw payload contains expected fields
		var djiMsg map[string]any
		err = json.Unmarshal(rawPayload, &djiMsg)
		require.NoError(t, err)
		assert.Equal(t, originalStdMsg.TID, djiMsg["tid"])
		assert.Equal(t, originalStdMsg.BID, djiMsg["bid"])
	})
}

// TestDJIAdapterErrorHandling tests error handling scenarios.
func TestDJIAdapterErrorHandling(t *testing.T) {
	dji.Register()
	defer adapter.Unregister(dji.VendorDJI)

	djiAdapter, err := adapter.Get(dji.VendorDJI)
	require.NoError(t, err)

	t.Run("invalid topic", func(t *testing.T) {
		_, err := djiAdapter.ParseRawMessage("invalid/topic", []byte(`{}`))
		assert.Error(t, err)
	})

	t.Run("invalid JSON payload", func(t *testing.T) {
		_, err := djiAdapter.ParseRawMessage("thing/product/DEVICE001/osd", []byte(`not json`))
		assert.Error(t, err)
	})

	t.Run("empty payload", func(t *testing.T) {
		// Empty payload should fail
		_, err := djiAdapter.ParseRawMessage("thing/product/DEVICE001/osd", []byte(``))
		assert.Error(t, err)
	})
}

// TestDJIAdapterWithFactory tests using the adapter factory.
func TestDJIAdapterWithFactory(t *testing.T) {
	dji.Register()
	defer adapter.Unregister(dji.VendorDJI)

	factory := adapter.NewFactory()

	t.Run("create DJI adapter via factory", func(t *testing.T) {
		djiAdapter, err := factory.NewAdapter(dji.VendorDJI)
		require.NoError(t, err)
		assert.Equal(t, dji.VendorDJI, djiAdapter.GetVendor())
	})

	t.Run("list available vendors", func(t *testing.T) {
		vendors := factory.ListVendors()
		assert.Contains(t, vendors, dji.VendorDJI)
	})

	t.Run("check vendor availability", func(t *testing.T) {
		assert.True(t, factory.HasVendor(dji.VendorDJI))
		assert.False(t, factory.HasVendor("unknown-vendor"))
	})
}

// TestSampleMessagesIntegration tests all sample messages from mocks package.
func TestSampleMessagesIntegration(t *testing.T) {
	dji.Register()
	defer adapter.Unregister(dji.VendorDJI)

	djiAdapter, err := adapter.Get(dji.VendorDJI)
	require.NoError(t, err)

	samples := mocks.SampleMessages()

	for name, sample := range samples {
		t.Run(name, func(t *testing.T) {
			// Parse raw message
			pm, err := djiAdapter.ParseRawMessage(sample.Topic, sample.Payload)
			require.NoError(t, err, "Failed to parse %s message", name)
			assert.NotEmpty(t, pm.Vendor)
			assert.NotEmpty(t, pm.DeviceSN)

			// Convert to standard message
			stdMsg, err := djiAdapter.ToStandardMessage(pm)
			require.NoError(t, err, "Failed to convert %s to standard message", name)
			assert.NotEmpty(t, stdMsg.DeviceSN)
			assert.NotEmpty(t, stdMsg.Action)
		})
	}
}
