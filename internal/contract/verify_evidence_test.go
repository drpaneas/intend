package contract

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func (s *scenarioState) verificationToolsAreAvailable() error {
	return s.installFakeVerificationTools("")
}

func (s *scenarioState) verificationToolsAreAvailableExcept(missing string) error {
	return s.installFakeVerificationTools(missing)
}

func (s *scenarioState) theVerificationLogContainsLinesInOrder(table *godog.Table) error {
	logLines, err := s.readVerificationLog()
	if err != nil {
		return err
	}

	if len(logLines) != len(table.Rows) {
		return fmt.Errorf("expected %d verification log lines, got %d (%v)", len(table.Rows), len(logLines), logLines)
	}

	for i, row := range table.Rows {
		if len(row.Cells) != 1 {
			return fmt.Errorf("expected single-cell table row, got %d cells", len(row.Cells))
		}

		want := row.Cells[0].Value
		if logLines[i] != want {
			return fmt.Errorf("expected log line %d to be %q, got %q", i, want, logLines[i])
		}
	}

	return nil
}

func (s *scenarioState) theVerificationLogIsEmpty() error {
	logLines, err := s.readVerificationLog()
	if err != nil {
		return err
	}

	if len(logLines) != 0 {
		return fmt.Errorf("expected empty verification log, got %v", logLines)
	}

	return nil
}

func InitializeVerifyScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^an empty working directory$", state.anEmptyWorkingDirectory)
	ctx.Step("^an initialized owned repository$", state.anInitializedOwnedRepository)
	ctx.Step("^a git repository$", state.aGitRepository)
	ctx.Step("^an existing bundle `([^`]+)`$", state.anExistingBundle)
	ctx.Step("^a locked bundle `([^`]+)`$", state.aLockedBundle)
	ctx.Step("^GitHub issue import is available$", state.gitHubIssueImportIsAvailable)
	ctx.Step("^a locked contribution bundle `([^`]+)`$", state.aLockedContributionBundle)
	ctx.Step("^verification tools are available$", state.verificationToolsAreAvailable)
	ctx.Step("^verification tools are available except `([^`]+)`$", state.verificationToolsAreAvailableExcept)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^I replace the contents of `([^`]+)`$", state.iReplaceTheContentsOf)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^stdout contains `([^`]+)`$", state.stdoutContains)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
	ctx.Step("^the verification log contains lines in order$", state.theVerificationLogContainsLinesInOrder)
	ctx.Step("^the verification log is empty$", state.theVerificationLogIsEmpty)
}

func TestVerifyEvidenceFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeVerifyScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "verify-evidence.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func (s *scenarioState) installFakeVerificationTools(missing string) error {
	if err := s.ensureWorkDir(); err != nil {
		return err
	}

	binDir := s.absPath(".fake-bin")
	if err := os.RemoveAll(binDir); err != nil {
		return fmt.Errorf("reset fake bin dir: %w", err)
	}
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return fmt.Errorf("create fake bin dir: %w", err)
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("locate git: %w", err)
	}

	gitWrapper := fmt.Sprintf("#!/bin/sh\nexec %q \"$@\"\n", gitPath)
	if err := os.WriteFile(filepath.Join(binDir, "git"), []byte(gitWrapper), 0o755); err != nil {
		return fmt.Errorf("write git wrapper: %w", err)
	}

	logPath := s.verificationLogPath()
	s.env = []string{
		"PATH=" + binDir,
		"INTEND_VERIFY_LOG=" + logPath,
	}

	script := "#!/bin/sh\n" +
		"tool=${0##*/}\n" +
		"printf '%s %s\\n' \"$tool\" \"$*\" >> \"$INTEND_VERIFY_LOG\"\n" +
		"exit 0\n"

	for _, tool := range []string{"go", "golangci-lint", "trufflehog", "gitleaks", "trivy"} {
		if tool == missing {
			continue
		}

		path := filepath.Join(binDir, tool)
		if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
			return fmt.Errorf("write fake tool %s: %w", tool, err)
		}
	}

	return nil
}

func (s *scenarioState) verificationLogPath() string {
	return s.absPath(".verification.log")
}

func (s *scenarioState) readVerificationLog() ([]string, error) {
	data, err := os.ReadFile(s.verificationLogPath())
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read verification log: %w", err)
	}

	content := strings.TrimSpace(string(data))
	if content == "" {
		return nil, nil
	}

	return strings.Split(content, "\n"), nil
}
