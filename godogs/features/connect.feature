Feature: connect to server

  @connect @old
  Scenario: Connect with server
    Given a server
    When worker starts
    Then should server receives "alohomora" message
    And should server sends "imperio" message

  @connect @old
  Scenario: Not allowed to connect with server
    Given there is a server
    When muggle worker starts
    Then should not server receives "alohomora" message
    And should server sends "avadakedavra" message

  @connect @new
  Scenario: Respond to health command
    Given a server
    And worker starts
    When server sends command "health"
    Then should worker respond "i am alive"
