// Package get implements the "get" subcommand.
package get

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
	"sigs.k8s.io/yaml"

	_get "github.com/slewiskelly/ock/internal/pkg/get"
	"github.com/slewiskelly/ock/internal/pkg/report"
)

// Get implements the "get" subcommand.
type Get struct {
	def      string
	expr     string
	format   string
	validate bool
	schema   string
}

// Name returns the name of the subcommand.
func (*Get) Name() string {
	return "get"
}

// Synopsis returns a one-line summary of the subcommand.
func (*Get) Synopsis() string {
	return "displays the metadata of file(s) under a given path"
}

// Usage returns a longer explanation and/or usage example(s) of the subcommand.
func (*Get) Usage() string {
	return `ock get [flags] <path>
`
}

// SetFlags sets the flags specific to the subcommand.
func (g *Get) SetFlags(f *flag.FlagSet) {
	f.StringVar(&g.expr, "e", "", "expression to filter files")
	f.StringVar(&g.format, "f", "yaml", "display format (json | yaml)")
	f.StringVar(&g.schema, "s", ".schemacue", "location of the schema file to validate against")
}

// Execute executes the subcommand.
func (g *Get) Execute(ctx context.Context, fs *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "No path provided.\n\nUsage: ")
		fs.Usage()
		return subcommands.ExitUsageError
	}

	// TODO(slewiskelly): Validate flags.

	if err := g.execute(ctx, fs, args...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (g *Get) execute(ctx context.Context, fs *flag.FlagSet, args ...interface{}) error {
	var opts []_get.Option

	if g.expr != "" {
		opts = append(opts, _get.Expr(g.expr))
	}

	r, err := _get.Get(fs.Arg(0), opts...)
	if err != nil {
		return err
	}

	return display(r, g.format)

}

func display(r report.Report, f string) error {
	switch f {
	case "json":
		return displayJSON(r)
	case "yaml":
		return displayYAML(r)
	default:
		return errors.New("unknown output format")
	}
}

func displayJSON(r report.Report) error {
	var b []byte
	var err error

	if len(r) == 1 {
		b, err = json.MarshalIndent(r[0].Metadata, "", "  ")
	} else {
		b, err = json.MarshalIndent(r, "", "  ")
	}
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func displayYAML(r report.Report) error {
	var b []byte
	var err error

	if len(r) == 1 {
		b, err = yaml.Marshal(r[0].Metadata)
	} else {
		b, err = yaml.Marshal(r)
	}
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
