package ax

import "context"

// Transaction is an interface for performing an atomic unit of work on the
// underlying storage system.
type Transaction interface {
	Commit() error
	Rollback() error
}

// Storage is an interface to a transactional data store.
type Storage interface {
	// Tx begins a new transaction.
	Tx(ctx context.Context) (Transaction, error)
}
