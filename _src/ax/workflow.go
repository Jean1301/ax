package ax

import (
	"github.com/jmalloc/ax/src/ax/internal/ident"
)

// WorkflowID is a unique identifier for a workflow instance.
type WorkflowID struct{ ident.ID }

// WorkflowInstance is the data associated with an active workflow.
type WorkflowInstance interface {
	// WorkflowID returns the workflow's unique identifier.
	// It panics if the workflow ID has not been set.
	WorkflowID() WorkflowID

	// MappingKey returns the "mapping key" used to correlate incoming messages
	// with the workflow instance.
	//
	// For each incoming message m, Workflow.MappingKey(m) is compared to
	// WorkflowInstance.MappingKey(TypeOf(m)). If the two values are equal,
	// Workflow.HandleMessageForInstance() is called with the instance and
	// message.
	MappingKey(MessageType) string

	// Apply updates the workflow state to reflect the receipt of a message.
	Apply(Message)
}

// Workflow is a stateful message handler used to coordinate "long-running"
// processes.
type Workflow interface {
	// Name returns a unique name for the workflow.
	Name() string

	// NewInstance returns a new workflow instance.
	NewInstance() WorkflowInstance

	// MessageTypes returns message types that the workflow handles.
	//
	// It returns two sets. tr is the set of message types that trigger a new
	// workflow instance. pr is the set of message types that 'progress'
	// existing workflows. A given message type may be a member of both sets.
	MessageTypes() (tr MessageTypeSet, pr MessageTypeSet)

	// MappingKey returns the "mapping key" used to correlate incoming messages
	// with workflow instances.
	//
	// For each incoming message m, Workflow.MappingKey(m) is compared to
	// WorkflowInstance.MappingKey(TypeOf(m)). If the two values are equal,
	// Workflow.HandleMessageForInstance() is called with the instance and
	// message.
	MappingKey(Message) string

	// Handle is invoked for each message that is routed to a workflow instance.
	Handle(HandlerContext, Message, WorkflowInstance) error

	// HandleNotFound is invoked when a message was unable to be routed to any
	// workflow instances.
	HandleNotFound(HandlerContext, Message) error
}

type WorkflowStore interface {
	Load(wf string, mt MessageType, k string, i WorkflowInstance) error
}

// WorkflowCoordinator is a handler that forwards messages to workflows.
type WorkflowCoordinator struct {
	Workflows []Workflow
	Store     WorkflowStore
}

// MessageTypes returns the message types that the handler listens to.
func (c *WorkflowCoordinator) MessageTypes() MessageTypeSet {
	var types MessageTypeSet

	for _, wf := range c.Workflows {
		tr, pr := wf.MessageTypes()
		types = types.Union(tr).Union(pr)
	}

	return types
}

// Handle handles a message.
func (c *WorkflowCoordinator) Handle(ctx HandlerContext, m Message) error {
	for _, wf := range c.Workflows {
		if err := c.handle(ctx, wf, m); err != nil {
			return err
		}
	}

	return nil
}

func (c *WorkflowCoordinator) handle(ctx HandlerContext, wf Workflow, m Message) error {
	mt := TypeOf(m)
	tr, pr := wf.MessageTypes()

	i := wf.NewInstance()

	if pr.Has(mt) {
		if err := c.Store.Load(
			wf.Name(),
			mt,
			wf.MappingKey(m),
			i,
		); err != nil {
			return err
		}

		if err := wf.Handle(ctx, m, i); err != nil {
			return err
		}
	} else if tr.Has(mt) {
		if err := wf.Handle(ctx, m, i); err != nil {
			return err
		}
	}

	return wf.HandleNotFound(ctx, m)
}
