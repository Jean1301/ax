package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// Repository is an interface for loading and saving saga instances to/from a
// persistent data store.
type Repository interface {
	// LoadSaga fetches a saga instance from the store based a mapping key
	// for a particular message type.
	//
	// ok is true if the instance is found, in which case si is populated with
	// data from the store.
	//
	// err is non-nil if there is a problem communicating with the store itself.
	LoadSaga(
		ctx context.Context,
		tx persistence.Transaction,
		mt ax.MessageType,
		mk MappingKey,
		si Instance,
	) (ok bool, err error)

	// SaveSaga persists a saga instance and its associated mapping table to the
	// store as part of tx.
	//
	// It returns an error if the saga instance has been modified since it was
	// loaded, or if there is a problem communicating with the store itself.
	//
	// Save() panics if the repository is not able to enlist in tx because it
	// uses a different underlying storage system.
	SaveSaga(
		ctx context.Context,
		tx persistence.Transaction,
		si Instance,
		mt MappingTable,
	) error
}
