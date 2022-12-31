Feature: Image SBOM
  Scenario: Running sbom command should return results with scanId
    Given an activated user
    When I run `veracode sbom image alpine:latest`
    Then the output should contain "artifacts"
    And the output should contain "alpine:latest"

  Scenario: Running sbom command on alpine:latest should match expected output
    Given an activated user
    When I run `veracode sbom image alpine:latest`
    And the output should match json in "alpine-latest.json"

  Scenario: Running sbom command on alpine:latest with user with invalid credentials
    Given an user with invalid credentials
    When I run `veracode sbom image alpine:latest`
    And the output should contain "HMAC credentials not associated with a valid Veracode account"
