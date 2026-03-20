Feature: Manage the core contract lifecycle
  In order to keep intent and acceptance explicit
  As a Go developer using `intend`
  I want to initialize a repo, create a bundle, lock it, detect drift, and amend it intentionally

  Scenario: Initialize an owned repository
    Given an empty working directory
    When I run `intend init`
    Then it exits with code 0
    And the directory `specs` exists
    And the directory `features` exists
    And the directory `.intend/trace` exists
    And the directory `.intend/locks` exists

  Scenario: Create a new bundle
    Given an initialized owned repository
    When I run `intend new health-check`
    Then it exits with code 0
    And the file `specs/health-check.md` exists
    And the file `features/health-check.feature` exists
    And the file `.intend/trace/health-check.json` exists

  Scenario: Lock an existing bundle
    Given an initialized owned repository
    And an existing bundle `health-check`
    When I run `intend lock health-check`
    Then it exits with code 0
    And the file `.intend/locks/health-check.json` exists
    And `intend trace health-check` succeeds

  Scenario: Detect drift after a locked spec changes
    Given an initialized owned repository
    And a locked bundle `health-check`
    When I replace the contents of `specs/health-check.md`
    And I run `intend trace health-check`
    Then it exits with code 1
    And stderr contains `contract drift`

  Scenario: Amend an intentional contract change
    Given an initialized owned repository
    And a locked bundle `health-check`
    When I replace the contents of `features/health-check.feature`
    And I run `intend amend health-check`
    Then it exits with code 0
    And the lock version for `health-check` is 2
    And `intend trace health-check` succeeds

  Scenario: Reject missing bundle name
    Given an initialized owned repository
    When I run `intend new`
    Then it exits with code 2
