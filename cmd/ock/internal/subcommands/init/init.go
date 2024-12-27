// Package init implements the "init" subcommand.
package init

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"

	_init "github.com/slewiskelly/ock/internal/pkg/init"
)

// Init implements the "init" subcommand.
type Init struct {
	def    string
	force  bool
	schema string
}

// Name returns the name of the subcommand.
func (*Init) Name() string {
	return "init"
}

// Synopsis returns a one-line summary of the subcommand.
func (*Init) Synopsis() string {
	return "initializes a default schema, in the current working directory"
}

// Usage returns a longer explanation and/or usage example(s) of the subcommand.
func (*Init) Usage() string {
	return `ock init [flags]
`
}

// SetFlags sets the flags specific to the subcommand.
func (i *Init) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&i.force, "force", false, "force initialization if a schema file already exists")
	f.StringVar(&i.schema, "schema", ".schema.cue", "location of the schema file to validate against")
}

// Execute executes the subcommand.
func (i *Init) Execute(ctx context.Context, fs *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	// TODO(slewiskelly): Validate flags.

	if err := i.execute(ctx, fs, args...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (i *Init) execute(ctx context.Context, fs *flag.FlagSet, args ...interface{}) error {
	return _init.Init(i.schema, i.def)
}
