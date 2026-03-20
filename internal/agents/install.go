package agents

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type managedFile struct {
	relativePath string
	content      string
}

var registry = map[string][]managedFile{
	"cursor": {
		{
			relativePath: filepath.ToSlash(filepath.Join(".cursor", "commands", "intend-workflow.md")),
			content: `# Intend Workflow

Use ` + "`intend`" + ` to keep changes contract-driven in this repository.

- Keep the order ` + "`spec -> feature -> tests -> implementation`" + `.
- Use ` + "`intend new <name>`" + ` to create a bundle.
- Use ` + "`intend lock <name>`" + ` after contract approval.
- Use ` + "`intend amend <name>`" + ` when the contract intentionally changes.
- Use ` + "`intend trace <name>`" + ` after contract edits.
- Use ` + "`intend verify`" + ` before claiming implementation is complete.
`,
		},
		{
			relativePath: filepath.ToSlash(filepath.Join(".cursor", "skills", "intend-idd-workflow", "SKILL.md")),
			content: `---
name: intend-idd-workflow
description: Use when changing Go code in this repo with intend
---
# Intend IDD Workflow

- Start from the reviewed intent and feature contract.
- Keep the sequence ` + "`spec -> feature -> tests -> implementation`" + `.
- If the contract changes, update the spec and feature before code.
- Run ` + "`intend trace`" + ` when contract files change.
- Run ` + "`intend verify`" + ` after implementation changes.
`,
		},
		{
			relativePath: filepath.ToSlash(filepath.Join(".cursor", "skills", "intend-go-core", "SKILL.md")),
			content: `---
name: intend-go-core
description: Use for core Go coding standards in this repo
---
# Intend Go Core

- Keep packages small and names literal.
- Return explicit errors; do not hide failures.
- Prefer simple structs and functions over speculative abstractions.
- Keep public APIs small and file-based where possible.
- Make changes traceable to the current contract bundle.
`,
		},
		{
			relativePath: filepath.ToSlash(filepath.Join(".cursor", "skills", "intend-go-testing", "SKILL.md")),
			content: `---
name: intend-go-testing
description: Use for Go testing, context, and safety standards in this repo
---
# Intend Go Testing

- Write or update tests before implementation when behavior changes.
- Keep tests tied to the approved feature contract.
- Use contexts deliberately and propagate cancellation.
- Avoid logging secrets, tokens, or raw credentials.
- Prefer deterministic tests over timing-sensitive behavior.
`,
		},
	},
}

func Install(root, agent string) error {
	files, ok := registry[agent]
	if !ok {
		return fmt.Errorf("unsupported agent: %s", agent)
	}

	if err := ensureManagedFilesUnchanged(root, files); err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(root, filepath.FromSlash(file.relativePath))
		if _, err := os.Stat(path); err == nil {
			continue
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat %s: %w", file.relativePath, err)
		}

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("create parent directory for %s: %w", file.relativePath, err)
		}

		if err := os.WriteFile(path, []byte(file.content), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", file.relativePath, err)
		}
	}

	return nil
}

func ensureManagedFilesUnchanged(root string, files []managedFile) error {
	for _, file := range files {
		path := filepath.Join(root, filepath.FromSlash(file.relativePath))
		data, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return fmt.Errorf("read %s: %w", file.relativePath, err)
		}
		if string(data) != file.content {
			return fmt.Errorf("managed file was modified: %s", file.relativePath)
		}
	}

	return nil
}
