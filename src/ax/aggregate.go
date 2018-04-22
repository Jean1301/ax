package ax

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jmalloc/ax/src/ax/ident"
)

// AggregateID uniquely identifies an aggregate.
type AggregateID struct {
	ident.ID
}

// Revision is the 'version' of an aggregate.
type Revision uint64

// Aggregate is an interface for application-defined aggregates.
type Aggregate interface {
	// AggregateID returns the aggregate's unique identifier.
	// It panics if the aggregate ID has not been set.
	AggregateID() AggregateID

	// SetAggregateID sets the ID of the aggregate represented by this value.
	// It panics if the aggregate ID has already been set.
	SetAggregateID(AggregateID)

	// Revision returns the revision at which the aggregate was loaded.
	Revision() Revision

	// SetRevision sets the aggregate's revision.
	SetRevision(Revision)
}

// AggregateCommandHandler is an interface for executing domain commands against
// aggregates.
type AggregateCommandHandler interface {
	// MessageTypes returns the set of commands that are routed to the an
	// aggregate of this type.
	MessageTypes() MessageTypeSet

	// MapToAggregate returns the aggregate ID that handles the given message.
	MapToAggregate(Command) AggregateID

	// InitialState returns a new aggregate value.
	InitialState() Aggregate

	// HandleMessage mutates an aggregate based on the given command.
	HandleMessage(AggregateCommandContext, Aggregate, Command) error
}

// AggregateCommandContext is a specialization of context.Context used by
// aggregate command handlers.
//
// It carries information about the messaging behing handled, and allows the
// handler to produce new messages.
//
// Note that unlike a generic message handler, an aggregate can not produce
// commands, only events.
type AggregateCommandContext interface {
	context.Context

	// MessageEnvelope returns the message envelope of the incoming message.
	MessageEnvelope() Envelope

	// PublishEvent enqueues events to be published.
	PublishEvent(Event)
}

// DomainError is returned by an aggregate command handler when executing the
// command would violate a domain invariant.
type DomainError string

func (e DomainError) Error() string {
	return "domain invariant violated, " + string(e)
}

// AggregateRepository is an interface for loading and saving aggregates
// to a persistent data store.
//
// It is a specialised form of TransactionRepository that allows for atomic
// changes to aggregate and transaction state.
type AggregateRepository interface {
	TransactionRepository

	// LoadAggregate populates an aggregate with data from the store.
	//
	// It returns an error only if there is a problem communicating with the
	// store. A non-existent aggregate is not an error.
	LoadAggregate(
		ctx context.Context,
		agg Aggregate,
	) error

	// SaveAggregateAndTx atomically persists an aggregate and a transaction to
	// the store.
	//
	// ok is false if either agg.Revision() or mtx.Revision do not match their
	// respective revisions in the store, in which case, no changes are
	// persisted.
	//
	// It returns an error if there is a problem communicating with the store.
	SaveAggregateAndTx(
		ctx context.Context,
		agg Aggregate,
		mtx *MessageTransaction,
	) (bool, error)
}

// AggregateBehavior is an embeddable struct that partially implements the
// Aggregate interface.
type AggregateBehavior struct {
	id  *AggregateID
	rev Revision
}

// AggregateID returns the aggregate's unique identifier.
// It panics if the aggregate ID has not been set.
func (a *AggregateBehavior) AggregateID() AggregateID {
	if a.id == nil {
		panic("aggregate ID has not been set")
	}

	return *a.id
}

// SetAggregateID sets the ID of the aggregate represented by this value.
func (a *AggregateBehavior) SetAggregateID(id AggregateID) {
	if a.id != nil {
		panic("aggregate ID has already been set")
	}

	a.id = &id
}

// Revision returns the revision at which the aggregate was loaded.
func (a *AggregateBehavior) Revision() Revision {
	return a.rev
}

// SetRevision sets the aggregate's revision.
func (a *AggregateBehavior) SetRevision(rev Revision) {
	a.rev = rev
}

// DescribeAggregate returns a string description of an aggregate.
func DescribeAggregate(a Aggregate) string {
	if s, ok := a.(fmt.Stringer); ok {
		return s.String()
	}

	return fmt.Sprintf(
		"aggregate<%s>[%s@%d]",
		reflect.TypeOf(a),
		a.AggregateID(),
		a.Revision(),
	)
}
