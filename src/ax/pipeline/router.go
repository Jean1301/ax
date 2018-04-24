package pipeline

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"go.uber.org/multierr"
)

// Router is an InboundStage that passes messages to the message handlers
// that are subscribed to the particular message type.
type Router struct {
	Routes RoutingTable
}

func (r *Router) DispatchMessage(
	ctx context.Context,
	o OutboundStage,
	m InboundMessage,
) error {
	mc := &messageContext{
		Context: ctx,
		inbound: m,
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
	// // TODO:
	panic("ni")
}
