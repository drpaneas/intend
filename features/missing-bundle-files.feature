Feature: missing-bundle-files
  Scenario: Reject amend when an owned spec file is missing
    Given a locked bundle `health-check`
    When I remove the path `specs/health-check.md`
    And I run `intend amend health-check`
    Then it exits with code 1
    And stderr contains `amend: digest specs/health-check.md`

  Scenario: Reject amend when a contribution issue snapshot is missing
    Given a locked contribution bundle `startup-fix`
    When I remove the path `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `amend: digest issue.json`

  Scenario: Stop verify when an owned feature file is missing
    Given a locked bundle `health-check`
    And verification tools are available
    When I remove the path `features/health-check.feature`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `verify: digest features/health-check.feature`
    And the verification log is empty

  Scenario: Stop verify when a contribution issue snapshot is missing
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I remove the path `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `verify: digest issue.json`
    And the verification log is empty
