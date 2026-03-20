Feature: trace-path-safety-validation-matrix
  Scenario Outline: Reject trace for owned trace metadata path escapes
    Given a locked bundle `health-check`
    When I replace the trace file `.intend/trace/health-check.json` so `<field>` becomes `<value>`
    And I run `intend trace health-check`
    Then it exits with code 1
    And stderr contains `<message>`

    Examples:
      | field       | value                       | message                                                |
      | specPath    | /tmp/intend-outside-spec.md | trace specPath must be relative to the workspace root  |
      | featurePath | ../outside.feature          | trace featurePath escapes the workspace root           |

  Scenario Outline: Reject trace for contribution trace metadata path escapes
    Given a locked contribution bundle `startup-fix`
    When I replace the trace file `.git/intend/contrib/startup-fix/trace/startup-fix.json` so `<field>` becomes `<value>`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `<message>`

    Examples:
      | field     | value                        | message                                             |
      | issuePath | /tmp/intend-outside-issue.json | trace issuePath must be relative to the bundle root |
      | issuePath | ../outside-issue.json        | trace issuePath escapes the bundle root             |

  Scenario: Reject trace for an owned spec symlink that resolves outside the workspace root
    Given a locked bundle `health-check`
    When I replace the path `specs/health-check.md` with a symlink to an external copy preserving its contents
    And I run `intend trace health-check`
    Then it exits with code 1
    And stderr contains `trace specPath resolves through a symlink outside the workspace root`

  Scenario: Reject trace for a contribution issue symlink that resolves outside the bundle root
    Given a locked contribution bundle `startup-fix`
    When I replace the path `.git/intend/contrib/startup-fix/issue.json` with a symlink to an external copy preserving its contents
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `trace issuePath resolves through a symlink outside the bundle root`

  Scenario: Reject amend for an owned trace JSON with an absolute spec path
    Given a locked bundle `health-check`
    When I replace the trace file `.intend/trace/health-check.json` so `specPath` becomes `/tmp/intend-outside-spec.md`
    And I run `intend amend health-check`
    Then it exits with code 1
    And stderr contains `trace specPath must be relative to the workspace root`

  Scenario: Reject amend for a contribution trace JSON with an issue path that escapes the bundle root
    Given a locked contribution bundle `startup-fix`
    When I replace the trace file `.git/intend/contrib/startup-fix/trace/startup-fix.json` so `issuePath` becomes `../outside-issue.json`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `trace issuePath escapes the bundle root`

  Scenario: Reject amend for an owned spec symlink that resolves outside the workspace root
    Given a locked bundle `health-check`
    When I replace the path `specs/health-check.md` with a symlink to an external copy preserving its contents
    And I run `intend amend health-check`
    Then it exits with code 1
    And stderr contains `lock file path resolves through a symlink outside the workspace root`

  Scenario: Reject amend for a contribution issue symlink that resolves outside the bundle root
    Given a locked contribution bundle `startup-fix`
    When I replace the path `.git/intend/contrib/startup-fix/issue.json` with a symlink to an external copy preserving its contents
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `lock file path resolves through a symlink outside the bundle root`

  Scenario: Stop verify for an owned trace JSON with a feature path that escapes the workspace root
    Given a locked bundle `health-check`
    And verification tools are available
    When I replace the trace file `.intend/trace/health-check.json` so `featurePath` becomes `../outside.feature`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `trace featurePath escapes the workspace root`
    And the verification log is empty

  Scenario: Stop verify for a contribution trace JSON with an absolute issue path
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I replace the trace file `.git/intend/contrib/startup-fix/trace/startup-fix.json` so `issuePath` becomes `/tmp/intend-outside-issue.json`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `trace issuePath must be relative to the bundle root`
    And the verification log is empty

  Scenario: Stop verify for an owned feature symlink that resolves outside the workspace root
    Given a locked bundle `health-check`
    And verification tools are available
    When I replace the path `features/health-check.feature` with a symlink to an external copy preserving its contents
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `trace featurePath resolves through a symlink outside the workspace root`
    And the verification log is empty

  Scenario: Stop verify for a contribution issue symlink that resolves outside the bundle root
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I replace the path `.git/intend/contrib/startup-fix/issue.json` with a symlink to an external copy preserving its contents
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `trace issuePath resolves through a symlink outside the bundle root`
    And the verification log is empty

  Scenario: Reject initial lock for an owned trace JSON with an absolute spec path
    Given an existing bundle `health-check`
    When I replace the trace file `.intend/trace/health-check.json` so `specPath` becomes `/tmp/intend-outside-spec.md`
    And I run `intend lock health-check`
    Then it exits with code 1
    And stderr contains `trace specPath must be relative to the workspace root`
    And the path `.intend/locks/health-check.json` does not exist

  Scenario: Reject initial lock for a contribution trace JSON with an issue path that escapes the bundle root
    Given an existing contribution bundle `startup-fix`
    When I replace the trace file `.git/intend/contrib/startup-fix/trace/startup-fix.json` so `issuePath` becomes `../outside-issue.json`
    And I run `intend lock --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `trace issuePath escapes the bundle root`
    And the path `.git/intend/contrib/startup-fix/locks/startup-fix.json` does not exist

  Scenario: Reject initial lock for an owned spec symlink that resolves outside the workspace root
    Given an existing bundle `health-check`
    When I replace the path `specs/health-check.md` with a symlink to an external copy preserving its contents
    And I run `intend lock health-check`
    Then it exits with code 1
    And stderr contains `trace specPath resolves through a symlink outside the workspace root`
    And the path `.intend/locks/health-check.json` does not exist

  Scenario: Reject initial lock for a contribution issue symlink that resolves outside the bundle root
    Given an existing contribution bundle `startup-fix`
    When I replace the path `.git/intend/contrib/startup-fix/issue.json` with a symlink to an external copy preserving its contents
    And I run `intend lock --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `trace issuePath resolves through a symlink outside the bundle root`
    And the path `.git/intend/contrib/startup-fix/locks/startup-fix.json` does not exist
