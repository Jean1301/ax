// Code generated by moq; DO NOT EDIT
// github.com/matryer/moq

package persistencetest

import (
	"github.com/jmalloc/ax/src/ax/persistence"
	"sync"
)

var (
	lockTxMockDataStore sync.RWMutex
)

// TxMock is a mock implementation of Tx.
//
//     func TestSomethingThatUsesTx(t *testing.T) {
//
//         // make and configure a mocked Tx
//         mockedTx := &TxMock{
//             DataStoreFunc: func() persistence.DataStore {
// 	               panic("TODO: mock out the DataStore method")
//             },
//         }
//
//         // TODO: use mockedTx in code that requires Tx
//         //       and then make assertions.
//
//     }
type TxMock struct {
	// DataStoreFunc mocks the DataStore method.
	DataStoreFunc func() persistence.DataStore

	// calls tracks calls to the methods.
	calls struct {
		// DataStore holds details about calls to the DataStore method.
		DataStore []struct {
		}
	}
}

// DataStore calls DataStoreFunc.
func (mock *TxMock) DataStore() persistence.DataStore {
	if mock.DataStoreFunc == nil {
		panic("moq: TxMock.DataStoreFunc is nil but Tx.DataStore was just called")
	}
	callInfo := struct {
	}{}
	lockTxMockDataStore.Lock()
	mock.calls.DataStore = append(mock.calls.DataStore, callInfo)
	lockTxMockDataStore.Unlock()
	return mock.DataStoreFunc()
}

// DataStoreCalls gets all the calls that were made to DataStore.
// Check the length with:
//     len(mockedTx.DataStoreCalls())
func (mock *TxMock) DataStoreCalls() []struct {
} {
	var calls []struct {
	}
	lockTxMockDataStore.RLock()
	calls = mock.calls.DataStore
	lockTxMockDataStore.RUnlock()
	return calls
}

var (
	lockCommitterMockCommit   sync.RWMutex
	lockCommitterMockRollback sync.RWMutex
)

// CommitterMock is a mock implementation of Committer.
//
//     func TestSomethingThatUsesCommitter(t *testing.T) {
//
//         // make and configure a mocked Committer
//         mockedCommitter := &CommitterMock{
//             CommitFunc: func() error {
// 	               panic("TODO: mock out the Commit method")
//             },
//             RollbackFunc: func() error {
// 	               panic("TODO: mock out the Rollback method")
//             },
//         }
//
//         // TODO: use mockedCommitter in code that requires Committer
//         //       and then make assertions.
//
//     }
type CommitterMock struct {
	// CommitFunc mocks the Commit method.
	CommitFunc func() error

	// RollbackFunc mocks the Rollback method.
	RollbackFunc func() error

	// calls tracks calls to the methods.
	calls struct {
		// Commit holds details about calls to the Commit method.
		Commit []struct {
		}
		// Rollback holds details about calls to the Rollback method.
		Rollback []struct {
		}
	}
}

// Commit calls CommitFunc.
func (mock *CommitterMock) Commit() error {
	if mock.CommitFunc == nil {
		panic("moq: CommitterMock.CommitFunc is nil but Committer.Commit was just called")
	}
	callInfo := struct {
	}{}
	lockCommitterMockCommit.Lock()
	mock.calls.Commit = append(mock.calls.Commit, callInfo)
	lockCommitterMockCommit.Unlock()
	return mock.CommitFunc()
}

// CommitCalls gets all the calls that were made to Commit.
// Check the length with:
//     len(mockedCommitter.CommitCalls())
func (mock *CommitterMock) CommitCalls() []struct {
} {
	var calls []struct {
	}
	lockCommitterMockCommit.RLock()
	calls = mock.calls.Commit
	lockCommitterMockCommit.RUnlock()
	return calls
}

// Rollback calls RollbackFunc.
func (mock *CommitterMock) Rollback() error {
	if mock.RollbackFunc == nil {
		panic("moq: CommitterMock.RollbackFunc is nil but Committer.Rollback was just called")
	}
	callInfo := struct {
	}{}
	lockCommitterMockRollback.Lock()
	mock.calls.Rollback = append(mock.calls.Rollback, callInfo)
	lockCommitterMockRollback.Unlock()
	return mock.RollbackFunc()
}

// RollbackCalls gets all the calls that were made to Rollback.
// Check the length with:
//     len(mockedCommitter.RollbackCalls())
func (mock *CommitterMock) RollbackCalls() []struct {
} {
	var calls []struct {
	}
	lockCommitterMockRollback.RLock()
	calls = mock.calls.Rollback
	lockCommitterMockRollback.RUnlock()
	return calls
}
