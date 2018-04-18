package app_test

import (
	"errors"
	"imperium-worker/app"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

//
// type clicon struct{
// 	grpc.ClientConn
// }

func TestShouldFailWhenNotHaveParamServerAddress(t *testing.T) {
	var address string = ""

	err := app.Start(address)

	assert.NotNil(t, err)
	assert.Equal(t, "missing server address", err.Error())
}

func TestShouldFailWhenServerAddressFormatInvalid(t *testing.T) {
	var address string = "helloworld"

	err := app.Start(address)

	assert.NotNil(t, err)
	assert.Equal(t, "server address invalid", err.Error())
}

func TestShouldFailWhenGrpcDialFail(t *testing.T) {
	var address string = "127.0.0.1:8000"
	oldGrpcDial := app.GrpcDial
	defer func() { app.GrpcDial = oldGrpcDial }()
	app.GrpcDial = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return nil, errors.New("Done!")
	}

	err := app.Start(address)

	assert.NotNil(t, err)
	assert.Equal(t, "connection fail", err.Error())
}

func TestShouldReturnNoErrorWhenGrpcDialSucceed(t *testing.T) {
	var address string = "127.0.0.1:8000"
	oldGrpcDial := app.GrpcDial
	defer func() { app.GrpcDial = oldGrpcDial }()
	app.GrpcDial = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		var connection = grpc.ClientConn{}
		return &connection, nil
	}

	err := app.Start(address)

	assert.Nil(t, err)
}
