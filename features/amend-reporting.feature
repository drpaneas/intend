Feature: amend-reporting
  Scenario: Report a no-op owned amend
    Given an initialized owned repository
    And a locked bundle `health-check`
    When I run `intend amend health-check`
    Then it exits with code 0
    And stdout contains `contract unchanged`
    And stdout does not contain `upgraded semantic lock metadata`
    And the lock version for `health-check` is 1
    And `intend trace health-check` succeeds

  Scenario: Report an owned amend after an intentional contract change
    Given an initialized owned repository
    And a locked bundle `health-check`
    When I replace the contents of `features/health-check.feature`
    And I run `intend amend health-check`
    Then it exits with code 0
    And stdout contains `amended health-check to version 2`
    And stdout does not contain `upgraded semantic lock metadata`
    And the lock version for `health-check` is 2
    And `intend trace health-check` succeeds

  Scenario: Report a no-op current contribution amend
    Given a locked contribution bundle `startup-fix`
    When I run `intend amend --mode contrib startup-fix`
    Then it exits with code 0
    And stdout contains `contract unchanged`
    And stdout does not contain `upgraded semantic lock metadata`
    And the contribution lock version for `startup-fix` is 1
    And `intend trace --mode contrib startup-fix` succeeds

  Scenario: Report an ordinary current contribution amend
    Given a locked contribution bundle `startup-fix`
    When I replace the contents of `.git/intend/contrib/startup-fix/features/startup-fix.feature`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 0
    And stdout contains `amended startup-fix to version 2`
    And stdout does not contain `upgraded semantic lock metadata`
    And the contribution lock version for `startup-fix` is 2
    And `intend trace --mode contrib startup-fix` succeeds

  Scenario: Report a contribution semantic lock upgrade amend
    Given a locked contribution bundle `startup-fix`
    And I remove the lock file semantic digests from `.git/intend/contrib/startup-fix/locks/startup-fix.json`
    When I replace the contents of `.git/intend/contrib/startup-fix/features/startup-fix.feature`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 0
    And stdout contains `amended startup-fix to version 2`
    And stdout contains `upgraded semantic lock metadata`
    And the contribution lock version for `startup-fix` is 2
    And `intend trace --mode contrib startup-fix` succeeds
