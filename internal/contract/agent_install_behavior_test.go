package contract

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func InitializeAgentInstallBehaviorScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^an initialized owned repository$", state.anInitializedOwnedRepository)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^the directory `([^`]+)` exists$", state.theDirectoryExists)
	ctx.Step("^the file `([^`]+)` exists$", state.theFileExists)
	ctx.Step("^the file `([^`]+)` contains `([^`]+)`$", state.theFileContains)
	ctx.Step("^I replace the contents of `([^`]+)`$", state.iReplaceTheContentsOf)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
}

func TestAgentInstallBehaviorFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeAgentInstallBehaviorScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "agent-install-behavior.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
