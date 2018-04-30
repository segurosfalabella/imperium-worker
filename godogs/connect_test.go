package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/gorilla/websocket"
	"github.com/segurosfalabella/imperium-worker/connection"
	"github.com/segurosfalabella/imperium-worker/godogs/drivers"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var Server *http.Server
var upgrader = websocket.Upgrader{}
var confirmError error
var respond string
var serverChannel = make(chan string)

type message struct {
	value string
}

type commandMessage struct {
	Command string
}

var receiveMessages []message

const addr = "127.0.0.1:7700"

func startServer(server *http.Server) {
	http.HandleFunc("/echo", echo)
	go http.ListenAndServe(addr, nil)
}

type websocketDialerShim struct {
	*websocket.Dialer
}

func (s websocketDialerShim) Dial(urlStr string) (connection.WsConn, error) {
	conn, _, err := s.Dialer.Dial(urlStr, nil)
	return conn, err
}

const mugglePassword = "expelliarmus"

func createMuggleWorker() error {
	conn, err := connection.Create(addr, new(websocketDialerShim))
	if err != nil {
		return errors.New("cannot create muggle worker")
	}
	conn.WriteMessage(websocket.TextMessage, []byte(mugglePassword))
	return nil
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
	respond = "avadakedavra"
	if string(m) == "alohomora" {
		respond = "imperio"
	}
	confirmError = c.WriteMessage(websocket.TextMessage, []byte(respond))
	go func() {
		select {
		case m := <-serverChannel:
			log.Info("command: ", m)
			c.WriteMessage(websocket.TextMessage, []byte(m))
		}
	}()
}

func aServer() error {
	startServer(Server)
	time.Sleep(1 * time.Millisecond)
	return nil
}

func workerStarts() error {
	drivers.RunApp()
	return nil
}

func shouldServerReceive(pattern string) error {
	if !stringInSlice(pattern, receiveMessages) {
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

func serverSendsCommand(command string) error {
	log.Info(command)
	message := &commandMessage{Command: command}
	bb, _ := json.Marshal(message)
	serverChannel <- string(bb)
	return nil
}

func shouldWorkerRespond(response string) error {
	log.Info(response)
	actualResponse := <-serverChannel
	if response != actualResponse {
		return fmt.Errorf("%s != %s", actualResponse, response)
	}
	return nil
}

func thereIsAServer() error {
	return nil
}

func muggleWorkerStarts() error {
	createMuggleWorker()
	return nil
}

func shouldNotServerReceivesMessage(pattern string) error {
	if !stringInSlice(pattern, receiveMessages) {
		return errors.New("should not server receives alohomora")
	}
	return nil
}

func stringInSlice(a string, list []message) bool {
	for _, b := range list {
		if b.value == a {
			return true
		}
	}
	return false
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^a server$`, aServer)
	s.Step(`^worker starts$`, workerStarts)
	s.Step(`^should server receives "(\w+)" message$`, shouldServerReceive)
	s.Step(`^should server sends "(\w+)" message$`, shouldServerSendAccepted)

	s.Step(`^there is a server$`, thereIsAServer)
	s.Step(`^muggle worker starts$`, muggleWorkerStarts)
	s.Step(`^should not server receives "(\w+)" message$`, shouldNotServerReceivesMessage)

	s.Step(`^server sends command "(\w+)"$`, serverSendsCommand)
	s.Step(`^should worker respond "([^"]*)"$`, shouldWorkerRespond)
}
