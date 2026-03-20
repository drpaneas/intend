package contract

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func (s *scenarioState) stdoutContainsLinesInOrder(table *godog.Table) error {
	output := strings.TrimSpace(s.stdout)
	if output == "" {
		return errors.New("stdout is empty")
	}

	lines := strings.Split(output, "\n")
	if len(lines) < len(table.Rows) {
		return fmt.Errorf("expected at least %d stdout lines, got %d (%v)", len(table.Rows), len(lines), lines)
	}

	for i, row := range table.Rows {
		if len(row.Cells) != 1 {
			return fmt.Errorf("expected single-cell table row, got %d cells", len(row.Cells))
		}

		want := row.Cells[0].Value
		if lines[i] != want {
			return fmt.Errorf("expected stdout line %d to be %q, got %q", i, want, lines[i])
		}
	}

	return nil
}

func InitializeCliUxScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^an initialized owned repository$", state.anInitializedOwnedRepository)
	ctx.Step("^a locked bundle `([^`]+)`$", state.aLockedBundle)
	ctx.Step("^verification tools are available$", state.verificationToolsAreAvailable)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^stdout contains lines in order$", state.stdoutContainsLinesInOrder)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
}

func TestCliUxReportingFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeCliUxScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "cli-ux-reporting.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
