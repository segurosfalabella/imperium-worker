package connection_test

import (
	"errors"
	"testing"

	"github.com/segurosfalabella/imperium-worker/connection"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDialer struct {
	mock.Mock
}

func (dialer *MockDialer) Dial(urlStr string) (connection.WsConn, error) {
	args := dialer.Called(urlStr)
	return args.Get(0).(connection.WsConn), args.Error(1)
}

type MockConn struct {
	mock.Mock
}

func (conn *MockConn) Close() error {
	return nil
}

func (conn *MockConn) ReadMessage() (messageType int, p []byte, err error) {
	return 0, nil, nil
}

func (conn *MockConn) WriteMessage(messageType int, data []byte) error {
	return nil
}

func TestShouldFailWhenNotHaveParamServerAddress(t *testing.T) {
	var address string

	conn, err := connection.Create(address, new(MockDialer))

	assert.Nil(t, conn)
	assert.NotNil(t, err)
	assert.Equal(t, "missing server address", err.Error())
}

func TestShouldFailWhenServerAddressFormatInvalid(t *testing.T) {
	var address = "helloworld"

	conn, err := connection.Create(address, new(MockDialer))

	assert.Nil(t, conn)
	assert.NotNil(t, err)
	assert.Equal(t, "server address invalid", err.Error())
}

func TestShouldFailWhenWebsocketDialFail(t *testing.T) {
	var address = "127.0.0.1:7700"
	mockDialer := new(MockDialer)
	mockDialer.On("Dial", "ws://127.0.0.1:7700/azkaban").Return(new(MockConn), errors.New("connection refused"))

	conn, err := connection.Create(address, mockDialer)

	assert.Nil(t, conn)
	assert.NotNil(t, err)
	assert.Equal(t, "connection refused", err.Error())
}

func TestShouldNotFailWithValidAddressWithIp(t *testing.T) {
	var address = "127.0.0.1:7700"
	mockConn := new(MockConn)
	mockDialer := new(MockDialer)
	mockDialer.On("Dial", "ws://127.0.0.1:7700/azkaban").Return(mockConn, nil)

	conn, err := connection.Create(address, mockDialer)

	assert.NotNil(t, conn)
	assert.Nil(t, err)
}

func TestShouldNotFailWithValidAddressWithHostname(t *testing.T) {
	var address = "localhost.tld:7700"
	mockConn := new(MockConn)
	testObj := new(MockDialer)
	testObj.On("Dial", "ws://localhost.tld:7700/azkaban").Return(mockConn, nil)

	conn, err := connection.Create(address, testObj)

	assert.NotNil(t, conn)
	assert.Nil(t, err)
}
