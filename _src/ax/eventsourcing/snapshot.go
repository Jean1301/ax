package eventsourcing

import (
	"context"

	"github.com/jmalloc/dogma/src/dogma"
)

// SnapshotStore loads and saves snapshots of event-sourced aggregates.
type SnapshotStore interface {
	// Load populates an aggregate from the latest available snapshot.
	Load(ctx context.Context, agg Aggregate) error

	// Save persists a new snapshot of an aggregate.
	Save(ctx context.Context, agg Aggregate) error
}

// SnapshotPolicy is a function that returns true if a new snapshot of an
// aggregate should be stored.
type SnapshotPolicy func(agg Aggregate, events []dogma.Event) bool

// SnapshotEveryN returns a SnapshotPolicy that saves snapshots after every N
// revisions.
func SnapshotEveryN(n uint64) SnapshotPolicy {
	return func(agg Aggregate, events []dogma.Event) bool {
		count := len(events)
		before := agg.Revision()
		after := before + uint64(count)
		return (before / n) != (after / n)
	}
}
