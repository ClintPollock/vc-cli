Feature: Image research
  Scenario: Running research command should return results
    Given an activated user
    When I run `veracode research image alpine:latest`
    Then the output should contain "target"
    And the output should contain "RepoTags"

  Scenario: Running research command on alpine:latest with user with invalid credentials
    Given an user with invalid credentials
    When I run `veracode research image alpine:latest`
    And the output should contain "HMAC credentials not associated with a valid Veracode account"