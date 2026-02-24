package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/utmos/utmos/internal/downlink"
	"github.com/utmos/utmos/internal/downlink/retry"
	"github.com/utmos/utmos/internal/downlink/router"
	"github.com/utmos/utmos/internal/shared/config"
	"github.com/utmos/utmos/internal/shared/server"
	"github.com/utmos/utmos/pkg/logger"
	djidownlink "github.com/utmos/utmos/pkg/adapter/dji/downlink"
	"github.com/utmos/utmos/pkg/metrics"
	"github.com/utmos/utmos/pkg/rabbitmq"
	"github.com/utmos/utmos/pkg/tracer"
)

const serviceName = "iot-downlink"

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

	// Initialize downlink service
	downlinkConfig := &downlink.Config{
		RetryConfig: &retry.Config{
			MaxRetries:       3,
			InitialDelay:     time.Second,
			MaxDelay:         30 * time.Second,
			Multiplier:       2.0,
			EnableDeadLetter: true,
		},
		RouterConfig: &router.Config{
			DefaultRoutingKey: router.RoutingKeyGatewayDownlink,
			EnableMetrics:     true,
		},
		EnableRetry:         true,
		EnableRouting:       true,
		RetryWorkerInterval: 5 * time.Second,
	}

	downlinkService := downlink.NewService(downlinkConfig, publisher, log.WithService(serviceName))

	// Register DJI dispatcher
	djiDispatcher := djidownlink.NewDispatcherAdapter(publisher, log.WithService(serviceName))
	downlinkService.RegisterAdapterDispatcher(djiDispatcher)

	// Start downlink service
	if err := downlinkService.Start(context.Background()); err != nil {
		log.WithService(serviceName).Fatalf("failed to start downlink service: %v", err)
	}

	// Setup Gin router for health checks
	if cfg.Logger.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	router.GET("/ready", func(c *gin.Context) {
		if rmqClient.IsConnected() && downlinkService.IsRunning() {
			c.JSON(http.StatusOK, gin.H{"status": "ready"})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
	})

	// Metrics endpoint
	router.GET(cfg.Metrics.Path, metrics.Handler(metricsCollector))

	// Create HTTP server for health checks
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Setup graceful shutdown
	shutdown := server.NewGracefulShutdown(30 * time.Second)
	shutdown.Register(func(ctx context.Context) error {
		log.WithService(serviceName).Info("Shutting down HTTP server")
		return srv.Shutdown(ctx)
	})
	shutdown.Register(func(_ context.Context) error {
		log.WithService(serviceName).Info("Stopping downlink service")
		return downlinkService.Stop()
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

	// Log service status
	vendors := downlinkService.GetRegisteredVendors()
	log.WithService(serviceName).Infof("Downlink service ready with vendors: %v", vendors)

	// Wait for shutdown signal
	if err := shutdown.Wait(); err != nil {
		log.WithService(serviceName).Errorf("shutdown error: %v", err)
	}
	log.WithService(serviceName).Info("Service stopped")
}
