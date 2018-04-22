package account

import (
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
)

// Handler handles commands for the Account aggregate.
type Handler struct{}

// CommandTypes returns the set of commands that the aggregate handles.
func (Handler) CommandTypes() ax.MessageTypeSet {
	return ax.TypesOf(
		&messages.OpenAccount{},
	)
}

// MapToAggregate returns the aggregate ID that handles the given message.
func (Handler) MapToAggregate(m ax.Command) (id ax.AggregateID) {
	type messageWithAccountID interface {
		GetAccountId() string
	}

	id.MustParse(
		m.(messageWithAccountID).GetAccountId(),
	)

	return
}

// InitialState returns a new aggregate value.
func (Handler) InitialState() ax.Aggregate {
	return &Account{}
}

// HandleCommand mutates an aggregate based on the given command.
func (Handler) HandleCommand(
	ctx ax.AggregateCommandContext,
	agg ax.Aggregate,
	m ax.Command,
) error {
	acct := agg.(*Account)

	switch c := m.(type) {
	case *messages.OpenAccount:
		return acct.DoOpen(ctx, c)
	default:
		return ax.UnexpectedMessageError{Message: m}
	}
}
