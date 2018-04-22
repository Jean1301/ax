package stateful

import "github.com/jmalloc/ax/src/ax"

// Revision is a state version, used for optimistic concurrency control.
type Revision uint64

// Instance is an instance of a stateful message handler's state data.
type Instance interface {
	// ID returns a unique identifier for the instance.
	// It panics if the ID has not been set.
	ID() ax.ID

	// SetID sets the ID for the instance.
	// It panics if the ID has already been set.
	SetID(ax.ID)

	// Revision returns the current revision of the instance.
	Revision() Revision

	// SetRevision sets the current revision of the instance.
	SetRevision(Revision)
}

// EventSourcedInstance is an instance that is built from a stream of events.
type EventSourcedInstance interface {
	Instance
	Apply(ax.Event)
}

// InstanceBehavior is an embeddable struct that implements the Instance
// interface.
type InstanceBehavior struct {
	id  *ax.ID
	rev Revision
}

// ID returns a unique identifier for the instance.
// It panics if the ID has not been set.
func (i *InstanceBehavior) ID() ax.ID {
	if i.id == nil {
		panic("instance ID has not been set")
	}

	return *i.id
}

// SetID sets the ID for the instance.
// It panics if the ID has already been set.
func (i *InstanceBehavior) SetID(id ax.ID) {
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
