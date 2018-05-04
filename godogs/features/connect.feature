Feature: connect to server

  @connect @old
  Scenario: Connect with server
    Given a server
    When worker starts
    Then should server receives "alohomora" message
    And should server sends "imperio" message

  @connect @old
  Scenario: Respond to health command
    Given a server
    And worker starts and login
    When server sends command "health"
    Then should worker respond "i am alive"

  @connect @old
  Scenario: Execute a job successfully
    Given a server
    And worker starts and login
    When server sends job with image "falabellacr/imperium-job-dummy" and arguments "0"
    Then worker should respond exit code "0"

  @connect @old
  Scenario: Execute a job with a failure
    Given a server
    And worker starts and login
    When server sends job with image "falabellacr/imperium-job-dummy" and arguments "5"
    Then worker should respond exit code "5"
