package app

import (
	"errors"
	"fmt"
	"regexp"

	"google.golang.org/grpc"
)

// GrpcDial ...
var GrpcDial = grpc.Dial

// Start method
func Start(address string) error {
	if address == "" {
		return errors.New("missing server address")
	}

	if match, err := regexp.MatchString(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d+`, address); err != nil || !match {
		return errors.New("server address invalid")
	}

	conn, err := GrpcDial(address, grpc.WithInsecure())
	if err != nil {
		return errors.New("connection fail")
	}

	fmt.Println(conn.GetState().String())

	fmt.Println("Hello World!")
	return nil
}
