# Reject trace path safety violations across lifecycle phases

Clarify how `intend` handles trace metadata and trace-tracked files that escape the allowed root through absolute paths, upward traversal, or symlink resolution.

This slice consolidates the current trace path safety behavior across direct tracing, downstream consumers, and initial locking. The trace file still parses successfully, but a tracked field either points outside the workspace or contribution bundle root by text, or resolves outside it through a symlink. Lifecycle commands must reject that input before doing digest, drift, or verifier work.

## Commands in scope

- `intend trace <name>`
- `intend trace --mode contrib <name>`
- `intend amend <name>`
- `intend amend --mode contrib <name>`
- `intend verify`
- `intend lock <name>`
- `intend lock --mode contrib <name>`

## Required behavior

- If trace metadata uses an absolute path or upward traversal that escapes the allowed root, the consuming command exits with code `1`.
- If a trace-tracked file resolves outside the allowed root through a symlink, the consuming command exits with code `1`.
- `intend verify` stops before running verification tools.
- `intend lock` leaves no lock file behind on failure.
- For `intend amend`, a symlink escape on a path already tracked in the lock may be reported by lock validation before trace validation.
- The failure messages identify the offending field or tracked path and the allowed root.

## Constraints

- Keep the checks file-based and local.
- Do not add recovery or auto-rewrite behavior in this slice.
- Keep lock metadata path validation as a separate contract family.

## Done when

- The feature contract for trace path safety validation matrix passes through `godog`.
- Trace metadata and trace-tracked files cannot escape their allowed roots across trace, amend, verify, and initial lock.
