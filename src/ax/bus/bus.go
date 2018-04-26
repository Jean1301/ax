package bus

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// MessageBus is used to "causal" messages.
type MessageBus struct {
	Dispatcher OutboundDispatcher
}

// ExecuteCommand enqueues a command to be executed.
func (b *MessageBus) ExecuteCommand(ctx context.Context, m ax.Command) error {
	return b.Dispatcher.DispatchMessage(ctx, OutboundMessage{
		Operation: OpSendUnicast,
		Envelope:  ax.NewCause(m),
	})
}

// PublishEvent enqueues events to be published.
func (b *MessageBus) PublishEvent(ctx context.Context, m ax.Event) error {
	return b.Dispatcher.DispatchMessage(ctx, OutboundMessage{
		Operation: OpSendMulticast,
		Envelope:  ax.NewCause(m),
	})
}
