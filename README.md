# intend

[![CI](https://github.com/drpaneas/intend/actions/workflows/ci.yml/badge.svg)](https://github.com/drpaneas/intend/actions/workflows/ci.yml)
[![Docs](https://github.com/drpaneas/intend/actions/workflows/docs.yml/badge.svg)](https://github.com/drpaneas/intend/actions/workflows/docs.yml)
[![Release](https://img.shields.io/github/v/release/drpaneas/intend?display_name=tag)](https://github.com/drpaneas/intend/releases)
[![License](https://img.shields.io/github/license/drpaneas/intend)](LICENSE)

`intend` is a Go 1.26 CLI for intent-driven development in Go projects.

```bash
go build ./cmd/intend
go run ./cmd/intend init
go run ./cmd/intend new <name>
go run ./cmd/intend lock <name>
go run ./cmd/intend trace <name>
go run ./cmd/intend amend <name>
go run ./cmd/intend verify
```

Docs: [site](https://drpaneas.github.io/intend/) | [source](docs/)
Contributing: [CONTRIBUTING.md](CONTRIBUTING.md)
