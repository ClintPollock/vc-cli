Feature: CLI init
  Scenario: Initiate Veracode CLI authentication
    Given a server that authenticates the CLI
    When I run `veracode init` interactively
    Then the output should contain "Configuring credentials"

  Scenario: Ask for api key and id, create credentials file
    Given a server that authenticates the CLI
    When I run `veracode init` interactively
    When I type key and secret
    Then the output should contain "Configuring credentials"
    And the output should contain "API ID"
    And the output should contain "API Secret Key"
    And the output should contain "Wrote configuration to"
    And the exit status should be 0
    And the file "~/.veracode/credentials" should exist
    And check key and secret present in the credentials file
  
  Scenario: Invalid credentials
    Given a server that authenticates the CLI
    When I run `veracode init` interactively
    When I type "foo"
    And I type "bar"
    Then the output should contain "Configuring credentials"
    And the output should contain "API ID"
    And the output should contain "API Secret Key"
    And the output should contain "HMAC credentials not associated with a valid Veracode account"
    And the exit status should be 0
