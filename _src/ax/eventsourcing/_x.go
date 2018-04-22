package eventsourcing

type Aggregate interface {
	ApplyEvent(ax.Event)
	HandleCommand(MessageContext, ax.Command) error
}

type AggregateX interface {
	CommandTypes() ax.MessageTypeSet
	InitialState() Aggregate
	MapToAggregateRoot(ax.Command) string
}

type AggregateHandler struct {
	X AggregateX
}
