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
