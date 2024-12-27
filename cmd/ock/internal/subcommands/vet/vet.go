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

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/google/subcommands"

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
	f.StringVar(&v.format, "f", "summary", "display format (json | report | summary)")
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

	x := cuecontext.New().BuildInstance(i).LookupPath(cue.ParsePath("#Metadata"))

	r, err := _vet.Vet(fs.Arg(0), &x) // TODO(slewiskelly): Options.
	if err != nil {
		return err
	}

	if len(r.Files) < 1 {
		return nil
	}

	return display(r, v.format)
}

func display(r *_vet.Report, f string) error {
	switch f {
	case "json":
		return displayJSON(r)
	case "report":
		return displayReport(r)
	case "summary":
		return displaySummary(r)
	default:
		return errors.New("unknown output format")
	}
}

func displayJSON(r *_vet.Report) error {
	b, err := json.MarshalIndent(r.Files, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func displayReport(r *_vet.Report) error {
	w := new(strings.Builder)
	tw := tabwriter.NewWriter(w, 0, 8, 1, '\t', 0)

	for f, e := range r.Files {
		fmt.Fprintln(w, f)

		for _, err := range e {
			// TODO(slewiskelly): Display level, field, message.
			fmt.Fprintf(tw, "\t%v\n", err)
		}

		tw.Flush()
		fmt.Println(w.String())
		w.Reset()
	}

	return nil
}

func displaySummary(r *_vet.Report) error {
	w := new(strings.Builder)
	tw := tabwriter.NewWriter(w, 0, 8, 1, '\t', 0)

	fmt.Fprintln(tw, "File\tError")
	fmt.Fprintln(tw, "----\t--------")

	for f, e := range r.Files {
		// TODO(slewiskelly): Display file, number of warnings, number of errors
		fmt.Fprintf(tw, "%s\t%v\n", f, e[0]) // Displays only the first error encountered.
	}

	tw.Flush()
	fmt.Println(w.String())

	return nil
}
