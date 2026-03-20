# Reject invalid GitHub issue payload data during contribution import

Clarify how `intend new --mode contrib` handles imported GitHub issue payloads that are malformed JSON or incomplete issue data.

This slice consolidates the current contribution import payload validation behavior. `gh issue view` may return malformed JSON or issue objects missing required fields, and contrib import must reject those payloads before any shadow bundle state is written.

## Command in scope

- `intend new --mode contrib --from-gh owner/repo#123 <name>`

## Required behavior

- If `gh issue view` returns malformed JSON, the command exits with code `1`.
- If `gh issue view` returns incomplete issue data, the command exits with code `1`.
- The failure message reports whether the problem was invalid JSON or incomplete issue data.
- On either failure path, `intend` does not create a partial shadow bundle under Git metadata.

## Constraints

- Keep the validation local and CLI-first.
- Validate the imported issue before writing bundle files.
- Do not add cleanup code for partially written bundles if validation can happen earlier.

## Done when

- The feature contract for contribution import data validation matrix passes through `godog`.
- Invalid GitHub issue payload data fails before `$(git rev-parse --git-dir)/intend/contrib/<name>/` exists.
