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

var oldWsDial func(string, http.Header) (*websocket.Conn, *http.Response, error)

type WebsocketConnectionDouble struct {
	mock.Mock
}

func TestShouldFailWhenNotHaveParamServerAddress(t *testing.T) {
	var address string

	err := app.Start(address)

	assert.NotNil(t, err)
	assert.Equal(t, "missing server address", err.Error())
}

func TestShouldFailWhenServerAddressFormatInvalid(t *testing.T) {
	var address = "helloworld"

	err := app.Start(address)

	assert.NotNil(t, err)
	assert.Equal(t, "server address invalid", err.Error())
}

func TestShouldFailWhenWebsocketDialFail(t *testing.T) {
	var address = "127.0.0.1:7700"
	createWsDialMock(nil, errors.New("connection refused"))
	defer restoreWsDialMock()

	err := app.Start(address)

	assert.NotNil(t, err)
	assert.Equal(t, "connection refused", err.Error())
}

func TestShouldNotFailWithValidAddressWithIp(t *testing.T) {
	var address = "127.0.0.1:7700"
	createWsDialMock(nil, nil)
	defer restoreWsDialMock()

	err := app.Start(address)

	assert.Nil(t, err)
}

func TestShouldNotFailWithValidAddressWithHostname(t *testing.T) {
	var address = "localhost.tld:7700"
	createWsDialMock(nil, nil)
	defer restoreWsDialMock()

	err := app.Start(address)

	assert.Nil(t, err)
}

// func TestShouldConnectionIsAliveWhenWsDialSuccess(t *testing.T) {
// 	var address = "localhost:7700"
// 	createWsDialMock(nil, nil)
// 	defer restoreWsDialMock()
//
// 	err := app.Start(address)
//
// 	assert.Nil(t, err)
// }

func createWsDialMock(connection *websocket.Conn, err error) {
	oldWsDial = app.WsDial
	app.WsDial = func(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error) {
		return connection, nil, err
	}
}

func restoreWsDialMock() {
	app.WsDial = oldWsDial
}
