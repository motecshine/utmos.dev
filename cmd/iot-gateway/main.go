package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/utmos/utmos/internal/gateway"
	"github.com/utmos/utmos/internal/gateway/bridge"
	"github.com/utmos/utmos/internal/gateway/mqtt"
	"github.com/utmos/utmos/internal/shared/config"
	"github.com/utmos/utmos/internal/shared/server"
	"github.com/utmos/utmos/pkg/logger"
	"github.com/utmos/utmos/pkg/metrics"
	"github.com/utmos/utmos/pkg/rabbitmq"
	"github.com/utmos/utmos/pkg/tracer"
)

const serviceName = "iot-gateway"

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

	// Create gateway service configuration
	mqttConfig := &mqtt.Config{
		Broker:           cfg.MQTT.Broker,
		Port:             cfg.MQTT.Port,
		ClientID:         cfg.MQTT.ClientID,
		Username:         cfg.MQTT.Username,
		Password:         cfg.MQTT.Password,
		CleanSession:     cfg.MQTT.CleanSession,
		AutoReconnect:    cfg.MQTT.AutoReconnect,
		ConnectTimeout:   cfg.MQTT.ConnectTimeout,
		KeepAlive:        cfg.MQTT.KeepAlive,
		PingTimeout:      cfg.MQTT.PingTimeout,
		MaxReconnectWait: cfg.MQTT.MaxReconnectWait,
		QoS:              byte(cfg.MQTT.QoS),
	}

	svcConfig := &gateway.ServiceConfig{
		MQTT:            mqttConfig,
		UplinkBridge:    bridge.DefaultUplinkBridgeConfig(),
		DownlinkBridge:  bridge.DefaultDownlinkBridgeConfig(),
		CleanupInterval: 5 * time.Minute,
		MaxStaleAge:     24 * time.Hour,
		SubscribeTopics: []string{
			"thing/product/+/+/#",
			"sys/product/+/#",
		},
	}

	// Create gateway service
	logEntry := log.WithService(serviceName)
	gatewaySvc := gateway.NewService(svcConfig, publisher, subscriber, metricsCollector, logEntry)

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
		// Check RabbitMQ and MQTT connections
		rmqReady := rmqClient.IsConnected()
		mqttReady := gatewaySvc.IsMQTTConnected()

		if rmqReady && mqttReady {
			c.JSON(http.StatusOK, gin.H{
				"status":         "ready",
				"rabbitmq":       "connected",
				"mqtt":           "connected",
				"online_devices": gatewaySvc.GetOnlineDeviceCount(),
			})
			return
		}

		status := gin.H{
			"status":   "not ready",
			"rabbitmq": "disconnected",
			"mqtt":     "disconnected",
		}
		if rmqReady {
			status["rabbitmq"] = "connected"
		}
		if mqttReady {
			status["mqtt"] = "connected"
		}
		c.JSON(http.StatusServiceUnavailable, status)
	})

	// Metrics endpoint
	router.GET(cfg.Metrics.Path, metrics.Handler(metricsCollector))

	// Device status endpoint
	router.GET("/devices/online", func(c *gin.Context) {
		connManager := gatewaySvc.GetConnectionManager()
		devices := connManager.GetOnlineDevices()

		deviceList := make([]gin.H, 0, len(devices))
		for _, d := range devices {
			deviceList = append(deviceList, gin.H{
				"device_sn":    d.DeviceSN,
				"client_id":    d.ClientID,
				"ip_address":   d.IPAddress,
				"connected_at": d.ConnectedAt,
				"last_seen_at": d.LastSeenAt,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"count":   len(devices),
			"devices": deviceList,
		})
	})

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
	shutdown.Register(func(ctx context.Context) error {
		log.WithService(serviceName).Info("Stopping gateway service")
		return gatewaySvc.Stop()
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

	// Start gateway service (connects to MQTT broker)
	ctx := context.Background()
	if err := gatewaySvc.Start(ctx); err != nil {
		log.WithService(serviceName).Warnf("failed to start gateway service: %v", err)
		log.WithService(serviceName).Info("Gateway running in degraded mode - MQTT not connected")
	} else {
		log.WithService(serviceName).Info("Gateway service started - MQTT ↔ RabbitMQ bridge active")
	}

	// Wait for shutdown signal
	if err := shutdown.Wait(); err != nil {
		log.WithService(serviceName).Errorf("shutdown error: %v", err)
	}
	log.WithService(serviceName).Info("Service stopped")
}
