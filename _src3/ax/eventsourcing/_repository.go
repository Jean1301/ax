package eventsourcing

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax"
)

// Repository loads and saves event sourced aggregates.
type Repository struct {
	EventStore        EventStore
	SnapshotStore     SnapshotStore
	SnapshotFrequency ax.Revision
}

// Load fetches an aggregate from the store.
func (r *Repository) Load(
	ctx ax.AggregateCommandContext,
	agg Aggregate,
	p *ax.MessageProgress,
) error {
	if r.SnapshotStore != nil {
		if err := r.SnapshotStore.Load(ctx, agg, p); err != nil {
			return err
		}
	}

	return r.loadEvents(ctx, agg, p)
}

// Save m new events for an aggregate.
//
// It returns false if another process has changed the aggregate since it was
// loaded, in which case agg is invalid.
func (r *Repository) Save(
	ctx ax.AggregateCommandContext,
	agg Aggregate,
	p *ax.MessageProgress,
	events []ax.Envelope,
) (bool, error) {
	rev := agg.Revision()

	ok, err := r.EventStore.Append(
		ctx,
		agg.AggregateID(),
		rev,
		events,
	)
	if !ok || err != nil {
		return ok, err
	}

	agg.SetRevision(
		rev + ax.Revision(len(events)),
	)

	return true, r.saveSnapshot(ctx, rev, agg)
}

// loadEvents reads events from the store and applies the to agg.
func (r *Repository) loadEvents(ctx context.Context, agg Aggregate, p *ax.MessageProgress) error {
	stream, err := r.EventStore.Open(ctx, agg.AggregateID(), agg.Revision())
	if err != nil {
		return err
	}
	defer stream.Close()

	var rev ax.Revision

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

		switch m := env.Message.(type) {
		case *CommandHandled:
			p.IsHandled = true
		case *CommandEventsPublished:
			p.IsPublished = true
			p.Outbox = nil
		case ax.Event:
			p.Outbox = append(p.Outbox, env)
			agg.Apply(m)
			rev++
		default:
			return errors.New("non-event persisted to event store") // TODO
		}
	}

	agg.SetRevision(rev)
	return nil
}

// saveSnapshot stores a snapshot of the aggregate if the changes cause the
// revision to increase by enough to cross a snashot frequency "boundary".
func (r *Repository) saveSnapshot(ctx context.Context, prev ax.Revision, agg Aggregate) error {
	if r.SnapshotStore == nil || r.SnapshotFrequency == 0 {
		return nil
	}

	if r.SnapshotFrequency > 1 {
		bucketBefore := prev / r.SnapshotFrequency
		bucketAfter := agg.Revision() / r.SnapshotFrequency

		if bucketBefore == bucketAfter {
			return nil
		}
	}

	return r.SnapshotStore.Save(ctx, agg)
}
