package pipeline

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/transport"
	"go.uber.org/multierr"
)

// Router is an InboundStage that passes messages to the message handlers
// that are subscribed to the particular message type.
type Router struct {
	Routes RoutingTable
}

// Initialize configures the transport to subscribe to all messages in the
// routing table.
func (r *Router) Initialize(ctx context.Context, t transport.Transport) error {
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
	o OutboundStage,
	m transport.InboundMessage,
) error {
	mc := &MessageContext{
		Context: ctx,
		In:      m,
	}

	mt := ax.TypeOf(m.Envelope.Message)

	for _, h := range r.Routes[mt] {
		if err := h.HandleMessage(mc, m.Envelope.Message); err != nil {
			return err
		}
	}

	return nil
}

// RoutingTable is a map of message type to the handlers that receive that
// message type.
type RoutingTable map[ax.MessageType][]ax.MessageHandler

// NewRoutingTable returns a new routing table for the given set of handlers.
//
// It returns an error if two message handlers have registered interest in the
// same command.
func NewRoutingTable(handlers []ax.MessageHandler) (RoutingTable, error) {
	rt := RoutingTable{}

	for _, h := range handlers {
		for _, mt := range h.MessageTypes().Members() {
			rt[mt] = append(rt[mt], h)
		}
	}

	var err error

	for mt, h := range rt {
		if mt.IsCommand() && len(h) > 1 {
			err = multierr.Append(
				err,
				DuplicateRoutesError{
					MessageType: mt,
					Handlers:    h,
				},
			)
		}
	}

	return rt, err
}

// DuplicateRoutesError is returned by NewRoutingTable if multiple handlers have
// tried to subscribe to the same command type.
type DuplicateRoutesError struct {
	MessageType ax.MessageType
	Handlers    []ax.MessageHandler
}

func (e DuplicateRoutesError) Error() string {
	return fmt.Sprintf(
		"can not build routing table, multiple message handlers are defined for the %s command",
		e.MessageType.Name,
	)
}
