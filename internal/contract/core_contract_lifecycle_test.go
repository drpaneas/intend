package contract

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"intend/internal/workflow"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

var intendRepoRoot = func() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}()

var intendBinaryPath string

type lockState struct {
	Version int `json:"version"`
}

type scenarioState struct {
	workDir    string
	stdout     string
	stderr     string
	exitCode   int
	commandRan bool
	env        []string
	extraPaths []string
}

func (s *scenarioState) reset(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	if err := s.removeWorkDir(); err != nil {
		return ctx, err
	}

	s.stdout = ""
	s.stderr = ""
	s.exitCode = 0
	s.commandRan = false
	s.env = nil
	s.extraPaths = nil
	return ctx, nil
}

func (s *scenarioState) cleanup(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
	return ctx, s.removeWorkDir()
}

func (s *scenarioState) anEmptyWorkingDirectory() error {
	return s.ensureWorkDir()
}

func (s *scenarioState) anInitializedOwnedRepository() error {
	if err := s.ensureInitToolsAvailable(); err != nil {
		return err
	}

	return s.runAndRequireSuccess("init")
}

func (s *scenarioState) anExistingBundle(name string) error {
	if err := s.anInitializedOwnedRepository(); err != nil {
		return err
	}

	return s.runAndRequireSuccess("new", name)
}

func (s *scenarioState) aLockedBundle(name string) error {
	if err := s.anExistingBundle(name); err != nil {
		return err
	}

	return s.runAndRequireSuccess("lock", name)
}

func (s *scenarioState) iRun(commandLine string) error {
	fields := strings.Fields(commandLine)
	if len(fields) == 0 {
		return errors.New("empty command line")
	}

	if fields[0] != "intend" {
		return fmt.Errorf("expected command to start with intend, got %q", fields[0])
	}

	if len(fields) > 1 && fields[1] == "init" {
		if err := s.ensureInitToolsAvailable(); err != nil {
			return err
		}
	}

	return s.runIntend(fields[1:]...)
}

func (s *scenarioState) itExitsWithCode(expected int) error {
	if !s.commandRan {
		return errors.New("intend command was not run")
	}

	if s.exitCode != expected {
		return fmt.Errorf("expected exit code %d, got %d (stdout=%q, stderr=%q)", expected, s.exitCode, s.stdout, s.stderr)
	}

	return nil
}

func (s *scenarioState) theFileExists(relativePath string) error {
	info, err := os.Stat(s.absPath(relativePath))
	if err != nil {
		return fmt.Errorf("stat %s: %w", relativePath, err)
	}

	if info.IsDir() {
		return fmt.Errorf("expected file %s, got directory", relativePath)
	}

	return nil
}

func (s *scenarioState) theFileContains(relativePath, want string) error {
	data, err := os.ReadFile(s.absPath(relativePath))
	if err != nil {
		return fmt.Errorf("read %s: %w", relativePath, err)
	}

	if !strings.Contains(string(data), want) {
		return fmt.Errorf("expected %s to contain %q, got %q", relativePath, want, string(data))
	}

	return nil
}

func (s *scenarioState) theDirectoryExists(relativePath string) error {
	info, err := os.Stat(s.absPath(relativePath))
	if err != nil {
		return fmt.Errorf("stat %s: %w", relativePath, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("expected directory %s, got file", relativePath)
	}

	return nil
}

func (s *scenarioState) thePathDoesNotExist(relativePath string) error {
	_, err := os.Stat(s.absPath(relativePath))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("stat %s: %w", relativePath, err)
	}

	return fmt.Errorf("expected %s to not exist", relativePath)
}

func (s *scenarioState) iRemoveThePath(relativePath string) error {
	path := s.absPath(relativePath)
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("stat %s: %w", relativePath, err)
	}

	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("remove %s: %w", relativePath, err)
	}

	return s.thePathDoesNotExist(relativePath)
}

func (s *scenarioState) iReplaceThePathWithASymlinkToAnExternalCopyPreservingItsContents(relativePath string) error {
	path := s.absPath(relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", relativePath, err)
	}

	externalDir, err := os.MkdirTemp("", "intend-external-*")
	if err != nil {
		return fmt.Errorf("create external temp dir: %w", err)
	}
	s.extraPaths = append(s.extraPaths, externalDir)

	externalPath := filepath.Join(externalDir, filepath.Base(path))
	if err := os.WriteFile(externalPath, data, 0o644); err != nil {
		return fmt.Errorf("write external copy for %s: %w", relativePath, err)
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("remove %s: %w", relativePath, err)
	}
	if err := os.Symlink(externalPath, path); err != nil {
		return fmt.Errorf("symlink %s to external copy: %w", relativePath, err)
	}

	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("lstat %s: %w", relativePath, err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("expected %s to be a symlink", relativePath)
	}

	return nil
}

func (s *scenarioState) iReplaceTheContentsOf(relativePath string) error {
	path := s.absPath(relativePath)
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("stat %s: %w", relativePath, err)
	}

	content := fmt.Sprintf("changed %s\n", relativePath)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relativePath, err)
	}

	return nil
}

func (s *scenarioState) iReplaceTheTraceFileWithValidJSONMissingRequiredPaths(relativePath string) error {
	path := s.absPath(relativePath)
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("stat %s: %w", relativePath, err)
	}

	content := `{"name":"broken-trace","mode":"owned"}` + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relativePath, err)
	}

	return nil
}

func (s *scenarioState) iReplaceTheTraceFileFieldWith(relativePath, field, value string) error {
	path := s.absPath(relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", relativePath, err)
	}

	var trace workflow.BundleTrace
	if err := json.Unmarshal(data, &trace); err != nil {
		return fmt.Errorf("decode %s: %w", relativePath, err)
	}

	if value == "<empty>" {
		value = ""
	}

	switch field {
	case "name":
		trace.Name = value
	case "mode":
		trace.Mode = value
	case "specPath":
		trace.SpecPath = value
	case "featurePath":
		trace.FeaturePath = value
	case "issueRef":
		trace.IssueRef = value
	case "issuePath":
		trace.IssuePath = value
	default:
		return fmt.Errorf("unsupported trace field %q", field)
	}

	updated, err := json.MarshalIndent(trace, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", relativePath, err)
	}
	updated = append(updated, '\n')

	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relativePath, err)
	}

	return nil
}

func (s *scenarioState) iReplaceTheIssueSnapshotFileFieldWith(relativePath, field, value string) error {
	path := s.absPath(relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", relativePath, err)
	}

	var snapshot struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		Body   string `json:"body"`
		URL    string `json:"url"`
	}
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("decode %s: %w", relativePath, err)
	}

	if value == "<empty>" {
		value = ""
	}

	switch field {
	case "number":
		number, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid issue snapshot number %q", value)
		}
		snapshot.Number = number
	case "title":
		snapshot.Title = value
	case "body":
		snapshot.Body = value
	case "url":
		snapshot.URL = value
	default:
		return fmt.Errorf("unsupported issue snapshot field %q", field)
	}

	updated, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", relativePath, err)
	}
	updated = append(updated, '\n')

	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relativePath, err)
	}

	return nil
}

func (s *scenarioState) iAddUnknownIssueSnapshotFieldWithValueTo(field, value, relativePath string) error {
	path := s.absPath(relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", relativePath, err)
	}

	var snapshot map[string]any
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("decode %s: %w", relativePath, err)
	}

	switch field {
	case "number", "title", "body", "url":
		return fmt.Errorf("issue snapshot field %q is reserved", field)
	}
	if _, exists := snapshot[field]; exists {
		return fmt.Errorf("issue snapshot field %q already exists", field)
	}

	if value == "<empty>" {
		value = ""
	}
	snapshot[field] = value

	updated, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", relativePath, err)
	}
	updated = append(updated, '\n')

	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relativePath, err)
	}

	return nil
}

func (s *scenarioState) iRewriteTheIssueSnapshotJSONWithPrettyFormatting(relativePath string) error {
	path := s.absPath(relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", relativePath, err)
	}

	var snapshot struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		Body   string `json:"body"`
		URL    string `json:"url"`
	}
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("decode %s: %w", relativePath, err)
	}

	updated, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", relativePath, err)
	}
	updated = append(updated, '\n')

	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relativePath, err)
	}

	return nil
}

func (s *scenarioState) iRewriteTheIssueSnapshotJSONWithSortedKeys(relativePath string) error {
	path := s.absPath(relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", relativePath, err)
	}

	var snapshot map[string]any
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("decode %s: %w", relativePath, err)
	}

	updated, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", relativePath, err)
	}
	updated = append(updated, '\n')

	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relativePath, err)
	}

	return nil
}

func (s *scenarioState) iReplaceTheLockFileFieldWith(relativePath, field, value string) error {
	path := s.absPath(relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", relativePath, err)
	}

	var lock workflow.BundleLock
	if err := json.Unmarshal(data, &lock); err != nil {
		return fmt.Errorf("decode %s: %w", relativePath, err)
	}

	switch field {
	case "name":
		lock.Name = value
	default:
		return fmt.Errorf("unsupported lock field %q", field)
	}

	updated, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", relativePath, err)
	}
	updated = append(updated, '\n')

	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relativePath, err)
	}

	return nil
}

func (s *scenarioState) iReplaceTheLockFileTrackedPathWith(relativePath, oldPath, newPath string) error {
	path := s.absPath(relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", relativePath, err)
	}

	var lock workflow.BundleLock
	if err := json.Unmarshal(data, &lock); err != nil {
		return fmt.Errorf("decode %s: %w", relativePath, err)
	}

	digest, ok := lock.Files[oldPath]
	if !ok {
		return fmt.Errorf("tracked path %q not found in %s", oldPath, relativePath)
	}

	delete(lock.Files, oldPath)
	lock.Files[newPath] = digest

	updated, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", relativePath, err)
	}
	updated = append(updated, '\n')

	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relativePath, err)
	}

	return nil
}

func (s *scenarioState) iAddTheLockFileSemanticDigestForTo(semanticPath, relativePath string) error {
	path := s.absPath(relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", relativePath, err)
	}

	var lock workflow.BundleLock
	if err := json.Unmarshal(data, &lock); err != nil {
		return fmt.Errorf("decode %s: %w", relativePath, err)
	}

	if lock.SemanticFiles == nil {
		lock.SemanticFiles = make(map[string]string)
	}

	digest, ok := lock.Files[semanticPath]
	if !ok {
		digest = "sha256:semantic-test-digest"
	}
	lock.SemanticFiles[semanticPath] = digest

	updated, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", relativePath, err)
	}
	updated = append(updated, '\n')

	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relativePath, err)
	}

	return nil
}

func (s *scenarioState) iRemoveTheLockFileSemanticDigestsFrom(relativePath string) error {
	path := s.absPath(relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", relativePath, err)
	}

	var lock workflow.BundleLock
	if err := json.Unmarshal(data, &lock); err != nil {
		return fmt.Errorf("decode %s: %w", relativePath, err)
	}

	lock.SemanticFiles = nil

	updated, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", relativePath, err)
	}
	updated = append(updated, '\n')

	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relativePath, err)
	}

	return nil
}

func (s *scenarioState) iReplaceTheLockFileWithValidJSONMissingRequiredFields(relativePath string) error {
	path := s.absPath(relativePath)
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("stat %s: %w", relativePath, err)
	}

	content := `{"name":"broken-lock"}` + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relativePath, err)
	}

	return nil
}

func (s *scenarioState) stderrContains(want string) error {
	if !strings.Contains(s.stderr, want) {
		return fmt.Errorf("expected stderr to contain %q, got %q", want, s.stderr)
	}

	return nil
}

func (s *scenarioState) initToolsAreAvailable() error {
	return s.installFakeVerificationTools("")
}

func (s *scenarioState) initToolsAreAvailableExcept(missing string) error {
	return s.installFakeVerificationTools(missing)
}

func (s *scenarioState) stdoutContains(want string) error {
	if !strings.Contains(s.stdout, want) {
		return fmt.Errorf("expected stdout to contain %q, got %q", want, s.stdout)
	}

	return nil
}

func (s *scenarioState) stdoutDoesNotContain(unwanted string) error {
	if strings.Contains(s.stdout, unwanted) {
		return fmt.Errorf("expected stdout not to contain %q, got %q", unwanted, s.stdout)
	}

	return nil
}

func (s *scenarioState) intendTraceSucceeds(name string) error {
	if err := s.runIntend("trace", name); err != nil {
		return err
	}

	if s.exitCode != 0 {
		return fmt.Errorf("expected intend trace %s to succeed, got exit code %d (stdout=%q, stderr=%q)", name, s.exitCode, s.stdout, s.stderr)
	}

	return nil
}

func (s *scenarioState) theLockVersionForIs(name string, expected int) error {
	data, err := os.ReadFile(s.absPath(filepath.Join(".intend", "locks", name+".json")))
	if err != nil {
		return fmt.Errorf("read lock file for %s: %w", name, err)
	}

	var lock lockState
	if err := json.Unmarshal(data, &lock); err != nil {
		return fmt.Errorf("decode lock file for %s: %w", name, err)
	}

	if lock.Version != expected {
		return fmt.Errorf("expected lock version %d, got %d", expected, lock.Version)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}

	ctx.Before(state.reset)
	ctx.After(state.cleanup)

	ctx.Step("^an empty working directory$", state.anEmptyWorkingDirectory)
	ctx.Step("^an initialized owned repository$", state.anInitializedOwnedRepository)
	ctx.Step("^an existing bundle `([^`]+)`$", state.anExistingBundle)
	ctx.Step("^a locked bundle `([^`]+)`$", state.aLockedBundle)
	ctx.Step("^I run `([^`]+)`$", state.iRun)
	ctx.Step("^it exits with code (\\d+)$", state.itExitsWithCode)
	ctx.Step("^the file `([^`]+)` exists$", state.theFileExists)
	ctx.Step("^the directory `([^`]+)` exists$", state.theDirectoryExists)
	ctx.Step("^I replace the contents of `([^`]+)`$", state.iReplaceTheContentsOf)
	ctx.Step("^stderr contains `([^`]+)`$", state.stderrContains)
	ctx.Step("^stdout contains `([^`]+)`$", state.stdoutContains)
	ctx.Step("^`intend trace ([^` ]+)` succeeds$", state.intendTraceSucceeds)
	ctx.Step("^the lock version for `([^`]+)` is (\\d+)$", state.theLockVersionForIs)
}

func TestFeatures(t *testing.T) {
	buildCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	binaryPath, cleanup, err := buildIntendBinary(buildCtx)
	if err != nil {
		t.Fatalf("build intend CLI: %v", err)
	}
	defer cleanup()

	intendBinaryPath = binaryPath

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join(intendRepoRoot, "features", "core-contract-lifecycle.feature")},
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func (s *scenarioState) runAndRequireSuccess(args ...string) error {
	if err := s.runIntend(args...); err != nil {
		return err
	}

	if s.exitCode != 0 {
		return fmt.Errorf("expected %q to succeed, got exit code %d (stdout=%q, stderr=%q)", strings.Join(args, " "), s.exitCode, s.stdout, s.stderr)
	}

	return nil
}

func (s *scenarioState) ensureInitToolsAvailable() error {
	if len(s.env) > 0 {
		return nil
	}

	return s.initToolsAreAvailable()
}

func (s *scenarioState) runIntend(args ...string) error {
	if err := s.ensureWorkDir(); err != nil {
		return err
	}

	if intendBinaryPath == "" {
		return errors.New("intend binary is not configured")
	}

	s.commandRan = true

	cmd := exec.Command(intendBinaryPath, args...)
	cmd.Dir = s.workDir
	if len(s.env) > 0 {
		cmd.Env = append(os.Environ(), s.env...)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	s.stdout = stdout.String()
	s.stderr = stderr.String()
	s.exitCode = 0

	if err == nil {
		return nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		s.exitCode = exitErr.ExitCode()
		return nil
	}

	return fmt.Errorf("run intend: %w", err)
}

func (s *scenarioState) ensureWorkDir() error {
	if s.workDir != "" {
		return nil
	}

	workDir, err := os.MkdirTemp("", "intend-contract-*")
	if err != nil {
		return fmt.Errorf("create temp work dir: %w", err)
	}

	s.workDir = workDir
	return nil
}

func (s *scenarioState) removeWorkDir() error {
	var errs []error
	for _, path := range s.extraPaths {
		if err := os.RemoveAll(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			errs = append(errs, fmt.Errorf("remove external path %s: %w", path, err))
		}
	}
	s.extraPaths = nil

	if s.workDir == "" {
		return errors.Join(errs...)
	}

	if err := os.RemoveAll(s.workDir); err != nil {
		errs = append(errs, err)
	}
	s.workDir = ""
	return errors.Join(errs...)
}

func (s *scenarioState) absPath(relativePath string) string {
	return filepath.Join(s.workDir, filepath.FromSlash(relativePath))
}

func buildIntendBinary(ctx context.Context) (string, func(), error) {
	tempDir, err := os.MkdirTemp("", "intend-binary-*")
	if err != nil {
		return "", nil, fmt.Errorf("create temp dir for intend binary: %w", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(tempDir)
	}

	binaryPath := filepath.Join(tempDir, "intend")
	build := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, "./cmd/intend")
	build.Dir = intendRepoRoot

	var stdout, stderr bytes.Buffer
	build.Stdout = &stdout
	build.Stderr = &stderr

	if err := build.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("build intend CLI: %w (stdout=%q, stderr=%q)", err, stdout.String(), stderr.String())
	}

	return binaryPath, cleanup, nil
}
