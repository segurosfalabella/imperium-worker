package app_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/segurosfalabella/imperium-worker/app"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDialer struct {
	mock.Mock
}

func (dialer *MockDialer) Dial(urlStr string, requestHeader http.Header) (app.WsConn, *http.Response, error) {
	args := dialer.Called(urlStr)
	return args.Get(0).(app.WsConn), nil, args.Error(2)
}

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

func TestShouldFailWhenNotHaveParamServerAddress(t *testing.T) {
	var address string

	err := app.Start(address, &MockDialer{})

	assert.NotNil(t, err)
	assert.Equal(t, "missing server address", err.Error())
}

func TestShouldFailWhenServerAddressFormatInvalid(t *testing.T) {
	var address = "helloworld"

	err := app.Start(address, &MockDialer{})

	assert.NotNil(t, err)
	assert.Equal(t, "server address invalid", err.Error())
}

func TestShouldFailWhenWebsocketDialFail(t *testing.T) {
	var address = "127.0.0.1:7700"
	testObj := new(MockDialer)
	testObj.On("Dial", "ws://127.0.0.1:7700/echo").Return(&websocket.Conn{}, nil, errors.New("connection refused"))

	err := app.Start(address, testObj)

	assert.NotNil(t, err)
	assert.Equal(t, "connection refused", err.Error())
}

func TestShouldNotFailWithValidAddressWithIp(t *testing.T) {
	var address = "127.0.0.1:7700"
	mockConn := new(MockConn)
	mockConn.On("ReadMessage").Return(websocket.TextMessage, []byte("Hello World"), nil)
	mockConn.On("WriteMessage", 1, []byte("e16b7b57-3eab-4866-805a-81ccc15a01ac")).Return(nil)
	mockDialer := new(MockDialer)
	mockDialer.On("Dial", "ws://127.0.0.1:7700/echo").Return(mockConn, nil, nil)

	err := app.Start(address, mockDialer)

	assert.Nil(t, err)
}

func TestShouldNotFailWithValidAddressWithHostname(t *testing.T) {
	var address = "localhost.tld:7700"
	mockConn := new(MockConn)
	mockConn.On("ReadMessage").Return(websocket.TextMessage, []byte("Hello World"), nil)
	mockConn.On("WriteMessage", 1, []byte("e16b7b57-3eab-4866-805a-81ccc15a01ac")).Return(nil)
	testObj := new(MockDialer)
	testObj.On("Dial", "ws://localhost.tld:7700/echo").Return(mockConn, nil, nil)

	err := app.Start(address, testObj)

	assert.Nil(t, err)
}
