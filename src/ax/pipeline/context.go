package pipeline

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

type messageContext struct {
	context.Context

	inbound  ax.InboundMessage
	outbound []ax.OutboundMessage
}

func (c *messageContext) MessageEnvelope() ax.Envelope {
	return c.env
}

func (c *messageContext) ExecuteCommand(m ax.Command) {
	c.outbound = append(c.outbound, OutboundMessage{
		Operation: OpExecute,
		Envelope:  c.env.New(m),
	})
}

func (c *messageContext) PublishEvent(m ax.Event) {
	c.outbound = append(c.outbound, OutboundMessage{
		Operation: OpPublish,
		Envelope:  c.env.New(m),
	})
}
