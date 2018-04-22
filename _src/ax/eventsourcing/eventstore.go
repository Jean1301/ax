package eventsourcing

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// EventStore is an interface for reading and writing events to/from a persisted
// store.
type EventStore interface {
	// Append stores new events for an aggregate at a specific revision.
	// If rev is not the current revision of the aggregate the append is
	// aborted and ok is false.
	Append(
		ctx context.Context,
		id string,
		rev uint64,
		events []ax.Event,
	) (ok bool, err error)

	// Open returns an EventStream used to read the events produced by a
	// specific aggregate as of a specific revision.
	Open(
		ctx context.Context,
		id string,
		rev uint64,
	) (EventStream, error)
}

// EventStream is an interface for reading an ordered stream of events from an
// event store.
type EventStream interface {
	// Next advances the stream to the next event.
	// It returns false if there are no more events.
	Next(ctx context.Context) (bool, error)

	// Get returns the event at the current location in the stream.
	Get(ctx context.Context) (ax.Envelope, error)

	// Close closes the stream.
	Close() error
}
