# Clarify contribution issue snapshot edit behavior under current semantic locks

Clarify how `intend` handles contribution `issue.json` edits when the bundle is using the current semantic-digest lock model.

This slice consolidates the current behavior matrix for contribution issue snapshot edits that preserve issue identity. Under current semantic locks, semantic content changes such as title, body, or extra unknown fields remain tracked contract content, while representation-only rewrites such as formatting or key order do not count as contract changes. Older contribution locks that do not yet contain `semanticFiles` remain covered separately by the backward-compatibility slice.

## Commands in scope

- `intend trace --mode contrib <name>`
- `intend amend --mode contrib <name>`
- `intend verify`
- `intend lock --mode contrib <name>`

## Required behavior

- If a locked contribution issue snapshot title or body changes while issue identity remains valid, `intend trace --mode contrib <name>` reports contract drift, `intend amend --mode contrib <name>` can intentionally re-lock the change, and `intend verify` stops before running verification tools.
- If a locked contribution issue snapshot gains an extra unknown JSON field while issue identity remains valid, that edit behaves like ordinary tracked contract content across `trace`, `amend`, and `verify`.
- If an existing contribution bundle's issue snapshot changes only in descriptive content or gains extra unknown fields before the first lock, `intend lock --mode contrib <name>` accepts the edited snapshot and writes the first lock.
- If a locked contribution issue snapshot is rewritten with formatting-only or key-order-only changes, `intend trace --mode contrib <name>` succeeds, `intend amend --mode contrib <name>` remains a no-op, and `intend verify` still runs verification tools.
- If an existing contribution bundle's issue snapshot is rewritten with formatting-only or key-order-only changes before the first lock, `intend lock --mode contrib <name>` accepts the rewritten snapshot and writes the first lock.

## Constraints

- Do not reinterpret non-identity content edits, unknown fields, or representation-only rewrites as identity validation failures.
- Keep the existing required-field, issue number, repository, URL, and semantic lock backward-compatibility behavior unchanged.
- Keep the behavior file-based and local.

## Done when

- The feature contract for contribution issue snapshot edit behavior passes through `godog`.
- Current semantic-lock behavior for contribution issue snapshot content edits, extra fields, and representation-only rewrites is explicit in one place.
