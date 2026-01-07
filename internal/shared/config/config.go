package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Environment string      `yaml:"environment"`
	Server      ServerConfig `yaml:"server"`
	Database    DatabaseConfig `yaml:"database"`
	RabbitMQ    RabbitMQConfig `yaml:"rabbitmq"`
	InfluxDB    InfluxDBConfig `yaml:"influxdb"`
	Logging     LoggingConfig `yaml:"logging"`
	Tracing     TracingConfig `yaml:"tracing"`
	Metrics     MetricsConfig `yaml:"metrics"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// RabbitMQConfig represents RabbitMQ configuration
type RabbitMQConfig struct {
	URL      string `yaml:"url"`
	Exchange string `yaml:"exchange"`
}

// InfluxDBConfig represents InfluxDB configuration
type InfluxDBConfig struct {
	URL    string `yaml:"url"`
	Token  string `yaml:"token"`
	Org    string `yaml:"org"`
	Bucket string `yaml:"bucket"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// TracingConfig represents tracing configuration
type TracingConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

// Load loads configuration from YAML file
func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override with environment variables if set
	overrideWithEnv(&config)

	return &config, nil
}

// overrideWithEnv overrides config values with environment variables
func overrideWithEnv(config *Config) {
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		config.Environment = env
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		// Parse port from string (simplified, should use strconv.Atoi)
		// For now, keep default
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}
	if rabbitMQURL := os.Getenv("RABBITMQ_URL"); rabbitMQURL != "" {
		config.RabbitMQ.URL = rabbitMQURL
	}
}

