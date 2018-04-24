package pipeline

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// OutboundStage is a step within an outbound message pipeline.
type OutboundStage interface {
	// DispatchMessage forwards an outbound message through the pipeline until
	// it ultimately is sent via a messaging transport.
	DispatchMessage(context.Context, OutboundMessage) error
}

// OutboundMessage is a container for a message being sent through the outbound
// pipeline.
type OutboundMessage struct {
	Operation OutboundOperation
	Envelope  ax.Envelope
}

// OutboundOperation is an enumeration of operations that can be performed to
// dispatch an outbound message.
type OutboundOperation int

const (
	// OpExecute is a transport operation that sends a command message to a
	// specific endpoint as determined by the routing configuration.
	OpExecute OutboundOperation = iota

	// OpPublish is a transport operation that sends an event message to its
	// subscribers.
	OpPublish
)
