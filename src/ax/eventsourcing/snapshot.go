package eventsourcing

// SnapshotStore loads and saves snapshots of event-sourced aggregates.
// type SnapshotStore interface {
// 	// Load populates an aggregate from the latest available snapshot.
// 	//
// 	// If there is no snapshot available, agg is not modified and nil
// 	// is returned.
// 	Load(ctx ax.AggregateCommandContext, agg Aggregate, p *ax.MessageProgress) error
//
// 	// Save persists a new snapshot of an aggregate.
// 	Save(ctx ax.AggregateCommandContext, agg Aggregate, p ax.MessageProgress) error
// }
