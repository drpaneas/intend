# Reject missing bundle files in amend and verify

Harden `intend amend` and `intend verify` so they fail clearly when valid bundle metadata points at files that no longer exist on disk.

This slice is about missing bundle content, not corrupted JSON. The trace and lock metadata remain valid, but a referenced spec, feature, or contribution issue snapshot has been removed, so amendment and verification must stop before proceeding.

## Commands in scope

- `intend amend <name>`
- `intend amend --mode contrib <name>`
- `intend verify`

## Required behavior

- If an owned bundle spec file is missing, `intend amend <name>` exits with code `1`.
- If a contribution bundle issue snapshot is missing, `intend amend --mode contrib <name>` exits with code `1`.
- If an owned bundle feature file is missing, `intend verify` exits with code `1` before running verification tools.
- If a contribution bundle issue snapshot is missing, `intend verify` exits with code `1` before running verification tools.
- The failure messages identify the missing bundle path being digested.

## Constraints

- Keep the check file-based and local.
- Do not add recovery or auto-rewrite behavior in this slice.
- Do not widen this slice to corrupted metadata JSON or drift-only cases.

## Done when

- The feature contract for missing bundle files passes through `godog`.
- `intend amend` and `intend verify` both stop clearly when referenced bundle files are missing.
