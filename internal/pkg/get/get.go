// Package get provides functionality for retrieving document metadata.
package get

import (
	"bufio"
	"bytes"
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

		b, err := frontmatter(p)
		if err != nil || b == nil {
			return err
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

func frontmatter(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var opened, closed int

	for i := 1; s.Scan(); i++ {
		l := s.Bytes()

		if len(bytes.TrimSpace(l)) == 0 {
			continue
		}

		if string(l) != "---" {
			return nil, nil
		}

		opened = i
		break
	}

	if opened == 0 {
		return nil, nil
	}

	var b []byte

	for i := opened + 1; s.Scan(); i++ {
		l := s.Bytes()

		if string(l) == "---" {
			closed = i
			break
		}

		b = append(b, append(l, '\n')...)
	}

	if closed <= opened {
		return nil, fmt.Errorf("%s: frontmatter not closed", p)
	}

	return b, nil
}
