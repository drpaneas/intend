Feature: agent-install-behavior
  Scenario: Install Cursor command and skills
    Given an initialized owned repository
    When I run `intend agent install cursor`
    Then it exits with code 0
    And the directory `.cursor/commands` exists
    And the file `.cursor/commands/intend-workflow.md` exists
    And the file `.cursor/skills/intend-idd-workflow/SKILL.md` exists
    And the file `.cursor/skills/intend-go-core/SKILL.md` exists
    And the file `.cursor/skills/intend-go-testing/SKILL.md` exists
    And the file `.cursor/commands/intend-workflow.md` contains `spec -> feature -> tests -> implementation`

  Scenario: Reject unsupported agent in v1
    Given an initialized owned repository
    When I run `intend agent install claude`
    Then it exits with code 1
    And stderr contains `unsupported agent: claude`

  Scenario: Re-running install with unchanged files succeeds
    Given an initialized owned repository
    When I run `intend agent install cursor`
    And I run `intend agent install cursor`
    Then it exits with code 0
    And the file `.cursor/commands/intend-workflow.md` contains `contract-driven in this repository.`

  Scenario: Refuse to overwrite a locally edited managed file
    Given an initialized owned repository
    And I run `intend agent install cursor`
    And I replace the contents of `.cursor/commands/intend-workflow.md`
    When I run `intend agent install cursor`
    Then it exits with code 1
    And stderr contains `managed file was modified: .cursor/commands/intend-workflow.md`
    And the file `.cursor/commands/intend-workflow.md` contains `changed .cursor/commands/intend-workflow.md`
