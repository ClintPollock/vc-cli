Feature: Image inspect
  Scenario: Running inspect command should return results
    Given an activated user
    When I run `veracode inspect image alpine:latest`
    Then the output should contain "alpine_latest"

  Scenario: Running inspect command on alpine:latest with user with invalid credentials
    Given an user with invalid credentials
    When I run `veracode inspect image alpine:latest`
    And the output should contain "HMAC credentials not associated with a valid Veracode account"