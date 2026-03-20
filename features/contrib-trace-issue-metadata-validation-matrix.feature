Feature: contrib-trace-issue-metadata-validation-matrix
  Scenario Outline: Reject trace for invalid contribution trace issue metadata
    Given a locked contribution bundle `startup-fix`
    When I replace the trace file `.git/intend/contrib/startup-fix/trace/startup-fix.json` so `<field>` becomes `<value>`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `<message>`

    Examples:
      | field     | value               | message                                                                            |
      | issueRef  | <empty>             | contribution trace is missing required issue metadata                              |
      | issuePath | <empty>             | contribution trace is missing required issue metadata                              |
      | issueRef  | not-a-github-ref    | contribution trace issueRef is invalid: not-a-github-ref                           |
      | issueRef  | owner/repo#0        | contribution trace issueRef is invalid: owner/repo#0                               |
      | issueRef  | owner/repo#-1       | contribution trace issueRef is invalid: owner/repo#-1                              |
      | issuePath | specs/startup-fix.md | contribution trace issuePath mismatch: expected issue.json, got specs/startup-fix.md |

  Scenario Outline: Reject amend for invalid contribution trace issue metadata
    Given a locked contribution bundle `startup-fix`
    When I replace the trace file `.git/intend/contrib/startup-fix/trace/startup-fix.json` so `<field>` becomes `<value>`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `<message>`

    Examples:
      | field     | value               | message                                                                            |
      | issueRef  | <empty>             | contribution trace is missing required issue metadata                              |
      | issuePath | <empty>             | contribution trace is missing required issue metadata                              |
      | issueRef  | not-a-github-ref    | contribution trace issueRef is invalid: not-a-github-ref                           |
      | issueRef  | owner/repo#0        | contribution trace issueRef is invalid: owner/repo#0                               |
      | issueRef  | owner/repo#-1       | contribution trace issueRef is invalid: owner/repo#-1                              |
      | issuePath | specs/startup-fix.md | contribution trace issuePath mismatch: expected issue.json, got specs/startup-fix.md |

  Scenario Outline: Stop verify for invalid contribution trace issue metadata
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I replace the trace file `.git/intend/contrib/startup-fix/trace/startup-fix.json` so `<field>` becomes `<value>`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `<message>`
    And the verification log is empty

    Examples:
      | field     | value               | message                                                                            |
      | issueRef  | <empty>             | contribution trace is missing required issue metadata                              |
      | issuePath | <empty>             | contribution trace is missing required issue metadata                              |
      | issueRef  | not-a-github-ref    | contribution trace issueRef is invalid: not-a-github-ref                           |
      | issueRef  | owner/repo#0        | contribution trace issueRef is invalid: owner/repo#0                               |
      | issueRef  | owner/repo#-1       | contribution trace issueRef is invalid: owner/repo#-1                              |
      | issuePath | specs/startup-fix.md | contribution trace issuePath mismatch: expected issue.json, got specs/startup-fix.md |

  Scenario Outline: Reject initial lock for invalid contribution trace issue metadata
    Given an existing contribution bundle `startup-fix`
    When I replace the trace file `.git/intend/contrib/startup-fix/trace/startup-fix.json` so `<field>` becomes `<value>`
    And I run `intend lock --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `<message>`
    And the path `.git/intend/contrib/startup-fix/locks/startup-fix.json` does not exist

    Examples:
      | field     | value               | message                                                                            |
      | issueRef  | <empty>             | contribution trace is missing required issue metadata                              |
      | issuePath | <empty>             | contribution trace is missing required issue metadata                              |
      | issueRef  | not-a-github-ref    | contribution trace issueRef is invalid: not-a-github-ref                           |
      | issueRef  | owner/repo#0        | contribution trace issueRef is invalid: owner/repo#0                               |
      | issueRef  | owner/repo#-1       | contribution trace issueRef is invalid: owner/repo#-1                              |
      | issuePath | specs/startup-fix.md | contribution trace issuePath mismatch: expected issue.json, got specs/startup-fix.md |
