Feature: Manage a contribution shadow bundle lifecycle
  In order to keep contribution intent explicit and reviewable
  As a contributor using `intend`
  I want to lock, trace, and amend a shadow bundle under Git metadata

  Scenario: Lock a contribution bundle
    Given a git repository
    And GitHub issue import is available
    And an existing contribution bundle `startup-fix`
    When I run `intend lock --mode contrib startup-fix`
    Then it exits with code 0
    And the file `.git/intend/contrib/startup-fix/locks/startup-fix.json` exists
    And `intend trace --mode contrib startup-fix` succeeds

  Scenario: Detect drift after a locked contribution spec changes
    Given a git repository
    And GitHub issue import is available
    And a locked contribution bundle `startup-fix`
    When I replace the contents of `.git/intend/contrib/startup-fix/specs/startup-fix.md`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `contract drift`

  Scenario: Amend an intentional contribution contract change
    Given a git repository
    And GitHub issue import is available
    And a locked contribution bundle `startup-fix`
    When I replace the contents of `.git/intend/contrib/startup-fix/features/startup-fix.feature`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 0
    And the contribution lock version for `startup-fix` is 2
    And `intend trace --mode contrib startup-fix` succeeds
