package pipeline

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/src/ax"
)

type UnitOfWorkManager struct {
	Persistence ax.Persistence
	Handlers    []ax.MessageHandler
}

func (u *UnitOfWorkManager) Do(ctx context.Context, m ax.Envelope) error {
	tx, err := u.Persistence.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	uow, err := tx.LoadUnitOfWork(m.MessageID)
	if err != nil {
		return err
	}

	if uow.Progress == ax.ProgressNone {
		err = u.handleMessage(ctx, tx, uow, m)
		if err != nil {
			return err
		}

		tx, err = u.Persistence.Tx(ctx)
		if err != nil {
			return err
		}
	}

	if uow.Progress == ax.ProgressHandled {
		defer tx.Rollback()

		if err := u.dispatchOutbox(ctx, tx, uow); err != nil {
			return err
		}
	}

	return nil
}

func (u *UnitOfWorkManager) handleMessage(
	ctx context.Context,
	tx ax.Transaction,
	uow ax.UnitOfWork,
	m ax.Envelope,
) error {
	mc := &messageContext{
		Context: ctx,
		tx:      tx,
		env:     m,
		uow:     &uow,
	}

	for _, h := range u.Handlers {
		if err := h.HandleMessage(mc, m.Message); err != nil {
			return err
		}
	}

	uow.Progress = ax.ProgressHandled

	if err := tx.SaveUnitOfWork(uow); err != nil {
		return err
	}

	return tx.Commit()

}

func (u *UnitOfWorkManager) dispatchOutbox(
	ctx context.Context,
	tx ax.Transaction,
	uow ax.UnitOfWork,
) error {
	for _, om := range uow.Outbox {
		fmt.Println("dispatch via transport", om)
	}

	uow.Progress = ax.ProgressComplete
	uow.Outbox = nil

	if err := tx.SaveUnitOfWork(uow); err != nil {
		return err
	}

	return tx.Commit()
}
