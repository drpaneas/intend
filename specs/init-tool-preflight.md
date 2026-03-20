# Preflight required tools before init

Harden `intend init` so it fails early when the required always-on verification tools are not installed.

This slice is only about preflight. The normal successful layout created by `intend init` is already covered by the core contract lifecycle slice. Here we want `init` to check its required external binaries before it writes any repository state.

## Command in scope

- `intend init`

## Required behavior

- Before creating any `intend` workspace directories, `intend init` checks for these required binaries:
  - `go`
  - `golangci-lint`
  - `trufflehog`
  - `gitleaks`
  - `trivy`
- If any required binary is missing, `intend init` exits with code `1`.
- The failure message names the missing tool.
- On this failure path, `intend init` does not create partial workspace directories.

## Constraints

- Keep the check CLI-first and local.
- Do not install tools automatically.
- Do not add optional tool profiles or configuration in this slice.

## Done when

- The feature contract for init preflight passes through `godog`.
- Missing required tools stop `intend init` before `.intend` or other workspace directories are created.
