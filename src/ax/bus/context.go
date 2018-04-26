package bus

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// MessageContext is the standard implementation of ax.MessageContext.
type MessageContext struct {
	context.Context

	Envelope   ax.Envelope
	Dispatcher OutboundDispatcher
}

// MessageEnvelope returns the envelope containing the message being handled.
func (c *MessageContext) MessageEnvelope() ax.Envelope {
	return c.Envelope
}

// ExecuteCommand enqueues a command to be executed.
func (c *MessageContext) ExecuteCommand(m ax.Command) error {
	return c.Dispatcher.DispatchMessage(c, OutboundMessage{
		Operation: OpSendUnicast,
		Envelope:  c.Envelope.NewEffect(m),
	})
}

// PublishEvent enqueues events to be published.
func (c *MessageContext) PublishEvent(m ax.Event) error {
	return c.Dispatcher.DispatchMessage(c, OutboundMessage{
		Operation: OpSendMulticast,
		Envelope:  c.Envelope.NewEffect(m),
	})
}
