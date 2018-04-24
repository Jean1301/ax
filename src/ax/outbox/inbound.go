package outbox

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/pipeline"
)

// InboundStage is a pipeline.InboundStage that provides message idempotency
// guarantees to the next pipeline stage by using the outbox pattern.
type InboundStage struct {
	Next       pipeline.InboundStage
	Storage    ax.Storage
	Repository Repository
}

// DispatchMessage passes m to each handler in s.Handler and persists the
// messages they produce to an outbox, before dispatching them via o.
func (s *InboundStage) DispatchMessage(
	ctx context.Context,
	o pipeline.OutboundStage,
	m pipeline.InboundMessage,
) error {
	outbox, ok, err := s.Repository.LoadOutbox(
		ctx,
		m.Envelope.MessageID,
	)
	if err != nil {
		return err
	}

	if !ok {
		outbox, err = s.forward(ctx, m)
		if err != nil {
			return err
		}
	}

	for _, om := range outbox {
		if err := s.dispatch(ctx, o, om); err != nil {
			return err
		}
	}

	return nil
}

func (s *InboundStage) forward(
	ctx context.Context,
	m pipeline.InboundMessage,
) ([]pipeline.OutboundMessage, error) {
	tx, err := s.Storage.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var o OutboundStage

	if err := s.Next.DispatchMessage(
		WithTransaction(ctx, tx),
		&o,
		m,
	); err != nil {
		return nil, err
	}

	if err := s.Repository.SaveOutbox(
		ctx,
		tx,
		m.Envelope.MessageID,
		o.Outbox,
	); err != nil {
		return nil, err
	}

	return o.Outbox, tx.Commit()
}

// dispatch dispatches m via o and then records the fact that the message has
// been dispatched in the outbox.
func (s *InboundStage) dispatch(
	ctx context.Context,
	o pipeline.OutboundStage,
	m pipeline.OutboundMessage,
) error {
	if err := o.DispatchMessage(ctx, m); err != nil {
		return err
	}

	tx, err := s.Storage.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.Repository.MarkAsDispatched(ctx, tx, m); err != nil {
		return err
	}

	return tx.Commit()
}
