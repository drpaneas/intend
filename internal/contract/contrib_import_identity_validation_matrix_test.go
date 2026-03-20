package contract

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func InitializeContribImportIdentityValidationMatrixScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^a git repository$", state.aGitRepository)
	ctx.Step("^GitHub issue import returns a different issue number$", state.gitHubIssueImportReturnsDifferentIssueNumber)
	ctx.Step("^GitHub issue import returns a URL for a different repository$", state.gitHubIssueImportReturnsURLForDifferentRepository)
	ctx.Step("^GitHub issue import returns a GitHub issue URL with a different issue number$", state.gitHubIssueImportReturnsGitHubIssueURLWithDifferentIssueNumber)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
	ctx.Step("^the path `([^`]+)` does not exist$", state.thePathDoesNotExist)
}

func TestContribImportIdentityValidationMatrixFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeContribImportIdentityValidationMatrixScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "contrib-import-identity-validation-matrix.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
