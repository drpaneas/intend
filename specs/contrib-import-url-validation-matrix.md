# Reject invalid GitHub issue URLs during contribution import

Clarify how `intend new --mode contrib` handles imported GitHub issue URLs that violate host, hostname, shape, or canonical URL rules.

This slice consolidates the current contribution import URL validation behavior. `gh issue view` may return URLs that are not valid canonical GitHub issue URLs, and contrib import must reject them before any shadow bundle state is written.

## Command in scope

- `intend new --mode contrib --from-gh owner/repo#123 <name>`

## Required behavior

- If `gh issue view` returns a non-GitHub issue URL, the command exits with code `1`.
- If `gh issue view` returns an alternate GitHub hostname, the command exits with code `1`.
- If `gh issue view` returns a GitHub URL that is not an issue URL, the command exits with code `1`.
- If `gh issue view` returns a GitHub issue URL with a query string or fragment, the command exits with code `1`.
- The failure message reports the returned URL or hostname that violated the URL policy.
- On each failure path, `intend` does not create a partial shadow bundle under Git metadata.

## Constraints

- Keep the validation local and CLI-first.
- Validate the imported issue before writing any bundle files.
- Do not add URL rewriting or normalization behavior in this slice.

## Done when

- The feature contract for contribution import URL validation matrix passes through `godog`.
- Invalid imported GitHub issue URLs fail before `$(git rev-parse --git-dir)/intend/contrib/<name>/` exists.
