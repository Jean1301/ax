package axmysql

import (
	"context"
	"database/sql"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/marshaling"
)

// loadTransaction populates mtx from the database.
// It returns false if no such transaction exists.
func loadTransaction(ctx context.Context, tx *sql.Tx) (ax.MessageTransaction, error) {
	row := tx.QueryRowContext(
		ctx,
		`SELECT
		"state",
		"revision"
		FROM "message_transaction"
		WHERE "causation_id" = ?`,
		mtx.MessageID,
	)

	if err := row.Scan(
		&mtx.State,
		&mtx.Revision,
	); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	rows, err := tx.QueryContext(
		ctx,
		`SELECT
		"message_id",
		"correlation_id",
		"causation_id",
		"time",
		"content_type",
		"body",
		FROM "message_outbox"
		WHERE "causation_id" = ?`,
		mtx.MessageID,
	)
	if err != nil {
		return false, err
	}

	for rows.Next() {
		env, err := unpackMessageEnvelope(rows)
		if err != nil {
			return false, err
		}

		mtx.Outbox = append(mtx.Outbox, env)
	}

	return true, nil
}

// saveTransaction persists mtx to the database.
// It returns false if mtx.Revision is not the current transaction revision.
func saveTransaction(
	ctx context.Context,
	tx *sql.Tx,
	mtx *ax.MessageTransaction,
) (bool, error) {
	var ok bool
	var err error

	if mtx.Revision == 0 {
		ok, err = insertTransaction(ctx, tx, mtx)
	} else {
		ok, err = updateTransaction(ctx, tx, mtx)
	}

	if !ok || err != nil {
		return false, err
	}

	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM "outbox_message"
		WHERE "causation_id" = ?`,
		mtx.MessageID,
	)
	if err != nil {
		return false, err
	}

	for _, env := range mtx.Outbox {
		err = insertOutboxMessage(ctx, tx, env)
		if err != nil {
			return false, err
		}
	}

	mtx.Revision++
	return true, nil
}

// insertTransaction inserts a new transaction to the database.
// It returns false if the transaction already exists.
func insertTransaction(
	ctx context.Context,
	tx *sql.Tx,
	mtx *ax.MessageTransaction,
) (bool, error) {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO "message_transaction" SET
		"causation_id" = ?,
		"revision" = 1,
		"state" = ?`,
		mtx.MessageID,
		mtx.State,
	)

	if err != nil {
		if isDuplicateKey(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// updateTransaction updates an existing transaction in the database. It returns
// false if the transaction does not exist or if mtx.Revision does not match the
// current transaction revision.
func updateTransaction(
	ctx context.Context,
	tx *sql.Tx,
	mtx *ax.MessageTransaction,
) (bool, error) {
	row := tx.QueryRowContext(
		ctx,
		`SELECT "revision"
		FROM "message_transaction"
		WHERE "causation_id" = ?
		FOR UPDATE`,
		mtx.MessageID,
	)

	var rev ax.Revision
	if err := row.Scan(&rev); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
	}

	if rev != mtx.Revision {
		return false, nil
	}

	res, err := tx.ExecContext(
		ctx,
		`UPDATE "message_transaction" SET
		"revision" = ?,
		"state" = ?
		WHERE "causation_id" = ?`,
		mtx.MessageID,
		mtx.Revision+1,
		mtx.State,
		mtx.Revision,
	)
	if err != nil {
		return false, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return n != 0, nil
}

// insertOutboxMessage inserts a message to the outbox.
func insertOutboxMessage(
	ctx context.Context,
	tx *sql.Tx,
	env ax.Envelope,
) error {
	contentType, body, err := marshaling.MarshalMessage(env.Message)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO "outbox_message" SET
		"message_id" = ?,
		"correlation_id" = ?,
		"causation_id" = ?,
		"time" = ?,
		"content_type" = ?,
		"body" = ?`,
		env.MessageID,
		env.CorrelationID,
		env.CausationID,
		env.Time,
		contentType,
		body,
	)

	return err
}
