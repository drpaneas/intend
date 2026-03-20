Feature: Create a contribution shadow bundle
  In order to work on upstream issues without polluting the upstream repository
  As a contributor using `intend`
  I want to import GitHub issue intent into a local shadow bundle under Git metadata

  Scenario: Create a contribution bundle from a GitHub issue
    Given a git repository
    And GitHub issue import is available
    When I run `intend new --mode contrib --from-gh owner/repo#123 startup-fix`
    Then it exits with code 0
    And the file `.git/intend/contrib/startup-fix/issue.json` exists
    And the file `.git/intend/contrib/startup-fix/specs/startup-fix.md` exists
    And the file `.git/intend/contrib/startup-fix/features/startup-fix.feature` exists
    And the file `.git/intend/contrib/startup-fix/trace/startup-fix.json` exists
    And the file `.git/intend/contrib/startup-fix/issue.json` contains `Fix crash on startup`
    And the file `.git/intend/contrib/startup-fix/trace/startup-fix.json` contains `"mode": "contrib"`

  Scenario: Reject contribution mode outside a git repository
    Given an empty working directory
    And GitHub issue import is available
    When I run `intend new --mode contrib --from-gh owner/repo#123 startup-fix`
    Then it exits with code 1
    And stderr contains `not inside a git repository`

  Scenario: Fail when gh is unavailable
    Given a git repository
    And GitHub issue import is unavailable
    When I run `intend new --mode contrib --from-gh owner/repo#123 startup-fix`
    Then it exits with code 1
    And stderr contains `required tool not found: gh`

  Scenario Outline: Reject a contribution bundle request with a non-positive issue number
    Given a git repository
    When I run `intend new --mode contrib --from-gh <issueRef> startup-fix`
    Then it exits with code 1
    And stderr contains `invalid GitHub issue reference "<issueRef>"`
    And the path `.git/intend/contrib/startup-fix` does not exist

    Examples:
      | issueRef      |
      | owner/repo#0  |
      | owner/repo#-1 |
