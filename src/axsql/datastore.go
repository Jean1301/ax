package axsql

import (
	"context"
	"database/sql"

	"github.com/jmalloc/ax/src/ax/persistence"
)

// DataStore is an implementation of persistence.DataStore that persists data in
// an SQL database.
type DataStore struct {
	Dialect Dialect
	DB      *sql.DB
}

// BeginTx starts a new transaction.
func (ds *DataStore) BeginTx(ctx context.Context) (persistence.Tx, persistence.Committer, error) {
	tx, err := ds.DB.BeginTx(ctx, ds.Dialect.TxOptions())
	if err != nil {
		return nil, nil, err
	}

	return &Tx{ds, tx}, tx, nil
}
