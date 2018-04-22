package ax

import "context"

// TransactionState is a enumeration of the progress of a message transaction.
type TransactionState int

const (
	// TxMessageReceived indicates that a message has been received but
	// has not yet been successfully handled by the message handler.
	TxMessageReceived TransactionState = iota

	// TxMessageHandled indicates that the message handler has been
	// invoked successfully. If the message handler returns an error, the
	// transaction state does not change, it is simply retried.
	TxMessageHandled

	// TxOutboxDispatched indicates that the message has been handled, and
	// that any messages produced by the handler have been dispatched.
	TxOutboxDispatched
)

// MessageTransaction describes the state of a the message handling process
// for a specific message and handler.
type MessageTransaction interface {
	// // Begin starts a new storage-level transaction. It is invoked before the
	// // application-defined message handler is called.
	// //
	// // It panics if the transaction state is not TxMessageReceived.
	// BeginStorageTx() error
	//
	// // Commit sets the transaction state to TxMessageHandled and commits the
	// // storage-level transaction. It is invoked after the application-defined
	// // message handler is called.
	// //
	// // It panics if Begin() has not been called.
	// CommitStorageTx() error
	//
	// State returns the current state of the message transaction.
	State() TransactionState
	//
	// SetState sets the current state of the message transaction.
	SetState(TransactionState)
	//
	// // Enqueue appends a message to the outbox.
	// //
	// // It panics if Begin() has not been called.
	// Enqueue(Envelope) error

	// Outbox() ([]Envelope, error)
	//
	// // MarkDispatched sets the transaction state to TxOutboxDispatched.
	// //
	// // It panics if Begin() has been called but Commit() has not, or if the
	// // transaction state is already TxOutboxDispatched.
	// MarkDispatched() error

	Commit() error
	Rollback() error
}

// Persistence is an interface for loading and saving message
// transactions to a persistent data store.
type Persistence interface {
	// Tx starts (or resumes) a message transaction for the message
	// with the given ID.
	Tx(ctx context.Context, id MessageID) (MessageTransaction, error)
}
