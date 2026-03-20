package commands

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"intend/internal/agents"
	"intend/internal/verify"
	"intend/internal/workflow"
)

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		if !printUsage(stderr) {
			return 1
		}
		return 2
	}

	root, err := os.Getwd()
	if err != nil {
		if !writef(stderr, "resolve working directory: %v\n", err) {
			return 1
		}
		return 1
	}

	switch args[0] {
	case "init":
		if len(args) != 1 {
			if !printUsage(stderr) {
				return 1
			}
			return 2
		}

		if err := verify.CheckRequiredTools(); err != nil {
			if !writef(stderr, "init: %v\n", err) {
				return 1
			}
			return 1
		}

		if err := workflow.Init(root); err != nil {
			if !writef(stderr, "init: %v\n", err) {
				return 1
			}
			return 1
		}

		if !writeln(stdout, "initialized intend workspace") {
			return 1
		}
		return 0
	case "new":
		flags := flag.NewFlagSet("new", flag.ContinueOnError)
		flags.SetOutput(stderr)

		mode := flags.String("mode", "owned", "bundle mode")
		fromGH := flags.String("from-gh", "", "GitHub issue reference")

		if err := flags.Parse(args[1:]); err != nil {
			return 2
		}

		if flags.NArg() != 1 {
			if !printUsage(stderr) {
				return 1
			}
			return 2
		}

		name := flags.Arg(0)
		switch *mode {
		case "owned":
			if err := workflow.CreateBundle(root, name); err != nil {
				if !writef(stderr, "new: %v\n", err) {
					return 1
				}
				return 1
			}

			if !writef(stdout, "created bundle %s\n", name) {
				return 1
			}
			return 0
		case "contrib":
			if *fromGH == "" {
				if !printUsage(stderr) {
					return 1
				}
				return 2
			}

			if err := workflow.CreateContribBundle(root, name, *fromGH); err != nil {
				if !writef(stderr, "new: %v\n", err) {
					return 1
				}
				return 1
			}

			if !writef(stdout, "created contribution bundle %s\n", name) {
				return 1
			}
			return 0
		default:
			if !writef(stderr, "new: unsupported mode: %s\n", *mode) {
				return 1
			}
			return 1
		}
	case "lock":
		mode, name, code := parseModeNameArgs("lock", args[1:], stderr)
		if code != 0 {
			return code
		}

		lock, err := workflow.LockBundleWithMode(root, mode, name)
		if err != nil {
			if !writef(stderr, "lock: %v\n", err) {
				return 1
			}
			return 1
		}

		if !writef(stdout, "locked %s at version %d\n", name, lock.Version) {
			return 1
		}
		return 0
	case "trace":
		mode, name, code := parseModeNameArgs("trace", args[1:], stderr)
		if code != 0 {
			return code
		}

		err := workflow.TraceBundleWithMode(root, mode, name)
		if err == nil {
			if !writef(stdout, "trace ok: %s\n", name) {
				return 1
			}
			return 0
		}

		if errors.Is(err, workflow.ErrContractDrift) {
			if !writef(stderr, "contract drift: %v\n", err) {
				return 1
			}
			return 1
		}

		if !writef(stderr, "trace: %v\n", err) {
			return 1
		}
		return 1
	case "amend":
		mode, name, code := parseModeNameArgs("amend", args[1:], stderr)
		if code != 0 {
			return code
		}

		lock, changed, upgradedSemanticLock, err := workflow.AmendBundleWithMode(root, mode, name)
		if err != nil {
			if !writef(stderr, "amend: %v\n", err) {
				return 1
			}
			return 1
		}

		if !changed {
			if !writef(stdout, "contract unchanged: %s remains at version %d\n", name, lock.Version) {
				return 1
			}
			return 0
		}

		if upgradedSemanticLock {
			if !writef(stdout, "amended %s to version %d and upgraded semantic lock metadata\n", name, lock.Version) {
				return 1
			}
			return 0
		}

		if !writef(stdout, "amended %s to version %d\n", name, lock.Version) {
			return 1
		}
		return 0
	case "verify":
		if len(args) != 1 {
			if !printVerifyUsage(stderr) {
				return 1
			}
			return 2
		}

		err := verify.Run(root, stdout)
		if err == nil {
			if !writeln(stdout, "verify ok") {
				return 1
			}
			return 0
		}

		if errors.Is(err, workflow.ErrContractDrift) {
			if !writef(stderr, "contract drift: %v\n", err) {
				return 1
			}
			return 1
		}

		if !writef(stderr, "verify: %v\n", err) {
			return 1
		}
		return 1
	case "agent":
		if len(args) != 3 || args[1] != "install" {
			if !printUsage(stderr) {
				return 1
			}
			return 2
		}

		if err := agents.Install(root, args[2]); err != nil {
			if !writef(stderr, "agent install: %v\n", err) {
				return 1
			}
			return 1
		}

		if !writef(stdout, "installed %s agent guidance\n", args[2]) {
			return 1
		}
		return 0
	default:
		if !printUsage(stderr) {
			return 1
		}
		return 2
	}
}

func printUsage(w io.Writer) bool {
	return writeln(w, "usage: intend <init|new|lock|trace|amend|verify|agent install> [name]")
}

func printVerifyUsage(w io.Writer) bool {
	return writeln(w, "usage: intend verify")
}

func parseModeNameArgs(command string, args []string, stderr io.Writer) (string, string, int) {
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	flags.SetOutput(stderr)

	mode := flags.String("mode", "owned", "bundle mode")
	if err := flags.Parse(args); err != nil {
		return "", "", 2
	}

	if flags.NArg() != 1 {
		if !printUsage(stderr) {
			return "", "", 1
		}
		return "", "", 2
	}

	return *mode, flags.Arg(0), 0
}

func writef(w io.Writer, format string, args ...any) bool {
	_, err := fmt.Fprintf(w, format, args...)
	return err == nil
}

func writeln(w io.Writer, args ...any) bool {
	_, err := fmt.Fprintln(w, args...)
	return err == nil
}
