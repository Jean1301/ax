package marshaling

import (
	"fmt"
	"mime"

	"github.com/jmalloc/ax/src/ax"
)

// MarshalMessage marshals m to a binary representation.
func MarshalMessage(m ax.Message) (contentType string, data []byte, err error) {
	return marshalProtobuf(m)
}

// UnmarshalMessage unmarshals a message from a binary representation.
// ct is the MIME content-type for the binary data.
func UnmarshalMessage(contentType string, data []byte) (ax.Message, error) {
	ct, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}

	if ct != protobufContentType {
		return nil, fmt.Errorf(
			"can not unmarshal '%s', content-type is not supported",
			contentType,
		)
	}

	pb, err := unmarshalProtobuf(ct, params, data)
	if err != nil {
		return nil, err
	}

	if m, ok := pb.(ax.Message); ok {
		return m, nil
	}

	return nil, fmt.Errorf(
		"can not unmarshal '%s', content-type is not a message",
		ct,
	)
}
