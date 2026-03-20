---
name: intend-idd-workflow
description: Use when changing Go code in this repo with intend
---
# Intend IDD Workflow

- Start from the reviewed intent and feature contract.
- Keep the sequence `spec -> feature -> tests -> implementation`.
- If the contract changes, update the spec and feature before code.
- Run `intend trace` when contract files change.
- Run `intend verify` after implementation changes.
