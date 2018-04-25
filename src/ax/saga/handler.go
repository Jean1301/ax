package saga

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/outbox"
)

// MessageHandler is an implementation of ax.MessageHandler that handles the
// persistence of saga instances before forwarding the message to a saga.
type MessageHandler struct {
	Repository Repository
	Saga       Saga
	Storage    ax.Storage

	triggers ax.MessageTypeSet
}

// MessageTypes returns the set of messages that the handler can handle.
//
// For sagas, this is the union of the message types that trigger new instances
// and the message types that are routed to existing instances.
func (h *MessageHandler) MessageTypes() ax.MessageTypeSet {
	triggers, others := h.Saga.MessageTypes()
	h.triggers = triggers

	return triggers.Union(others)
}

// HandleMessage loads a saga instance, passes m to the saga to be handled, and
// saves the changes to the saga instance.
//
// Changes to the saga are persisted within the outbox transaction if one is
// present in ctx. Otherwise, a new transaction is started using h.Storage.
func (h *MessageHandler) HandleMessage(ctx ax.MessageContext, m ax.Message) error {
	mt := ax.TypeOf(m)
	mk := h.Saga.MapMessage(m)
	si := h.Saga.InitialState()

	// attempt to find an existing saga instance from the message mapping key.
	ok, err := h.Repository.LoadSagaInstance(ctx, mt, mk, si)
	if err != nil {
		return err
	}

	// fetch the outbox transaction from the ctx, or start a new transaction
	// if none is present.
	tx, ownTx, err := h.getTransaction(ctx)
	if err != nil {
		return err
	}

	if ownTx {
		defer tx.Rollback()
	}

	hctx := WithTransaction(ctx, tx)

	// if no existing instance is found, and this message type does not produce
	// new instances, then the not-found handler is called.
	if !ok && !h.triggers.Has(mt) {
		return h.Saga.HandleNotFound(ctx, m)
	}


	// pass the message to the saga for handling.
	if err := h.Saga.HandleMessage(ctx, m, si); err != nil {
		return err
	}

	// save the changes to the saga and its mapping table.
	if err := h.Repository.SaveSagaInstance(
		ctx,
		tx,
		si,
		buildMappingTable(h.Saga, si),
	); err != nil {
		return err
	}

	if ownTx {
		return tx.Commit()
	}

	return nil
}

func (h *MessageHandler) getTransaction(ctx context.Context) (ax.Transaction, bool, error) {
	if tx, ok := outbox.GetTransaction(ctx); ok {
		return tx, false, nil
	}

	if h.storage == nil {
		return nil, false, errors.New(
			"cannot persist saga state, ",
		)
	}
		tx, err = h.Storage.Tx(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()
	}
}
