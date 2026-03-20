# Create a contribution-mode shadow bundle

Build the first contribution-mode command for the narrowed v1 CLI.

This slice only covers importing a GitHub issue into a local shadow bundle under Git metadata. It should help contributors use `intend` locally without polluting the upstream working tree.

## Command in scope

- `intend new --mode contrib --from-gh owner/repo#123 <name>`

## Required behavior

- The command requires a Git repository because contribution bundles live under Git metadata.
- The requested GitHub issue reference must use a positive issue number.
- It resolves the Git metadata directory through `git rev-parse --git-dir`.
- It imports the referenced GitHub issue through `gh issue view`.
- It writes a shadow bundle under `$(git rev-parse --git-dir)/intend/contrib/<name>/`.
- The shadow bundle contains:
  - `issue.json`
  - `specs/<name>.md`
  - `features/<name>.feature`
  - `trace/<name>.json`
- `issue.json` stores the imported GitHub issue snapshot.
- The generated trace metadata records that the bundle mode is `contrib`.

## CLI rules

- Exit `0` on success.
- Exit `1` when the issue reference is invalid, when the directory is not a Git repository, when `gh` is missing, or when issue import fails.
- Exit `2` for invalid CLI usage.

## Constraints

- Keep contribution state local and file-based.
- Do not write the contribution bundle into the upstream working tree.
- Use the `git` and `gh` CLIs instead of embedding Git or GitHub APIs in this slice.

## Done when

- The feature contract for contribution-mode bundle creation passes through `godog`.
- Running the command creates the expected shadow files under Git metadata, and invalid non-positive issue refs fail before any shadow bundle state is written.
