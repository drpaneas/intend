package contract

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func InitializeContribImportURLValidationMatrixScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^a git repository$", state.aGitRepository)
	ctx.Step("^GitHub issue import returns a non-GitHub issue URL$", state.gitHubIssueImportReturnsNonGitHubIssueURL)
	ctx.Step("^GitHub issue import returns a GitHub issue URL with an alternate hostname$", state.gitHubIssueImportReturnsGitHubIssueURLWithAlternateHostname)
	ctx.Step("^GitHub issue import returns a GitHub pull request URL$", state.gitHubIssueImportReturnsGitHubPullRequestURL)
	ctx.Step("^GitHub issue import returns a GitHub issue URL with a query string$", state.gitHubIssueImportReturnsGitHubIssueURLWithQueryString)
	ctx.Step("^GitHub issue import returns a GitHub issue URL with a fragment$", state.gitHubIssueImportReturnsGitHubIssueURLWithFragment)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
	ctx.Step("^the path `([^`]+)` does not exist$", state.thePathDoesNotExist)
}

func TestContribImportURLValidationMatrixFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeContribImportURLValidationMatrixScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "contrib-import-url-validation-matrix.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
