package bus

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// InboundMessage is a container for a message that is received via a message
// transport.
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

// OutboundMessage is a container for a message that is sent via a message
// transport.
type OutboundMessage struct {
	Operation       OutboundOperation
	UnicastEndpoint string
	Envelope        ax.Envelope
}

// OutboundOperation is an enumeration of operations that can be performed to
// dispatch an outbound message.
type OutboundOperation int

const (
	// OpSendUnicast is an outbound transport operation that sends a message to a
	// specific endpoint as determined by the outbound message's UnicastEndpoint property.
	OpSendUnicast OutboundOperation = iota

	// OpSendMulticast is an outbound transport operation that sends a message to all of its
	// subscribers.
	OpSendMulticast
)
