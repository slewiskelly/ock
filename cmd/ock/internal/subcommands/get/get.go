// Package get implements the "get" subcommand.
package get

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/yaml"
	"github.com/google/subcommands"

	_get "github.com/slewiskelly/ock/internal/pkg/get"
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
	return "displays the metadata of a given file"
}

// Usage returns a longer explanation and/or usage example(s) of the subcommand.
func (*Get) Usage() string {
	return `ock get [flags] <path>
`
}

// SetFlags sets the flags specific to the subcommand.
func (g *Get) SetFlags(f *flag.FlagSet) {
	f.StringVar(&g.expr, "e", "", "expression to filter files")
	f.StringVar(&g.format, "f", "summary", "display format (json | yaml)")
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
		opts = append(opts, _get.WithExpr(g.expr))
	}

	v, err := _get.Get(fs.Arg(0), opts...)
	if err != nil {
		return err
	}

	return display(v, g.format)

}

func display(v *cue.Value, f string) error {
	switch f {
	case "json":
		return displayJSON(v)
	case "summary":
		return displayYAML(v)
	default:
		return errors.New("unknown output format")
	}
}

func displayJSON(v *cue.Value) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func displayYAML(v *cue.Value) error {
	b, err := yaml.Encode(*v)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
