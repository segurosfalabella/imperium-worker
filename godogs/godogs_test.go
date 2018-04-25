package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DATA-DOG/godog"
	"github.com/segurosfalabella/imperium-worker/connection"
)

func aServer() error {
	Server = &http.Server{Addr: *addr}
	err := Server.ListenAndServe()

	if err != nil {
		return nil
	}
	return errors.New("there is no server in godogs var")
}

func workerStarts() error {
	conn, err := connection.Create(*addr, new(WebsocketDialerShim))
	if (err != nil) && (conn != nil) {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func itShouldConnect() error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^a server$`, aServer)
	s.Step(`^worker starts$`, workerStarts)
	s.Step(`^it should connect$`, itShouldConnect)
}
