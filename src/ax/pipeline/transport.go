package pipeline

import (
	"context"

	"github.com/jmalloc/ax/src/ax/transport"
)

type X struct {
	Transport transport.Transport
}

func (x *X) DispatchMessage(context.Context, transport.OutboundMessage) error {
}
