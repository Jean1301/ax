package outbox

import (
	"github.com/jmalloc/ax/src/ax"
)

type MessageContext struct {
	ax.MessageContext
	Outbox *Outbox
}

// ExecuteCommand enqueues a command to be executed.
func (c MessageContext) ExecuteCommand(m ax.Command) error {
	c.Outbox.Operations = append(c.Outbox.Operations, Operation{
		Op:      OpExecute,
		Message: m,
	})

	return nil
}

// PublishEvent enqueues events to be published.
func (c MessageContext) PublishEvent(m ax.Event) error {
	c.Outbox.Operations = append(c.Outbox.Operations, Operation{
		Op:      OpPublish,
		Message: m,
	})

	return nil
}
