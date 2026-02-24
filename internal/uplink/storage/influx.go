// Package storage provides data storage functionality for iot-uplink
package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/sirupsen/logrus"

	"github.com/utmos/utmos/pkg/adapter"
)

// Config holds InfluxDB configuration
type Config struct {
	URL           string
	Token         string
	Org           string
	Bucket        string
	BatchSize     int
	FlushInterval time.Duration
}

// DefaultConfig returns default InfluxDB configuration
func DefaultConfig() *Config {
	return &Config{
		URL:           "http://localhost:8086",
		Token:         "",
		Org:           "utmos",
		Bucket:        "iot",
		BatchSize:     1000,
		FlushInterval: time.Second,
	}
}

// Storage provides InfluxDB storage functionality
type Storage struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
	config   *Config
	logger   *logrus.Entry
	mu       sync.RWMutex
	closed   bool
}

// NewStorage creates a new InfluxDB storage
func NewStorage(config *Config, logger *logrus.Entry) *Storage {
	if config == nil {
		config = DefaultConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	// Create InfluxDB client with options
	options := influxdb2.DefaultOptions().
		SetBatchSize(uint(config.BatchSize)).
		SetFlushInterval(uint(config.FlushInterval.Milliseconds()))

	client := influxdb2.NewClientWithOptions(config.URL, config.Token, options)

	// Get non-blocking write API
	writeAPI := client.WriteAPI(config.Org, config.Bucket)

	// Setup error handler
	errorsCh := writeAPI.Errors()
	go func() {
		for err := range errorsCh {
			logger.WithError(err).Error("InfluxDB write error")
		}
	}()

	return &Storage{
		client:   client,
		writeAPI: writeAPI,
		config:   config,
		logger:   logger.WithField("component", "influx-storage"),
	}
}

// WriteProcessedMessage writes a processed message to InfluxDB
func (s *Storage) WriteProcessedMessage(ctx context.Context, msg *adapter.ProcessedMessage) error {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return fmt.Errorf("storage is closed")
	}
	s.mu.RUnlock()

	if msg == nil {
		return fmt.Errorf("message is nil")
	}

	// Create measurement name based on message type
	measurement := s.getMeasurement(msg)

	// Create tags
	tags := map[string]string{
		"device_sn": msg.DeviceSN,
		"vendor":    msg.Vendor,
	}

	// Convert timestamp
	timestamp := time.UnixMilli(msg.Timestamp)

	// Write properties as fields
	if len(msg.Properties) > 0 {
		fields := s.convertToFields(msg.Properties)
		if len(fields) > 0 {
			point := write.NewPoint(measurement, tags, fields, timestamp)
			s.writeAPI.WritePoint(point)
		}
	}

	// Write events
	for _, event := range msg.Events {
		eventTags := make(map[string]string)
		for k, v := range tags {
			eventTags[k] = v
		}
		eventTags["event_name"] = event.Name

		fields := s.convertToFields(event.Params)
		if len(fields) == 0 {
			fields["event"] = event.Name
		}

		point := write.NewPoint("events", eventTags, fields, timestamp)
		s.writeAPI.WritePoint(point)
	}

	s.logger.WithFields(logrus.Fields{
		"device_sn":   msg.DeviceSN,
		"measurement": measurement,
		"properties":  len(msg.Properties),
		"events":      len(msg.Events),
	}).Debug("Wrote message to InfluxDB")

	return nil
}

// WriteTelemetry writes telemetry data to InfluxDB
func (s *Storage) WriteTelemetry(ctx context.Context, deviceSN, vendor, measurement string, fields map[string]any, tags map[string]string, timestamp time.Time) error {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return fmt.Errorf("storage is closed")
	}
	s.mu.RUnlock()

	// Merge default tags
	allTags := map[string]string{
		"device_sn": deviceSN,
		"vendor":    vendor,
	}
	for k, v := range tags {
		allTags[k] = v
	}

	// Convert fields to proper types
	convertedFields := s.convertToFields(fields)
	if len(convertedFields) == 0 {
		return fmt.Errorf("no valid fields to write")
	}

	point := write.NewPoint(measurement, allTags, convertedFields, timestamp)
	s.writeAPI.WritePoint(point)

	return nil
}

// WritePoint writes a single point to InfluxDB
func (s *Storage) WritePoint(point *write.Point) error {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return fmt.Errorf("storage is closed")
	}
	s.mu.RUnlock()

	s.writeAPI.WritePoint(point)
	return nil
}

// Flush forces a flush of buffered data
func (s *Storage) Flush() {
	s.writeAPI.Flush()
}

// Close closes the InfluxDB connection
func (s *Storage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	s.writeAPI.Flush()
	s.client.Close()
	s.logger.Info("InfluxDB storage closed")

	return nil
}

// Health checks the InfluxDB connection health
func (s *Storage) Health(ctx context.Context) error {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return fmt.Errorf("storage is closed")
	}
	s.mu.RUnlock()

	health, err := s.client.Health(ctx)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	if health.Status != "pass" {
		return fmt.Errorf("unhealthy status: %s", health.Status)
	}

	return nil
}

// getMeasurement returns the measurement name for a message type
func (s *Storage) getMeasurement(msg *adapter.ProcessedMessage) string {
	switch msg.MessageType {
	case adapter.MessageTypeProperty:
		return "telemetry"
	case adapter.MessageTypeEvent:
		return "events"
	case adapter.MessageTypeService:
		return "services"
	case adapter.MessageTypeStatus:
		return "status"
	default:
		return "telemetry"
	}
}

// convertToFields converts a map to InfluxDB fields, filtering out unsupported types
func (s *Storage) convertToFields(data map[string]any) map[string]any {
	fields := make(map[string]any)

	for key, value := range data {
		switch v := value.(type) {
		case float64, float32:
			fields[key] = v
		case int, int8, int16, int32, int64:
			fields[key] = v
		case uint, uint8, uint16, uint32, uint64:
			fields[key] = v
		case string:
			fields[key] = v
		case bool:
			fields[key] = v
		case nil:
			// Skip nil values
		default:
			// For complex types, convert to string
			fields[key] = fmt.Sprintf("%v", v)
		}
	}

	return fields
}

// TelemetryPoint represents a telemetry data point
type TelemetryPoint struct {
	DeviceSN    string
	Vendor      string
	Measurement string
	Tags        map[string]string
	Fields      map[string]any
	Timestamp   time.Time
}

// NewTelemetryPoint creates a new telemetry point
func NewTelemetryPoint(deviceSN, vendor, measurement string) *TelemetryPoint {
	return &TelemetryPoint{
		DeviceSN:    deviceSN,
		Vendor:      vendor,
		Measurement: measurement,
		Tags:        make(map[string]string),
		Fields:      make(map[string]any),
		Timestamp:   time.Now(),
	}
}

// AddTag adds a tag to the point
func (p *TelemetryPoint) AddTag(key, value string) *TelemetryPoint {
	p.Tags[key] = value
	return p
}

// AddField adds a field to the point
func (p *TelemetryPoint) AddField(key string, value any) *TelemetryPoint {
	p.Fields[key] = value
	return p
}

// SetTimestamp sets the timestamp
func (p *TelemetryPoint) SetTimestamp(t time.Time) *TelemetryPoint {
	p.Timestamp = t
	return p
}

// ToInfluxPoint converts to an InfluxDB point
func (p *TelemetryPoint) ToInfluxPoint() *write.Point {
	tags := map[string]string{
		"device_sn": p.DeviceSN,
		"vendor":    p.Vendor,
	}
	for k, v := range p.Tags {
		tags[k] = v
	}

	return write.NewPoint(p.Measurement, tags, p.Fields, p.Timestamp)
}
