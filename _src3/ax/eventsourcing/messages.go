package eventsourcing

import (
	"fmt"

	"github.com/jmalloc/ax/src/ax/ident"
)

// IsEvent marks the message as an event.
func (m *CommandHandled) IsEvent() {}

// Description returns a human-readable description of the message.
func (m *CommandHandled) Description() string {
	return fmt.Sprintf("command %s handled", ident.FormatID(m.CommandId))
}

// IsEvent marks the message as an event.
func (m *CommandEventsPublished) IsEvent() {}

// Description returns a human-readable description of the message.
func (m *CommandEventsPublished) Description() string {
	return fmt.Sprintf("events for command %s published", ident.FormatID(m.CommandId))
}
