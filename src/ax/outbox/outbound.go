package outbox

import (
	"context"

	"github.com/jmalloc/ax/src/ax/pipeline"
)

// OutboundStage is a pipeline.OutboundStage that keeps a collection of the
// dispatched messages in memory.
type OutboundStage struct {
	Outbox []pipeline.OutboundMessage
}

// DispatchMessage adds m to s.Outbox.
func (s *OutboundStage) DispatchMessage(
	ctx context.Context,
	m pipeline.OutboundMessage,
) error {
	s.Outbox = append(s.Outbox, m)
	return nil
}
