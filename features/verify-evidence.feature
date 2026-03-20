Feature: Run verification evidence
  In order to prove that a Go repository is healthy
  As a developer using `intend`
  I want `intend verify` to check contract traceability first and then run the repo verification tools

  Scenario: Verify a repository successfully
    Given an initialized owned repository
    And a locked bundle `health-check`
    And verification tools are available
    When I run `intend verify`
    Then it exits with code 0
    And stdout contains `verify ok`
    And the verification log contains lines in order
      | go test ./...              |
      | golangci-lint run         |
      | trufflehog filesystem .   |
      | gitleaks dir .            |
      | trivy fs .                |

  Scenario: Fail when a required verification tool is missing
    Given an initialized owned repository
    And a locked bundle `health-check`
    And verification tools are available except `trivy`
    When I run `intend verify`
    Then it exits with code 1
    And stderr contains `required tool not found: trivy`

  Scenario: Stop verification when contract drift exists
    Given an initialized owned repository
    And a locked bundle `health-check`
    And verification tools are available
    When I replace the contents of `specs/health-check.md`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `contract drift`
    And the verification log is empty

  Scenario: Stop verification when contribution contract drift exists
    Given a git repository
    And GitHub issue import is available
    And a locked contribution bundle `startup-fix`
    And verification tools are available
    When I replace the contents of `.git/intend/contrib/startup-fix/specs/startup-fix.md`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `contract drift`
    And the verification log is empty
