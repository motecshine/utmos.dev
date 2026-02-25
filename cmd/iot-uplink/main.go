package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/utmos/utmos/internal/shared/config"
	"github.com/utmos/utmos/internal/shared/server"
	"github.com/utmos/utmos/internal/uplink"
	"github.com/utmos/utmos/internal/uplink/router"
	"github.com/utmos/utmos/internal/uplink/storage"
	djiuplink "github.com/utmos/utmos/pkg/adapter/dji/uplink"
	"github.com/utmos/utmos/pkg/logger"
	"github.com/utmos/utmos/pkg/metrics"
	"github.com/utmos/utmos/pkg/rabbitmq"
	"github.com/utmos/utmos/pkg/tracer"
)

const serviceName = "iot-uplink"

func main() {
	// Load configuration
	cfg, err := config.LoadFromEnv("dev")
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	cfg.Tracer.ServiceName = serviceName

	// Initialize logger
	log := logger.New(&cfg.Logger)
	log.WithService(serviceName).Info("Starting service")

	// Initialize tracer
	tracerProvider, err := tracer.NewProvider(&cfg.Tracer)
	if err != nil {
		log.WithService(serviceName).Fatalf("failed to initialize tracer: %v", err)
	}

	// Initialize metrics collector
	metricsCollector := metrics.NewCollector(cfg.Metrics.Namespace)

	// Initialize RabbitMQ client
	rmqClient := rabbitmq.NewClient(&cfg.RabbitMQ)
	if err := rmqClient.Connect(context.Background()); err != nil {
		log.WithService(serviceName).Warnf("failed to connect to RabbitMQ: %v", err)
	}

	// Setup Exchange and Queue
	if rmqClient.IsConnected() {
		if err := rmqClient.SetupExchange(); err != nil {
			log.WithService(serviceName).Warnf("failed to setup exchange: %v", err)
		}
	}

	// Initialize RabbitMQ subscriber and publisher
	subscriber := rabbitmq.NewSubscriber(rmqClient)
	publisher := rabbitmq.NewPublisher(rmqClient)

	// Create uplink service configuration
	influxConfig := &storage.Config{
		URL:           cfg.Database.InfluxDB.URL,
		Token:         cfg.Database.InfluxDB.Token,
		Org:           cfg.Database.InfluxDB.Org,
		Bucket:        cfg.Database.InfluxDB.Bucket,
		BatchSize:     1000,
		FlushInterval: time.Second,
	}

	routerConfig := &router.Config{
		Exchange:         cfg.RabbitMQ.ExchangeName,
		EnableWSRouting:  true,
		EnableAPIRouting: true,
	}

	svcConfig := &uplink.ServiceConfig{
		QueueName: "iot.uplink.messages",
		RoutingKeys: []string{
			"iot.dji.#",
			"iot.raw.*.uplink",
		},
		Influx:        influxConfig,
		Router:        routerConfig,
		EnableStorage: true,
		EnableRouting: true,
	}

	// Create uplink service
	logEntry := log.WithService(serviceName)
	uplinkSvc := uplink.NewService(svcConfig, subscriber, publisher, metricsCollector, logEntry)

	// Register DJI processor
	djiProcessor := djiuplink.NewProcessorAdapter(logEntry)
	uplinkSvc.RegisterProcessor(djiProcessor)

	// Setup Gin router for health checks
	if cfg.Logger.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())

	// Health check endpoints
	ginRouter.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	ginRouter.GET("/ready", func(c *gin.Context) {
		rmqReady := rmqClient.IsConnected()
		svcRunning := uplinkSvc.IsRunning()

		if rmqReady && svcRunning {
			stats := uplinkSvc.GetStats()
			c.JSON(http.StatusOK, gin.H{
				"status":            "ready",
				"rabbitmq":          "connected",
				"service":           "running",
				"registered_vendors": stats.RegisteredVendors,
				"storage_enabled":   stats.StorageEnabled,
				"routing_enabled":   stats.RoutingEnabled,
			})
			return
		}

		status := gin.H{
			"status":   "not ready",
			"rabbitmq": "disconnected",
			"service":  "stopped",
		}
		if rmqReady {
			status["rabbitmq"] = "connected"
		}
		if svcRunning {
			status["service"] = "running"
		}
		c.JSON(http.StatusServiceUnavailable, status)
	})

	// Metrics endpoint
	ginRouter.GET(cfg.Metrics.Path, metrics.Handler(metricsCollector))

	// Stats endpoint
	ginRouter.GET("/stats", func(c *gin.Context) {
		stats := uplinkSvc.GetStats()
		c.JSON(http.StatusOK, gin.H{
			"running":            stats.Running,
			"registered_vendors": stats.RegisteredVendors,
			"storage_enabled":    stats.StorageEnabled,
			"routing_enabled":    stats.RoutingEnabled,
		})
	})

	// Create HTTP server for health checks
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      ginRouter,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Setup graceful shutdown
	shutdown := server.NewGracefulShutdown(30 * time.Second)
	shutdown.Register(func(ctx context.Context) error {
		log.WithService(serviceName).Info("Shutting down HTTP server")
		return srv.Shutdown(ctx)
	})
	shutdown.Register(func(ctx context.Context) error {
		log.WithService(serviceName).Info("Stopping uplink service")
		return uplinkSvc.Stop()
	})
	shutdown.Register(func(_ context.Context) error {
		log.WithService(serviceName).Info("Stopping RabbitMQ subscriber")
		subscriber.UnsubscribeAll()
		return nil
	})
	shutdown.Register(func(_ context.Context) error {
		log.WithService(serviceName).Info("Closing RabbitMQ connection")
		return rmqClient.Close()
	})
	shutdown.Register(func(ctx context.Context) error {
		log.WithService(serviceName).Info("Shutting down tracer")
		return tracerProvider.Shutdown(ctx)
	})

	// Start HTTP server for health checks
	go func() {
		log.WithService(serviceName).Infof("Health check server listening on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithService(serviceName).Fatalf("failed to start server: %v", err)
		}
	}()

	// Start uplink service
	ctx := context.Background()
	if err := uplinkSvc.Start(ctx); err != nil {
		log.WithService(serviceName).Warnf("failed to start uplink service: %v", err)
		log.WithService(serviceName).Info("Uplink running in degraded mode")
	} else {
		log.WithService(serviceName).Info("Uplink service started - processing messages")
	}

	// Wait for shutdown signal
	if err := shutdown.Wait(); err != nil {
		log.WithService(serviceName).Errorf("shutdown error: %v", err)
	}
	log.WithService(serviceName).Info("Service stopped")
}
