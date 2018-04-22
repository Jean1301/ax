package ax

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/src/ax/persistence"
)

// MessageHandler is an interface for types that handle messages.
type MessageHandler interface {
	// MessageTypes returns the set of messages that the handler can handle.
	MessageTypes() MessageTypeSet

	// HandleMessage handles a message.
	HandleMessage(MessageContext, Message) error
}

// MessageContext is a specialization of context.Context used by message
// handlers.
//
// It carries information about the messaging behing handled, and allows the
// handler to produce new messages.
type MessageContext interface {
	context.Context

	// MessageEnvelope returns the envelope containing the message being handled.
	MessageEnvelope() Envelope

	// MessageTransaction returns the underlying message transaction.
	Transaction() persistence.Transaction

	// ExecuteCommand enqueues a command to be executed.
	ExecuteCommand(Command) error

	// PublishEvent enqueues events to be published.
	PublishEvent(Event) error
}

// UnexpectedMessageError is an error returned by a handler when it receives a
// message of an unexpected type.
type UnexpectedMessageError struct {
	Message Message
}

func (e UnexpectedMessageError) Error() string {
	return fmt.Sprintf(
		"unexpected %s message, %s",
		TypeOf(e.Message).Name,
		e.Message.Description(),
	)
}
