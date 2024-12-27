// Package init provides functionality to initialize a default schema file.
package init

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// Init initializes a default schema with the given filename.
//
// By default, if a file with the given name already exists, it will not attempt
// to overwrite it, and an error will occur.
func Init(filename, def string, opts ...Option) error {
	o := &options{}

	for _, opt := range opts {
		opt.apply(o)
	}

	_, err := os.Stat(filename)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	if err == nil && !o.force {
		return fmt.Errorf("file %s already exists, specify --force to overwrite", filename)
	}

	b, err := f.ReadFile("schema.cue")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, b, 0o644)
}

//go:embed schema.cue
var f embed.FS
