Feature: contrib-issue-snapshot-identity-validation-matrix
  Scenario Outline: Reject trace for contribution issue snapshot identity mismatch <kind>
    Given a locked contribution bundle `startup-fix`
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `<field>` becomes `<value>`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `<error>`

    Examples:
      | kind                    | field  | value                                     | error                                                                     |
      | with a different number | number | 124                                       | contribution issue snapshot number mismatch: expected 123, got 124        |
      | with a different repo   | url    | https://github.com/other/repo/issues/123  | contribution issue snapshot repo mismatch: expected owner/repo, got other/repo |

  Scenario Outline: Reject amend for contribution issue snapshot identity mismatch <kind>
    Given a locked contribution bundle `startup-fix`
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `<field>` becomes `<value>`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `<error>`

    Examples:
      | kind                    | field  | value                                     | error                                                                     |
      | with a different number | number | 124                                       | contribution issue snapshot number mismatch: expected 123, got 124        |
      | with a different repo   | url    | https://github.com/other/repo/issues/123  | contribution issue snapshot repo mismatch: expected owner/repo, got other/repo |

  Scenario Outline: Stop verify for contribution issue snapshot identity mismatch <kind>
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `<field>` becomes `<value>`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `<error>`
    And the verification log is empty

    Examples:
      | kind                    | field  | value                                     | error                                                                     |
      | with a different number | number | 124                                       | contribution issue snapshot number mismatch: expected 123, got 124        |
      | with a different repo   | url    | https://github.com/other/repo/issues/123  | contribution issue snapshot repo mismatch: expected owner/repo, got other/repo |

  Scenario Outline: Reject initial lock for contribution issue snapshot identity mismatch <kind>
    Given an existing contribution bundle `startup-fix`
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `<field>` becomes `<value>`
    And I run `intend lock --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `<error>`
    And the path `.git/intend/contrib/startup-fix/locks/startup-fix.json` does not exist

    Examples:
      | kind                    | field  | value                                     | error                                                                     |
      | with a different number | number | 124                                       | contribution issue snapshot number mismatch: expected 123, got 124        |
      | with a different repo   | url    | https://github.com/other/repo/issues/123  | contribution issue snapshot repo mismatch: expected owner/repo, got other/repo |
