package ax

import "context"

// Handler listens to and handles arbitrary messages.
type Handler interface {
	// MessageTypes returns the message types that the handler listens to.
	MessageTypes() MessageTypeSet

	// Handle handles a message.
	Handle(HandlerContext, Message) error
}

// HandlerContext provides access to the messaging system within the context
// of handling a particular message.
type HandlerContext interface {
	context.Context

	// Envelope returns the message envelope.
	Envelope() Envelope

	// Execute enqueues a command for execution.
	Execute(Command)

	// Publish publishes an event.
	Publish(Event)
}

// UnexpectedMessage is returned by a message handler when it is passed a
// message that it was not expecting to receive.
type UnexpectedMessage struct {
	Message Message
}

func (e UnexpectedMessage) Error() string {
	return "unexpected message type: " + TypeOf(e.Message).Name
}
