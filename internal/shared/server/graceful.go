// Package server provides server utilities including graceful shutdown.
package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// WaitForShutdown waits for shutdown signals and returns a context that is canceled on shutdown.
func WaitForShutdown(_ time.Duration) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	return ctx
}

// GracefulShutdown handles graceful shutdown with a timeout.
type GracefulShutdown struct {
	cleanups []func(context.Context) error
	timeout  time.Duration
}

// NewGracefulShutdown creates a new GracefulShutdown handler.
func NewGracefulShutdown(timeout time.Duration) *GracefulShutdown {
	return &GracefulShutdown{
		timeout:  timeout,
		cleanups: make([]func(context.Context) error, 0),
	}
}

// Register registers a cleanup function to be called during shutdown.
// Cleanup functions are called in reverse order of registration.
func (g *GracefulShutdown) Register(cleanup func(context.Context) error) {
	g.cleanups = append(g.cleanups, cleanup)
}

// Wait waits for shutdown signal and executes cleanup functions.
func (g *GracefulShutdown) Wait() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	return g.Shutdown()
}

// Shutdown executes all cleanup functions with timeout.
func (g *GracefulShutdown) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()

	var lastErr error
	// Execute cleanups in reverse order
	for i := len(g.cleanups) - 1; i >= 0; i-- {
		if err := g.cleanups[i](ctx); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// ShutdownContext returns a context that is canceled when a shutdown signal is received.
func ShutdownContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	return ctx, cancel
}
