package ax

import (
	"fmt"
)

// MessageHandler is an interface for types that handle messages.
type MessageHandler interface {
	// MessageTypes returns the set of messages that the handler can handle.
	MessageTypes() MessageTypeSet

	// HandleMessage handles a message.
	HandleMessage(MessageContext, Message) error
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
