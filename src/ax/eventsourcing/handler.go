package eventsourcing

// // MessageHandler is an ax.MessageHandler dispatches commands to an
// // ax.AggregateMessageHandler and persists the changes as a stream of events.
// type MessageHandler struct {
// 	Repository *Repository
// 	Handler    ax.AggregateCommandHandler
// }
//
// // MessageTypes returns the set of messages that the handler can handle.
// func (h *MessageHandler) MessageTypes() ax.MessageTypeSet {
// 	return h.Handler.MessageTypes()
// }
//
// // Handle handles a message.
// func (h *MessageHandler) HandleMessage(ctx ax.MessageContext, m ax.Message) error {
// 	c, ok := m.(ax.Command)
// 	if !ok {
// 		return ax.UnexpectedMessageError{Message: m}
// 	}
//
// 	ac := &aggregateCommandContext{
// 		Context: ctx,
// 		cause:   ctx.MessageEnvelope(),
// 	}
//
// 	id := h.Handler.MapToAggregate(c)
// 	agg := h.Handler.InitialState().(Aggregate)
// 	agg.SetAggregateID(id)
//
// 	if err := h.Repository.Load(ac, agg, &ac.progress); err != nil {
// 		return err
// 	}
//
// 	if err := h.Handler.HandleMessage(ac, agg, c); err != nil {
// 		if _, ok := err.(ax.DomainError); ok {
// 			// TODO: log
// 			return nil
// 		}
//
// 		return err
// 	}
//
// 	events := make([]ax.Envelope, len(ac.progress.Outbox)+1)
// 	copy(events, ac.progress.Outbox)
//
// 	events = append(
// 		events,
// 		ac.cause.New(
// 			&CommandHandled{
// 				CommandId: ac.cause.MessageID.Get(),
// 			},
// 		),
// 	)
//
// 	ac.progress.IsHandled = true
//
// 	ok, err := h.Repository.Save(ctx, agg, &ac.progress, events)
// 	if err != nil {
// 		return err
// 	}
//
// 	if !ok {
// 		// TODO: conflict occurred, retry
// 		return errors.New("conflict")
// 	}
//
// 	for _, env := range ac.progress.Outbox {
// 		// rmq.Publish(env)
// 		fmt.Println(env)
// 	}
//
// 	ac.progress.IsPublished = true
// 	ok, err = h.Repository.Save(ctx, agg, &ac.progress, []ax.Envelope{
// 		ac.cause.New(
// 			&CommandEventsPublished{
// 				CommandId: ac.cause.MessageID.Get(),
// 			},
// 		),
// 	})
// 	if err != nil {
// 		return err
// 	}
//
// 	if !ok {
// 		// TODO: conflict occurred, retry
// 		return errors.New("conflict")
// 	}
//
// 	// rmq.Ack(msg)
// 	fmt.Println("ACK")
//
// 	return nil
// }
//
// // aggregateCommandContext is an implementation of ax.AggregateCommandContext
// // for event-sourced aggregates.
// type aggregateCommandContext struct {
// 	context.Context
//
// 	cause    ax.Envelope
// 	agg      Aggregate
// 	progress ax.MessageProgress
// }
//
// // TODO: does it make sense to put the progress on all message contexts?
// // MessageEnvelope returns the message envelope of the command.
// // func (c *aggregateCommandContext) MessageProgress() *ax.MessageProgress {
// // 	return c.progress
// // }
//
// // MessageEnvelope returns the message envelope of the command.
// func (c *aggregateCommandContext) MessageEnvelope() ax.Envelope {
// 	return c.cause
// }
//
// // PublishEvent enqueues events to be published.
// func (c *aggregateCommandContext) PublishEvent(m ax.Event) {
// 	c.agg.Apply(m)
// 	c.progress.Outbox = append(
// 		c.progress.Outbox,
// 		c.cause.New(m),
// 	)
// }
