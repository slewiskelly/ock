// Package vet provides functionalty to validate a file's metadata.
package vet

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"cuelang.org/go/cue"

	"github.com/slewiskelly/ock/internal/pkg/get"
	"github.com/slewiskelly/ock/internal/pkg/report"
)

// Vet validates all files rooted at the given path, against the given schema.
//
// The returned report contains all files which failed validation, along with
// their corresponding error(s).
func Vet(path string, schema cue.Value, opts ...Option) (report.Report, error) {
	o := &options{}

	for _, opt := range opts {
		opt.apply(o)
	}

	if err := schema.Err(); err != nil {
		return nil, fmt.Errorf("invalid schema: %w", err)
	}

	schema = schema.LookupPath(cue.ParsePath("#Metadata"))

	var r report.Report

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.Type().IsRegular() || filepath.Ext(p) != ".md" {
			return nil
		}

		v, err := get.Get(p)
		if err != nil {
			r = append(r, &report.File{Name: p, Errors: []string{err.Error()}})
			return nil
		}

		// TODO(slewiskelly): Signals a lack of metadata, should be considered a
		// failure?
		if len(v) < 1 {
			return nil
		}

		if errs := validate(v[0].Metadata.Unify(schema)); len(errs) > 0 {
			r = append(r, &report.File{Name: p, Errors: errs})
		}

		return nil
	})

	return r, err
}

func validate(v cue.Value) []string {
	var errs []string

	i, err := v.Fields()
	if err != nil {
		return []string{fmt.Sprintf("failed to validate: %v", err)} // TODO(slewiskelly): Reconsider.
	}

	for i.Next() {
		x := i.Value()

		if err := x.Validate(cue.Concrete(true)); err != nil {
			if a := x.Attribute("error"); a.NumArgs() > 0 {
				errs = append(errs, fmt.Sprintf("%s: %s", x.Path(), a.Contents()))
				continue
			}

			errs = append(errs, err.Error())
		}
	}

	return errs
}
