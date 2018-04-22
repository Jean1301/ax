package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/axrmq"
	"github.com/streadway/amqp"
)

type Handler struct {
}

func (h *Handler) MessageTypes() ax.MessageTypeSet {
	return ax.TypesOf(
		&messages.OpenAccount{},
	)
}

func (h *Handler) Handle(ctx ax.HandlerContext, m ax.Message) error {
	spew.Dump(m)
	return errors.New("oh shit")
}

func main() {
	broker, err := amqp.Dial("amqp://localhost")
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	transport := axrmq.Transport{
		Conn: broker,
	}
	defer transport.Close()

	ep := endpoint.Endpoint{
		Name:      "ax.examples.banking",
		Transport: &transport,
		Handler:   &Handler{},
	}

	ctx := context.Background()

	go do(&ep)

	if err := ep.Run(ctx); err != nil {
		panic(err)
	}
}

func do(ep *endpoint.Endpoint) {
	time.Sleep(1 * time.Second)

	fmt.Println("doing")

	ctx := context.Background()

	m := &messages.OpenAccount{
		AccountId: "billy-bob",
	}

	if err := ep.Send(ctx, m); err != nil {
		fmt.Println(err)
	}
}
