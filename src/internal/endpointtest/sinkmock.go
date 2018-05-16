// Code generated by moq; DO NOT EDIT
// github.com/matryer/moq

package endpointtest

import (
	"context"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"sync"
)

var (
	lockMessageSinkMockAccept sync.RWMutex
)

// MessageSinkMock is a mock implementation of MessageSink.
//
//     func TestSomethingThatUsesMessageSink(t *testing.T) {
//
//         // make and configure a mocked MessageSink
//         mockedMessageSink := &MessageSinkMock{
//             AcceptFunc: func(ctx context.Context, env endpoint.OutboundEnvelope) error {
// 	               panic("TODO: mock out the Accept method")
//             },
//         }
//
//         // TODO: use mockedMessageSink in code that requires MessageSink
//         //       and then make assertions.
//
//     }
type MessageSinkMock struct {
	// AcceptFunc mocks the Accept method.
	AcceptFunc func(ctx context.Context, env endpoint.OutboundEnvelope) error

	// calls tracks calls to the methods.
	calls struct {
		// Accept holds details about calls to the Accept method.
		Accept []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Env is the env argument value.
			Env endpoint.OutboundEnvelope
		}
	}
}

// Accept calls AcceptFunc.
func (mock *MessageSinkMock) Accept(ctx context.Context, env endpoint.OutboundEnvelope) error {
	if mock.AcceptFunc == nil {
		panic("moq: MessageSinkMock.AcceptFunc is nil but MessageSink.Accept was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Env endpoint.OutboundEnvelope
	}{
		Ctx: ctx,
		Env: env,
	}
	lockMessageSinkMockAccept.Lock()
	mock.calls.Accept = append(mock.calls.Accept, callInfo)
	lockMessageSinkMockAccept.Unlock()
	return mock.AcceptFunc(ctx, env)
}

// AcceptCalls gets all the calls that were made to Accept.
// Check the length with:
//     len(mockedMessageSink.AcceptCalls())
func (mock *MessageSinkMock) AcceptCalls() []struct {
	Ctx context.Context
	Env endpoint.OutboundEnvelope
} {
	var calls []struct {
		Ctx context.Context
		Env endpoint.OutboundEnvelope
	}
	lockMessageSinkMockAccept.RLock()
	calls = mock.calls.Accept
	lockMessageSinkMockAccept.RUnlock()
	return calls
}