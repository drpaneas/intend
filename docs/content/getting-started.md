---
title: "Getting started"
weight: 10
summary: "Build the tool, initialize a workspace, and create the first bundle."
---

`intend` is small on purpose. Write the intent. Write the feature. Write the tests. Then write the code.

## Build

```bash
go build ./cmd/intend
```

## Initialize a workspace

```bash
./intend init
```

This creates:

```text
specs/
features/
.intend/trace/
.intend/locks/
```

## Create the first bundle

```bash
./intend new hello-world
```

This creates:

```text
specs/hello-world.md
features/hello-world.feature
.intend/trace/hello-world.json
```

## Lock and check it

```bash
./intend lock hello-world
./intend trace hello-world
./intend verify
```

If the contract changes on purpose, record it:

```bash
./intend amend hello-world
```
