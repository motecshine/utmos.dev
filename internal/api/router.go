// Package api provides the HTTP API for iot-api service
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/utmos/utmos/internal/api/handler"
	"github.com/utmos/utmos/internal/api/middleware"
	"github.com/utmos/utmos/internal/downlink/dispatcher"
	"github.com/utmos/utmos/pkg/metrics"

	// Import swagger docs
	_ "github.com/utmos/utmos/docs/swagger"
)

// Config holds router configuration
type Config struct {
	// APIKeys for authentication
	APIKeys []string
	// EnableAuth enables authentication middleware
	EnableAuth bool
	// EnableTrace enables tracing middleware
	EnableTrace bool
	// ServiceName for tracing
	ServiceName string
	// TelemetryConfig for telemetry handler
	TelemetryConfig *handler.TelemetryConfig
}

// DefaultConfig returns default router configuration
func DefaultConfig() *Config {
	return &Config{
		EnableAuth:  true,
		EnableTrace: true,
		ServiceName: "iot-api",
	}
}

// Router wraps gin.Engine with additional functionality
type Router struct {
	engine           *gin.Engine
	config           *Config
	logger           *logrus.Entry
	db               *gorm.DB
	deviceHandler    *handler.Device
	serviceHandler   *handler.Service
	telemetryHandler *handler.Telemetry
}

// NewRouter creates a new API router
func NewRouter(
	config *Config,
	db *gorm.DB,
	dispatchHandler *dispatcher.DispatchHandler,
	metricsCollector *metrics.Collector,
	logger *logrus.Entry,
) *Router {
	if config == nil {
		config = DefaultConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	// Add recovery middleware
	engine.Use(gin.Recovery())

	// Add trace middleware
	if config.EnableTrace {
		traceMiddleware := middleware.NewTraceMiddleware(&middleware.TraceConfig{
			ServiceName: config.ServiceName,
			TracerName:  config.ServiceName,
			SkipPaths:   []string{"/health", "/ready", "/metrics"},
		}, logger)
		engine.Use(traceMiddleware.Handler())
	} else {
		// At minimum, add request logging
		engine.Use(middleware.RequestLogger(logger))
	}

	// Create handlers
	deviceHandler := handler.NewDevice(db, logger)
	serviceHandler := handler.NewService(db, dispatchHandler, logger)

	var telemetryHandler *handler.Telemetry
	if config.TelemetryConfig != nil {
		telemetryHandler = handler.NewTelemetry(config.TelemetryConfig, logger)
	}

	router := &Router{
		engine:           engine,
		config:           config,
		logger:           logger.WithField("component", "router"),
		db:               db,
		deviceHandler:    deviceHandler,
		serviceHandler:   serviceHandler,
		telemetryHandler: telemetryHandler,
	}

	// Setup routes
	router.setupHealthRoutes()
	router.setupMetricsRoute(metricsCollector)
	router.setupAPIRoutes()

	return router
}

// setupHealthRoutes sets up health check routes
func (r *Router) setupHealthRoutes() {
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.engine.GET("/ready", func(c *gin.Context) {
		checks := make(map[string]string)
		allReady := true

		// Check database connection
		if r.db != nil {
			sqlDB, err := r.db.DB()
			if err != nil {
				checks["database"] = "error: " + err.Error()
				allReady = false
			} else if err := sqlDB.Ping(); err != nil {
				checks["database"] = "error: " + err.Error()
				allReady = false
			} else {
				checks["database"] = "ok"
			}
		} else {
			checks["database"] = "not configured"
		}

		// Check telemetry handler (InfluxDB)
		if r.telemetryHandler != nil {
			checks["telemetry"] = "ok"
		} else {
			checks["telemetry"] = "not configured"
		}

		if allReady {
			c.JSON(http.StatusOK, gin.H{
				"status": "ready",
				"checks": checks,
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"checks": checks,
			})
		}
	})

	// Swagger documentation endpoint
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// setupMetricsRoute sets up the metrics endpoint
func (r *Router) setupMetricsRoute(collector *metrics.Collector) {
	if collector != nil {
		r.engine.GET("/metrics", metrics.Handler(collector))
	}
}

// setupAPIRoutes sets up API routes
func (r *Router) setupAPIRoutes() {
	api := r.engine.Group("/api/v1")

	// Add auth middleware if enabled
	if r.config.EnableAuth && len(r.config.APIKeys) > 0 {
		authMiddleware := middleware.NewAuthMiddleware(&middleware.AuthConfig{
			APIKeys:    r.config.APIKeys,
			HeaderName: "X-API-Key",
			SkipPaths:  []string{},
		}, r.logger)
		api.Use(authMiddleware.Handler())
	}

	// Device routes
	devices := api.Group("/devices")
	{
		devices.POST("", r.deviceHandler.Create)
		devices.GET("", r.deviceHandler.List)
		devices.GET("/:id", r.deviceHandler.Get)
		devices.GET("/sn/:sn", r.deviceHandler.GetBySN)
		devices.PUT("/:id", r.deviceHandler.Update)
		devices.DELETE("/:id", r.deviceHandler.Delete)
	}

	// Service call routes
	services := api.Group("/services")
	{
		services.POST("/call", r.serviceHandler.Call)
		services.GET("/calls/:id", r.serviceHandler.Get)
		services.GET("/calls/device/:device_sn", r.serviceHandler.ListByDevice)
		services.POST("/calls/:id/cancel", r.serviceHandler.Cancel)

		// Note: Vendor-specific routes (e.g., /dji/takeoff) have been removed.
		// Use the generic /call endpoint with vendor and method parameters instead.
	}

	// Telemetry routes
	if r.telemetryHandler != nil {
		telemetry := api.Group("/telemetry")
		{
			telemetry.GET("/:device_sn", r.telemetryHandler.Query)
			telemetry.GET("/:device_sn/latest", r.telemetryHandler.Latest)
			telemetry.GET("/:device_sn/aggregate", r.telemetryHandler.Aggregate)
		}
	}
}

// Engine returns the underlying gin.Engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}

// Close closes any resources held by the router
func (r *Router) Close() {
	if r.telemetryHandler != nil {
		r.telemetryHandler.Close()
	}
}
