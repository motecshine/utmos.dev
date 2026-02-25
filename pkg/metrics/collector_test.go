package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

// DefaultBuckets for histogram metrics
var DefaultBuckets = prometheus.DefBuckets

func TestNewCollector(t *testing.T) {
	collector := NewCollector("test")
	if collector == nil {
		t.Fatal("expected non-nil collector")
	}

	if collector.namespace != "test" {
		t.Errorf("expected namespace test, got %s", collector.namespace)
	}

	if collector.registry == nil {
		t.Error("expected non-nil registry")
	}
}

func TestCollector_Registry(t *testing.T) {
	collector := NewCollector("test")
	registry := collector.Registry()

	if registry == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestCollector_NewCounter(t *testing.T) {
	collector := NewCollector("test")

	counter := collector.NewCounter("requests", "Total requests", []string{"method", "status"})
	if counter == nil {
		t.Fatal("expected non-nil counter")
	}

	// Should be able to use the counter
	counter.WithLabelValues("GET", "200").Inc()
}

func TestCollector_NewHistogram(t *testing.T) {
	collector := NewCollector("test")

	histogram := collector.NewHistogram(
		"request_duration",
		"Request duration in seconds",
		[]string{"method"},
		[]float64{0.1, 0.5, 1.0, 2.5, 5.0},
	)
	if histogram == nil {
		t.Fatal("expected non-nil histogram")
	}

	// Should be able to observe values
	histogram.WithLabelValues("GET").Observe(0.25)
}

func TestCollector_NewGauge(t *testing.T) {
	collector := NewCollector("test")

	gauge := collector.NewGauge("connections", "Active connections", []string{"service"})
	if gauge == nil {
		t.Fatal("expected non-nil gauge")
	}

	// Should be able to set values
	gauge.WithLabelValues("api").Set(10)
	gauge.WithLabelValues("api").Inc()
	gauge.WithLabelValues("api").Dec()
}

func TestCollector_MultipleMetrics(_ *testing.T) {
	collector := NewCollector("iot")

	// Create multiple metrics
	counter := collector.NewCounter("messages_total", "Total messages", []string{"type"})
	histogram := collector.NewHistogram("processing_seconds", "Processing time", []string{"service"}, DefaultBuckets)
	gauge := collector.NewGauge("active_devices", "Active devices", []string{"vendor"})

	// Use all metrics
	counter.WithLabelValues("uplink").Add(100)
	histogram.WithLabelValues("gateway").Observe(0.05)
	gauge.WithLabelValues("dji").Set(42)

	// All should work without panicking
}

func TestCollector_DuplicateMetricName(t *testing.T) {
	collector := NewCollector("test")

	// Create first counter
	_ = collector.NewCounter("duplicate", "First counter", []string{"label"})

	// Creating same metric again should panic (Prometheus behavior)
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on duplicate metric registration")
		}
	}()

	_ = collector.NewCounter("duplicate", "Second counter", []string{"label"})
}

func TestDefaultBucketsOrder(t *testing.T) {
	if len(DefaultBuckets) == 0 {
		t.Error("expected non-empty default buckets")
	}

	// Verify buckets are in ascending order
	for i := 1; i < len(DefaultBuckets); i++ {
		if DefaultBuckets[i] <= DefaultBuckets[i-1] {
			t.Errorf("buckets not in ascending order at index %d", i)
		}
	}
}

func TestLabelConstants(t *testing.T) {
	// Verify label constants are defined
	labels := []string{LabelService, LabelVendor, LabelMessageType, LabelStatus}

	for _, label := range labels {
		if label == "" {
			t.Error("expected non-empty label constant")
		}
	}
}

func TestCollector_NewSummary(t *testing.T) {
	collector := NewCollector("test")

	summary := collector.NewSummary(
		"request_latency",
		"Request latency in seconds",
		[]string{"method"},
	)
	if summary == nil {
		t.Fatal("expected non-nil summary")
	}

	// Should be able to observe values
	summary.WithLabelValues("GET").Observe(0.25)
}

func TestCollector_NewHistogramWithNilBuckets(t *testing.T) {
	collector := NewCollector("test")

	// Should use default buckets when nil is passed
	histogram := collector.NewHistogram(
		"request_duration_default",
		"Request duration in seconds",
		[]string{"method"},
		nil,
	)
	if histogram == nil {
		t.Fatal("expected non-nil histogram")
	}

	histogram.WithLabelValues("GET").Observe(0.25)
}
