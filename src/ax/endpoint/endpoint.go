package endpoint

import (
	"context"
	"sync"

	"github.com/jmalloc/ax/src/ax/pipeline"
	"github.com/jmalloc/ax/src/ax/transport"
)

// Endpoint is a named endpoint participating in the exchange of messages.
type Endpoint struct {
	Name             string
	Transport        transport.Transport
	InboundPipeline  pipeline.InboundStage
	OutboundPipeline pipeline.OutboundStage

	wg sync.WaitGroup
}

// Run processes messages until ctx is canceled.
func (ep *Endpoint) Run(ctx context.Context) error {
	if err := ep.Transport.Initialize(ctx, ep.Name); err != nil {
		return err
	}

	if ep.InboundPipeline != nil {
		if err := ep.InboundPipeline.Initialize(ctx, ep.Transport); err != nil {
			return err
		}
	}

	if ep.OutboundPipeline != nil {
		if err := ep.OutboundPipeline.Initialize(ctx, ep.Transport); err != nil {
			return err
		}
	}

	var (
		m   transport.InboundMessage
		err error
	)

	for {
		m, err = ep.Transport.Receive(ctx)
		if err != nil {
			break
		}

		ep.wg.Add(1)
		go ep.process(ctx, m)
	}

	ep.wg.Wait()

	return err
}

func (ep *Endpoint) process(ctx context.Context, m transport.InboundMessage) {
	defer ep.wg.Done()

	if err := ep.InboundPipeline.DispatchMessage(ctx, ep.OutboundPipeline, m); err != nil {
		panic(err)
	}
}
