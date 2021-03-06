// Code generated by moq; DO NOT EDIT
// github.com/matryer/moq

package routingtest

import (
	"context"
	"github.com/jmalloc/ax/src/ax"
	"sync"
)

var (
	lockMessageHandlerMockHandleMessage sync.RWMutex
	lockMessageHandlerMockMessageTypes  sync.RWMutex
)

// MessageHandlerMock is a mock implementation of MessageHandler.
//
//     func TestSomethingThatUsesMessageHandler(t *testing.T) {
//
//         // make and configure a mocked MessageHandler
//         mockedMessageHandler := &MessageHandlerMock{
//             HandleMessageFunc: func(ctx context.Context, s ax.Sender, env ax.Envelope) error {
// 	               panic("TODO: mock out the HandleMessage method")
//             },
//             MessageTypesFunc: func() ax.MessageTypeSet {
// 	               panic("TODO: mock out the MessageTypes method")
//             },
//         }
//
//         // TODO: use mockedMessageHandler in code that requires MessageHandler
//         //       and then make assertions.
//
//     }
type MessageHandlerMock struct {
	// HandleMessageFunc mocks the HandleMessage method.
	HandleMessageFunc func(ctx context.Context, s ax.Sender, env ax.Envelope) error

	// MessageTypesFunc mocks the MessageTypes method.
	MessageTypesFunc func() ax.MessageTypeSet

	// calls tracks calls to the methods.
	calls struct {
		// HandleMessage holds details about calls to the HandleMessage method.
		HandleMessage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// S is the s argument value.
			S ax.Sender
			// Env is the env argument value.
			Env ax.Envelope
		}
		// MessageTypes holds details about calls to the MessageTypes method.
		MessageTypes []struct {
		}
	}
}

// HandleMessage calls HandleMessageFunc.
func (mock *MessageHandlerMock) HandleMessage(ctx context.Context, s ax.Sender, env ax.Envelope) error {
	if mock.HandleMessageFunc == nil {
		panic("moq: MessageHandlerMock.HandleMessageFunc is nil but MessageHandler.HandleMessage was just called")
	}
	callInfo := struct {
		Ctx context.Context
		S   ax.Sender
		Env ax.Envelope
	}{
		Ctx: ctx,
		S:   s,
		Env: env,
	}
	lockMessageHandlerMockHandleMessage.Lock()
	mock.calls.HandleMessage = append(mock.calls.HandleMessage, callInfo)
	lockMessageHandlerMockHandleMessage.Unlock()
	return mock.HandleMessageFunc(ctx, s, env)
}

// HandleMessageCalls gets all the calls that were made to HandleMessage.
// Check the length with:
//     len(mockedMessageHandler.HandleMessageCalls())
func (mock *MessageHandlerMock) HandleMessageCalls() []struct {
	Ctx context.Context
	S   ax.Sender
	Env ax.Envelope
} {
	var calls []struct {
		Ctx context.Context
		S   ax.Sender
		Env ax.Envelope
	}
	lockMessageHandlerMockHandleMessage.RLock()
	calls = mock.calls.HandleMessage
	lockMessageHandlerMockHandleMessage.RUnlock()
	return calls
}

// MessageTypes calls MessageTypesFunc.
func (mock *MessageHandlerMock) MessageTypes() ax.MessageTypeSet {
	if mock.MessageTypesFunc == nil {
		panic("moq: MessageHandlerMock.MessageTypesFunc is nil but MessageHandler.MessageTypes was just called")
	}
	callInfo := struct {
	}{}
	lockMessageHandlerMockMessageTypes.Lock()
	mock.calls.MessageTypes = append(mock.calls.MessageTypes, callInfo)
	lockMessageHandlerMockMessageTypes.Unlock()
	return mock.MessageTypesFunc()
}

// MessageTypesCalls gets all the calls that were made to MessageTypes.
// Check the length with:
//     len(mockedMessageHandler.MessageTypesCalls())
func (mock *MessageHandlerMock) MessageTypesCalls() []struct {
} {
	var calls []struct {
	}
	lockMessageHandlerMockMessageTypes.RLock()
	calls = mock.calls.MessageTypes
	lockMessageHandlerMockMessageTypes.RUnlock()
	return calls
}
