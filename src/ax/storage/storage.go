package storage

import "context"

// Storage is an interface for accessing a transactional data store.
type Storage interface {
	// BeginTx starts a new transaction.
	BeginTx(ctx context.Context) (Tx, error)
}
