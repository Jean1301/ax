package transport

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Transport is an interface for sending and receiving messages
// between endpoints.
type Transport interface {
	// Initialize prepares the transport for communication.
	//
	// ep is the name of the endpoint. s is the set of message types
	// that the endpoint wishes to receive.
	Initialize(ctx context.Context, ep string, s ax.MessageTypeSet) error

	// Receive returns the next message from the transport.
	// It blocks until a message is availble, or ctx is canceled.
	Receive(ctx context.Context) (*InboundMessage, error)

	// Send sends a message to a specific endpoint.
	Send(ctx context.Context, ep string, m *OutboundMessage) error

	// Publish multicasts a message to any endpoints that are interested in
	// receiving it.
	Publish(ctx context.Context, m *OutboundMessage) error
}

// InboundMessage is a message received via a transport.
type InboundMessage struct {
	Endpoint string
	Envelope ax.Envelope

	Done  func(context.Context) error
	Retry func(context.Context) error
}

// OutboundMessage is a message sent via a transport.
type OutboundMessage struct {
	Parent   InboundMessage
	Envelope ax.Envelope
}
