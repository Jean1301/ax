package messages

import (
	"fmt"

	"github.com/jmalloc/ax/src/ax/ident"
)

// IsCommand marks the message as a command.
func (m *OpenAccount) IsCommand() {}

// Description returns a human-readable description of the message.
func (m *OpenAccount) Description() string {
	return fmt.Sprintf("open account %s", ident.FormatID(m.AccountId))
}

// IsEvent marks the message as an event.
func (ev *AccountOpened) IsEvent() {}

// Description returns a human-readable description of the message.
func (ev *AccountOpened) Description() string {
	return fmt.Sprintf("opened account %s", ident.FormatID(ev.AccountId))
}
