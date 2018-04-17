package ax

import (
	"reflect"

	"github.com/golang/protobuf/proto"
)

// Envelope is a container for a message and its meta-data.
type Envelope struct {
	MessageID string
	Headers   map[string]string
	Message   Message
}

// Message is a unit of communication.
type Message interface {
	proto.Message

	// Description returns a human-readable description of the message.
	//
	// Assume that the description will be used inside log messages or displayed
	// in audit logs.
	//
	// Follow the same conventions as for error messages:
	// https://github.com/golang/go/wiki/CodeReviewComments#error-strings
	//
	// Messages that implement DomainMessage should return a description that
	// makes sense to a person familiar with the business domain.
	Description() string
}

// Command is a message that requests some action take place.
//
// Commands are always sent to a single handler within a single end-point.
// Commands may optionally have an associated reply message.
type Command interface {
	Message

	// IsCommand() is a "marker method" used to indicate that a message is
	// intended to be used as a command.
	IsCommand()
}

// Event is a message that indicates some action has already taken place.
//
// Events are published by one endpoint and (potentially) consumed by many.
type Event interface {
	Message

	// IsEvent() is a "marker method" used to indicate that a message is
	// intended to be used as an event.
	IsEvent()
}

var (
	commandType = reflect.TypeOf((*Command)(nil)).Elem()
	eventType   = reflect.TypeOf((*Event)(nil)).Elem()
)
