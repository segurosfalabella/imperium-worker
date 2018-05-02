package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/segurosfalabella/imperium-worker/godogs/drivers"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var serverRequestChannel = make(chan string)
var serverResponseChannel = make(chan string)

type commandMessage struct {
	Command string
}

func aServer() error {
	drivers.StartServer(serverRequestChannel, serverResponseChannel)
	time.Sleep(100 * time.Microsecond)
	return nil
}

func workerStarts() error {
	drivers.RunApp()
	time.Sleep(100 * time.Microsecond)
	return nil
}

func workerStartsAndLogin() error {
	drivers.RunApp()
	time.Sleep(10 * time.Microsecond)
	<-serverResponseChannel
	return nil
}

func shouldServerReceive(pattern string) error {
	if drivers.NotExistsPattern(pattern) {
		return errors.New("should server receive fail match")
	}
	return nil
}

func shouldServerSendAccepted(pattern string) error {
	actualResponse := <-serverResponseChannel
	if actualResponse != pattern || drivers.HasError() {
		return errors.New("should server send imperio command")
	}
	return nil
}

func serverSendsCommand(command string) error {
	message := &commandMessage{Command: command}
	bb, _ := json.Marshal(message)
	serverRequestChannel <- string(bb)
	return nil
}

func shouldWorkerRespond(response string) error {
	actualResponse := <-serverResponseChannel
	if response != actualResponse {
		return fmt.Errorf("%s != %s", actualResponse, response)
	}
	return nil
}

func afterScenario(arg1 interface{}, arg2 error) {
	drivers.CloseServer()
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^a server$`, aServer)
	s.Step(`^worker starts$`, workerStarts)
	s.Step(`^worker starts and login$`, workerStartsAndLogin)
	s.Step(`^should server receives "(\w+)" message$`, shouldServerReceive)
	s.Step(`^should server sends "(\w+)" message$`, shouldServerSendAccepted)
	s.Step(`^server sends command "([^"]*)"$`, serverSendsCommand)
	s.Step(`^should worker respond "([^"]*)"$`, shouldWorkerRespond)

	s.AfterScenario(afterScenario)
}
