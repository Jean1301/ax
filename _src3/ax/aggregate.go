package ax

//
// // AggregateID uniquely identifies an aggregate.
// type AggregateID struct {
// 	ident.ID
// }
//
// // Aggregate is an interface for application-defined aggregates.
// type Aggregate interface {
// 	// AggregateID returns the aggregate's unique identifier.
// 	// It panics if the aggregate ID has not been set.
// 	AggregateID() AggregateID
//
// 	// SetAggregateID sets the ID of the aggregate represented by this value.
// 	// It panics if the aggregate ID has already been set.
// 	SetAggregateID(AggregateID)
// }
//
// // AggregateCommandHandler is an interface for executing domain commands against
// // aggregates.
// type AggregateCommandHandler interface {
// 	// MessageTypes returns the set of commands that are routed to the an
// 	// aggregate of this type.
// 	MessageTypes() MessageTypeSet
//
// 	// MapToAggregate returns the aggregate ID that handles the given message.
// 	MapToAggregate(Command) AggregateID
//
// 	// InitialState returns a new aggregate value.
// 	InitialState() Aggregate
//
// 	// HandleMessage mutates an aggregate based on the given command.
// 	HandleMessage(AggregateCommandContext, Aggregate, Command) error
// }
//
// // AggregateCommandContext is a specialization of context.Context used by
// // aggregate command handlers.
// //
// // It carries information about the messaging behing handled, and allows the
// // handler to produce new messages.
// //
// // Note that unlike a generic message handler, an aggregate can not produce
// // commands, only events.
// type AggregateCommandContext interface {
// 	context.Context
//
// 	// MessageEnvelope returns the message envelope of the incoming message.
// 	MessageEnvelope() Envelope
//
// 	// MessageTransaction returns the underlying message transaction.
// 	MessageTransaction() MessageTransaction
//
// 	// PublishEvent enqueues events to be published.
// 	PublishEvent(Event)
// }
//
// // DomainError is returned by an aggregate command handler when executing the
// // command would violate a domain invariant.
// type DomainError string
//
// func (e DomainError) Error() string {
// 	return "domain invariant violated, " + string(e)
// }
//
// // AggregateRepository is an interface for loading and saving aggregates
// // to a persistent data store.
// type AggregateRepository interface {
// 	// Load loads an aggregate from the store.
// 	//
// 	// agg is populated with the aggregate data according to it's ID.
// 	//
// 	// It returns an error only if there is a problem communicating with the
// 	// store. A non-existent aggregate is not an error.
// 	Load(ctx context.Context, agg Aggregate) error
//
// 	// Save persists aggregate changes to the store, as part of tx.
// 	//
// 	// It returns true if the aggregate is persisted correctly, or false if
// 	// aggregate has been changed by another process.
// 	//
// 	// It returns an error if there is a problem communicating with the store.
// 	//
// 	// It panics if the repository can not participate in tx because of a
// 	// difference in the storage layer.
// 	Save(ctx context.Context, tx MessageTransaction, agg Aggregate) (bool, error)
// }
//
// // AggregateHandlerAdaptor is a MessageHandler that dispatches commands to an
// // AggregateCommandHandler.
// type AggregateHandlerAdaptor struct {
// 	Repository AggregateRepository
// 	Next       AggregateCommandHandler
// }
//
// // MessageTypes returns the set of messages that the handler can handle.
// func (a *AggregateHandlerAdaptor) MessageTypes() MessageTypeSet {
// 	set := a.Next.MessageTypes()
//
// 	for _, mt := range set.Members() {
// 		if !mt.IsCommand() {
// 			panic(fmt.Sprintf(
// 				"aggregate command handler %s is attempting to listen to non-command %s",
// 				reflect.TypeOf(a.Next),
// 				mt.Name,
// 			))
// 		}
// 	}
//
// 	return set
// }
//
// func (a *AggregateHandlerAdaptor) HandleMessage(ctx MessageContext, m Message) error {
// 	cmd, ok := m.(Command)
// 	if !ok {
// 		return UnexpectedMessageError{Message: m}
// 	}
//
// 	id := a.Next.MapToAggregate(cmd)
//
// 	for {
// 		ok, err := a.handleCommand(ctx, id, cmd)
// 		if ok || err != nil {
// 			return err
// 		}
// 	}
// }
//
// func (a *AggregateHandlerAdaptor) handleCommand(
// 	ctx MessageContext,
// 	id AggregateID,
// 	m Command,
// ) (bool, error) {
// 	agg := a.Next.InitialState()
// 	agg.SetAggregateID(id)
//
// 	if err := a.Repository.LoadAggregate(ctx, agg); err != nil {
// 		return false, err
// 	}
//
// 	var ac AggregateCommandContext
//
// 	if err := a.Next.HandleMessage(ac, agg, m); err != nil {
// 		if _, ok := err.(DomainError); ok {
// 			// TODO: log
// 		} else {
// 			return false, err
// 		}
// 	}
//
// 	return a.Repository.SaveAggregate(
// 		ctx,
// 		ctx.MessageTransaction(),
// 		agg,
// 	)
// }
//
// // AggregateBehavior is an embeddable struct that partially implements the
// // Aggregate interface.
// type AggregateBehavior struct {
// 	id *AggregateID
// 	// rev Revision
// }
//
// // AggregateID returns the aggregate's unique identifier.
// // It panics if the aggregate ID has not been set.
// func (a *AggregateBehavior) AggregateID() AggregateID {
// 	if a.id == nil {
// 		panic("aggregate ID has not been set")
// 	}
//
// 	return *a.id
// }
//
// // SetAggregateID sets the ID of the aggregate represented by this value.
// func (a *AggregateBehavior) SetAggregateID(id AggregateID) {
// 	if a.id != nil {
// 		panic("aggregate ID has already been set")
// 	}
//
// 	a.id = &id
// }
//
// // // Revision returns the revision at which the aggregate was loaded.
// // func (a *AggregateBehavior) Revision() Revision {
// // 	return a.rev
// // }
// //
// // // SetRevision sets the aggregate's revision.
// // func (a *AggregateBehavior) SetRevision(rev Revision) {
// // 	a.rev = rev
// // }
//
// // DescribeAggregate returns a string description of an aggregate.
// func DescribeAggregate(a Aggregate) string {
// 	if s, ok := a.(fmt.Stringer); ok {
// 		return s.String()
// 	}
//
// 	return fmt.Sprintf(
// 		"aggregate<%s>[%s@%d]",
// 		reflect.TypeOf(a),
// 		a.AggregateID(),
// 		a.Revision(),
// 	)
// }
