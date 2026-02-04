package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Load loads configuration from a YAML file with environment variable substitution.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Substitute environment variables
	content := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults
	applyDefaults(&cfg)

	return &cfg, nil
}

// LoadFromEnv loads configuration from environment variables with a given prefix.
func LoadFromEnv(env string) (*Config, error) {
	configPath := fmt.Sprintf("configs/config.%s.yaml", strings.ToLower(env))
	return Load(configPath)
}

// applyDefaults applies default values to the configuration.
func applyDefaults(cfg *Config) {
	applyServerDefaults(cfg)
	applyDatabaseDefaults(cfg)
	applyRabbitMQDefaults(cfg)
	applyTracerDefaults(cfg)
	applyMetricsDefaults(cfg)
	applyLoggerDefaults(cfg)
}

func applyServerDefaults(cfg *Config) {
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 30 * time.Second
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 30 * time.Second
	}
}

func applyDatabaseDefaults(cfg *Config) {
	// PostgreSQL defaults
	if cfg.Database.Postgres.Host == "" {
		cfg.Database.Postgres.Host = "localhost"
	}
	if cfg.Database.Postgres.Port == 0 {
		cfg.Database.Postgres.Port = 5432
	}
	if cfg.Database.Postgres.SSLMode == "" {
		cfg.Database.Postgres.SSLMode = "disable"
	}
	if cfg.Database.Postgres.MaxIdleConns == 0 {
		cfg.Database.Postgres.MaxIdleConns = 10
	}
	if cfg.Database.Postgres.MaxOpenConns == 0 {
		cfg.Database.Postgres.MaxOpenConns = 100
	}
	if cfg.Database.Postgres.ConnMaxLifetime == 0 {
		cfg.Database.Postgres.ConnMaxLifetime = time.Hour
	}

	// InfluxDB defaults
	if cfg.Database.InfluxDB.URL == "" {
		cfg.Database.InfluxDB.URL = "http://localhost:8086"
	}
}

func applyRabbitMQDefaults(cfg *Config) {
	if cfg.RabbitMQ.URL == "" {
		cfg.RabbitMQ.URL = "amqp://guest:guest@localhost:5672/"
	}
	if cfg.RabbitMQ.ExchangeName == "" {
		cfg.RabbitMQ.ExchangeName = "iot"
	}
	if cfg.RabbitMQ.ExchangeType == "" {
		cfg.RabbitMQ.ExchangeType = "topic"
	}
	if cfg.RabbitMQ.PrefetchCount == 0 {
		cfg.RabbitMQ.PrefetchCount = 10
	}

	// Retry defaults
	if cfg.RabbitMQ.Retry.MaxRetries == 0 {
		cfg.RabbitMQ.Retry.MaxRetries = 10
	}
	if cfg.RabbitMQ.Retry.InitialDelay == 0 {
		cfg.RabbitMQ.Retry.InitialDelay = time.Second
	}
	if cfg.RabbitMQ.Retry.MaxDelay == 0 {
		cfg.RabbitMQ.Retry.MaxDelay = 30 * time.Second
	}
	if cfg.RabbitMQ.Retry.Multiplier == 0 {
		cfg.RabbitMQ.Retry.Multiplier = 2.0
	}
}

func applyTracerDefaults(cfg *Config) {
	if cfg.Tracer.Endpoint == "" {
		cfg.Tracer.Endpoint = "http://localhost:4318/v1/traces"
	}
	if cfg.Tracer.SamplingRate == 0 {
		cfg.Tracer.SamplingRate = 1.0
	}
	if cfg.Tracer.BatchTimeout == 0 {
		cfg.Tracer.BatchTimeout = 5 * time.Second
	}
	if cfg.Tracer.MaxQueueSize == 0 {
		cfg.Tracer.MaxQueueSize = 2048
	}
}

func applyMetricsDefaults(cfg *Config) {
	if cfg.Metrics.Path == "" {
		cfg.Metrics.Path = "/metrics"
	}
	if cfg.Metrics.Namespace == "" {
		cfg.Metrics.Namespace = "iot"
	}
}

func applyLoggerDefaults(cfg *Config) {
	if cfg.Logger.Level == "" {
		cfg.Logger.Level = "info"
	}
	if cfg.Logger.Format == "" {
		cfg.Logger.Format = "json"
	}
	if cfg.Logger.Output == "" {
		cfg.Logger.Output = "stdout"
	}
}
