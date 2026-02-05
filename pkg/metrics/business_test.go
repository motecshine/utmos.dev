package metrics

import (
	"strings"
	"testing"
)

func TestBusinessMetrics_Counter(t *testing.T) {
	collector := NewCollector("iot")

	// Test counter creation with proper naming convention
	counter := collector.NewCounter(
		"gateway_messages_total",
		"Total messages processed by gateway",
		[]string{LabelVendor, LabelMessageType},
	)

	if counter == nil {
		t.Fatal("expected non-nil counter")
	}

	// Test incrementing
	counter.WithLabelValues("dji", "uplink").Inc()
	counter.WithLabelValues("dji", "uplink").Add(10)
	counter.WithLabelValues("tuya", "downlink").Inc()
}

func TestBusinessMetrics_Histogram(t *testing.T) {
	collector := NewCollector("iot")

	// Test histogram for message processing duration
	histogram := collector.NewHistogram(
		"message_processing_seconds",
		"Message processing duration in seconds",
		[]string{LabelService, LabelVendor},
		[]float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0},
	)

	if histogram == nil {
		t.Fatal("expected non-nil histogram")
	}

	// Test observing values
	histogram.WithLabelValues("uplink", "dji").Observe(0.015)
	histogram.WithLabelValues("uplink", "generic").Observe(0.008)
	histogram.WithLabelValues("downlink", "tuya").Observe(0.025)
}

func TestBusinessMetrics_Gauge(t *testing.T) {
	collector := NewCollector("iot")

	// Test gauge for active connections
	gauge := collector.NewGauge(
		"active_connections",
		"Number of active connections",
		[]string{LabelService},
	)

	if gauge == nil {
		t.Fatal("expected non-nil gauge")
	}

	// Test setting, incrementing, decrementing
	gauge.WithLabelValues("gateway").Set(100)
	gauge.WithLabelValues("gateway").Inc()
	gauge.WithLabelValues("gateway").Dec()
	gauge.WithLabelValues("websocket").Add(50)
	gauge.WithLabelValues("websocket").Sub(10)
}

func TestBusinessMetrics_NamingConvention(t *testing.T) {
	// Verify naming convention: iot_{component}_{metric_type}_{unit}
	validNames := []string{
		"gateway_messages_total",
		"rabbitmq_connection_total",
		"message_processing_seconds",
		"postgres_query_duration_seconds",
		"active_devices_count",
	}

	for _, name := range validNames {
		// Names should be lowercase with underscores
		if strings.ToLower(name) != name {
			t.Errorf("metric name should be lowercase: %s", name)
		}

		// Names should not contain dashes
		if strings.Contains(name, "-") {
			t.Errorf("metric name should not contain dashes: %s", name)
		}
	}
}

func TestBusinessMetrics_DeviceMetrics(_ *testing.T) {
	collector := NewCollector("iot")

	// Device online/offline counter
	deviceStatus := collector.NewCounter(
		"device_status_changes_total",
		"Total device status changes",
		[]string{LabelVendor, LabelStatus},
	)

	deviceStatus.WithLabelValues("dji", "online").Inc()
	deviceStatus.WithLabelValues("dji", "offline").Inc()
	deviceStatus.WithLabelValues("tuya", "online").Add(5)

	// Active devices gauge
	activeDevices := collector.NewGauge(
		"devices_active",
		"Number of currently active devices",
		[]string{LabelVendor},
	)

	activeDevices.WithLabelValues("dji").Set(150)
	activeDevices.WithLabelValues("tuya").Set(300)
	activeDevices.WithLabelValues("generic").Set(50)
}

func TestBusinessMetrics_MessageMetrics(_ *testing.T) {
	collector := NewCollector("iot")

	// Message throughput
	messageCounter := collector.NewCounter(
		"messages_processed_total",
		"Total messages processed",
		[]string{LabelService, LabelVendor, LabelMessageType},
	)

	messageCounter.WithLabelValues("gateway", "dji", "property_report").Add(1000)
	messageCounter.WithLabelValues("gateway", "dji", "event").Add(50)
	messageCounter.WithLabelValues("uplink", "generic", "property_report").Add(500)

	// Message size histogram
	messageSizeHistogram := collector.NewHistogram(
		"message_size_bytes",
		"Message size in bytes",
		[]string{LabelMessageType},
		[]float64{64, 128, 256, 512, 1024, 2048, 4096, 8192},
	)

	messageSizeHistogram.WithLabelValues("property_report").Observe(256)
	messageSizeHistogram.WithLabelValues("event").Observe(512)
	messageSizeHistogram.WithLabelValues("command").Observe(128)
}

func TestBusinessMetrics_ErrorMetrics(_ *testing.T) {
	collector := NewCollector("iot")

	// Error counter
	errorCounter := collector.NewCounter(
		"errors_total",
		"Total errors by type",
		[]string{LabelService, "error_code"},
	)

	errorCounter.WithLabelValues("gateway", "1001").Inc()
	errorCounter.WithLabelValues("uplink", "2001").Inc()
	errorCounter.WithLabelValues("downlink", "3001").Add(3)
}

func TestBusinessMetrics_LatencyMetrics(_ *testing.T) {
	collector := NewCollector("iot")

	// End-to-end latency
	latencyHistogram := collector.NewHistogram(
		"e2e_latency_seconds",
		"End-to-end message latency",
		[]string{LabelVendor, "direction"},
		[]float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0},
	)

	latencyHistogram.WithLabelValues("dji", "uplink").Observe(0.035)
	latencyHistogram.WithLabelValues("dji", "downlink").Observe(0.045)
	latencyHistogram.WithLabelValues("tuya", "uplink").Observe(0.028)
}

func TestNewMessageMetrics(t *testing.T) {
	collector := NewCollector("iot")
	metrics := NewMessageMetrics(collector)

	if metrics == nil {
		t.Fatal("expected non-nil MessageMetrics")
	}
	if metrics.ProcessedTotal == nil {
		t.Error("expected non-nil ProcessedTotal")
	}
	if metrics.ProcessDuration == nil {
		t.Error("expected non-nil ProcessDuration")
	}
	if metrics.ErrorTotal == nil {
		t.Error("expected non-nil ErrorTotal")
	}
	if metrics.QueueSize == nil {
		t.Error("expected non-nil QueueSize")
	}

	// Test using the metrics
	metrics.ProcessedTotal.WithLabelValues("gateway", "dji", "property", "success").Inc()
	metrics.ProcessDuration.WithLabelValues("gateway", "dji", "property").Observe(0.05)
	metrics.ErrorTotal.WithLabelValues("gateway", "dji", "property").Inc()
	metrics.QueueSize.WithLabelValues("gateway").Set(100)
}

func TestNewDeviceMetrics(t *testing.T) {
	collector := NewCollector("iot")
	metrics := NewDeviceMetrics(collector)

	if metrics == nil {
		t.Fatal("expected non-nil DeviceMetrics")
	}
	if metrics.OnlineTotal == nil {
		t.Error("expected non-nil OnlineTotal")
	}
	if metrics.EventTotal == nil {
		t.Error("expected non-nil EventTotal")
	}
	if metrics.PropertyTotal == nil {
		t.Error("expected non-nil PropertyTotal")
	}

	// Test using the metrics
	metrics.OnlineTotal.WithLabelValues("gateway", "dji").Set(50)
	metrics.EventTotal.WithLabelValues("gateway", "dji").Inc()
	metrics.PropertyTotal.WithLabelValues("gateway", "dji").Add(10)
}

func TestNewHTTPMetrics(t *testing.T) {
	collector := NewCollector("iot")
	metrics := NewHTTPMetrics(collector)

	if metrics == nil {
		t.Fatal("expected non-nil HTTPMetrics")
	}
	if metrics.RequestTotal == nil {
		t.Error("expected non-nil RequestTotal")
	}
	if metrics.RequestDuration == nil {
		t.Error("expected non-nil RequestDuration")
	}
	if metrics.RequestSize == nil {
		t.Error("expected non-nil RequestSize")
	}
	if metrics.ResponseSize == nil {
		t.Error("expected non-nil ResponseSize")
	}

	// Test using the metrics
	metrics.RequestTotal.WithLabelValues("api", "GET", "/health", "200").Inc()
	metrics.RequestDuration.WithLabelValues("api", "GET", "/health").Observe(0.01)
	metrics.RequestSize.WithLabelValues("api", "POST", "/devices").Observe(1024)
	metrics.ResponseSize.WithLabelValues("api", "GET", "/devices").Observe(2048)
}
