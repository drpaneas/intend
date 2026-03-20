# Clarify backward compatibility for contribution locks without semantic digests

Clarify how `intend` behaves when a contribution bundle lock predates the `semanticFiles` metadata used for semantic issue snapshot comparison.

This slice makes the upgrade behavior explicit instead of implicit. Older contribution locks without `semanticFiles` must remain readable. Until such a lock is explicitly rewritten by `intend amend --mode contrib`, representation-only rewrites of `issue.json` still compare by raw file bytes and therefore behave like ordinary drift. An explicit amend upgrades the lock to the current semantic-digest format, after which later representation-only rewrites no longer count as drift.

## Commands in scope

- `intend trace --mode contrib <name>`
- `intend amend --mode contrib <name>`
- `intend verify`

## Required behavior

- If a locked contribution bundle's lock file does not contain `semanticFiles`, a formatting-only rewrite of `issue.json` still causes `intend trace --mode contrib <name>` to exit with code `1` for contract drift.
- If a locked contribution bundle's lock file does not contain `semanticFiles`, a formatting-only rewrite of `issue.json` still causes `intend verify` to exit with code `1` for contract drift before running verification tools.
- If a locked contribution bundle's lock file does not contain `semanticFiles`, `intend amend --mode contrib <name>` can intentionally rewrite the lock after a representation-only `issue.json` change and bumps the contribution lock version.
- After that explicit amend rewrites the lock into the current semantic-digest format, later representation-only rewrites of `issue.json` no longer cause drift.

## Constraints

- Keep older contribution locks readable without requiring a manual file migration.
- Do not silently upgrade older locks during `trace` or `verify`.
- Keep the upgrade path explicit and local to `intend amend --mode contrib`.
- Keep contribution issue identity validation unchanged.

## Done when

- The feature contract for contribution semantic lock backward compatibility passes through `godog`.
- Older contribution locks remain usable, and their upgrade semantics are explicit across trace, amend, and verify.
