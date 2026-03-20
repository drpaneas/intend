# Delete an owned or contribution bundle

Add an explicit delete command to the narrowed v1 CLI so bundles can be removed intentionally instead of by hand.

This slice is about clean removal, not contract recovery. It should let a developer abandon an owned bundle or a contribution shadow bundle while keeping the default behavior safe around locked contract state.

## Commands in scope

- `intend delete <name>`
- `intend delete --force <name>`
- `intend delete --mode contrib <name>`
- `intend delete --mode contrib --force <name>`

## Required behavior

- `intend delete <name>` removes an unlocked owned bundle by deleting:
  - `specs/<name>.md`
  - `features/<name>.feature`
  - `.intend/trace/<name>.json`
  - `.intend/locks/<name>.json` if it exists
- `intend delete <name>` refuses to remove a locked owned bundle unless `--force` is provided.
- `intend delete --force <name>` removes a locked owned bundle and its lock file.
- `intend delete --mode contrib <name>` removes an unlocked contribution shadow bundle rooted at:
  - `$(git rev-parse --git-dir)/intend/contrib/<name>/`
- `intend delete --mode contrib <name>` refuses to remove a locked contribution bundle unless `--force` is provided.
- After a bundle is deleted, `intend verify` no longer traces it.

## CLI rules

- Exit `0` on success.
- Exit `1` for locked bundles without `--force`, missing bundle state, missing Git context, or command execution failures.
- Exit `2` for invalid CLI usage.

## Constraints

- Keep deletion file-based and local.
- Do not add interactive prompts in this slice.
- Do not silently unlock a bundle before deleting it.
- Do not delete the owned workspace directories created by `intend init`.

## Done when

- The feature contract for bundle deletion passes through `godog`.
- Owned and contribution bundles can be removed intentionally through the CLI.
- Deleted bundles are no longer part of verification trace enumeration.
