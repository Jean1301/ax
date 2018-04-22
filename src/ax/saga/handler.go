package saga

import (
	"github.com/jmalloc/ax/src/ax"
)

// MessageHandler is an implementation of ax.MessageHandler that loads a saga
// instance from a repository and passes control to a saga.
type MessageHandler struct {
	Repository Repository
	Saga       Saga

	triggers     ax.MessageTypeSet
	messageTypes ax.MessageTypeSet
}

// MessageTypes returns the set of messages that the handler can handle.
func (h *MessageHandler) MessageTypes() ax.MessageTypeSet {
	triggers, others := h.Saga.MessageTypes()
	h.triggers = triggers
	h.messageTypes = triggers.Union(others)

	return h.messageTypes
}

// HandleMessage handles a message.
func (h *MessageHandler) HandleMessage(ctx ax.MessageContext, m ax.Message) error {
	mt := ax.TypeOf(m)

	if !h.messageTypes.Has(mt) {
		return ax.UnexpectedMessageError{
			Message: m,
		}
	}

	mk := h.Saga.MapMessage(m)
	si := h.Saga.InitialState()

	ok, err := h.Repository.Load(ctx, mt, mk, si)
	if err != nil {
		return err
	}

	if !ok && !h.triggers.Has(mt) {
		return h.Saga.HandleNotFound(ctx, m)
	}

	// pass the message to the saga for handling.
	if err := h.Saga.HandleMessage(ctx, m, si); err != nil {
		return err
	}

	// save the changes to the saga and its mapping table.
	return h.Repository.Save(
		ctx,
		ctx.MessageTransaction(),
		si,
		buildMappingTable(h.Saga, si),
	)
}
