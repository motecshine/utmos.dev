package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/sirupsen/logrus"
)

// Telemetry handles telemetry query API requests
type Telemetry struct {
	client   influxdb2.Client
	queryAPI api.QueryAPI
	org      string
	bucket   string
	logger   *logrus.Entry
}

// TelemetryConfig holds InfluxDB configuration for telemetry handler
type TelemetryConfig struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

// NewTelemetry creates a new telemetry handler
func NewTelemetry(config *TelemetryConfig, logger *logrus.Entry) *Telemetry {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	var client influxdb2.Client
	var queryAPI api.QueryAPI

	if config != nil && config.URL != "" {
		client = influxdb2.NewClient(config.URL, config.Token)
		queryAPI = client.QueryAPI(config.Org)
	}

	return &Telemetry{
		client:   client,
		queryAPI: queryAPI,
		org:      config.Org,
		bucket:   config.Bucket,
		logger:   logger.WithField("handler", "telemetry"),
	}
}

// TelemetryPoint represents a single telemetry data point
type TelemetryPoint struct {
	Time   string            `json:"time"`
	Fields map[string]any    `json:"fields"`
	Tags   map[string]string `json:"tags,omitempty"`
}

// TelemetryQueryResponse represents the response for telemetry queries
type TelemetryQueryResponse struct {
	DeviceSN string           `json:"device_sn"`
	Points   []TelemetryPoint `json:"points"`
	Total    int              `json:"total"`
}

// LatestTelemetryResponse represents the response for latest telemetry
type LatestTelemetryResponse struct {
	DeviceSN  string         `json:"device_sn"`
	Timestamp string         `json:"timestamp"`
	Data      map[string]any `json:"data"`
}

// requireQueryAPI checks that the InfluxDB query API is available.
// Returns true if available. On failure it writes a 503 response and returns false.
func (h *Telemetry) requireQueryAPI(c *gin.Context) bool {
	return requireDependency(c, h.queryAPI, "Telemetry service")
}

// requireDeviceSN validates the device_sn URL parameter.
// Returns the value and true on success.
func requireDeviceSN(c *gin.Context) (string, bool) {
	return requireStringParam(c, "device_sn", "INVALID_DEVICE_SN", "Device serial number is required")
}

// Query queries telemetry data for a device
// @Summary Query telemetry data
// @Description Query telemetry data for a device within a time range
// @Tags telemetry
// @Produce json
// @Param device_sn path string true "Device Serial Number"
// @Param start query string false "Start time (RFC3339)" default(-1h)
// @Param stop query string false "Stop time (RFC3339)" default(now)
// @Param measurement query string false "Measurement name"
// @Param limit query int false "Limit" default(100)
// @Success 200 {object} TelemetryQueryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/telemetry/{device_sn} [get]
func (h *Telemetry) Query(c *gin.Context) {
	deviceSN, ok := requireDeviceSN(c)
	if !ok {
		return
	}

	if !h.requireQueryAPI(c) {
		return
	}

	// Parse query parameters
	start := c.DefaultQuery("start", "-1h")
	stop := c.DefaultQuery("stop", "now()")
	measurement := c.Query("measurement")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	if limit < 1 || limit > 1000 {
		limit = 100
	}

	// Build Flux query
	query := h.buildQuery(deviceSN, start, stop, measurement, limit)

	// Execute query
	result, err := h.queryAPI.Query(context.Background(), query)
	if err != nil {
		respondInternalError(c, h.logger, err, "Failed to query telemetry", "Failed to query telemetry data")
		return
	}
	defer func() { _ = result.Close() }()

	// Parse results
	points := make([]TelemetryPoint, 0)
	for result.Next() {
		record := result.Record()
		point := TelemetryPoint{
			Time: record.Time().Format(time.RFC3339),
			Fields: map[string]any{
				record.Field(): record.Value(),
			},
			Tags: make(map[string]string),
		}

		// Add tags
		for k, v := range record.Values() {
			if k != "_time" && k != "_value" && k != "_field" && k != "_measurement" {
				if str, ok := v.(string); ok {
					point.Tags[k] = str
				}
			}
		}

		points = append(points, point)
	}

	if result.Err() != nil {
		h.logger.WithError(result.Err()).Error("Error reading telemetry results")
	}

	c.JSON(http.StatusOK, TelemetryQueryResponse{
		DeviceSN: deviceSN,
		Points:   points,
		Total:    len(points),
	})
}

// Latest gets the latest telemetry data for a device
// @Summary Get latest telemetry
// @Description Get the latest telemetry data for a device
// @Tags telemetry
// @Produce json
// @Param device_sn path string true "Device Serial Number"
// @Success 200 {object} LatestTelemetryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/telemetry/{device_sn}/latest [get]
func (h *Telemetry) Latest(c *gin.Context) {
	deviceSN, ok := requireDeviceSN(c)
	if !ok {
		return
	}

	if !h.requireQueryAPI(c) {
		return
	}

	// Build query for latest data
	query := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: -1h)
			|> filter(fn: (r) => r["device_sn"] == "%s")
			|> last()
			|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
	`, h.bucket, deviceSN)

	result, err := h.queryAPI.Query(context.Background(), query)
	if err != nil {
		respondInternalError(c, h.logger, err, "Failed to query latest telemetry", "Failed to query telemetry data")
		return
	}
	defer func() { _ = result.Close() }()

	// Get the first (and should be only) result
	if result.Next() {
		record := result.Record()
		data := make(map[string]any)

		for k, v := range record.Values() {
			if k != "_time" && k != "_start" && k != "_stop" && k != "_measurement" && k != "result" && k != "table" {
				data[k] = v
			}
		}

		c.JSON(http.StatusOK, LatestTelemetryResponse{
			DeviceSN:  deviceSN,
			Timestamp: record.Time().Format(time.RFC3339),
			Data:      data,
		})
		return
	}

	if result.Err() != nil {
		h.logger.WithError(result.Err()).Error("Error reading telemetry results")
	}

	respondNotFound(c, "NO_DATA", "No telemetry data found for device")
}

// Aggregate queries aggregated telemetry data
// @Summary Query aggregated telemetry
// @Description Query aggregated telemetry data for a device
// @Tags telemetry
// @Produce json
// @Param device_sn path string true "Device Serial Number"
// @Param start query string false "Start time (RFC3339)" default(-24h)
// @Param stop query string false "Stop time (RFC3339)" default(now)
// @Param window query string false "Aggregation window" default(1h)
// @Param fn query string false "Aggregation function (mean, max, min, sum)" default(mean)
// @Param field query string true "Field to aggregate"
// @Success 200 {object} TelemetryQueryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/telemetry/{device_sn}/aggregate [get]
func (h *Telemetry) Aggregate(c *gin.Context) {
	deviceSN, ok := requireDeviceSN(c)
	if !ok {
		return
	}

	field := c.Query("field")
	if field == "" {
		respondBadRequest(c, "INVALID_FIELD", "Field parameter is required")
		return
	}

	if !h.requireQueryAPI(c) {
		return
	}

	start := c.DefaultQuery("start", "-24h")
	stop := c.DefaultQuery("stop", "now()")
	window := c.DefaultQuery("window", "1h")
	fn := c.DefaultQuery("fn", "mean")

	// Validate aggregation function
	validFns := map[string]bool{"mean": true, "max": true, "min": true, "sum": true, "count": true}
	if !validFns[fn] {
		respondBadRequest(c, "INVALID_FUNCTION", "Invalid aggregation function. Use: mean, max, min, sum, count")
		return
	}

	// Build aggregation query
	query := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["device_sn"] == "%s")
			|> filter(fn: (r) => r["_field"] == "%s")
			|> aggregateWindow(every: %s, fn: %s, createEmpty: false)
	`, h.bucket, start, stop, deviceSN, field, window, fn)

	result, err := h.queryAPI.Query(context.Background(), query)
	if err != nil {
		respondInternalError(c, h.logger, err, "Failed to query aggregated telemetry", "Failed to query telemetry data")
		return
	}
	defer func() { _ = result.Close() }()

	points := make([]TelemetryPoint, 0)
	for result.Next() {
		record := result.Record()
		point := TelemetryPoint{
			Time: record.Time().Format(time.RFC3339),
			Fields: map[string]any{
				field: record.Value(),
			},
		}
		points = append(points, point)
	}

	if result.Err() != nil {
		h.logger.WithError(result.Err()).Error("Error reading aggregated results")
	}

	c.JSON(http.StatusOK, TelemetryQueryResponse{
		DeviceSN: deviceSN,
		Points:   points,
		Total:    len(points),
	})
}

// buildQuery builds a Flux query for telemetry data
func (h *Telemetry) buildQuery(deviceSN, start, stop, measurement string, limit int) string {
	query := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["device_sn"] == "%s")
	`, h.bucket, start, stop, deviceSN)

	if measurement != "" {
		query += fmt.Sprintf(`|> filter(fn: (r) => r["_measurement"] == "%s")
		`, measurement)
	}

	query += fmt.Sprintf(`|> limit(n: %d)`, limit)

	return query
}

// Close closes the InfluxDB client
func (h *Telemetry) Close() {
	if h.client != nil {
		h.client.Close()
	}
}
