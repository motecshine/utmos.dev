// Package init provides initialization for the DJI adapter.
package init

import (
	"fmt"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/adapter/dji/handler"
	"github.com/utmos/utmos/pkg/adapter/dji/router"
)

// registryAdapter wraps handler.Registry to implement dji.HandlerRegistry.
type registryAdapter struct {
	registry *handler.Registry
}

// Get returns the handler for a topic type.
func (r *registryAdapter) Get(topicType dji.TopicType) (dji.MessageHandler, error) {
	return r.registry.Get(topicType)
}

// registerServiceCommands registers all service commands to the router.
func registerServiceCommands(sr *router.ServiceRouter) error {
	registrations := []struct {
		name string
		fn   func(*router.ServiceRouter) error
	}{
		{"device commands", router.RegisterDeviceCommands},
		{"camera commands", router.RegisterCameraCommands},
		{"wayline commands", router.RegisterWaylineCommands},
		{"drc commands", router.RegisterDRCCommands},
		{"file commands", router.RegisterFileCommands},
		{"firmware commands", router.RegisterFirmwareCommands},
		{"live commands", router.RegisterLiveCommands},
		{"config commands", router.RegisterConfigCommands},
	}

	for _, r := range registrations {
		if err := r.fn(sr); err != nil {
			return fmt.Errorf("register %s: %w", r.name, err)
		}
	}
	return nil
}

// registerEvents registers all events to the router.
func registerEvents(er *router.EventRouter) error {
	registrations := []struct {
		name string
		fn   func(*router.EventRouter) error
	}{
		{"core events", router.RegisterCoreEvents},
		{"wayline events", router.RegisterWaylineEvents},
		{"drc events", router.RegisterDRCEvents},
		{"file events", router.RegisterFileEvents},
		{"firmware events", router.RegisterFirmwareEvents},
	}

	for _, r := range registrations {
		if err := r.fn(er); err != nil {
			return fmt.Errorf("register %s: %w", r.name, err)
		}
	}
	return nil
}

// registerHandlers registers all handlers to the registry.
func registerHandlers(registry *handler.Registry, sr *router.ServiceRouter, er *router.EventRouter) error {
	handlers := []handler.Handler{
		handler.NewOSDHandler(),
		handler.NewStateHandler(),
		handler.NewStatusHandler(),
		handler.NewServiceHandler(sr),
		handler.NewEventHandler(er),
		handler.NewRequestHandler(sr),
		handler.NewDRCHandler(sr, er),
	}

	for _, h := range handlers {
		if err := registry.Register(h); err != nil {
			return fmt.Errorf("register %T: %w", h, err)
		}
	}
	return nil
}

// InitializeAdapter initializes the DJI adapter with all handlers.
// Returns an error if any registration fails.
func InitializeAdapter(a *dji.Adapter) error {
	registry := handler.NewRegistry()
	serviceRouter := router.NewServiceRouter()
	eventRouter := router.NewEventRouter()

	if err := registerServiceCommands(serviceRouter); err != nil {
		return err
	}
	if err := registerEvents(eventRouter); err != nil {
		return err
	}
	if err := registerHandlers(registry, serviceRouter, eventRouter); err != nil {
		return err
	}

	a.SetHandlerRegistry(&registryAdapter{registry: registry})
	return nil
}

// NewInitializedAdapter creates a new DJI adapter with all handlers initialized.
// This function panics if initialization fails and should only be used during
// program initialization (e.g., in init() or main()).
// For error handling, use dji.NewAdapter() followed by InitializeAdapter().
func NewInitializedAdapter() *dji.Adapter {
	a := dji.NewAdapter()
	if err := InitializeAdapter(a); err != nil {
		panic(fmt.Errorf("NewInitializedAdapter: %w", err))
	}
	return a
}

// MustInitializeAdapter initializes the adapter and panics on error.
// This should only be used during program initialization (e.g., in init() or main()).
// For runtime initialization with error handling, use InitializeAdapter().
func MustInitializeAdapter(a *dji.Adapter) {
	if err := InitializeAdapter(a); err != nil {
		panic(fmt.Errorf("MustInitializeAdapter: %w", err))
	}
}
