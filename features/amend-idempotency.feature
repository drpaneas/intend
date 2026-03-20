Feature: Keep amend idempotent
  In order to record only real contract changes
  As a developer using `intend`
  I want unchanged contracts to produce a no-op amend result

  Scenario: Do not increment the version for an unchanged owned bundle
    Given an initialized owned repository
    And a locked bundle `health-check`
    When I run `intend amend health-check`
    Then it exits with code 0
    And stdout contains `contract unchanged`
    And the lock version for `health-check` is 1
    And `intend trace health-check` succeeds

  Scenario: Do not increment the version for an unchanged contribution bundle
    Given a git repository
    And GitHub issue import is available
    And a locked contribution bundle `startup-fix`
    When I run `intend amend --mode contrib startup-fix`
    Then it exits with code 0
    And stdout contains `contract unchanged`
    And the contribution lock version for `startup-fix` is 1
    And `intend trace --mode contrib startup-fix` succeeds
