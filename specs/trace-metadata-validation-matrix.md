# Reject invalid trace metadata across lifecycle phases

Clarify how `intend` handles trace metadata when the trace file is corrupted JSON or valid JSON missing required bundle paths.

This slice consolidates the current trace metadata validity behavior across direct tracing, downstream consumers, and initial locking. The trace file exists, but it is either unreadable JSON or incomplete metadata, so lifecycle commands must reject that input before lock, drift, or verifier work proceeds.

## Commands in scope

- `intend trace <name>`
- `intend trace --mode contrib <name>`
- `intend amend <name>`
- `intend amend --mode contrib <name>`
- `intend verify`
- `intend lock <name>`
- `intend lock --mode contrib <name>`

## Required behavior

- If a trace file contains invalid JSON, the consuming command exits with code `1`.
- If a trace file is valid JSON but missing required paths, the consuming command exits with code `1`.
- `intend verify` stops before running verification tools.
- `intend lock` leaves no lock file behind on failure.
- The failure messages report either the trace decode failure or the missing required paths.

## Constraints

- Keep the checks file-based and local.
- Do not add recovery or auto-rewrite behavior in this slice.
- Preserve the existing precedence for trace metadata problems ahead of lock or drift evaluation.
- Keep lock metadata validation as a separate contract family.

## Done when

- The feature contract for trace metadata validation matrix passes through `godog`.
- Trace JSON corruption and missing required paths are explicitly covered across trace, amend, verify, and initial lock.
