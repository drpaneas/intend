# Reject mismatched contribution issue snapshots across all lifecycle phases

Clarify how `intend` validates the semantic identity of the contribution `issue.json` snapshot against `trace.issueRef`.

This slice consolidates the current identity-validation behavior for edited contribution issue snapshots across trace, amend, verify, and initial lock. The snapshot still has valid paths and valid JSON, but it no longer points at the traced GitHub issue because the issue number or repository changed. Every lifecycle command that consumes the snapshot must reject that mismatch explicitly.

## Commands in scope

- `intend trace --mode contrib <name>`
- `intend amend --mode contrib <name>`
- `intend verify`
- `intend lock --mode contrib <name>`

## Required behavior

- If the contribution issue snapshot `number` differs from the issue number in `trace.issueRef`, each command exits with code `1`.
- If the contribution issue snapshot URL points to a different repository than `trace.issueRef`, each command exits with code `1`.
- `intend verify` stops before running verification tools.
- `intend lock --mode contrib <name>` leaves no lock file behind on failure.
- The failure message reports the specific mismatch.

## Constraints

- Keep the validation local and file-based.
- Preserve existing path-validation and missing-file precedence for `issue.json`.
- Do not reinterpret identity mismatches as generic contract drift.
- Do not widen this slice to malformed snapshot data or invalid URL shape, which are covered elsewhere.

## Done when

- The feature contract for contribution issue snapshot identity validation matrix passes through `godog`.
- Contribution issue snapshot identity mismatches are explicitly covered across trace, amend, verify, and initial lock.
