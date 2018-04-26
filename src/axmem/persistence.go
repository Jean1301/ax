package axmem

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/jmalloc/ax/src/ax/persistence"
)

// DataStore is an implementation of persistence.DataStore that stores data in
// memory.
type DataStore struct{}

// BeginTx starts a new transaction.
func (ds *DataStore) BeginTx(ctx context.Context) (persistence.Tx, error) {
	return &Tx{}, nil
}

// Tx represents an atomic unit-of-work performed on a DataStore.
type Tx struct {
	m       sync.Mutex
	isEnded bool
	actions []func(bool)
	locks   map[sync.Locker]struct{}
}

// Commit applies the changes to the data store.
func (tx *Tx) Commit() error {
	return tx.end(true)
}

// Rollback discards the changes without applying them to the data store.
func (tx *Tx) Rollback() error {
	return tx.end(false)
}

// UnderlyingTx returns a data-store-specific value that represents the
// transaction. For example, for an SQL-based data store this may be the
// *sql.Tx.
func (tx *Tx) UnderlyingTx() interface{} {
	return tx
}

// Enlist registers a function to be called when this transaction is committed
// or rolled back.
//
// lock is acquired immediately (if it has not already been enlisted in this
// transaction) and released after all enlisted functions have been called.
func (tx *Tx) Enlist(
	lock sync.Locker,
	fn func(bool),
) error {
	tx.m.Lock()
	defer tx.m.Unlock()

	if tx.isEnded {
		return errors.New("transaction has already ended")
	}

	if tx.locks == nil {
		tx.locks = map[sync.Locker]struct{}{}
	}

	if _, ok := tx.locks[lock]; !ok {
		lock.Lock()
		tx.locks[lock] = struct{}{}
	}

	tx.actions = append(tx.actions, fn)

	return nil
}

func (tx *Tx) end(commit bool) error {
	tx.m.Lock()
	defer tx.m.Unlock()

	if tx.isEnded {
		return errors.New("transaction has already ended")
	}

	defer func() {
		tx.isEnded = true
		tx.actions = nil

		for l := range tx.locks {
			l.Unlock()
		}
	}()

	for _, fn := range tx.actions {
		fn(commit)
	}

	return nil
}

func unwrapTx(tx persistence.Tx) (*Tx, error) {
	v := tx.UnderlyingTx()

	memtx, ok := v.(*Tx)
	if ok {
		return memtx, nil
	}

	return nil, fmt.Errorf(
		"can not enlist in %s, expected underlying transaction to be %s",
		reflect.TypeOf(tx),
		reflect.TypeOf(memtx),
	)
}
