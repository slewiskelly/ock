// Package vet provides functionalty to validate a file's metadata.
package vet

import (
	"io/fs"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"

	"github.com/slewiskelly/ock/internal/pkg/get"
	"github.com/slewiskelly/ock/internal/pkg/report"
)

// Vet validates all files rooted at the given path, against the given schema.
//
// The returned report contains all files which failed validation, along with
// their corresponding error(s).
func Vet(path string, schema *cue.Value, opts ...Option) (report.Report, error) {
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

		if err := v[0].Metadata.Unify(*schema).Validate(); err != nil {
			r = append(r, &report.File{Name: p, Errors: errs(err)})
		}

		return nil
	})

	return r, err
}

func errs(err error) []string {
	var s []string

	for _, e := range errors.Errors(err) {
		s = append(s, e.Error())
	}

	return s
}
