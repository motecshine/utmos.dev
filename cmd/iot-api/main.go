package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/utmos/utmos/internal/api"
	"github.com/utmos/utmos/internal/api/handler"
	"github.com/utmos/utmos/internal/downlink/dispatcher"
	"github.com/utmos/utmos/internal/shared/config"
	"github.com/utmos/utmos/internal/shared/database"
	"github.com/utmos/utmos/internal/shared/server"
	"github.com/utmos/utmos/pkg/logger"
	djidownlink "github.com/utmos/utmos/pkg/adapter/dji/downlink"
	"github.com/utmos/utmos/pkg/metrics"
	"github.com/utmos/utmos/pkg/models"
	"github.com/utmos/utmos/pkg/rabbitmq"
	"github.com/utmos/utmos/pkg/tracer"
)

// @title UMOS IoT Platform API
// @version 1.0
// @description UMOS IoT Platform RESTful API for device management and service calls
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.utmos.dev/support
// @contact.email support@utmos.dev

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

const serviceName = "iot-api"

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

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database.Postgres)
	if err != nil {
		log.WithService(serviceName).Fatalf("failed to connect to database: %v", err)
	}

	// Run database migrations
	if err := models.AutoMigrate(db); err != nil {
		log.WithService(serviceName).Fatalf("failed to run migrations: %v", err)
	}

	// Initialize RabbitMQ client for service calls
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

	// Initialize RabbitMQ publisher for service calls
	publisher := rabbitmq.NewPublisher(rmqClient)

	// Initialize dispatcher registry and handler
	dispatcherRegistry := dispatcher.NewDispatcherRegistry(log.WithService(serviceName))
	djiDispatcher := djidownlink.NewDispatcherAdapter(publisher, log.WithService(serviceName))
	dispatcherRegistry.Register(dispatcher.NewAdapterDispatcher(djiDispatcher))
	dispatchHandler := dispatcher.NewDispatchHandler(dispatcherRegistry, log.WithService(serviceName))

	// Get API keys from environment
	apiKeys := getAPIKeys()

	// Create router configuration
	routerConfig := &api.Config{
		APIKeys:     apiKeys,
		EnableAuth:  len(apiKeys) > 0,
		EnableTrace: true,
		ServiceName: serviceName,
		TelemetryConfig: &handler.TelemetryConfig{
			URL:    cfg.Database.InfluxDB.URL,
			Token:  cfg.Database.InfluxDB.Token,
			Org:    cfg.Database.InfluxDB.Org,
			Bucket: cfg.Database.InfluxDB.Bucket,
		},
	}

	// Create API router
	apiRouter := api.NewRouter(
		routerConfig,
		db,
		dispatchHandler,
		metricsCollector,
		log.WithService(serviceName),
	)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      apiRouter,
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
		log.WithService(serviceName).Info("Closing API router")
		apiRouter.Close()
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
	shutdown.Register(func(_ context.Context) error {
		log.WithService(serviceName).Info("Closing database connection")
		return database.Close(db)
	})

	// Start server
	go func() {
		log.WithService(serviceName).Infof("Server listening on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithService(serviceName).Fatalf("failed to start server: %v", err)
		}
	}()

	// Log startup info
	log.WithService(serviceName).Infof("API service ready (auth: %v, trace: %v)", routerConfig.EnableAuth, routerConfig.EnableTrace)

	// Wait for shutdown signal
	if err := shutdown.Wait(); err != nil {
		log.WithService(serviceName).Errorf("shutdown error: %v", err)
	}
	log.WithService(serviceName).Info("Service stopped")
}

// getAPIKeys returns API keys from environment
func getAPIKeys() []string {
	keysStr := os.Getenv("API_KEYS")
	if keysStr == "" {
		return nil
	}

	keys := strings.Split(keysStr, ",")
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key != "" {
			result = append(result, key)
		}
	}
	return result
}
