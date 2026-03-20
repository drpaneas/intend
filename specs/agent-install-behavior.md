# Install Cursor guidance safely

Clarify how `intend agent install cursor` behaves for initial installation, unsupported agents, safe re-runs, and locally modified managed files.

This slice consolidates the current Cursor agent installation behavior. The command should install the managed Cursor workflow and built-in skills, reject unsupported agents, remain idempotent when managed files are unchanged, and refuse to overwrite locally edited managed files.

## Command in scope

- `intend agent install cursor`
- `intend agent install <other-agent>`

## Required behavior

- `intend agent install cursor` creates the managed Cursor command and skill files when they are missing.
- Re-running `intend agent install cursor` with unchanged managed files exits with code `0`.
- If a managed Cursor file exists and its contents differ from the built-in managed content, `intend agent install cursor` exits with code `1`.
- The failure message for a locally edited file names the modified managed file.
- The locally edited file remains unchanged after the failed install attempt.
- In v1, any agent name other than `cursor` is rejected with a clear error.

## Content requirements

- The workflow command explains the `intend` contract flow briefly and points at the core verbs.
- The built-in skills stay compact and procedural.
- The Go skills cover naming and package shape, error handling, testing discipline, context usage, and logging and secret safety.

## Constraints

- Keep the installation file-based and repo-local.
- Do not silently overwrite locally edited managed files.
- Do not add a `--force` flag in this slice.
- Do not add support for Claude, Codex, Copilot, or Gemini in this slice.

## Done when

- The feature contract for agent install behavior passes through `godog`.
- Running `intend agent install cursor` produces the expected managed files, remains safe to re-run, and protects local edits from silent overwrite.
