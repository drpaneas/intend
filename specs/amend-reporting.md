# Report amend outcomes clearly

Clarify the plain-text operator feedback for `intend amend` across owned bundles, current contribution bundles, and older contribution locks that require a semantic metadata upgrade.

This slice makes amend reporting explicit in one place. The message should tell a human whether amend was a no-op, an ordinary version bump, or a contribution lock rewrite that also upgraded semantic lock metadata.

## Commands in scope

- `intend amend <name>`
- `intend amend --mode contrib <name>`

## Required behavior

- If an owned amend is a no-op, `intend amend <name>` exits with code `0` and reports `contract unchanged`.
- If an owned amend records an intentional contract change, `intend amend <name>` exits with code `0` and reports `amended <name> to version <n>`.
- If a current contribution amend is a no-op, `intend amend --mode contrib <name>` exits with code `0`, reports `contract unchanged`, and does not report a semantic lock metadata upgrade.
- If a current contribution amend records an intentional contract change, `intend amend --mode contrib <name>` exits with code `0`, reports `amended <name> to version <n>`, and does not report a semantic lock metadata upgrade.
- If an older contribution lock is rewritten during amend because it lacks semantic lock metadata, `intend amend --mode contrib <name>` exits with code `0`, reports `amended <name> to version <n>`, and also reports that semantic lock metadata was upgraded.

## Constraints

- Keep reporting plain text and line-oriented.
- Do not change amend exit codes in this slice.
- Do not change traceability, locking, or semantic digest rules in this slice.

## Done when

- The feature contract for amend reporting passes through `godog`.
- A human can distinguish no-op amend, ordinary amend, and contribution semantic metadata upgrade amend from stdout alone.
