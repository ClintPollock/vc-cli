Feature: Clear Cache
  Scenario: Running clear command should clear Veracode cli cache
    Given existing cache
    When I run `veracode clear`
    Then the cache directory should be cleared