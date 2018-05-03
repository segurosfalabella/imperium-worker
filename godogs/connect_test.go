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
	Command   string
	Image     string
	Arguments string
	ExitCode  int
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
	time.Sleep(100 * time.Microsecond)
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
	binary, _ := json.Marshal(message)
	serverRequestChannel <- string(binary)
	return nil
}

func shouldWorkerRespond(response string) error {
	actualResponse := <-serverResponseChannel
	if response != actualResponse {
		return fmt.Errorf("%s != %s", actualResponse, response)
	}
	return nil
}

func serverSendsJobWithImageAndArguments(image string, args string) error {
	message := &commandMessage{
		Image:     image,
		Arguments: args,
	}
	binary, _ := json.Marshal(message)
	serverRequestChannel <- string(binary)
	return nil
}

func workerShouldRespondExitCode(code int) error {
	actualResponse := <-serverResponseChannel
	response := &commandMessage{}
	json.Unmarshal([]byte(actualResponse), &response)
	if code != response.ExitCode {
		return fmt.Errorf("%d != %d", code, response.ExitCode)
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
	s.Step(`^server sends job with image "([^"]*)" and arguments "([^"]*)"$`, serverSendsJobWithImageAndArguments)
	s.Step(`^worker should respond exit code "([^"]*)"$`, workerShouldRespondExitCode)

	s.AfterScenario(afterScenario)
}
