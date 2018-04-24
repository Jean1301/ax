package ax

type MessageFlow struct {
}

// UnitOfWork describes the process of handling a single message.
//
// The message is first passed to a MessageHandler which performs application
// defined logic and may optionally produce additional messages.
//
// Once the message handler returns successfully, the messages it produced are
// stored in the UnitOfWork's "outbox" until they are dispatched via a transport.
//
// Finally, once the produced messages have been dispatched, the UnitOfWork is
// marked as complete and the original message is acknowledged.
type UnitOfWork struct {
	// MessageID is the ID "causal" message, that is, the message being handled.
	MessageID MessageID

	// Progress is the current "step" that is to be performed.
	Progress UnitOfWorkProgress

	// Outbox is the set of messages produced by the message handler.
	//
	// It is populated when the unit progresses to ProgressHandled,
	// and later cleared when the unit of work progresses to ProgressComplete.
	Outbox []OutboxMessage

	// Revision is the current revision of the unit of work, used to implement
	// optimistic concurrency control.
	Revision uint64
}

// UnitOfWorkProgress is an enumeration of the steps that form a UnitOfWork.
type UnitOfWorkProgress int

func x() {
	o, ok, err := outbox.Load(id)

	if !ok {
		handler.Handle(m)

		outbox.Save(m, outbox)
	}

	for _, m := range outbox {
		send(m)
		outbox.MarkSent(m)
	}

	m.Ack()
}

const (
	// ProgressNone is the state of a unit of work before the MessageHandler has
	// been invoked.
	ProgressNone UnitOfWorkProgress = iota

	// ProgressHandled is the state of a unit of work after the MessageHandler
	// has been invoked and returned successfully, but before the messages it
	// produces have been dispatched.
	ProgressHandled

	// ProgressComplete is the state of a unit of work after the produced
	// messages have been dispatched.
	ProgressComplete
)
