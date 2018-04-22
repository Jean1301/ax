package ax

import "context"

// TransactionState is a enumeration of the progress of a message transaction.
type TransactionState int

const (
	// TxUnknown indicates that the state of the transaction is unknown.
	// Probably because the value has not been initialized.
	TxUnknown TransactionState = iota

	// TxInboundMessageReceived indicates that a message has been received but
	// has not yet been successfully handled by the message handler.
	TxInboundMessageReceived

	// TxInboundMessageHandled indicates that the message handler has been
	// invoked successfully. If the message handler returns an error, the
	// transaction state does not change, it is simply retried.
	TxInboundMessageHandled

	// TxOutboundMessagesSent indicates that the message has been handled, and
	// that any messages produced by the handler have been sent.
	TxOutboundMessagesSent
)

// MessageTransaction describes the state of a the message handling process
// for a specific message and handler.
type MessageTransaction struct {
	// MessageID is the ID of the message that started the transaction.
	MessageID MessageID

	// Revision is the current revision of the transaction, used to implement
	// optimistic concurrency control.
	Revision Revision

	// State is the current transaction state.
	State TransactionState

	// Outbox is a collection of messages produced by the message handler.
	//
	// The outbox is only populated if the transaction state is
	// InboundMessageHandled. A message handler does not necessarily produce
	// any outbound messages.
	Outbox []Envelope
}

// TransactionRepository is an interface for loading and saving message
// transactions to a persistent data store.
type TransactionRepository interface {
	// LoadTx populates a transaction with data from the store.
	//
	// It returns an error only if there is a problem communicating with the
	// store. A non-existent transaction is not an error; rather, mtx.State is
	// set to TxInboundMessageReceived and mtx.Revision is set to 0.
	LoadTx(
		ctx context.Context,
		mtx *MessageTransaction,
	) error

	// SaveTx persists a transaction to the store.
	//
	// It returns an error if there is a problem communicating with the store.
	// ok is false if mtx.Revision does not match the transaction revision in
	// the store.
	SaveTx(
		ctx context.Context,
		mtx *MessageTransaction,
	) (ok bool, err error)
}
