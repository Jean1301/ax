package outbox

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// Repository is an interface for manipulating the outgoing messages that
// comprise an incoming message's outbox.
type Repository interface {
	// LoadOutbox loads the unsent outbound messages that were produced when the
	// message identified by id was first delivered.
	//
	// ok is false if the message has not yet been successfully delivered.
	LoadOutbox(
		ctx context.Context,
		ds persistence.DataStore,
		id ax.MessageID,
	) (envs []endpoint.OutboundEnvelope, ok bool, err error)

	// SaveOutbox saves a set of unsent outbound messages that were produced
	// when the message identified by id was delivered.
	SaveOutbox(
		ctx context.Context,
		tx persistence.Tx,
		id ax.MessageID,
		envs []endpoint.OutboundEnvelope,
	) error

	// MarkAsSent marks a message as sent, removing it from the outbox.
	MarkAsSent(
		ctx context.Context,
		tx persistence.Tx,
		env endpoint.OutboundEnvelope,
	) error
}
