Feature: Report verify progress clearly
  In order to understand what `intend verify` is doing
  As a developer running local verification
  I want `intend verify` to show trace and tool progress in plain text

  Scenario: Show verify progress in order
    Given an initialized owned repository
    And a locked bundle `health-check`
    And verification tools are available
    When I run `intend verify`
    Then it exits with code 0
    And stdout contains lines in order
      | trace: checking owned bundle health-check |
      | verify: running go test ./...            |
      | verify: running golangci-lint run        |
      | verify: running trufflehog filesystem .  |
      | verify: running gitleaks dir .           |
      | verify: running trivy fs .               |
      | verify ok                                |

  Scenario: Show command-specific usage for verify
    Given an initialized owned repository
    When I run `intend verify extra`
    Then it exits with code 2
    And stderr contains `usage: intend verify`
