package endpoint

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

type Context struct {
	context.Context

	envelope ax.Envelope
	commands []ax.Command
	events   []ax.Event
}

// Envelope returns the message envelope.
func (c *Context) Envelope() ax.Envelope {
	return c.envelope
}

// Execute enqueues a command for execution.
func (c *Context) Execute(m ax.Command) {
	c.commands = append(c.commands, m)
}

// Publish publishes an event.
func (c *Context) Publish(m ax.Event) {
	c.events = append(c.events, m)

}

//
// func (c *HandlerContext) Execute(m ax.Command) error {
// 	mt := ax.TypeOf(m)
//
// 	if ep, ok := c.Routing.EndpointFor(mt); ok {
// 		return c.Transport.Send(
// 			c.Context,
// 			ep,
// 			transport.OutboundMessage{
// 				Parent:  c.Parent,
// 				Message: m,
// 			},
// 		)
// 	}
//
// 	return fmt.Errorf("no route for %s", mt.Name)
// }
//
// func (c *HandlerContext) Publish(m ax.Event) error {
// 	return c.Transport.Publish(
// 		c.Context,
// 		transport.OutboundMessage{
// 			Message: m,
// 			Parent:  c.Parent,
// 		},
// 	)
// }
