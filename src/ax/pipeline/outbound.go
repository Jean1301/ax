package pipeline

import (
	"context"

	"github.com/jmalloc/ax/src/ax/transport"
)

// OutboundStage is a step within an outbound message pipeline.
type OutboundStage interface {
	// TODO
	Initialize(context.Context, transport.Transport) error

	// DispatchMessage forwards an outbound message through the pipeline until
	// it ultimately is sent via a messaging transport.
	DispatchMessage(context.Context, transport.OutboundMessage) error
}
