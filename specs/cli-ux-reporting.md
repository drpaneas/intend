# Report verification progress clearly

Improve the narrowed v1 CLI so `intend verify` is easier for humans to understand while it runs.

This slice is only about plain-text operator feedback. It should not change the verification model, the tool set, or the pass/fail rules.

## Commands in scope

- `intend verify`

## Required behavior

- `intend verify` prints a progress line before each contract or tool check it runs.
- For owned bundles, the progress line names the bundle being traced.
- For each verification tool, the progress line prints the exact command form that `intend` is about to run.
- Successful verification still ends with `verify ok`.
- Invalid CLI usage for `intend verify` prints a command-specific usage line:
  - `usage: intend verify`

## Constraints

- Keep reporting plain text and line-oriented.
- Do not add colors, spinners, tables, JSON output, or background execution in this slice.
- Do not change which tools `intend verify` runs.

## CLI rules

- Exit codes stay unchanged from the existing `verify` contract.

## Done when

- The feature contract for verify reporting passes through `godog`.
- A human can see which bundle and which verification command `intend verify` is processing.
