package storage

import (
	"context"
)

// WithTx returns a new context derived from parent that contains a storage
// transaction.
//
// The transaction can be retrieved from the context with GetTx().
func WithTx(parent context.Context, tx Tx) context.Context {
	return context.WithValue(parent, contextKey("tx"), tx)
}

// GetTx returns the transaction stored in ctx.
//
// If ctx does not contain a transaction then ok is false.
//
// Transactions are made available via the context so that application-defined
// message handlers can optionally perform some additional storage operations
// within the same transaction as infrastructure features such as the outbox
// system.
//
// Care should be taken not to commit or rollback the transaction within
// the message handler.
func GetTx(ctx context.Context) (tx Tx, ok bool) {
	v := ctx.Value(contextKey("tx"))

	if v != nil {
		tx = v.(Tx)
		ok = true
	}

	return
}

// WithStorage returns a new context derived from parent that contains
// a storage instance.
//
// The storage instance can be retrieved from the context with GetStorage().
func WithStorage(parent context.Context, s Storage) context.Context {
	return context.WithValue(parent, contextKey("storage"), s)
}

// GetStorage returns the Storage instance stored in ctx.
//
// If ctx does not contain

type contextKey string
