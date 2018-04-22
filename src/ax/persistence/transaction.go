package persistence

// TransactionState is a enumeration of the progress of a message transaction.
type TransactionState string

const (
	// TxMessageReceived indicates that a message has been received but
	// has not yet been successfully handled by the message handler.
	TxMessageReceived TransactionState = "message-received"

	// TxMessageHandled indicates that the message handler has been
	// invoked successfully. If the message handler returns an error, the
	// transaction state does not change, it is simply retried.
	TxMessageHandled TransactionState = "message-handled"

	// TxOutboxDispatched indicates that the message has been handled, and
	// that any messages produced by the handler have been dispatched.
	TxOutboxDispatched TransactionState = "outbox-dispatched"
)

type Transaction interface {
	State() TransactionState
	SetState(TransactionState)
	Commit() error
	Rollback() error
}
