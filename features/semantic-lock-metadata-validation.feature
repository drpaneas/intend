Feature: semantic-lock-metadata-validation
  Scenario: Reject trace for owned lock metadata with semantic digests
    Given a locked bundle `health-check`
    When I add the lock file semantic digest for `specs/health-check.md` to `.intend/locks/health-check.json`
    And I run `intend trace health-check`
    Then it exits with code 1
    And stderr contains `lock file semantic path is unsupported for owned bundle: specs/health-check.md`

  Scenario: Reject amend for contribution lock metadata with a non-issue semantic path
    Given a locked contribution bundle `startup-fix`
    When I add the lock file semantic digest for `specs/startup-fix.md` to `.git/intend/contrib/startup-fix/locks/startup-fix.json`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `lock file semantic path mismatch: expected issue.json, got specs/startup-fix.md`

  Scenario: Stop verify for contribution lock metadata with a non-issue semantic path
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I add the lock file semantic digest for `features/startup-fix.feature` to `.git/intend/contrib/startup-fix/locks/startup-fix.json`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `lock file semantic path mismatch: expected issue.json, got features/startup-fix.feature`
    And the verification log is empty
