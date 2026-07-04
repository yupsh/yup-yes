// Command yup-yes is the CLI wrapper around github.com/gloo-foo/cmd-yes.
package main

import (
	"strings"

	clix "github.com/gloo-foo/cli"
	command "github.com/gloo-foo/cmd-yes"
	urf "github.com/urfave/cli/v3"
)

// version is the build version. It defaults to "dev" for local builds and is
// overridden at release time via the linker: -ldflags "-X main.version=<v>".
var version = "dev"

const (
	name      = "yes"
	flagCount = "count"
)

// synopsis is the multi-line --help usage block; urfave/cli indents it three
// spaces, so the lines stay flush-left.
const synopsis = `yes [OPTIONS] [STRING]...

repeatedly output a line with all specified STRING(s), or 'y'.`

// spec declares the yes wrapper. yes is a source command: it produces its
// repeated line directly, so build returns it as the whole pipeline (a nil
// filter).
var spec = clix.Spec{
	Name:     name,
	Summary:  "output a string repeatedly until killed",
	Synopsis: synopsis,
	Build:    build,
	Flags:    flags(),
}

// flags builds a fresh set of the wrapper's flags. It is a constructor rather
// than a package var so each parse gets independent flag values (urfave/cli
// records IsSet state on the flag itself, which would otherwise leak between
// invocations that share the pointers).
func flags() []urf.Flag {
	return []urf.Flag{
		&urf.IntFlag{
			Name:    flagCount,
			Aliases: []string{"n"},
			Usage:   "output COUNT lines instead of repeating forever",
		},
	}
}

// build maps the invocation to yes's pipeline: the STRING operands and flags
// produce the repeating source, with no filter.
func build(inv clix.Invocation) (clix.Source, clix.Command, error) {
	return command.Yes(options(inv.Args)...), nil, nil
}

// options folds the operands and parsed flags into yes's option values.
func options(c *urf.Command) []any {
	opts := []any{command.YesText(text(c))}
	if c.IsSet(flagCount) {
		opts = append(opts, command.YesCount(c.Int(flagCount)))
	}
	return opts
}

// text is the string yes repeats: the space-joined operands, or "y" when none
// are given.
func text(c *urf.Command) string {
	if c.NArg() == 0 {
		return "y"
	}
	return strings.Join(c.Args().Slice(), " ")
}

// runMain is an indirection seam so main's wiring is testable without spawning
// the process; a test swaps it and restores it.
var runMain = clix.Main

func main() { runMain(spec, version) }
