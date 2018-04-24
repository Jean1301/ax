package transport

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// InboundMessage is a container for a message being sent through the inbound
// pipeline.
type InboundMessage struct {
	DeliveryCount uint // zero = unknown
	Envelope      ax.Envelope

	Done func(context.Context, InboundOperation) error
}

// InboundOperation is an enumeration of operations that can be performed to
// an inbound message.
type InboundOperation int

const (
	// OpAck is an inbound transport operation that causes the inbound message
	// to be removed from the endpoint's queue.
	OpAck InboundOperation = iota

	// OpRetry is an inbound transport operation that causes the inbound message
	// to be retried.
	OpRetry

	// OpReject is an inbound transport operation that causes the inbound
	// message to be rejected. Depending on the transport, the message may be
	// moved to some form of error queue, or dropped completely.
	OpReject
)
