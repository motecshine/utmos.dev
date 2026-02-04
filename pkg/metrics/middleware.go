package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// MiddlewareMetrics provides metrics for infrastructure middleware.
type MiddlewareMetrics struct {
	RabbitMQ   *RabbitMQMetrics
	PostgreSQL *PostgreSQLMetrics
	InfluxDB   *InfluxDBMetrics
}

// RabbitMQMetrics provides RabbitMQ-related metrics.
type RabbitMQMetrics struct {
	ConnectionTotal   *prometheus.GaugeVec
	MessageTotal      *prometheus.CounterVec
	MessageDuration   *prometheus.HistogramVec
	PublishTotal      *prometheus.CounterVec
	ConsumeTotal      *prometheus.CounterVec
	ErrorTotal        *prometheus.CounterVec
}

// PostgreSQLMetrics provides PostgreSQL-related metrics.
type PostgreSQLMetrics struct {
	ConnectionPoolSize *prometheus.GaugeVec
	ConnectionPoolUsed *prometheus.GaugeVec
	QueryDuration      *prometheus.HistogramVec
	QueryTotal         *prometheus.CounterVec
	ErrorTotal         *prometheus.CounterVec
}

// InfluxDBMetrics provides InfluxDB-related metrics.
type InfluxDBMetrics struct {
	WriteDuration *prometheus.HistogramVec
	WriteTotal    *prometheus.CounterVec
	ErrorTotal    *prometheus.CounterVec
}

// NewMiddlewareMetrics creates middleware metrics.
func NewMiddlewareMetrics(collector *Collector) *MiddlewareMetrics {
	return &MiddlewareMetrics{
		RabbitMQ:   newRabbitMQMetrics(collector),
		PostgreSQL: newPostgreSQLMetrics(collector),
		InfluxDB:   newInfluxDBMetrics(collector),
	}
}

func newRabbitMQMetrics(collector *Collector) *RabbitMQMetrics {
	return &RabbitMQMetrics{
		ConnectionTotal: collector.NewGauge(
			"rabbitmq_connection_total",
			"Total number of RabbitMQ connections",
			[]string{LabelService},
		),
		MessageTotal: collector.NewCounter(
			"rabbitmq_message_total",
			"Total number of RabbitMQ messages",
			[]string{LabelService, LabelStatus},
		),
		MessageDuration: collector.NewHistogram(
			"rabbitmq_message_duration_seconds",
			"RabbitMQ message processing duration in seconds",
			[]string{LabelService},
			[]float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		),
		PublishTotal: collector.NewCounter(
			"rabbitmq_publish_total",
			"Total number of published messages",
			[]string{LabelService, LabelStatus},
		),
		ConsumeTotal: collector.NewCounter(
			"rabbitmq_consume_total",
			"Total number of consumed messages",
			[]string{LabelService, LabelStatus},
		),
		ErrorTotal: collector.NewCounter(
			"rabbitmq_error_total",
			"Total number of RabbitMQ errors",
			[]string{LabelService},
		),
	}
}

func newPostgreSQLMetrics(collector *Collector) *PostgreSQLMetrics {
	return &PostgreSQLMetrics{
		ConnectionPoolSize: collector.NewGauge(
			"postgres_connection_pool_size",
			"PostgreSQL connection pool size",
			[]string{LabelService},
		),
		ConnectionPoolUsed: collector.NewGauge(
			"postgres_connection_pool_used",
			"PostgreSQL connections in use",
			[]string{LabelService},
		),
		QueryDuration: collector.NewHistogram(
			"postgres_query_duration_seconds",
			"PostgreSQL query duration in seconds",
			[]string{LabelService},
			[]float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		),
		QueryTotal: collector.NewCounter(
			"postgres_query_total",
			"Total number of PostgreSQL queries",
			[]string{LabelService, LabelStatus},
		),
		ErrorTotal: collector.NewCounter(
			"postgres_error_total",
			"Total number of PostgreSQL errors",
			[]string{LabelService},
		),
	}
}

func newInfluxDBMetrics(collector *Collector) *InfluxDBMetrics {
	return &InfluxDBMetrics{
		WriteDuration: collector.NewHistogram(
			"influxdb_write_duration_seconds",
			"InfluxDB write duration in seconds",
			[]string{LabelService},
			[]float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		),
		WriteTotal: collector.NewCounter(
			"influxdb_write_total",
			"Total number of InfluxDB writes",
			[]string{LabelService, LabelStatus},
		),
		ErrorTotal: collector.NewCounter(
			"influxdb_error_total",
			"Total number of InfluxDB errors",
			[]string{LabelService},
		),
	}
}
