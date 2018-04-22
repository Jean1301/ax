package ax

import "context"

type Transport interface {
	Initialize(context.Context, *Endpoint) error
	Subscribe(context.Context, MessageType) error
	Send(context.Context, OutboundMessage)
	Receive(context.Context) (InboundMessage, error)
}

type InboundMessage struct {
	Message Message
	Ack     func() error
}

type OutboundMessage struct {
	Message Message
	Ack     func() error
}
