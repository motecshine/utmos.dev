package observability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/pkg/metrics"
)

func TestNewMetrics(t *testing.T) {
	collector := metrics.NewCollector("test")
	m := NewMetrics(collector)

	require.NotNil(t, m)
	assert.NotNil(t, m.MessagesReceived)
	assert.NotNil(t, m.MessagesSent)
	assert.NotNil(t, m.MessagesErrors)
	assert.NotNil(t, m.ProcessingDuration)
	assert.NotNil(t, m.ActiveDevices)
}

func TestMetrics_RecordMessageReceived(t *testing.T) {
	collector := metrics.NewCollector("test")
	m := NewMetrics(collector)

	// Should not panic
	m.RecordMessageReceived("osd", "success")
	m.RecordMessageReceived("state", "success")
	m.RecordMessageReceived("services", "error")
}

func TestMetrics_RecordMessageSent(t *testing.T) {
	collector := metrics.NewCollector("test")
	m := NewMetrics(collector)

	// Should not panic
	m.RecordMessageSent("services", "success")
	m.RecordMessageSent("events_reply", "success")
}

func TestMetrics_RecordError(t *testing.T) {
	collector := metrics.NewCollector("test")
	m := NewMetrics(collector)

	// Should not panic
	m.RecordError("osd", "parse_error")
	m.RecordError("services", "timeout")
}

func TestMetrics_RecordProcessingDuration(t *testing.T) {
	collector := metrics.NewCollector("test")
	m := NewMetrics(collector)

	// Should not panic
	m.RecordProcessingDuration("osd", 0.001)
	m.RecordProcessingDuration("services", 0.05)
}

func TestMetrics_SetActiveDevices(t *testing.T) {
	collector := metrics.NewCollector("test")
	m := NewMetrics(collector)

	// Should not panic
	m.SetActiveDevices("dock", 10)
	m.SetActiveDevices("aircraft", 5)
}

func TestTracer_StartMessageSpan(t *testing.T) {
	tracer := NewTracer()
	ctx := context.Background()

	ctx, span := tracer.StartMessageSpan(ctx, "osd", "property.report", "device-001")
	defer span.End()

	assert.NotNil(t, ctx)
	assert.NotNil(t, span)
}

func TestTracer_StartServiceCallSpan(t *testing.T) {
	tracer := NewTracer()
	ctx := context.Background()

	ctx, span := tracer.StartServiceCallSpan(ctx, "flighttask_prepare", "device-001")
	defer span.End()

	assert.NotNil(t, ctx)
	assert.NotNil(t, span)
}

func TestTracer_StartEventSpan(t *testing.T) {
	tracer := NewTracer()
	ctx := context.Background()

	ctx, span := tracer.StartEventSpan(ctx, "flighttask_progress", "device-001")
	defer span.End()

	assert.NotNil(t, ctx)
	assert.NotNil(t, span)
}

func TestGetTraceID_NoSpan(t *testing.T) {
	ctx := context.Background()
	traceID := GetTraceID(ctx)
	assert.Empty(t, traceID)
}

func TestGetSpanID_NoSpan(t *testing.T) {
	ctx := context.Background()
	spanID := GetSpanID(ctx)
	assert.Empty(t, spanID)
}
