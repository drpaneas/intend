# Require locks before trace and verify

Harden the narrowed v1 CLI so draft bundles are not treated as approved contracts.

This slice closes the semantic gap where a bundle can exist without any lock metadata but still be treated as traceable or verifiable. In v1, `trace` and `verify` should only accept locked bundles as evidence-bearing contracts.

## Commands in scope

- `intend trace <name>`
- `intend trace --mode contrib <name>`
- `intend verify`

## Required behavior

- If an owned bundle exists but has no lock file, `intend trace <name>` exits with code `1`.
- If a contribution shadow bundle exists but has no lock file, `intend trace --mode contrib <name>` exits with code `1`.
- If any discovered bundle is unlocked, `intend verify` exits before running external verification tools.
- Error output should clearly state that the contract is not locked.

## Constraints

- Keep the behavior file-based.
- Do not infer approval from the presence of a spec, feature, or trace file alone.
- Do not introduce a special draft mode in this slice.

## CLI rules

- Exit `1` for unlocked bundles in `trace` and `verify`.
- Existing success and usage behavior stays unchanged.

## Done when

- The feature contract for lock-required tracing passes through `godog`.
- Unlocked owned and contribution bundles are rejected consistently by `trace`, and `verify` stops before external tools when it discovers one.
