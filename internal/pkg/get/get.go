// Package get provides functionality for retrieving document metadata.
package get

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/yaml"

	"github.com/slewiskelly/ock/internal/pkg/report"
)

// Get retrieves metadata from (markdown) files rooted at the given path.
func Get(path string, opts ...Option) (report.Report, error) {
	o := &options{}

	for _, opt := range opts {
		opt.apply(o)
	}

	var r report.Report

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.Type().IsRegular() || filepath.Ext(p) != ".md" {
			return nil
		}

		var b []byte

		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		s := bufio.NewScanner(f)

		var closed bool

		for i := 0; s.Scan(); i++ {
			l := s.Bytes()

			if i == 0 {
				if string(l) != "---" {
					return nil
				}
				continue
			}

			if i > 0 && string(l) == "---" {
				closed = true
				break
			}

			b = append(b, append(l, '\n')...)
		}

		if !closed {
			return fmt.Errorf("%s: frontmatter not closed", p)
		}

		y, err := yaml.Extract(path, b)
		if err != nil {
			return err
		}

		v := cuecontext.New().BuildFile(y)
		if err := v.Err(); err != nil {
			return err
		}

		r = append(r, &report.File{
			Name:     p,
			Metadata: v,
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	if o.expr == "" {
		return r, nil
	}

	// TODO(slewiskelly): Support filtering via expression.

	return r, nil
}
