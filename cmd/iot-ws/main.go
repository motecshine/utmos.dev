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
	"github.com/utmos/utmos/pkg/logger"
	"github.com/utmos/utmos/internal/ws"
	"github.com/utmos/utmos/internal/ws/hub"
	"github.com/utmos/utmos/internal/ws/push"
	"github.com/utmos/utmos/pkg/metrics"
	"github.com/utmos/utmos/pkg/rabbitmq"
	"github.com/utmos/utmos/pkg/tracer"
)

const serviceName = "iot-ws"

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

	// Setup Exchange
	if rmqClient.IsConnected() {
		if err := rmqClient.SetupExchange(); err != nil {
			log.WithService(serviceName).Warnf("failed to setup exchange: %v", err)
		}
	}

	// Initialize RabbitMQ subscriber
	subscriber := rabbitmq.NewSubscriber(rmqClient)

	// Create WebSocket service configuration
	wsConfig := &ws.ServiceConfig{
		HubConfig: &hub.Config{
			MaxConnections: 10000,
			WriteTimeout:   10 * time.Second,
			ReadTimeout:    60 * time.Second,
			PingInterval:   30 * time.Second,
			PongTimeout:    30 * time.Second,
		},
		ClientConfig: hub.DefaultClientConfig(),
		PusherConfig: &push.Config{
			WorkerCount: 4,
			QueueSize:   10000,
		},
		AllowedOrigins: []string{"*"},
	}

	// Create WebSocket service
	wsSvc := ws.NewService(wsConfig, metricsCollector, log.WithService(serviceName))
	wsSvc.SetSubscriber(subscriber)

	// Start WebSocket service
	if err := wsSvc.Start(context.Background()); err != nil {
		log.WithService(serviceName).Fatalf("failed to start WebSocket service: %v", err)
	}

	// Setup Gin router
	if cfg.Logger.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(tracer.HTTPMiddleware(tracerProvider.Tracer(serviceName)))

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	router.GET("/ready", func(c *gin.Context) {
		if wsSvc.IsRunning() {
			c.JSON(http.StatusOK, gin.H{"status": "ready"})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
	})

	// Metrics endpoint
	router.GET(cfg.Metrics.Path, metrics.Handler(metricsCollector))

	// Stats endpoint
	router.GET("/stats", func(c *gin.Context) {
		stats := wsSvc.GetStats()
		c.JSON(http.StatusOK, stats)
	})

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		wsSvc.HandleWebSocket(c.Writer, c.Request)
	})

	// Create HTTP server
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
		log.WithService(serviceName).Info("Stopping WebSocket service")
		return wsSvc.Stop()
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

	// Start server
	go func() {
		log.WithService(serviceName).Infof("Server listening on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithService(serviceName).Fatalf("failed to start server: %v", err)
		}
	}()

	// Log startup info
	log.WithService(serviceName).Info("WebSocket service ready")

	// Wait for shutdown signal
	if err := shutdown.Wait(); err != nil {
		log.WithService(serviceName).Errorf("shutdown error: %v", err)
	}
	log.WithService(serviceName).Info("Service stopped")
}
