package receiver_test

import (
	"errors"
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

func (job *MockJob) FromJSON(text string) {
	job.Called(text)
}

func (job *MockJob) ToJSON() string {
	args := job.Called()
	return args.String(0)
}

func TestShouldFailAuthWhenPasswordNotMatch(t *testing.T) {
	mockConn := new(MockConn)
	mockConn.On("WriteMessage", websocket.TextMessage, mock.Anything).Return(nil)
	mockConn.On("ReadMessage").Return(websocket.TextMessage, []byte("bad-password"), nil)
	mockJob := new(MockJob)

	receiver.Start(mockConn, mockJob)

	mockConn.AssertNumberOfCalls(t, "WriteMessage", 1)
	mockConn.AssertNumberOfCalls(t, "ReadMessage", 1)
	mockJob.AssertNotCalled(t, "Execute")
}

func TestShouldNotExecuteWhenParseJobFail(t *testing.T) {
	mockConn := new(MockConn)
	mockConn.On("WriteMessage", websocket.TextMessage, mock.Anything).Return(nil)
	mockConn.On("ReadMessage").Return(websocket.TextMessage, []byte("imperio"), nil).Once()
	mockConn.On("ReadMessage").Return(-1, []byte(""), nil)
	mockJob := new(MockJob)
	mockJob.On("Execute").Return(errors.New("dummy error"))

	receiver.Start(mockConn, mockJob)

	mockConn.AssertNumberOfCalls(t, "WriteMessage", 1)
	mockConn.AssertNumberOfCalls(t, "ReadMessage", 2)
	mockJob.AssertNotCalled(t, "Execute")
}

func TestShouldRespondErrorWhenExecuteFail(t *testing.T) {
	mockConn := new(MockConn)
	mockConn.On("WriteMessage", websocket.TextMessage, mock.Anything).Return(nil)
	mockConn.On("ReadMessage").Return(websocket.TextMessage, []byte("imperio"), nil).Once()
	mockConn.On("ReadMessage").Return(websocket.TextMessage, []byte(`{"name":"dummy","description":"dummy description","command":"exit"}`), nil)
	mockJob := new(MockJob)
	mockJob.On("FromJSON", mock.Anything).Return()
	mockJob.On("Execute").Return(errors.New("dummy error"))
	mockJob.On("ToJSON").Return("6a41ee8c-f942-42c9-8904-5fba1b5854d7")

	receiver.Start(mockConn, mockJob)

	mockConn.AssertNumberOfCalls(t, "WriteMessage", 2)
	mockConn.AssertNumberOfCalls(t, "ReadMessage", 2)
	mockJob.AssertCalled(t, "Execute")
}

func TestShouldRespondResponseWhenExecuteSucceed(t *testing.T) {
	mockConn := new(MockConn)
	mockConn.On("WriteMessage", websocket.TextMessage, mock.Anything).Return(nil)
	mockConn.On("ReadMessage").Return(websocket.TextMessage, []byte("imperio"), nil).Once()
	mockConn.On("ReadMessage").Return(websocket.TextMessage, []byte(`{"name":"dummy","description":"dummy description","command":"exit"}`), nil)
	mockJob := new(MockJob)
	mockJob.On("FromJSON", mock.Anything).Return()
	mockJob.On("Execute").Return(nil)
	mockJob.On("ToJSON").Return("6a41ee8c-f942-42c9-8904-5fba1b5854d7")

	receiver.Start(mockConn, mockJob)

	mockConn.AssertNumberOfCalls(t, "WriteMessage", 2)
	mockConn.AssertNumberOfCalls(t, "ReadMessage", 2)
	mockConn.AssertCalled(t, "WriteMessage", websocket.TextMessage, []byte("6a41ee8c-f942-42c9-8904-5fba1b5854d7"))
	mockJob.AssertCalled(t, "Execute")
	mockJob.AssertCalled(t, "ToJSON")
}
