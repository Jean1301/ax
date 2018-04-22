package ax

import "context"

type Endpoint struct {
	Name             string
	Transport        Transport
	InboundPipeline  InboundPipelineStage
	OutboundPipeline OutboundPipelineStage
}

func (ep *Endpoint) Run(ctx context.Context) error {
	if err := ep.Transport.Initialize(ctx, ep); err != nil {
		return err
	}

	if ep.InboundPipeline != nil {
		if err := ep.InboundPipeline.Initialize(ctx, ep); err != nil {
			return err
		}
	}

	if ep.OutboundPipeline != nil {
		if err := ep.OutboundPipeline.Initialize(ctx, ep); err != nil {
			return err
		}
	}

	for {
		m, err := ep.Transport.Receive(ctx)
		if err != nil {
			return err
		}

		go ep.process(ctx, m)
	}
}

func (ep *Endpoint) process(ctx context.Context, m InboundMessage) {
	if err := ep.InboundPipeline.Process(ctx, m); err != nil {
		panic(err)
	}
}
