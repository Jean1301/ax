package pipeline

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/outbox"
	"github.com/jmalloc/ax/src/ax/persistence"
)

type Executor struct {
	Out              Sink
	Persistence      persistence.Driver
	OutboxRepository outbox.Repository
	Transport        int
	Handlers         []ax.MessageHandler
}

func (x *Executor) Dispatch(ctx context.Context, m ax.Envelope) error {
	for {
		ok, err := x.step(ctx, m)
		if ok || err != nil {
			return err
		}
	}
}

func (x *Executor) step(ctx context.Context, m ax.Envelope) (bool, error) {
	tx, err := x.Persistence.Tx(ctx, m.MessageID)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	switch tx.State() {
	case persistence.TxMessageReceived:
		return false, x.handle(ctx, tx, m)
	case persistence.TxMessageHandled:
		return x.publish(ctx, tx)
	case persistence.TxOutboxDispatched:
		return true, nil
	}

	return false, fmt.Errorf("unrecognised transaction state: %d", tx.State())
}

func (x *Executor) handle(
	ctx context.Context,
	tx persistence.Transaction,
	m ax.Envelope,
) error {
	var mc ax.MessageContext

	for _, h := range d.Handlers {
		if err := h.HandleMessage(mc, m.Message); err != nil {
			return false, err
		}
	}

	tx.SetState(TxMessageHandled)
	return false, tx.Commit()
}

func (x *Executor) publish(ctx context.Context, tx ax.MessageTransaction) (bool, error) {
}
