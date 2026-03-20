Feature: Require locks before trace and verify
  In order to distinguish approved contracts from drafts
  As a developer using `intend`
  I want unlocked bundles to fail `trace` and stop `verify`

  Scenario: Reject trace for an unlocked owned bundle
    Given an initialized owned repository
    And an existing bundle `health-check`
    When I run `intend trace health-check`
    Then it exits with code 1
    And stderr contains `contract is not locked`

  Scenario: Reject trace for an unlocked contribution bundle
    Given a git repository
    And GitHub issue import is available
    And an existing contribution bundle `startup-fix`
    When I run `intend trace --mode contrib startup-fix`
    Then it exits with code 1
    And stderr contains `contract is not locked`

  Scenario: Stop verify when an owned bundle is unlocked
    Given an initialized owned repository
    And an existing bundle `health-check`
    And verification tools are available
    When I run `intend verify`
    Then it exits with code 1
    And stderr contains `contract is not locked`
    And the verification log is empty

  Scenario: Stop verify when a contribution bundle is unlocked
    Given a git repository
    And GitHub issue import is available
    And an existing contribution bundle `startup-fix`
    And verification tools are available
    When I run `intend verify`
    Then it exits with code 1
    And stderr contains `contract is not locked`
    And the verification log is empty
