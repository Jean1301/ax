package axmysql

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/marshaling"
)

type scanner interface {
	Scan(dst ...interface{}) error
}

func unpackMessageEnvelope(sc scanner) (ax.Envelope, error) {
	var env ax.Envelope
	var contentType string
	var body []byte

	err := sc.Scan(
		&env.MessageID,
		&env.CorrelationID,
		&env.CausationID,
		&env.Time,
		&contentType,
		&body,
	)

	if err == nil {
		env.Message, err = marshaling.UnmarshalMessage(contentType, body)
	}

	return env, err
}
