# Reject malformed semantic lock metadata

Harden `intend` so malformed `semanticFiles` lock metadata is rejected explicitly instead of being treated as ordinary drift or silently normalized by amendment.

This slice covers the new semantic digest metadata introduced for contribution issue snapshots. That metadata is only valid for the contribution issue snapshot path itself. Owned bundles must not carry semantic digests, and contribution bundles must not attach semantic digests to tracked files other than `issue.json`.

## Commands in scope

- `intend trace <name>`
- `intend amend --mode contrib <name>`
- `intend verify`

## Required behavior

- If an owned lock file includes a `semanticFiles` entry, `intend trace <name>` exits with code `1`.
- If a contribution lock file includes a `semanticFiles` entry for a tracked file other than `issue.json`, `intend amend --mode contrib <name>` exits with code `1`.
- If a contribution lock file includes a `semanticFiles` entry for a tracked file other than `issue.json`, `intend verify` exits with code `1` before running verification tools.
- Failures report malformed semantic lock metadata rather than contract drift.

## Constraints

- Keep semantic digests backward-compatible for valid contribution `issue.json` entries.
- Do not let `intend amend` rewrite malformed semantic lock metadata into a new valid lock.
- Keep the validation local and lock-file based.

## Done when

- The feature contract for semantic lock metadata validation passes through `godog`.
- Owned and contribution consumers reject unsupported `semanticFiles` paths before drift detection or verifier execution.
