package easycmd

// Package easycmd provides an easy way to define a group of commands.
//
// Commands are grouped in a tree. Commands and common setup are
// functions registered in the tree. Commands are functions on the
// leaves, and setup are in the interior.
//
// A command is a typical cli command, a bunch of flags, flag.Parse(),
// and then the main body. A setup function is usually definition of
// common flags.
import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

type cmdInfo struct {
	desc string
	fn   func()
	// empty children is end cmd
	children map[string]*cmdInfo
}

var root = &cmdInfo{children: make(map[string]*cmdInfo)}

// Handle registers a cmd. Cmd is a space separated cmd chain,
// e.g. "db create".
func Handle(cmd string, fn func(), desc string) {
	if fn == nil {
		log.Fatal("empty fn")
	}
	chain := strings.Fields(strings.ToLower(cmd))
	if len(chain) == 0 {
		log.Fatalf("empty cmd:%q to handle", cmd)
	}

	err := addCmd(root, chain, fn, desc)
	if err != nil {
		log.Fatal(err)
	}
}

func addCmd(ci *cmdInfo, chain []string, fn func(), desc string) error {
	for _, c := range chain {
		if strings.HasPrefix(c, "-") {
			return fmt.Errorf("cmd:%q cannot begin with -", c)
		}

		cci, ok := ci.children[c]
		if !ok {
			cci = &cmdInfo{children: make(map[string]*cmdInfo)}
			ci.children[c] = cci
		}
		ci = cci
	}

	if ci.fn != nil {
		return fmt.Errorf("cmd:%q is a duplicate", chain)
	}
	ci.fn = fn
	ci.desc = desc
	return nil
}

// Main runs the command.
func Main() {
	ci, fns, chain := findCmd(root, os.Args)
	cmd := strings.Join(chain, " ")
	os.Args = append([]string{cmd}, os.Args[len(chain):]...)
	if !runCmd(ci, fns, cmd) {
		os.Exit(2)
	}
}

func findCmd(ci *cmdInfo, args []string) (_ *cmdInfo, fns []func(), chain []string) {
	chain = append(chain, args[0])
	if ci.fn != nil {
		fns = append(fns, ci.fn)
	}

	for _, c := range args[1:] {
		if strings.HasPrefix(c, "-") {
			break
		}
		cci, ok := ci.children[c]
		if !ok {
			break
		}
		ci = cci
		chain = append(chain, c)
		if ci.fn != nil {
			fns = append(fns, ci.fn)
		}
	}
	return ci, fns, chain
}

func runCmd(ci *cmdInfo, fns []func(), cmd string) bool {
	if len(ci.children) > 0 {
		printHelp(cmd, ci)
		return false
	}

	for _, fn := range fns {
		fn()
	}
	return true
}

func printHelp(cmd string, ci *cmdInfo) {
	stderr := flag.CommandLine.Output()
	fmt.Fprintf(stderr, "Usage of %s:\n", cmd)

	tw := tabwriter.NewWriter(stderr, 2, 4, 2, ' ', 0)
	for k, ci := range ci.children {
		fmt.Fprintf(tw, "\t%s\t%s\n", k, ci.desc)
	}
	tw.Flush()
}