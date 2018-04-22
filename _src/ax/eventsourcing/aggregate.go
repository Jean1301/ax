package eventsourcing

import (
	"github.com/jmalloc/ax/src/ax"
)

type Aggregate interface {
	// Apply updates the aggregate's state to reflect the occurrence of ev.
	Apply(ax.Event)

	Handle(MessageContext, ax.Command) error
}

type MessageContext interface {
	Envelope() ax.Envelope
	Record(ax.Event)
}
