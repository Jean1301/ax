package axrmq

import (
	"context"
	"fmt"
	"runtime"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
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
func (t *Transport) Initialize(ctx context.Context, ep string) error {
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

	err = setupTopology(ch, ep)
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

// Subscribe configures the transport to listen to messages of the given types.
func (t *Transport) Subscribe(ctx context.Context, mt ax.MessageTypeSet) error {
	return setupSubscriptionBindings(t.channel, t.endpoint, mt)
}

// Close stops and closes the transport.
func (t *Transport) Close() error {
	if t.channel == nil {
		return nil
	}

	return t.channel.Close()
}

// Receive returns the next message from the transport.
// It blocks until a message is available, or ctx is canceled.
func (t *Transport) Receive(ctx context.Context) (bus.InboundMessage, error) {
	for {
		select {
		case del := <-t.messages:
			m, ok, err := t.receive(ctx, del)
			if ok || err != nil {
				return m, err
			}
		case err := <-t.closed:
			return bus.InboundMessage{}, err
		case <-ctx.Done():
			return bus.InboundMessage{}, ctx.Err()
		}
	}
}

func (t *Transport) receive(
	ctx context.Context,
	del amqp.Delivery,
) (bus.InboundMessage, bool, error) {
	m := bus.InboundMessage{
		Done: func(_ context.Context, op bus.InboundOperation) error {
			switch op {
			case bus.OpAck:
				return del.Ack(false) // false = single message
			case bus.OpRetry:
				return del.Reject(true) // true = requeue
			case bus.OpReject:
				return del.Reject(false) // false = don't requeue
			default:
				panic(fmt.Sprintf("unrecognized inbound operation: %d", op))
			}
		},
	}

	if err := unmarshalMessage(del, &m.Envelope); err != nil {
		// TODO: sentry, etc?
		return bus.InboundMessage{}, false, del.Reject(false)
	}

	return m, true, nil
}

// Send sends a message to a specific endpoint.
func (t *Transport) Send(ctx context.Context, m bus.OutboundMessage) error {
	var pub amqp.Publishing

	if err := marshalMessage(m.Envelope, &pub); err != nil {
		fmt.Println(err)
		return err
	}

	switch m.Operation {
	case bus.OpSendUnicast:
		return t.sendUnicast(ctx, m.UnicastEndpoint, pub)
	case bus.OpSendMulticast:
		return t.sendMulticast(ctx, pub)
	default:
		panic(fmt.Sprintf("unrecognized outbound operation: %d", m.Operation))
	}
}

func (t *Transport) sendUnicast(
	ctx context.Context,
	ep string,
	pub amqp.Publishing,
) error {
	return t.publisher.Publish(
		ctx,
		unicastExchange,
		ep,
		true,
		pub,
	)
}

func (t *Transport) sendMulticast(
	ctx context.Context,
	pub amqp.Publishing,
) error {
	return t.publisher.Publish(
		ctx,
		multicastExchange,
		pub.Type,
		false,
		pub,
	)
}
