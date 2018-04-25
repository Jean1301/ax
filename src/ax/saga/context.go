package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/internal/transactionx"
)

// WithTransaction returns a new context derived from parent that contains
// the saga storage transaction.
//
// The transaction can be retrieved from the returned context with
// GetTransaction().
func WithTransaction(parent context.Context, tx ax.Transaction) context.Context {
	return transactionx.WithTransaction(parent, "saga", tx)
}

// GetTransaction returns the saga storage transaction stored in ctx.
//
// It returns false if ctx does not contain a saga transaction.
//
// The saga transaction is made available via the context so that
// application-defined message handlers always have a way to perform some
// additional storage operations atomically with the
// Repository.SaveSagaInstance() operation.
//
// Care should be taken not to commit or rollback the transaction within
// the message handler.
func GetTransaction(ctx context.Context) (ax.Transaction, bool) {
	return transactionx.GetTransaction(ctx, "saga")
}
