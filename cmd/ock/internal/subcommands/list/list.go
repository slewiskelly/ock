// Package list implements the "list" subcommand.
package list

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/google/subcommands"

	_list "github.com/slewiskelly/ock/internal/pkg/list"
)

// List implements the "list" subcommand.
type List struct {
	expr   string
	format string
}

// Name returns the name of the subcommand.
func (*List) Name() string {
	return "list"
}

// Synopsis returns a one-line summary of the subcommand.
func (*List) Synopsis() string {
	return "lists files under a given path"
}

// Usage returns a longer explanation and/or usage example(s) of the subcommand.
func (*List) Usage() string {
	return `ock list [flags] <path>
`
}

// SetFlags sets the flags specific to the subcommand.
func (l *List) SetFlags(f *flag.FlagSet) {
	f.StringVar(&l.expr, "e", "", "expression to filter files")
	f.StringVar(&l.format, "f", "summary", "display format (json | summary)")
}

// Execute executes the subcommand.
func (l *List) Execute(ctx context.Context, fs *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "No path provided.\n\nUsage: ")
		fs.Usage()
		return subcommands.ExitUsageError
	}

	// TODO(slewiskelly): Validate flags.

	if err := l.execute(ctx, fs, args...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (l *List) execute(ctx context.Context, fs *flag.FlagSet, args ...interface{}) error {
	var opts []_list.Option

	if l.expr != "" {
		opts = append(opts, _list.Expr(l.expr))
	}

	f, err := _list.List(fs.Arg(0), opts...)
	if err != nil {
		return err
	}

	return display(f, l.format)
}

func display(s []string, f string) error {
	switch f {
	case "json":
		return displayJSON(s)
	case "summary":
		return displaySummary(s)
	default:
		return errors.New("unknown formatting option")
	}
}

func displayJSON(s []string) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func displaySummary(s []string) error {
	w := new(strings.Builder)
	tw := tabwriter.NewWriter(w, 0, 8, 1, '\t', 0)

	fmt.Fprintln(tw, "File")
	fmt.Fprintln(tw, "----")

	for _, t := range s {
		fmt.Fprintf(tw, "%s\n", t)
	}

	tw.Flush()
	fmt.Println(w.String())

	return nil
}
