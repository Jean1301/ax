package account

import (
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
)

// Account is a bank account.
type Account struct {
	ax.AggregateBehavior

	// IsOpen is true if the account has been opened.
	IsOpen bool
}

// DoOpen opens a new bank account.
func (a *Account) DoOpen(ctx ax.AggregateCommandContext, m *messages.OpenAccount) error {
	if !a.IsOpen {
		ctx.PublishEvent(&messages.AccountOpened{
			AccountId: m.AccountId,
		})
	}

	return nil
}

// WhenOpened is called when an account is opened.
func (a *Account) WhenOpened(m *messages.AccountOpened) {
	a.IsOpen = true
}

// Apply updates the aggregate state to reflect the occurence of m.
func (a *Account) Apply(m ax.Event) {
	switch e := m.(type) {
	case *messages.AccountOpened:
		a.WhenOpened(e)
	default:
		panic(ax.UnexpectedMessageError{Message: m})
	}
}
