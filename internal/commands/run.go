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
		printUsage(stderr)
		return 2
	}

	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stderr, "resolve working directory: %v\n", err)
		return 1
	}

	switch args[0] {
	case "init":
		if len(args) != 1 {
			printUsage(stderr)
			return 2
		}

		if err := verify.CheckRequiredTools(); err != nil {
			fmt.Fprintf(stderr, "init: %v\n", err)
			return 1
		}

		if err := workflow.Init(root); err != nil {
			fmt.Fprintf(stderr, "init: %v\n", err)
			return 1
		}

		fmt.Fprintln(stdout, "initialized intend workspace")
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
			printUsage(stderr)
			return 2
		}

		name := flags.Arg(0)
		switch *mode {
		case "owned":
			if err := workflow.CreateBundle(root, name); err != nil {
				fmt.Fprintf(stderr, "new: %v\n", err)
				return 1
			}

			fmt.Fprintf(stdout, "created bundle %s\n", name)
			return 0
		case "contrib":
			if *fromGH == "" {
				printUsage(stderr)
				return 2
			}

			if err := workflow.CreateContribBundle(root, name, *fromGH); err != nil {
				fmt.Fprintf(stderr, "new: %v\n", err)
				return 1
			}

			fmt.Fprintf(stdout, "created contribution bundle %s\n", name)
			return 0
		default:
			fmt.Fprintf(stderr, "new: unsupported mode: %s\n", *mode)
			return 1
		}
	case "lock":
		mode, name, code := parseModeNameArgs("lock", args[1:], stderr)
		if code != 0 {
			return code
		}

		lock, err := workflow.LockBundleWithMode(root, mode, name)
		if err != nil {
			fmt.Fprintf(stderr, "lock: %v\n", err)
			return 1
		}

		fmt.Fprintf(stdout, "locked %s at version %d\n", name, lock.Version)
		return 0
	case "trace":
		mode, name, code := parseModeNameArgs("trace", args[1:], stderr)
		if code != 0 {
			return code
		}

		err := workflow.TraceBundleWithMode(root, mode, name)
		if err == nil {
			fmt.Fprintf(stdout, "trace ok: %s\n", name)
			return 0
		}

		if errors.Is(err, workflow.ErrContractDrift) {
			fmt.Fprintf(stderr, "contract drift: %v\n", err)
			return 1
		}

		fmt.Fprintf(stderr, "trace: %v\n", err)
		return 1
	case "amend":
		mode, name, code := parseModeNameArgs("amend", args[1:], stderr)
		if code != 0 {
			return code
		}

		lock, changed, upgradedSemanticLock, err := workflow.AmendBundleWithMode(root, mode, name)
		if err != nil {
			fmt.Fprintf(stderr, "amend: %v\n", err)
			return 1
		}

		if !changed {
			fmt.Fprintf(stdout, "contract unchanged: %s remains at version %d\n", name, lock.Version)
			return 0
		}

		if upgradedSemanticLock {
			fmt.Fprintf(stdout, "amended %s to version %d and upgraded semantic lock metadata\n", name, lock.Version)
			return 0
		}

		fmt.Fprintf(stdout, "amended %s to version %d\n", name, lock.Version)
		return 0
	case "verify":
		if len(args) != 1 {
			printVerifyUsage(stderr)
			return 2
		}

		err := verify.Run(root, stdout)
		if err == nil {
			fmt.Fprintln(stdout, "verify ok")
			return 0
		}

		if errors.Is(err, workflow.ErrContractDrift) {
			fmt.Fprintf(stderr, "contract drift: %v\n", err)
			return 1
		}

		fmt.Fprintf(stderr, "verify: %v\n", err)
		return 1
	case "agent":
		if len(args) != 3 || args[1] != "install" {
			printUsage(stderr)
			return 2
		}

		if err := agents.Install(root, args[2]); err != nil {
			fmt.Fprintf(stderr, "agent install: %v\n", err)
			return 1
		}

		fmt.Fprintf(stdout, "installed %s agent guidance\n", args[2])
		return 0
	default:
		printUsage(stderr)
		return 2
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: intend <init|new|lock|trace|amend|verify|agent install> [name]")
}

func printVerifyUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: intend verify")
}

func parseModeNameArgs(command string, args []string, stderr io.Writer) (string, string, int) {
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	flags.SetOutput(stderr)

	mode := flags.String("mode", "owned", "bundle mode")
	if err := flags.Parse(args); err != nil {
		return "", "", 2
	}

	if flags.NArg() != 1 {
		printUsage(stderr)
		return "", "", 2
	}

	return *mode, flags.Arg(0), 0
}
