Feature: contrib-import-identity-validation-matrix
  Scenario: Reject a returned issue snapshot with the wrong number
    Given a git repository
    And GitHub issue import returns a different issue number
    When I run `intend new --mode contrib --from-gh owner/repo#123 startup-fix`
    Then it exits with code 1
    And stderr contains `gh issue view returned issue #124, expected #123`
    And the path `.git/intend/contrib/startup-fix` does not exist

  Scenario: Reject a returned issue URL for a different repository
    Given a git repository
    And GitHub issue import returns a URL for a different repository
    When I run `intend new --mode contrib --from-gh owner/repo#123 startup-fix`
    Then it exits with code 1
    And stderr contains `gh issue view returned URL for repo other/repo, expected owner/repo`
    And the path `.git/intend/contrib/startup-fix` does not exist

  Scenario: Reject a returned GitHub issue URL with a different issue number
    Given a git repository
    And GitHub issue import returns a GitHub issue URL with a different issue number
    When I run `intend new --mode contrib --from-gh owner/repo#123 startup-fix`
    Then it exits with code 1
    And stderr contains `gh issue view returned URL issue #124, expected #123`
    And the path `.git/intend/contrib/startup-fix` does not exist
