package axrmq

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/streadway/amqp"
)

const sendExchange = "ax/send"
const publishExchange = "ax/publish"

// queueNames returns the name of the inbox and error queue to use for the
// endpoint named ep.
func queueNames(ep string) (string, string) {
	return ep + "/inbox", ep + "/error"
}

// setupTopology declares all exchanges, queues and bindings for the endpoint
// named ep.
func setupTopology(ch *amqp.Channel, ep string, subscriptions ax.MessageTypeSet) error {
	if err := setupExchanges(ch); err != nil {
		return err
	}

	if err := setupQueues(ch, ep); err != nil {
		return err
	}

	if err := setupBindings(ch, ep, subscriptions); err != nil {
		return err
	}

	return nil
}

func setupExchanges(ch *amqp.Channel) error {
	if err := ch.ExchangeDeclare(
		sendExchange,
		"direct",
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait,
		nil,   // args,
	); err != nil {
		return err
	}

	if err := ch.ExchangeDeclare(
		publishExchange,
		"direct",
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait,
		nil,   // args,
	); err != nil {
		return err
	}

	return nil
}

// setupQueues declares the inbox and error queues for the endpoint named ep.
func setupQueues(ch *amqp.Channel, ep string) error {
	inbox, errors := queueNames(ep)

	if _, err := ch.QueueDeclare(
		inbox,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": errors,
		},
	); err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(
		errors,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	); err != nil {
		return err
	}

	return nil
}

func setupBindings(ch *amqp.Channel, ep string, subscriptions ax.MessageTypeSet) error {
	inbox, _ := queueNames(ep)

	if err := ch.QueueBind(
		inbox,
		ep,
		sendExchange,
		false, // noWait
		nil,   // args
	); err != nil {
		return err
	}

	for _, mt := range subscriptions.Members() {
		if mt.IsEvent() {
			if err := ch.QueueBind(
				inbox,
				mt.Name,
				publishExchange,
				false, // noWait
				nil,   // args
			); err != nil {
				return err
			}
		}
	}

	return nil
}
