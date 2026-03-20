Feature: contrib-issue-snapshot-edit-behavior
  Scenario: Detect drift after a locked contribution issue snapshot title changes
    Given a locked contribution bundle `startup-fix`
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `title` becomes `Fix startup crash on Linux`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `contract drift`

  Scenario: Amend an intentional contribution issue snapshot body change
    Given a locked contribution bundle `startup-fix`
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `body` becomes `The app crashes during startup on Linux only.`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 0
    And the contribution lock version for `startup-fix` is 2
    And `intend trace --mode contrib startup-fix` succeeds

  Scenario: Stop verify after a locked contribution issue snapshot title changes
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `title` becomes `Fix startup crash on Linux`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `contract drift`
    And the verification log is empty

  Scenario: Lock a contribution bundle after an issue snapshot title change
    Given an existing contribution bundle `startup-fix`
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `title` becomes `Fix startup crash on Linux`
    And I run `intend lock --mode contrib startup-fix`
    Then it exits with code 0
    And the file `.git/intend/contrib/startup-fix/locks/startup-fix.json` exists
    And `intend trace --mode contrib startup-fix` succeeds

  Scenario: Lock a contribution bundle after an issue snapshot body change
    Given an existing contribution bundle `startup-fix`
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `body` becomes `The app crashes during startup on Linux only.`
    And I run `intend lock --mode contrib startup-fix`
    Then it exits with code 0
    And the file `.git/intend/contrib/startup-fix/locks/startup-fix.json` exists
    And `intend trace --mode contrib startup-fix` succeeds

  Scenario: Detect drift after a locked contribution issue snapshot gains an extra field
    Given a locked contribution bundle `startup-fix`
    When I add the unknown issue snapshot field `platform` with value `linux` to `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `contract drift`

  Scenario: Amend an intentional contribution issue snapshot extra field
    Given a locked contribution bundle `startup-fix`
    When I add the unknown issue snapshot field `platform` with value `linux` to `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 0
    And the contribution lock version for `startup-fix` is 2
    And `intend trace --mode contrib startup-fix` succeeds

  Scenario: Stop verify after a locked contribution issue snapshot gains an extra field
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I add the unknown issue snapshot field `platform` with value `linux` to `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `contract drift`
    And the verification log is empty

  Scenario: Lock a contribution bundle after an issue snapshot gains an extra field
    Given an existing contribution bundle `startup-fix`
    When I add the unknown issue snapshot field `platform` with value `linux` to `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend lock --mode contrib startup-fix`
    Then it exits with code 0
    And the file `.git/intend/contrib/startup-fix/locks/startup-fix.json` exists
    And `intend trace --mode contrib startup-fix` succeeds

  Scenario: Ignore formatting-only rewrites during trace
    Given a locked contribution bundle `startup-fix`
    When I rewrite the issue snapshot JSON with pretty formatting at `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 0

  Scenario: Keep the same lock version after a representation-only amend
    Given a locked contribution bundle `startup-fix`
    When I rewrite the issue snapshot JSON with sorted keys at `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 0
    And stdout contains `contract unchanged`
    And the contribution lock version for `startup-fix` is 1
    And `intend trace --mode contrib startup-fix` succeeds

  Scenario: Continue verify after a formatting-only rewrite
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I rewrite the issue snapshot JSON with pretty formatting at `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend verify`
    Then it exits with code 0
    And the verification log contains lines in order:
      | go test ./...            |
      | golangci-lint run       |
      | trufflehog filesystem . |
      | gitleaks dir .          |
      | trivy fs .              |

  Scenario: Lock a contribution bundle after an issue snapshot key-order rewrite
    Given an existing contribution bundle `startup-fix`
    When I rewrite the issue snapshot JSON with sorted keys at `.git/intend/contrib/startup-fix/issue.json`
    And I run `intend lock --mode contrib startup-fix`
    Then it exits with code 0
    And the file `.git/intend/contrib/startup-fix/locks/startup-fix.json` exists
    And `intend trace --mode contrib startup-fix` succeeds
