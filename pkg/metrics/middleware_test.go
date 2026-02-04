package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestMiddlewareMetrics_RabbitMQ(_ *testing.T) {
	collector := NewCollector("iot")

	// RabbitMQ connection metrics
	connectionTotal := collector.NewCounter(
		"rabbitmq_connection_total",
		"Total RabbitMQ connection attempts",
		[]string{LabelStatus},
	)

	connectionTotal.WithLabelValues("success").Inc()
	connectionTotal.WithLabelValues("failure").Inc()

	// RabbitMQ message metrics
	messageTotal := collector.NewCounter(
		"rabbitmq_message_total",
		"Total RabbitMQ messages",
		[]string{"direction", LabelStatus},
	)

	messageTotal.WithLabelValues("publish", "success").Add(100)
	messageTotal.WithLabelValues("publish", "failure").Add(2)
	messageTotal.WithLabelValues("consume", "success").Add(98)
	messageTotal.WithLabelValues("consume", "failure").Add(1)

	// RabbitMQ message duration
	messageDuration := collector.NewHistogram(
		"rabbitmq_message_duration_seconds",
		"RabbitMQ message processing duration",
		[]string{"direction"},
		prometheus.DefBuckets,
	)

	messageDuration.WithLabelValues("publish").Observe(0.005)
	messageDuration.WithLabelValues("consume").Observe(0.015)
}

func TestMiddlewareMetrics_PostgreSQL(_ *testing.T) {
	collector := NewCollector("iot")

	// PostgreSQL connection pool
	poolSize := collector.NewGauge(
		"postgres_connection_pool_size",
		"Current PostgreSQL connection pool size",
		[]string{"state"},
	)

	poolSize.WithLabelValues("idle").Set(5)
	poolSize.WithLabelValues("in_use").Set(3)
	poolSize.WithLabelValues("max").Set(20)

	// PostgreSQL query duration
	queryDuration := collector.NewHistogram(
		"postgres_query_duration_seconds",
		"PostgreSQL query execution duration",
		[]string{"operation"},
		[]float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0},
	)

	queryDuration.WithLabelValues("select").Observe(0.003)
	queryDuration.WithLabelValues("insert").Observe(0.008)
	queryDuration.WithLabelValues("update").Observe(0.012)
	queryDuration.WithLabelValues("delete").Observe(0.006)

	// PostgreSQL errors
	errorTotal := collector.NewCounter(
		"postgres_error_total",
		"Total PostgreSQL errors",
		[]string{"operation", "error_type"},
	)

	errorTotal.WithLabelValues("select", "timeout").Inc()
	errorTotal.WithLabelValues("insert", "constraint_violation").Inc()
}

func TestMiddlewareMetrics_InfluxDB(_ *testing.T) {
	collector := NewCollector("iot")

	// InfluxDB write duration
	writeDuration := collector.NewHistogram(
		"influxdb_write_duration_seconds",
		"InfluxDB write operation duration",
		[]string{"measurement"},
		[]float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5},
	)

	writeDuration.WithLabelValues("device_telemetry").Observe(0.008)
	writeDuration.WithLabelValues("device_events").Observe(0.005)

	// InfluxDB errors
	errorTotal := collector.NewCounter(
		"influxdb_error_total",
		"Total InfluxDB errors",
		[]string{"operation"},
	)

	errorTotal.WithLabelValues("write").Inc()
	errorTotal.WithLabelValues("query").Inc()

	// InfluxDB batch size
	batchSize := collector.NewHistogram(
		"influxdb_batch_size",
		"InfluxDB write batch size",
		[]string{"measurement"},
		[]float64{1, 10, 50, 100, 500, 1000, 5000},
	)

	batchSize.WithLabelValues("device_telemetry").Observe(100)
}

func TestMiddlewareMetrics_HTTP(_ *testing.T) {
	collector := NewCollector("iot")

	// HTTP request total
	requestTotal := collector.NewCounter(
		"http_requests_total",
		"Total HTTP requests",
		[]string{"method", "path", "status_code"},
	)

	requestTotal.WithLabelValues("GET", "/health", "200").Inc()
	requestTotal.WithLabelValues("GET", "/ready", "200").Inc()
	requestTotal.WithLabelValues("POST", "/api/v1/devices", "201").Inc()
	requestTotal.WithLabelValues("GET", "/api/v1/devices", "200").Add(50)

	// HTTP request duration
	requestDuration := collector.NewHistogram(
		"http_request_duration_seconds",
		"HTTP request duration",
		[]string{"method", "path"},
		[]float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5},
	)

	requestDuration.WithLabelValues("GET", "/health").Observe(0.001)
	requestDuration.WithLabelValues("POST", "/api/v1/devices").Observe(0.025)

	// HTTP response size
	responseSize := collector.NewHistogram(
		"http_response_size_bytes",
		"HTTP response size in bytes",
		[]string{"method", "path"},
		[]float64{100, 500, 1000, 5000, 10000, 50000},
	)

	responseSize.WithLabelValues("GET", "/api/v1/devices").Observe(2500)
}

func TestMiddlewareMetrics_WebSocket(_ *testing.T) {
	collector := NewCollector("iot")

	// WebSocket connections
	wsConnections := collector.NewGauge(
		"websocket_connections_active",
		"Active WebSocket connections",
		[]string{LabelVendor},
	)

	wsConnections.WithLabelValues("dji").Set(25)
	wsConnections.WithLabelValues("tuya").Set(50)
	wsConnections.WithLabelValues("generic").Set(10)

	// WebSocket messages
	wsMessages := collector.NewCounter(
		"websocket_messages_total",
		"Total WebSocket messages",
		[]string{"direction", LabelVendor},
	)

	wsMessages.WithLabelValues("sent", "dji").Add(1000)
	wsMessages.WithLabelValues("received", "dji").Add(500)

	// WebSocket connection duration
	wsConnectionDuration := collector.NewHistogram(
		"websocket_connection_duration_seconds",
		"WebSocket connection duration",
		[]string{LabelVendor},
		[]float64{1, 10, 60, 300, 600, 1800, 3600},
	)

	wsConnectionDuration.WithLabelValues("dji").Observe(300)
}

func TestMiddlewareMetrics_MQTT(_ *testing.T) {
	collector := NewCollector("iot")

	// MQTT connections
	mqttConnections := collector.NewGauge(
		"mqtt_connections_active",
		"Active MQTT connections",
		[]string{},
	)

	mqttConnections.WithLabelValues().Set(1)

	// MQTT messages
	mqttMessages := collector.NewCounter(
		"mqtt_messages_total",
		"Total MQTT messages",
		[]string{"direction", "qos"},
	)

	mqttMessages.WithLabelValues("publish", "0").Add(5000)
	mqttMessages.WithLabelValues("publish", "1").Add(1000)
	mqttMessages.WithLabelValues("subscribe", "1").Add(6000)
}
