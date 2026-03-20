Feature: lock-metadata-validation-matrix
  Scenario Outline: Reject trace for corrupted owned lock JSON
    Given a locked bundle `<name>`
    When I replace the contents of `<path>`
    And I run `intend trace <name>`
    Then it exits with code 1
    And stderr contains `decode lock file for <name>`

    Examples:
      | name         | path                            |
      | health-check | .intend/locks/health-check.json |

  Scenario Outline: Reject trace for corrupted contribution lock JSON
    Given a locked contribution bundle `<name>`
    When I replace the contents of `<path>`
    And I run `intend trace --mode contrib <name>`
    Then it exits with code 1
    And stderr contains `decode lock file for <name>`

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/locks/startup-fix.json |

  Scenario Outline: Reject trace for owned lock JSON missing required fields
    Given a locked bundle `<name>`
    When I replace the lock file `<path>` with valid JSON missing required fields
    And I run `intend trace <name>`
    Then it exits with code 1
    And stderr contains `lock file is missing required fields`

    Examples:
      | name         | path                            |
      | health-check | .intend/locks/health-check.json |

  Scenario Outline: Reject trace for contribution lock JSON missing required fields
    Given a locked contribution bundle `<name>`
    When I replace the lock file `<path>` with valid JSON missing required fields
    And I run `intend trace --mode contrib <name>`
    Then it exits with code 1
    And stderr contains `lock file is missing required fields`

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/locks/startup-fix.json |

  Scenario Outline: Reject amend for corrupted owned lock JSON
    Given a locked bundle `<name>`
    When I replace the contents of `<path>`
    And I run `intend amend <name>`
    Then it exits with code 1
    And stderr contains `decode lock file for <name>`

    Examples:
      | name         | path                            |
      | health-check | .intend/locks/health-check.json |

  Scenario Outline: Reject amend for corrupted contribution lock JSON
    Given a locked contribution bundle `<name>`
    When I replace the contents of `<path>`
    And I run `intend amend --mode contrib <name>`
    Then it exits with code 1
    And stderr contains `decode lock file for <name>`

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/locks/startup-fix.json |

  Scenario Outline: Reject amend for owned lock JSON missing required fields
    Given a locked bundle `<name>`
    When I replace the lock file `<path>` with valid JSON missing required fields
    And I run `intend amend <name>`
    Then it exits with code 1
    And stderr contains `lock file is missing required fields`

    Examples:
      | name         | path                            |
      | health-check | .intend/locks/health-check.json |

  Scenario Outline: Reject amend for contribution lock JSON missing required fields
    Given a locked contribution bundle `<name>`
    When I replace the lock file `<path>` with valid JSON missing required fields
    And I run `intend amend --mode contrib <name>`
    Then it exits with code 1
    And stderr contains `lock file is missing required fields`

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/locks/startup-fix.json |

  Scenario Outline: Stop verify for corrupted owned lock JSON
    Given a locked bundle `<name>`
    And verification tools are available
    When I replace the contents of `<path>`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `decode lock file for <name>`
    And the verification log is empty

    Examples:
      | name         | path                            |
      | health-check | .intend/locks/health-check.json |

  Scenario Outline: Stop verify for corrupted contribution lock JSON
    Given a locked contribution bundle `<name>`
    And verification tools are available
    When I replace the contents of `<path>`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `decode lock file for <name>`
    And the verification log is empty

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/locks/startup-fix.json |

  Scenario Outline: Stop verify for owned lock JSON missing required fields
    Given a locked bundle `<name>`
    And verification tools are available
    When I replace the lock file `<path>` with valid JSON missing required fields
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `lock file is missing required fields`
    And the verification log is empty

    Examples:
      | name         | path                            |
      | health-check | .intend/locks/health-check.json |

  Scenario Outline: Stop verify for contribution lock JSON missing required fields
    Given a locked contribution bundle `<name>`
    And verification tools are available
    When I replace the lock file `<path>` with valid JSON missing required fields
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `lock file is missing required fields`
    And the verification log is empty

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/locks/startup-fix.json |
