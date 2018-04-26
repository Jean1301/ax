package persistence

import (
	"context"

	"github.com/jmalloc/ax/src/ax/bus"
)

// Injector is an implementation of bus.InboundPipeline that injects a data
// store into the context.
type Injector struct {
	DataStore DataStore
	Next      bus.InboundPipeline
}

// Initialize is called when the transport is initialized.
func (i *Injector) Initialize(ctx context.Context, t bus.Transport) error {
	return i.Next.Initialize(ctx, t)
}

// DispatchMessage calls i.Next.DispatchMessage() with a context containing
// i.DataStore.
func (i *Injector) DispatchMessage(
	ctx context.Context,
	out bus.OutboundDispatcher,
	m bus.InboundMessage,
) error {
	return i.Next.DispatchMessage(
		WithDataStore(ctx, i.DataStore),
		out,
		m,
	)
}
