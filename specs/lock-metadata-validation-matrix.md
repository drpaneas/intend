# Reject invalid lock metadata across lifecycle phases

Clarify how `intend` handles lock metadata when the lock file is corrupted JSON or valid JSON missing required fields.

This slice consolidates the current lock metadata validity behavior across direct tracing and downstream consumers. The lock file exists, but it is either unreadable JSON or incomplete metadata, so lifecycle commands that consume lock metadata must reject that input before drift or verifier work proceeds.

## Commands in scope

- `intend trace <name>`
- `intend trace --mode contrib <name>`
- `intend amend <name>`
- `intend amend --mode contrib <name>`
- `intend verify`

## Required behavior

- If a lock file contains invalid JSON, the consuming command exits with code `1`.
- If a lock file is valid JSON but missing required fields, the consuming command exits with code `1`.
- `intend verify` stops before running verification tools.
- The failure messages report either the lock decode failure or the missing required fields.

## Constraints

- Keep the checks file-based and local.
- Do not add recovery or auto-rewrite behavior in this slice.
- Preserve the existing precedence for lock metadata problems ahead of drift evaluation.
- Keep trace metadata validation as a separate contract family.
- Do not widen this slice to initial locking, because `intend lock` does not consume an existing lock file.

## Done when

- The feature contract for lock metadata validation matrix passes through `godog`.
- Lock JSON corruption and missing required fields are explicitly covered across trace, amend, and verify.
