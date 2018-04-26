package receiver_test

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/segurosfalabella/imperium-worker/receiver"
	"github.com/stretchr/testify/mock"
)

type MockConn struct {
	mock.Mock
}

func (conn *MockConn) Close() error {
	args := conn.Called()
	return args.Error(0)
}

func (conn *MockConn) ReadMessage() (messageType int, p []byte, err error) {
	args := conn.Called()
	return args.Int(0), args.Get(1).([]byte), args.Error(2)
}

func (conn *MockConn) WriteMessage(messageType int, data []byte) error {
	args := conn.Called(messageType, data)
	return args.Error(0)
}

type MockJob struct {
	mock.Mock
}

func (job *MockJob) Execute() error {
	args := job.Called()
	return args.Error(0)
}

func TestShouldFailAuthWhenPasswordToMatch(t *testing.T) {
	mockConn := new(MockConn)
	mockConn.On("WriteMessage", websocket.TextMessage, mock.Anything).Return(nil)
	mockConn.On("ReadMessage").Return(websocket.TextMessage, []byte("bad-password"), nil)
	mockJob := new(MockJob)

	receiver.Start(mockConn, mockJob)

	mockConn.AssertNumberOfCalls(t, "WriteMessage", 1)
	mockConn.AssertNumberOfCalls(t, "ReadMessage", 1)
	mockJob.AssertNotCalled(t, "Execute")
}

func TestShouldExecuteJobWhenMessageParseSuccess(t *testing.T) {
	mockConn := new(MockConn)
	mockConn.On("WriteMessage", websocket.TextMessage, mock.Anything).Return(nil)
	mockConn.On("ReadMessage").Return(websocket.TextMessage, []byte("imperio"), nil).Once()
	mockConn.On("ReadMessage").Return(websocket.TextMessage, []byte(`{"name":"dummy","description":"dummy description","command":"exit"}`), nil)
	mockJob := new(MockJob)
	mockJob.On("Execute").Return(nil)

	receiver.Start(mockConn, mockJob)

	mockConn.AssertNumberOfCalls(t, "WriteMessage", 1)
	mockConn.AssertNumberOfCalls(t, "ReadMessage", 2)
	mockJob.AssertCalled(t, "Execute")
}
