---
title: "Intent-Driven Development"
weight: 10
summary: "Why IDD exists, what problem it solves, and how intend turns it into a practical workflow."
---

Code generation got cheap. Correctness did not.

That is the problem Intent-Driven Development tries to solve.

An AI model can produce a lot of code very quickly. What it cannot do reliably is decide, on its own, whether the code matches the change you actually wanted. That part still needs a contract.

## The problem

The default AI coding loop usually looks like this:

```text
describe the change -> ask for code -> review the diff -> run tests -> hope nothing important was missed
```

This works for small changes. It breaks down when the prompt is vague, the edge cases are not written down, or the code looks plausible enough to slip past review.

By the time you are staring at a large diff, the cheapest moment to fix the misunderstanding is already gone.

## The core idea

IDD inserts one explicit step before implementation: define how the change will be judged before the code is written.

For `intend`, the working order is:

```text
spec -> feature -> tests -> implementation
```

That order matters.

The spec explains the intent. The feature file explains the expected behavior. The tests make that behavior executable. Only then should implementation begin.

## The contract in files

`intend` keeps that contract in plain files:

| File | Purpose |
| --- | --- |
| `specs/<name>.md` | What the change is for and what it should do |
| `features/<name>.feature` | The human-readable behavior contract |
| `.intend/trace/<name>.json` | The mapping between the bundle name and its tracked files |
| `.intend/locks/<name>.json` | The locked baseline used to detect drift |

You can read every part of it. There is no hidden database. There is no silent state in an editor plugin.

## Why this helps

IDD separates two things that often get mixed together:

- **Intent**: what do we want?
- **Acceptance**: what would prove that we got it?

Without that split, code review turns into archaeology. Reviewers have to infer the real target from the implementation. With IDD, the target is already written down.

This changes how AI fits into the process:

- the human approves the contract
- the tool records and checks the contract
- the model writes candidate implementation inside those boundaries

That is a better division of labor.

## The workflow

In practice, the loop is short.

1. Write or edit `specs/<name>.md`.
2. Write or review `features/<name>.feature`.
3. Lock the approved contract with `intend lock <name>`.
4. Implement the code and tests.
5. Run `intend trace <name>` to confirm the locked contract still matches.
6. Run `intend verify` to trace all bundles and run the repository checks.

If the contract itself must change, do not quietly rewrite it after the fact. Update it deliberately, then record that new baseline with `intend amend <name>`.

That is the key discipline. The definition of "correct" should not drift to match whatever code happened to be produced.

## What `intend` adds

`intend` does not invent testing. It does not replace review. It does not replace `go test`.

What it adds is a small amount of structure:

- a named bundle for each change
- a fixed place for intent and behavior
- a locked baseline
- a simple drift check
- a repository-level verification command

The tool stays small on purpose. It uses normal files and normal command-line tools.

## What IDD is not

IDD is not bureaucracy for its own sake.

It is not a demand that every idea become a giant specification. It is not a promise that feature files will magically remove human judgment. And it is not limited to AI-assisted work, though AI makes the need more obvious.

The point is simpler than that:

Write down what the change means before you let the code define it for you.

## A small example

A new bundle called `hello-http` starts like this:

```text
specs/hello-http.md
features/hello-http.feature
.intend/trace/hello-http.json
```

After review, `intend lock hello-http` adds:

```text
.intend/locks/hello-http.json
```

From that point on, `intend trace hello-http` can tell you whether the contract drifted.

## Next step

If the idea makes sense, continue to [Getting started](../getting-started/).

If you want the day-to-day command flow, read [Workflow](../workflow/).
