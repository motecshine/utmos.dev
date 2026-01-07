package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// defaultLogger is the default logger instance
	defaultLogger *logrus.Logger
)

// Init initializes the logger with JSON formatter
func Init(level string) error {
	defaultLogger = logrus.New()
	defaultLogger.SetOutput(os.Stdout)
	defaultLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	defaultLogger.SetLevel(logLevel)

	return nil
}

// GetLogger returns the default logger instance
func GetLogger() *logrus.Logger {
	if defaultLogger == nil {
		Init("info")
	}
	return defaultLogger
}

// WithTraceID creates a logger entry with trace_id field
func WithTraceID(traceID string) *logrus.Entry {
	return GetLogger().WithField("trace_id", traceID)
}

// WithSpanID creates a logger entry with span_id field
func WithSpanID(spanID string) *logrus.Entry {
	return GetLogger().WithField("span_id", spanID)
}

// WithTraceContext creates a logger entry with trace_id and span_id fields
func WithTraceContext(traceID, spanID string) *logrus.Entry {
	return GetLogger().WithFields(logrus.Fields{
		"trace_id": traceID,
		"span_id":  spanID,
	})
}

