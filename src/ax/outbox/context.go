package outbox

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

func WithTransaction(ctx context.Context, tx ax.Transaction) context.Context {
	return context.WithValue(ctx, contextKey("tx"), tx)
}

func GetTransaction(ctx context.Context) (ax.Transaction, bool) {
	tx := ctx.Value(contextKey("tx"))

	if tx == nil {
		return nil, false
	}

	return tx.(ax.Transaction), true
}

type contextKey string
