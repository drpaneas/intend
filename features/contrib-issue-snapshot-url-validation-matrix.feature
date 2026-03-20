Feature: contrib-issue-snapshot-url-validation-matrix
  Scenario Outline: Reject trace for contribution issue snapshot URL <kind>
    Given a locked contribution bundle `startup-fix`
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `url` becomes `<url>`
    And I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `<error>`

    Examples:
      | kind                             | url                                                 | error                                                                                             |
      | on a non-GitHub host            | https://gitlab.com/owner/repo/-/issues/123          | contribution issue snapshot URL is not a GitHub issue URL: https://gitlab.com/owner/repo/-/issues/123 |
      | on an alternate GitHub hostname | https://www.github.com/owner/repo/issues/123        | contribution issue snapshot URL uses unsupported GitHub hostname: https://www.github.com/owner/repo/issues/123 |
      | with a non-issue GitHub path    | https://github.com/owner/repo/pull/123              | contribution issue snapshot URL is invalid: https://github.com/owner/repo/pull/123              |
      | with a query string             | https://github.com/owner/repo/issues/123?tab=comments | contribution issue snapshot URL is non-canonical: https://github.com/owner/repo/issues/123?tab=comments |
      | with a fragment                 | https://github.com/owner/repo/issues/123#comment-1  | contribution issue snapshot URL is non-canonical: https://github.com/owner/repo/issues/123#comment-1 |

  Scenario Outline: Reject amend for contribution issue snapshot URL <kind>
    Given a locked contribution bundle `startup-fix`
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `url` becomes `<url>`
    And I run `intend amend --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `<error>`

    Examples:
      | kind                             | url                                                 | error                                                                                             |
      | on a non-GitHub host            | https://gitlab.com/owner/repo/-/issues/123          | contribution issue snapshot URL is not a GitHub issue URL: https://gitlab.com/owner/repo/-/issues/123 |
      | on an alternate GitHub hostname | https://www.github.com/owner/repo/issues/123        | contribution issue snapshot URL uses unsupported GitHub hostname: https://www.github.com/owner/repo/issues/123 |
      | with a non-issue GitHub path    | https://github.com/owner/repo/pull/123              | contribution issue snapshot URL is invalid: https://github.com/owner/repo/pull/123              |
      | with a query string             | https://github.com/owner/repo/issues/123?tab=comments | contribution issue snapshot URL is non-canonical: https://github.com/owner/repo/issues/123?tab=comments |
      | with a fragment                 | https://github.com/owner/repo/issues/123#comment-1  | contribution issue snapshot URL is non-canonical: https://github.com/owner/repo/issues/123#comment-1 |

  Scenario Outline: Stop verify for contribution issue snapshot URL <kind>
    Given a locked contribution bundle `startup-fix`
    And verification tools are available
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `url` becomes `<url>`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `<error>`
    And the verification log is empty

    Examples:
      | kind                             | url                                                 | error                                                                                             |
      | on a non-GitHub host            | https://gitlab.com/owner/repo/-/issues/123          | contribution issue snapshot URL is not a GitHub issue URL: https://gitlab.com/owner/repo/-/issues/123 |
      | on an alternate GitHub hostname | https://www.github.com/owner/repo/issues/123        | contribution issue snapshot URL uses unsupported GitHub hostname: https://www.github.com/owner/repo/issues/123 |
      | with a non-issue GitHub path    | https://github.com/owner/repo/pull/123              | contribution issue snapshot URL is invalid: https://github.com/owner/repo/pull/123              |
      | with a query string             | https://github.com/owner/repo/issues/123?tab=comments | contribution issue snapshot URL is non-canonical: https://github.com/owner/repo/issues/123?tab=comments |
      | with a fragment                 | https://github.com/owner/repo/issues/123#comment-1  | contribution issue snapshot URL is non-canonical: https://github.com/owner/repo/issues/123#comment-1 |

  Scenario Outline: Reject initial lock for contribution issue snapshot URL <kind>
    Given an existing contribution bundle `startup-fix`
    When I replace the issue snapshot file `.git/intend/contrib/startup-fix/issue.json` so `url` becomes `<url>`
    And I run `intend lock --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `<error>`
    And the path `.git/intend/contrib/startup-fix/locks/startup-fix.json` does not exist

    Examples:
      | kind                             | url                                                 | error                                                                                             |
      | on a non-GitHub host            | https://gitlab.com/owner/repo/-/issues/123          | contribution issue snapshot URL is not a GitHub issue URL: https://gitlab.com/owner/repo/-/issues/123 |
      | on an alternate GitHub hostname | https://www.github.com/owner/repo/issues/123        | contribution issue snapshot URL uses unsupported GitHub hostname: https://www.github.com/owner/repo/issues/123 |
      | with a non-issue GitHub path    | https://github.com/owner/repo/pull/123              | contribution issue snapshot URL is invalid: https://github.com/owner/repo/pull/123              |
      | with a query string             | https://github.com/owner/repo/issues/123?tab=comments | contribution issue snapshot URL is non-canonical: https://github.com/owner/repo/issues/123?tab=comments |
      | with a fragment                 | https://github.com/owner/repo/issues/123#comment-1  | contribution issue snapshot URL is non-canonical: https://github.com/owner/repo/issues/123#comment-1 |
