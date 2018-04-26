package main

import (
	"errors"
	"net/http"

	"github.com/DATA-DOG/godog"
	"github.com/gorilla/websocket"
	"github.com/segurosfalabella/imperium-worker/godogs/drivers"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var Server *http.Server
var upgrader = websocket.Upgrader{}
var confirmError error
var respond string

type message struct {
	value string
}

var receiveMessages []message

const addr = "127.0.0.1:7700"

func startServer(server *http.Server) {
	http.HandleFunc("/echo", echo)
	go http.ListenAndServe(addr, nil)
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Info("upgrade:", err)
		return
	}
	defer c.Close()
	_, m, _ := c.ReadMessage()
	receiveMessages = append(receiveMessages, message{value: string(m)})
	respond = "bad-password"
	if string(m) == "alohomora" {
		respond = "imperio"
	}
	confirmError = c.WriteMessage(websocket.TextMessage, []byte(respond))
}

func aServer() error {
	startServer(Server)
	return nil
}

func workerStarts() error {
	drivers.RunApp()
	return nil
}

func shouldServerReceive(pattern string) error {
	if receiveMessages[0].value != pattern {
		return errors.New("should server receive fail match")
	}
	return nil
}

func shouldServerSendAccepted(pattern string) error {
	if respond != pattern || confirmError != nil {
		return errors.New("should server send imperio command")
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^a server$`, aServer)
	s.Step(`^worker starts$`, workerStarts)
	s.Step(`^should server receives "(\w+)" message$`, shouldServerReceive)
	s.Step(`^should server sends "(\w+)" message$`, shouldServerSendAccepted)
}
