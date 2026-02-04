// Package metrics provides unified Prometheus metrics collection.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Standard label names
const (
	LabelService     = "service"
	LabelVendor      = "vendor"
	LabelMessageType = "message_type"
	LabelStatus      = "status"
	LabelMethod      = "method"
	LabelPath        = "path"
	LabelCode        = "code"
)

// Collector provides metrics collection and registration.
type Collector struct {
	registry  *prometheus.Registry
	namespace string
}

// NewCollector creates a new metrics collector.
func NewCollector(namespace string) *Collector {
	registry := prometheus.NewRegistry()
	// Register default collectors
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGoCollector())

	return &Collector{
		registry:  registry,
		namespace: namespace,
	}
}

// Registry returns the Prometheus registry.
func (c *Collector) Registry() *prometheus.Registry {
	return c.registry
}

// NewCounter creates and registers a new counter.
func (c *Collector) NewCounter(name, help string, labels []string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: c.namespace,
			Name:      name,
			Help:      help,
		},
		labels,
	)
	c.registry.MustRegister(counter)
	return counter
}

// NewHistogram creates and registers a new histogram.
func (c *Collector) NewHistogram(name, help string, labels []string, buckets []float64) *prometheus.HistogramVec {
	if buckets == nil {
		buckets = prometheus.DefBuckets
	}
	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: c.namespace,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		},
		labels,
	)
	c.registry.MustRegister(histogram)
	return histogram
}

// NewGauge creates and registers a new gauge.
func (c *Collector) NewGauge(name, help string, labels []string) *prometheus.GaugeVec {
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: c.namespace,
			Name:      name,
			Help:      help,
		},
		labels,
	)
	c.registry.MustRegister(gauge)
	return gauge
}

// NewSummary creates and registers a new summary.
func (c *Collector) NewSummary(name, help string, labels []string) *prometheus.SummaryVec {
	summary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: c.namespace,
			Name:      name,
			Help:      help,
		},
		labels,
	)
	c.registry.MustRegister(summary)
	return summary
}
