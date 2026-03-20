# Keep amend idempotent when nothing changed

Harden the narrowed v1 CLI so `intend amend` only records intentional contract changes.

This slice closes the gap where running `amend` twice with identical contract files still increments the lock version. In v1, unchanged contracts should produce a clear no-op result instead.

## Commands in scope

- `intend amend <name>`
- `intend amend --mode contrib <name>`

## Required behavior

- If an owned bundle is locked and its contract files are unchanged, `intend amend <name>`:
  - exits with code `0`
  - reports that the contract is unchanged
  - leaves the lock version unchanged
- If a contribution bundle is locked and its contract files are unchanged, `intend amend --mode contrib <name>`:
  - exits with code `0`
  - reports that the contract is unchanged
  - leaves the lock version unchanged

## Constraints

- Do not rewrite the lock just to bump a version number.
- Do not treat a successful no-op as an error.
- Keep the owned and contribution semantics consistent.

## CLI rules

- Exit code stays `0` for a no-op amend.

## Done when

- The feature contract for amend idempotency passes through `godog`.
- Re-running `amend` without contract changes leaves the version unchanged for both owned and contribution bundles.
