package bus

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Transport is an interface for communicating messages between endpoints.
type Transport interface {
	// Initialize sets up the transport to communicate as an endpoint named ep.
	Initialize(ctx context.Context, ep string) error

	// Subscribe instructs the transport to listen to multicast messages of the
	// given type.
	Subscribe(ctx context.Context, mt ax.MessageTypeSet) error

	// Send sends a message via the transport.
	Send(ctx context.Context, m OutboundMessage) error

	// Receive returns the next message that has been delivered to the endpoint.
	Receive(ctx context.Context) (InboundMessage, error)
}

// TransportPipelineStage is an outbound pipeline stage that sends a message
// over a transport.
type TransportPipelineStage struct {
	Transport Transport
}

// Initialize is called when the transport is initialized.
func (s *TransportPipelineStage) Initialize(ctx context.Context, t Transport) error {
	s.Transport = t
	return nil
}

// DispatchMessage forwards an outbound message through the pipeline until
// it ultimately is sent via a messaging transport.
func (s *TransportPipelineStage) DispatchMessage(ctx context.Context, m OutboundMessage) error {
	if m.Operation == OpSendUnicast && m.UnicastEndpoint == "" {
		m.UnicastEndpoint = unicastEndpointFor(m.Envelope.Message)
	}

	return s.Transport.Send(ctx, m)
}

// unicastEndpointFor returns the target endpoint for the given command.
//
// This is a temporary implementation that routes to an endpoint with the same
// name as the Protocol Buffers package.
//
// At some point it may become necessary to add an outbound pipeline stage that
// uses more complex routing logic.
func unicastEndpointFor(m ax.Message) string {
	return ax.TypeOf(m).PackageName()
}
