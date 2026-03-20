Feature: contrib-import-data-validation-matrix
  Scenario: Reject malformed GitHub issue JSON
    Given a git repository
    And GitHub issue import returns malformed JSON
    When I run `intend new --mode contrib --from-gh owner/repo#123 startup-fix`
    Then it exits with code 1
    And stderr contains `gh issue view returned invalid JSON`
    And the path `.git/intend/contrib/startup-fix` does not exist

  Scenario: Reject incomplete GitHub issue data
    Given a git repository
    And GitHub issue import returns incomplete issue data
    When I run `intend new --mode contrib --from-gh owner/repo#123 startup-fix`
    Then it exits with code 1
    And stderr contains `gh issue view returned incomplete issue data`
    And the path `.git/intend/contrib/startup-fix` does not exist
