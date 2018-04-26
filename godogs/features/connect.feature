Feature: connect to server

  Scenario: Connect with server
    Given a server
    When worker starts
    Then should server receives "alohomora" message
    And should server sends "imperio" message
