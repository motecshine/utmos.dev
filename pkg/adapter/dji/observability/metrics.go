// Package observability provides metrics, logging, and tracing for the DJI adapter.
package observability

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/utmos/utmos/pkg/metrics"
)

// Metrics holds all DJI adapter metrics.
type Metrics struct {
	// Message counters
	MessagesReceived *prometheus.CounterVec
	MessagesSent     *prometheus.CounterVec
	MessagesErrors   *prometheus.CounterVec

	// Processing duration
	ProcessingDuration *prometheus.HistogramVec

	// Active connections
	ActiveDevices *prometheus.GaugeVec
}

// NewMetrics creates and registers DJI adapter metrics.
func NewMetrics(collector *metrics.Collector) *Metrics {
	return &Metrics{
		MessagesReceived: collector.NewCounter(
			"dji_messages_received_total",
			"Total number of messages received from DJI devices",
			[]string{metrics.LabelMessageType, metrics.LabelStatus},
		),
		MessagesSent: collector.NewCounter(
			"dji_messages_sent_total",
			"Total number of messages sent to DJI devices",
			[]string{metrics.LabelMessageType, metrics.LabelStatus},
		),
		MessagesErrors: collector.NewCounter(
			"dji_messages_errors_total",
			"Total number of message processing errors",
			[]string{metrics.LabelMessageType, "error_type"},
		),
		ProcessingDuration: collector.NewHistogram(
			"dji_message_processing_duration_seconds",
			"Duration of message processing in seconds",
			[]string{metrics.LabelMessageType},
			[]float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0},
		),
		ActiveDevices: collector.NewGauge(
			"dji_active_devices",
			"Number of active DJI devices",
			[]string{"device_type"},
		),
	}
}

// RecordMessageReceived records a received message.
func (m *Metrics) RecordMessageReceived(messageType, status string) {
	m.MessagesReceived.WithLabelValues(messageType, status).Inc()
}

// RecordMessageSent records a sent message.
func (m *Metrics) RecordMessageSent(messageType, status string) {
	m.MessagesSent.WithLabelValues(messageType, status).Inc()
}

// RecordError records a message processing error.
func (m *Metrics) RecordError(messageType, errorType string) {
	m.MessagesErrors.WithLabelValues(messageType, errorType).Inc()
}

// RecordProcessingDuration records the duration of message processing.
func (m *Metrics) RecordProcessingDuration(messageType string, durationSeconds float64) {
	m.ProcessingDuration.WithLabelValues(messageType).Observe(durationSeconds)
}

// SetActiveDevices sets the number of active devices.
func (m *Metrics) SetActiveDevices(deviceType string, count float64) {
	m.ActiveDevices.WithLabelValues(deviceType).Set(count)
}
