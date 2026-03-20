package contract

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func InitializeContribIssueSnapshotEditBehaviorScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^an existing contribution bundle `([^`]+)`$", state.anExistingContributionBundle)
	ctx.Step("^a locked contribution bundle `([^`]+)`$", state.aLockedContributionBundle)
	ctx.Step("^verification tools are available$", state.verificationToolsAreAvailable)
	ctx.Step("^I replace the issue snapshot file `([^`]+)` so `([^`]+)` becomes `([^`]+)`$", state.iReplaceTheIssueSnapshotFileFieldWith)
	ctx.Step("^I add the unknown issue snapshot field `([^`]+)` with value `([^`]+)` to `([^`]+)`$", state.iAddUnknownIssueSnapshotFieldWithValueTo)
	ctx.Step("^I rewrite the issue snapshot JSON with pretty formatting at `([^`]+)`$", state.iRewriteTheIssueSnapshotJSONWithPrettyFormatting)
	ctx.Step("^I rewrite the issue snapshot JSON with sorted keys at `([^`]+)`$", state.iRewriteTheIssueSnapshotJSONWithSortedKeys)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^stdout contains `([^`]+)`$", state.stdoutContains)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
	ctx.Step("^the contribution lock version for `([^`]+)` is (\\d+)$", state.theContributionLockVersionForIs)
	ctx.Step("^the file `([^`]+)` exists$", state.theFileExists)
	ctx.Step("^`intend trace --mode ([^` ]+) ([^`]+)` succeeds$", state.intendTraceModeSucceeds)
	ctx.Step("^the verification log contains lines in order:$", state.theVerificationLogContainsLinesInOrder)
	ctx.Step("^the verification log is empty$", state.theVerificationLogIsEmpty)
}

func TestContribIssueSnapshotEditBehaviorFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeContribIssueSnapshotEditBehaviorScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "contrib-issue-snapshot-edit-behavior.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
