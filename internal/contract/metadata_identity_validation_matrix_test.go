package contract

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func InitializeMetadataIdentityValidationMatrixScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^an existing bundle `([^`]+)`$", state.anExistingBundle)
	ctx.Step("^an existing contribution bundle `([^`]+)`$", state.anExistingContributionBundle)
	ctx.Step("^a locked bundle `([^`]+)`$", state.aLockedBundle)
	ctx.Step("^a locked contribution bundle `([^`]+)`$", state.aLockedContributionBundle)
	ctx.Step("^verification tools are available$", state.verificationToolsAreAvailable)
	ctx.Step("^I replace the trace file `([^`]+)` so `([^`]+)` becomes `([^`]+)`$", state.iReplaceTheTraceFileFieldWith)
	ctx.Step("^I replace the lock file `([^`]+)` so `([^`]+)` becomes `([^`]+)`$", state.iReplaceTheLockFileFieldWith)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
	ctx.Step("^the verification log is empty$", state.theVerificationLogIsEmpty)
	ctx.Step("^the path `([^`]+)` does not exist$", state.thePathDoesNotExist)
}

func TestMetadataIdentityValidationMatrixFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeMetadataIdentityValidationMatrixScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "metadata-identity-validation-matrix.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
