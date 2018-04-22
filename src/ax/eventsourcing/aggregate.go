package eventsourcing

import "github.com/jmalloc/ax/src/ax"

// Aggregate is an aggregate that uses event sourcing.
type Aggregate interface {
	ax.Aggregate

	// Apply updates the aggregate's state to reflect the occurrence of m.
	Apply(m ax.Event)
}
