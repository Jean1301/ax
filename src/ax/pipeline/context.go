package pipeline

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/transport"
)

// MessageContext is an implementation of ax.MessageHandler that
// appends outbound messages to a slice.
type MessageContext struct {
	context.Context

	In  transport.InboundMessage
	Out []transport.OutboundMessage
}

// MessageEnvelope returns the envelope containing the message being handled.
func (c *MessageContext) MessageEnvelope() ax.Envelope {
	return c.In.Envelope
}

// ExecuteCommand enqueues a command to be executed.
func (c *MessageContext) ExecuteCommand(m ax.Command) {
	c.Out = append(c.Out, transport.OutboundMessage{
		Operation:       transport.OpSendUnicast,
		Envelope:        c.MessageEnvelope().New(m),
		UnicastEndpoint: ax.TypeOf(m).PackageName(),
	})
}

// PublishEvent enqueues events to be published.
func (c *MessageContext) PublishEvent(m ax.Event) {
	c.Out = append(c.Out, transport.OutboundMessage{
		Operation: transport.OpSendMulticast,
		Envelope:  c.MessageEnvelope().New(m),
	})
}
