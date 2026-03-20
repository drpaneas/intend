Feature: init-tool-preflight
  Scenario: Refuse init when a required verification tool is missing
    Given an empty working directory
    And init tools are available except `trivy`
    When I run `intend init`
    Then it exits with code 1
    And stderr contains `required tool not found: trivy`
    And the path `.intend` does not exist
    And the path `specs` does not exist
