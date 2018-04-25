package storage

// Tx is an interface for performing an atomic unit of work on the underlying
// storage system.
type Tx interface {
	Commit() error
	Rollback() error
}
