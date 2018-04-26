package bus

import (
	"context"
)

// MessageCollector is an OutboundDispatcher that keeps a collection of the
// dispatched messages in memory.
type MessageCollector struct {
	Messages []OutboundMessage
}

// DispatchMessage adds m to c.Messages.
func (c *MessageCollector) DispatchMessage(ctx context.Context, m OutboundMessage) error {
	c.Messages = append(c.Messages, m)
	return nil
}
