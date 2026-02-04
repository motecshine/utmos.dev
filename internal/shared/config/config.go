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
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	SSLMode         string        `yaml:"sslmode"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	Port            int           `yaml:"port"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
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
	Endpoint     string        `yaml:"endpoint"`
	ServiceName  string        `yaml:"service_name"`
	BatchTimeout time.Duration `yaml:"batch_timeout"`
	SamplingRate float64       `yaml:"sampling_rate"`
	MaxQueueSize int           `yaml:"max_queue_size"`
	Enabled      bool          `yaml:"enabled"`
}

// MetricsConfig holds Prometheus metrics configuration.
type MetricsConfig struct {
	Path      string `yaml:"path"`
	Namespace string `yaml:"namespace"`
	Port      int    `yaml:"port"`
	Enabled   bool   `yaml:"enabled"`
}

// LoggerConfig holds logger configuration.
type LoggerConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	Output   string `yaml:"output"`
	FilePath string `yaml:"file_path,omitempty"`
}
