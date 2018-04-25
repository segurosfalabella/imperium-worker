# file: $GOPATH/src/imperium-worker/godogs/features/dummy-worker.feature
Feature: connect to server
  In order to be connected
  As a worker
  I need to be able to connect with server

  Scenario: Connect with server
    Given a server
    When worker starts
    Then it should connect
