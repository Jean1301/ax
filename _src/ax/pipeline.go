package ax

import "context"

type InboundPipelineStage interface {
	Initialize(context.Context, *Endpoint) error
	Process(context.Context, InboundMessage) error
}

type OutboundPipelineStage interface {
	Initialize(context.Context, *Endpoint) error
	Process(context.Context, OutboundMessage) error
}
