package outbox

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/pipeline"
)

// Repository is an interface for manipulating the outgoing messages that
// comprise an incoming message's outbox.
type Repository interface {
	// LoadOutbox loads the undispatched outbound messages that were generated
	// when the given message was handled.
	//
	// ok is false if the message has not yet been handled.
	LoadOutbox(
		ctx context.Context,
		id ax.MessageID,
	) (ob []pipeline.OutboundMessage, ok bool, err error)

	// SaveOutbox saves a set of undispatched outbound messages that were
	// generated when the given message was handled. list of pending messages.
	SaveOutbox(
		ctx context.Context,
		tx ax.Transaction,
		id ax.MessageID,
		ob []pipeline.OutboundMessage,
	) error

	// MarkAsDispatched marks an OutboxMessage as dispatched, removing it from the
	// list of pending messages.
	MarkAsDispatched(
		ctx context.Context,
		tx ax.Transaction,
		m pipeline.OutboundMessage,
	) error
}
