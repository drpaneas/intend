package contract

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func InitializeContribSemanticLockBackwardCompatibilityScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^a locked contribution bundle `([^`]+)`$", state.aLockedContributionBundle)
	ctx.Step("^verification tools are available$", state.verificationToolsAreAvailable)
	ctx.Step("^I remove the lock file semantic digests from `([^`]+)`$", state.iRemoveTheLockFileSemanticDigestsFrom)
	ctx.Step("^I rewrite the issue snapshot JSON with pretty formatting at `([^`]+)`$", state.iRewriteTheIssueSnapshotJSONWithPrettyFormatting)
	ctx.Step("^I rewrite the issue snapshot JSON with sorted keys at `([^`]+)`$", state.iRewriteTheIssueSnapshotJSONWithSortedKeys)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
	ctx.Step("^the contribution lock version for `([^`]+)` is (\\d+)$", state.theContributionLockVersionForIs)
	ctx.Step("^the verification log is empty$", state.theVerificationLogIsEmpty)
}

func TestContribSemanticLockBackwardCompatibilityFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeContribSemanticLockBackwardCompatibilityScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "contrib-semantic-lock-backward-compatibility.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
