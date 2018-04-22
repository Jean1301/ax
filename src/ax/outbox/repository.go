package outbox

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// Repository is an interface for loading and saving message outboxes from/to a
// persistent data store.
type Repository interface {
	// LoadOutbox loads the outbox for the given message ID.
	LoadOutbox(
		ctx context.Context,
		tx persistence.Transaction,
		id ax.MessageID,
	) (Outbox, error)

	// SaveOutbox saves the outbox for the given message ID.
	SaveOutbox(
		ctx context.Context,
		tx persistence.Transaction,
		id ax.MessageID,
		ob Outbox,
	) error

	// ClearOutbox removes all messages in the outbox for the given message ID.
	ClearOutbox(
		ctx context.Context,
		tx persistence.Transaction,
		id ax.MessageID,
	) error
}
