package eventsourcing

import "github.com/jmalloc/ax/src/ax"

// Aggregate is an aggregate that uses event sourcing.
type Aggregate interface {
	ax.Aggregate

	// Revision returns the revision at which the aggregate was loaded.
	Revision() Revision

	// SetRevision sets the aggregate's revision.
	SetRevision(Revision)

	// Apply updates the aggregate's state to reflect the occurrence of m.
	Apply(m ax.Event)
}

// Revision is the 'version' of an aggregate.
type Revision uint64
