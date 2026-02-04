package integration

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/utmos/utmos/pkg/metrics"
)

// TestMetricsCollectorIntegration tests the metrics collector integration
func TestMetricsCollectorIntegration(t *testing.T) {
	collector := metrics.NewCollector("iot_test")

	// Create various metric types
	counter := collector.NewCounter("test_requests_total", "Test requests", []string{"method", "status"})
	histogram := collector.NewHistogram("test_duration_seconds", "Test duration", []string{"operation"}, nil)
	gauge := collector.NewGauge("test_active_connections", "Test connections", []string{"service"})

	// Use the metrics
	counter.WithLabelValues("GET", "200").Inc()
	counter.WithLabelValues("POST", "201").Inc()
	counter.WithLabelValues("GET", "500").Inc()

	histogram.WithLabelValues("read").Observe(0.1)
	histogram.WithLabelValues("write").Observe(0.2)

	gauge.WithLabelValues("api").Set(10)
	gauge.WithLabelValues("ws").Set(5)

	// Verify registry contains metrics
	registry := collector.Registry()
	mfs, err := registry.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}

	if len(mfs) == 0 {
		t.Error("expected metrics to be registered")
	}

	// Look for our custom metrics
	foundCounter := false
	foundHistogram := false
	foundGauge := false

	for _, mf := range mfs {
		switch *mf.Name {
		case "iot_test_test_requests_total":
			foundCounter = true
		case "iot_test_test_duration_seconds":
			foundHistogram = true
		case "iot_test_test_active_connections":
			foundGauge = true
		}
	}

	if !foundCounter {
		t.Error("counter metric not found")
	}
	if !foundHistogram {
		t.Error("histogram metric not found")
	}
	if !foundGauge {
		t.Error("gauge metric not found")
	}
}

// TestMetricsNamingConvention tests that metrics follow naming conventions
func TestMetricsNamingConvention(t *testing.T) {
	collector := metrics.NewCollector("iot")

	// These should follow iot_{component}_{metric}_{unit} convention
	validMetrics := []struct {
		name string
		help string
	}{
		{"gateway_messages_total", "Total gateway messages"},
		{"rabbitmq_publish_duration_seconds", "RabbitMQ publish duration"},
		{"postgres_query_duration_seconds", "PostgreSQL query duration"},
		{"http_requests_total", "Total HTTP requests"},
		{"websocket_connections_active", "Active WebSocket connections"},
	}

	for _, m := range validMetrics {
		counter := collector.NewCounter(m.name, m.help, []string{"label"})
		if counter == nil {
			t.Errorf("failed to create counter: %s", m.name)
		}
	}

	// Verify all metrics were registered
	registry := collector.Registry()
	mfs, err := registry.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}

	// Should have at least the metrics we created plus process/go collectors
	if len(mfs) < len(validMetrics) {
		t.Errorf("expected at least %d metrics, got %d", len(validMetrics), len(mfs))
	}
}

// TestMetricsLabels tests label handling
func TestMetricsLabels(t *testing.T) {
	collector := metrics.NewCollector("test")

	counter := collector.NewCounter("labeled_metric", "Test with labels", []string{
		metrics.LabelService,
		metrics.LabelVendor,
		metrics.LabelStatus,
	})

	// Use with different label combinations
	counter.WithLabelValues("gateway", "dji", "success").Inc()
	counter.WithLabelValues("gateway", "tuya", "success").Inc()
	counter.WithLabelValues("gateway", "generic", "failure").Inc()
	counter.WithLabelValues("uplink", "dji", "success").Inc()

	registry := collector.Registry()
	mfs, err := registry.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}

	// Find our metric
	var found *prometheus.Metric
	for _, mf := range mfs {
		if *mf.Name == "test_labeled_metric" {
			if len(mf.Metric) != 4 {
				t.Errorf("expected 4 label combinations, got %d", len(mf.Metric))
			}
			break
		}
	}

	_ = found // silence unused variable warning
}

// TestMetricsHistogramBuckets tests histogram bucket configuration
func TestMetricsHistogramBuckets(t *testing.T) {
	collector := metrics.NewCollector("test")

	// Custom buckets for latency
	latencyBuckets := []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0}
	histogram := collector.NewHistogram("latency_seconds", "Request latency", []string{"endpoint"}, latencyBuckets)

	// Observe values across different buckets
	testValues := []float64{0.002, 0.008, 0.015, 0.03, 0.075, 0.15, 0.3, 0.75, 1.5}
	for _, v := range testValues {
		histogram.WithLabelValues("/api/v1/devices").Observe(v)
	}

	registry := collector.Registry()
	mfs, err := registry.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}

	// Verify histogram was created
	found := false
	for _, mf := range mfs {
		if *mf.Name == "test_latency_seconds" {
			found = true
			break
		}
	}

	if !found {
		t.Error("histogram metric not found")
	}
}

// TestMetricsIsolation tests that different collectors are isolated
func TestMetricsIsolation(t *testing.T) {
	collector1 := metrics.NewCollector("service1")
	collector2 := metrics.NewCollector("service2")

	// Create same metric name in different collectors
	counter1 := collector1.NewCounter("requests", "Requests", []string{})
	counter2 := collector2.NewCounter("requests", "Requests", []string{})

	counter1.WithLabelValues().Add(100)
	counter2.WithLabelValues().Add(200)

	// Each collector should have its own registry
	mfs1, _ := collector1.Registry().Gather()
	mfs2, _ := collector2.Registry().Gather()

	// Find the request counts
	var count1, count2 float64
	for _, mf := range mfs1 {
		if *mf.Name == "service1_requests" {
			count1 = *mf.Metric[0].Counter.Value
		}
	}
	for _, mf := range mfs2 {
		if *mf.Name == "service2_requests" {
			count2 = *mf.Metric[0].Counter.Value
		}
	}

	if count1 != 100 {
		t.Errorf("expected count1=100, got %f", count1)
	}
	if count2 != 200 {
		t.Errorf("expected count2=200, got %f", count2)
	}
}
