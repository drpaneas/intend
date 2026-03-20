# Reject invalid contribution trace issue metadata across all lifecycle phases

Clarify how `intend` validates contribution trace `issueRef` and `issuePath` metadata when those fields are missing, invalid, or non-canonical.

This slice consolidates the current contribution trace issue-metadata validation behavior across trace, amend, verify, and initial lock. Regardless of which lifecycle command consumes the trace metadata, missing issue fields, invalid GitHub issue references, or non-canonical issue snapshot paths must be rejected as explicit metadata errors rather than being treated as ordinary contract drift.

## Commands in scope

- `intend trace --mode contrib <name>`
- `intend amend --mode contrib <name>`
- `intend verify`
- `intend lock --mode contrib <name>`

## Required behavior

- If a contribution trace file is missing `issueRef`, each command exits with code `1`.
- If a contribution trace file is missing `issuePath`, each command exits with code `1`.
- If a contribution trace file contains an invalid `issueRef`, each command exits with code `1`.
- If a contribution trace file uses an in-bundle `issuePath` other than `issue.json`, each command exits with code `1`.
- `intend verify` stops before running verification tools.
- `intend lock --mode contrib <name>` leaves no lock file behind on failure.
- The failure messages report the specific contribution issue metadata problem.

## Constraints

- Keep the validation file-based and local.
- Do not add recovery or auto-rewrite behavior in this slice.
- Preserve existing precedence for malformed JSON, missing required spec or feature paths, metadata identity mismatch, and issue path safety errors.

## Done when

- The feature contract for contribution trace issue metadata validation matrix passes through `godog`.
- Contribution trace issue metadata errors are explicitly covered across trace, amend, verify, and initial lock.
