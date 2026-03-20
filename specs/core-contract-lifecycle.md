# Build the core contract lifecycle for `intend`

Build a Go 1.26 CLI named `intend` for Go projects.

This first implementation slice is only about the core contract lifecycle. It must prove that `intend` can create a bundle, lock it, detect drift, and record an intentional amendment without hidden state.

## Commands in scope

- `intend init`
- `intend new <name>`
- `intend lock <name>`
- `intend trace <name>`
- `intend amend <name>`

## Required behavior

- `intend init` creates the plain-file layout for an owned repository:
  - `specs/`
  - `features/`
  - `.intend/trace/`
  - `.intend/locks/`
- `intend new <name>` creates:
  - `specs/<name>.md`
  - `features/<name>.feature`
  - `.intend/trace/<name>.json`
- `intend lock <name>` writes a lock file at `.intend/locks/<name>.json`.
- The lock file stores digests for the locked contract artifacts:
  - `specs/<name>.md`
  - `features/<name>.feature`
  - `.intend/trace/<name>.json`
- `intend trace <name>` validates that the trace metadata points at existing files and, when a lock exists, reports contract drift if any locked artifact changes.
- `intend amend <name>` records an intentional contract change by rewriting the lock file and incrementing its version.

## CLI rules

- Exit `0` on success.
- Exit `1` for contract drift or command execution failures.
- Exit `2` for invalid CLI usage.

## Constraints

- Keep the implementation file-based and inspectable.
- Do not depend on Git history or Git diffs to detect drift.
- Do not add a TUI, performance workflow, dependency workflow, local Actions execution, or non-Cursor agent support in this slice.
- Do not hide state in a database, daemon, or remote service.

## Done when

- The feature contract for this slice passes through `godog`.
- The CLI can initialize a repo, create a bundle, lock it, detect a changed locked artifact, and amend the lock after an intentional change.
