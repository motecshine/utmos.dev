package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Business metrics helpers follow the naming convention: iot_{component}_{metric_type}_{unit}

// MessageMetrics provides message-related metrics.
type MessageMetrics struct {
	ProcessedTotal   *prometheus.CounterVec
	ProcessDuration  *prometheus.HistogramVec
	ErrorTotal       *prometheus.CounterVec
	QueueSize        *prometheus.GaugeVec
}

// NewMessageMetrics creates message metrics.
func NewMessageMetrics(collector *Collector) *MessageMetrics {
	return &MessageMetrics{
		ProcessedTotal: collector.NewCounter(
			"message_processed_total",
			"Total number of processed messages",
			[]string{LabelService, LabelVendor, LabelMessageType, LabelStatus},
		),
		ProcessDuration: collector.NewHistogram(
			"message_process_duration_seconds",
			"Message processing duration in seconds",
			[]string{LabelService, LabelVendor, LabelMessageType},
			[]float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		),
		ErrorTotal: collector.NewCounter(
			"message_error_total",
			"Total number of message processing errors",
			[]string{LabelService, LabelVendor, LabelMessageType},
		),
		QueueSize: collector.NewGauge(
			"message_queue_size",
			"Current message queue size",
			[]string{LabelService},
		),
	}
}

// DeviceMetrics provides device-related metrics.
type DeviceMetrics struct {
	OnlineTotal   *prometheus.GaugeVec
	EventTotal    *prometheus.CounterVec
	PropertyTotal *prometheus.CounterVec
}

// NewDeviceMetrics creates device metrics.
func NewDeviceMetrics(collector *Collector) *DeviceMetrics {
	return &DeviceMetrics{
		OnlineTotal: collector.NewGauge(
			"device_online_total",
			"Total number of online devices",
			[]string{LabelService, LabelVendor},
		),
		EventTotal: collector.NewCounter(
			"device_event_total",
			"Total number of device events",
			[]string{LabelService, LabelVendor},
		),
		PropertyTotal: collector.NewCounter(
			"device_property_total",
			"Total number of device property updates",
			[]string{LabelService, LabelVendor},
		),
	}
}

// HTTPMetrics provides HTTP-related metrics.
type HTTPMetrics struct {
	RequestTotal    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	RequestSize     *prometheus.HistogramVec
	ResponseSize    *prometheus.HistogramVec
}

// NewHTTPMetrics creates HTTP metrics.
func NewHTTPMetrics(collector *Collector) *HTTPMetrics {
	return &HTTPMetrics{
		RequestTotal: collector.NewCounter(
			"http_request_total",
			"Total number of HTTP requests",
			[]string{LabelService, LabelMethod, LabelPath, LabelCode},
		),
		RequestDuration: collector.NewHistogram(
			"http_request_duration_seconds",
			"HTTP request duration in seconds",
			[]string{LabelService, LabelMethod, LabelPath},
			[]float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		),
		RequestSize: collector.NewHistogram(
			"http_request_size_bytes",
			"HTTP request size in bytes",
			[]string{LabelService, LabelMethod, LabelPath},
			prometheus.ExponentialBuckets(100, 10, 7),
		),
		ResponseSize: collector.NewHistogram(
			"http_response_size_bytes",
			"HTTP response size in bytes",
			[]string{LabelService, LabelMethod, LabelPath},
			prometheus.ExponentialBuckets(100, 10, 7),
		),
	}
}
