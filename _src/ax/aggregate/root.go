package aggregate

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Root is a domain aggregate root.
//
// Aggregate roots are specialized message handlers that handle "domain commands"
// to produce some change in state within a business domain.
type Root interface {
	// CommandTypes returns the domain command types that the handler root
	// is responsible for executing.
	CommandTypes() ax.MessageTypeSet

	// InitialState returns a new aggregate instance.
	InitialState() ax.Instance

	// Execute mutates an aggregate based on the given command.
	Execute(ExecutionContext, DomainCommand, ax.Instance) error
}

// ExecutionContext provides access to the messaging system within the context
// of handling a particular domain command.
type ExecutionContext interface {
	context.Context

	// Publish publishes an event.
	Publish(ax.Event)
}
