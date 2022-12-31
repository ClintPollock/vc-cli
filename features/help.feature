Feature: Help
  Scenario: Running help command should return a list of available commands
    Given an API server 
    When I run `veracode help`
    Then the output should contain "Available Commands"