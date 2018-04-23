package receiver_test

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/segurosfalabella/imperium-worker/receiver"
	"github.com/stretchr/testify/assert"
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

func TestShouldExecuteJobWhenMessageParseSuccess(t *testing.T) {
	mockConn := new(MockConn)
	mockConn.On("ReadMessage").Return(websocket.TextMessage, []byte(`{"name":"dummy","description":"dummy description","command":"exit"}`), nil)
	mockJob := new(MockJob)
	mockJob.On("Execute").Return(nil)

	receiver.Start(mockConn, mockJob)
	assert.True(t, mockJob.AssertCalled(t, "Execute"))
}
