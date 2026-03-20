Feature: contrib-issue-snapshot-data-validation-matrix
  Scenario: Reject trace for malformed contribution issue snapshot JSON
    Given a locked contribution bundle `startup-fix`
    When I replace the contents of `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `decode contribution issue snapshot`

  Scenario Outline: Reject trace for contribution issue snapshot JSON missing required fields
    Given a locked contribution bundle `startup-fix`
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `<field>` becomes `<value>`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `contribution issue snapshot is missing required fields`

    Examples:
      | field  | value   |
      | number | 0       |
      | title  | <empty> |
      | url    | <empty> |

  Scenario: Reject amend for malformed contribution issue snapshot JSON
    Given a locked contribution bundle `startup-fix`
    When I replace the contents of `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `decode contribution issue snapshot`

  Scenario Outline: Reject amend for contribution issue snapshot JSON missing required fields
    Given a locked contribution bundle `startup-fix`
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `<field>` becomes `<value>`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `contribution issue snapshot is missing required fields`

    Examples:
      | field  | value   |
      | number | 0       |
      | title  | <empty> |
      | url    | <empty> |

  Scenario: Stop verify for malformed contribution issue snapshot JSON
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I replace the contents of `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `decode contribution issue snapshot`
    And the verification log is empty

  Scenario Outline: Stop verify for contribution issue snapshot JSON missing required fields
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `<field>` becomes `<value>`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `contribution issue snapshot is missing required fields`
    And the verification log is empty

    Examples:
      | field  | value   |
      | number | 0       |
      | title  | <empty> |
      | url    | <empty> |

  Scenario: Reject initial lock for malformed contribution issue snapshot JSON
    Given an existing contribution bundle `startup-fix`
    When I replace the contents of `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend lock --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `decode contribution issue snapshot`
    And the path `.git/intend/contrib/startup-fix/locks/startup-fix.json` does not exist

  Scenario Outline: Reject initial lock for contribution issue snapshot JSON missing required fields
    Given an existing contribution bundle `startup-fix`
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `<field>` becomes `<value>`
    And I run `intend lock --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `contribution issue snapshot is missing required fields`
    And the path `.git/intend/contrib/startup-fix/locks/startup-fix.json` does not exist

    Examples:
      | field  | value   |
      | number | 0       |
      | title  | <empty> |
      | url    | <empty> |
