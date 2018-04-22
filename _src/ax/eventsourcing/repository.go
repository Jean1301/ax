package eventsourcing

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/dogma/src/dogma"
)

// Repository loads and saves event sourced aggregates.
type Repository struct {
	EventStore     EventStore
	SnapshotStore  SnapshotStore
	SnapshotPolicy SnapshotPolicy
}

// Load fetches an aggregate from the store.
func (r *Repository) Load(ctx context.Context, agg Aggregate) error {
	if r.SnapshotStore != nil {
		if err := r.SnapshotStore.Load(ctx, agg); err != nil {
			return err
		}
	}

	return r.loadEvents(ctx, agg)
}

// Save persists an aggregate to the store.
func (r *Repository) Save(ctx context.Context, agg Aggregate, events []ax.Event) error {
	if len(events) == 0 {
		return nil
	}

	ok, err := r.EventStore.Append(ctx, agg.AggregateID(), agg.Revision(), events)
	if err != nil {
		return err
	} else if !ok {
		return dogma.ConflictError{Aggregate: agg}
	}

	n := len(events)
	agg.SetRevision(
		agg.Revision() + uint64(n),
	)

	return r.saveSnapshot(ctx, agg, events)
}

func (r *Repository) loadEvents(ctx context.Context, agg Aggregate) error {
	stream, err := r.EventStore.Open(ctx, agg.AggregateID(), agg.Revision())
	if err != nil {
		return err
	}
	defer stream.Close()

	var rev uint64

	for {
		ok, err := stream.Next(ctx)
		if err != nil {
			return err
		} else if !ok {
			break
		}

		env, err := stream.Get(ctx)
		if err != nil {
			return err
		}

		agg.Apply(env.Event)
		rev++
	}

	agg.SetRevision(rev)
	return nil
}

func (r *Repository) saveSnapshot(ctx context.Context, agg Aggregate, events []dogma.Event) error {
	if r.SnapshotStore == nil {
		return nil
	}

	if r.SnapshotPolicy == nil {
		return nil
	}

	if !r.SnapshotPolicy(agg, events) {
		return nil
	}

	// TODO: do we really want to propagate this error, and if not, what if it
	// keeps happening?
	return r.SnapshotStore.Save(ctx, agg)
}
