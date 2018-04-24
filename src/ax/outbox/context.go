package outbox

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// WithTransaction returns a new context derived from parent that contains
// the outbox storage transaction.
//
// The transaction can be retrieved from the returned context with
// GetTransaction().
func WithTransaction(parent context.Context, tx ax.Transaction) context.Context {
	return context.WithValue(parent, contextKey("tx"), tx)
}

// GetTransaction returns the outbox storage transaction stored in ctx.
//
// It returns false if ctx does not contain an outbox transaction.
//
// The outbox transaction is made available via the context so that
// application-defined message handlers always have a way to perform some
// additional storage operations atomically with the
// Repository.SaveOutbox() operation.
//
// Care should be taken not to commit or rollback the transaction within
// the message handler.
func GetTransaction(ctx context.Context) (ax.Transaction, bool) {
	tx := ctx.Value(contextKey("tx"))

	if tx == nil {
		return nil, false
	}

	return tx.(ax.Transaction), true
}

type contextKey string
