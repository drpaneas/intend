# Reject mismatched metadata identity across lifecycle phases

Clarify how `intend` handles trace and lock metadata whose embedded bundle identity no longer matches the bundle and mode being operated on.

This slice consolidates the current metadata identity validation behavior across direct tracing, downstream consumers, and initial locking where applicable. The metadata files still parse successfully, but the embedded `name` or `mode` points at a different bundle or mode than the requested one, so lifecycle commands must reject that input explicitly before doing further work.

## Commands in scope

- `intend trace <name>`
- `intend trace --mode contrib <name>`
- `intend amend <name>`
- `intend amend --mode contrib <name>`
- `intend verify`
- `intend lock <name>`
- `intend lock --mode contrib <name>`

## Required behavior

- If trace metadata identity is mismatched, the lifecycle command consuming that trace file exits with code `1`.
- If lock metadata `name` is mismatched, the lifecycle command consuming that lock file exits with code `1`.
- `intend verify` stops before running verification tools.
- `intend lock` leaves no lock file behind on failure.
- The failure messages report the expected and actual metadata identity.

## Constraints

- Keep the check file-based and local.
- Do not add recovery or auto-rewrite behavior in this slice.
- Preserve existing precedence for malformed JSON and missing required paths or fields.

## Done when

- The feature contract for metadata identity validation matrix passes through `godog`.
- Metadata identity mismatches are explicitly covered across trace, amend, verify, and initial lock where applicable.
