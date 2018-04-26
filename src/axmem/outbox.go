package axmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// OutboxRepository is an implementation of outbox.Repository that stores data
// within an in-memory transaction.
type OutboxRepository struct {
	m     sync.RWMutex
	boxes map[ax.MessageID]*outbox
}

type outbox struct {
	m        sync.RWMutex
	exists   bool
	messages []bus.OutboundMessage
}

// LoadOutbox loads the undispatched outbound messages that were generated
// when the given message was handled.
func (r *OutboxRepository) LoadOutbox(
	ctx context.Context,
	id ax.MessageID,
) ([]bus.OutboundMessage, bool, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	box, ok := r.boxes[id]
	if !ok {
		return nil, false, nil
	}

	box.m.RLock()
	defer box.m.RUnlock()

	return box.messages, box.exists, nil
}

// SaveOutbox saves a set of undispatched outbound messages that were
// generated when the given message was handled. list of pending messages.
func (r *OutboxRepository) SaveOutbox(
	ctx context.Context,
	tx persistence.Tx,
	id ax.MessageID,
	ob []bus.OutboundMessage,
) error {
	mtx, err := unwrapTx(tx)
	if err != nil {
		return err
	}

	r.m.Lock()
	defer r.m.Unlock()

	box, ok := r.boxes[id]
	if !ok {
		if r.boxes == nil {
			r.boxes = map[ax.MessageID]*outbox{}
		}

		box = &outbox{}
		r.boxes[id] = box
	}

	return mtx.Enlist(&box.m, func(commit bool) {
		if commit {
			box.exists = true
			box.messages = ob
		} else if !box.exists {
			r.m.Lock()
			defer r.m.Unlock()
			delete(r.boxes, id)
		}
	})
}

// MarkAsDispatched marks an OutboxMessage as dispatched, removing it from the
// list of pending messages.
func (r *OutboxRepository) MarkAsDispatched(
	ctx context.Context,
	tx persistence.Tx,
	m bus.OutboundMessage,
) error {
	mtx, err := unwrapTx(tx)
	if err != nil {
		return err
	}

	r.m.Lock()
	defer r.m.Unlock()

	box, ok := r.boxes[m.Envelope.CausationID]
	if !ok {
		return fmt.Errorf(
			"can not mark message %s as dispatched, no outbox is stored for causal message %s",
			m.Envelope.MessageID,
			m.Envelope.CausationID,
		)
	}

	return mtx.Enlist(&box.m, func(commit bool) {
		if commit {
			for i, v := range box.messages {
				if v.Envelope.MessageID == m.Envelope.MessageID {
					box.messages = append(
						box.messages[:i],
						box.messages[i+1:]...,
					)
				}
			}
		}
	})
}
