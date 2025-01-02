// Package list provides functionality to list files.
package list

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/slewiskelly/ock/internal/pkg/get"
)

// List lists (markdown) files rooted at the given path.
func List(path string, opts ...Option) ([]string, error) {
	o := &options{}

	for _, opt := range opts {
		opt.apply(o)
	}

	var files []string

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.Type().IsRegular() || filepath.Ext(p) != ".md" {
			return nil
		}

		if o.expr == "" {
			files = append(files, p)
			return nil
		}

		v, err := get.Get(p, get.Expr(o.expr))
		if err != nil {
			return fmt.Errorf("failed to extract metadata from %q: %w", p, err)
		}

		if v != nil {
			files = append(files, p)
		}

		return nil
	})

	return files, err
}
