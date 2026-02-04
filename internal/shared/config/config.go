// Package config provides configuration management for UMOS IoT services.
package config

import "time"

// Config is the main application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
	Tracer   TracerConfig   `yaml:"tracer"`
	Metrics  MetricsConfig  `yaml:"metrics"`
	Logger   LoggerConfig   `yaml:"logger"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// DatabaseConfig holds database configuration.
type DatabaseConfig struct {
	Postgres PostgresConfig `yaml:"postgres"`
	InfluxDB InfluxDBConfig `yaml:"influxdb"`
}

// PostgresConfig holds PostgreSQL configuration.
type PostgresConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	SSLMode         string        `yaml:"sslmode"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// InfluxDBConfig holds InfluxDB configuration.
type InfluxDBConfig struct {
	URL    string `yaml:"url"`
	Token  string `yaml:"token"`
	Org    string `yaml:"org"`
	Bucket string `yaml:"bucket"`
}

// RabbitMQConfig holds RabbitMQ configuration.
type RabbitMQConfig struct {
	URL           string      `yaml:"url"`
	ExchangeName  string      `yaml:"exchange_name"`
	ExchangeType  string      `yaml:"exchange_type"`
	PrefetchCount int         `yaml:"prefetch_count"`
	Retry         RetryConfig `yaml:"retry"`
}

// RetryConfig holds retry configuration.
type RetryConfig struct {
	MaxRetries   int           `yaml:"max_retries"`
	InitialDelay time.Duration `yaml:"initial_delay"`
	MaxDelay     time.Duration `yaml:"max_delay"`
	Multiplier   float64       `yaml:"multiplier"`
}

// TracerConfig holds distributed tracing configuration.
type TracerConfig struct {
	Enabled      bool          `yaml:"enabled"`
	Endpoint     string        `yaml:"endpoint"`
	ServiceName  string        `yaml:"service_name"`
	SamplingRate float64       `yaml:"sampling_rate"`
	BatchTimeout time.Duration `yaml:"batch_timeout"`
	MaxQueueSize int           `yaml:"max_queue_size"`
}

// MetricsConfig holds Prometheus metrics configuration.
type MetricsConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Path      string `yaml:"path"`
	Port      int    `yaml:"port"`
	Namespace string `yaml:"namespace"`
}

// LoggerConfig holds logger configuration.
type LoggerConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	Output   string `yaml:"output"`
	FilePath string `yaml:"file_path,omitempty"`
}
