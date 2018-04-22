package aggregate

import (
	"reflect"

	"github.com/jmalloc/ax/src/ax"
)

// DomainCommand is a Command that is routed to an aggregate to request some
// change of state within the domain model.
type DomainCommand interface {
	ax.Command

	// AggregateID returns the ID of the aggregate that the command is routed
	// to. The command is said to "target" this aggregate.
	AggregateID() ax.ID
}

var domainCommandType = reflect.TypeOf((*DomainCommand)(nil)).Elem()
