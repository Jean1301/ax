package persistence

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

type Driver interface {
	Tx(ctx context.Context, id ax.MessageID) (Transaction, error)
}
