package verify

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"intend/internal/workflow"
)

type toolCommand struct {
	name string
	args []string
}

var defaultCommands = []toolCommand{
	{name: "go", args: []string{"test", "./..."}},
	{name: "golangci-lint", args: []string{"run"}},
	{name: "trufflehog", args: []string{"filesystem", "."}},
	{name: "gitleaks", args: []string{"dir", "."}},
	{name: "trivy", args: []string{"fs", "."}},
}

func CheckRequiredTools() error {
	for _, command := range defaultCommands {
		if _, err := lookupTool(command.name); err != nil {
			return err
		}
	}

	return nil
}

func Run(root string, stdout io.Writer) error {
	refs, err := workflow.ListBundleRefs(root)
	if err != nil {
		return err
	}

	for _, ref := range refs {
		if _, err := fmt.Fprintf(stdout, "trace: checking %s bundle %s\n", ref.Mode, ref.Name); err != nil {
			return err
		}
		if err := workflow.TraceBundleWithMode(root, ref.Mode, ref.Name); err != nil {
			return err
		}
	}

	for _, command := range defaultCommands {
		if _, err := fmt.Fprintf(stdout, "verify: running %s\n", formatCommand(command)); err != nil {
			return err
		}
		if err := runCommand(root, command); err != nil {
			return err
		}
	}

	return nil
}

func runCommand(root string, command toolCommand) error {
	path, err := lookupTool(command.name)
	if err != nil {
		return err
	}

	cmd := exec.Command(path, command.args...)
	cmd.Dir = root

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w (stdout=%q, stderr=%q)", command.name, err, stdout.String(), stderr.String())
	}

	return nil
}

func lookupTool(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("required tool not found: %s", name)
	}

	return path, nil
}

func formatCommand(command toolCommand) string {
	return strings.TrimSpace(command.name + " " + strings.Join(command.args, " "))
}
