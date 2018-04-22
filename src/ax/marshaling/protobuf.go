package marshaling

import (
	"errors"
	"fmt"
	"mime"
	"reflect"

	"github.com/golang/protobuf/proto"
)

const (
	protobufContentType = "application/vnd.google.protobuf"
)

func marshalProtobuf(m proto.Message) (string, []byte, error) {
	n := proto.MessageName(m)
	if n == "" {
		return "", nil, fmt.Errorf(
			"can not marshal '%s', protocol is not registered",
			reflect.TypeOf(m),
		)
	}

	contentType := mime.FormatMediaType(
		protobufContentType,
		map[string]string{"proto": n},
	)

	data, err := proto.Marshal(m)
	return contentType, data, err
}

func unmarshalProtobuf(
	contentType string,
	params map[string]string,
	data []byte,
) (proto.Message, error) {
	n, ok := params["proto"]
	if !ok {
		return nil, errors.New(
			"can not unmarshal message, protocol is not specified",
		)
	}

	t := proto.MessageType(n)
	if t == nil {
		return nil, fmt.Errorf(
			"can not unmarshal '%s', protocol is not registered",
			n,
		)
	}

	m := reflect.New(
		t.Elem(),
	).Interface().(proto.Message)

	return m, proto.Unmarshal(data, m)
}
