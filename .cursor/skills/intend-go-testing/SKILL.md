---
name: intend-go-testing
description: Use for Go testing, context, and safety standards in this repo
---
# Intend Go Testing

- Write or update tests before implementation when behavior changes.
- Keep tests tied to the approved feature contract.
- Use contexts deliberately and propagate cancellation.
- Avoid logging secrets, tokens, or raw credentials.
- Prefer deterministic tests over timing-sensitive behavior.
