Feature: metadata-identity-validation-matrix
  Scenario: Reject owned trace JSON with a mismatched name during trace
    Given a locked bundle `health-check`
    When I replace the trace file `.intend/trace/health-check.json` so `name` becomes `other-bundle`
    And I run `intend trace health-check`
    Then it exits with code 1
    And stderr contains `trace file name mismatch: expected health-check, got other-bundle`

  Scenario: Reject owned trace JSON with a mismatched mode during trace
    Given a locked bundle `health-check`
    When I replace the trace file `.intend/trace/health-check.json` so `mode` becomes `contrib`
    And I run `intend trace health-check`
    Then it exits with code 1
    And stderr contains `trace file mode mismatch: expected owned, got contrib`

  Scenario: Reject contribution trace JSON with a mismatched name during trace
    Given a locked contribution bundle `startup-fix`
    When I replace the trace file `.git/intend/contrib/startup-fix/trace/startup-fix.json` so `name` becomes `other-bundle`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `trace file name mismatch: expected startup-fix, got other-bundle`

  Scenario: Reject contribution trace JSON with a mismatched mode during trace
    Given a locked contribution bundle `startup-fix`
    When I replace the trace file `.git/intend/contrib/startup-fix/trace/startup-fix.json` so `mode` becomes `owned`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `trace file mode mismatch: expected contrib, got owned`

  Scenario: Reject owned lock JSON with a mismatched name during trace
    Given a locked bundle `health-check`
    When I replace the lock file `.intend/locks/health-check.json` so `name` becomes `other-bundle`
    And I run `intend trace health-check`
    Then it exits with code 1
    And stderr contains `lock file name mismatch: expected health-check, got other-bundle`

  Scenario: Reject contribution lock JSON with a mismatched name during trace
    Given a locked contribution bundle `startup-fix`
    When I replace the lock file `.git/intend/contrib/startup-fix/locks/startup-fix.json` so `name` becomes `other-bundle`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `lock file name mismatch: expected startup-fix, got other-bundle`

  Scenario: Reject amend for owned trace JSON with a mismatched name
    Given a locked bundle `health-check`
    When I replace the trace file `.intend/trace/health-check.json` so `name` becomes `other-bundle`
    And I run `intend amend health-check`
    Then it exits with code 1
    And stderr contains `trace file name mismatch: expected health-check, got other-bundle`

  Scenario: Reject amend for contribution trace JSON with a mismatched mode
    Given a locked contribution bundle `startup-fix`
    When I replace the trace file `.git/intend/contrib/startup-fix/trace/startup-fix.json` so `mode` becomes `owned`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `trace file mode mismatch: expected contrib, got owned`

  Scenario: Reject amend for owned lock JSON with a mismatched name
    Given a locked bundle `health-check`
    When I replace the lock file `.intend/locks/health-check.json` so `name` becomes `other-bundle`
    And I run `intend amend health-check`
    Then it exits with code 1
    And stderr contains `lock file name mismatch: expected health-check, got other-bundle`

  Scenario: Stop verify for owned trace JSON with a mismatched mode
    Given a locked bundle `health-check`
    And verification tools are available
    When I replace the trace file `.intend/trace/health-check.json` so `mode` becomes `contrib`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `trace file mode mismatch: expected owned, got contrib`
    And the verification log is empty

  Scenario: Stop verify for contribution trace JSON with a mismatched name
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I replace the trace file `.git/intend/contrib/startup-fix/trace/startup-fix.json` so `name` becomes `other-bundle`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `trace file name mismatch: expected startup-fix, got other-bundle`
    And the verification log is empty

  Scenario: Stop verify for contribution lock JSON with a mismatched name
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I replace the lock file `.git/intend/contrib/startup-fix/locks/startup-fix.json` so `name` becomes `other-bundle`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `lock file name mismatch: expected startup-fix, got other-bundle`
    And the verification log is empty

  Scenario: Reject initial lock for owned trace JSON with a mismatched name
    Given an existing bundle `health-check`
    When I replace the trace file `.intend/trace/health-check.json` so `name` becomes `other-bundle`
    And I run `intend lock health-check`
    Then it exits with code 1
    And stderr contains `trace file name mismatch: expected health-check, got other-bundle`
    And the path `.intend/locks/health-check.json` does not exist

  Scenario: Reject initial lock for contribution trace JSON with a mismatched mode
    Given an existing contribution bundle `startup-fix`
    When I replace the trace file `.git/intend/contrib/startup-fix/trace/startup-fix.json` so `mode` becomes `owned`
    And I run `intend lock --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `trace file mode mismatch: expected contrib, got owned`
    And the path `.git/intend/contrib/startup-fix/locks/startup-fix.json` does not exist
