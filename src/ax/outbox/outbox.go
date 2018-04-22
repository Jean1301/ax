package outbox

import "github.com/jmalloc/ax/src/ax"

type OperationType int

const (
	OpExecute OperationType = iota
	OpPublish
)

type Outbox struct {
	Operations []Operation
}

type Operation struct {
	Op      OperationType
	Message ax.Message
}
