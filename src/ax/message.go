package ax

import (
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax/ident"
)

// MessageID uniquely identifies a message.
type MessageID struct {
	ident.ID
}

// Envelope is a container for a message and its meta-data.
type Envelope struct {
	MessageID     MessageID
	CorrelationID MessageID
	CausationID   MessageID
	Time          time.Time
	Message       Message
}

// New returns a new message envelope that contains m.
// m is "caused by" e.Message.
func (e Envelope) New(m Message) Envelope {
	env := Envelope{
		CorrelationID: e.CorrelationID,
		CausationID:   e.MessageID,
		Time:          time.Now(),
	}

	env.MessageID.GenerateUUID()

	return env
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
