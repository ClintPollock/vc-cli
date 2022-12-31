Feature: Image Scan
  Scenario: Running an image scan should return results with policy-results
    Given an activated user
    When I run `veracode scan image alpine:latest`
    Then the output should contain "policy-results"

  Scenario: Running scan command on alpine:latest with user with invalid credentials
    Given an user with invalid credentials
    When I run `veracode scan image alpine:latest`
    And the output should contain "HMAC credentials not associated with a valid Veracode account"