package axrmq

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/marshaling"
	"github.com/streadway/amqp"
)

func marshalMessage(env ax.Envelope, pub *amqp.Publishing) error {
	mt := ax.TypeOf(env.Message)
	pub.MessageId = env.MessageID.Get()
	pub.Type = mt.Name

	var err error
	pub.ContentType, pub.Body, err = marshaling.MarshalMessage(env.Message)

	return err
}

func unmarshalMessage(del amqp.Delivery, env *ax.Envelope) error {
	err := env.MessageID.Parse(del.MessageId)
	if err != nil {
		return err
	}

	env.Message, err = marshaling.UnmarshalMessage(del.ContentType, del.Body)

	return err
}
