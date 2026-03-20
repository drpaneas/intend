# Run verification evidence for a Go repository

Build `intend verify` for the narrowed v1 CLI.

This command is responsible for proving that a Go repository is healthy after contract traceability is intact. It must remain plain-file and CLI-first, and it must use external verification tools instead of embedding their logic.

## Command in scope

- `intend verify`

## Required behavior

- `intend verify` first validates all known contract bundles:
  - owned bundles under `.intend/trace/`
  - contribution shadow bundles under `$(git rev-parse --git-dir)/intend/contrib/`
- If any bundle has contract drift, `intend verify` exits before running external verification tools.
- After traceability passes, `intend verify` runs these evidence commands from the repository root:
  - `go test ./...`
  - `golangci-lint run`
  - `trufflehog filesystem .`
  - `gitleaks dir .`
  - `trivy fs .`
- If a required tool is missing, `intend verify` fails with a clear installation-oriented message naming the missing tool.
- If any verification command fails, `intend verify` exits with code `1`.
- If all verification commands succeed, `intend verify` exits with code `0`.

## CLI rules

- Exit `0` on success.
- Exit `1` on contract drift, missing required tools, or verification command failures.
- Exit `2` for invalid CLI usage.

## Constraints

- Keep `intend verify` repo-level and CLI-first.
- Do not invent a custom linter or scanner API.
- Do not add `act`, profiling, benchmarking, or dependency maintenance to this command in v1.
- Do not require Git history to decide whether the contract drifted.

## Done when

- The feature contract for `intend verify` passes through `godog`.
- The implementation proves both the trace-first behavior and the external verification orchestration.
