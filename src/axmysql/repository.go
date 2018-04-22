package axmysql

import (
	"database/sql"
	"errors"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/eventsourcing"
)

// EventSourcedRepository is an implementation of ax.Repository that stores
// event-sourced aggregates in a MySQL database.
type EventSourcedRepository struct {
	DB *sql.DB
}

// LoadAggregate populates agg with data fetched from the store and sets its
// revision.
//
// It returns an error if there is a problem communicating with the store.
//
// It is NOT an error if the aggregate does not exist, the aggregate will
// simply remain at revision zero.
func (r *EventSourcedRepository) LoadAggregate(
	ctx ax.AggregateCommandContext,
	agg ax.Aggregate,
) error {
	a := agg.(eventsourcing.Aggregate)
	mtx := ctx.MessageTransaction()

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = loadTransaction(ctx, tx, mtx)
	if err != nil {
		return err
	}

	rows, err := tx.QueryContext(
		ctx,
		`SELECT
        "message_id",
        "correlation_id",
        "causation_id",
        "time",
        "content_type",
        "body"
        FROM "aggregate_event"
        WHERE "aggregate_id" = ?
        AND "revision" > ?
        ORDER BY "revision"`,
		agg.AggregateID(),
		agg.Revision(),
	)

	if err != nil {
		return err
	}

	for rows.Next() {
		select {
		case <-ctx.Done():
			// TODO: is this neccessary? need some definitive information about
			// whether the context passed to QueryContext() is checked when reading
			// the rows.
			return ctx.Err()
		default:
		}

		env, err := unpackMessageEnvelope(rows)
		if err != nil {
			return err
		}

		m, ok := env.Message.(ax.Event)
		if !ok {
			return errors.New("unexpected message, not an event") // TODO
		}

		a.Apply(m)
	}

	return nil
}

// SaveAggregate persists an aggregate to the store.
//
// It returns an error if there is a problem communicating with the store.
// ok is true if the changes are persisted successfully, or false if the
// change was rejected because agg.Revision() is out of date.
func (r *EventSourcedRepository) SaveAggregate(
	ctx ax.AggregateCommandContext,
	agg ax.Aggregate,
) (ok bool, err error) {
}

func (r *EventSourcedRepository) MarkComplete(ctx ax.AggregateCommandContext) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	markTransactionComplete(
		ctx,
		tx,
		ctx.MessageTransaction(),
	)
}
