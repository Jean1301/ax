package pipeline

import (
	"context"

	"github.com/jmalloc/ax/src/ax/transport"
)

// InboundStage is a step within an inbound message pipeline.
type InboundStage interface {
	// TODO
	Initialize(context.Context, transport.Transport) error

	// DispatchMessage forwards an inbound message through the pipeline until
	// it ultimately is handled by zero or more message handlers.
	DispatchMessage(
		context.Context,
		OutboundStage,
		transport.InboundMessage,
	) error
}
