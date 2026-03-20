---
title: "Getting started"
weight: 20
summary: "Build intend, create the first bundle, lock its contract, and run verification."
---

Let us create your first `intend` bundle.

The goal is not just to run a command. The goal is to put one change under contract, lock it, and verify it with normal Go tooling.

## Before you start

`intend init` checks for these tools before it writes anything:

- `go`
- `golangci-lint`
- `trufflehog`
- `gitleaks`
- `trivy`

If one of them is missing, `init` stops early. That is deliberate.

## Build the CLI

```bash
go build ./cmd/intend
```

This produces `./intend` in the repository root.

## Initialize a workspace

```bash
./intend init
```

Example output:

```text
initialized intend workspace
```

This creates:

```text
specs/
features/
.intend/trace/
.intend/locks/
```

The layout is plain on purpose. You can inspect every file that `intend` manages.

## Create the first bundle

```bash
./intend new hello-world
```

Example output:

```text
created bundle hello-world
```

This creates:

```text
specs/hello-world.md
features/hello-world.feature
.intend/trace/hello-world.json
```

Now stop and edit the contract files before you write implementation code.

- `specs/hello-world.md` explains the change in plain language.
- `features/hello-world.feature` records the expected behavior.

## Lock the approved contract

```bash
./intend lock hello-world
```

Example output:

```text
locked hello-world at version 1
```

This writes:

```text
.intend/locks/hello-world.json
```

The lock file records the current baseline of the contract files.

## Trace the bundle

```bash
./intend trace hello-world
```

Example output:

```text
trace ok: hello-world
```

If a locked contract file changes, `intend trace hello-world` reports contract drift instead of silently accepting it.

## Implement and verify

Now write the Go code and tests for the change.

When you are ready, run:

```bash
./intend verify
```

`intend verify` traces all known bundles first. Then it runs the repository checks:

```text
go test ./...
golangci-lint run
trufflehog filesystem .
gitleaks dir .
trivy fs .
```

That is the heart of the workflow: first the contract, then the code, then the evidence.

## Change the contract on purpose

Sometimes the contract itself needs to move. When that happens, record the new baseline explicitly:

```bash
./intend amend hello-world
```

If nothing changed, `intend` says so. If the contract changed intentionally, the lock version is incremented.

## Next steps

- Read [Workflow](../workflow/) for the full owned and contribution-mode lifecycle.
- Read [Reference](../reference/) for commands, files, and exit codes.
