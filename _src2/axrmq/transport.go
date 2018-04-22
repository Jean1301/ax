package axrmq

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/transport"
	"github.com/streadway/amqp"
)

// DefaultConcurrency is the default number of messages to handle concurrently.
var DefaultConcurrency = runtime.NumCPU() * 2

// Transport is an implementation of transport.Transport that communicates using
// a RabbitMQ AMQP broker.
type Transport struct {
	Conn        *amqp.Connection
	Concurrency int

	endpoint  string
	publisher *Publisher

	channel  *amqp.Channel
	messages <-chan amqp.Delivery
	closed   chan *amqp.Error
}

// Initialize sets up the transport to operate for a particular endpoint.
func (t *Transport) Initialize(ctx context.Context, ep string, subscriptions ax.MessageTypeSet) error {
	ch, err := t.Conn.Channel()
	if err != nil {
		return err
	}
	defer func() {
		// close the channel if it has not been "captured" by the transport
		// for continued use.
		if t.channel != ch {
			_ = ch.Close()
		}
	}()

	err = setupTopology(ch, ep, subscriptions)
	if err != nil {
		return err
	}

	concurrency := t.Concurrency
	if concurrency == 0 {
		concurrency = DefaultConcurrency
	}

	err = ch.Qos(concurrency, 0, false)
	if err != nil {
		return err
	}

	t.publisher = NewPublisher(t.Conn, concurrency)

	inbox, _ := queueNames(ep)

	messages, err := ch.Consume(
		inbox,
		"inbox", // consumer tag
		false,   // autoAck
		false,   // exclusive
		false,   // noLocal
		false,   // noWait
		nil,     // args
	)
	if err != nil {
		return err
	}

	t.closed = make(chan *amqp.Error)
	ch.NotifyClose(t.closed)

	t.endpoint = ep
	t.messages = messages
	t.channel = ch

	return nil
}

// Close stops and closes the transport.
func (t *Transport) Close() error {
	if t.channel == nil {
		return nil
	}

	return t.channel.Close()
}

// Receive returns the next message from the transport.
// It blocks until a message is availble, or ctx is canceled.
func (t *Transport) Receive(ctx context.Context) (*transport.InboundMessage, error) {
	for {
		select {
		case del := <-t.messages:
			env, err := t.receive(ctx, del)
			if env != nil || err != nil {
				return env, err
			}
		case err := <-t.closed:
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (t *Transport) receive(ctx context.Context, del amqp.Delivery) (*transport.InboundMessage, error) {
	m := transport.InboundMessage{
		Endpoint: del.ReplyTo,
		Done: func(context.Context) error {
			return del.Ack(false) // false = single message
		},
		Retry: func(context.Context) error {
			if del.Redelivered {
				time.Sleep(1 * time.Second)
			}
			return del.Reject(true) // false = don't requeue
		},
	}

	if err := unmarshalMessage(del, &m.Envelope); err != nil {
		// TODO: log/sentry/etc

		if err := del.Reject(false); err != nil {
			return nil, err
		}

		return nil, nil
	}

	return &m, nil
}

// Send sends a message to a specific endpoint.
func (t *Transport) Send(ctx context.Context, ep string, m *transport.OutboundMessage) error {
	var pub amqp.Publishing

	if err := marshalMessage(m.Envelope, &pub); err != nil {
		fmt.Println(err)
		return err
	}

	return t.publisher.Publish(
		ctx,
		sendExchange,
		ep,
		true,
		pub,
	)
}

// Publish multicasts a message to any endpoints that are interested in
// receiving it.
func (t *Transport) Publish(ctx context.Context, m *transport.OutboundMessage) error {
	var pub amqp.Publishing

	if err := marshalMessage(m.Envelope, &pub); err != nil {
		return err
	}

	return t.publisher.Publish(
		ctx,
		publishExchange,
		pub.Type,
		false,
		pub,
	)
}
