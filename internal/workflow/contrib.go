package workflow

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type issueSnapshot struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	URL    string `json:"url"`
}

type issueURLStatus int

const (
	issueURLInvalidHost issueURLStatus = iota
	issueURLAlternateHostname
	issueURLInvalidShape
	issueURLNonCanonical
	issueURLValid
)

func CreateContribBundle(root, name, issueRef string) error {
	if err := validateBundleName(name); err != nil {
		return err
	}

	repo, number, err := parseIssueRef(issueRef)
	if err != nil {
		return err
	}

	gitDir, err := resolveGitDir(root)
	if err != nil {
		return err
	}

	issueJSON, err := importIssue(root, repo, number)
	if err != nil {
		return err
	}

	bundleDir := filepath.Join(gitDir, "intend", "contrib", name)
	specRel := filepath.ToSlash(filepath.Join("specs", name+".md"))
	featureRel := filepath.ToSlash(filepath.Join("features", name+".feature"))
	traceRel := filepath.ToSlash(filepath.Join("trace", name+".json"))
	issueRel := "issue.json"

	for _, rel := range []string{issueRel, specRel, featureRel, traceRel} {
		path := filepath.Join(bundleDir, filepath.FromSlash(rel))
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists", path)
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat %s: %w", path, err)
		}
	}

	if err := os.MkdirAll(filepath.Join(bundleDir, "specs"), 0o755); err != nil {
		return fmt.Errorf("create specs dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(bundleDir, "features"), 0o755); err != nil {
		return fmt.Errorf("create features dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(bundleDir, "trace"), 0o755); err != nil {
		return fmt.Errorf("create trace dir: %w", err)
	}

	if err := os.WriteFile(filepath.Join(bundleDir, issueRel), issueJSON, 0o644); err != nil {
		return fmt.Errorf("write issue snapshot: %w", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, filepath.FromSlash(specRel)), []byte(specTemplate(name)), 0o644); err != nil {
		return fmt.Errorf("write contrib spec: %w", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, filepath.FromSlash(featureRel)), []byte(featureTemplate(name)), 0o644); err != nil {
		return fmt.Errorf("write contrib feature: %w", err)
	}

	trace := BundleTrace{
		Name:        name,
		Mode:        "contrib",
		SpecPath:    specRel,
		FeaturePath: featureRel,
		IssueRef:    issueRef,
		IssuePath:   issueRel,
	}

	if err := writeJSONFile(filepath.Join(bundleDir, filepath.FromSlash(traceRel)), trace); err != nil {
		return fmt.Errorf("write contrib trace: %w", err)
	}

	return nil
}

func parseIssueRef(issueRef string) (string, int, error) {
	repo, issueNumber, ok := strings.Cut(issueRef, "#")
	if !ok || repo == "" || issueNumber == "" {
		return "", 0, fmt.Errorf("invalid GitHub issue reference %q", issueRef)
	}

	number, err := strconv.Atoi(issueNumber)
	if err != nil || number < 1 {
		return "", 0, fmt.Errorf("invalid GitHub issue reference %q", issueRef)
	}

	return repo, number, nil
}

func resolveGitDir(root string) (string, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return "", errors.New("required tool not found: git")
	}

	cmd := exec.Command(gitPath, "rev-parse", "--git-dir")
	cmd.Dir = root

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", errors.New("not inside a git repository")
	}

	gitDir := strings.TrimSpace(stdout.String())
	if gitDir == "" {
		return "", errors.New("not inside a git repository")
	}

	if filepath.IsAbs(gitDir) {
		return gitDir, nil
	}

	return filepath.Join(root, filepath.FromSlash(gitDir)), nil
}

func importIssue(root, repo string, number int) ([]byte, error) {
	ghPath, err := exec.LookPath("gh")
	if err != nil {
		return nil, errors.New("required tool not found: gh")
	}

	cmd := exec.Command(ghPath, "issue", "view", strconv.Itoa(number), "--repo", repo, "--json", "number,title,body,url")
	cmd.Dir = root

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("gh issue view failed: %w (stdout=%q, stderr=%q)", err, stdout.String(), stderr.String())
	}

	data := bytes.TrimSpace(stdout.Bytes())
	if len(data) == 0 {
		return nil, errors.New("gh issue view returned empty output")
	}

	var snapshot issueSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, errors.New("gh issue view returned invalid JSON")
	}
	if snapshot.Number <= 0 || strings.TrimSpace(snapshot.Title) == "" || strings.TrimSpace(snapshot.URL) == "" {
		return nil, errors.New("gh issue view returned incomplete issue data")
	}
	if snapshot.Number != number {
		return nil, fmt.Errorf("gh issue view returned issue #%d, expected #%d", snapshot.Number, number)
	}
	repoFromURL, issueNumberFromURL, status := issueRepoFromURL(snapshot.URL)
	if status == issueURLInvalidHost {
		return nil, fmt.Errorf("gh issue view returned non-GitHub issue URL: %s", snapshot.URL)
	}
	if status == issueURLAlternateHostname {
		parsed, _ := url.Parse(snapshot.URL)
		return nil, fmt.Errorf("gh issue view returned unsupported GitHub hostname: %s", parsed.Hostname())
	}
	if status == issueURLInvalidShape {
		return nil, fmt.Errorf("gh issue view returned invalid GitHub issue URL: %s", snapshot.URL)
	}
	if status == issueURLNonCanonical {
		return nil, fmt.Errorf("gh issue view returned non-canonical GitHub issue URL: %s", snapshot.URL)
	}
	if repoFromURL != repo {
		return nil, fmt.Errorf("gh issue view returned URL for repo %s, expected %s", repoFromURL, repo)
	}
	if issueNumberFromURL != number {
		return nil, fmt.Errorf("gh issue view returned URL issue #%d, expected #%d", issueNumberFromURL, number)
	}

	return append(data, '\n'), nil
}

func issueRepoFromURL(rawURL string) (string, int, issueURLStatus) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", 0, issueURLInvalidHost
	}
	hostname := strings.ToLower(parsed.Hostname())
	if hostname != "github.com" {
		if strings.HasSuffix(hostname, ".github.com") {
			return "", 0, issueURLAlternateHostname
		}
		return "", 0, issueURLInvalidHost
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", 0, issueURLNonCanonical
	}

	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) != 4 || parts[2] != "issues" {
		return "", 0, issueURLInvalidShape
	}
	issueNumber, err := strconv.Atoi(parts[3])
	if err != nil {
		return "", 0, issueURLInvalidShape
	}

	return parts[0] + "/" + parts[1], issueNumber, issueURLValid
}
