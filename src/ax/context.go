package ax

import (
	"context"
)

// MessageContext is a specialization of context.Context used by message
// handlers.
//
// It carries information about the messaging behing handled, and allows the
// handler to produce new messages.
type MessageContext interface {
	context.Context

	// MessageEnvelope returns the envelope containing the message being handled.
	MessageEnvelope() Envelope

	// ExecuteCommand enqueues a command to be executed.
	ExecuteCommand(Command) error

	// PublishEvent enqueues events to be published.
	PublishEvent(Event) error
}

// WithContext returns a new MessageContext that forwards messages to mc,
// but uses ctx for context.Context related operations.
func WithContext(ctx context.Context, mc MessageContext) MessageContext {
	return contextProxy{ctx, mc}
}

type contextProxy struct {
	context.Context
	parent MessageContext
}

// MessageEnvelope returns the envelope containing the message being handled.
func (c contextProxy) MessageEnvelope() Envelope {
	return c.parent.MessageEnvelope()
}

// ExecuteCommand enqueues a command to be executed.
func (c contextProxy) ExecuteCommand(m Command) error {
	return c.parent.ExecuteCommand(m)
}

// PublishEvent enqueues events to be published.
func (c contextProxy) PublishEvent(m Event) error {
	return c.parent.PublishEvent(m)
}
