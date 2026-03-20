---
title: "intend"
---

Welcome to `intend`.

`intend` is a Go CLI for Intent-Driven Development. It keeps the contract for a change in plain files, not in prompts, screenshots, or memory.

You write:

- intent in `specs/`
- behavior in `features/`
- trace metadata in `.intend/trace/`
- locked contract state in `.intend/locks/`

The result is simple. A change has a name, a contract, and a recorded baseline. No hidden state. No guessing what was approved.

## What is `intend` for?

- Go projects that want a contract-first workflow
- teams using AI to draft code, but not to decide what "correct" means
- repositories that prefer plain files and normal command-line tools

If you can run `go test`, you can use `intend`.

## What is in this site?

- [Intent-Driven Development](intent-driven-development/) explains the idea and the mental model.
- [Getting started](getting-started/) walks through the first bundle from `init` to `verify`.
- [Workflow](workflow/) shows the day-to-day lifecycle, including `contrib` mode.
- [Reference](reference/) lists commands, files, exit codes, and tool requirements.

## Quick example

```bash
go build ./cmd/intend
./intend init
./intend new hello-http
./intend lock hello-http
./intend trace hello-http
```

Those five commands create a bundle, lock its contract, and confirm that the locked files still match.

## Next step

If the idea is new, start with [Intent-Driven Development](intent-driven-development/).

If you want to type commands right away, jump to [Getting started](getting-started/).
