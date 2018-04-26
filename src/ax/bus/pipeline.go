package bus

import (
	"context"
)

// InboundPipeline is an interface for a message pipeline that processes
// messages received from the message transport.
//
// A "stage" within the pipeline is simply an implementation of the
// InboundPipeline interface that forwards messages to another pipeline.
type InboundPipeline interface {
	// Initialize is called when the transport is initialized.
	Initialize(context.Context, Transport) error

	// DispatchMessage forwards an inbound message through the pipeline until
	// it is handled by some application-defined message handler(s).
	DispatchMessage(context.Context, OutboundDispatcher, InboundMessage) error
}

// OutboundPipeline is an interface for a message pipeline that processes
// messages that are sent via the message transport.
//
// A "stage" within the pipeline is simply an implementation of the
// OutboundPipeline interface that forwards messages to another pipeline.
type OutboundPipeline interface {
	OutboundDispatcher

	// Initialize is called when the transport is initialized.
	Initialize(context.Context, Transport) error
}

// OutboundDispatcher is an interface for dispatching an outbound message.
type OutboundDispatcher interface {
	// DispatchMessage forwards an outbound message through the pipeline until
	// it ultimately is sent via a messaging transport.
	DispatchMessage(context.Context, OutboundMessage) error
}
