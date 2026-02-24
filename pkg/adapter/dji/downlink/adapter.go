// Package downlink provides DJI downlink message dispatching functionality
package downlink

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// DispatcherAdapter adapts the DJI Dispatcher to the public adapter.DownlinkDispatcher interface.
// This adapter does NOT depend on internal/ packages.
type DispatcherAdapter struct {
	dispatcher *Dispatcher
}

// NewDispatcherAdapter creates a new dispatcher adapter
func NewDispatcherAdapter(publisher *rabbitmq.Publisher, logger *logrus.Entry) *DispatcherAdapter {
	return &DispatcherAdapter{
		dispatcher: NewDispatcher(publisher, logger),
	}
}

// GetVendor returns the vendor name
func (a *DispatcherAdapter) GetVendor() string {
	return a.dispatcher.GetVendor()
}

// CanDispatch checks if this dispatcher can handle the given vendor
func (a *DispatcherAdapter) CanDispatch(vendor string) bool {
	return vendor == a.dispatcher.vendor
}

// Dispatch dispatches a service call and returns the result in the public adapter format
func (a *DispatcherAdapter) Dispatch(ctx context.Context, call *adapter.ServiceCall) (*adapter.DispatchResult, error) {
	// Convert to DJI format
	djiCall := &ServiceCall{
		ID:       call.ID,
		DeviceSN: call.DeviceSN,
		Vendor:   call.Vendor,
		Method:   call.Method,
		Params:   call.Params,
	}

	result, err := a.dispatcher.Dispatch(ctx, djiCall)
	if err != nil {
		return &adapter.DispatchResult{
			CallID:   call.ID,
			DeviceSN: call.DeviceSN,
			Vendor:   call.Vendor,
			Method:   call.Method,
			Success:  false,
			Error:    err.Error(),
		}, err
	}

	return &adapter.DispatchResult{
		CallID:     result.CallID,
		DeviceSN:   result.DeviceSN,
		Vendor:     result.Vendor,
		Method:     result.Method,
		Success:    result.Success,
		Error:      result.ErrorMsg,
		RoutingKey: result.RoutingKey,
	}, nil
}

// Ensure DispatcherAdapter implements adapter.DownlinkDispatcher
var _ adapter.DownlinkDispatcher = (*DispatcherAdapter)(nil)
