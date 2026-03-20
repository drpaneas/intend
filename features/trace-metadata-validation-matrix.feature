Feature: trace-metadata-validation-matrix
  Scenario Outline: Reject trace for corrupted owned trace JSON
    Given a locked bundle `<name>`
    When I replace the contents of `<path>`
    And I run `intend trace <name>`
    Then it exits with code 1
    And stderr contains `decode trace file for <name>`

    Examples:
      | name         | path                            |
      | health-check | .intend/trace/health-check.json |

  Scenario Outline: Reject trace for corrupted contribution trace JSON
    Given a locked contribution bundle `<name>`
    When I replace the contents of `<path>`
    And I run `intend trace --mode contrib <name>`
    Then it exits with code 1
    And stderr contains `decode trace file for <name>`

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/trace/startup-fix.json |

  Scenario Outline: Reject trace for owned trace JSON missing required paths
    Given a locked bundle `<name>`
    When I replace the trace file `<path>` with valid JSON missing required paths
    And I run `intend trace <name>`
    Then it exits with code 1
    And stderr contains `trace file is missing required paths`

    Examples:
      | name         | path                            |
      | health-check | .intend/trace/health-check.json |

  Scenario Outline: Reject trace for contribution trace JSON missing required paths
    Given a locked contribution bundle `<name>`
    When I replace the trace file `<path>` with valid JSON missing required paths
    And I run `intend trace --mode contrib <name>`
    Then it exits with code 1
    And stderr contains `trace file is missing required paths`

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/trace/startup-fix.json |

  Scenario Outline: Reject amend for corrupted owned trace JSON
    Given a locked bundle `<name>`
    When I replace the contents of `<path>`
    And I run `intend amend <name>`
    Then it exits with code 1
    And stderr contains `decode trace file for <name>`

    Examples:
      | name         | path                            |
      | health-check | .intend/trace/health-check.json |

  Scenario Outline: Reject amend for corrupted contribution trace JSON
    Given a locked contribution bundle `<name>`
    When I replace the contents of `<path>`
    And I run `intend amend --mode contrib <name>`
    Then it exits with code 1
    And stderr contains `decode trace file for <name>`

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/trace/startup-fix.json |

  Scenario Outline: Reject amend for owned trace JSON missing required paths
    Given a locked bundle `<name>`
    When I replace the trace file `<path>` with valid JSON missing required paths
    And I run `intend amend <name>`
    Then it exits with code 1
    And stderr contains `trace file is missing required paths`

    Examples:
      | name         | path                            |
      | health-check | .intend/trace/health-check.json |

  Scenario Outline: Reject amend for contribution trace JSON missing required paths
    Given a locked contribution bundle `<name>`
    When I replace the trace file `<path>` with valid JSON missing required paths
    And I run `intend amend --mode contrib <name>`
    Then it exits with code 1
    And stderr contains `trace file is missing required paths`

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/trace/startup-fix.json |

  Scenario Outline: Stop verify for corrupted owned trace JSON
    Given a locked bundle `<name>`
    And verification tools are available
    When I replace the contents of `<path>`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `decode trace file for <name>`
    And the verification log is empty

    Examples:
      | name         | path                            |
      | health-check | .intend/trace/health-check.json |

  Scenario Outline: Stop verify for corrupted contribution trace JSON
    Given a locked contribution bundle `<name>`
    And verification tools are available
    When I replace the contents of `<path>`
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `decode trace file for <name>`
    And the verification log is empty

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/trace/startup-fix.json |

  Scenario Outline: Stop verify for owned trace JSON missing required paths
    Given a locked bundle `<name>`
    And verification tools are available
    When I replace the trace file `<path>` with valid JSON missing required paths
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `trace file is missing required paths`
    And the verification log is empty

    Examples:
      | name         | path                            |
      | health-check | .intend/trace/health-check.json |

  Scenario Outline: Stop verify for contribution trace JSON missing required paths
    Given a locked contribution bundle `<name>`
    And verification tools are available
    When I replace the trace file `<path>` with valid JSON missing required paths
    And I run `intend verify`
    Then it exits with code 1
    And stderr contains `trace file is missing required paths`
    And the verification log is empty

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/trace/startup-fix.json |

  Scenario Outline: Reject initial lock for corrupted owned trace JSON
    Given an existing bundle `<name>`
    When I replace the contents of `<path>`
    And I run `intend lock <name>`
    Then it exits with code 1
    And stderr contains `decode trace file for <name>`
    And the path `.intend/locks/<name>.json` does not exist

    Examples:
      | name         | path                            |
      | health-check | .intend/trace/health-check.json |

  Scenario Outline: Reject initial lock for corrupted contribution trace JSON
    Given an existing contribution bundle `<name>`
    When I replace the contents of `<path>`
    And I run `intend lock --mode contrib <name>`
    Then it exits with code 1
    And stderr contains `decode trace file for <name>`
    And the path `.git/intend/contrib/<name>/locks/<name>.json` does not exist

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/trace/startup-fix.json |

  Scenario Outline: Reject initial lock for owned trace JSON missing required paths
    Given an existing bundle `<name>`
    When I replace the trace file `<path>` with valid JSON missing required paths
    And I run `intend lock <name>`
    Then it exits with code 1
    And stderr contains `trace file is missing required paths`
    And the path `.intend/locks/<name>.json` does not exist

    Examples:
      | name         | path                            |
      | health-check | .intend/trace/health-check.json |

  Scenario Outline: Reject initial lock for contribution trace JSON missing required paths
    Given an existing contribution bundle `<name>`
    When I replace the trace file `<path>` with valid JSON missing required paths
    And I run `intend lock --mode contrib <name>`
    Then it exits with code 1
    And stderr contains `trace file is missing required paths`
    And the path `.git/intend/contrib/<name>/locks/<name>.json` does not exist

    Examples:
      | name        | path                                                   |
      | startup-fix | .git/intend/contrib/startup-fix/trace/startup-fix.json |
