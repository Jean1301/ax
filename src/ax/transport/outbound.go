package transport

import "github.com/jmalloc/ax/src/ax"

// OutboundMessage is a container for a message being sent through the outbound
// pipeline.
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
