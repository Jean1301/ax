package main

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/ax/bus/router"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/outbox"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axmem"
	"github.com/jmalloc/ax/src/axrmq"
	"github.com/streadway/amqp"
)

type handler struct {
}

func (h *handler) MessageTypes() ax.MessageTypeSet {
	return ax.TypesOf(
		&messages.OpenAccount{},
		&messages.AccountOpened{},
	)
}

func (h *handler) HandleMessage(ctx ax.MessageContext, m ax.Message) error {
	spew.Dump(m)
	return nil
	// return errors.New("oh shit")
}

func main() {
	routes, err := router.NewRoutingTable([]ax.MessageHandler{
		&handler{},
	})
	if err != nil {
		panic(err)
	}

	broker, err := amqp.Dial("amqp://localhost")
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	transport := &axrmq.Transport{
		Conn: broker,
	}
	defer transport.Close()

	ep := endpoint.Endpoint{
		Name:      "ax.examples.banking",
		Transport: transport,
		InboundPipeline: &persistence.Injector{
			DataStore: &axmem.DataStore{},
			Next: &outbox.InboundStage{
				Repository: &axmem.OutboxRepository{},
				Next: &router.Router{
					Routes: routes,
				},
			},
		},
	}

	go do(&ep)

	ctx := context.Background()
	if err := ep.Run(ctx); err != nil {
		panic(err)
	}
}

func do(ep *endpoint.Endpoint) {
	time.Sleep(1 * time.Second)

	fmt.Println("doing")

	ctx := context.Background()

	var env ax.Envelope
	env.MessageID.GenerateUUID()
	env.Message = &messages.OpenAccount{
		AccountId: "billy-bob",
	}

	o := bus.OutboundMessage{
		Operation:       bus.OpSendUnicast,
		UnicastEndpoint: "ax.examples.banking",
		Envelope:        env,
	}

	if err := ep.Transport.Send(ctx, o); err != nil {
		fmt.Println(err)
	}
}
