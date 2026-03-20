package contract

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func InitializeBundleDeleteScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^an existing bundle `([^`]+)`$", state.anExistingBundle)
	ctx.Step("^a locked bundle `([^`]+)`$", state.aLockedBundle)
	ctx.Step("^an existing contribution bundle `([^`]+)`$", state.anExistingContributionBundle)
	ctx.Step("^verification tools are available$", state.verificationToolsAreAvailable)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^stdout contains `([^`]+)`$", state.stdoutContains)
	ctx.Step("^stdout does not contain `([^`]+)`$", state.stdoutDoesNotContain)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
	ctx.Step("^the file `([^`]+)` exists$", state.theFileExists)
	ctx.Step("^the path `([^`]+)` does not exist$", state.thePathDoesNotExist)
	ctx.Step("^the verification log contains lines in order$", state.theVerificationLogContainsLinesInOrder)
}

func TestBundleDeleteFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeBundleDeleteScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "bundle-delete.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
