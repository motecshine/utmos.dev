// Package logger provides structured logging with trace context support.
package logger

import (
	"context"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"

	"github.com/utmos/utmos/internal/shared/config"
)

// Logger is a wrapper around logrus.Logger with trace context support.
type Logger struct {
	*logrus.Logger
}

// New creates a new Logger instance based on the provided configuration.
func New(cfg *config.LoggerConfig) *Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set log format
	if cfg.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FullTimestamp:   true,
		})
	}

	// Set output
	var output io.Writer
	switch cfg.Output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	case "file":
		if cfg.FilePath != "" {
			file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				output = os.Stdout
			} else {
				output = file
			}
		} else {
			output = os.Stdout
		}
	default:
		output = os.Stdout
	}
	logger.SetOutput(output)

	return &Logger{Logger: logger}
}

// WithTrace extracts trace context from context and returns a logrus Entry with trace fields.
func (l *Logger) WithTrace(ctx context.Context) *logrus.Entry {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		return l.WithFields(logrus.Fields{
			"trace_id": spanCtx.TraceID().String(),
			"span_id":  spanCtx.SpanID().String(),
		})
	}
	return l.WithFields(logrus.Fields{})
}

// WithContext returns a logrus Entry with context fields including trace information.
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	entry := l.WithTrace(ctx)
	return entry
}

// WithService returns a logrus Entry with service name field.
func (l *Logger) WithService(service string) *logrus.Entry {
	return l.WithField("service", service)
}

// WithDevice returns a logrus Entry with device serial number field.
func (l *Logger) WithDevice(deviceSN string) *logrus.Entry {
	return l.WithField("device_sn", deviceSN)
}

// WithVendor returns a logrus Entry with vendor field.
func (l *Logger) WithVendor(vendor string) *logrus.Entry {
	return l.WithField("vendor", vendor)
}

// WithTID returns a logrus Entry with transaction ID field.
func (l *Logger) WithTID(tid string) *logrus.Entry {
	return l.WithField("tid", tid)
}

// WithBID returns a logrus Entry with business ID field.
func (l *Logger) WithBID(bid string) *logrus.Entry {
	return l.WithField("bid", bid)
}

// Default creates a default logger with standard settings.
func Default() *Logger {
	return New(&config.LoggerConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	})
}
