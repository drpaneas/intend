# Reject lock path safety violations across lifecycle phases

Clarify how `intend` handles lock metadata and lock-tracked files that escape the allowed root through absolute paths, upward traversal, or symlink resolution.

This slice consolidates the current lock path safety behavior across direct tracing and downstream consumers. The lock file still parses successfully, but a tracked path either points outside the workspace or contribution bundle root by text, or resolves outside it through a symlink. Lifecycle commands that consume lock metadata must reject that input before drift or verifier work proceeds.

## Commands in scope

- `intend trace <name>`
- `intend trace --mode contrib <name>`
- `intend amend <name>`
- `intend amend --mode contrib <name>`
- `intend verify`

## Required behavior

- If lock metadata uses an absolute path or upward traversal that escapes the allowed root, the consuming command exits with code `1`.
- If a lock-tracked file resolves outside the allowed root through a symlink, the consuming command exits with code `1`.
- `intend verify` stops before running verification tools.
- The failure messages identify the offending tracked path and the allowed root.

## Constraints

- Keep the checks file-based and local.
- Do not add recovery or auto-rewrite behavior in this slice.
- Keep trace metadata path validation as a separate contract family.
- Do not widen this slice to initial locking, because `intend lock` does not consume an existing lock file.

## Done when

- The feature contract for lock path safety validation matrix passes through `godog`.
- Lock metadata and lock-tracked files cannot escape their allowed roots across trace, amend, and verify.
