package ax

import "time"

// Envelope is a container for a message and its meta-data.
type Envelope struct {
	MessageID     MessageID
	CorrelationID MessageID
	CausationID   MessageID
	Time          time.Time
	Message       Message
}

// NewCause returns a new message envelope that contains m.
func NewCause(m Message) Envelope {
	env := Envelope{
		Time:    time.Now(),
		Message: m,
	}

	env.MessageID.GenerateUUID()
	env.CausationID = env.MessageID
	env.CorrelationID = env.MessageID

	return env
}

// NewEffect returns a new message envelope that contains m.
// m is configured as an "effect" of the "causal" message e.Message.
func (e Envelope) NewEffect(m Message) Envelope {
	env := Envelope{
		CorrelationID: e.CorrelationID,
		CausationID:   e.MessageID,
		Time:          time.Now(),
		Message:       m,
	}

	env.MessageID.GenerateUUID()

	return env
}
