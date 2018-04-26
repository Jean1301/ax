package router

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
)

// Router is an implementation of bus.InboundPipeline that passes messages
// to the message handlers that are subscribed to the particular message type.
type Router struct {
	Routes RoutingTable
}

// Initialize configures the transport to subscribe to all messages in the
// routing table.
func (r *Router) Initialize(ctx context.Context, t bus.Transport) error {
	var set ax.MessageTypeSet

	for mt := range r.Routes {
		if mt.IsEvent() {
			set = set.Add(mt)
		}
	}

	return t.Subscribe(ctx, set)
}

// DispatchMessage passes m to the appropriate message handlers according to the
// routing table.
func (r *Router) DispatchMessage(
	ctx context.Context,
	out bus.OutboundDispatcher,
	m bus.InboundMessage,
) error {
	mc := &bus.MessageContext{
		Context:    ctx,
		Envelope:   m.Envelope,
		Dispatcher: out,
	}

	mt := ax.TypeOf(m.Envelope.Message)

	for _, h := range r.Routes[mt] {
		if err := h.HandleMessage(mc, m.Envelope.Message); err != nil {
			return err
		}
	}

	return nil
}
