package ax

// StatefulHandler is a message handler that has some form of persistent
// state which is managed by the framework.
//
// Each stateful handler can manage many "instances" of the state.
type StatefulHandler interface {
	// MessageTypes returns the message types that the handler listens to.
	//
	// Two sets are returned. The first is the set of messages that trigger a
	// new "instance" of the state to be created. The second is the set of
	// messages that should be routed to existing instances.
	MessageTypes() (trig MessageTypeSet, res MessageTypeSet)

	// InitialState returns a new instance of the state managed by the handler.
	InitialState() Instance

	// InstanceKey and MessageKey return "mapping keys" used to correlate
	// incoming messages the instances that need to receive them.
	//
	// For each incoming message m, MappingKeyForMessage(m) is compared to
	// MappingKeyForInstance(TypeOf(m), i). If the two values are equal,
	// HandleMessage() is called with the message and instance.
	InstanceKey(MessageType, Instance) (string, bool)
	MessageKey(Message) (string, bool)

	HandleMessage(HandlerContext, Message, Instance) error
	HandleNotFound(HandlerContext, Message) error
}

// StatefulAdaptor wraps a StatefulHandler to provide the Handler interface.
type StatefulAdaptor struct {
	Next StatefulHandler
}

// MessageTypes returns the message types that the handler listens to.
func (a *StatefulAdaptor) MessageTypes() MessageTypeSet {
	trig, res := a.Next.MessageTypes()
	return trig.Union(res)
}

// Handle handles a message.
func (a *StatefulAdaptor) Handle(ctx HandlerContext, m Message) error {
	mt := TypeOf(m)
	trig, res := a.Next.MessageTypes()

	if trig.Has(mt) {
		i := a.Next.InitialState()
		a.Next.HandleMessage(ctx, m, i)
	}

	if res.Has(mt) {

	}

	return UnexpectedMessage{m}
}

// Revision is a state version, used for optimistic concurrency control.
type Revision uint64

// Instance is an instance of a stateful message handler's state data.
type Instance interface {
	// ID returns a unique identifier for the instance.
	// It panics if the ID has not been set.
	ID() ID

	// SetID sets the ID for the instance.
	// It panics if the ID has already been set.
	SetID(ID)

	// Revision returns the current revision of the instance.
	Revision() Revision

	// SetRevision sets the current revision of the instance.
	SetRevision(Revision)
}

// MessageSourcedInstance is an instance that can be rebuilt from a stream of
// messages.
type MessageSourcedInstance interface {
	Instance
	Apply(Message)
}

// InstanceBehavior is an embeddable struct that implements the Instance
// interface.
type InstanceBehavior struct {
	id  *ID
	rev Revision
}

// ID returns a unique identifier for the instance.
// It panics if the ID has not been set.
func (i *InstanceBehavior) ID() ID {
	if i.id == nil {
		panic("instance ID has not been set")
	}

	return *i.id
}

// SetID sets the ID for the instance.
// It panics if the ID has already been set.
func (i *InstanceBehavior) SetID(id ID) {
	if i.id != nil {
		panic("aggregate ID has already been set")
	}

	i.id = &id
}

// Revision returns the current revision of the instance.
func (i *InstanceBehavior) Revision() Revision {
	return i.rev
}

// SetRevision sets the current revision of the instance.
func (i *InstanceBehavior) SetRevision(rev Revision) {
	i.rev = rev
}
