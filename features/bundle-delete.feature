Feature: Delete a bundle
  In order to remove abandoned contract work safely
  As a developer using `intend`
  I want to delete owned and contribution bundles without leaving stale trace state behind

  Scenario: Delete an unlocked owned bundle
    Given an existing bundle `health-check`
    When I run `intend delete health-check`
    Then it exits with code 0
    And stdout contains `deleted bundle health-check`
    And the path `specs/health-check.md` does not exist
    And the path `features/health-check.feature` does not exist
    And the path `.intend/trace/health-check.json` does not exist
    And the path `.intend/locks/health-check.json` does not exist

  Scenario: Refuse to delete a locked owned bundle without force
    Given a locked bundle `health-check`
    When I run `intend delete health-check`
    Then it exits with code 1
    And stderr contains `bundle health-check is locked`
    And the file `specs/health-check.md` exists
    And the file `features/health-check.feature` exists
    And the file `.intend/trace/health-check.json` exists
    And the file `.intend/locks/health-check.json` exists

  Scenario: Delete a locked owned bundle with force
    Given a locked bundle `health-check`
    When I run `intend delete --force health-check`
    Then it exits with code 0
    And stdout contains `deleted bundle health-check`
    And the path `specs/health-check.md` does not exist
    And the path `features/health-check.feature` does not exist
    And the path `.intend/trace/health-check.json` does not exist
    And the path `.intend/locks/health-check.json` does not exist

  Scenario: Delete a contribution bundle
    Given an existing contribution bundle `startup-fix`
    When I run `intend delete --mode contrib startup-fix`
    Then it exits with code 0
    And stdout contains `deleted contribution bundle startup-fix`
    And the path `.git/intend/contrib/startup-fix` does not exist

  Scenario: Stop tracing a deleted bundle during verify
    Given a locked bundle `health-check`
    And a locked bundle `other-check`
    And verification tools are available
    When I run `intend delete --force health-check`
    And I run `intend verify`
    Then it exits with code 0
    And stdout contains `trace: checking owned bundle other-check`
    And stdout does not contain `trace: checking owned bundle health-check`
    And the verification log contains lines in order
      | go test ./...            |
      | golangci-lint run        |
      | trufflehog filesystem .  |
      | gitleaks dir .           |
      | trivy fs .               |
