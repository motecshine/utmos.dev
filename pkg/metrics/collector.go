// Package metrics provides unified Prometheus metrics collection.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())

	return &Collector{
		registry:  registry,
		namespace: namespace,
	}
}

// Registry returns the Prometheus registry.
func (c *Collector) Registry() *prometheus.Registry {
	return c.registry
}

// registerMetric registers a prometheus.Collector and returns it.
func registerMetric[T prometheus.Collector](c *Collector, metric T) T {
	c.registry.MustRegister(metric)
	return metric
}

// NewCounter creates and registers a new counter.
func (c *Collector) NewCounter(name, help string, labels []string) *prometheus.CounterVec {
	return registerMetric(c, prometheus.NewCounterVec(
		prometheus.CounterOpts{Namespace: c.namespace, Name: name, Help: help},
		labels,
	))
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
	return registerMetric(c, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Namespace: c.namespace, Name: name, Help: help},
		labels,
	))
}

// NewSummary creates and registers a new summary.
func (c *Collector) NewSummary(name, help string, labels []string) *prometheus.SummaryVec {
	return registerMetric(c, prometheus.NewSummaryVec(
		prometheus.SummaryOpts{Namespace: c.namespace, Name: name, Help: help},
		labels,
	))
}
