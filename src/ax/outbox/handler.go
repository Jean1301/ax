package outbox

import "github.com/jmalloc/ax/src/ax"

type MessageHandler struct {
	Repository Repository
	Next       ax.MessageHandler
}

// MessageTypes returns the set of messages that the handler can handle.
func (h *MessageHandler) MessageTypes() ax.MessageTypeSet {
	return h.Next.MessageTypes()
}

// HandleMessage handles a message.
func (h *MessageHandler) HandleMessage(ctx ax.MessageContext, m ax.Message) error {
	outbox, err := h.Repository.LoadOutbox(
		ctx,
		ctx.Transaction(),
		ctx.MessageEnvelope().MessageID,
	)
	if err != nil {
		return err
	}

	if err := h.Next.HandleMessage(
		&MessageContext{
			MessageContext: ctx,
			Outbox:         &outbox,
		},
		m,
	); err != nil {
		return err
	}
}
