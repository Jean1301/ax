package axrmq

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax"
	"github.com/streadway/amqp"
)

func marshalMessage(env ax.Envelope, pub *amqp.Publishing) error {
	buf, err := proto.Marshal(env.Message)
	if err != nil {
		return err
	}

	mt := ax.TypeOf(env.Message)
	pub.MessageId = env.MessageID
	pub.Type = mt.Name
	pub.Body = buf
	pub.Headers = marshalHeaders(env.Headers)

	return nil
}

func unmarshalMessage(del amqp.Delivery, env *ax.Envelope) error {
	mt, ok := ax.TypeByName(del.Type)
	if !ok {
		return fmt.Errorf("unrecognised message type: %s", mt.Name)
	}

	env.Message = mt.New()

	err := proto.Unmarshal(del.Body, env.Message)
	if err != nil {
		return err
	}

	env.MessageID = del.MessageId
	env.Headers, err = unmarshalHeaders(del.Headers)

	return err
}

// marshal headers converts a string map into an AMQP table.
func marshalHeaders(h map[string]string) amqp.Table {
	n := len(h)

	if len(h) == 0 {
		return nil
	}

	t := make(amqp.Table, n)
	for k, v := range h {
		t[k] = v
	}

	return t
}

// unmarshalHeaders converts an AMQP table into a string map.
func unmarshalHeaders(t amqp.Table) (map[string]string, error) {
	if len(t) == 0 {
		return nil, nil
	}

	h := make(map[string]string, len(t))

	for k, v := range t {
		if s, ok := v.(string); ok {
			h[k] = s
		} else {
			return nil, fmt.Errorf("%s header is not a string", k)
		}
	}

	return h, nil
}
