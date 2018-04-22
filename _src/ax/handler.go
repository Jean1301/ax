package ax

import "context"

type HandlerContext interface {
	context.Context

	Send(Command) error
	Publish(Event) error
}

type MessageHandler interface {
	MessageTypes() []MessageType
	Handle(HandlerContext, Message) error
}

type MessageHandlerStage struct {
	Handlers []MessageHandler
}

func (s *MessageHandlerStage) Initialize(ctx context.Context, ep *Endpoint) error {
	for _, h := range s.Handlers {
		for _, mt := range h.MessageTypes() {
			if err := ep.Transport.Subscribe(ctx, mt); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *MessageHandlerStage) Process(ctx context.Context, m InboundMessage) error {
	var hc HandlerContext

	for _, h := range s.Handlers {
		if err := h.Handle(hc, m.Message); err != nil {
			return err
		}
	}

	return nil
}
