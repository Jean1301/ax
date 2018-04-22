package ax

import "errors"

// type AggregateExecutor struct {
// 	Handler AggregateCommandHandler
// }
//
// func (x AggregateExecutor) HandleMessage(ctx AggregateCommandContext, c Command) error {
// 	var p AggregatePersistence
//
// 	id := x.Handler.MapToAggregate(c)
// 	agg := x.Handler.InitialState()
// 	agg.SetAggregateID(id)
//
// 	if err := p.StartOrResumeTransaction(ctx, agg); err != nil {
// 		return err
// 	}
//
// 	switch ctx.MessageTransaction().State {
// 	case TxInboundMessageReceived:
// 		if err := x.Handler.HandleMessage(ctx, agg, c); err != nil {
// 			return err
// 		}
//
// 		p.PersistAggregateChanges(ctx, agg)
// 		fallthrough
//
// 	case TxInboundMessageHandled:
// 		// publish to rabbit
// 		return p.CompleteTransaction(ctx)
//
// 	case TxOutboundMessagesSent:
// 		return nil
//
// 	default:
// 		return errors.New("unknown transaction state")
// 	}
// }

type AggregateExecutor struct {
	Handler    AggregateCommandHandler
	Repository Repository
}

func (x AggregateExecutor) HandleMessage(ctx AggregateCommandContext, c Command) error {
	id := x.Handler.MapToAggregate(c)
	agg := x.Handler.InitialState()
	agg.SetAggregateID(id)

	mtx := ctx.MessageTransaction()
	mtx.MessageID = ctx.MessageEnvelope().MessageID

	if err := x.Repository.LoadAggregate(ctx, agg); err != nil {
		return err
	}

	switch mtx.State {
	case TxInboundMessageReceived:
		if err := x.Handler.HandleMessage(ctx, agg, c); err != nil {
			return err
		}

		mtx.State = TxInboundMessageHandled

		ok, err := x.Repository.SaveAggregate(ctx, agg)
		if err != nil {
			return err
		}
		if !ok {
			// conflict
		}

		fallthrough

	case TxInboundMessageHandled:
		// publish to rabbit rmq.Publish(mtx.Outbox...)
		mtx.State = TxOutboundMessagesSent

		ok, err := x.Repository.SaveAggregate(ctx, agg)
		if err != nil {
			return err
		}
		if !ok {
			// conflict
		}

	case TxOutboundMessagesSent:
		return nil

	default:
		return errors.New("unknown transaction state")
	}

	return nil
}
