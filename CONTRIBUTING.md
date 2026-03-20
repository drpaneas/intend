# Contributing to `intend`

Thanks for the help.

Keep changes small. Prefer the order `spec -> feature -> tests -> implementation`.

## Local checks

Run the core checks before opening a pull request:

```bash
go test ./... -count=1
golangci-lint run
go build ./...
```

If you change the docs site, make sure Hugo can rebuild it:

```bash
cd docs
hugo
```

## Working style

- Explain the change in plain language.
- Add or update tests when behavior changes.
- Update docs when user-facing behavior changes.
- Keep bundle names in lowercase kebab-case.

## Contract-driven changes

When a change affects contract behavior, use `intend` itself:

```bash
go run ./cmd/intend new <name>
go run ./cmd/intend lock <name>
go run ./cmd/intend trace <name>
go run ./cmd/intend amend <name>
go run ./cmd/intend verify
```
