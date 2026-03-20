package contract

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func InitializeLockRequiredTraceScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^an initialized owned repository$", state.anInitializedOwnedRepository)
	ctx.Step("^an existing bundle `([^`]+)`$", state.anExistingBundle)
	ctx.Step("^a git repository$", state.aGitRepository)
	ctx.Step("^GitHub issue import is available$", state.gitHubIssueImportIsAvailable)
	ctx.Step("^an existing contribution bundle `([^`]+)`$", state.anExistingContributionBundle)
	ctx.Step("^verification tools are available$", state.verificationToolsAreAvailable)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
	ctx.Step("^the verification log is empty$", state.theVerificationLogIsEmpty)
}

func TestLockRequiredTraceFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeLockRequiredTraceScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "lock-required-trace.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
