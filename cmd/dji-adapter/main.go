// Package main provides the DJI protocol adapter service.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/utmos/utmos/internal/shared/config"
	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/adapter/dji"
	djiinit "github.com/utmos/utmos/pkg/adapter/dji/init"
	pkgconfig "github.com/utmos/utmos/pkg/config"
	"github.com/utmos/utmos/pkg/logger"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

const serviceName = "dji-adapter"

// Prometheus metrics
var (
	messagesProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dji_adapter_messages_processed_total",
			Help: "Total number of messages processed by the DJI adapter",
		},
		[]string{"direction", "topic_type", "status"},
	)

	messageProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dji_adapter_message_processing_duration_seconds",
			Help:    "Duration of message processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"direction", "topic_type"},
	)

	parseErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dji_adapter_parse_errors_total",
			Help: "Total number of parse errors",
		},
		[]string{"error_type"},
	)
)

func main() {
	// Load configuration
	cfg, err := config.Load(serviceName)
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	log := logger.New(&cfg.Logger)
	log.WithField("service", serviceName).Info("Starting DJI adapter service")

	// Register DJI adapter
	dji.Register()
	log.Info("DJI adapter registered")

	// Get adapter instance
	protoAdapter, err := adapter.Get(dji.VendorDJI)
	if err != nil {
		log.WithError(err).Fatal("Failed to get DJI adapter")
	}

	// Type assert to get the concrete *dji.Adapter for full initialization
	djiAdapter, ok := protoAdapter.(*dji.Adapter)
	if !ok {
		log.Fatal("Failed to get concrete DJI adapter")
	}

	// Initialize full adapter with handler registry (ServiceRouter, EventRouter, all handlers)
	// This enables full protocol handling: OSD, State, Status, Service, Event, Request, DRC
	if err := djiinit.InitializeAdapter(djiAdapter); err != nil {
		log.WithError(err).Fatal("Failed to initialize DJI adapter with handlers")
	}
	log.Info("DJI adapter initialized with full handler registry")

	// Initialize RabbitMQ client
	rmqClient := rabbitmq.NewClient(&cfg.RabbitMQ)

	// Connect to RabbitMQ
	ctx := context.Background()
	if err := rmqClient.Connect(ctx); err != nil {
		log.WithError(err).Fatal("Failed to connect to RabbitMQ")
	}
	defer func() {
		if err := rmqClient.Close(); err != nil {
			log.WithError(err).Error("Failed to close RabbitMQ client")
		}
	}()

	// Declare exchange
	if err := rmqClient.DeclareExchange(cfg.RabbitMQ.ExchangeName, cfg.RabbitMQ.ExchangeType); err != nil {
		log.WithError(err).Fatal("Failed to declare exchange")
	}

	// Create RabbitMQ publisher for trace-injected publishing
	publisher := rabbitmq.NewPublisher(rmqClient)

	// Create context for graceful shutdown
	shutdownCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start message processing
	go processUplinkMessages(shutdownCtx, log, rmqClient, publisher, djiAdapter, &cfg.RabbitMQ)
	go processDownlinkMessages(shutdownCtx, log, rmqClient, publisher, djiAdapter, &cfg.RabbitMQ)

	// Setup HTTP server for health check and metrics
	router := setupRouter(log, rmqClient)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Start HTTP server
	go func() {
		log.Info("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Fatal("HTTP server failed")
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down DJI adapter service")

	// Cancel context to stop message processing
	cancel()

	// Shutdown HTTP server
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer timeoutCancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.WithError(err).Error("HTTP server shutdown failed")
	}

	log.Info("DJI adapter service stopped")
}

func setupRouter(_ *logger.Logger, rmqClient *rabbitmq.Client) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		status := "healthy"
		rmqStatus := "connected"

		if !rmqClient.IsConnected() {
			status = "unhealthy"
			rmqStatus = "disconnected"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   status,
			"service":  serviceName,
			"rabbitmq": rmqStatus,
		})
	})

	// Readiness probe
	router.GET("/ready", func(c *gin.Context) {
		if rmqClient.IsConnected() {
			c.JSON(http.StatusOK, gin.H{"status": "ready"})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
		}
	})

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router
}

func processUplinkMessages(ctx context.Context, log *logger.Logger, rmqClient *rabbitmq.Client, publisher *rabbitmq.Publisher, djiAdapter adapter.ProtocolAdapter, rmqCfg *pkgconfig.RabbitMQConfig) {
	log.Info("Starting uplink message processor")

	// Declare and bind queue for raw DJI uplink messages
	rawRoutingKey := rabbitmq.NewRawRoutingKey(dji.VendorDJI, rabbitmq.DirectionUplink)
	queueName := "dji-adapter-uplink"

	if _, err := rmqClient.DeclareQueue(queueName, true); err != nil {
		log.WithError(err).Error("Failed to declare uplink queue")
		return
	}

	if err := rmqClient.BindQueue(queueName, rawRoutingKey.String(), rmqCfg.ExchangeName); err != nil {
		log.WithError(err).Error("Failed to bind uplink queue")
		return
	}

	// Get channel and consume messages
	channel := rmqClient.Channel()
	if channel == nil {
		log.Error("Channel is nil")
		return
	}

	msgs, err := channel.Consume(
		queueName,
		"",    // consumer tag
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.WithError(err).Error("Failed to start consuming uplink messages")
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping uplink message processor")
			return
		case msg, ok := <-msgs:
			if !ok {
				log.Warn("Uplink message channel closed")
				return
			}
			processUplinkMessage(ctx, log, rmqClient, publisher, djiAdapter, rmqCfg, msg)
		}
	}
}

func processUplinkMessage(ctx context.Context, log *logger.Logger, rmqClient *rabbitmq.Client, publisher *rabbitmq.Publisher, djiAdapter adapter.ProtocolAdapter, rmqCfg *pkgconfig.RabbitMQConfig, msg amqp.Delivery) {
	start := time.Now()

	// Extract topic from message headers
	var topic string
	if t, ok := msg.Headers["original_topic"].(string); ok {
		topic = t
	}
	if topic == "" {
		parseErrors.WithLabelValues("missing_topic").Inc()
		log.Warn("Message missing original_topic header")
		_ = msg.Nack(false, false)
		return
	}

	// Parse raw message
	pm, err := djiAdapter.ParseRawMessage(topic, msg.Body)
	if err != nil {
		parseErrors.WithLabelValues("parse_error").Inc()
		log.WithError(err).WithField("topic", topic).Error("Failed to parse raw message")
		messagesProcessed.WithLabelValues("uplink", "unknown", "error").Inc()
		_ = msg.Nack(false, false)
		return
	}

	// Convert to standard message
	stdMsg, err := djiAdapter.ToStandardMessage(pm)
	if err != nil {
		parseErrors.WithLabelValues("conversion_error").Inc()
		log.WithError(err).Error("Failed to convert to standard message")
		messagesProcessed.WithLabelValues("uplink", string(pm.MessageType), "error").Inc()
		_ = msg.Nack(false, false)
		return
	}

	// Build routing key for standard message
	routingKey := rabbitmq.NewRoutingKey(dji.VendorDJI, "device", stdMsg.Action)

	// Publish to exchange using publisher which auto-injects W3C trace context
	err = publisher.Publish(ctx, routingKey.String(), stdMsg)
	if err != nil {
		log.WithError(err).Error("Failed to publish standard message")
		messagesProcessed.WithLabelValues("uplink", string(pm.MessageType), "error").Inc()
		_ = msg.Nack(false, true)
		return
	}

	duration := time.Since(start).Seconds()
	messageProcessingDuration.WithLabelValues("uplink", string(pm.MessageType)).Observe(duration)
	messagesProcessed.WithLabelValues("uplink", string(pm.MessageType), "success").Inc()

	_ = msg.Ack(false)

	log.WithFields(map[string]any{
		"tid":      stdMsg.TID,
		"device":   stdMsg.DeviceSN,
		"action":   stdMsg.Action,
		"duration": duration,
	}).Debug("Processed uplink message")
}

func processDownlinkMessages(ctx context.Context, log *logger.Logger, rmqClient *rabbitmq.Client, publisher *rabbitmq.Publisher, djiAdapter adapter.ProtocolAdapter, rmqCfg *pkgconfig.RabbitMQConfig) {
	log.Info("Starting downlink message processor")

	// Declare and bind queue for standard messages for DJI devices (service calls)
	// Uses canonical iot.{vendor}.{service}.{action} format with service=downlink
	bindingPattern := rabbitmq.BuildBindingPattern(dji.VendorDJI, "downlink", "")
	queueName := "dji-adapter-downlink"

	if _, err := rmqClient.DeclareQueue(queueName, true); err != nil {
		log.WithError(err).Error("Failed to declare downlink queue")
		return
	}

	if err := rmqClient.BindQueue(queueName, bindingPattern, rmqCfg.ExchangeName); err != nil {
		log.WithError(err).Error("Failed to bind downlink queue")
		return
	}

	// Get channel and consume messages
	channel := rmqClient.Channel()
	if channel == nil {
		log.Error("Channel is nil")
		return
	}

	msgs, err := channel.Consume(
		queueName,
		"",    // consumer tag
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.WithError(err).Error("Failed to start consuming downlink messages")
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping downlink message processor")
			return
		case msg, ok := <-msgs:
			if !ok {
				log.Warn("Downlink message channel closed")
				return
			}
			processDownlinkMessage(ctx, log, rmqClient, publisher, djiAdapter, rmqCfg, msg)
		}
	}
}

func processDownlinkMessage(ctx context.Context, log *logger.Logger, rmqClient *rabbitmq.Client, publisher *rabbitmq.Publisher, djiAdapter adapter.ProtocolAdapter, rmqCfg *pkgconfig.RabbitMQConfig, msg amqp.Delivery) {
	start := time.Now()

	// Parse standard message
	var stdMsg rabbitmq.StandardMessage
	if err := json.Unmarshal(msg.Body, &stdMsg); err != nil {
		parseErrors.WithLabelValues("unmarshal_error").Inc()
		log.WithError(err).Error("Failed to unmarshal standard message")
		messagesProcessed.WithLabelValues("downlink", "unknown", "error").Inc()
		_ = msg.Nack(false, false)
		return
	}

	// Convert to protocol message
	pm, err := djiAdapter.FromStandardMessage(&stdMsg)
	if err != nil {
		parseErrors.WithLabelValues("conversion_error").Inc()
		log.WithError(err).Error("Failed to convert from standard message")
		messagesProcessed.WithLabelValues("downlink", "unknown", "error").Inc()
		_ = msg.Nack(false, false)
		return
	}

	// Get raw payload for MQTT publish
	payload, err := djiAdapter.GetRawPayload(pm)
	if err != nil {
		log.WithError(err).Error("Failed to get raw payload")
		messagesProcessed.WithLabelValues("downlink", string(pm.MessageType), "error").Inc()
		_ = msg.Nack(false, false)
		return
	}

	// Build standard message for gateway bridge
	// Gateway's DownlinkBridge.handleStandardMessage expects StandardMessage with
	// Data containing {topic, payload, qos} fields
	dataPayload := map[string]any{
		"topic":   pm.Topic,
		"payload": json.RawMessage(payload),
		"qos":     1,
	}

	// Create new StandardMessage preserving tid/bid from original, with corrected data format
	// for gateway bridge which expects {topic, payload, qos} in Data
	downlinkMsg, err := rabbitmq.NewStandardMessageWithIDs(
		stdMsg.TID,
		stdMsg.BID,
		"dji-adapter",
		"downlink",
		pm.DeviceSN,
		dataPayload,
	)
	if err != nil {
		log.WithError(err).Error("Failed to create standard message")
		messagesProcessed.WithLabelValues("downlink", string(pm.MessageType), "error").Inc()
		_ = msg.Nack(false, false)
		return
	}

	// Build raw routing key for downlink: iot.raw.dji.downlink
	rawRoutingKey := rabbitmq.NewRawRoutingKey(dji.VendorDJI, rabbitmq.DirectionDownlink)

	// Publish to exchange using rabbitmq publisher for trace injection
	channel := rmqClient.Channel()
	if channel == nil {
		log.Error("Channel is nil")
		_ = msg.Nack(false, true)
		return
	}

	// Publish to exchange using publisher which auto-injects W3C trace context
	err = publisher.Publish(ctx, rawRoutingKey.String(), downlinkMsg)
	if err != nil {
		log.WithError(err).Error("Failed to publish raw message")
		messagesProcessed.WithLabelValues("downlink", string(pm.MessageType), "error").Inc()
		_ = msg.Nack(false, true)
		return
	}

	duration := time.Since(start).Seconds()
	messageProcessingDuration.WithLabelValues("downlink", string(pm.MessageType)).Observe(duration)
	messagesProcessed.WithLabelValues("downlink", string(pm.MessageType), "success").Inc()

	_ = msg.Ack(false)

	log.WithFields(map[string]any{
		"tid":      downlinkMsg.TID,
		"device":   downlinkMsg.DeviceSN,
		"action":   downlinkMsg.Action,
		"duration": duration,
	}).Debug("Processed downlink message")
}
