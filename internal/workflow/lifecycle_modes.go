package workflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

type bundlePaths struct {
	baseDir   string
	mode      string
	rootLabel string
	traceRel  string
	tracePath string
	lockRel   string
	lockPath  string
}

type BundleRef struct {
	Mode string
	Name string
}

func LockBundleWithMode(root, mode, name string) (BundleLock, error) {
	if err := validateBundleName(name); err != nil {
		return BundleLock{}, err
	}

	paths, err := resolveBundlePaths(root, mode, name)
	if err != nil {
		return BundleLock{}, err
	}

	if _, err := os.Stat(paths.lockPath); err == nil {
		return BundleLock{}, errors.New("bundle is already locked; use `intend amend` after an intentional change")
	} else if !errors.Is(err, os.ErrNotExist) {
		return BundleLock{}, fmt.Errorf("stat lock file: %w", err)
	}

	return writeModeLock(paths, name, 1)
}

func AmendBundleWithMode(root, mode, name string) (BundleLock, bool, bool, error) {
	if err := validateBundleName(name); err != nil {
		return BundleLock{}, false, false, err
	}

	paths, err := resolveBundlePaths(root, mode, name)
	if err != nil {
		return BundleLock{}, false, false, err
	}

	lock, err := readModeLock(paths, name)
	if err != nil {
		return BundleLock{}, false, false, err
	}

	trace, err := readModeTrace(paths, name)
	if err != nil {
		return BundleLock{}, false, false, err
	}

	current, err := digestsForModeBundle(paths, trace)
	if err != nil {
		return BundleLock{}, false, false, err
	}
	if err := validateContribIssueSnapshotIdentity(paths, trace); err != nil {
		return BundleLock{}, false, false, err
	}
	currentSemantic, err := semanticDigestsForModeBundle(paths, trace)
	if err != nil {
		return BundleLock{}, false, false, err
	}
	upgradedSemanticLock := contribSemanticLockUpgradePending(lock, paths, trace)

	if lockDigestsEqual(lock, current, currentSemantic) {
		return lock, false, false, nil
	}

	updated, err := writeModeLock(paths, name, lock.Version+1)
	if err != nil {
		return BundleLock{}, false, false, err
	}

	return updated, true, upgradedSemanticLock, nil
}

func TraceBundleWithMode(root, mode, name string) error {
	if err := validateBundleName(name); err != nil {
		return err
	}

	paths, err := resolveBundlePaths(root, mode, name)
	if err != nil {
		return err
	}

	trace, err := readModeTrace(paths, name)
	if err != nil {
		return err
	}

	if _, err := os.Stat(paths.lockPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w: %s bundle %s", ErrContractUnlocked, modeLabel(mode), name)
	} else if err != nil {
		return fmt.Errorf("stat lock file: %w", err)
	}

	lock, err := readModeLock(paths, name)
	if err != nil {
		return err
	}

	current, err := digestsForModeBundle(paths, trace)
	if err != nil {
		return err
	}
	if err := validateContribIssueSnapshotIdentity(paths, trace); err != nil {
		return err
	}
	currentSemantic, err := semanticDigestsForModeBundle(paths, trace)
	if err != nil {
		return err
	}

	if driftPath, ok := firstDriftedLockPath(lock, current, currentSemantic); ok {
		return fmt.Errorf("%w: %s", ErrContractDrift, driftPath)
	}

	return nil
}

func ListBundleRefs(root string) ([]BundleRef, error) {
	var refs []BundleRef

	for _, mode := range []string{"owned", "contrib"} {
		names, err := traceBundleNames(root, mode)
		if err != nil {
			return nil, err
		}

		for _, name := range names {
			refs = append(refs, BundleRef{Mode: mode, Name: name})
		}
	}

	return refs, nil
}

func modeLabel(mode string) string {
	if mode == "" {
		return "owned"
	}

	return mode
}

func resolveBundlePaths(root, mode, name string) (bundlePaths, error) {
	switch mode {
	case "", "owned":
		return bundlePaths{
			baseDir:   root,
			mode:      "owned",
			rootLabel: "workspace root",
			traceRel:  traceRelPath(name),
			tracePath: filepath.Join(root, filepath.FromSlash(traceRelPath(name))),
			lockRel:   lockRelPath(name),
			lockPath:  filepath.Join(root, filepath.FromSlash(lockRelPath(name))),
		}, nil
	case "contrib":
		gitDir, err := resolveGitDir(root)
		if err != nil {
			return bundlePaths{}, err
		}

		baseDir := filepath.Join(gitDir, "intend", "contrib", name)
		traceRel := filepath.ToSlash(filepath.Join("trace", name+".json"))
		lockRel := filepath.ToSlash(filepath.Join("locks", name+".json"))

		return bundlePaths{
			baseDir:   baseDir,
			mode:      "contrib",
			rootLabel: "bundle root",
			traceRel:  traceRel,
			tracePath: filepath.Join(baseDir, filepath.FromSlash(traceRel)),
			lockRel:   lockRel,
			lockPath:  filepath.Join(baseDir, filepath.FromSlash(lockRel)),
		}, nil
	default:
		return bundlePaths{}, fmt.Errorf("unsupported mode: %s", mode)
	}
}

func writeModeLock(paths bundlePaths, name string, version int) (BundleLock, error) {
	trace, err := readModeTrace(paths, name)
	if err != nil {
		return BundleLock{}, err
	}

	digests, err := digestsForModeBundle(paths, trace)
	if err != nil {
		return BundleLock{}, err
	}
	if err := validateContribIssueSnapshotIdentity(paths, trace); err != nil {
		return BundleLock{}, err
	}
	semanticDigests, err := semanticDigestsForModeBundle(paths, trace)
	if err != nil {
		return BundleLock{}, err
	}

	lock := BundleLock{
		Name:          name,
		Version:       version,
		Files:         digests,
		SemanticFiles: semanticDigests,
	}

	if err := os.MkdirAll(filepath.Dir(paths.lockPath), 0o755); err != nil {
		return BundleLock{}, fmt.Errorf("create lock directory: %w", err)
	}

	if err := writeJSONFile(paths.lockPath, lock); err != nil {
		return BundleLock{}, fmt.Errorf("write lock file: %w", err)
	}

	return lock, nil
}

func readModeTrace(paths bundlePaths, name string) (BundleTrace, error) {
	data, err := os.ReadFile(paths.tracePath)
	if err != nil {
		return BundleTrace{}, fmt.Errorf("read trace file for %s: %w", name, err)
	}

	var trace BundleTrace
	if err := json.Unmarshal(data, &trace); err != nil {
		return BundleTrace{}, fmt.Errorf("decode trace file for %s: %w", name, err)
	}

	if trace.SpecPath == "" || trace.FeaturePath == "" {
		return BundleTrace{}, errors.New("trace file is missing required paths")
	}
	if trace.Name != name {
		return BundleTrace{}, fmt.Errorf("trace file name mismatch: expected %s, got %s", name, trace.Name)
	}
	if trace.Mode != paths.mode {
		return BundleTrace{}, fmt.Errorf("trace file mode mismatch: expected %s, got %s", paths.mode, trace.Mode)
	}
	if paths.mode == "contrib" {
		if trace.IssueRef == "" || trace.IssuePath == "" {
			return BundleTrace{}, errors.New("contribution trace is missing required issue metadata")
		}
		if _, _, err := parseIssueRef(trace.IssueRef); err != nil {
			return BundleTrace{}, fmt.Errorf("contribution trace issueRef is invalid: %s", trace.IssueRef)
		}
	}
	if err := validateTracePath(paths, "specPath", trace.SpecPath); err != nil {
		return BundleTrace{}, err
	}
	if err := validateTracePath(paths, "featurePath", trace.FeaturePath); err != nil {
		return BundleTrace{}, err
	}
	if trace.IssuePath != "" {
		if err := validateTracePath(paths, "issuePath", trace.IssuePath); err != nil {
			return BundleTrace{}, err
		}
		if paths.mode == "contrib" && trace.IssuePath != "issue.json" {
			return BundleTrace{}, fmt.Errorf("contribution trace issuePath mismatch: expected issue.json, got %s", trace.IssuePath)
		}
	}

	return trace, nil
}

func validateContribIssueSnapshotIdentity(paths bundlePaths, trace BundleTrace) error {
	if paths.mode != "contrib" || trace.IssuePath == "" {
		return nil
	}

	expectedRepo, expectedNumber, err := parseIssueRef(trace.IssueRef)
	if err != nil {
		return fmt.Errorf("contribution trace issueRef is invalid: %s", trace.IssueRef)
	}

	data, err := os.ReadFile(filepath.Join(paths.baseDir, filepath.FromSlash(trace.IssuePath)))
	if err != nil {
		return fmt.Errorf("read contribution issue snapshot: %w", err)
	}

	var snapshot issueSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("decode contribution issue snapshot: %w", err)
	}
	if snapshot.Number <= 0 || strings.TrimSpace(snapshot.Title) == "" || strings.TrimSpace(snapshot.URL) == "" {
		return errors.New("contribution issue snapshot is missing required fields")
	}
	if snapshot.Number != expectedNumber {
		return fmt.Errorf("contribution issue snapshot number mismatch: expected %d, got %d", expectedNumber, snapshot.Number)
	}

	repoFromURL, issueNumberFromURL, status := issueRepoFromURL(snapshot.URL)
	switch status {
	case issueURLInvalidHost:
		return fmt.Errorf("contribution issue snapshot URL is not a GitHub issue URL: %s", snapshot.URL)
	case issueURLAlternateHostname:
		return fmt.Errorf("contribution issue snapshot URL uses unsupported GitHub hostname: %s", snapshot.URL)
	case issueURLInvalidShape:
		return fmt.Errorf("contribution issue snapshot URL is invalid: %s", snapshot.URL)
	case issueURLNonCanonical:
		return fmt.Errorf("contribution issue snapshot URL is non-canonical: %s", snapshot.URL)
	}
	if repoFromURL != expectedRepo {
		return fmt.Errorf("contribution issue snapshot repo mismatch: expected %s, got %s", expectedRepo, repoFromURL)
	}
	if issueNumberFromURL != expectedNumber {
		return fmt.Errorf("contribution issue snapshot URL number mismatch: expected %d, got %d", expectedNumber, issueNumberFromURL)
	}

	return nil
}

func validateTracePath(paths bundlePaths, field, relPath string) error {
	resolvedPath, err := validateBundleRelativePath(paths.baseDir, paths.rootLabel, "trace "+field, relPath)
	if err != nil {
		return err
	}
	if _, err := os.Lstat(resolvedPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("stat trace %s: %w", field, err)
	}

	evaluatedPath, err := filepath.EvalSymlinks(resolvedPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("resolve trace %s symlinks: %w", field, err)
	}

	relativeEvaluatedPath, err := filepath.Rel(filepath.Clean(paths.baseDir), evaluatedPath)
	if err != nil {
		return fmt.Errorf("resolve trace %s symlinks: %w", field, err)
	}
	if relativeEvaluatedPath == ".." || strings.HasPrefix(relativeEvaluatedPath, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("trace %s resolves through a symlink outside the %s: %s", field, paths.rootLabel, relPath)
	}

	return nil
}

func validateLockPath(paths bundlePaths, relPath string) error {
	resolvedPath, err := validateBundleRelativePath(paths.baseDir, paths.rootLabel, "lock file path", relPath)
	if err != nil {
		return err
	}
	if _, err := os.Lstat(resolvedPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("stat lock file path: %w", err)
	}

	evaluatedPath, err := filepath.EvalSymlinks(resolvedPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("resolve lock file path symlinks: %w", err)
	}

	relativeEvaluatedPath, err := filepath.Rel(filepath.Clean(paths.baseDir), evaluatedPath)
	if err != nil {
		return fmt.Errorf("resolve lock file path symlinks: %w", err)
	}
	if relativeEvaluatedPath == ".." || strings.HasPrefix(relativeEvaluatedPath, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("lock file path resolves through a symlink outside the %s: %s", paths.rootLabel, relPath)
	}

	return nil
}

func validateBundleRelativePath(baseDir, rootLabel, pathLabel, relPath string) (string, error) {
	if path.IsAbs(relPath) {
		return "", fmt.Errorf("%s must be relative to the %s: %s", pathLabel, rootLabel, relPath)
	}

	baseDir = filepath.Clean(baseDir)
	resolvedPath := filepath.Join(baseDir, filepath.FromSlash(relPath))
	relativePath, err := filepath.Rel(baseDir, resolvedPath)
	if err != nil {
		return "", fmt.Errorf("resolve %s: %w", pathLabel, err)
	}
	if relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("%s escapes the %s: %s", pathLabel, rootLabel, relPath)
	}

	return resolvedPath, nil
}

func readModeLock(paths bundlePaths, name string) (BundleLock, error) {
	data, err := os.ReadFile(paths.lockPath)
	if err != nil {
		return BundleLock{}, fmt.Errorf("read lock file for %s: %w", name, err)
	}

	var lock BundleLock
	if err := json.Unmarshal(data, &lock); err != nil {
		return BundleLock{}, fmt.Errorf("decode lock file for %s: %w", name, err)
	}
	if lock.Name == "" || lock.Version < 1 || len(lock.Files) == 0 {
		return BundleLock{}, errors.New("lock file is missing required fields")
	}
	if lock.Name != name {
		return BundleLock{}, fmt.Errorf("lock file name mismatch: expected %s, got %s", name, lock.Name)
	}
	for relPath := range lock.Files {
		if err := validateLockPath(paths, relPath); err != nil {
			return BundleLock{}, err
		}
	}
	semanticPaths := make([]string, 0, len(lock.SemanticFiles))
	for relPath := range lock.SemanticFiles {
		semanticPaths = append(semanticPaths, relPath)
	}
	slices.Sort(semanticPaths)
	for _, relPath := range semanticPaths {
		digest := lock.SemanticFiles[relPath]
		if digest == "" {
			return BundleLock{}, fmt.Errorf("lock file semantic digest is missing for %s", relPath)
		}
		if err := validateLockPath(paths, relPath); err != nil {
			return BundleLock{}, err
		}
		if _, ok := lock.Files[relPath]; !ok {
			return BundleLock{}, fmt.Errorf("lock file semantic path is not tracked in files: %s", relPath)
		}
		if paths.mode != "contrib" {
			return BundleLock{}, fmt.Errorf("lock file semantic path is unsupported for %s bundle: %s", paths.mode, relPath)
		}
		if relPath != "issue.json" {
			return BundleLock{}, fmt.Errorf("lock file semantic path mismatch: expected issue.json, got %s", relPath)
		}
	}

	return lock, nil
}

func digestsForModeBundle(paths bundlePaths, trace BundleTrace) (map[string]string, error) {
	files := []string{
		trace.SpecPath,
		trace.FeaturePath,
	}
	if trace.IssuePath != "" {
		files = append(files, trace.IssuePath)
	}
	files = append(files, paths.traceRel)

	digests := make(map[string]string, len(files))
	for _, relPath := range files {
		digest, err := fileDigest(filepath.Join(paths.baseDir, filepath.FromSlash(relPath)))
		if err != nil {
			return nil, fmt.Errorf("digest %s: %w", relPath, err)
		}

		digests[relPath] = digest
	}

	return digests, nil
}

func semanticDigestsForModeBundle(paths bundlePaths, trace BundleTrace) (map[string]string, error) {
	if paths.mode != "contrib" || trace.IssuePath == "" {
		return nil, nil
	}

	digest, err := semanticJSONDigest(filepath.Join(paths.baseDir, filepath.FromSlash(trace.IssuePath)))
	if err != nil {
		return nil, fmt.Errorf("semantic digest %s: %w", trace.IssuePath, err)
	}

	return map[string]string{trace.IssuePath: digest}, nil
}

func semanticJSONDigest(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return "", err
	}

	canonical, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return digestBytes(canonical), nil
}

func contribSemanticLockUpgradePending(lock BundleLock, paths bundlePaths, trace BundleTrace) bool {
	if paths.mode != "contrib" || trace.IssuePath == "" {
		return false
	}

	_, ok := lock.SemanticFiles[trace.IssuePath]
	return !ok
}

func lockDigestsEqual(lock BundleLock, currentFiles, currentSemanticFiles map[string]string) bool {
	_, drifted := firstDriftedLockPath(lock, currentFiles, currentSemanticFiles)
	return !drifted
}

func firstDriftedLockPath(lock BundleLock, currentFiles, currentSemanticFiles map[string]string) (string, bool) {
	if len(lock.Files) != len(currentFiles) {
		for relPath := range lock.Files {
			if _, ok := currentFiles[relPath]; !ok {
				return relPath, true
			}
		}
		for relPath := range currentFiles {
			if _, ok := lock.Files[relPath]; !ok {
				return relPath, true
			}
		}
	}

	for relPath, lockedDigest := range lock.Files {
		if lockedSemanticDigest, ok := lock.SemanticFiles[relPath]; ok {
			currentSemanticDigest, ok := currentSemanticFiles[relPath]
			if !ok || currentSemanticDigest != lockedSemanticDigest {
				return relPath, true
			}
			continue
		}

		currentDigest, ok := currentFiles[relPath]
		if !ok || currentDigest != lockedDigest {
			return relPath, true
		}
	}

	return "", false
}

func traceBundleNames(root, mode string) ([]string, error) {
	var names []string

	switch mode {
	case "owned":
		entries, err := os.ReadDir(filepath.Join(root, ".intend", "trace"))
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("read trace directory: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
				continue
			}

			names = append(names, strings.TrimSuffix(entry.Name(), ".json"))
		}
	case "contrib":
		gitDir, err := resolveGitDir(root)
		if err != nil {
			return nil, nil
		}

		contribRoot := filepath.Join(gitDir, "intend", "contrib")
		entries, err := os.ReadDir(contribRoot)
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("read contribution bundle directory: %w", err)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			name := entry.Name()
			tracePath := filepath.Join(contribRoot, name, "trace", name+".json")
			if _, err := os.Stat(tracePath); err == nil {
				names = append(names, name)
			} else if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("stat contribution trace file for %s: %w", name, err)
			}
		}
	default:
		return nil, fmt.Errorf("unsupported mode: %s", mode)
	}

	slices.Sort(names)
	return names, nil
}
