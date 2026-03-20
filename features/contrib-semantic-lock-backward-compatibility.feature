Feature: contrib-semantic-lock-backward-compatibility
  Scenario: Detect drift for a formatting-only rewrite against an older contribution lock
    Given a locked contribution bundle `startup-fix`
    And I remove the lock file semantic digests from `.git/intend/contrib/startup-fix/locks/startup-fix.json`
    When I rewrite the issue snapshot JSON with pretty formatting at `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `contract drift`

  Scenario: Stop verify for a formatting-only rewrite against an older contribution lock
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    And I remove the lock file semantic digests from `.git/intend/contrib/startup-fix/locks/startup-fix.json`
    When I rewrite the issue snapshot JSON with pretty formatting at `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `contract drift`
    And the verification log is empty

  Scenario: Amend upgrades an older contribution lock after a representation-only rewrite
    Given a locked contribution bundle `startup-fix`
    And I remove the lock file semantic digests from `.git/intend/contrib/startup-fix/locks/startup-fix.json`
    When I rewrite the issue snapshot JSON with sorted keys at `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 0
    And the contribution lock version for `startup-fix` is 2

  Scenario: Ignore later representation-only rewrites after an older contribution lock is upgraded
    Given a locked contribution bundle `startup-fix`
    And I remove the lock file semantic digests from `.git/intend/contrib/startup-fix/locks/startup-fix.json`
    And I rewrite the issue snapshot JSON with sorted keys at `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend amend --mode contrib startup-fix`
    When I rewrite the issue snapshot JSON with pretty formatting at `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 0
    And the contribution lock version for `startup-fix` is 2
