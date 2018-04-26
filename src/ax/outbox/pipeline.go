package outbox

import (
	"context"

	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// InboundStage is a pipeline.InboundStage that provides message idempotency
// guarantees to the next pipeline stage by using the outbox pattern.
type InboundStage struct {
	Repository Repository
	Next       bus.InboundPipeline
}

// Initialize is called when the transport is initialized.
func (s *InboundStage) Initialize(ctx context.Context, t bus.Transport) error {
	return s.Next.Initialize(ctx, t)
}

// DispatchMessage passes m to each handler in s.Handler and persists the
// messages they produce to an outbox, before dispatching them via o.
func (s *InboundStage) DispatchMessage(
	ctx context.Context,
	out bus.OutboundDispatcher,
	m bus.InboundMessage,
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
		if err := s.dispatch(ctx, out, om); err != nil {
			return err
		}
	}

	return nil
}

func (s *InboundStage) forward(
	ctx context.Context,
	m bus.InboundMessage,
) ([]bus.OutboundMessage, error) {
	tx, err := persistence.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var out bus.MessageCollector

	if err := s.Next.DispatchMessage(
		persistence.WithTx(ctx, tx),
		&out,
		m,
	); err != nil {
		return nil, err
	}

	if err := s.Repository.SaveOutbox(
		ctx,
		tx,
		m.Envelope.MessageID,
		out.Messages,
	); err != nil {
		return nil, err
	}

	return out.Messages, tx.Commit()
}

// dispatch dispatches m via o and then records the fact that the message has
// been dispatched in the outbox.
func (s *InboundStage) dispatch(
	ctx context.Context,
	out bus.OutboundDispatcher,
	m bus.OutboundMessage,
) error {
	if err := out.DispatchMessage(ctx, m); err != nil {
		return err
	}

	tx, err := persistence.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.Repository.MarkAsDispatched(ctx, tx, m); err != nil {
		return err
	}

	return tx.Commit()
}
