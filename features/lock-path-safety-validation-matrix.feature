Feature: lock-path-safety-validation-matrix
  Scenario Outline: Reject trace for owned lock metadata path escapes
    Given a locked bundle `health-check`
    When I replace the lock file `.intend/locks/health-check.json` so tracked path `<trackedPath>` becomes `<value>`
    And I run `intend trace health-check`
    Then it exits with code 1
    And stderr contains `<message>`

    Examples:
      | trackedPath               | value                       | message                                               |
      | specs/health-check.md     | /tmp/intend-outside-spec.md | lock file path must be relative to the workspace root |
      | features/health-check.feature | ../outside.feature      | lock file path escapes the workspace root             |

  Scenario Outline: Reject trace for contribution lock metadata path escapes
    Given a locked contribution bundle `startup-fix`
    When I replace the lock file `.git/intend/contrib/startup-fix/locks/startup-fix.json` so tracked path `issue.json` becomes `<value>`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `<message>`

    Examples:
      | value                        | message                                            |
      | /tmp/intend-outside-issue.json | lock file path must be relative to the bundle root |
      | ../outside-issue.json        | lock file path escapes the bundle root            |

  Scenario: Reject trace for an owned lock-tracked trace file symlink outside the workspace root
    Given a locked bundle `health-check`
    When I replace the path `.intend/trace/health-check.json` with a symlink to an external copy preserving its contents
    And I run `intend trace health-check`
    Then it exits with code 1
    And stderr contains `lock file path resolves through a symlink outside the workspace root`

  Scenario: Reject trace for a contribution lock-tracked trace file symlink outside the bundle root
    Given a locked contribution bundle `startup-fix`
    When I replace the path `.git/intend/contrib/startup-fix/trace/startup-fix.json` with a symlink to an external copy preserving its contents
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `lock file path resolves through a symlink outside the bundle root`

  Scenario: Reject amend for an owned lock JSON with an absolute tracked path
    Given a locked bundle `health-check`
    When I replace the lock file `.intend/locks/health-check.json` so tracked path `specs/health-check.md` becomes `/tmp/intend-outside-spec.md`
    And I run `intend amend health-check`
    Then it exits with code 1
    And stderr contains `lock file path must be relative to the workspace root`

  Scenario: Reject amend for a contribution lock JSON with a tracked path that escapes the bundle root
    Given a locked contribution bundle `startup-fix`
    When I replace the lock file `.git/intend/contrib/startup-fix/locks/startup-fix.json` so tracked path `issue.json` becomes `../outside-issue.json`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `lock file path escapes the bundle root`

  Scenario: Reject amend for an owned lock-tracked trace file symlink outside the workspace root
    Given a locked bundle `health-check`
    When I replace the path `.intend/trace/health-check.json` with a symlink to an external copy preserving its contents
    And I run `intend amend health-check`
    Then it exits with code 1
    And stderr contains `lock file path resolves through a symlink outside the workspace root`

  Scenario: Reject amend for a contribution lock-tracked trace file symlink outside the bundle root
    Given a locked contribution bundle `startup-fix`
    When I replace the path `.git/intend/contrib/startup-fix/trace/startup-fix.json` with a symlink to an external copy preserving its contents
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `lock file path resolves through a symlink outside the bundle root`

  Scenario: Stop verify for an owned lock JSON with a tracked path that escapes the workspace root
    Given a locked bundle `health-check`
    And verification tools are available
    When I replace the lock file `.intend/locks/health-check.json` so tracked path `features/health-check.feature` becomes `../outside.feature`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `lock file path escapes the workspace root`
    And the verification log is empty

  Scenario: Stop verify for a contribution lock JSON with an absolute tracked path
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I replace the lock file `.git/intend/contrib/startup-fix/locks/startup-fix.json` so tracked path `issue.json` becomes `/tmp/intend-outside-issue.json`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `lock file path must be relative to the bundle root`
    And the verification log is empty

  Scenario: Stop verify for an owned lock-tracked trace file symlink outside the workspace root
    Given a locked bundle `health-check`
    And verification tools are available
    When I replace the path `.intend/trace/health-check.json` with a symlink to an external copy preserving its contents
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `lock file path resolves through a symlink outside the workspace root`
    And the verification log is empty

  Scenario: Stop verify for a contribution lock-tracked trace file symlink outside the bundle root
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I replace the path `.git/intend/contrib/startup-fix/trace/startup-fix.json` with a symlink to an external copy preserving its contents
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `lock file path resolves through a symlink outside the bundle root`
    And the verification log is empty
