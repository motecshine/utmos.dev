// Package observability provides metrics, logging, and tracing for the DJI adapter.
package observability

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/utmos/utmos/pkg/logger"
)

// Logger wraps the shared logger with DJI-specific fields.
type Logger struct {
	log *logger.Logger
}

// NewLogger creates a new DJI adapter logger.
func NewLogger(log *logger.Logger) *Logger {
	return &Logger{log: log}
}

// DefaultLogger creates a default logger.
func DefaultLogger() *Logger {
	return &Logger{log: logger.Default()}
}

// WithContext returns a log entry with trace context.
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	return l.log.WithTrace(ctx).WithField("vendor", "dji")
}

// WithMessage returns a log entry with message context.
func (l *Logger) WithMessage(ctx context.Context, deviceSN, messageType, method string) *logrus.Entry {
	return l.log.WithTrace(ctx).WithFields(logrus.Fields{
		"vendor":       "dji",
		"device_sn":    deviceSN,
		"message_type": messageType,
		"method":       method,
	})
}

// withDJIFields returns a log entry with vendor "dji" and the given extra fields.
func (l *Logger) withDJIFields(ctx context.Context, fields logrus.Fields) *logrus.Entry {
	fields["vendor"] = "dji"
	return l.log.WithTrace(ctx).WithFields(fields)
}

// WithDevice returns a log entry with device context.
func (l *Logger) WithDevice(ctx context.Context, deviceSN, gatewaySN string) *logrus.Entry {
	return l.withDJIFields(ctx, logrus.Fields{"device_sn": deviceSN, "gateway_sn": gatewaySN})
}

// WithTID returns a log entry with transaction ID.
func (l *Logger) WithTID(ctx context.Context, tid, bid string) *logrus.Entry {
	return l.withDJIFields(ctx, logrus.Fields{"tid": tid, "bid": bid})
}

// Info logs an info message with context.
func (l *Logger) Info(ctx context.Context, msg string) {
	l.WithContext(ctx).Info(msg)
}

// Error logs an error message with context.
func (l *Logger) Error(ctx context.Context, msg string, err error) {
	l.WithContext(ctx).WithError(err).Error(msg)
}

// Debug logs a debug message with context.
func (l *Logger) Debug(ctx context.Context, msg string) {
	l.WithContext(ctx).Debug(msg)
}

// Warn logs a warning message with context.
func (l *Logger) Warn(ctx context.Context, msg string) {
	l.WithContext(ctx).Warn(msg)
}
