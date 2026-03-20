# Manage the contribution contract lifecycle

Extend the narrowed v1 CLI so contribution-mode shadow bundles can use the same contract lifecycle verbs as owned bundles.

This slice is about local contract discipline for upstream work. It should let a contributor lock a shadow bundle, detect drift in it, and intentionally amend it, all without writing upstream workflow files into the working tree.

## Commands in scope

- `intend trace --mode contrib <name>`
- `intend lock --mode contrib <name>`
- `intend amend --mode contrib <name>`

## Required behavior

- `intend trace --mode contrib <name>` reads the shadow trace file from:
  - `$(git rev-parse --git-dir)/intend/contrib/<name>/trace/<name>.json`
- `intend lock --mode contrib <name>` writes:
  - `$(git rev-parse --git-dir)/intend/contrib/<name>/locks/<name>.json`
- The contribution lock stores digests for the locked shadow contract artifacts:
  - `issue.json`
  - `specs/<name>.md`
  - `features/<name>.feature`
  - `trace/<name>.json`
- `intend trace --mode contrib <name>` reports contract drift when any locked shadow artifact changes.
- `intend amend --mode contrib <name>` rewrites the shadow lock file and increments its version.

## CLI rules

- Exit `0` on success.
- Exit `1` for contract drift, missing Git context, or command execution failures.
- Exit `2` for invalid CLI usage.

## Constraints

- Keep contribution lifecycle state under Git metadata, not in the upstream working tree.
- Use file digests, not Git history, to detect drift.
- Do not add a TUI, remote service, or embedded Git implementation.

## Done when

- The feature contract for contribution lifecycle passes through `godog`.
- A contribution shadow bundle can be locked, traced, and amended locally.
