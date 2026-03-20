package contract

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func InitializeAmendIdempotencyScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^an initialized owned repository$", state.anInitializedOwnedRepository)
	ctx.Step("^a locked bundle `([^`]+)`$", state.aLockedBundle)
	ctx.Step("^a git repository$", state.aGitRepository)
	ctx.Step("^GitHub issue import is available$", state.gitHubIssueImportIsAvailable)
	ctx.Step("^a locked contribution bundle `([^`]+)`$", state.aLockedContributionBundle)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^stdout contains `([^`]+)`$", state.stdoutContains)
	ctx.Step("^the lock version for `([^`]+)` is (\\d+)$", state.theLockVersionForIs)
	ctx.Step("^the contribution lock version for `([^`]+)` is (\\d+)$", state.theContributionLockVersionForIs)
	ctx.Step("^`intend trace ([^` ]+)` succeeds$", state.intendTraceSucceeds)
	ctx.Step("^`intend trace --mode ([^` ]+) ([^`]+)` succeeds$", state.intendTraceModeSucceeds)
}

func TestAmendIdempotencyFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeAmendIdempotencyScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "amend-idempotency.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
