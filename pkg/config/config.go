// Package config provides public configuration types for UMOS IoT services.
// This package contains configuration structures that can be used by external packages.
package config

import "time"

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

// LoggerConfig holds logger configuration.
type LoggerConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	Output   string `yaml:"output"`
	FilePath string `yaml:"file_path,omitempty"`
}
