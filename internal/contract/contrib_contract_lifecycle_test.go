package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func (s *scenarioState) anExistingContributionBundle(name string) error {
	if err := s.aGitRepository(); err != nil {
		return err
	}
	if err := s.gitHubIssueImportIsAvailable(); err != nil {
		return err
	}

	return s.runAndRequireSuccess("new", "--mode", "contrib", "--from-gh", "owner/repo#123", name)
}

func (s *scenarioState) aLockedContributionBundle(name string) error {
	if err := s.anExistingContributionBundle(name); err != nil {
		return err
	}

	return s.runAndRequireSuccess("lock", "--mode", "contrib", name)
}

func (s *scenarioState) intendTraceModeSucceeds(mode, name string) error {
	if err := s.runIntend("trace", "--mode", mode, name); err != nil {
		return err
	}

	if s.exitCode != 0 {
		return fmt.Errorf("expected intend trace --mode %s %s to succeed, got exit code %d (stdout=%q, stderr=%q)", mode, name, s.exitCode, s.stdout, s.stderr)
	}

	return nil
}

func (s *scenarioState) theContributionLockVersionForIs(name string, expected int) error {
	data, err := os.ReadFile(s.absPath(filepath.Join(".git", "intend", "contrib", name, "locks", name+".json")))
	if err != nil {
		return fmt.Errorf("read contribution lock file for %s: %w", name, err)
	}

	var lock lockState
	if err := json.Unmarshal(data, &lock); err != nil {
		return fmt.Errorf("decode contribution lock file for %s: %w", name, err)
	}

	if lock.Version != expected {
		return fmt.Errorf("expected contribution lock version %d, got %d", expected, lock.Version)
	}

	return nil
}

func InitializeContribContractLifecycleScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^a git repository$", state.aGitRepository)
	ctx.Step("^GitHub issue import is available$", state.gitHubIssueImportIsAvailable)
	ctx.Step("^an existing contribution bundle `([^`]+)`$", state.anExistingContributionBundle)
	ctx.Step("^a locked contribution bundle `([^`]+)`$", state.aLockedContributionBundle)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^I replace the contents of `([^`]+)`$", state.iReplaceTheContentsOf)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^the file `([^`]+)` exists$", state.theFileExists)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
	ctx.Step("^`intend trace --mode ([^` ]+) ([^`]+)` succeeds$", state.intendTraceModeSucceeds)
	ctx.Step("^the contribution lock version for `([^`]+)` is (\\d+)$", state.theContributionLockVersionForIs)
}

func TestContribContractLifecycleFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeContribContractLifecycleScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "contrib-contract-lifecycle.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
