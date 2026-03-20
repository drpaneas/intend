package workflow

import (
	"errors"
	"fmt"
	"os"
)

func DeleteBundleWithMode(root, mode, name string, force bool) error {
	if err := validateBundleName(name); err != nil {
		return err
	}

	paths, err := resolveBundlePaths(root, mode, name)
	if err != nil {
		return err
	}

	if !force {
		locked, err := bundleIsLocked(paths)
		if err != nil {
			return err
		}
		if locked {
			return fmt.Errorf("%s bundle %s is locked; re-run with --force to delete it", modeLabel(mode), name)
		}
	}

	trace, err := readModeTrace(paths, name)
	if err != nil {
		return err
	}

	if paths.mode == "contrib" {
		return deleteContribBundle(paths)
	}

	return deleteOwnedBundle(paths, trace)
}

func bundleIsLocked(paths bundlePaths) (bool, error) {
	if _, err := os.Stat(paths.lockPath); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, fmt.Errorf("stat lock file: %w", err)
	}
}

func deleteOwnedBundle(paths bundlePaths, trace BundleTrace) error {
	for _, relPath := range []string{
		trace.SpecPath,
		trace.FeaturePath,
		paths.traceRel,
		paths.lockRel,
	} {
		if err := deleteFileIfPresent(paths.baseDir, paths.rootLabel, relPath); err != nil {
			return err
		}
	}

	return nil
}

func deleteContribBundle(paths bundlePaths) error {
	if err := os.RemoveAll(paths.baseDir); err != nil {
		return fmt.Errorf("remove contribution bundle root: %w", err)
	}

	return nil
}

func deleteFileIfPresent(baseDir, rootLabel, relPath string) error {
	if relPath == "" {
		return nil
	}

	resolvedPath, err := validateBundleRelativePath(baseDir, rootLabel, "delete path", relPath)
	if err != nil {
		return err
	}

	info, err := os.Lstat(resolvedPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("stat delete path %s: %w", relPath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("delete path is a directory: %s", relPath)
	}

	if err := os.Remove(resolvedPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove %s: %w", relPath, err)
	}
	return nil
}
