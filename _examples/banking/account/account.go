package account

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/dogma/example/domain/events"
	"github.com/jmalloc/dogma/example/domain/types"
	"github.com/jmalloc/dogma/src/dogma"
)

// Account is a bank account.
type Account struct {
	ax.AggregateBehavior

	IsOpen  bool
	Balance types.Amount
}

// Apply updates the aggregate's state to reflect the occurrence of ev.
func (a *Account) Apply(ev ax.DomainEvent) {
	switch e := ev.(type) {
	case *events.AccountOpened:
		a.IsOpen = true
	case *events.AccountCredited:
		a.Balance += types.Amount(e.Cents)
	case *events.AccountDebited:
		a.Balance -= types.Amount(e.Cents)
	}
}

func (a *Account) doOpen(ctx ax.AggregateContext, cmd *Open) error {
	if !a.IsOpen {
		ctx.Record(&events.AccountOpened{})
	}

	return nil
}

func (a *Account) doDebit(ctx ax.AggregateContext, cmd *Debit) error {
	if err := a.checkIsOpen(); err != nil {
		return err
	}

	if a.Balance >= cmd.Amount {
		ctx.Record(&events.AccountDebited{
			TransactionId: cmd.TransactionID.Get(),
			Cents:         cmd.Amount.ToCents(),
		})
	} else {
		ctx.Record(&events.DebitRejectedDueToInsufficientFunds{
			TransactionId: cmd.TransactionID.Get(),
			Cents:         cmd.Amount.ToCents(),
		})
	}

	return nil
}

func (a *Account) doCredit(ctx ax.AggregateContext, cmd *Credit) error {
	if err := a.checkIsOpen(); err != nil {
		return err
	}

	ctx.Record(&events.AccountCredited{
		TransactionId: cmd.TransactionID.Get(),
		Cents:         cmd.Amount.ToCents(),
	})

	return nil
}

func (a *Account) checkIsOpen() error {
	if a.IsOpen {
		return nil
	}

	return dogma.InvariantViolated("account has not been opened")
}

// Root is the root of the Account aggregate.
type Root struct{}

// InitialState returns a new aggregate value.
func (Root) InitialState() dogma.Aggregate {
	return &Account{}
}

// Executes returns the set of command types executed by this aggregate root.
func (Root) Executes() []dogma.CommandType {
	return dogma.CommandsT(
		&Open{},
		&Credit{},
		&Debit{},
	)
}

// Produces returns the set of event types produced by this aggregate root.
func (Root) Produces() []dogma.EventType {
	return dogma.EventsT(
		&events.AccountOpened{},
		&events.AccountCredited{},
		&events.AccountDebited{},
		&events.DebitRejectedDueToInsufficientFunds{},
	)
}

// Execute mutates an aggregate based on the given command.
func (Root) Execute(ctx ax.AggregateContext, env dogma.CommandEnvelope) error {
	a := ctx.Aggregate().(*Account)

	switch cmd := env.Command.(type) {
	case *Open:
		return a.doOpen(ctx, cmd)
	case *Debit:
		return a.doDebit(ctx, cmd)
	case *Credit:
		return a.doCredit(ctx, cmd)
	}

	return dogma.ErrCommandUnsupported
}
