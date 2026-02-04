package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/utmos/utmos/internal/shared/config"
	"github.com/utmos/utmos/internal/shared/database"
	"github.com/utmos/utmos/internal/shared/logger"
	"github.com/utmos/utmos/internal/shared/server"
	"github.com/utmos/utmos/pkg/metrics"
	"github.com/utmos/utmos/pkg/models"
	"github.com/utmos/utmos/pkg/tracer"
)

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
		// Check database connection
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// Metrics endpoint
	router.GET(cfg.Metrics.Path, metrics.Handler(metricsCollector))

	// API routes (placeholder)
	v1 := router.Group("/api/v1")
	v1.GET("/devices", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "list devices"})
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

	// Wait for shutdown signal
	if err := shutdown.Wait(); err != nil {
		log.WithService(serviceName).Errorf("shutdown error: %v", err)
	}
	log.WithService(serviceName).Info("Service stopped")
}
