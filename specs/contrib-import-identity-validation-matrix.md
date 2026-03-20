# Reject mismatched GitHub issue identity during contribution import

Clarify how `intend new --mode contrib` handles imported GitHub issue payloads whose returned issue identity does not match the requested ref.

This slice consolidates the current contribution import identity validation behavior. `gh issue view` may return the wrong issue number, a URL for a different repository, or a GitHub issue URL with the wrong issue number, and contrib import must reject those mismatches before any shadow bundle state is written.

## Command in scope

- `intend new --mode contrib --from-gh owner/repo#123 <name>`

## Required behavior

- If `gh issue view` returns a different issue number in the JSON payload, the command exits with code `1`.
- If `gh issue view` returns a URL for a different repository, the command exits with code `1`.
- If `gh issue view` returns a GitHub issue URL with a different issue number, the command exits with code `1`.
- The failure message reports the returned identity and the expected requested identity.
- On each failure path, `intend` does not create a partial shadow bundle under Git metadata.

## Constraints

- Keep the validation local and CLI-first.
- Validate the imported issue before writing any bundle files.
- Keep broader URL host and shape policy validation as a separate contract family.

## Done when

- The feature contract for contribution import identity validation matrix passes through `godog`.
- Mismatched imported GitHub issue identity fails before `$(git rev-parse --git-dir)/intend/contrib/<name>/` exists.
