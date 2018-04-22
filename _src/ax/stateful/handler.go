package stateful

import "github.com/jmalloc/ax/src/ax"

// MessageHandler is a message handler that has some form of persistent
// state which is managed by the framework.
//
// Each stateful handler can manage many "instances" of the state.
type MessageHandler interface {
	// MessageTypes returns the message types that the handler listens to.
	//
	// Two sets are returned. The first is the set of messages that cause a
	// new "instance" of the state to be created. The second is the set of
	// messages that should be routed to existing instances.
	MessageTypes() (ax.MessageTypeSet, ax.MessageTypeSet)

	// InitialState returns a new instance of the state managed by the handler.
	InitialState() Instance

	// InstanceKey and MessageKey return "mapping keys" used to correlate
	// incoming messages the instances that need to receive them.
	//
	// For each incoming message m, MappingKeyForMessage(m) is compared to
	// MappingKeyForInstance(TypeOf(m), i). If the two values are equal,
	// HandleMessage() is called with the message and instance.
	InstanceKey(ax.MessageType, Instance) (string, bool)
	MessageKey(ax.Message) (string, bool)

	HandleMessage(ax.HandlerContext, ax.Message, Instance) error
	HandleNotFound(ax.HandlerContext, ax.Message) error
}
