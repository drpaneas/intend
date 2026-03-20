Feature: contrib-import-url-validation-matrix
  Scenario: Reject a returned issue URL on a non-GitHub host
    Given a git repository
    And GitHub issue import returns a non-GitHub issue URL
    When I run `intend new --mode contrib --from-gh owner/repo#123 startup-fix`
    Then it exits with code 1
    And stderr contains `gh issue view returned non-GitHub issue URL: https://gitlab.com/owner/repo/-/issues/123`
    And the path `.git/intend/contrib/startup-fix` does not exist

  Scenario: Reject a returned GitHub issue URL with an alternate hostname
    Given a git repository
    And GitHub issue import returns a GitHub issue URL with an alternate hostname
    When I run `intend new --mode contrib --from-gh owner/repo#123 startup-fix`
    Then it exits with code 1
    And stderr contains `gh issue view returned unsupported GitHub hostname: www.github.com`
    And the path `.git/intend/contrib/startup-fix` does not exist

  Scenario: Reject a returned GitHub URL that is not an issue URL
    Given a git repository
    And GitHub issue import returns a GitHub pull request URL
    When I run `intend new --mode contrib --from-gh owner/repo#123 startup-fix`
    Then it exits with code 1
    And stderr contains `gh issue view returned invalid GitHub issue URL: https://github.com/owner/repo/pull/123`
    And the path `.git/intend/contrib/startup-fix` does not exist

  Scenario: Reject a returned GitHub issue URL with a query string
    Given a git repository
    And GitHub issue import returns a GitHub issue URL with a query string
    When I run `intend new --mode contrib --from-gh owner/repo#123 startup-fix`
    Then it exits with code 1
    And stderr contains `gh issue view returned non-canonical GitHub issue URL: https://github.com/owner/repo/issues/123?tab=comments`
    And the path `.git/intend/contrib/startup-fix` does not exist

  Scenario: Reject a returned GitHub issue URL with a fragment
    Given a git repository
    And GitHub issue import returns a GitHub issue URL with a fragment
    When I run `intend new --mode contrib --from-gh owner/repo#123 startup-fix`
    Then it exits with code 1
    And stderr contains `gh issue view returned non-canonical GitHub issue URL: https://github.com/owner/repo/issues/123#issuecomment-1`
    And the path `.git/intend/contrib/startup-fix` does not exist
