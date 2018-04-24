package transport

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Transport is an interface for communicating messages between endpoints.
type Transport interface {
	Initialize(ctx context.Context, ep string) error
	Subscribe(ctx context.Context, mt ax.MessageTypeSet) error
	Send(ctx context.Context, m OutboundMessage) error
	Receive(ctx context.Context) (InboundMessage, error)
}
