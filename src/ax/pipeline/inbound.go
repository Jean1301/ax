package pipeline

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// InboundStage is a step within an inbound message pipeline.
type InboundStage interface {
	// DispatchMessage forwards an inbound message through the pipeline until
	// it ultimately is handled by zero or more message handlers.
	DispatchMessage(
		context.Context,
		OutboundStage,
		InboundMessage,
	) error
}

// InboundMessage is a container for a message being sent through the inbound
// pipeline.
type InboundMessage struct {
	Envelope ax.Envelope
}
