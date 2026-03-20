package workflow

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ErrContractDrift    = errors.New("locked contract changed")
	ErrContractUnlocked = errors.New("contract is not locked")
	bundleNameRE        = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)
)

type BundleTrace struct {
	Name        string `json:"name"`
	Mode        string `json:"mode"`
	SpecPath    string `json:"specPath"`
	FeaturePath string `json:"featurePath"`
	IssueRef    string `json:"issueRef,omitempty"`
	IssuePath   string `json:"issuePath,omitempty"`
}

type BundleLock struct {
	Name          string            `json:"name"`
	Version       int               `json:"version"`
	Files         map[string]string `json:"files"`
	SemanticFiles map[string]string `json:"semanticFiles,omitempty"`
}

func Init(root string) error {
	for _, dir := range []string{
		"specs",
		"features",
		".intend/trace",
		".intend/locks",
	} {
		if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(dir)), 0o755); err != nil {
			return fmt.Errorf("create %s: %w", dir, err)
		}
	}

	return nil
}

func CreateBundle(root, name string) error {
	if err := validateBundleName(name); err != nil {
		return err
	}

	if !isInitialized(root) {
		return errors.New("repository is not initialized; run `intend init` first")
	}

	specRel := specRelPath(name)
	featureRel := featureRelPath(name)
	traceRel := traceRelPath(name)

	for _, rel := range []string{specRel, featureRel, traceRel} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err == nil {
			return fmt.Errorf("%s already exists", rel)
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat %s: %w", rel, err)
		}
	}

	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(specRel)), []byte(specTemplate(name)), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", specRel, err)
	}

	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(featureRel)), []byte(featureTemplate(name)), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", featureRel, err)
	}

	trace := BundleTrace{
		Name:        name,
		Mode:        "owned",
		SpecPath:    specRel,
		FeaturePath: featureRel,
	}

	if err := writeJSONFile(filepath.Join(root, filepath.FromSlash(traceRel)), trace); err != nil {
		return fmt.Errorf("write %s: %w", traceRel, err)
	}

	return nil
}

func LockBundle(root, name string) (BundleLock, error) {
	return LockBundleWithMode(root, "owned", name)
}

func AmendBundle(root, name string) (BundleLock, bool, bool, error) {
	return AmendBundleWithMode(root, "owned", name)
}

func TraceBundle(root, name string) error {
	return TraceBundleWithMode(root, "owned", name)
}

func TraceAllBundles(root string) error {
	refs, err := ListBundleRefs(root)
	if err != nil {
		return err
	}

	for _, ref := range refs {
		if err := TraceBundleWithMode(root, ref.Mode, ref.Name); err != nil {
			return err
		}
	}

	return nil
}

func writeLock(root, name string, version int) (BundleLock, error) {
	trace, err := readTrace(root, name)
	if err != nil {
		return BundleLock{}, err
	}

	digests, err := digestsForBundle(root, name, trace)
	if err != nil {
		return BundleLock{}, err
	}

	lock := BundleLock{
		Name:    name,
		Version: version,
		Files:   digests,
	}

	if err := writeJSONFile(filepath.Join(root, filepath.FromSlash(lockRelPath(name))), lock); err != nil {
		return BundleLock{}, fmt.Errorf("write lock file: %w", err)
	}

	return lock, nil
}

func readTrace(root, name string) (BundleTrace, error) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(traceRelPath(name))))
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

	return trace, nil
}

func readLock(root, name string) (BundleLock, error) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(lockRelPath(name))))
	if err != nil {
		return BundleLock{}, fmt.Errorf("read lock file for %s: %w", name, err)
	}

	var lock BundleLock
	if err := json.Unmarshal(data, &lock); err != nil {
		return BundleLock{}, fmt.Errorf("decode lock file for %s: %w", name, err)
	}

	return lock, nil
}

func digestsForBundle(root, name string, trace BundleTrace) (map[string]string, error) {
	files := []string{
		trace.SpecPath,
		trace.FeaturePath,
		traceRelPath(name),
	}

	digests := make(map[string]string, len(files))
	for _, relPath := range files {
		digest, err := fileDigest(filepath.Join(root, filepath.FromSlash(relPath)))
		if err != nil {
			return nil, fmt.Errorf("digest %s: %w", relPath, err)
		}

		digests[relPath] = digest
	}

	return digests, nil
}

func fileDigest(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return digestBytes(data), nil
}

func digestBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func writeJSONFile(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func validateBundleName(name string) error {
	if !bundleNameRE.MatchString(name) {
		return fmt.Errorf("invalid bundle name %q", name)
	}

	return nil
}

func isInitialized(root string) bool {
	info, err := os.Stat(filepath.Join(root, ".intend"))
	return err == nil && info.IsDir()
}

func specRelPath(name string) string {
	return filepath.ToSlash(filepath.Join("specs", name+".md"))
}

func featureRelPath(name string) string {
	return filepath.ToSlash(filepath.Join("features", name+".feature"))
}

func traceRelPath(name string) string {
	return filepath.ToSlash(filepath.Join(".intend", "trace", name+".json"))
}

func lockRelPath(name string) string {
	return filepath.ToSlash(filepath.Join(".intend", "locks", name+".json"))
}

func specTemplate(name string) string {
	return fmt.Sprintf("# %s\n\nDescribe the intent for `%s` here.\n", name, name)
}

func featureTemplate(name string) string {
	return fmt.Sprintf("Feature: %s\n  Scenario: Review and refine this contract\n    Given the intent is not finished\n    Then edit this feature before implementation\n", name)
}
