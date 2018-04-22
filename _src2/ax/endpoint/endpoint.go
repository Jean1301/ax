package endpoint

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/transport"
	uuid "github.com/satori/go.uuid"
)

type Endpoint struct {
	Name      string
	Transport transport.Transport
	Handlers  []ax.Handler
}

func (ep *Endpoint) Run(ctx context.Context) error {
	var subscriptions ax.MessageTypeSet
	for _, h := range ep.Handlers {
		subscriptions = subscriptions.Union(
			h.MessageTypes(),
		)
	}

	if err := ep.Transport.Initialize(ctx, ep.Name, subscriptions); err != nil {
		return err
	}

	for {
		m, err := ep.Transport.Receive(ctx)
		if err != nil {
			return err
		}

		go ep.process(ctx, m)
	}
}

func (ep *Endpoint) Send(ctx context.Context, m ax.Command) error {
	om := &transport.OutboundMessage{
		Envelope: ax.Envelope{
			MessageID: uuid.NewV4().String(),
			Message:   m,
		},
	}

	return ep.Transport.Send(
		ctx,
		ax.TypeOf(m).Package(),
		om,
	)
}

func (ep *Endpoint) process(ctx context.Context, m *transport.InboundMessage) {
	g, ctx := errgroup.WithContext(ctx)
	mt := ax.TypeOf(m.Envelope.Message)

	for _, h := range ep.Handlers {
		if !h.MessageTypes().Has(mt) {
			continue
		}
	}

	hc := &Context{
		Context:  ctx,
		envelope: m.Envelope,
	}

	err := ep.Handler.Handle(hc, m.Envelope.Message)

	if err != nil {
		fmt.Println(err)
		err = m.Retry(ctx)
	} else {
		err = m.Done(ctx)
	}

	if err != nil {
		panic(err)
	}
}
