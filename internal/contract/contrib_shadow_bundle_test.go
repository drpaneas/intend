package contract

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

const validGitHubIssueJSON = `{"number":123,"title":"Fix crash on startup","body":"The app crashes during startup.","url":"https://github.com/owner/repo/issues/123"}`
const wrongNumberGitHubIssueJSON = `{"number":124,"title":"Fix crash on startup","body":"The app crashes during startup.","url":"https://github.com/owner/repo/issues/124"}`
const wrongRepoURLGitHubIssueJSON = `{"number":123,"title":"Fix crash on startup","body":"The app crashes during startup.","url":"https://github.com/other/repo/issues/123"}`
const nonGitHubURLGitHubIssueJSON = `{"number":123,"title":"Fix crash on startup","body":"The app crashes during startup.","url":"https://gitlab.com/owner/repo/-/issues/123"}`
const nonIssueGitHubURLGitHubIssueJSON = `{"number":123,"title":"Fix crash on startup","body":"The app crashes during startup.","url":"https://github.com/owner/repo/pull/123"}`
const wrongURLNumberGitHubIssueJSON = `{"number":123,"title":"Fix crash on startup","body":"The app crashes during startup.","url":"https://github.com/owner/repo/issues/124"}`
const queryGitHubIssueJSON = `{"number":123,"title":"Fix crash on startup","body":"The app crashes during startup.","url":"https://github.com/owner/repo/issues/123?tab=comments"}`
const fragmentGitHubIssueJSON = `{"number":123,"title":"Fix crash on startup","body":"The app crashes during startup.","url":"https://github.com/owner/repo/issues/123#issuecomment-1"}`
const alternateHostnameGitHubIssueJSON = `{"number":123,"title":"Fix crash on startup","body":"The app crashes during startup.","url":"https://www.github.com/owner/repo/issues/123"}`

func (s *scenarioState) aGitRepository() error {
	if err := s.ensureWorkDir(); err != nil {
		return err
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("locate git: %w", err)
	}

	cmd := exec.Command(gitPath, "init")
	cmd.Dir = s.workDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git init: %w (output=%q)", err, string(output))
	}

	return nil
}

func (s *scenarioState) gitHubIssueImportIsAvailable() error {
	return s.configureGitHubImportOutput(validGitHubIssueJSON, true)
}

func (s *scenarioState) gitHubIssueImportIsUnavailable() error {
	return s.configureGitHubImportOutput("", false)
}

func (s *scenarioState) gitHubIssueImportReturnsMalformedJSON() error {
	return s.configureGitHubImportOutput("{not-json", true)
}

func (s *scenarioState) gitHubIssueImportReturnsIncompleteIssueData() error {
	return s.configureGitHubImportOutput(`{"number":123,"body":"The app crashes during startup."}`, true)
}

func (s *scenarioState) gitHubIssueImportReturnsDifferentIssueNumber() error {
	return s.configureGitHubImportOutput(wrongNumberGitHubIssueJSON, true)
}

func (s *scenarioState) gitHubIssueImportReturnsURLForDifferentRepository() error {
	return s.configureGitHubImportOutput(wrongRepoURLGitHubIssueJSON, true)
}

func (s *scenarioState) gitHubIssueImportReturnsNonGitHubIssueURL() error {
	return s.configureGitHubImportOutput(nonGitHubURLGitHubIssueJSON, true)
}

func (s *scenarioState) gitHubIssueImportReturnsGitHubPullRequestURL() error {
	return s.configureGitHubImportOutput(nonIssueGitHubURLGitHubIssueJSON, true)
}

func (s *scenarioState) gitHubIssueImportReturnsGitHubIssueURLWithDifferentIssueNumber() error {
	return s.configureGitHubImportOutput(wrongURLNumberGitHubIssueJSON, true)
}

func (s *scenarioState) gitHubIssueImportReturnsGitHubIssueURLWithQueryString() error {
	return s.configureGitHubImportOutput(queryGitHubIssueJSON, true)
}

func (s *scenarioState) gitHubIssueImportReturnsGitHubIssueURLWithFragment() error {
	return s.configureGitHubImportOutput(fragmentGitHubIssueJSON, true)
}

func (s *scenarioState) gitHubIssueImportReturnsGitHubIssueURLWithAlternateHostname() error {
	return s.configureGitHubImportOutput(alternateHostnameGitHubIssueJSON, true)
}

func InitializeContribShadowBundleScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^an empty working directory$", state.anEmptyWorkingDirectory)
	ctx.Step("^a git repository$", state.aGitRepository)
	ctx.Step("^GitHub issue import is available$", state.gitHubIssueImportIsAvailable)
	ctx.Step("^GitHub issue import is unavailable$", state.gitHubIssueImportIsUnavailable)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^the file `([^`]+)` exists$", state.theFileExists)
	ctx.Step("^the file `([^`]+)` contains `([^`]+)`$", state.theFileContains)
	ctx.Step("^the path `([^`]+)` does not exist$", state.thePathDoesNotExist)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
}

func TestContribShadowBundleFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeContribShadowBundleScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "contrib-shadow-bundle.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func (s *scenarioState) configureGitHubImportOutput(output string, available bool) error {
	if err := s.ensureWorkDir(); err != nil {
		return err
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("locate git: %w", err)
	}

	binDir := s.absPath(".fake-gh-bin")
	if err := os.RemoveAll(binDir); err != nil {
		return fmt.Errorf("reset fake gh bin dir: %w", err)
	}
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return fmt.Errorf("create fake gh bin dir: %w", err)
	}

	gitWrapper := fmt.Sprintf("#!/bin/sh\nexec %q \"$@\"\n", gitPath)
	if err := os.WriteFile(filepath.Join(binDir, "git"), []byte(gitWrapper), 0o755); err != nil {
		return fmt.Errorf("write git wrapper: %w", err)
	}

	if available {
		ghScript := fmt.Sprintf("#!/bin/sh\nprintf '%%s\\n' %s\n", strconv.Quote(output))
		if err := os.WriteFile(filepath.Join(binDir, "gh"), []byte(ghScript), 0o755); err != nil {
			return fmt.Errorf("write gh wrapper: %w", err)
		}
	}

	s.env = []string{"PATH=" + binDir}
	return nil
}
