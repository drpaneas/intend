# Reject invalid contribution issue snapshot URLs across all lifecycle phases

Clarify how `intend` validates the contribution `issue.json` URL when the snapshot no longer uses a canonical GitHub issue URL shape.

This slice consolidates the current URL-validation behavior for contribution issue snapshots across trace, amend, verify, and initial lock. The snapshot may still decode correctly, but its URL is no longer canonical because it uses a non-GitHub host, an unsupported GitHub hostname, a non-issue path shape, or a non-canonical query string or fragment. Every lifecycle command that consumes the snapshot must reject those URL errors explicitly.

## Commands in scope

- `intend trace --mode contrib <name>`
- `intend amend --mode contrib <name>`
- `intend verify`
- `intend lock --mode contrib <name>`

## Required behavior

- If the contribution issue snapshot URL is not on `github.com`, each command exits with code `1`.
- If the contribution issue snapshot URL uses an alternate GitHub hostname, each command exits with code `1`.
- If the contribution issue snapshot URL is not in the canonical `/owner/repo/issues/<n>` shape, each command exits with code `1`.
- If the contribution issue snapshot URL includes a query string or fragment, each command exits with code `1`.
- `intend verify` stops before running verification tools.
- `intend lock --mode contrib <name>` leaves no lock file behind on failure.
- The failure message reports the specific URL problem.

## Constraints

- Keep the validation local and file-based.
- Preserve the existing number and repository mismatch behavior, which is covered elsewhere.
- Do not reinterpret invalid URLs as generic contract drift.

## Done when

- The feature contract for contribution issue snapshot URL validation matrix passes through `godog`.
- Contribution issue snapshot URL validation is explicitly covered across trace, amend, verify, and initial lock.
