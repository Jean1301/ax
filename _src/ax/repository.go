package ax

// Repository is an interface for loading and saving aggregates from a
// persistent data store.
type Repository interface {
	// LoadAggregate populates agg with data fetched from the store and sets its
	// revision.
	//
	// It returns an error if there is a problem communicating with the store.
	//
	// It is NOT an error if the aggregate does not exist, the aggregate will
	// simply remain at revision zero.
	LoadAggregate(ctx AggregateCommandContext, agg Aggregate) error

	// SaveAggregate persists an aggregate to the store.
	//
	// It returns an error if there is a problem communicating with the store.
	// ok is true if the changes are persisted successfully, or false if the
	// change was rejected because agg.Revision() is out of date.
	SaveAggregate(ctx AggregateCommandContext, agg Aggregate) (ok bool, err error)
}
