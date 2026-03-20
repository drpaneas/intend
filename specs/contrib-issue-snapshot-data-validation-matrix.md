# Reject malformed or incomplete contribution issue snapshots across all lifecycle phases

Clarify how `intend` validates the on-disk contribution `issue.json` snapshot when its JSON is malformed or missing required fields.

This slice consolidates the current data-shape validation behavior for contribution issue snapshots across trace, amend, verify, and initial lock. Regardless of which lifecycle command consumes the snapshot, malformed JSON or missing required issue fields must be rejected as snapshot-data errors rather than being treated as ordinary contract drift.

## Commands in scope

- `intend trace --mode contrib <name>`
- `intend amend --mode contrib <name>`
- `intend verify`
- `intend lock --mode contrib <name>`

## Required behavior

- If the contribution issue snapshot is malformed JSON, each command exits with code `1`.
- If the contribution issue snapshot is missing a required issue number, title, or URL, each command exits with code `1`.
- `intend verify` stops before running verification tools.
- `intend lock --mode contrib <name>` leaves no lock file behind on failure.
- The failure message reports either a snapshot decode failure or missing required fields.

## Constraints

- Keep the validation local and file-based.
- Preserve the existing missing-file precedence for `issue.json`.
- Do not reinterpret malformed or incomplete snapshot data as identity mismatch or generic contract drift.

## Done when

- The feature contract for contribution issue snapshot data validation matrix passes through `godog`.
- Contribution issue snapshot data-shape errors are explicitly covered across trace, amend, verify, and initial lock.
