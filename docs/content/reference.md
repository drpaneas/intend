---
title: "Reference"
weight: 40
summary: "Command forms, file layout, exit codes, and required tools."
---

This page is for day-to-day lookup.

## Commands

| Command | Meaning | Notes |
| --- | --- | --- |
| `intend init` | Initialize the owned workspace layout | Checks required tools before writing files |
| `intend new <name>` | Create an owned bundle | Creates spec, feature, and trace |
| `intend new --mode contrib --from-gh owner/repo#123 <name>` | Create a contribution shadow bundle | Requires `git` and `gh` |
| `intend lock <name>` | Lock an owned bundle for the first time | Fails if it is already locked |
| `intend lock --mode contrib <name>` | Lock a contribution bundle | Writes under Git metadata |
| `intend trace <name>` | Validate and compare an owned bundle to its lock | Reports contract drift if the baseline changed |
| `intend trace --mode contrib <name>` | Validate and compare a contribution bundle | Reads from the shadow bundle |
| `intend amend <name>` | Record an intentional contract change | Increments the lock version |
| `intend amend --mode contrib <name>` | Amend a contribution bundle | Same rule, different storage location |
| `intend delete <name>` | Delete an owned bundle | Refuses locked bundles unless `--force` is used |
| `intend delete --mode contrib <name>` | Delete a contribution bundle | Removes the shadow bundle, refuses locked bundles unless `--force` is used |
| `intend verify` | Trace all bundles, then run repo checks | Uses external tools |
| `intend agent install cursor` | Install managed Cursor guidance | Safe to re-run, refuses overwrite of edited managed files |

## Exit codes

| Code | Meaning |
| --- | --- |
| `0` | success |
| `1` | contract drift or command failure |
| `2` | invalid CLI usage |

## Owned workspace layout

After `intend init`:

```text
specs/
features/
.intend/trace/
.intend/locks/
```

For an owned bundle called `add-login`:

```text
specs/add-login.md
features/add-login.feature
.intend/trace/add-login.json
.intend/locks/add-login.json
```

## Contribution bundle layout

For a contribution bundle called `fix-issue-123`:

```text
$(git rev-parse --git-dir)/intend/contrib/fix-issue-123/
├── issue.json
├── specs/fix-issue-123.md
├── features/fix-issue-123.feature
├── trace/fix-issue-123.json
└── locks/fix-issue-123.json
```

The contribution bundle stays under Git metadata, not in the tracked working tree.

## Required tools

### Required by `intend init`

- `go`
- `golangci-lint`
- `trufflehog`
- `gitleaks`
- `trivy`

### Required by `intend verify`

- `go`
- `golangci-lint`
- `trufflehog`
- `gitleaks`
- `trivy`

`intend verify` runs:

```text
go test ./...
golangci-lint run
trufflehog filesystem .
gitleaks dir .
trivy fs .
```

### Required by contribution mode

- `git` for all contribution-mode operations
- `gh` for `intend new --mode contrib --from-gh ...`

## Naming

Bundle names use lowercase kebab-case:

```text
good: add-login
good: fix-issue-123
bad: AddLogin
bad: add_login
```

## Notes

- `trace` expects the bundle to already be locked.
- `delete` removes bundle files from the active mode and requires `--force` for locked bundles.
- `verify` traces bundles before running repository checks.
- Current agent support is Cursor-only.
