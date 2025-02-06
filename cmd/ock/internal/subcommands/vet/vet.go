// Package vet implements the "vet" subcommand.
package vet

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/google/subcommands"

	"github.com/slewiskelly/ock/internal/pkg/report"
	_vet "github.com/slewiskelly/ock/internal/pkg/vet"
)

// Vet implements the "vet" subcommand.
type Vet struct {
	def    string
	format string
	schema string
}

// Name returns the name of the subcommand.
func (*Vet) Name() string {
	return "vet"
}

// Synopsis returns a one-line summary of the subcommand.
func (*Vet) Synopsis() string {
	return "validates file(s) under a given path, according to the schema"
}

// Usage returns a longer explanation and/or usage example(s) of the subcommand.
func (*Vet) Usage() string {
	return `ock vet [flags] <path>
`
}

// SetFlags sets the flags specific to the subcommand.
func (v *Vet) SetFlags(f *flag.FlagSet) {
	f.StringVar(&v.format, "f", "summary", "display format (json | summary)")
	f.StringVar(&v.schema, "schema", ".schema.cue", "location of the schema file to validate against")
}

// Execute executes the subcommand.
func (v *Vet) Execute(ctx context.Context, fs *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "No path provided.\n\nUsage: ")
		fs.Usage()
		return subcommands.ExitUsageError
	}

	// TODO(slewiskelly): Validate flags.

	if err := v.execute(ctx, fs, args...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (v *Vet) execute(ctx context.Context, fs *flag.FlagSet, args ...interface{}) error {
	i := load.Instances([]string{v.schema}, nil)[0]
	if err := i.Err; err != nil {
		return err
	}

	r, err := _vet.Vet(fs.Arg(0), cuecontext.New().BuildInstance(i)) // TODO(slewiskelly): Options.
	if err != nil {
		return err
	}

	if len(r) < 1 {
		return nil
	}

	return display(r, v.format)
}

func display(r report.Report, f string) error {
	switch f {
	case "json":
		return displayJSON(r)
	case "summary":
		return displaySummary(r)
	default:
		return errors.New("unknown output format")
	}
}

func displayJSON(r report.Report) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func displaySummary(r report.Report) error {
	w := new(strings.Builder)
	tw := tabwriter.NewWriter(w, 0, 8, 1, '\t', 0)

	var ec, wc int

	for _, x := range r {
		if len(x.Errors) > 0 || len(x.Warnings) > 0 {
			ec += len(x.Errors)
			wc += len(x.Warnings)

			fmt.Fprintln(tw, x.Name)

			for _, e := range x.Errors {
				fmt.Fprintf(tw, "\033[38;2;255;0;0mERROR\033[0m\t%v\n", e)
			}

			for _, e := range x.Warnings {
				fmt.Fprintf(tw, "\033[38;2;255;128;0mWARN\033[0m\t%v\n", e)
			}

			fmt.Fprintln(tw)
		}
	}

	tw.Flush()
	fmt.Println(w.String())
	fmt.Printf("\033[38;2;255;0;0m%d errors\033[0m, \033[38;2;255;128;0m%d warnings\033[0m\n", ec, wc)

	return nil
}
