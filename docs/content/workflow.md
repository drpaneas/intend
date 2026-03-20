---
title: "Workflow"
weight: 30
summary: "The day-to-day lifecycle for owned bundles, contribution bundles, and verification."
---

`intend` is a small tool. Most of its value is in the workflow it encourages.

This page shows that workflow from start to finish.

## The bundle

Every change lives in a bundle.

For an owned bundle called `add-login`, the important files are:

| File | Role |
| --- | --- |
| `specs/add-login.md` | intent |
| `features/add-login.feature` | behavior contract |
| `.intend/trace/add-login.json` | tracked file mapping |
| `.intend/locks/add-login.json` | locked baseline |

That is the whole unit of work.

## Owned bundles

### 1. Start the bundle

```bash
./intend new add-login
```

This creates the spec, feature, and trace files.

At this point there is no lock yet. The contract is still being drafted.

### 2. Review the contract before code

Edit:

- `specs/add-login.md`
- `features/add-login.feature`

Keep the order clear:

```text
spec -> feature -> tests -> implementation
```

The spec says what the change is for. The feature file says how the behavior will be judged.

### 3. Lock the approved baseline

```bash
./intend lock add-login
```

This writes the lock file and records version `1`.

From here on, the bundle has a baseline.

### 4. Trace while you work

```bash
./intend trace add-login
```

If the locked files still match, `trace` succeeds.

If one of them changed, `trace` reports contract drift. That is not an error in the tool. That is the tool doing its job.

### 5. Amend on purpose

Sometimes the contract itself must change. That happens. The important thing is to make the change explicit.

```bash
./intend amend add-login
```

Use `amend` when the contract changed on purpose after review.

Do not use it to hide accidental drift.

### 6. Delete a bundle you are abandoning

Sometimes you decide not to carry a bundle forward.

```bash
./intend delete add-login
```

This removes the spec, feature, trace, and any lock file for that owned bundle.

If the bundle is already locked, `delete` refuses by default. Use `--force` when you really mean to remove that locked contract state:

```bash
./intend delete --force add-login
```

### 7. Verify the repository

```bash
./intend verify
```

`verify` traces all known bundles first. Only then does it run the repository checks.

This keeps contract drift and implementation failures separate.

## Contribution mode

Owned bundles live in the working tree. Contribution bundles do not.

Contribution mode creates a local shadow bundle under Git metadata. That is useful when you want to work against an upstream issue without writing the contract files into the tracked repository tree.

### Create a contribution bundle

```bash
./intend new --mode contrib --from-gh owner/repo#123 fix-issue-123
```

This requires:

- a Git repository
- `gh`
- a positive issue number

The bundle is written under:

```text
$(git rev-parse --git-dir)/intend/contrib/fix-issue-123/
```

with:

```text
issue.json
specs/fix-issue-123.md
features/fix-issue-123.feature
trace/fix-issue-123.json
```

### Use the same lifecycle verbs

```bash
./intend lock --mode contrib fix-issue-123
./intend trace --mode contrib fix-issue-123
./intend amend --mode contrib fix-issue-123
./intend delete --mode contrib fix-issue-123
```

The idea is the same. The only difference is where the bundle lives.

If the contribution bundle is locked, add `--force` to delete it.

## Cursor guidance

`intend` can also install managed Cursor guidance for this repository:

```bash
./intend agent install cursor
```

This installs the built-in command and skill files under `.cursor/`.

Re-running the command is safe when those managed files are unchanged. If one was edited locally, `intend` refuses to overwrite it.

## Practical rules

If you only remember five things, remember these:

- Write the spec before the implementation.
- Review the feature contract before locking it.
- Use `trace` to detect drift, not to guess.
- Use `amend` only when the contract changed intentionally.
- Run `verify` before you claim the work is done.

## Next steps

- Read [Reference](../reference/) for command forms and exit codes.
- Go back to [Getting started](../getting-started/) if you want the smallest complete example.
