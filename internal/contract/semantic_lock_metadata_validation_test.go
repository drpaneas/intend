package contract

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func InitializeSemanticLockMetadataValidationScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^a locked bundle `([^`]+)`$", state.aLockedBundle)
	ctx.Step("^a locked contribution bundle `([^`]+)`$", state.aLockedContributionBundle)
	ctx.Step("^verification tools are available$", state.verificationToolsAreAvailable)
	ctx.Step("^I add the lock file semantic digest for `([^`]+)` to `([^`]+)`$", state.iAddTheLockFileSemanticDigestForTo)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
	ctx.Step("^the verification log is empty$", state.theVerificationLogIsEmpty)
}

func TestSemanticLockMetadataValidationFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeSemanticLockMetadataValidationScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "semantic-lock-metadata-validation.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
