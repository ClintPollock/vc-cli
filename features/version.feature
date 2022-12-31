Feature: Print Version
  Scenario: Running version command should provide Veracode version
    When I run `veracode version`
    Then the output should contain the version and hash