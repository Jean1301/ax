package messages

import (
	"fmt"
)

// Description returns a human-readable description of the event.
func (ev *TransactionStarted) Description() string {
	switch d := ev.Details.(type) {
	case *TransactionStarted_Transfer:
		return fmt.Sprintf(
			"started funds transfer %s of %d from %s to %s",
			ev.TransactionId,
			ev.Cents,
			d.Transfer.FromAccountId,
			d.Transfer.ToAccountId,
		)
	}

	return fmt.Sprintf(
		"started unrecognised transaction %s for the amount of %d",
		ev.TransactionId,
		ev.Cents,
	)
}
